package jvm

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"math"
	"os"
	"path/filepath"
)

type BuildPackage struct {
	Package         *ast.Package
	classes         map[string]*cg.ClassHighLevel
	mainClass       *cg.ClassHighLevel
	BuildExpression BuildExpression
}

func (this *BuildPackage) Make(p *ast.Package) {
	this.Package = p
	mainClass := &cg.ClassHighLevel{}
	this.mainClass = mainClass
	mainClass.AccessFlags |= cg.AccClassPublic
	mainClass.AccessFlags |= cg.AccClassFinal
	mainClass.AccessFlags |= cg.AccClassSynthetic
	mainClass.SuperClass = ast.JavaRootClass
	mainClass.Name = p.Name + "/main"
	if p.Block.Functions != nil {
		for _, v := range p.Block.Functions {
			mainClass.InsertSourceFile(v.Pos.Filename)
			break
		}
	}
	mainClass.Fields = make(map[string]*cg.FieldHighLevel)
	this.mkClassDefaultConstruction(this.mainClass)
	this.BuildExpression.BuildPackage = this
	this.classes = make(map[string]*cg.ClassHighLevel)
	this.mkGlobalConstants()
	this.mkGlobalTypeAlias()
	this.mkGlobalVariables()
	this.mkGlobalFunctions()
	this.mkInitFunctions()
	for _, v := range p.Block.Classes {
		this.putClass(this.buildClass(v))
	}
	for _, v := range p.Block.Enums {
		this.putClass(this.mkEnum(v))
	}
	err := this.DumpClass()
	if err != nil {
		panic(fmt.Sprintf("dump to file failed,err:%v\n", err))
	}
}

func (this *BuildPackage) newClassName(prefix string) (autoName string) {
	for i := 0; i < math.MaxInt16; i++ {
		if i == 0 {
			//use prefix only
			autoName = prefix
		} else {
			autoName = fmt.Sprintf("%s$%d", prefix, i)
		}
		if _, exists := this.Package.Block.NameExists(autoName); exists {
			continue
		}
		autoName = this.Package.Name + "/" + autoName
		if this.classes != nil && this.classes[autoName] != nil {
			continue
		} else {
			return autoName
		}
	}
	panic("new class name overflow") // impossible
}

func (this *BuildPackage) putClass(class *cg.ClassHighLevel) {
	if class.Name == "" {
		panic("missing name")
	}
	name := class.Name
	if name == this.mainClass.Name {
		panic("cannot have main class`s name")
	}
	if this.classes == nil {
		this.classes = make(map[string]*cg.ClassHighLevel)
	}
	if _, ok := this.classes[name]; ok {
		panic(fmt.Sprintf("name:'%s' already been token", name))
	}
	this.classes[name] = class
}

func (this *BuildPackage) mkEnum(e *ast.Enum) *cg.ClassHighLevel {
	class := &cg.ClassHighLevel{}
	class.Name = e.Name
	class.InsertSourceFile(e.Pos.Filename)
	class.AccessFlags = e.AccessFlags
	class.SuperClass = ast.JavaRootClass
	class.Fields = make(map[string]*cg.FieldHighLevel)
	class.Class.AttributeLucyEnum = &cg.AttributeLucyEnum{}
	if e.Comment != "" {
		class.Class.AttributeLucyComment = &cg.AttributeLucyComment{
			Comment: e.Comment,
		}
	}
	for _, v := range e.Enums {
		field := &cg.FieldHighLevel{}
		if e.AccessFlags&cg.AccClassPublic != 0 {
			field.AccessFlags |= cg.AccFieldPublic
		} else {
			field.AccessFlags |= cg.AccFieldPrivate
		}
		if v.Comment != "" {
			field.AttributeLucyComment = &cg.AttributeLucyComment{
				Comment: v.Comment,
			}
		}
		field.Name = v.Name
		field.Descriptor = "I"
		field.AttributeConstantValue = &cg.AttributeConstantValue{}
		field.AttributeConstantValue.Index = class.Class.InsertIntConst(v.Value)
		class.Fields[v.Name] = field
	}
	return class
}

func (this *BuildPackage) mkGlobalConstants() {
	for k, v := range this.Package.Block.Constants {
		f := &cg.FieldHighLevel{}
		f.AccessFlags |= cg.AccFieldStatic
		f.AccessFlags |= cg.AccFieldFinal
		if v.AccessFlags&cg.AccFieldPublic != 0 {
			f.AccessFlags |= cg.AccFieldPublic
		}
		f.Name = v.Name
		f.AttributeConstantValue = &cg.AttributeConstantValue{}
		f.AttributeConstantValue.Index = this.insertDefaultValue(this.mainClass, v.Type, v.Value)
		f.AttributeLucyConst = &cg.AttributeLucyConst{}
		if v.Comment != "" {
			f.AttributeLucyComment = &cg.AttributeLucyComment{
				Comment: v.Comment,
			}
		}
		f.Descriptor = Descriptor.typeDescriptor(v.Type)
		this.mainClass.Fields[k] = f
	}
}
func (this *BuildPackage) mkGlobalTypeAlias() {
	for name, v := range this.Package.Block.TypeAliases {
		t := &cg.AttributeLucyTypeAlias{}
		t.Alias = LucyTypeAliasParser.Encode(name, v)
		if v.Alias != nil {
			t.Comment = v.Alias.Comment
		}
		this.mainClass.Class.TypeAlias = append(this.mainClass.Class.TypeAlias, t)
	}
}

func (this *BuildPackage) mkGlobalVariables() {
	for k, v := range this.Package.Block.Variables {
		f := &cg.FieldHighLevel{}
		f.AccessFlags |= cg.AccFieldStatic
		f.Descriptor = Descriptor.typeDescriptor(v.Type)
		if v.AccessFlags&cg.AccFieldPublic != 0 {
			f.AccessFlags |= cg.AccFieldPublic
		}
		f.AccessFlags |= cg.AccFieldVolatile
		if LucyFieldSignatureParser.Need(v.Type) {
			f.AttributeLucyFieldDescriptor = &cg.AttributeLucyFieldDescriptor{}
			f.AttributeLucyFieldDescriptor.Descriptor = LucyFieldSignatureParser.Encode(v.Type)
			if v.Type.Type == ast.VariableTypeFunction {
				f.AttributeLucyFieldDescriptor.MethodAccessFlag |=
					cg.AccMethodVarargs
			}
		}
		if v.Comment != "" {
			f.AttributeLucyComment = &cg.AttributeLucyComment{
				Comment: v.Comment,
			}
		}
		f.Name = v.Name
		this.mainClass.Fields[k] = f
	}
}

func (this *BuildPackage) mkInitFunctions() {
	if len(this.Package.InitFunctions) == 0 {
		needTrigger := false
		for _, v := range this.Package.LoadedPackages {
			if v.TriggerPackageInitMethodName != "" {
				needTrigger = true
				break
			}
		}
		if needTrigger == false {
			return
		}
	}
	blockMethods := []*cg.MethodHighLevel{}
	for _, v := range this.Package.InitFunctions {
		method := &cg.MethodHighLevel{}
		blockMethods = append(blockMethods, method)
		method.AccessFlags |= cg.AccMethodStatic
		method.AccessFlags |= cg.AccMethodFinal
		method.AccessFlags |= cg.AccMethodPrivate
		method.Name = this.mainClass.NewMethodName("block")
		method.Class = this.mainClass
		method.Descriptor = "()V"
		method.Code = &cg.AttributeCode{}
		this.buildFunction(this.mainClass, nil, method, v)
		this.mainClass.AppendMethod(method)
	}
	method := &cg.MethodHighLevel{}
	method.AccessFlags |= cg.AccMethodStatic
	method.Name = "<clinit>"
	method.Descriptor = "()V"
	codes := make([]byte, 65536)
	codeLength := int(0)
	method.Code = &cg.AttributeCode{}
	for _, v := range this.Package.LoadedPackages {
		if v.TriggerPackageInitMethodName == "" {
			continue
		}
		codes[codeLength] = cg.OP_invokestatic
		this.mainClass.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      v.Name + "/main", // main class
			Method:     v.TriggerPackageInitMethodName,
			Descriptor: "()V",
		}, codes[codeLength+1:codeLength+3])
		codeLength += 3
	}
	for _, v := range blockMethods {
		codes[codeLength] = cg.OP_invokestatic
		this.mainClass.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
			Class:      this.mainClass.Name,
			Method:     v.Name,
			Descriptor: "()V",
		}, codes[codeLength+1:codeLength+3])
		codeLength += 3
	}
	codes[codeLength] = cg.OP_return
	codeLength++
	codes = codes[0:codeLength]
	method.Code.Codes = codes
	method.Code.CodeLength = codeLength
	this.mainClass.AppendMethod(method)

	// trigger init
	trigger := &cg.MethodHighLevel{}
	trigger.Name = this.mainClass.NewMethodName("triggerPackageInit")
	trigger.AccessFlags |= cg.AccMethodPublic
	trigger.AccessFlags |= cg.AccMethodBridge
	trigger.AccessFlags |= cg.AccMethodStatic
	trigger.AccessFlags |= cg.AccMethodSynthetic
	trigger.Descriptor = "()V"
	trigger.Code = &cg.AttributeCode{}
	trigger.Code.Codes = make([]byte, 1)
	trigger.Code.Codes[0] = cg.OP_return
	trigger.Code.CodeLength = 1
	trigger.AttributeLucyTriggerPackageInitMethod = &cg.AttributeLucyTriggerPackageInitMethod{}
	this.mainClass.AppendMethod(trigger)
	this.mainClass.TriggerPackageInitMethod = trigger
}

func (this *BuildPackage) insertDefaultValue(c *cg.ClassHighLevel, t *ast.Type, v interface{}) (index uint16) {
	switch t.Type {
	case ast.VariableTypeBool:
		if v.(bool) {
			index = c.Class.InsertIntConst(1)
		} else {
			index = c.Class.InsertIntConst(0)
		}
	case ast.VariableTypeByte:
		index = c.Class.InsertIntConst(int32(v.(int64)))
	case ast.VariableTypeShort:
		index = c.Class.InsertIntConst(int32(v.(int64)))
	case ast.VariableTypeInt:
		index = c.Class.InsertIntConst(int32(v.(int64)))
	case ast.VariableTypeLong:
		index = c.Class.InsertLongConst(v.(int64))
	case ast.VariableTypeFloat:
		index = c.Class.InsertFloatConst(v.(float32))
	case ast.VariableTypeDouble:
		index = c.Class.InsertDoubleConst(v.(float64))
	case ast.VariableTypeString:
		index = c.Class.InsertStringConst(v.(string))
	default:
		panic("no match")
	}
	return
}

func (this *BuildPackage) buildClass(astClass *ast.Class) *cg.ClassHighLevel {
	class := &cg.ClassHighLevel{}
	class.Name = astClass.Name
	class.InsertSourceFile(astClass.Pos.Filename)
	class.AccessFlags = astClass.AccessFlags
	if astClass.SuperClass != nil {
		class.SuperClass = astClass.SuperClass.Name
	} else {
		class.SuperClass = astClass.SuperClassName.Name
	}
	if astClass.Comment != "" {
		class.Class.AttributeLucyComment = &cg.AttributeLucyComment{
			Comment: astClass.Comment,
		}
	}
	if len(astClass.Block.Constants) > 0 {
		attr := &cg.AttributeLucyClassConst{}
		for _, v := range astClass.Block.Constants {
			c := &cg.LucyClassConst{}
			c.Name = v.Name
			c.Comment = v.Comment
			c.Descriptor = Descriptor.typeDescriptor(v.Type)
			c.ValueIndex = this.insertDefaultValue(class, v.Type, v.Value)
			attr.Constants = append(attr.Constants, c)
		}
		class.Class.AttributeLucyClassConst = attr
	}
	class.Fields = make(map[string]*cg.FieldHighLevel)
	class.Methods = make(map[string][]*cg.MethodHighLevel)
	for _, v := range astClass.Interfaces {
		class.Interfaces = append(class.Interfaces, v.Name)
	}
	for _, v := range astClass.Fields {
		f := &cg.FieldHighLevel{}
		f.Name = v.Name
		f.AccessFlags = v.AccessFlags
		if v.IsStatic() && v.DefaultValue != nil {
			f.AttributeConstantValue = &cg.AttributeConstantValue{}
			f.AttributeConstantValue.Index = this.insertDefaultValue(class, v.Type,
				v.DefaultValue)
		}
		f.Descriptor = Descriptor.typeDescriptor(v.Type)
		if LucyFieldSignatureParser.Need(v.Type) {
			t := &cg.AttributeLucyFieldDescriptor{}
			t.Descriptor = LucyFieldSignatureParser.Encode(v.Type)
			f.AttributeLucyFieldDescriptor = t
		}
		if v.Comment != "" {
			f.AttributeLucyComment = &cg.AttributeLucyComment{
				Comment: v.Comment,
			}
		}
		class.Fields[v.Name] = f
	}
	for name, v := range astClass.Methods {
		vv := v[0]
		method := &cg.MethodHighLevel{}
		method.Name = name
		method.AccessFlags = vv.Function.AccessFlags
		if vv.Function.Type.VArgs != nil {
			method.AccessFlags |= cg.AccMethodVarargs
		}
		if vv.IsCompilerAuto {
			method.AccessFlags |= cg.AccMethodSynthetic
		}
		method.Class = class
		method.Descriptor = Descriptor.methodDescriptor(&vv.Function.Type)
		method.IsConstruction = name == specialMethodInit
		if vv.IsAbstract() == false {
			method.Code = &cg.AttributeCode{}
			this.buildFunction(class, astClass, method, vv.Function)
		}
		class.AppendMethod(method)
	}
	return class
}

func (this *BuildPackage) mkGlobalFunctions() {
	ms := make(map[string]*cg.MethodHighLevel)
	for k, f := range this.Package.Block.Functions { // first round
		if f.TemplateFunction != nil {
			this.mainClass.TemplateFunctions = append(this.mainClass.TemplateFunctions,
				&cg.AttributeTemplateFunction{
					Name:        f.Name,
					Filename:    f.Pos.Filename,
					StartLine:   uint16(f.Pos.Line),
					StartColumn: uint16(f.Pos.Column),
					Code:        string(f.SourceCode),
				})
			continue
		}
		if f.IsBuildIn { //
			continue
		}
		class := this.mainClass
		method := &cg.MethodHighLevel{}
		method.Class = class
		method.Name = f.Name
		if f.Name == ast.MainFunctionName {
			method.Descriptor = "([Ljava/lang/String;)V"
		} else {
			method.Descriptor = Descriptor.methodDescriptor(&f.Type)
		}
		method.AccessFlags = 0
		method.AccessFlags |= cg.AccMethodStatic
		if f.AccessFlags&cg.AccMethodPublic != 0 || f.Name == ast.MainFunctionName {
			method.AccessFlags |= cg.AccMethodPublic
		}
		if f.Comment != "" {
			method.AttributeLucyComment = &cg.AttributeLucyComment{
				Comment: f.Comment,
			}
		}
		if f.Type.VArgs != nil {
			method.AccessFlags |= cg.AccMethodVarargs
		}
		ms[k] = method
		f.Entrance = method
		method.Code = &cg.AttributeCode{}
		this.mainClass.AppendMethod(method)
	}
	for k, f := range this.Package.Block.Functions {
		if f.IsBuildIn || f.TemplateFunction != nil { //
			continue
		}
		this.buildFunction(ms[k].Class, nil, ms[k], f)
	}
}

func (this *BuildPackage) DumpClass() error {
	//dump main class
	f, err := os.OpenFile("main.class", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if err := this.mainClass.ToLow().OutPut(f); err != nil {
		f.Close()
		return err
	}
	f.Close()
	for _, c := range this.classes {
		f, err = os.OpenFile(filepath.Base(c.Name)+".class", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		if err = c.ToLow().OutPut(f); err != nil {
			f.Close()
			return err
		} else {
			f.Close()
		}
	}
	return nil
}

/*
	make_node_objects a default construction
*/
func (this *BuildPackage) mkClassDefaultConstruction(class *cg.ClassHighLevel) {
	method := &cg.MethodHighLevel{}
	method.Name = specialMethodInit
	method.Descriptor = "()V"
	method.AccessFlags |= cg.AccMethodPublic
	method.Code = &cg.AttributeCode{}
	method.Code.Codes = make([]byte, 5)
	method.Code.CodeLength = 5
	method.Code.MaxLocals = 1
	method.Code.Codes[0] = cg.OP_aload_0
	method.Code.Codes[1] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.ConstantInfoMethodrefHighLevel{
		Class:      class.SuperClass,
		Method:     specialMethodInit,
		Descriptor: "()V",
	}, method.Code.Codes[2:4])
	method.Code.MaxStack = 1
	method.Code.Codes[4] = cg.OP_return
	class.AppendMethod(method)
}

func (this *BuildPackage) storeGlobalVariable(class *cg.ClassHighLevel, code *cg.AttributeCode,
	v *ast.Variable) {
	code.Codes[code.CodeLength] = cg.OP_putstatic
	if v.JvmDescriptor == "" {
		v.JvmDescriptor = Descriptor.typeDescriptor(v.Type)
	}
	class.InsertFieldRefConst(cg.ConstantInfoFieldrefHighLevel{
		Class:      this.mainClass.Name,
		Field:      v.Name,
		Descriptor: v.JvmDescriptor,
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
}

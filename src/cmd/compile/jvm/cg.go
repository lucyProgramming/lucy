package jvm

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type BuildPackage struct {
	Package         *ast.Package
	classes         map[string]*cg.ClassHighLevel
	mainClass       *cg.ClassHighLevel
	BuildExpression BuildExpression
}

func (buildPackage *BuildPackage) newClassName(prefix string) (autoName string) {
	for i := 0; i < math.MaxInt16; i++ {
		autoName = fmt.Sprintf("%s$%d", prefix, i)
		if _, exists := buildPackage.Package.Block.NameExists(autoName); exists {
			continue
		}
		autoName = buildPackage.Package.Name + "/" + autoName
		if buildPackage.classes != nil && buildPackage.classes[autoName] != nil {
			continue
		} else {
			return autoName
		}
	}
	panic("new class name overflow")
}

func (buildPackage *BuildPackage) putClass(class *cg.ClassHighLevel) {
	if class.Name == "" {
		panic("missing name")
	}
	name := class.Name
	if name == buildPackage.mainClass.Name {
		panic("cannot have main class`s name")
	}
	if buildPackage.classes == nil {
		buildPackage.classes = make(map[string]*cg.ClassHighLevel)
	}
	if _, ok := buildPackage.classes[name]; ok {
		panic(fmt.Sprintf("name:'%s' already been token", name))
	}
	buildPackage.classes[name] = class
}

func (buildPackage *BuildPackage) Make(p *ast.Package) {
	buildPackage.Package = p
	mainClass := &cg.ClassHighLevel{}
	buildPackage.mainClass = mainClass
	mainClass.AccessFlags |= cg.ACC_CLASS_PUBLIC
	mainClass.AccessFlags |= cg.ACC_CLASS_FINAL
	mainClass.AccessFlags |= cg.ACC_CLASS_SYNTHETIC
	mainClass.SuperClass = ast.JavaRootClass
	mainClass.Name = p.Name + "/main"
	mainClass.Fields = make(map[string]*cg.FieldHighLevel)
	buildPackage.mkClassDefaultConstruction(buildPackage.mainClass, nil)
	buildPackage.BuildExpression.BuildPackage = buildPackage
	buildPackage.classes = make(map[string]*cg.ClassHighLevel)
	buildPackage.mkGlobalConstants()
	buildPackage.mkGlobalTypeAlias()
	buildPackage.mkGlobalVariables()
	buildPackage.mkGlobalFunctions()
	buildPackage.mkInitFunctions()
	for _, v := range p.Block.Classes {
		buildPackage.classes[v.Name] = buildPackage.buildClass(v)
	}
	for _, v := range p.Block.Enums {
		buildPackage.classes[v.Name] = buildPackage.mkEnum(v)
	}
	err := buildPackage.DumpClass()
	if err != nil {
		panic(fmt.Sprintf("dump to file failed,err:%v\n", err))
	}
}

func (buildPackage *BuildPackage) mkEnum(e *ast.Enum) *cg.ClassHighLevel {
	class := &cg.ClassHighLevel{}
	class.Name = e.Name
	class.SourceFiles = make(map[string]struct{})
	class.SourceFiles[e.Pos.Filename] = struct{}{}
	class.AccessFlags = e.AccessFlags
	class.SuperClass = ast.JavaRootClass
	class.Fields = make(map[string]*cg.FieldHighLevel)
	class.Class.AttributeLucyEnum = &cg.AttributeLucyEnum{}
	for _, v := range e.Enums {
		field := &cg.FieldHighLevel{}
		if e.AccessFlags&cg.ACC_CLASS_PUBLIC != 0 {
			field.AccessFlags |= cg.ACC_FIELD_PUBLIC
		} else {
			field.AccessFlags |= cg.ACC_FIELD_PRIVATE
		}
		field.Name = v.Name
		field.Descriptor = "I"
		field.AttributeConstantValue = &cg.AttributeConstantValue{}
		field.AttributeConstantValue.Index = class.Class.InsertIntConst(v.Value)
		class.Fields[v.Name] = field
	}
	return class
}

func (buildPackage *BuildPackage) mkGlobalConstants() {
	for k, v := range buildPackage.Package.Block.Constants {
		f := &cg.FieldHighLevel{}
		f.AccessFlags |= cg.ACC_FIELD_STATIC
		f.AccessFlags |= cg.ACC_FIELD_SYNTHETIC
		f.AccessFlags |= cg.ACC_FIELD_FINAL
		if v.AccessFlags&cg.ACC_FIELD_PUBLIC != 0 {
			f.AccessFlags |= cg.ACC_FIELD_PUBLIC
		}
		f.Name = v.Name
		f.AttributeConstantValue = &cg.AttributeConstantValue{}
		f.AttributeConstantValue.Index = buildPackage.insertDefaultValue(buildPackage.mainClass, v.Type, v.Value)
		f.AttributeLucyConst = &cg.AttributeLucyConst{}
		f.Descriptor = Descriptor.typeDescriptor(v.Type)
		buildPackage.mainClass.Fields[k] = f
	}
}
func (buildPackage *BuildPackage) mkGlobalTypeAlias() {
	for name, v := range buildPackage.Package.Block.TypeAliases {
		t := &cg.AttributeLucyTypeAlias{}
		t.Alias = LucyTypeAliasParser.Encode(name, v)
		buildPackage.mainClass.Class.TypeAlias = append(buildPackage.mainClass.Class.TypeAlias, t)
	}
}

func (buildPackage *BuildPackage) mkGlobalVariables() {
	for k, v := range buildPackage.Package.Block.Variables {
		f := &cg.FieldHighLevel{}
		f.AccessFlags |= v.AccessFlags
		f.AccessFlags |= cg.ACC_FIELD_STATIC
		f.Descriptor = Descriptor.typeDescriptor(v.Type)
		if v.AccessFlags&cg.ACC_FIELD_PUBLIC != 0 {
			f.AccessFlags |= cg.ACC_FIELD_PUBLIC
		}
		f.AccessFlags |= cg.ACC_FIELD_VOLATILE
		if LucyFieldSignatureParser.Need(v.Type) {
			f.AttributeLucyFieldDescriptor = &cg.AttributeLucyFieldDescriptor{}
			f.AttributeLucyFieldDescriptor.Descriptor = LucyFieldSignatureParser.Encode(v.Type)
		}
		f.Name = v.Name
		buildPackage.mainClass.Fields[k] = f
	}
}

func (buildPackage *BuildPackage) mkInitFunctions() {
	if len(buildPackage.Package.InitFunctions) == 0 {
		needTrigger := false
		for _, v := range buildPackage.Package.LoadedPackages {
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
	for _, v := range buildPackage.Package.InitFunctions {
		method := &cg.MethodHighLevel{}
		blockMethods = append(blockMethods, method)
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		method.AccessFlags |= cg.ACC_METHOD_FINAL
		method.AccessFlags |= cg.ACC_METHOD_PRIVATE
		method.Name = buildPackage.mainClass.NewFunctionName("block")
		method.Class = buildPackage.mainClass
		method.Descriptor = "()V"
		method.Code = &cg.AttributeCode{}
		buildPackage.buildFunction(buildPackage.mainClass, nil, method, v)
		buildPackage.mainClass.AppendMethod(method)
	}
	method := &cg.MethodHighLevel{}
	method.AccessFlags |= cg.ACC_METHOD_STATIC
	method.Name = "<clinit>"
	method.Descriptor = "()V"
	codes := make([]byte, 65536)
	codeLength := int(0)
	method.Code = &cg.AttributeCode{}
	for _, v := range buildPackage.Package.LoadedPackages {
		if v.TriggerPackageInitMethodName == "" {
			continue
		}
		codes[codeLength] = cg.OP_invokestatic
		buildPackage.mainClass.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      v.Name + "/main", // main class
			Method:     v.TriggerPackageInitMethodName,
			Descriptor: "()V",
		}, codes[codeLength+1:codeLength+3])
		codeLength += 3
	}
	for _, v := range blockMethods {
		codes[codeLength] = cg.OP_invokestatic
		buildPackage.mainClass.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      buildPackage.mainClass.Name,
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
	buildPackage.mainClass.AppendMethod(method)

	// trigger init
	trigger := &cg.MethodHighLevel{}
	trigger.Name = buildPackage.mainClass.NewFunctionName("triggerPackageInit")
	trigger.AccessFlags |= cg.ACC_METHOD_PUBLIC
	trigger.AccessFlags |= cg.ACC_METHOD_BRIDGE
	trigger.AccessFlags |= cg.ACC_METHOD_STATIC
	trigger.Descriptor = "()V"
	trigger.Code = &cg.AttributeCode{}
	trigger.Code.Codes = make([]byte, 1)
	trigger.Code.Codes[0] = cg.OP_return
	trigger.Code.CodeLength = 1
	trigger.AttributeLucyTriggerPackageInitMethod = &cg.AttributeLucyTriggerPackageInitMethod{}
	buildPackage.mainClass.AppendMethod(trigger)
	buildPackage.mainClass.TriggerPackageInitMethod = trigger
}
func (buildPackage *BuildPackage) insertDefaultValue(c *cg.ClassHighLevel, t *ast.Type, v interface{}) (index uint16) {
	switch t.Type {
	case ast.VariableTypeBool:
		if v.(bool) {
			index = c.Class.InsertIntConst(1)
		} else {
			index = c.Class.InsertIntConst(0)
		}
	case ast.VariableTypeByte:
		index = c.Class.InsertIntConst(int32(v.(byte)))
	case ast.VariableTypeShort:
		fallthrough
	case ast.VariableTypeInt:
		index = c.Class.InsertIntConst(v.(int32))
	case ast.VariableTypeLong:
		index = c.Class.InsertLongConst(v.(int64))
	case ast.VariableTypeFloat:
		index = c.Class.InsertFloatConst(v.(float32))
	case ast.VariableTypeDouble:
		index = c.Class.InsertDoubleConst(v.(float64))
	case ast.VariableTypeString:
		index = c.Class.InsertStringConst(v.(string))
	}
	return
}
func (buildPackage *BuildPackage) buildClass(c *ast.Class) *cg.ClassHighLevel {
	class := &cg.ClassHighLevel{}
	class.Name = c.Name
	class.SourceFiles = make(map[string]struct{})
	class.SourceFiles[c.Pos.Filename] = struct{}{}
	class.AccessFlags = c.AccessFlags
	if c.SuperClass != nil {
		class.SuperClass = c.SuperClass.Name
	} else {
		class.SuperClass = c.SuperClassName
	}
	class.Fields = make(map[string]*cg.FieldHighLevel)
	class.Methods = make(map[string][]*cg.MethodHighLevel)
	for _, v := range c.Interfaces {
		class.Interfaces = append(class.Interfaces, v.Name)
	}
	for _, v := range c.Fields {
		f := &cg.FieldHighLevel{}
		f.Name = v.Name
		f.AccessFlags = v.AccessFlags
		if v.IsStatic() && v.DefaultValue != nil {
			f.AttributeConstantValue = &cg.AttributeConstantValue{}
			f.AttributeConstantValue.Index = buildPackage.insertDefaultValue(class, v.Type, v.DefaultValue)
		}
		f.Descriptor = Descriptor.typeDescriptor(v.Type)
		if LucyFieldSignatureParser.Need(v.Type) {
			t := &cg.AttributeLucyFieldDescriptor{}
			t.Descriptor = LucyFieldSignatureParser.Encode(v.Type)
			f.AttributeLucyFieldDescriptor = t
		}
		class.Fields[v.Name] = f
	}

	for k, v := range c.Methods {
		if k == ast.SpecialMethodInit && c.IsInterface() == false {
			continue
		}
		vv := v[0]
		method := &cg.MethodHighLevel{}
		method.Name = vv.Function.Name
		method.AccessFlags = vv.Function.AccessFlags
		if c.IsInterface() {
			method.AccessFlags |= cg.ACC_METHOD_PUBLIC
			method.AccessFlags |= cg.ACC_METHOD_ABSTRACT
		}
		method.Class = class
		method.Descriptor = Descriptor.methodDescriptor(&vv.Function.Type)
		if c.IsInterface() == false {
			method.Code = &cg.AttributeCode{}
			buildPackage.buildFunction(class, nil, method, vv.Function)
		}
		class.AppendMethod(method)
	}
	if c.IsInterface() == false {
		//construction
		if t := c.Methods[ast.SpecialMethodInit]; t != nil && len(t) > 0 {
			method := &cg.MethodHighLevel{}
			method.Name = "<init>"
			method.AccessFlags = t[0].Function.AccessFlags
			method.Class = class
			method.Descriptor = Descriptor.methodDescriptor(&t[0].Function.Type)
			method.IsConstruction = true
			method.Code = &cg.AttributeCode{}
			buildPackage.buildFunction(class, c, method, t[0].Function)
			class.AppendMethod(method)

		}
	}
	return class
}

func (buildPackage *BuildPackage) mkGlobalFunctions() {
	ms := make(map[string]*cg.MethodHighLevel)
	for k, f := range buildPackage.Package.Block.Functions { // first round
		if f.TemplateFunction != nil {
			buildPackage.mainClass.TemplateFunctions = append(buildPackage.mainClass.TemplateFunctions, &cg.AttributeTemplateFunction{
				Name:        f.Name,
				Filename:    f.Pos.Filename,
				StartLine:   uint16(f.Pos.StartLine),
				StartColumn: uint16(f.Pos.StartColumn),
				Code:        string(f.SourceCodes),
			})
			continue
		}
		if f.IsBuildIn { //
			continue
		}
		class := buildPackage.mainClass
		method := &cg.MethodHighLevel{}
		method.Class = class
		method.Name = f.Name
		if f.Name == ast.MainFunctionName {
			method.Descriptor = "([Ljava/lang/String;)V"
		} else {
			method.Descriptor = Descriptor.methodDescriptor(&f.Type)
		}
		method.AccessFlags = 0
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		if f.AccessFlags&cg.ACC_METHOD_PUBLIC != 0 || f.Name == ast.MainFunctionName {
			method.AccessFlags |= cg.ACC_METHOD_PUBLIC
		}
		ms[k] = method
		f.ClassMethod = method
		method.Code = &cg.AttributeCode{}
		buildPackage.mainClass.AppendMethod(method)
	}
	for k, f := range buildPackage.Package.Block.Functions { // first round
		if f.IsBuildIn || f.TemplateFunction != nil { //
			continue
		}
		buildPackage.buildFunction(ms[k].Class, nil, ms[k], f)
	}
}

func (buildPackage *BuildPackage) DumpClass() error {
	//dump main class
	f, err := os.OpenFile("main.class", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if err := buildPackage.mainClass.ToLow(common.CompileFlags.JvmVersion).OutPut(f); err != nil {
		f.Close()
		return err
	}
	f.Close()
	for _, c := range buildPackage.classes {
		f, err = os.OpenFile(filepath.Base(c.Name)+".class", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		if err = c.ToLow(common.CompileFlags.JvmVersion).OutPut(f); err != nil {
			f.Close()
			return err
		} else {
			f.Close()
		}
	}
	return nil
}

/*
	make a default construction
*/
func (buildPackage *BuildPackage) mkClassDefaultConstruction(class *cg.ClassHighLevel, astClass *ast.Class) {
	method := &cg.MethodHighLevel{}
	method.Name = specialMethodInit
	method.Descriptor = "()V"
	method.AccessFlags |= cg.ACC_METHOD_PUBLIC
	method.Code = &cg.AttributeCode{}
	method.Code.Codes = make([]byte, 65536)
	defer func() {
		method.Code.Codes = method.Code.Codes[0:method.Code.CodeLength]
	}()
	method.Code.Codes[0] = cg.OP_aload_0
	method.Code.Codes[1] = cg.OP_invokespecial
	class.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
		Class:      class.SuperClass,
		Method:     specialMethodInit,
		Descriptor: "()V",
	}, method.Code.Codes[2:4])
	if 1 > method.Code.MaxStack {
		method.Code.MaxStack = 1
	}
	method.Code.CodeLength = 4
	method.Code.Codes[method.Code.CodeLength] = cg.OP_return
	method.Code.CodeLength += 1
	method.Code.MaxLocals = 1
	class.AppendMethod(method)
}

func (buildPackage *BuildPackage) storeGlobalVariable(class *cg.ClassHighLevel, code *cg.AttributeCode,
	v *ast.Variable) {
	code.Codes[code.CodeLength] = cg.OP_putstatic
	class.InsertFieldRefConst(cg.CONSTANT_Fieldref_info_high_level{
		Class:      buildPackage.mainClass.Name,
		Field:      v.Name,
		Descriptor: Descriptor.typeDescriptor(v.Type),
	}, code.Codes[code.CodeLength+1:code.CodeLength+3])
	code.CodeLength += 3
}

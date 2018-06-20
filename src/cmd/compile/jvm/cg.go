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

type MakeClass struct {
	Package        *ast.Package
	classes        map[string]*cg.ClassHighLevel
	mainClass      *cg.ClassHighLevel
	makeExpression MakeExpression
}

func (makeClass *MakeClass) newClassName(prefix string) (autoName string) {
	for i := 0; i < math.MaxInt16; i++ {
		autoName = fmt.Sprintf("%s_%d", prefix, i)
		if _, exists := makeClass.Package.Block.NameExists(autoName); exists {
			continue
		}
		autoName = makeClass.Package.Name + "/" + autoName
		if makeClass.classes != nil && makeClass.classes[autoName] != nil {
			continue
		} else {
			return autoName
		}
	}
	panic("new class name overflow")
}

func (makeClass *MakeClass) putClass(name string, class *cg.ClassHighLevel) {
	if name == makeClass.mainClass.Name {
		panic("cannot have main class`s name")
	}
	if makeClass.classes == nil {
		makeClass.classes = make(map[string]*cg.ClassHighLevel)
	}
	if _, ok := makeClass.classes[name]; ok {
		panic(fmt.Sprintf("name:'%s' already been token", name))
	}
	makeClass.classes[name] = class
}

func (makeClass *MakeClass) Make(p *ast.Package) {
	makeClass.Package = p
	mainClass := &cg.ClassHighLevel{}
	makeClass.mainClass = mainClass
	mainClass.AccessFlags |= cg.ACC_CLASS_PUBLIC
	mainClass.AccessFlags |= cg.ACC_CLASS_FINAL
	mainClass.AccessFlags |= cg.ACC_CLASS_SYNTHETIC
	mainClass.SuperClass = ast.JAVA_ROOT_CLASS
	mainClass.Name = p.Name + "/main"
	mainClass.Fields = make(map[string]*cg.FieldHighLevel)
	makeClass.mkClassDefaultConstruction(makeClass.mainClass, nil)
	makeClass.makeExpression.MakeClass = makeClass
	makeClass.classes = make(map[string]*cg.ClassHighLevel)
	makeClass.mkGlobalConstants()
	makeClass.mkGlobalTypeAlias()
	makeClass.mkGlobalVariables()
	makeClass.mkGlobalFunctions()
	makeClass.mkInitFunctions()
	for _, v := range p.Block.Classes {
		makeClass.classes[v.Name] = makeClass.buildClass(v)
	}
	for _, v := range p.Block.Enums {
		makeClass.classes[v.Name] = makeClass.mkEnum(v)
	}
	err := makeClass.DumpClass()
	if err != nil {
		panic(fmt.Sprintf("dump to file failed,err:%v\n", err))
	}
}

func (makeClass *MakeClass) mkEnum(e *ast.Enum) *cg.ClassHighLevel {
	class := &cg.ClassHighLevel{}
	class.Name = e.Name
	class.SourceFiles = make(map[string]struct{})
	class.SourceFiles[e.Pos.Filename] = struct{}{}
	class.AccessFlags = e.AccessFlags
	class.SuperClass = ast.JAVA_ROOT_CLASS
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

func (makeClass *MakeClass) mkGlobalConstants() {
	for k, v := range makeClass.Package.Block.Constants {
		f := &cg.FieldHighLevel{}
		f.AccessFlags |= cg.ACC_FIELD_STATIC
		f.AccessFlags |= cg.ACC_FIELD_SYNTHETIC
		f.AccessFlags |= cg.ACC_FIELD_FINAL
		if v.AccessFlags&cg.ACC_FIELD_PUBLIC != 0 {
			f.AccessFlags |= cg.ACC_FIELD_PUBLIC
		}
		f.Name = v.Name
		f.AttributeConstantValue = &cg.AttributeConstantValue{}
		f.AttributeConstantValue.Index = makeClass.insertDefaultValue(makeClass.mainClass, v.Type, v.Value)
		f.AttributeLucyConst = &cg.AttributeLucyConst{}
		f.Descriptor = Descriptor.typeDescriptor(v.Type)
		makeClass.mainClass.Fields[k] = f
	}
}
func (makeClass *MakeClass) mkGlobalTypeAlias() {
	for name, v := range makeClass.Package.Block.TypeAliases {
		t := &cg.AttributeLucyTypeAlias{}
		t.Alias = LucyTypeAliasParser.Encode(name, v)
		makeClass.mainClass.Class.TypeAlias = append(makeClass.mainClass.Class.TypeAlias, t)
	}
}

func (makeClass *MakeClass) mkGlobalVariables() {
	for k, v := range makeClass.Package.Block.Variables {
		f := &cg.FieldHighLevel{}
		f.AccessFlags |= v.AccessFlags
		f.AccessFlags |= cg.ACC_FIELD_STATIC
		f.Descriptor = Descriptor.typeDescriptor(v.Type)
		if v.AccessFlags&cg.ACC_FIELD_PUBLIC != 0 {
			f.AccessFlags |= cg.ACC_FIELD_PUBLIC
		}
		if LucyFieldSignatureParser.Need(v.Type) {
			f.AttributeLucyFieldDescriptor = &cg.AttributeLucyFieldDescriptor{}
			f.AttributeLucyFieldDescriptor.Descriptor = LucyFieldSignatureParser.Encode(v.Type)
		}
		f.Name = v.Name
		makeClass.mainClass.Fields[k] = f
	}
}

func (makeClass *MakeClass) mkInitFunctions() {
	if len(makeClass.Package.InitFunctions) == 0 {
		needTrigger := false
		for _, v := range makeClass.Package.LoadedPackages {
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
	for _, v := range makeClass.Package.InitFunctions {
		method := &cg.MethodHighLevel{}
		blockMethods = append(blockMethods, method)
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		method.AccessFlags |= cg.ACC_METHOD_FINAL
		method.AccessFlags |= cg.ACC_METHOD_PRIVATE
		method.Name = makeClass.mainClass.NewFunctionName("block")
		method.Class = makeClass.mainClass
		method.Descriptor = "()V"
		method.Code = &cg.AttributeCode{}
		makeClass.buildFunction(makeClass.mainClass, nil, method, v)
		makeClass.mainClass.AppendMethod(method)
	}
	method := &cg.MethodHighLevel{}
	method.AccessFlags |= cg.ACC_METHOD_STATIC
	method.Name = "<clinit>"
	method.Descriptor = "()V"
	codes := make([]byte, 65536)
	codeLength := int(0)
	method.Code = &cg.AttributeCode{}
	for _, v := range makeClass.Package.LoadedPackages {
		if v.TriggerPackageInitMethodName == "" {
			continue
		}
		codes[codeLength] = cg.OP_invokestatic
		makeClass.mainClass.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      v.Name + "/main", // main class
			Method:     v.TriggerPackageInitMethodName,
			Descriptor: "()V",
		}, codes[codeLength+1:codeLength+3])
		codeLength += 3
	}
	for _, v := range blockMethods {
		codes[codeLength] = cg.OP_invokestatic
		makeClass.mainClass.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      makeClass.mainClass.Name,
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
	makeClass.mainClass.AppendMethod(method)

	// trigger init
	trigger := &cg.MethodHighLevel{}
	trigger.Name = makeClass.mainClass.NewFunctionName("triggerPackageInit")
	trigger.AccessFlags |= cg.ACC_METHOD_PUBLIC
	trigger.AccessFlags |= cg.ACC_METHOD_BRIDGE
	trigger.AccessFlags |= cg.ACC_METHOD_STATIC
	trigger.Descriptor = "()V"
	trigger.Code = &cg.AttributeCode{}
	trigger.Code.Codes = make([]byte, 1)
	trigger.Code.Codes[0] = cg.OP_return
	trigger.Code.CodeLength = 1
	trigger.AttributeLucyTriggerPackageInitMethod = &cg.AttributeLucyTriggerPackageInitMethod{}
	makeClass.mainClass.AppendMethod(trigger)
	makeClass.mainClass.TriggerPackageInitMethod = trigger
}
func (makeClass *MakeClass) insertDefaultValue(c *cg.ClassHighLevel, t *ast.Type, v interface{}) (index uint16) {
	switch t.Type {
	case ast.VARIABLE_TYPE_BOOL:
		if v.(bool) {
			index = c.Class.InsertIntConst(1)
		} else {
			index = c.Class.InsertIntConst(0)
		}
	case ast.VARIABLE_TYPE_BYTE:
		index = c.Class.InsertIntConst(int32(v.(byte)))
	case ast.VARIABLE_TYPE_SHORT:
		fallthrough
	case ast.VARIABLE_TYPE_INT:
		index = c.Class.InsertIntConst(v.(int32))
	case ast.VARIABLE_TYPE_LONG:
		index = c.Class.InsertLongConst(v.(int64))
	case ast.VARIABLE_TYPE_FLOAT:
		index = c.Class.InsertFloatConst(v.(float32))
	case ast.VARIABLE_TYPE_DOUBLE:
		index = c.Class.InsertDoubleConst(v.(float64))
	case ast.VARIABLE_TYPE_STRING:
		index = c.Class.InsertStringConst(v.(string))
	}
	return
}
func (makeClass *MakeClass) buildClass(c *ast.Class) *cg.ClassHighLevel {
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
			f.AttributeConstantValue.Index = makeClass.insertDefaultValue(class, v.Type, v.DefaultValue)
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
		if k == ast.CONSTRUCTION_METHOD_NAME && c.IsInterface() == false {
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
		method.Descriptor = Descriptor.methodDescriptor(vv.Function)
		if c.IsInterface() == false {
			method.Code = &cg.AttributeCode{}
			makeClass.buildFunction(class, nil, method, vv.Function)
		}
		class.AppendMethod(method)
	}
	if c.IsInterface() == false {
		//construction
		if t := c.Methods[ast.CONSTRUCTION_METHOD_NAME]; t != nil && len(t) > 0 {
			method := &cg.MethodHighLevel{}
			method.Name = "<init>"
			method.AccessFlags = t[0].Function.AccessFlags
			method.Class = class
			method.Descriptor = Descriptor.methodDescriptor(t[0].Function)
			method.IsConstruction = true
			method.Code = &cg.AttributeCode{}
			makeClass.buildFunction(class, c, method, t[0].Function)
			class.AppendMethod(method)
			if len(t[0].Function.Type.ParameterList) > 0 {
				makeClass.mkClassDefaultConstruction(class, c)
			}
		} else {
			makeClass.mkClassDefaultConstruction(class, c)
		}
	}
	return class
}

func (makeClass *MakeClass) mkGlobalFunctions() {
	ms := make(map[string]*cg.MethodHighLevel)
	for k, f := range makeClass.Package.Block.Functions { // first round
		if f.TemplateFunction != nil {
			makeClass.mainClass.TemplateFunctions = append(makeClass.mainClass.TemplateFunctions, &cg.AttributeTemplateFunction{
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
		class := makeClass.mainClass
		method := &cg.MethodHighLevel{}
		method.Class = class
		method.Name = f.Name
		method.Descriptor = Descriptor.methodDescriptor(f)
		method.AccessFlags = 0
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		if f.AccessFlags&cg.ACC_METHOD_PUBLIC != 0 || f.Name == ast.MAIN_FUNCTION_NAME {
			method.AccessFlags |= cg.ACC_METHOD_PUBLIC
		}
		ms[k] = method
		f.ClassMethod = method
		method.Code = &cg.AttributeCode{}
		makeClass.mainClass.AppendMethod(method)
	}
	for k, f := range makeClass.Package.Block.Functions { // first round
		if f.IsBuildIn || f.TemplateFunction != nil { //
			continue
		}
		makeClass.buildFunction(ms[k].Class, nil, ms[k], f)
	}
}

func (makeClass *MakeClass) DumpClass() error {
	//dump main class
	f, err := os.OpenFile("main.class", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if err := makeClass.mainClass.ToLow(common.CompileFlags.JvmVersion).OutPut(f); err != nil {
		f.Close()
		return err
	}
	f.Close()
	for _, c := range makeClass.classes {
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
func (makeClass *MakeClass) mkClassDefaultConstruction(class *cg.ClassHighLevel, astClass *ast.Class) {
	method := &cg.MethodHighLevel{}
	method.Name = special_method_init
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
		Method:     special_method_init,
		Descriptor: "()V",
	}, method.Code.Codes[2:4])
	if 1 > method.Code.MaxStack {
		method.Code.MaxStack = 1
	}
	method.Code.CodeLength = 4
	if astClass != nil {
		makeClass.mkFieldDefaultValue(class, method.Code, &Context{class: astClass}, nil)
	}
	method.Code.Codes[method.Code.CodeLength] = cg.OP_return
	method.Code.CodeLength += 1
	method.Code.MaxLocals = 1
	class.AppendMethod(method)
}

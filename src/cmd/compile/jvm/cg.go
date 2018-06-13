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
	p              *ast.Package
	Classes        map[string]*cg.ClassHighLevel
	mainClass      *cg.ClassHighLevel
	MakeExpression MakeExpression
}

func (m *MakeClass) newClassName(prefix string) (autoName string) {
	for i := 0; i < math.MaxInt16; i++ {
		autoName = fmt.Sprintf("%s_%d", prefix, i)
		if _, exists := m.p.Block.NameExists(autoName); exists {
			continue
		}
		autoName = m.p.Name + "/" + autoName
		if m.Classes != nil && m.Classes[autoName] != nil {
			continue
		} else {
			return autoName
		}
	}
	panic("new class name overflow")
}

func (m *MakeClass) putClass(name string, class *cg.ClassHighLevel) {
	if name == m.mainClass.Name {
		panic("cannot have main class`s name")
	}
	if m.Classes == nil {
		m.Classes = make(map[string]*cg.ClassHighLevel)
	}
	if _, ok := m.Classes[name]; ok {
		panic(fmt.Sprintf("name:'%s' already been token", name))
	}
	m.Classes[name] = class
}

func (m *MakeClass) Make(p *ast.Package) {
	m.p = p
	mainClass := &cg.ClassHighLevel{}
	m.mainClass = mainClass
	mainClass.AccessFlags |= cg.ACC_CLASS_PUBLIC
	mainClass.AccessFlags |= cg.ACC_CLASS_FINAL
	mainClass.AccessFlags |= cg.ACC_CLASS_SYNTHETIC
	mainClass.SuperClass = ast.JAVA_ROOT_CLASS
	mainClass.Name = p.Name + "/main"
	mainClass.Fields = make(map[string]*cg.FieldHighLevel)
	m.mkClassDefaultConstruction(m.mainClass, nil)
	m.MakeExpression.MakeClass = m
	m.Classes = make(map[string]*cg.ClassHighLevel)
	m.mkConsts()
	m.mkTypes()
	m.mkVars()
	m.mkFuncs()
	m.mkInitFunctions()
	for _, v := range p.Block.Classes {
		m.Classes[v.Name] = m.buildClass(v)
	}
	for _, v := range p.Block.Enums {
		m.Classes[v.Name] = m.mkEnum(v)
	}
	err := m.Dump()
	if err != nil {
		panic(fmt.Sprintf("dump to file failed,err:%v\n", err))
	}
}

func (m *MakeClass) mkEnum(e *ast.Enum) *cg.ClassHighLevel {
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

func (m *MakeClass) mkConsts() {
	for k, v := range m.p.Block.Consts {
		f := &cg.FieldHighLevel{}
		f.AccessFlags |= cg.ACC_FIELD_STATIC
		f.AccessFlags |= cg.ACC_FIELD_SYNTHETIC
		f.AccessFlags |= cg.ACC_FIELD_FINAL
		if v.AccessFlags&cg.ACC_FIELD_PUBLIC != 0 {
			f.AccessFlags |= cg.ACC_FIELD_PUBLIC
		}
		f.Name = v.Name
		f.AttributeConstantValue = &cg.AttributeConstantValue{}
		f.AttributeConstantValue.Index = m.insertDefaultValue(m.mainClass, v.Typ, v.Value)
		f.AttributeLucyConst = &cg.AttributeLucyConst{}
		f.Descriptor = Descriptor.typeDescriptor(v.Typ)
		m.mainClass.Fields[k] = f
	}
}
func (m *MakeClass) mkTypes() {
	for name, v := range m.p.Block.Types {
		t := &cg.AttributeLucyTypeAlias{}
		t.Alias = LucyTypeAliasParser.Encode(name, v)
		m.mainClass.Class.TypeAlias = append(m.mainClass.Class.TypeAlias, t)
	}
}

func (m *MakeClass) mkVars() {
	for k, v := range m.p.Block.Vars {
		f := &cg.FieldHighLevel{}
		f.AccessFlags |= v.AccessFlags
		f.AccessFlags |= cg.ACC_FIELD_STATIC
		f.Descriptor = Descriptor.typeDescriptor(v.Typ)
		if v.AccessFlags&cg.ACC_FIELD_PUBLIC != 0 {
			f.AccessFlags |= cg.ACC_FIELD_PUBLIC
		}
		if LucyFieldSignatureParser.Need(v.Typ) {
			f.AttributeLucyFieldDescritor = &cg.AttributeLucyFieldDescriptor{}
			f.AttributeLucyFieldDescritor.Descriptor = LucyFieldSignatureParser.Encode(v.Typ)
		}
		f.Name = v.Name
		m.mainClass.Fields[k] = f
	}
}

func (m *MakeClass) mkInitFunctions() {
	if len(m.p.InitFunctions) == 0 {
		needTrigger := false
		for _, v := range m.p.LoadedPackages {
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
	for _, v := range m.p.InitFunctions {
		method := &cg.MethodHighLevel{}
		blockMethods = append(blockMethods, method)
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		method.AccessFlags |= cg.ACC_METHOD_FINAL
		method.AccessFlags |= cg.ACC_METHOD_PRIVATE
		method.Name = m.mainClass.NewFunctionName("block")
		method.Class = m.mainClass
		method.Descriptor = "()V"
		method.Code = &cg.AttributeCode{}
		m.buildFunction(m.mainClass, nil, method, v)
		m.mainClass.AppendMethod(method)
	}
	method := &cg.MethodHighLevel{}
	method.AccessFlags |= cg.ACC_METHOD_STATIC
	method.Name = "<clinit>"
	method.Descriptor = "()V"
	codes := make([]byte, 65536)
	codeLength := int(0)
	method.Code = &cg.AttributeCode{}
	for _, v := range m.p.LoadedPackages {
		if v.TriggerPackageInitMethodName == "" {
			continue
		}
		codes[codeLength] = cg.OP_invokestatic
		m.mainClass.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      v.Name + "/main", // main class
			Method:     v.TriggerPackageInitMethodName,
			Descriptor: "()V",
		}, codes[codeLength+1:codeLength+3])
		codeLength += 3
	}
	for _, v := range blockMethods {
		codes[codeLength] = cg.OP_invokestatic
		m.mainClass.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      m.mainClass.Name,
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
	m.mainClass.AppendMethod(method)

	// trigger init
	trigger := &cg.MethodHighLevel{}
	trigger.Name = m.mainClass.NewFunctionName("triggerPackageInit")
	trigger.AccessFlags |= cg.ACC_METHOD_PUBLIC
	trigger.AccessFlags |= cg.ACC_METHOD_BRIDGE
	trigger.AccessFlags |= cg.ACC_METHOD_STATIC
	trigger.Descriptor = "()V"
	trigger.Code = &cg.AttributeCode{}
	trigger.Code.Codes = make([]byte, 1)
	trigger.Code.Codes[0] = cg.OP_return
	trigger.Code.CodeLength = 1
	trigger.AttributeLucyTriggerPackageInitMethod = &cg.AttributeLucyTriggerPackageInitMethod{}
	m.mainClass.AppendMethod(trigger)
	m.mainClass.TriggerCLinit = trigger
}
func (m *MakeClass) insertDefaultValue(c *cg.ClassHighLevel, t *ast.VariableType, v interface{}) (index uint16) {
	switch t.Typ {
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
func (m *MakeClass) buildClass(c *ast.Class) *cg.ClassHighLevel {
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
			f.AttributeConstantValue.Index = m.insertDefaultValue(class, v.Typ, v.DefaultValue)
		}
		f.Descriptor = Descriptor.typeDescriptor(v.Typ)
		if LucyFieldSignatureParser.Need(v.Typ) {
			t := &cg.AttributeLucyFieldDescriptor{}
			t.Descriptor = LucyFieldSignatureParser.Encode(v.Typ)
			f.AttributeLucyFieldDescritor = t
		}
		class.Fields[v.Name] = f
	}

	for k, v := range c.Methods {
		if k == ast.CONSTRUCTION_METHOD_NAME && c.IsInterface() == false {
			continue
		}
		vv := v[0]
		method := &cg.MethodHighLevel{}
		method.Name = vv.Func.Name
		method.AccessFlags = vv.Func.AccessFlags
		if c.IsInterface() {
			method.AccessFlags |= cg.ACC_METHOD_PUBLIC
			method.AccessFlags |= cg.ACC_METHOD_ABSTRACT
		}
		method.Class = class
		method.Descriptor = Descriptor.methodDescriptor(vv.Func)
		if c.IsInterface() == false {
			method.Code = &cg.AttributeCode{}
			m.buildFunction(class, nil, method, vv.Func)
		}
		class.AppendMethod(method)
	}
	if c.IsInterface() == false {
		//construction
		if t := c.Methods[ast.CONSTRUCTION_METHOD_NAME]; t != nil && len(t) > 0 {
			method := &cg.MethodHighLevel{}
			method.Name = "<init>"
			method.AccessFlags = t[0].Func.AccessFlags
			method.Class = class
			method.Descriptor = Descriptor.methodDescriptor(t[0].Func)
			method.IsConstruction = true
			method.Code = &cg.AttributeCode{}
			m.buildFunction(class, c, method, t[0].Func)
			class.AppendMethod(method)
			if len(t[0].Func.Typ.ParameterList) > 0 {
				m.mkClassDefaultConstruction(class, c)
			}
		} else {
			m.mkClassDefaultConstruction(class, c)
		}
	}
	return class
}

func (m *MakeClass) mkFuncs() {
	ms := make(map[string]*cg.MethodHighLevel)
	for k, f := range m.p.Block.Funcs { // fisrt round
		if f.TemplateFunction != nil {
			m.mainClass.TemplateFunctions = append(m.mainClass.TemplateFunctions, &cg.AttributeTemplateFunction{
				Name:        f.Name,
				Filename:    f.Pos.Filename,
				StartLine:   uint16(f.Pos.StartLine),
				StartColumn: uint16(f.Pos.StartColumn),
				Code:        string(f.SourceCode),
			})
			continue
		}
		if f.IsBuildIn { //
			continue
		}
		class := m.mainClass
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
		m.mainClass.AppendMethod(method)
	}
	for k, f := range m.p.Block.Funcs { // fisrt round
		if f.IsBuildIn || f.TemplateFunction != nil { //
			continue
		}
		m.buildFunction(ms[k].Class, nil, ms[k], f)
	}
}

func (m *MakeClass) Dump() error {
	//dump main class
	f, err := os.OpenFile("main.class", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if err := m.mainClass.ToLow(common.CompileFlags.JvmVersion).OutPut(f); err != nil {
		f.Close()
		return err
	}
	f.Close()
	for _, c := range m.Classes {
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
func (m *MakeClass) mkClassDefaultConstruction(class *cg.ClassHighLevel, astClass *ast.Class) {
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
		m.mkFieldDefaultValue(class, method.Code, &Context{class: astClass}, nil)
	}
	method.Code.Codes[method.Code.CodeLength] = cg.OP_return
	method.Code.CodeLength += 1
	method.Code.MaxLocals = 1
	class.AppendMethod(method)
}

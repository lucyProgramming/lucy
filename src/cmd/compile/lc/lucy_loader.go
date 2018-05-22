package lc

import (
	"encoding/binary"
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (loader *RealNameLoader) loadAsLucy(c *cg.Class) (*ast.Class, error) {
	if t := c.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_COMPILTER_AUTO); t != nil && len(t) > 0 {
		return nil, nil
	}
	//name
	astClass := &ast.Class{}
	{
		nameindex := binary.BigEndian.Uint16(c.ConstPool[c.ThisClass].Info)
		astClass.Name = string(c.ConstPool[nameindex].Info)
		if astClass.Name != ast.JAVA_ROOT_CLASS {
			nameindex = binary.BigEndian.Uint16(c.ConstPool[c.SuperClass].Info)
			astClass.SuperClassName = string(c.ConstPool[nameindex].Info)
		}
	}
	astClass.AccessFlags = c.AccessFlag
	astClass.LoadFromOutSide = true
	var err error
	astClass.Fields = make(map[string]*ast.ClassField)
	for _, v := range c.Fields {
		f := &ast.ClassField{}
		f.Name = string(c.ConstPool[v.NameIndex].Info)
		f.Descriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		f.LoadFromOutSide = true
		_, f.Typ, err = jvm.Descriptor.ParseType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
		}
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_FIELD_DESCRIPTOR); t != nil && len(t) > 0 {
			index := binary.BigEndian.Uint16(t[0].Info)
			_, f.Typ, err = jvm.LucyFieldSignatureParser.Decode(c.ConstPool[index].Info)
			if err != nil {
				return nil, err
			}
		}
		if f.Typ.Typ == ast.VARIABLE_TYPE_ENUM {
			loadEnumForVariableType(f.Typ)
		}
		f.AccessFlags = v.AccessFlags
		astClass.Fields[f.Name] = f
	}
	astClass.Methods = make(map[string][]*ast.ClassMethod)
	for _, v := range c.Methods {
		m := &ast.ClassMethod{}
		m.Func = &ast.Function{}
		m.Func.Name = string(c.ConstPool[v.NameIndex].Info)
		m.Func.Typ, err = jvm.Descriptor.ParseFunctionType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
		}
		m.Func.AccessFlags = v.AccessFlags
		m.LoadFromOutSide = true
		m.Func.Descriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_METHOD_DESCRIPTOR); t != nil && len(t) > 0 {
			index := binary.BigEndian.Uint16(t[0].Info)
			err = jvm.LucyMethodSignatureParser.Decode(m.Func, c.ConstPool[index].Info)
			if err != nil {
				return nil, err
			}
		}
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_DEFAULT_PARAMETERS); t != nil && len(t) > 0 {
			dp := &cg.AttributeDefaultParameters{}
			dp.FromBytes(t[0].Info)
			jvm.FunctionDefaultValueParser.Decode(c, m.Func, dp)
		}
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_METHOD_PARAMETERS); t != nil && len(t) > 0 {
			parseMethodParameter(c, t[0].Info, m.Func)
		}
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_RETURNLIST_NAMES); t != nil && len(t) > 0 {
			parseReturnListNames(c, t[0].Info, m.Func)
		}
		err = loadEnumForFunction(m.Func)
		if err != nil {
			return nil, err
		}
		if astClass.Methods[m.Func.Name] == nil {
			astClass.Methods[m.Func.Name] = []*ast.ClassMethod{m}
		} else {
			astClass.Methods[m.Func.Name] = append(astClass.Methods[m.Func.Name], m)
		}
	}
	return astClass, nil
}

func (loader *RealNameLoader) loadLucyEnum(c *cg.Class) (*ast.Enum, error) {
	e := &ast.Enum{}
	{
		nameindex := binary.BigEndian.Uint16(c.ConstPool[c.ThisClass].Info)
		e.Name = string(c.ConstPool[nameindex].Info)
	}
	e.AccessFlags = c.AccessFlag
	for _, v := range c.Fields {
		en := &ast.EnumName{}
		name := string(c.ConstPool[v.NameIndex].Info)
		en.Name = name
		en.Enum = e
		constValue := v.AttributeGroupedByName[cg.ATTRIBUTE_NAME_CONST_VALUE][0] // must have this attribute
		en.Value = int32(binary.BigEndian.Uint32(c.ConstPool[binary.BigEndian.Uint16(constValue.Info)].Info))
		e.Enums = append(e.Enums, en)
	}
	return e, nil
}

func (loader *RealNameLoader) loadLucyMainClass(pack *ast.Package, c *cg.Class) error {
	var err error
	mainClassName := &cg.ClassHighLevel{}
	mainClassName.Name = pack.Name + "/main"
	pack.Block.Vars = make(map[string]*ast.VariableDefinition)
	pack.Block.Consts = make(map[string]*ast.Const)
	pack.Block.Funcs = make(map[string]*ast.Function)
	for _, f := range c.Fields {
		name := string(c.ConstPool[f.NameIndex].Info)
		constValue := f.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_CONST_VALUE)
		if len(constValue) > 1 {
			return fmt.Errorf("constant value length greater than  1 at class 'main'  field '%s'", name)
		}
		_, typ, err := jvm.Descriptor.ParseType(c.ConstPool[f.DescriptorIndex].Info)
		if err != nil {
			return err
		}
		if len(f.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_CONST)) > 0 {
			//const
			cos := &ast.Const{}
			cos.Name = name
			cos.AccessFlags = f.AccessFlags
			cos.Typ = typ
			_, cos.Typ, err = jvm.Descriptor.ParseType(c.ConstPool[f.DescriptorIndex].Info)
			if err != nil {
				return err
			}
			valueIndex := binary.BigEndian.Uint16(constValue[0].Info)
			switch cos.Typ.Typ {
			case ast.VARIABLE_TYPE_BOOL:
				cos.Value = binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info) != 0
			case ast.VARIABLE_TYPE_BYTE:
				cos.Value = byte(binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_SHORT:
				fallthrough
			case ast.VARIABLE_TYPE_INT:
				cos.Value = int32(binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_LONG:
				cos.Value = int64(binary.BigEndian.Uint64(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_FLOAT:
				cos.Value = float32(binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_DOUBLE:
				cos.Value = float64(binary.BigEndian.Uint64(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_STRING:
				valueIndex = binary.BigEndian.Uint16(c.ConstPool[valueIndex].Info) // const_string_info
				cos.Value = string(c.ConstPool[valueIndex].Info)                   // utf 8
			}
			pack.Block.Consts[name] = cos
		} else {
			//global vars
			vd := &ast.VariableDefinition{}
			vd.Name = name
			vd.AccessFlags = f.AccessFlags
			vd.Descriptor = string(c.ConstPool[f.DescriptorIndex].Info)
			vd.Typ = typ
			vd.IsGlobal = true
			pack.Block.Vars[name] = vd
			if t := f.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_FIELD_DESCRIPTOR); t != nil && len(t) > 0 {
				index := binary.BigEndian.Uint16(t[0].Info)
				_, vd.Typ, err = jvm.LucyFieldSignatureParser.Decode(c.ConstPool[index].Info)
				if err != nil {
					return err
				}
			}
			if typ.Typ == ast.VARIABLE_TYPE_ENUM {
				loadEnumForVariableType(typ)
			}
		}
	}
	for _, m := range c.Methods {
		if t := m.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_TRIGGER_PACKAGE_INIT); t != nil && len(t) > 0 {
			pack.TriggerPackageInitMethodName = string(c.ConstPool[m.NameIndex].Info)
			continue
		}
		if m.AccessFlags&cg.ACC_METHOD_BRIDGE != 0 {
			continue
		}
		name := string(c.ConstPool[m.NameIndex].Info)
		if name == ast.MAIN_FUNCTION_NAME {
			// this is main function
			continue
		}
		function := &ast.Function{}
		function.Name = name
		function.AccessFlags = m.AccessFlags
		function.Descriptor = string(c.ConstPool[m.DescriptorIndex].Info)
		function.Typ, err = jvm.Descriptor.ParseFunctionType(c.ConstPool[m.DescriptorIndex].Info)
		if err != nil {
			return err
		}
		if t := m.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_METHOD_DESCRIPTOR); t != nil && len(t) > 0 {
			index := binary.BigEndian.Uint16(t[0].Info)
			err = jvm.LucyMethodSignatureParser.Decode(function, c.ConstPool[index].Info)
			if err != nil {
				return err
			}
		}
		err = loadEnumForFunction(function)
		if err != nil {
			return err
		}
		if t := m.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_METHOD_PARAMETERS); t != nil && len(t) > 0 {
			parseMethodParameter(c, t[0].Info, function)
		}
		if t := m.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_RETURNLIST_NAMES); t != nil && len(t) > 0 {
			parseReturnListNames(c, t[0].Info, function)
		}
		function.ClassMethod = &cg.MethodHighLevel{}
		function.ClassMethod.Name = function.Name
		function.ClassMethod.Class = mainClassName
		function.ClassMethod.Descriptor = function.Descriptor
		function.IsGlobal = true
		pack.Block.Funcs[name] = function
	}
	if pack.Block.Types == nil {
		pack.Block.Types = make(map[string]*ast.VariableType)
	}
	for _, v := range c.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_TYPE_ALIAS) {
		index := binary.BigEndian.Uint16(v.Info)
		name, typ, err := jvm.LucyTypeAliasParser.Decode(c.ConstPool[index].Info)
		if err != nil {
			return err
		}
		typ.Alias = name
		pack.Block.Types[name] = typ
		if typ.Typ == ast.VARIABLE_TYPE_ENUM {
			err = loadEnumForVariableType(typ)
			if err != nil {
				return err
			}
		}
	}
	for _, v := range c.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_TEMPLATE_FUNCTION) {
		attr := &cg.AttributeTemplateFunction{}
		attr.FromBytes(c, v.Info)
		f, es := ParseFunctionHandler([]byte(attr.Code), &ast.Pos{
			Filename:    attr.Filename,
			StartLine:   int(attr.StartLine),
			StartColumn: int(attr.StartColumn),
		})
		if len(es) > 0 { // looks impossible
			return es[0]
		}
		f.TemplateFunction = &ast.TemplateFunction{}
		pack.Block.Funcs[attr.Name] = f
	}
	return nil
}

package lc

import (
	"encoding/binary"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (loader *RealNameLoader) loadAsLucy(c *cg.Class) (*ast.Class, error) {
	if t := c.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_CLOSURE_FUNCTION_CLASS); t != nil && len(t) > 0 {
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
	astClass.Access = c.AccessFlag
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
			index := binary.BigEndian.Uint64(t[0].Info)
			_, f.Typ, err = jvm.LucyFieldSignatureParser.Decode(c.ConstPool[index].Info)
			if err != nil {
				return nil, err
			}
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
		m.LoadFromOutSide = true
		m.Func.Descriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_METHOD_DESCRIPTOR); t != nil && len(t) > 0 {
			index := binary.BigEndian.Uint64(t[0].Info)
			err = jvm.LucyMethodSignatureParser.Deocde(c.ConstPool[index].Info, m.Func)
			if err != nil {
				return nil, err
			}
		}
		if m.Func.Name == "<init>" {
			astClass.Constructors = append(astClass.Constructors, m)
		} else {
			if astClass.Methods[m.Func.Name] == nil {
				astClass.Methods[m.Func.Name] = []*ast.ClassMethod{m}
			} else {
				astClass.Methods[m.Func.Name] = append(astClass.Methods[m.Func.Name], m)
			}
		}
	}
	return astClass, nil
}

func (loader *RealNameLoader) loadLucyMainClass(pack *ast.Package, c *cg.Class) error {
	var err error
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
		if constValue != nil && len(constValue) > 0 {
			//const
			cos := &ast.Const{}
			cos.Name = name
			cos.AccessFlags = f.AccessFlags
			cos.Typ = typ
			_, cos.Typ, err = jvm.Descriptor.ParseType(c.ConstPool[f.DescriptorIndex].Info)
			if err != nil {
				return err
			}
			cos.Descriptor = string(c.ConstPool[f.DescriptorIndex].Info)
			if pack.Block.Consts == nil {
				pack.Block.Consts = make(map[string]*ast.Const)
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
			if pack.Block.Vars == nil {
				pack.Block.Vars = make(map[string]*ast.VariableDefinition)
			}
			pack.Block.Vars[name] = vd
		}
	}
	for _, m := range c.Methods {
		if t := m.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_INNER_STATIC_METHOD); t != nil && len(t) > 0 {
			//innsert static method cannot called from outside
			continue
		}
		function := &ast.Function{}
		function.Name = string(c.ConstPool[m.NameIndex].Info)
		function.AccessFlags = m.AccessFlags
		function.Typ, err = jvm.Descriptor.ParseFunctionType(c.ConstPool[m.DescriptorIndex].Info)
		if err != nil {
			return err
		}
		function.IsGlobal = true
		if pack.Block.Funcs == nil {
			pack.Block.Funcs = make(map[string]*ast.Function)
		}
		pack.Block.Funcs[function.Name] = function
	}

	if pack.Block.Types == nil {
		pack.Block.Types = make(map[string]*ast.VariableType)
	}
	for _, v := range c.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_TYPE_ALIAS) {
		index := binary.BigEndian.Uint64(v.Info)
		name, typ, err := jvm.LucyTypeAliasParser.Decode(c.ConstPool[index].Info)
		if err != nil {
			return err
		}
		pack.Block.Types[name] = typ
	}
	return nil
}

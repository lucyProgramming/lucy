package lc

import (
	"encoding/binary"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

/*
	lucy and java have no difference
*/
func (loader *RealNameLoader) loadInterfaces(astClass *ast.Class, c *cg.Class) error {
	astClass.InterfaceNames = make([]*ast.NameWithPos, len(c.Interfaces))
	for k, v := range c.Interfaces {
		astClass.InterfaceNames[k] = &ast.NameWithPos{
			Name: string(c.ConstPool[v].Info),
		}
	}
	astClass.Interfaces = make([]*ast.Class, len(astClass.InterfaceNames))
	for k, c := range astClass.InterfaceNames {
		i := &ast.Class{}
		i.Name = c.Name
		i.NotImportedYet = true
		astClass.Interfaces[k] = i
	}
	return nil
}

func (loader *RealNameLoader) loadAsJava(c *cg.Class) (*ast.Class, error) {
	//name
	if t := c.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_SIGNATURE); t != nil && len(t) > 0 {
		//TODO:: support signature???
	}
	astClass := &ast.Class{}
	{
		nameIndex := binary.BigEndian.Uint16(c.ConstPool[c.ThisClass].Info)
		astClass.Name = string(c.ConstPool[nameIndex].Info)
		if astClass.Name != ast.JAVA_ROOT_CLASS {
			nameIndex = binary.BigEndian.Uint16(c.ConstPool[c.SuperClass].Info)
			astClass.SuperClassName = string(c.ConstPool[nameIndex].Info)
		}
	}
	loader.loadInterfaces(astClass, c)
	astClass.AccessFlags = c.AccessFlag
	astClass.IsJava = true // class compiled from java
	var err error
	astClass.Fields = make(map[string]*ast.ClassField)
	astClass.LoadFromOutSide = true
	for _, v := range c.Fields {
		f := &ast.ClassField{}
		f.LoadFromOutSide = true
		f.AccessFlags = v.AccessFlags
		f.Descriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		f.Name = string(c.ConstPool[v.NameIndex].Info)
		_, f.Typ, err = jvm.Descriptor.ParseType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
		}
		astClass.Fields[f.Name] = f
	}
	astClass.Methods = make(map[string][]*ast.ClassMethod)
	for _, v := range c.Methods {
		m := &ast.ClassMethod{}
		m.Func = &ast.Function{}
		m.LoadFromOutSide = true
		m.Func.Name = string(c.ConstPool[v.NameIndex].Info)
		m.Func.Descriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		m.Func.AccessFlags = v.AccessFlags
		m.Func.Typ, err = jvm.Descriptor.ParseFunctionType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
		}
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_METHOD_PARAMETERS); t != nil && len(t) > 0 {
			parseMethodParameter(c, t[0].Info, m.Func)
		}
		if t := v.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_RETURNLIST_NAMES); t != nil && len(t) > 0 {
			parseReturnListNames(c, t[0].Info, m.Func)
		}
		if astClass.Methods[m.Func.Name] == nil {
			astClass.Methods[m.Func.Name] = []*ast.ClassMethod{m}
		} else {
			astClass.Methods[m.Func.Name] = append(astClass.Methods[m.Func.Name], m)
		}
	}
	return astClass, nil
}

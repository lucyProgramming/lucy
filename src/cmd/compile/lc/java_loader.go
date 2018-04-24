package lc

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type NotSupportTypeSignatureError struct {
}

func (e *NotSupportTypeSignatureError) Error() string {
	return "lucy does not support typed parameter currently"
}

func (this *RealNameLoader) loadAsJava(c *cg.Class) (*ast.Class, error) {
	//name
	//if t := c.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_SIGNATURE); t != nil && len(t) > 0 {
	//	return nil, &NotSupportTypeSignatureError{}
	//}
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

		if astClass.Methods[m.Func.Name] == nil {
			astClass.Methods[m.Func.Name] = []*ast.ClassMethod{m}
		} else {
			astClass.Methods[m.Func.Name] = append(astClass.Methods[m.Func.Name], m)
		}

	}
	return astClass, nil
}

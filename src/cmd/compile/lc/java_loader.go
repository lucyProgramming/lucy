package lc

import (
	"encoding/binary"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (this *RealNameLoader) loadAsJava(c *cg.Class) (*ast.Class, error) {
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
		f.LoadFromOutSide = true
		f.Descriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		f.Name = string(c.ConstPool[v.NameIndex].Info)
		_, f.Typ, err = jvm.Descriptor.ParseType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
		}
		f.AccessFlags = v.AccessFlags
		astClass.Fields[f.Name] = f
	}
	astClass.Methods = make(map[string][]*ast.ClassMethod)
	for _, v := range c.Methods {
		m := &ast.ClassMethod{}
		m.Func = &ast.Function{}
		m.LoadFromOutSide = true
		m.Func.Name = string(c.ConstPool[v.NameIndex].Info)
		m.Func.Descriptor = string(c.ConstPool[v.DescriptorIndex].Info)
		m.Func.Typ, err = jvm.Descriptor.ParseFunctionType(c.ConstPool[v.DescriptorIndex].Info)
		if err != nil {
			return nil, err
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

package lc

import (
	"encoding/binary"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func (loader *RealNameLoader) loadAsLucy(c *cg.Class) (*ast.Class, error) {

	return nil, nil
}

func (loader *RealNameLoader) loadLucyMainClass(p *ast.Package, c *cg.Class) error {
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
			valueIndex := binary.BigEndian.Uint16(constValue[0].Info)
			switch cos.Typ.Typ {
			case ast.VARIABLE_TYPE_BOOL:
				cos.Data = binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info) != 0
			case ast.VARIABLE_TYPE_BYTE:
				cos.Data = byte(binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_SHORT:
				cos.Data = binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info)
			case ast.VARIABLE_TYPE_INT:
				cos.Data = binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info)
			case ast.VARIABLE_TYPE_LONG:
				cos.Data = int64(binary.BigEndian.Uint64(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_FLOAT:
				cos.Data = float32(binary.BigEndian.Uint32(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_DOUBLE:
				cos.Data = float64(binary.BigEndian.Uint64(c.ConstPool[valueIndex].Info))
			case ast.VARIABLE_TYPE_STRING:
				cos.Data = string(c.ConstPool[valueIndex].Info)
			}
			if loader.Package.Block.Consts == nil {
				loader.Package.Block.Consts = make(map[string]*ast.Const)
			}
			loader.Package.Block.Consts[name] = cos
		} else {
			//global vars
			vd := &ast.VariableDefinition{}
			vd.Name = name
			vd.AccessFlags = f.AccessFlags
			vd.Typ = typ
			if loader.Package.Block.Vars == nil {
				loader.Package.Block.Vars = make(map[string]*ast.VariableDefinition)
			}
			loader.Package.Block.Vars[name] = vd
		}
	}
	var err error
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
		if p.Block.Funcs == nil {
			p.Block.Funcs = make(map[string]*ast.Function)
		}
		p.Block.Funcs[function.Name] = function
	}

	return nil
}

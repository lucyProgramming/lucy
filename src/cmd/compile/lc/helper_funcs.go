package lc

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

func parseMethodParameter(class *cg.Class, bs []byte, f *ast.Function) {
	a := &cg.AttributeMethodParameters{}
	a.FromBs(class, bs)
	for k, v := range a.Parameters {
		f.Typ.ParameterList[k].Name = v.Name
	}
}
func parseReturnListNames(class *cg.Class, bs []byte, f *ast.Function) {
	a := &cg.AttributeMethodParameters{}
	a.FromBs(class, bs)
	for k, v := range a.Parameters {
		f.Typ.ReturnList[k].Name = v.Name
	}
}

func loadEnumForFunction(f *ast.Function) error {
	for _, v := range f.Typ.ParameterList {
		if v.Typ.Typ == ast.VARIABLE_TYPE_ENUM {
			err := loadEnumForVariableType(v.Typ)
			if err != nil {
				return err
			}
		}
	}
	for _, v := range f.Typ.ReturnList {
		if v.Typ.Typ == ast.VARIABLE_TYPE_ENUM {
			err := loadEnumForVariableType(v.Typ)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func loadEnumForVariableType(v *ast.VariableType) error {
	t, err := loader.LoadName(v.Enum.Name)
	if err != nil {
		return err
	}
	if tt, ok := t.(*ast.Enum); ok && tt != nil {
		v.Enum = tt
	} else {
		return fmt.Errorf("'%s' is not a enum", v.Enum.Name)
	}
	return nil
}

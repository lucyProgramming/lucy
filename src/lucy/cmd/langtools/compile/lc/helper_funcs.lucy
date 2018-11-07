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
		f.Type.ParameterList[k].Name = v.Name
	}
}
func parseReturnListNames(class *cg.Class, bs []byte, f *ast.Function) {
	a := &cg.AttributeMethodParameters{}
	a.FromBs(class, bs)
	for k, v := range a.Parameters {
		f.Type.ReturnList[k].Name = v.Name
	}
}

func loadEnumForFunction(f *ast.Function) error {
	for _, v := range f.Type.ParameterList {
		if v.Type.Type == ast.VariableTypeEnum {
			err := loadEnumForVariableType(v.Type)
			if err != nil {
				return err
			}
		}
	}
	for _, v := range f.Type.ReturnList {
		if v.Type.Type == ast.VariableTypeEnum {
			err := loadEnumForVariableType(v.Type)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func loadEnumForVariableType(v *ast.Type) error {
	t, err := loader.LoadImport(v.Enum.Name)
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

package lc

import (
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

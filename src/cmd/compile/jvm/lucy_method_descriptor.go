package jvm

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type LucyMethodSignature struct {
}

func (signature *LucyMethodSignature) Need(ft *ast.FunctionType) bool {
	for _, v := range ft.ParameterList {
		if LucyFieldSignatureParser.Need(v.Type) {
			return true
		}
	}
	for _, v := range ft.ReturnList {
		if LucyFieldSignatureParser.Need(v.Type) {
			return true
		}
	}
	if len(ft.ReturnList) > 1 {
		return true
	}
	return false
}

func (signature *LucyMethodSignature) Encode(ft *ast.FunctionType) (descriptor string) {
	descriptor = "("
	for _, v := range ft.ParameterList {
		descriptor += LucyFieldSignatureParser.Encode(v.Type)
	}
	descriptor += ")"
	if ft.NoReturnValue() {
		descriptor += "V"
	} else {
		descriptor += "("
		for _, v := range ft.ReturnList {
			descriptor += LucyFieldSignatureParser.Encode(v.Type)
		}
		descriptor += ")"
	}
	return descriptor
}

//rewrite types
func (signature *LucyMethodSignature) Decode(ft *ast.FunctionType, bs []byte) ([]byte, error) {
	bs = bs[1:] // skip (
	var err error
	for i := 0; i < len(ft.ParameterList); i++ {
		bs, ft.ParameterList[i].Type, err = LucyFieldSignatureParser.Decode(bs)
		if err != nil {
			return bs, err
		}
	}
	if bs[0] != ')' {
		return bs, fmt.Errorf("function type format wrong")
	}
	bs = bs[1:] // skip )
	if bs[0] == '(' {
		bs = bs[1:]
		ft.ReturnList = []*ast.Variable{}
		for bs[0] != ')' {
			v := &ast.Variable{}
			var t *ast.Type
			bs, t, err = LucyFieldSignatureParser.Decode(bs)
			if err != nil {
				return bs, err
			}
			v.Type = t
			ft.ReturnList = append(ft.ReturnList, v)
		}
		bs = bs[1:] // skip )
	} else if bs[0] == 'V' {
		bs = bs[1:] // skip V
		ft.ReturnList = make([]*ast.Variable, 1)
		ft.ReturnList[0] = &ast.Variable{
			Name: "returnValue",
			Type: &ast.Type{
				Type: ast.VariableTypeVoid,
			},
		}
	} else {
		return bs, fmt.Errorf("function type format wrong")
	}

	return bs, nil
}

package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type LucyMethodSignature struct {
}

func (signature *LucyMethodSignature) Need(functionType *ast.FunctionType) bool {
	for _, v := range functionType.ParameterList {
		if LucyFieldSignatureParser.Need(v.Type) {
			return true
		}
	}
	for _, v := range functionType.ReturnList {
		if LucyFieldSignatureParser.Need(v.Type) {
			return true
		}
	}
	if len(functionType.ReturnList) > 1 {
		return true
	}
	return false
}

func (signature *LucyMethodSignature) Encode(f *ast.Function) (descriptor string) {
	descriptor = "("
	for _, v := range f.Type.ParameterList {
		descriptor += LucyFieldSignatureParser.Encode(v.Type)
	}
	descriptor += ")"
	if f.NoReturnValue() {
		descriptor += "V"
	} else {
		for _, v := range f.Type.ReturnList {
			descriptor += LucyFieldSignatureParser.Encode(v.Type)
		}
	}
	return descriptor
}

//rewrite types
func (signature *LucyMethodSignature) Decode(f *ast.Function, bs []byte) error {
	bs = bs[1:] // skip (
	var err error
	for i := 0; i < len(f.Type.ParameterList); i++ {
		bs, f.Type.ParameterList[i].Type, err = LucyFieldSignatureParser.Decode(bs)
		if err != nil {
			return err
		}
	}
	bs = bs[1:] // skip )
	f.Type.ReturnList = []*ast.VariableDefinition{}
	i := 1
	for len(bs) > 0 {
		var t *ast.VariableType
		bs, t, err = LucyFieldSignatureParser.Decode(bs)
		if err != nil {
			return err
		}
		vd := &ast.VariableDefinition{}
		vd.Type = t
		f.Type.ReturnList = append(f.Type.ReturnList, vd)
		i++
	}

	return nil
}

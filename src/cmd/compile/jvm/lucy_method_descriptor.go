package jvm

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type LucyMethodSignatureParse struct {
}

func (parser *LucyMethodSignatureParse) Need(functionType *ast.FunctionType) bool {
	for _, v := range functionType.ParameterList {
		if LucyFieldSignatureParser.Need(v.Typ) {
			return true
		}
	}
	for _, v := range functionType.ReturnList {
		if LucyFieldSignatureParser.Need(v.Typ) {
			return true
		}
	}
	if len(functionType.ReturnList) > 1 {
		return true
	}
	return false
}

func (parser *LucyMethodSignatureParse) Encode(f *ast.Function) (descriptor string) {
	descriptor = "("
	for _, v := range f.Typ.ParameterList {
		descriptor += LucyFieldSignatureParser.Encode(v.Typ)
	}
	descriptor += ")"
	if f.NoReturnValue() {
		descriptor += "V"
	} else {
		for _, v := range f.Typ.ReturnList {
			descriptor += LucyFieldSignatureParser.Encode(v.Typ)
		}
	}
	return descriptor
}

//rewrite types
func (parser *LucyMethodSignatureParse) Decode(f *ast.Function, bs []byte) error {
	bs = bs[1:] // skip (
	var err error
	for i := 0; i < len(f.Typ.ParameterList); i++ {
		bs, f.Typ.ParameterList[i].Typ, err = LucyFieldSignatureParser.Decode(bs)
		if err != nil {
			return err
		}
	}
	bs = bs[1:] // skip )
	f.Typ.ReturnList = []*ast.VariableDefinition{}
	i := 1
	for len(bs) > 0 {
		var t *ast.VariableType
		bs, t, err = LucyFieldSignatureParser.Decode(bs)
		if err != nil {
			return err
		}
		vd := &ast.VariableDefinition{}
		vd.Typ = t
		f.Typ.ReturnList = append(f.Typ.ReturnList, vd)
		i++
	}

	return nil
}

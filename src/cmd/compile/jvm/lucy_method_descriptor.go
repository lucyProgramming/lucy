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
func (parser *LucyMethodSignatureParse) Deocde(bs []byte, f *ast.Function) error {
	bs = bs[1:] // skip (
	var err error
	for i := 0; i < len(f.Typ.ParameterList); i++ {
		bs, f.Typ.ParameterList[i].Typ, err = LucyFieldSignatureParser.Decode(bs)
		if err != nil {
			return err
		}
	}
	bs = bs[1:] // skip )
	for i := 0; i < len(f.Typ.ReturnList); i++ {
		bs, f.Typ.ReturnList[i].Typ, err = LucyFieldSignatureParser.Decode(bs)
		if err != nil {
			return err
		}
	}
	return nil
}

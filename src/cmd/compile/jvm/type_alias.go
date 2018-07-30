package jvm

import (
	"bytes"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type LucyTypeAlias struct {
}

func (LucyTypeAlias) Encode(name string, variableType *ast.Type) string {
	name += ";"
	name += LucyFieldSignatureParser.Encode(variableType)
	return name
}

func (LucyTypeAlias) Decode(bs []byte) (name string, variableType *ast.Type, err error) {
	index := bytes.Index(bs, []byte{';'})
	name = string(bs[0:index])
	_, variableType, err = LucyFieldSignatureParser.Decode(bs[index+1:])
	return
}

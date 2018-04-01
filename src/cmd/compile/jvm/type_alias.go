package jvm

import (
	"bytes"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type LucyTypeAliasParse struct {
}

func (l *LucyTypeAliasParse) Encode(name string, variableType *ast.VariableType) string {
	name += ";"
	name += LucyFieldSignatureParser.Encode(variableType)
	return name
}

func (l *LucyTypeAliasParse) Decode(bs []byte) (name string, variableType *ast.VariableType, err error) {
	index := bytes.Index(bs, []byte{';'})
	name = string(bs[0:index])
	_, variableType, err = LucyFieldSignatureParser.Decode(bs[index+1:])
	return
}

package signature

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
)

type LucyFunctionSignatureParser struct {
}

func (LucyFunctionSignatureParser) Parse(bs []byte) (*ast.FunctionType, error) {
	if bs[0] != '(' {
		return nil, fmt.Errorf("signature dose not beging with '('")
	}
	bs = bs[1:]
	ret := &ast.FunctionType{}
	i := 1
	for bs[0] != ')' {
		switch bs[0] {
		case 'B':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_BYTE
			bs = bs[1:]
		case 'C':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_SHORT
			bs = bs[1:]
		case 'D':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_DOUBLE
			bs = bs[1:]
		case 'F':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_FLOAT
			bs = bs[1:]
		case 'I':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_INT
			bs = bs[1:]
		case 'J':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_LONG
			bs = bs[1:]
		case 'S':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_SHORT
			bs = bs[1:]
		case 'Z':
			vd := &ast.VariableDefinition{}
			vd.Name = fmt.Sprintf("var_%d", i)
			vd.Typ = &ast.VariableType{}
			vd.Typ.Typ = ast.VARIABLE_TYPE_BOOL
			bs = bs[1:]
		case '[':
		case 'L':
		default:

		}

		i++
	}
}

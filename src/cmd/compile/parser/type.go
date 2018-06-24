package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (parser *Parser) parseType(pre ...*ast.Type) (*ast.Type, error) {
	var err error
	var ret *ast.Type
	switch parser.token.Type {
	case lex.TOKEN_LB:
		pos := parser.mkPos()
		if len(pre) > 0 {
			parser.Next()
			if parser.token.Type != lex.TOKEN_RB {
				// [ and ] not match
				err = fmt.Errorf("%s '[' and ']' not match", parser.errorMsgPrefix())
				parser.errs = append(parser.errs, err)
				return nil, err
			}
			parser.Next() //skip ]
			ret = &ast.Type{
				Pos:       pos,
				Type:      ast.VARIABLE_TYPE_JAVA_ARRAY,
				ArrayType: pre[0],
			}
			break
		}
		parser.Next()
		if parser.token.Type != lex.TOKEN_RB {
			// [ and ] not match
			err = fmt.Errorf("%s '[' and ']' not match", parser.errorMsgPrefix())
			parser.errs = append(parser.errs, err)
			return nil, err
		}
		//lookahead
		parser.Next() //skip ]
		t, err := parser.parseType()
		if err != nil {
			return nil, err
		}
		ret = &ast.Type{}
		ret.Pos = pos
		ret.Type = ast.VARIABLE_TYPE_ARRAY
		ret.ArrayType = t
	case lex.TOKEN_BOOL:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VARIABLE_TYPE_BOOL,
			Pos:  pos,
		}
	case lex.TOKEN_BYTE:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VARIABLE_TYPE_BYTE,
			Pos:  pos,
		}
	case lex.TOKEN_SHORT:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VARIABLE_TYPE_SHORT,
			Pos:  pos,
		}
	case lex.TOKEN_INT:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VARIABLE_TYPE_INT,
			Pos:  pos,
		}
	case lex.TOKEN_FLOAT:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VARIABLE_TYPE_FLOAT,
			Pos:  pos,
		}
	case lex.TOKEN_DOUBLE:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VARIABLE_TYPE_DOUBLE,
			Pos:  pos,
		}
	case lex.TOKEN_LONG:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VARIABLE_TYPE_LONG,
			Pos:  pos,
		}
	case lex.TOKEN_STRING:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VARIABLE_TYPE_STRING,
			Pos:  pos,
		}
	case lex.TOKEN_IDENTIFIER:
		ret, err = parser.parseIdentifierType()

	case lex.TOKEN_MAP:
		pos := parser.mkPos()
		parser.Next() // skip map key word
		if parser.token.Type != lex.TOKEN_LC {
			return nil, fmt.Errorf("%s expect '{',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description)
		}
		parser.Next() // skip {
		var k, v *ast.Type
		k, err = parser.parseType()
		if err != nil {
			return nil, err
		}
		if parser.token.Type != lex.TOKEN_ARROW {
			return nil, fmt.Errorf("%s expect '->',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description)
		}
		parser.Next() // skip ->
		v, err := parser.parseType()
		if err != nil {
			return nil, err
		}
		if parser.token.Type != lex.TOKEN_RC {
			return nil, fmt.Errorf("%s expect '}',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description)
		}
		parser.Next()
		m := &ast.Map{
			Key:   k,
			Value: v,
		}
		ret = &ast.Type{
			Type: ast.VARIABLE_TYPE_MAP,
			Map:  m,
			Pos:  pos,
		}
	case lex.TOKEN_T:
		pos := parser.mkPos()
		ret = &ast.Type{
			Type: ast.VARIABLE_TYPE_T,
			Pos:  pos,
			Name: parser.token.Data.(string),
		}
		parser.Next()
	case lex.TOKEN_FUNCTION:
		pos := parser.mkPos()
		ft, err := parser.parseFunctionType()
		if err != nil {
			return nil, err
		}
		ret = &ast.Type{
			Type:         ast.VARIABLE_TYPE_FUNCTION_POINTER,
			Pos:          pos,
			FunctionType: &ft,
		}
	default:
		err = fmt.Errorf("%s unkown type,begining token is '%s'",
			parser.errorMsgPrefix(), parser.token.Description)
	}
	if err != nil {
		parser.errs = append(parser.errs, err)
		return nil, err
	}
	if parser.token.Type == lex.TOKEN_LB {
		return parser.parseType(ret)
	} else {
		return ret, err
	}
}

func (parser *Parser) isValidTypeBegin() bool {
	return parser.token.Type == lex.TOKEN_LB ||
		parser.token.Type == lex.TOKEN_BOOL ||
		parser.token.Type == lex.TOKEN_BYTE ||
		parser.token.Type == lex.TOKEN_SHORT ||
		parser.token.Type == lex.TOKEN_INT ||
		parser.token.Type == lex.TOKEN_FLOAT ||
		parser.token.Type == lex.TOKEN_DOUBLE ||
		parser.token.Type == lex.TOKEN_LONG ||
		parser.token.Type == lex.TOKEN_STRING ||
		parser.token.Type == lex.TOKEN_MAP ||
		parser.token.Type == lex.TOKEN_IDENTIFIER ||
		parser.token.Type == lex.TOKEN_T

}
func (parser *Parser) parseIdentifierType() (*ast.Type, error) {
	name := parser.token.Data.(string)
	ret := &ast.Type{
		Pos:  parser.mkPos(),
		Type: ast.VARIABLE_TYPE_NAME,
	}
	parser.Next() // skip name identifier
	for parser.token.Type == lex.TOKEN_DOT {
		parser.Next() // skip .
		if parser.token.Type != lex.TOKEN_IDENTIFIER {
			return nil, fmt.Errorf("%s not a identifier after dot",
				parser.errorMsgPrefix())
		}
		name += "." + parser.token.Data.(string)
		parser.Next() // if
	}
	ret.Name = name
	return ret, nil
}

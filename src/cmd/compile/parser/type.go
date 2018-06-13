package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (p *Parser) parseType(pre ...*ast.VariableType) (*ast.VariableType, error) {
	var err error
	var ret *ast.VariableType
	switch p.token.Type {
	case lex.TOKEN_LB:
		pos := p.mkPos()
		if len(pre) > 0 {
			p.Next()
			if p.token.Type != lex.TOKEN_RB {
				// [ and ] not match
				err = fmt.Errorf("%s '[' and ']' not match", p.errorMsgPrefix())
				p.errs = append(p.errs, err)
				return nil, err
			}
			p.Next() //skip ]
			ret = &ast.VariableType{
				Pos:       pos,
				Typ:       ast.VARIABLE_TYPE_JAVA_ARRAY,
				ArrayType: pre[0],
			}
			break
		}
		p.Next()
		if p.token.Type != lex.TOKEN_RB {
			// [ and ] not match
			err = fmt.Errorf("%s '[' and ']' not match", p.errorMsgPrefix())
			p.errs = append(p.errs, err)
			return nil, err
		}
		//lookahead
		p.Next() //skip ]
		t, err := p.parseType()
		if err != nil {
			return nil, err
		}
		ret = &ast.VariableType{}
		ret.Pos = pos
		ret.Typ = ast.VARIABLE_TYPE_ARRAY
		ret.ArrayType = t
	case lex.TOKEN_BOOL:
		pos := p.mkPos()
		p.Next()
		ret = &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BOOL,
			Pos: pos,
		}
	case lex.TOKEN_BYTE:
		pos := p.mkPos()
		p.Next()
		ret = &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BYTE,
			Pos: pos,
		}
	case lex.TOKEN_SHORT:
		pos := p.mkPos()
		p.Next()
		ret = &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_SHORT,
			Pos: pos,
		}
	case lex.TOKEN_INT:
		pos := p.mkPos()
		p.Next()
		ret = &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_INT,
			Pos: pos,
		}
	case lex.TOKEN_FLOAT:
		pos := p.mkPos()
		p.Next()
		ret = &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_FLOAT,
			Pos: pos,
		}
	case lex.TOKEN_DOUBLE:
		pos := p.mkPos()
		p.Next()
		ret = &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_DOUBLE,
			Pos: pos,
		}
	case lex.TOKEN_LONG:
		pos := p.mkPos()
		p.Next()
		ret = &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_LONG,
			Pos: pos,
		}
	case lex.TOKEN_STRING:
		pos := p.mkPos()
		p.Next()
		ret = &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_STRING,
			Pos: pos,
		}
	case lex.TOKEN_IDENTIFIER:
		ret, err = p.parseIdentifierType()

	case lex.TOKEN_MAP:
		pos := p.mkPos()
		p.Next() // skip map key word
		if p.token.Type != lex.TOKEN_LC {
			return nil, fmt.Errorf("%s expect '{',but '%s'",
				p.errorMsgPrefix(), p.token.Description)
		}
		p.Next() // skip {
		var k, v *ast.VariableType
		k, err = p.parseType()
		if err != nil {
			return nil, err
		}
		if p.token.Type != lex.TOKEN_ARROW {
			return nil, fmt.Errorf("%s expect '->',but '%s'",
				p.errorMsgPrefix(), p.token.Description)
		}
		p.Next() // skip ->
		v, err := p.parseType()
		if err != nil {
			return nil, err
		}
		if p.token.Type != lex.TOKEN_RC {
			return nil, fmt.Errorf("%s expect '}',but '%s'",
				p.errorMsgPrefix(), p.token.Description)
		}
		p.Next()
		m := &ast.Map{
			K: k,
			V: v,
		}
		ret = &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_MAP,
			Map: m,
			Pos: pos,
		}
	case lex.TOKEN_T:
		pos := p.mkPos()
		ret = &ast.VariableType{
			Typ:  ast.VARIABLE_TYPE_T,
			Pos:  pos,
			Name: p.token.Data.(string),
		}
		p.Next()
	default:
		err = fmt.Errorf("%s unkown type,begining token is '%s'",
			p.errorMsgPrefix(), p.token.Description)
	}
	if err != nil {
		p.errs = append(p.errs, err)
		return nil, err
	}
	if p.token.Type == lex.TOKEN_LB {
		return p.parseType(ret)
	} else {
		return ret, err
	}

}

func (p *Parser) isValidTypeBegin() bool {
	return p.token.Type == lex.TOKEN_LB ||
		p.token.Type == lex.TOKEN_BOOL ||
		p.token.Type == lex.TOKEN_BYTE ||
		p.token.Type == lex.TOKEN_SHORT ||
		p.token.Type == lex.TOKEN_INT ||
		p.token.Type == lex.TOKEN_FLOAT ||
		p.token.Type == lex.TOKEN_DOUBLE ||
		p.token.Type == lex.TOKEN_LONG ||
		p.token.Type == lex.TOKEN_STRING ||
		p.token.Type == lex.TOKEN_MAP ||
		p.token.Type == lex.TOKEN_IDENTIFIER ||
		p.token.Type == lex.TOKEN_T

}
func (p *Parser) parseIdentifierType() (*ast.VariableType, error) {
	name := p.token.Data.(string)
	ret := &ast.VariableType{
		Pos: p.mkPos(),
		Typ: ast.VARIABLE_TYPE_NAME,
	}
	p.Next() // skip name identifier
	for p.token.Type == lex.TOKEN_DOT && p.token.Type != lex.TOKEN_EOF {
		p.Next() // skip .
		if p.token.Type != lex.TOKEN_IDENTIFIER {
			return nil, fmt.Errorf("%s not a identifier after dot",
				p.errorMsgPrefix())
		}
		name += "." + p.token.Data.(string)
		p.Next() // if
	}
	ret.Name = name
	return ret, nil
}

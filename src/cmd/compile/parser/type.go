package parser

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

func (p *Parser) parseType() (*ast.VariableType, error) {
	var err error
	switch p.token.Type {
	case lex.TOKEN_LB:
		pos := p.mkPos()
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
		tt := &ast.VariableType{}
		tt.Pos = pos
		tt.Typ = ast.VARIABLE_TYPE_ARRAY
		tt.CombinationType = t
		return tt, nil
	case lex.TOKEN_BOOL:
		pos := p.mkPos()
		p.Next()
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BOOL,
			Pos: pos,
		}, nil
	case lex.TOKEN_BYTE:
		pos := p.mkPos()
		p.Next()
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BYTE,
			Pos: pos,
		}, nil
	case lex.TOKEN_SHORT:
		pos := p.mkPos()
		p.Next()
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_SHORT,
			Pos: pos,
		}, nil
	case lex.TOKEN_INT:
		pos := p.mkPos()
		p.Next()
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_INT,
			Pos: pos,
		}, nil
	case lex.TOKEN_FLOAT:
		pos := p.mkPos()
		p.Next()
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_FLOAT,
			Pos: pos,
		}, nil
	case lex.TOKEN_DOUBLE:
		pos := p.mkPos()
		p.Next()
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_DOUBLE,
			Pos: pos,
		}, nil
	case lex.TOKEN_LONG:
		pos := p.mkPos()
		p.Next()
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_LONG,
			Pos: pos,
		}, nil
	case lex.TOKEN_STRING:
		pos := p.mkPos()
		p.Next()
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_STRING,
			Pos: pos,
		}, nil
	case lex.TOKEN_IDENTIFIER:
		return p.parseIdentifierType()
	case lex.TOKEN_MAP:
		pos := p.mkPos()
		p.Next() // skip map key word
		if p.token.Type != lex.TOKEN_LB {
			return nil, fmt.Errorf("%s expect '[',but '%s'", p.errorMsgPrefix(), p.token.Desp)
		}
		p.Next() // skip [
		k, err := p.parseType()
		if err != nil {
			return nil, err
		}
		if p.token.Type != lex.TOKEN_ARROW {
			return nil, fmt.Errorf("%s expect '->',but '%s'", p.errorMsgPrefix(), p.token.Desp)
		}
		p.Next() // skip ->
		v, err := p.parseType()
		if err != nil {
			return nil, err
		}
		if p.token.Type != lex.TOKEN_RB {
			return nil, fmt.Errorf("%s expect ']',but '%s'", p.errorMsgPrefix(), p.token.Desp)
		}
		p.Next()
		m := &ast.Map{
			K: k,
			V: v,
		}
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_MAP,
			Map: m,
			Pos: pos,
		}, nil
	}
	err = fmt.Errorf("%s unkown type,first token:%s", p.errorMsgPrefix(), p.token.Desp)
	p.errs = append(p.errs, err)
	return nil, err
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
		p.token.Type == lex.TOKEN_IDENTIFIER
}
func (p *Parser) parseIdentifierType() (*ast.VariableType, error) {
	name := p.token.Data.(string)
	ret := &ast.VariableType{
		Pos: p.mkPos(),
		Typ: ast.VARIABLE_TYPE_NAME,
	}
	p.Next() // skip name identifier
	for p.token.Type == lex.TOKEN_DOT && !p.eof {
		p.Next() // skip .
		if p.token.Type != lex.TOKEN_IDENTIFIER {
			return nil, fmt.Errorf("%s not a identifier after dot", p.errorMsgPrefix())
		}
		name += "." + p.token.Data.(string)
		p.Next() // if
	}
	ret.Name = name
	return ret, nil
}

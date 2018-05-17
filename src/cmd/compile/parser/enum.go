package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (p *Parser) parseEnum(ispublic bool) (e *ast.Enum, err error) {
	p.Next() // skip enum

	if p.token.Type != lex.TOKEN_IDENTIFIER {
		err = fmt.Errorf("%s expect 'identifier', but '%s'",
			p.errorMsgPrefix(), p.token.Desp)
		p.errs = append(p.errs, err)
		return nil, err
	}
	enumName := &ast.NameWithPos{
		Name: p.token.Data.(string),
		Pos:  p.mkPos(),
	}
	p.Next() // skip enum name
	if p.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s expect '{',but '%s'",
			p.errorMsgPrefix(), p.token.Desp)
		p.errs = append(p.errs, err)
		return nil, err
	}
	p.Next() // skip {
	e = &ast.Enum{}
	e.Name = enumName.Name
	e.Pos = enumName.Pos
	//first name
	if p.token.Type != lex.TOKEN_IDENTIFIER {
		err = fmt.Errorf("%s expect 'identifier',but '%s'",
			p.errorMsgPrefix(), p.token.Desp)
		p.errs = append(p.errs, err)
		return nil, err
	}
	names := []*ast.NameWithPos{
		&ast.NameWithPos{
			Name: p.token.Data.(string),
			Pos:  p.mkPos(),
		},
	}
	p.Next()
	var initExpression *ast.Expression
	if p.token.Type == lex.TOKEN_ASSIGN { // first value defined here
		p.Next() // skip assign
		initExpression, err = p.Expression.parseExpression(false)
		if err != nil {
			p.errs = append(p.errs, err)
			return nil, err
		}
	}
	if p.token.Type == lex.TOKEN_COMMA {
		p.Next() // skip ,should be a identifier after  commna
		ns, err := p.parseNameList()
		if err != nil {
			return nil, err
		}
		names = append(names, ns...)
	}
	if p.token.Type != lex.TOKEN_RC {
		err = fmt.Errorf("%s expect '}',but '%s'", p.token.Desp, p.token.Desp)
		p.errs = append(p.errs, err)
		return nil, err
	}
	p.Next() // skip }
	e.Init = initExpression
	for _, v := range names {
		t := &ast.EnumName{}
		t.Name = v.Name
		t.Pos = v.Pos
		t.Enum = e
		e.Enums = append(e.Enums, t)
	}
	e.AccessFlags = 0
	if ispublic {
		e.AccessFlags |= cg.ACC_CLASS_PUBLIC
	}
	return e, err
}

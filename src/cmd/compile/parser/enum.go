package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (p *Parser) parseEnum(ispublic bool) (e *ast.Enum, err error) {
	p.Next() // skip enum
	if p.eof {
		err = p.mkUnexpectedEofErr()
		p.errs = append(p.errs, err)
		return nil, err
	}

	if p.token.Type != lex.TOKEN_IDENTIFIER {
		err = fmt.Errorf("%s enum type have no name", p.errorMsgPrefix())
		p.errs = append(p.errs, err)
		return nil, err
	}
	enumName := &ast.NameWithPos{
		Name: p.token.Data.(string),
		Pos:  p.mkPos(),
	}
	p.Next() // skip enum name
	if p.eof {
		err = p.mkUnexpectedEofErr()
		p.errs = append(p.errs, err)
		return nil, err
	}
	if p.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s enum type have no '{' after it`s name defined,but %s", p.errorMsgPrefix(), p.token.Desp)
		p.errs = append(p.errs, err)
		return nil, err
	}
	p.Next() // skip {
	if p.eof {
		err = p.mkUnexpectedEofErr()
		return nil, err
	}
	e = &ast.Enum{}
	e.Name = enumName.Name
	e.Pos = enumName.Pos
	e.NamesMap = make(map[string]*ast.EnumName)
	if p.token.Type == lex.TOKEN_RC {
		return e, nil
	}
	//first name
	if p.token.Type != lex.TOKEN_IDENTIFIER {
		err = fmt.Errorf("%s no enum names defined after {", p.errorMsgPrefix())
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
	if p.eof {
		err = p.mkUnexpectedEofErr()
		p.errs = append(p.errs, err)
		return nil, err
	}
	var initExpression *ast.Expression
	if p.token.Type == lex.TOKEN_ASSIGN { // first value defined here
		p.Next()
		if p.eof {
			err = p.mkUnexpectedEofErr()
			p.errs = append(p.errs, err)
			return nil, err
		}
		initExpression, err = p.ExpressionParser.parseExpression(true)
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
		err = fmt.Errorf("%s unexcept token(%s) after a enum name", p.token.Desp)
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
		e.Names = append(e.Names, t)
		if e.NamesMap[v.Name] != nil {
			p.errs = append(p.errs, fmt.Errorf("%s enumname %s already declared", p.errorMsgPrefix(v.Pos)))
		} else {
			e.NamesMap[v.Name] = t
		}
	}
	e.Access = 0
	if ispublic {
		e.Access |= cg.ACC_CLASS_PUBLIC
	}
	return e, err
}

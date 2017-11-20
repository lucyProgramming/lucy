package parser

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
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
		initExpression, err = p.ExpressionParser.parseExpression()
		if err != nil {
			p.errs = append(p.errs, err)
			return nil, err
		}
	}
	if p.token.Type == lex.TOKEN_COMMA {
		p.Next()
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
	e = &ast.Enum{}
	e.Name = enumName.Name
	e.Init = initExpression
	e.Pos = enumName.Pos
	for _, v := range names {
		t := &ast.EnumName{}
		t.Name = v.Name
		t.Pos = v.Pos
		t.Enum = e
	}
	if ispublic {
		e.Access = ast.ACCESS_PUBLIC
	} else {
		e.Access = ast.ACCESS_PRIVATE
	}
	return e, err
}

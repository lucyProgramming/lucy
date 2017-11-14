package parser

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

func (p *Parser) parseEnum(ispublic bool) (e *ast.Enum) {
	p.Next()
	if p.eof {
		p.unexpectedErr()
		return
	}
	if p.token.Type != lex.TOKEN_IDENTIFIER {
		p.errs = append(p.errs, fmt.Errorf("%s enum type have no name", p.errorMsgPrefix()))
		p.consume(lex.TOKEN_RC)
		p.Next()
		return
	}
	enumName := &ast.NameWithPos{
		Name: p.token.Data.(string),
		Pos:  p.mkPos(),
	}
	p.Next()
	if p.eof {
		p.unexpectedErr()
		return
	}
	if p.token.Type != lex.TOKEN_LC {
		p.errs = append(p.errs, fmt.Errorf("%s enum type have no '{' after it`s name defined,but %s", p.errorMsgPrefix(), p.token.Desp))
		p.consume(lex.TOKEN_RC)
		p.Next()
		return
	}
	p.Next()
	if p.eof {
		p.unexpectedErr()
		return
	}
	//first name
	if p.token.Type != lex.TOKEN_IDENTIFIER {
		p.errs = append(p.errs, fmt.Errorf("%s no enum names defined after {", p.errorMsgPrefix()))
		p.consume(lex.TOKEN_RC)
		p.Next()
		return
	}
	names := []*ast.NameWithPos{
		&ast.NameWithPos{
			Name: p.token.Data.(string),
			Pos:  p.mkPos(),
		},
	}
	p.Next()
	if p.eof {
		p.unexpectedErr()
		return
	}
	var initExpression *ast.Expression
	var err error
	if p.token.Type == lex.TOKEN_ASSIGN { // first value defined here
		p.Next()
		if p.eof {
			return
		}
		initExpression, err = p.ExpressionParser.parseExpression()
		if err != nil {
			p.errs = append(p.errs, err)
			p.consume(lex.TOKEN_RC)
			p.Next()
			return
		}
	}
	if p.token.Type == lex.TOKEN_COMMA {
		p.Next()
		ns, err := p.parseNameList()
		if err != nil {
			p.errs = append(p.errs, err)
			p.consume(lex.TOKEN_RC)
			p.Next()
			return
		}
		if p.token.Type != lex.TOKEN_RC {
			p.errs = append(p.errs, fmt.Errorf("%s except } after name list,but %s", p.errorMsgPrefix(), p.token.Desp))
			p.consume(lex.TOKEN_RC)
			p.Next()
			return
		}
		names = append(names, ns...)
	} else if p.token.Type == lex.TOKEN_RC { //  enmu define ended
	} else {
		p.errs = append(p.errs, fmt.Errorf("%s unexcept token(%s) after a enum name", p.token.Desp))
		p.consume(lex.TOKEN_RC)
		p.Next()
		return
	}
	p.Next()
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
		e.AccessProperty.Access = ast.ACCESS_PUBLIC
	} else {
		e.AccessProperty.Access = ast.ACCESS_PRIVATE
	}
	return
}

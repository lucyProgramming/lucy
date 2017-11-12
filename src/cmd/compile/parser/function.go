package parser

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
	"github.com/open-falcon/nodata/g"
)

func (p *Parser) parseFunction(ispublic bool) (f *ast.Function, err error) {
	p.Next()
	if p.eof {
		return nil, p.mkUnexpectedErr()
	}
	f = &ast.Function{}
	if p.token.Type == lex.TOKEN_LP {
		f.Typ.Returns, err = p.parseTypedNames()
		if err != nil {
			p.consume(lex.TOKEN_RC)
			p.Next()
			return
		}
		if p.token.Type != lex.TOKEN_RP {
			p.consume(lex.TOKEN_RC)
			p.Next()
			err = fmt.Errorf("%s except ) but %s", p.errorMsgPrefix(), p.token.Desp)
			return
		}
	}
	if p.token.Type == lex.TOKEN_IDENTIFIER {
		f.Name = p.token.Data.(string)
		p.Next()
	}
	if p.token.Type != lex.TOKEN_LC {
		return nil, fmt.Errorf("%s except { but %s", p.token.Desp)
	}
	p.Next() //
	if p.eof {
		return nil, p.mkUnexpectedErr()
	}

	f.Block = &ast.Block{}
	f.Block.Statements = p.parseStatementList()
	p.Next()
	return f, nil
}
func (p *Parser) parseStatementList(ispublic bool) []*ast.Statement {
	ret := []*ast.Statement{}
	for {
		switch p.token {
		case lex.TOKEN_RC:
			return ret
		case lex.TOKEN_IDENTIFIER:
			e, err := p.ExpressionParser.parseExpression(false)
			if err != nil {
				p.errs = append(p.errs)
			}
			ret = append(ret, &ast.Statement{
				Typ:        ast.STATEMENT_TYPE_EXPRESSION,
				Expression: e,
			})
		case lex.TOKEN_IF:
		case lex.TOKEN_FOR:
		case lex.TOKEN_SWITCH:
		default:
			p.errs = append(p.errs, fmt.Errorf("%s unkown begining of a statement", p.errorMsgPrefix()))
			p.consume(lex.TOKEN_SEMICOLON)
			p.Next()
		}
	}
}

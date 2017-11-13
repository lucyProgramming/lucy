package parser

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

type Function struct {
	parser *Parser
}

func (p *Function) Next() {
	p.parser.Next()
}

func (p *Function) consume(untils ...int) {
	p.parser.consume(untils...)
}

func (p *Function) parse(ispublic bool) (f *ast.Function, err error) {
	p.Next()
	if p.parser.eof {
		return nil, p.parser.mkUnexpectedEofErr()
	}
	f = &ast.Function{}
	if p.parser.token.Type == lex.TOKEN_LP {
		f.Typ.Returns, err = p.parser.parseTypedNames()
		if err != nil {
			p.parser.consume(lex.TOKEN_RC)
			p.Next()
			return
		}
		if p.parser.token.Type != lex.TOKEN_RP {
			p.parser.consume(lex.TOKEN_RC)
			p.Next()
			err = fmt.Errorf("%s except ) but %s", p.parser.errorMsgPrefix(), p.parser.token.Desp)
			return
		}
	}
	if p.parser.token.Type == lex.TOKEN_IDENTIFIER {
		f.Name = p.parser.token.Data.(string)
		p.Next()
	}
	if p.parser.token.Type != lex.TOKEN_LC {
		return nil, fmt.Errorf("%s except { but %s", p.parser.token.Desp)
	}
	p.Next() //
	if p.parser.eof {
		return nil, p.parser.mkUnexpectedEofErr()
	}
	if ispublic {
		f.Access = ast.ACCESS_PUBLIC
	} else {
		f.Access = ast.ACCESS_PRIVATE
	}
	f.Block = &ast.Block{}
	f.Block.Statements = p.parseStatementList(f.Block)
	p.Next()
	return f, nil
}
func (p *Function) parseStatementList(block *ast.Block) []*ast.Statement {
	ret := []*ast.Statement{}
	for {
		switch p.parser.token.Type {
		case lex.TOKEN_RC: // end
			return ret
		case lex.TOKEN_IDENTIFIER:
			e, err := p.parser.ExpressionParser.parseExpression(false)
			if err != nil {
				p.parser.errs = append(p.parser.errs, err)
			}
			ret = append(ret, &ast.Statement{
				Typ:        ast.STATEMENT_TYPE_EXPRESSION,
				Expression: e,
			})
		case lex.TOKEN_VAR:
			errs := p.parser.insertVariableIntoBlock(block, p.parser.parseVarDefinition())
			if errs != nil {
				p.parser.errs = append(p.parser.errs, errs...)
			}
		case lex.TOKEN_IF:
		case lex.TOKEN_FOR:
		case lex.TOKEN_SWITCH:
		default:
			p.parser.errs = append(p.parser.errs, fmt.Errorf("%s unkown begining of a statement", p.parser.errorMsgPrefix()))
			p.consume(lex.TOKEN_SEMICOLON)
			p.Next()
		}
	}
}

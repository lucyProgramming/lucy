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

func (p *Function) consume(untils map[int]bool) {
	p.parser.consume(untils)
}

func (p *Function) parse(ispublic bool) (f *ast.Function, err error) {
	p.Next()
	if p.parser.eof {
		return nil, p.parser.mkUnexpectedEofErr()
	}
	f = &ast.Function{}
	f.Pos = p.parser.mkPos()
	if p.parser.token.Type == lex.TOKEN_IDENTIFIER {
		f.Name = p.parser.token.Data.(string)
		p.Next()
	}
	if p.parser.token.Type != lex.TOKEN_LP {
		err = fmt.Errorf("%s fn declared wrong,missing (", p.parser.errorMsgPrefix())
		p.parser.errs = append(p.parser.errs, err)
		return
	}
	p.Next()
	if p.parser.eof {
		return nil, p.parser.mkUnexpectedEofErr()
	}
	if p.parser.token.Type != lex.TOKEN_RP { // not (
		f.Typ.Parameters, err = p.parser.parseTypedNames()
		if err != nil {
			return
		}
	}
	if p.parser.token.Type != lex.TOKEN_RP { // not (
		err = fmt.Errorf("%s fn declared wrong,missing ),but %s", p.parser.errorMsgPrefix(), p.parser.token.Desp)
		p.parser.errs = append(p.parser.errs, err)
		return
	}
	p.Next()
	if p.parser.token.Type == lex.TOKEN_ARROW { // ->
		p.Next()
		if p.parser.token.Type != lex.TOKEN_LP {
			err = fmt.Errorf("%s fn declared wrong, not ( after ->", p.parser.errorMsgPrefix())
			p.parser.errs = append(p.parser.errs, err)
			return
		}
		p.Next()
		f.Typ.Returns, err = p.parser.parseTypedNames()
		if err != nil {
			return
		}
		if p.parser.token.Type != lex.TOKEN_RP {
			err = fmt.Errorf("%s fn declared wrong, ( and ) not match", p.parser.errorMsgPrefix())
			p.parser.errs = append(p.parser.errs, err)
			return
		}
		p.Next()
	}
	if p.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s except { but %s", p.parser.errorMsgPrefix(), p.parser.token.Desp)
		p.parser.errs = append(p.parser.errs, err)
		return
	}
	if ispublic {
		f.Access = ast.ACCESS_PUBLIC
	} else {
		f.Access = ast.ACCESS_PRIVATE
	}
	f.Block = &ast.Block{}
	err = p.parser.Block.parse(f.Block)
	return f, err
}

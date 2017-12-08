package parser

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
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
	p.Next() // skip fn key word
	if p.parser.eof {
		err = p.parser.mkUnexpectedEofErr()
		p.parser.errs = append(p.parser.errs, err)
		return nil, err
	}
	f = &ast.Function{}
	f.Pos = p.parser.mkPos()
	if p.parser.token.Type == lex.TOKEN_IDENTIFIER {
		f.Name = p.parser.token.Data.(string)
		p.Next()
	}
	f.Typ, err = p.parser.parseFunctionType()
	if err != nil {
		p.consume(untils_lc)
	}
	if p.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s except { but %s", p.parser.errorMsgPrefix(), p.parser.token.Desp)
		p.parser.errs = append(p.parser.errs, err)
		return
	}
	if ispublic {
		f.AccessFlags |= cg.ACC_METHOD_PUBLIC
	} else {
		f.AccessFlags |= cg.ACC_METHOD_PRIVATE
	}
	f.AccessFlags |= cg.ACC_METHOD_STATIC
	f.AccessFlags |= cg.ACC_METHOD_FINAL
	f.Block = &ast.Block{}
	err = p.parser.Block.parse(f.Block)
	return f, err
}

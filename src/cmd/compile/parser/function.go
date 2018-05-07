package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
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

func (p *Function) parse(needName bool) (f *ast.Function, err error) {
	p.Next() // skip fn key word
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
		err = fmt.Errorf("%s except '{' but '%s'", p.parser.errorMsgPrefix(), p.parser.token.Desp)
		p.parser.errs = append(p.parser.errs, err)
		return
	}
	f.Block.IsFunctionTopBlock = true
	p.Next()
	err = p.parser.Block.parse(&f.Block, false, lex.TOKEN_RC)
	return f, err
}

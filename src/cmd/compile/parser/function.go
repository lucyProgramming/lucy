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
	f = &ast.Function{}
	{
		offset := p.parser.scanner.GetOffSet()
		defer func() {
			if f.Block.EndPos != nil {
				f.SourceCode = p.parser.bs[offset:f.Block.EndPos.Offset]
			}
		}()
	}
	p.Next() // skip fn key word
	f.Pos = p.parser.mkPos()
	if needName {
		if p.parser.token.Type != lex.TOKEN_IDENTIFIER {
			err := fmt.Errorf("%s expect function name,but '%s'",
				p.parser.errorMsgPrefix(), p.parser.token.Desp)
			p.parser.errs = append(p.parser.errs, err)
			return nil, err
		}
	}
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
		p.consume(untils_lc)
	}
	f.Block.IsFunctionTopBlock = true
	p.Next()
	err = p.parser.Block.parse(&f.Block, false, false, lex.TOKEN_RC)
	return f, err
}

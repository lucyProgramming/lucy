package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type FunctionParser struct {
	parser *Parser
}

func (p *FunctionParser) Next() {
	p.parser.Next()
}

func (p *FunctionParser) consume(untils map[int]bool) {
	p.parser.consume(untils)
}

func (p *FunctionParser) parse(needName bool) (f *ast.Function, err error) {
	f = &ast.Function{}
	var offset int
	//	if p.parser.token == nil {
	//		offset = 0 // template function
	//	} else {
	offset = p.parser.token.Offset
	//	}

	p.Next() // skip fn key word
	f.Pos = p.parser.mkPos()
	if needName {
		if p.parser.token.Type != lex.TOKEN_IDENTIFIER {
			err := fmt.Errorf("%s expect function name,but '%s'",
				p.parser.errorMsgPrefix(), p.parser.token.Description)
			p.parser.errs = append(p.parser.errs, err)
			if p.parser.token.Type != lex.TOKEN_LC {
				return nil, err
			}
		}
	}
	if p.parser.token.Type == lex.TOKEN_IDENTIFIER {
		f.Name = p.parser.token.Data.(string)
		p.Next()
	}
	f.Type, err = p.parser.parseFunctionType()
	if err != nil {
		p.consume(untils_lc)
	}
	if p.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s except '{' but '%s'", p.parser.errorMsgPrefix(), p.parser.token.Description)
		p.parser.errs = append(p.parser.errs, err)
		p.consume(untils_lc)
	}
	f.Block.IsFunctionTopBlock = true
	p.Next()
	p.parser.BlockParser.parseStatementList(&f.Block, false)
	if p.parser.token.Type != lex.TOKEN_RC {
		err = fmt.Errorf("%s expect '}', but '%s'",
			p.parser.errorMsgPrefix(), p.parser.token.Description)
	} else {
		f.SourceCodes = p.parser.bs[offset : p.parser.token.Offset+1]
		p.Next()
	}

	return f, err
}

package parser

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func init() {
	ast.ParseFunctionHandler = ParseFunction
}

var (
	untilLc = map[int]bool{
		lex.TokenLc: true,
	}
	untilGt = map[int]bool{
		lex.TokenGt: true,
	}
	untilRc = map[int]bool{
		lex.TokenRc: true,
	}
	untilSemicolon = map[int]bool{
		lex.TokenSemicolon: true,
	}
	untilRcAndSemicolon = map[int]bool{
		lex.TokenRc:        true,
		lex.TokenSemicolon: true,
	}
)

func ParseFunction(bs []byte, pos *ast.Position) (*ast.Function, []error) {
	p := &Parser{}
	p.filename = pos.Filename
	p.nErrors2Stop = 10
	tops := []*ast.Top{}
	p.tops = &tops
	p.bs = bs
	p.FunctionParser = &FunctionParser{
		parser: p,
	}
	p.BlockParser = &BlockParser{
		parser: p,
	}
	p.ExpressionParser = &ExpressionParser{
		parser: p,
	}
	p.ClassParser = &ClassParser{
		parser: p,
	}
	p.InterfaceParser = &InterfaceParser{
		parser: p,
	}
	p.scanner = lex.New(p.bs, pos.StartLine, pos.StartColumn)
	p.Next() // parse fn
	f, err := p.FunctionParser.parse(true)
	if err != nil {
		p.errs = append(p.errs, err)
	}
	return f, p.errs
}

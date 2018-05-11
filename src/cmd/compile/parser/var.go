package parser

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func init() {
	ast.ParseFunctionHandler = ParseFunction
}

var (
	untils_lc = map[int]bool{
		lex.TOKEN_LC: true,
	}
	untils_gt = map[int]bool{
		lex.TOKEN_GT: true,
	}
	untils_rc = map[int]bool{
		lex.TOKEN_RC: true,
	}
	untils_semicolon = map[int]bool{
		lex.TOKEN_SEMICOLON: true,
	}
	untils_rc_semicolon = map[int]bool{
		lex.TOKEN_RC:        true,
		lex.TOKEN_SEMICOLON: true,
	}
)

func ParseFunction(bs []byte, pos *ast.Pos) (*ast.Function, []error) {
	p := &Parser{}
	p.filename = pos.Filename
	p.nerr = 10
	tops := []*ast.Node{}
	p.tops = &tops
	p.bs = bs
	p.Function = &Function{
		parser: p,
	}
	p.Block = &Block{
		parser: p,
	}
	p.Expression = &Expression{
		parser: p,
	}
	p.Class = &Class{
		parser: p,
	}
	p.Interface = &Interface{
		parser: p,
	}
	p.scanner = lex.New(p.bs, pos.StartLine, pos.StartColumn)
	f, err := p.Function.parse(true)
	if err != nil {
		p.errs = append(p.errs, err)
	}
	return f, p.errs
}

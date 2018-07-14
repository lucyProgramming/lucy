package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func init() {
	ast.ParseFunctionHandler = ParseFunction
}

var (
	untilLp = map[int]bool{
		lex.TokenLp: true,
	}
	untilRp = map[int]bool{
		lex.TokenRp: true,
	}
	untilGt = map[int]bool{
		lex.TokenGt: true,
	}
	untilLc = map[int]bool{
		lex.TokenLc: true,
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
	parser := &Parser{}
	parser.filename = pos.Filename
	parser.nErrors2Stop = 10
	parser.bs = bs
	parser.initParser()
	parser.scanner = lex.New(parser.bs, pos.StartLine, pos.StartColumn)
	parser.Next() //
	f, err := parser.FunctionParser.parse(true)
	if err != nil {
		parser.errs = append(parser.errs, err)
	}
	return f, parser.errs
}

var (
	autoNameIndex = 1
)

func compileAutoName() string {
	s := fmt.Sprintf("autoName$%d", autoNameIndex)
	autoNameIndex++
	return s
}

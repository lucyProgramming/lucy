package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

const (
	lfIsToken  = true
	lfNotToken = false
)

func init() {
	ast.ParseFunctionHandler = ParseFunction
}

var (
	untilLp = map[lex.TokenKind]bool{
		lex.TokenLp: true,
	}
	untilRp = map[lex.TokenKind]bool{
		lex.TokenRp: true,
	}
	untilGt = map[lex.TokenKind]bool{
		lex.TokenGt: true,
	}
	untilLc = map[lex.TokenKind]bool{
		lex.TokenLc: true,
	}
	untilRc = map[lex.TokenKind]bool{
		lex.TokenRc: true,
	}
	untilComma = map[lex.TokenKind]bool{
		lex.TokenComma: true,
	}
	untilSemicolonOrLf = map[lex.TokenKind]bool{
		lex.TokenSemicolon: true,
		lex.TokenLf:        true,
	}
)

func ParseFunction(bs []byte, pos *ast.Pos) (*ast.Function, []error) {
	parser := &Parser{}
	parser.filename = pos.Filename
	parser.nErrors2Stop = 10
	parser.bs = bs
	parser.initParser()
	parser.lexer = lex.New(parser.bs, pos.StartLine, pos.StartColumn)
	parser.Next(lfNotToken) //
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

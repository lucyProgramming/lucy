package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//atBeginningOfFile bool
func (parser *Parser) parseImports() []*ast.Import {
	ret := []*ast.Import{}
	for parser.token.Type == lex.TokenImport ||
		parser.token.Type == lex.TokenComment ||
		parser.token.Type == lex.TokenCommentMultiLine {
		if parser.token.Type == lex.TokenComment ||
			parser.token.Type == lex.TokenCommentMultiLine {
			parser.Next(lfNotToken)
			continue
		}
		parser.Next(lfIsToken) // skip import key word
		parser.unExpectNewLineAndSkip()
		if parser.token.Type != lex.TokenLiteralString {
			parser.errs = append(parser.errs, fmt.Errorf("%s expect 'package' after import,but '%s'",
				parser.errMsgPrefix(), parser.token.Description))
			parser.consume(untilSemicolonOrLf)
			parser.Next(lfNotToken)
			continue
		}
		i := &ast.Import{}
		i.Pos = parser.mkPos()
		i.Import = parser.token.Data.(string)
		ret = append(ret, i)
		parser.Next(lfIsToken) // skip name
		if parser.token.Type == lex.TokenAs {
			/*
				import "xxxxxxxxxxx" as yyy
			*/
			parser.Next(lfNotToken) // skip as
			if parser.token.Type != lex.TokenIdentifier {
				parser.errs = append(parser.errs, fmt.Errorf("%s expect 'identifier' after 'as',but '%s'",
					parser.errMsgPrefix(), parser.token.Description))
				parser.consume(untilSemicolonOrLf)
				parser.Next(lfNotToken)
				continue
			} else {
				i.Alias = parser.token.Data.(string)
				parser.Next(lfIsToken) // skip identifier
			}
		}
		parser.validStatementEnding()
		parser.Next(lfNotToken)
	}
	return ret
}

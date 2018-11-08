package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//atBeginningOfFile bool
func (this *Parser) parseImports() []*ast.Import {
	var ret []*ast.Import
	for this.token.Type == lex.TokenImport ||
		this.token.Type == lex.TokenComment ||
		this.token.Type == lex.TokenMultiLineComment {
		if this.token.Type == lex.TokenComment ||
			this.token.Type == lex.TokenMultiLineComment {
			this.Next(lfNotToken)
			continue
		}
		this.Next(lfIsToken) // skip import key word
		this.unExpectNewLineAndSkip()
		if this.token.Type != lex.TokenLiteralString {
			this.errs = append(this.errs, fmt.Errorf("%s expect 'package' after import,but '%s'",
				this.errMsgPrefix(), this.token.Description))
			this.consume(untilSemicolonOrLf)
			this.Next(lfNotToken)
			continue
		}
		i := &ast.Import{}
		i.Pos = this.mkPos()
		i.Import = this.token.Data.(string)
		ret = append(ret, i)
		this.Next(lfIsToken) // skip name
		if this.token.Type == lex.TokenAs {
			/*
				import "xxxxxxxxxxx" as yyy
			*/
			this.Next(lfNotToken) // skip as
			if this.token.Type != lex.TokenIdentifier {
				this.errs = append(this.errs, fmt.Errorf("%s expect 'identifier' after 'as',but '%s'",
					this.errMsgPrefix(), this.token.Description))
				this.consume(untilSemicolonOrLf)
				this.Next(lfNotToken)
				continue
			} else {
				i.Alias = this.token.Data.(string)
				this.Next(lfIsToken) // skip identifier
			}
		}
		this.validStatementEnding()
		this.Next(lfNotToken)
	}
	return ret
}

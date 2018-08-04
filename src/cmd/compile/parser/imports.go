package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//atBeginningOfFile bool
func (parser *Parser) parseImports() {

	if parser.token.Type != lex.TokenImport {
		// not a import
		return
	}
	parser.Next(lfIsToken) // skip import key word
	if err := parser.unExpectNewLine(); err != nil {
		parser.consume(untilSemicolonOrLf)
		parser.Next(lfNotToken)
		parser.parseImports()
		return
	}
	if parser.token.Type != lex.TokenLiteralString {
		parser.errs = append(parser.errs, fmt.Errorf("%s expect 'package' after import,but '%s'",
			parser.errorMsgPrefix(), parser.token.Description))
		parser.consume(untilSemicolonOrLf)
		parser.Next(lfNotToken)
		parser.parseImports()
		return
	}
	i := &ast.Import{}
	i.Pos = parser.mkPos()
	i.Import = parser.token.Data.(string)
	parser.Next(lfIsToken) // skip name
	if parser.token.Type == lex.TokenAs {
		/*
			import "xxxxxxxxxxx" as yyy
		*/
		parser.Next(lfNotToken) // skip as
		if parser.token.Type != lex.TokenIdentifier {
			parser.insertImports(i)
			parser.errs = append(parser.errs, fmt.Errorf("%s expect 'identifier' after 'as',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description))
			parser.consume(untilSemicolonOrLf)
			parser.Next(lfNotToken)
			parser.parseImports()
			return
		} else {
			i.AccessName = parser.token.Data.(string)
			parser.Next(lfIsToken) // skip identifier
		}
	}
	parser.validStatementEnding()
	parser.Next(lfNotToken)
	parser.insertImports(i)
	parser.parseImports()
}

func (parser *Parser) insertImports(im *ast.Import) {
	if parser.importsByAccessName == nil {
		parser.importsByAccessName = make(map[string]*ast.Import)
	}
	if parser.importsByResourceName == nil {
		parser.importsByResourceName = make(map[string]*ast.Import)
	}
	err := im.MkAccessName()
	if err != nil {
		parser.errs = append(parser.errs, fmt.Errorf("%s %v", parser.errorMsgPrefix(im.Pos), err))
		return
	}
	*parser.tops = append(*parser.tops, &ast.Top{
		Data: im,
	})
	if im.AccessName != ast.NoNameIdentifier {
		if parser.importsByAccessName[im.AccessName] != nil {
			parser.errs = append(parser.errs, fmt.Errorf("%s '%s' reImported",
				parser.errorMsgPrefix(im.Pos), im.AccessName))
			return
		}
		parser.importsByAccessName[im.AccessName] = im
	}
	if parser.importsByResourceName[im.Import] != nil {
		parser.errs = append(parser.errs, fmt.Errorf("%s '%s' reImported",
			parser.errorMsgPrefix(im.Pos), im.Import))
		return
	}
	parser.importsByResourceName[im.Import] = im
}

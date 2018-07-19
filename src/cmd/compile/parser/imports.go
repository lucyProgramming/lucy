package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (parser *Parser) parseImports() {
	if parser.token.Type != lex.TokenImport {
		// not a import
		return
	}
	parser.Next(lfIsToken) // skip import key word
	if err := parser.unExpectNewLine(); err != nil {
		parser.consume(untilSemicolonAndLf)
		parser.Next(lfNotToken)
		parser.parseImports()
		return
	}
	if parser.token.Type != lex.TokenLiteralString {
		parser.errs = append(parser.errs, fmt.Errorf("%s expect 'package' after import,but '%s'",
			parser.errorMsgPrefix(), parser.token.Description))
		parser.consume(untilSemicolonAndLf)
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
			parser.consume(untilSemicolonAndLf)
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
	if parser.imports == nil {
		parser.imports = make(map[string]*ast.Import)
	}
	access, err := im.GetAccessName()
	if err != nil {
		parser.errs = append(parser.errs, fmt.Errorf("%s %v", parser.errorMsgPrefix(im.Pos), err))
		return
	}
	if parser.imports[access] != nil {
		parser.errs = append(parser.errs, fmt.Errorf("%s package '%s' reImported",
			parser.errorMsgPrefix(im.Pos), access))
		return
	}
	parser.imports[access] = im
	*parser.tops = append(*parser.tops, &ast.Top{
		Data: im,
	})
}

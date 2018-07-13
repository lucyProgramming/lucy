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
	parser.Next() // skip import key word
	if parser.token.Type != lex.TokenLiteralString {
		parser.errs = append(parser.errs, fmt.Errorf("%s expect 'package' after import,but '%s'",
			parser.errorMsgPrefix(), parser.token.Description))
		parser.consume(untilSemicolon)
		parser.Next()
		parser.parseImports()
		return
	}
	i := &ast.Import{}
	i.Pos = parser.mkPos()
	i.Import = parser.token.Data.(string)
	parser.Next() // skip name
	if parser.token.Type == lex.TokenAs {
		/*
			import "xxxxxxxxxxx" as yyy
		*/
		parser.Next() // skip as
		if parser.token.Type != lex.TokenIdentifier {
			parser.insertImports(i)
			parser.errs = append(parser.errs, fmt.Errorf("%s expect 'identifier' after 'as',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description))
			parser.consume(untilSemicolon)
			parser.Next()
			parser.parseImports()
			return
		} else {
			i.AccessName = parser.token.Data.(string)
			parser.Next() // skip identifier
		}
	}
	if parser.token.Type != lex.TokenSemicolon {
		parser.errs = append(parser.errs, fmt.Errorf("%s expect semicolon, but '%s'",
			parser.errorMsgPrefix(), parser.token.Description))
		if parser.token.Type != lex.TokenImport { // next token is not import
			parser.consume(untilSemicolon)
		}
	}
	parser.Next() // skip ;
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

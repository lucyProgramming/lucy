package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//imports,always call next
func (parser *Parser) parseImports() {
	if parser.token.Type != lex.TokenImport {
		// not a import
		return
	}
	parser.Next()
	if parser.token.Type != lex.TokenLiteralString {
		parser.consume(untilSemicolon)
		parser.errs = append(parser.errs, fmt.Errorf("%s expect 'string_literal' after import,but '%s'",
			parser.errorMsgPrefix(), parser.token.Description))
		parser.parseImports()
		return
	}
	packageName := parser.token.Data.(string)
	parser.Next()
	i := &ast.Import{}
	i.ImportName = packageName
	i.Pos = parser.mkPos()
	if parser.token.Type == lex.TokenAs { // import "xxxxxxxxxxx" as yyy ;
		parser.Next() // skip as
		if parser.token.Type != lex.TokenIdentifier {
			parser.errs = append(parser.errs, fmt.Errorf("%s expect 'identifier' after 'as',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description))
			parser.consume(untilSemicolon)
			parser.Next()
			parser.insertImports(i)
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
		parser.consume(untilSemicolon)
		parser.consume(untilSemicolon)
		parser.Next()
		parser.insertImports(i)
		parser.parseImports()
		return
	}
	parser.insertImports(i)
	parser.Next() // skip ;
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
		parser.errs = append(parser.errs, fmt.Errorf("%s package '%s' reimported",
			parser.errorMsgPrefix(im.Pos), access))
		return
	}
	parser.imports[access] = im
	*parser.tops = append(*parser.tops, &ast.Top{
		Data: im,
	})
}

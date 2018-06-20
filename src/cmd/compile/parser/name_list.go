package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//at least one name
func (parser *Parser) parseNameList() (names []*ast.NameWithPos, err error) {
	if parser.token.Type != lex.TOKEN_IDENTIFIER {
		err = fmt.Errorf("%s expect identifier,but '%s'",
			parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return nil, err
	}
	names = []*ast.NameWithPos{}
	for parser.token.Type == lex.TOKEN_IDENTIFIER {
		names = append(names, &ast.NameWithPos{
			Name: parser.token.Data.(string),
			Pos:  parser.mkPos(),
		})
		parser.Next()
		if parser.token.Type != lex.TOKEN_COMMA {
			// not a ,
			break
		}
		parser.Next()
		if parser.token.Type != lex.TOKEN_IDENTIFIER {
			err = fmt.Errorf("%s not a 'identifier' after a comma,but '%s'",
				parser.errorMsgPrefix(), parser.token.Description)
			parser.errs = append(parser.errs, err)
			return names, err
		}
	}
	return
}

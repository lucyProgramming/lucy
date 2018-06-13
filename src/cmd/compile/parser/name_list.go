package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//at least one name
func (p *Parser) parseNameList() (names []*ast.NameWithPos, err error) {
	if p.token.Type != lex.TOKEN_IDENTIFIER {
		err = fmt.Errorf("%s expect identifier,but '%s'",
			p.errorMsgPrefix(), p.token.Description)
		p.errs = append(p.errs, err)
		return nil, err
	}
	names = []*ast.NameWithPos{}
	for p.token.Type == lex.TOKEN_IDENTIFIER {
		names = append(names, &ast.NameWithPos{
			Name: p.token.Data.(string),
			Pos:  p.mkPos(),
		})
		p.Next()
		if p.token.Type != lex.TOKEN_COMMA {
			// not a ,
			break
		}
		p.Next()
		if p.token.Type != lex.TOKEN_IDENTIFIER {
			err = fmt.Errorf("%s not a 'identifier' after a comma,but '%s'",
				p.errorMsgPrefix(), p.token.Description)
			p.errs = append(p.errs, err)
			return names, err
		}
	}
	return
}

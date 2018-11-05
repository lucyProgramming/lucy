package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//at least one name
func (this *Parser) parseNameList() (names []*ast.NameWithPos, err error) {
	if this.token.Type != lex.TokenIdentifier {
		err = fmt.Errorf("%s expect identifier,but '%s'",
			this.errMsgPrefix(), this.token.Description)
		this.errs = append(this.errs, err)
		return nil, err
	}
	names = []*ast.NameWithPos{}
	for this.token.Type == lex.TokenIdentifier {
		names = append(names, &ast.NameWithPos{
			Name: this.token.Data.(string),
			Pos:  this.mkPos(),
		})
		this.Next(lfIsToken)
		if this.token.Type != lex.TokenComma {
			// not a ,
			break
		} else {
			this.Next(lfNotToken) // skip comma
			if this.token.Type != lex.TokenIdentifier {
				err = fmt.Errorf("%s not a 'identifier' after a comma,but '%s'",
					this.errMsgPrefix(), this.token.Description)
				this.errs = append(this.errs, err)
				return names, err
			}
		}
	}
	return
}

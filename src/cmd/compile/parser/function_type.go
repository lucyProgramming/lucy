package parser

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"

	"github.com/756445638/lucy/src/cmd/compile/lex"
)

//(a,b int)->(total int)
func (p *Parser) parseFunctionType() (t *ast.FunctionType, err error) {
	t = &ast.FunctionType{}
	if p.token.Type != lex.TOKEN_LP {
		err = fmt.Errorf("%s fn declared wrong,missing (,but %s", p.errorMsgPrefix(), p.token.Desp)
		p.errs = append(p.errs, err)
		return
	}
	p.Next()                          // skip (
	if p.token.Type != lex.TOKEN_RP { // not (
		t.ParameterList, err = p.parseTypedNames()
		if err != nil {
			return nil, err
		}
	}
	if p.token.Type != lex.TOKEN_RP { // not )
		err = fmt.Errorf("%s fn declared wrong,missing ),but %s", p.errorMsgPrefix(), p.token.Desp)
		p.errs = append(p.errs, err)
		return
	}
	p.Next()
	if p.token.Type == lex.TOKEN_ARROW { // ->
		p.Next() // skip ->
		if p.token.Type != lex.TOKEN_LP {
			err = fmt.Errorf("%s fn declared wrong, not ( after ->", p.errorMsgPrefix())
			p.errs = append(p.errs, err)
			return
		}
		p.Next()
		t.ReturnList, err = p.parseTypedNames()
		if err != nil {
			return
		}
		if p.token.Type != lex.TOKEN_RP {
			err = fmt.Errorf("%s fn declared wrong, ( and ) not match", p.errorMsgPrefix())
			p.errs = append(p.errs, err)
			return
		}
		p.Next()
	}
	return t, err
}

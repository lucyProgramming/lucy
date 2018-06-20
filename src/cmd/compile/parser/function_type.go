package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//(a,b int)->(total int)
func (p *Parser) parseFunctionType() (t ast.FunctionType, err error) {
	t = ast.FunctionType{}
	if p.token.Type != lex.TOKEN_LP {
		err = fmt.Errorf("%s fn declared wrong,missing '(',but '%s'",
			p.errorMsgPrefix(), p.token.Description)
		p.errs = append(p.errs, err)
		return
	}
	p.Next()                          // skip (
	if p.token.Type != lex.TOKEN_RP { // not (
		t.ParameterList, err = p.parseReturnLists()
		if err != nil {
			return t, err
		}
	}
	if p.token.Type != lex.TOKEN_RP { // not )
		err = fmt.Errorf("%s fn declared wrong,missing ')',but '%s'",
			p.errorMsgPrefix(), p.token.Description)
		p.errs = append(p.errs, err)

		return
	}
	p.Next()
	if p.token.Type == lex.TOKEN_ARROW { // ->
		p.Next() // skip ->
		if p.token.Type != lex.TOKEN_LP {
			err = fmt.Errorf("%s fn declared wrong, not '(' after '->'",
				p.errorMsgPrefix())
			p.errs = append(p.errs, err)
			return
		}
		p.Next() // skip (
		if p.token.Type != lex.TOKEN_RP {
			t.ReturnList, err = p.parseReturnLists()
			if err != nil { // skip until next (,continue to analyse
				p.consume(map[int]bool{
					lex.TOKEN_RP: true,
				})
				p.Next()
			}
		}
		if p.token.Type != lex.TOKEN_RP {
			err = fmt.Errorf("%s fn declared wrong,expected ')',but '%s'",
				p.errorMsgPrefix(), p.token.Description)
			p.errs = append(p.errs, err)
			return
		}
		p.Next()
	}
	return t, err
}

func (p *Parser) parseReturnList() (returnList []*ast.Variable, err error) {
	returnList, err = p.parseTypedName()
	if p.token.Type != lex.TOKEN_ASSIGN {
		return
	}
	p.Next() // skip =
	for k, v := range returnList {
		var er error
		v.Expression, er = p.ExpressionParser.parseExpression(false)
		if er != nil {
			p.errs = append(p.errs, err)
			p.consume(map[int]bool{
				lex.TOKEN_COMMA: true,
			})
			err = er
			p.Next()
			continue
		}
		if p.token.Type != lex.TOKEN_COMMA || k == len(returnList)-1 {
			break
		} else {
			p.Next() // skip ,
		}
	}
	return returnList, nil
}
func (p *Parser) parseReturnLists() (returnList []*ast.Variable, err error) {
	for p.token.Type == lex.TOKEN_IDENTIFIER {
		v, err := p.parseReturnList()
		if v != nil {
			returnList = append(returnList, v...)
		}
		if err != nil {
			break
		}
		if p.token.Type == lex.TOKEN_COMMA {
			p.Next()
		} else {
			break
		}
	}
	return
}

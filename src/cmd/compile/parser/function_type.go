package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//(a,b int)->(total int)
func (parser *Parser) parseFunctionType() (t ast.FunctionType, err error) {
	t = ast.FunctionType{}
	if parser.token.Type != lex.TokenLp {
		err = fmt.Errorf("%s fn declared wrong,missing '(',but '%s'",
			parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return
	}
	parser.Next()                         // skip (
	if parser.token.Type != lex.TokenRp { // not (
		t.ParameterList, err = parser.parseReturnLists()
		if err != nil {
			return t, err
		}
	}
	if parser.token.Type != lex.TokenRp { // not )
		err = fmt.Errorf("%s fn declared wrong,missing ')',but '%s'",
			parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return
	}
	parser.Next()
	if parser.token.Type == lex.TokenArrow { // ->
		parser.Next() // skip ->
		if parser.token.Type != lex.TokenLp {
			err = fmt.Errorf("%s fn declared wrong, not '(' after '->'",
				parser.errorMsgPrefix())
			parser.errs = append(parser.errs, err)
			return
		}
		parser.Next() // skip (
		if parser.token.Type != lex.TokenRp {
			t.ReturnList, err = parser.parseReturnLists()
			if err != nil { // skip until next (,continue to analyse
				parser.consume(map[int]bool{
					lex.TokenRp: true,
				})
				parser.Next()
			}
		}
		if parser.token.Type != lex.TokenRp {
			err = fmt.Errorf("%s fn declared wrong,expected ')',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description)
			parser.errs = append(parser.errs, err)
			return
		}
		parser.Next()
	}
	return t, err
}

func (parser *Parser) parseReturnList() (returnList []*ast.Variable, err error) {
	returnList, err = parser.parseTypedName()
	if parser.token.Type != lex.TokenAssign {
		return
	}
	parser.Next() // skip =
	for k, v := range returnList {
		var er error
		v.Expression, er = parser.ExpressionParser.parseTernaryExpression()
		if er != nil {
			parser.errs = append(parser.errs, err)
			parser.consume(map[int]bool{
				lex.TokenComma: true,
			})
			err = er
			parser.Next()
			continue
		}
		if parser.token.Type != lex.TokenComma || k == len(returnList)-1 {
			break
		} else {
			parser.Next() // skip ,
		}
	}
	return returnList, nil
}
func (parser *Parser) parseReturnLists() (returnList []*ast.Variable, err error) {
	for parser.token.Type == lex.TokenIdentifier {
		v, err := parser.parseReturnList()
		if v != nil {
			returnList = append(returnList, v...)
		}
		if err != nil {
			break
		}
		if parser.token.Type == lex.TokenComma {
			parser.Next()
		} else {
			break
		}
	}
	return
}

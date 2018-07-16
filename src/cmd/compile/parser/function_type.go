package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//(a,b int)->(total int)
func (parser *Parser) parseFunctionType() (functionType ast.FunctionType, err error) {
	functionType = ast.FunctionType{}
	if parser.token.Type != lex.TokenLp {
		err = fmt.Errorf("%s fn declared wrong,missing '(',but '%s'",
			parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return
	}
	parser.Next()                         // skip (
	if parser.token.Type != lex.TokenRp { // not )
		functionType.ParameterList, err = parser.parseReturnLists()
		if err != nil {
			parser.consume(untilRp)
			parser.Next()
		}
	}
	if parser.token.Type != lex.TokenRp { // not )
		err = fmt.Errorf("%s fn declared wrong,missing ')',but '%s'",
			parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return
	}
	parser.Next()                            // skip )
	if parser.token.Type == lex.TokenArrow { // ->  parse return list
		parser.Next() // skip ->
		if parser.token.Type != lex.TokenLp {
			err = fmt.Errorf("%s fn declared wrong, not '(' after '->'",
				parser.errorMsgPrefix())
			parser.errs = append(parser.errs, err)
			return
		}
		parser.Next() // skip (
		if parser.token.Type != lex.TokenRp {
			functionType.ReturnList, err = parser.parseReturnLists()
			if err != nil { // skip until next (,continue to analyse
				parser.consume(untilRp)
				parser.Next()
			}
		}
		if parser.token.Type != lex.TokenRp {
			err = fmt.Errorf("%s fn declared wrong,expected ')',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description)
			parser.errs = append(parser.errs, err)
			return
		}
		parser.Next() // skip )
	}
	return functionType, err
}

func (parser *Parser) parseTypedNameForReturnVar() (returnList []*ast.Variable, err error) {
	returnList, err = parser.parseTypedName()
	if parser.token.Type != lex.TokenAssign {
		return
	}
	parser.Next() // skip =
	for k, v := range returnList {
		var er error
		v.Expression, er = parser.ExpressionParser.parseExpression(false)
		if er != nil {
			parser.errs = append(parser.errs, err)
			parser.consume(untilComma)
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
	return returnList, err
}
func (parser *Parser) parseReturnLists() (returnList []*ast.Variable, err error) {
	for parser.token.Type == lex.TokenIdentifier {
		v, err := parser.parseTypedNameForReturnVar()
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

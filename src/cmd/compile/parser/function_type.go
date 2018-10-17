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
			parser.errMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return
	}
	parser.Next(lfNotToken)               // skip (
	if parser.token.Type != lex.TokenRp { // not )
		functionType.ParameterList, err = parser.parseParameterOrReturnList()
		if err != nil {
			parser.consume(untilRp)
			parser.Next(lfNotToken)
		}
	}
	parser.ifTokenIsLfThenSkip()
	if parser.token.Type != lex.TokenRp { // not )
		err = fmt.Errorf("%s fn declared wrong,missing ')',but '%s'",
			parser.errMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return
	}
	parser.Next(lfIsToken)                   // skip )
	if parser.token.Type == lex.TokenArrow { // ->  parse return list
		parser.Next(lfNotToken) // skip ->
		if parser.token.Type != lex.TokenLp {
			err = fmt.Errorf("%s fn declared wrong, not '(' after '->'",
				parser.errMsgPrefix())
			parser.errs = append(parser.errs, err)
			return
		}
		parser.Next(lfNotToken) // skip (
		if parser.token.Type != lex.TokenRp {
			functionType.ReturnList, err = parser.parseParameterOrReturnList()
			if err != nil { // skip until next (,continue to analyse
				parser.consume(untilRp)
				parser.Next(lfIsToken)
			}
		}
		parser.ifTokenIsLfThenSkip()
		if parser.token.Type != lex.TokenRp {
			err = fmt.Errorf("%s fn declared wrong,expected ')',but '%s'",
				parser.errMsgPrefix(), parser.token.Description)
			parser.errs = append(parser.errs, err)
			return
		}
		parser.Next(lfIsToken) // skip )
	} else {
		functionType.ReturnList = make([]*ast.Variable, 1)
		functionType.ReturnList[0] = &ast.Variable{}
		functionType.ReturnList[0].Pos = parser.mkPos()
		functionType.ReturnList[0].Type = &ast.Type{}
		functionType.ReturnList[0].Type.Pos = parser.mkPos()
		functionType.ReturnList[0].Type.Type = ast.VariableTypeVoid
	}
	return functionType, err
}

/*
	parse default value
	a int = ""
	int = 1

*/
func (parser *Parser) parseTypedNameDefaultValue() (returnList []*ast.Variable, err error) {
	returnList, err = parser.parseTypedName()
	if parser.token.Type != lex.TokenAssign {
		return
	}
	parser.Next(lfIsToken) // skip =
	for k, v := range returnList {
		var er error
		v.DefaultValueExpression, er = parser.ExpressionParser.parseExpression(false)
		if er != nil {
			parser.consume(untilComma)
			err = er
			parser.Next(lfNotToken)
			continue
		}
		if parser.token.Type != lex.TokenComma ||
			k == len(returnList)-1 {
			break
		} else {
			parser.Next(lfNotToken) // skip ,
		}
	}
	return returnList, err
}
func (parser *Parser) parseParameterOrReturnList() (returnList []*ast.Variable, err error) {
	for parser.token.Type != lex.TokenRp {
		if parser.token.Type == lex.TokenComma {
			parser.errs = append(parser.errs, fmt.Errorf("%s extra comma",
				parser.errMsgPrefix()))
			parser.Next(lfNotToken)
			continue
		}
		v, err := parser.parseTypedNameDefaultValue()
		if v != nil {
			returnList = append(returnList, v...)
		}
		if err != nil {
			break
		}
		if parser.token.Type == lex.TokenComma {
			parser.Next(lfNotToken)
		} else {
			break
		}
	}
	return
}

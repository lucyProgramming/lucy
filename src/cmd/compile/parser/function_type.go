package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//(a,b int)->(total int)
func (this *Parser) parseFunctionType() (functionType ast.FunctionType, err error) {
	functionType = ast.FunctionType{}
	if this.token.Type == lex.TokenLt {
		this.Next(lfNotToken)
		var err error
		functionType.TemplateNames, err = this.parseNameList()
		if err != nil {
			this.consume(untilLp)
			goto skipTemplateNames
		}
		this.errs = append(this.errs, functionType.CheckTemplateNameDuplication()...)
		this.Next(lfIsToken)
	}
	this.unExpectNewLineAndSkip()
skipTemplateNames:
	if this.token.Type != lex.TokenLp {
		err = fmt.Errorf("%s fn declared wrong,missing '(',but '%s'",
			this.errMsgPrefix(), this.token.Description)
		this.errs = append(this.errs, err)
		return
	}
	this.Next(lfNotToken)               // skip (
	if this.token.Type != lex.TokenRp { // not )
		functionType.ParameterList, err = this.parseParameterOrReturnList()
		if err != nil {
			this.consume(untilRp)
			this.Next(lfNotToken)
		}
	}
	this.ifTokenIsLfThenSkip()
	if this.token.Type != lex.TokenRp { // not )
		err = fmt.Errorf("%s fn declared wrong,missing ')',but '%s'",
			this.errMsgPrefix(), this.token.Description)
		this.errs = append(this.errs, err)
		return
	}
	this.Next(lfIsToken)                   // skip )
	if this.token.Type == lex.TokenArrow { // ->  parse return list
		this.Next(lfNotToken) // skip ->
		if this.token.Type != lex.TokenLp {
			err = fmt.Errorf("%s fn declared wrong, not '(' after '->'",
				this.errMsgPrefix())
			this.errs = append(this.errs, err)
			return
		}
		this.Next(lfNotToken) // skip (
		if this.token.Type != lex.TokenRp {
			functionType.ReturnList, err = this.parseParameterOrReturnList()
			if err != nil { // skip until next (,continue to analyse
				this.consume(untilRp)
				this.Next(lfIsToken)
			}
		}
		this.ifTokenIsLfThenSkip()
		if this.token.Type != lex.TokenRp {
			err = fmt.Errorf("%s fn declared wrong,expected ')',but '%s'",
				this.errMsgPrefix(), this.token.Description)
			this.errs = append(this.errs, err)
			return
		}
		this.Next(lfIsToken) // skip )
	} else {
		functionType.ReturnList = make([]*ast.Variable, 1)
		functionType.ReturnList[0] = &ast.Variable{}
		functionType.ReturnList[0].Pos = this.mkPos()
		functionType.ReturnList[0].Type = &ast.Type{}
		functionType.ReturnList[0].Type.Pos = this.mkPos()
		functionType.ReturnList[0].Type.Type = ast.VariableTypeVoid
	}
	return functionType, err
}

/*
	parse default value
	a int = ""
	int = 1

*/
func (this *Parser) parseTypedNameDefaultValue() (returnList []*ast.Variable, err error) {
	returnList, err = this.parseTypedName()
	if this.token.Type != lex.TokenAssign {
		return
	}
	this.Next(lfIsToken) // skip =
	for k, v := range returnList {
		var er error
		v.DefaultValueExpression, er = this.ExpressionParser.parseExpression(false)
		if er != nil {
			this.consume(untilComma)
			err = er
			this.Next(lfNotToken)
			continue
		}
		if this.token.Type != lex.TokenComma ||
			k == len(returnList)-1 {
			break
		} else {
			this.Next(lfNotToken) // skip ,
		}
	}
	return returnList, err
}
func (this *Parser) parseParameterOrReturnList() (returnList []*ast.Variable, err error) {
	for this.token.Type != lex.TokenRp {
		if this.token.Type == lex.TokenComma {
			this.errs = append(this.errs, fmt.Errorf("%s extra comma",
				this.errMsgPrefix()))
			this.Next(lfNotToken)
			continue
		}
		v, err := this.parseTypedNameDefaultValue()
		if v != nil {
			returnList = append(returnList, v...)
		}
		if err != nil {
			break
		}
		if this.token.Type == lex.TokenComma {
			this.Next(lfNotToken)
		} else {
			break
		}
	}
	return
}

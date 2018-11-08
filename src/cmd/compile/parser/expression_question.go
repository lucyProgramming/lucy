package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

/*
	true ? 1 : 2
*/
func (ep *ExpressionParser) parseQuestionExpression() (*ast.Expression, error) {
	left, err := ep.parseLogicalOrExpression()
	if err != nil {
		return left, err
	}
	if ep.parser.token.Type != lex.TokenQuestion {
		return left, nil
	}
	pos := ep.parser.mkPos()
	ep.Next(lfNotToken) // skip ?
	True, err := ep.parseLogicalOrExpression()
	if err != nil {
		return left, nil
	}
	ep.parser.unExpectNewLineAndSkip()
	if ep.parser.token.Type != lex.TokenColon {
		err := fmt.Errorf("%s expect ':' ,but '%s'",
			ep.parser.errMsgPrefix(), ep.parser.token.Description)
		ep.parser.errs = append(ep.parser.errs, err)
		return left, err
	}
	ep.Next(lfNotToken) // skip :
	False, err := ep.parseLogicalOrExpression()
	if err != nil {
		return left, nil
	}
	newExpression := &ast.Expression{}
	newExpression.Op = "question"
	newExpression.Pos = pos
	newExpression.Type = ast.ExpressionTypeQuestion
	question := &ast.ExpressionQuestion{}
	question.Selection = left
	question.True = True
	question.False = False
	newExpression.Data = question
	return newExpression, nil
}

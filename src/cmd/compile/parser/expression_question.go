package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (expressionParser *ExpressionParser) parseQuestionExpression() (*ast.Expression, error) {
	left, err := expressionParser.parseLogicalOrExpression()
	if err != nil {
		return left, err
	}
	if expressionParser.parser.token.Type != lex.TokenQuestion {
		return left, nil
	}
	expressionParser.Next(lfNotToken) // skip ?
	True, err := expressionParser.parseLogicalOrExpression()
	if err != nil {
		return left, nil
	}
	expressionParser.parser.unExpectNewLineAndSkip()
	if expressionParser.parser.token.Type != lex.TokenColon {
		return left, fmt.Errorf("%s expect ':',but '%s'",
			expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
	}
	expressionParser.Next(lfNotToken) // skip :
	False, err := expressionParser.parseLogicalOrExpression()
	if err != nil {
		return left, nil
	}
	newExpression := &ast.Expression{}
	newExpression.Description = "question"
	newExpression.Pos = expressionParser.parser.mkPos()
	newExpression.Type = ast.ExpressionTypeQuestion
	question := &ast.ExpressionQuestion{}
	question.Selection = left
	question.True = True
	question.False = False
	newExpression.Data = question
	return newExpression, nil
}

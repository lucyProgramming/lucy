package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (expressionParser *ExpressionParser) parseTernaryExpression() (*ast.Expression, error) {
	left, err := expressionParser.parseLogicalOrExpression()
	if err != nil {
		return left, err
	}
	if expressionParser.parser.token.Type != lex.TokenQuestion {
		return left, nil
	}
	expressionParser.Next() // skip ?
	True, err := expressionParser.parseExpression(false)
	if err != nil {
		return left, nil
	}
	if expressionParser.parser.token.Type != lex.TokenColon {
		return left, fmt.Errorf("%s expect ':',but '%s'",
			expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
	}
	expressionParser.Next() // skip :
	False, err := expressionParser.parseExpression(false)
	if err != nil {
		return left, nil
	}
	newExpression := &ast.Expression{}
	newExpression.Pos = expressionParser.parser.mkPos()
	newExpression.Type = ast.ExpressionTypeTernary
	ternary := &ast.ExpressionTernary{}
	ternary.Selection = left
	ternary.True = True
	ternary.False = False
	newExpression.Data = ternary
	return newExpression, nil
}

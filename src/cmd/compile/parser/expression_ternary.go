package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (ep *ExpressionParser) parseTernaryExpression() (*ast.Expression, error) {
	left, err := ep.parseLogicalOrExpression()
	if err != nil {
		return left, err
	}
	if ep.parser.token.Type != lex.TOKEN_QUESTION {
		return left, nil
	}
	newe := &ast.Expression{}
	newe.Pos = ep.parser.mkPos()
	newe.Type = ast.EXPRESSION_TYPE_TERNARY
	ep.Next() // skip ?
	True, err := ep.parseExpression(false)
	if err != nil {
		return left, nil
	}
	if ep.parser.token.Type != lex.TOKEN_COLON {
		return left, fmt.Errorf("%s expect ':',but '%s'",
			ep.parser.errorMsgPrefix(), ep.parser.token.Description)
	}
	ep.Next() // skip :
	False, err := ep.parseExpression(false)
	if err != nil {
		return left, nil
	}
	ternary := &ast.ExpressionTernary{}
	ternary.Selection = left
	ternary.True = True
	ternary.False = False
	newe.Data = ternary
	return newe, nil
}

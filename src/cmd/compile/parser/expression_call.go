package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (expressionParser *ExpressionParser) parseCallExpression(e *ast.Expression) (*ast.Expression, error) {
	var err error
	pos := expressionParser.parser.mkPos()
	expressionParser.Next() // skip (
	args := []*ast.Expression{}
	if expressionParser.parser.token.Type != lex.TOKEN_RP { //a(123)
		args, err = expressionParser.parseExpressions()
		if err != nil {
			return nil, err
		}
	}

	if expressionParser.parser.token.Type != lex.TOKEN_RP {
		return nil, fmt.Errorf("%s except ')' ,but '%s'",
			expressionParser.parser.errorMsgPrefix(),
			expressionParser.parser.token.Description)
	}
	var result ast.Expression
	if e.Type == ast.EXPRESSION_TYPE_IDENTIFIER {
		result.Type = ast.EXPRESSION_TYPE_FUNCTION_CALL
		call := &ast.ExpressionFunctionCall{}
		call.Expression = e
		call.Args = args
		result.Data = call
		result.Pos = expressionParser.parser.mkPos()
	} else if e.Type == ast.EXPRESSION_TYPE_SELECTION {
		result.Type = ast.EXPRESSION_TYPE_METHOD_CALL
		call := &ast.ExpressionMethodCall{}
		index := e.Data.(*ast.ExpressionSelection)
		call.Expression = index.Expression
		call.Args = args
		call.Name = index.Name
		result.Data = call
		result.Pos = expressionParser.parser.mkPos()
	} else {
		return nil, fmt.Errorf("%s can`t make call on '%s'",
			expressionParser.parser.errorMsgPrefix(), e.OpName())
	}
	expressionParser.Next() // skip )
	if expressionParser.parser.token.Type == lex.TOKEN_LT {
		expressionParser.Next() // skip <
		ts, err := expressionParser.parser.parseTypes()
		if err != nil {
			expressionParser.parser.consume(untilGt)
			expressionParser.Next()
		} else {
			if expressionParser.parser.token.Type != lex.TOKEN_GT {
				expressionParser.parser.errs = append(expressionParser.parser.errs, fmt.Errorf("%s '<' and '>' not match",
					expressionParser.parser.errorMsgPrefix()))
				expressionParser.parser.consume(untilGt)
			}
			expressionParser.Next()
			if result.Type == ast.EXPRESSION_TYPE_FUNCTION_CALL {
				result.Data.(*ast.ExpressionFunctionCall).ParameterTypes = ts
			} else {
				result.Data.(*ast.ExpressionMethodCall).ParameterTypes = ts
			}
		}
	}
	result.Pos = pos
	return &result, nil
}

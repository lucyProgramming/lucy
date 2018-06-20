package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (ep *ExpressionParser) parseCallExpression(e *ast.Expression) (*ast.Expression, error) {
	var err error
	pos := ep.parser.mkPos()
	ep.Next() // skip (
	args := []*ast.Expression{}
	if ep.parser.token.Type != lex.TOKEN_RP { //a(123)
		args, err = ep.parseExpressions()
		if err != nil {
			return nil, err
		}
	}

	if ep.parser.token.Type != lex.TOKEN_RP {
		return nil, fmt.Errorf("%s except ')' ,but '%s'",
			ep.parser.errorMsgPrefix(),
			ep.parser.token.Description)
	}
	var result ast.Expression
	if e.Type == ast.EXPRESSION_TYPE_IDENTIFIER {
		result.Type = ast.EXPRESSION_TYPE_FUNCTION_CALL
		call := &ast.ExpressionFunctionCall{}
		call.Expression = e
		call.Args = args
		result.Data = call
		result.Pos = ep.parser.mkPos()
	} else if e.Type == ast.EXPRESSION_TYPE_SELECT {
		result.Type = ast.EXPRESSION_TYPE_METHOD_CALL
		call := &ast.ExpressionMethodCall{}
		index := e.Data.(*ast.ExpressionSelection)
		call.Expression = index.Expression
		call.Args = args
		call.Name = index.Name
		result.Data = call
		result.Pos = ep.parser.mkPos()
	} else {
		return nil, fmt.Errorf("%s can`t make call on '%s'",
			ep.parser.errorMsgPrefix(), e.OpName())
	}
	ep.Next() // skip )
	if ep.parser.token.Type == lex.TOKEN_LT {
		ep.Next() // skip <
		ts, err := ep.parser.parseTypes()
		if err != nil {
			ep.parser.consume(untilGt)
			ep.Next()
		} else {
			if ep.parser.token.Type != lex.TOKEN_GT {
				ep.parser.errs = append(ep.parser.errs, fmt.Errorf("%s '<' and '>' not match",
					ep.parser.errorMsgPrefix()))
				ep.parser.consume(untilGt)
			}
			ep.Next()
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

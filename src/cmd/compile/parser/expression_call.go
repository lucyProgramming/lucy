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
	if ep.parser.eof {
		return nil, ep.parser.mkUnexpectedEofErr()
	}
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
			ep.parser.token.Desp)
	}
	var result ast.Expression
	if e.Typ == ast.EXPRESSION_TYPE_IDENTIFIER {
		result.Typ = ast.EXPRESSION_TYPE_FUNCTION_CALL
		call := &ast.ExpressionFunctionCall{}
		call.Expression = e
		call.Args = args
		result.Data = call
		result.Pos = ep.parser.mkPos()
	} else if e.Typ == ast.EXPRESSION_TYPE_DOT {
		result.Typ = ast.EXPRESSION_TYPE_METHOD_CALL
		call := &ast.ExpressionMethodCall{}
		index := e.Data.(*ast.ExpressionDot)
		call.Expression = index.Expression
		call.Args = args
		call.Name = index.Name
		result.Data = call
		result.Pos = ep.parser.mkPos()
	} else {
		return nil, fmt.Errorf("%s can`t make call on '%s'", ep.parser.errorMsgPrefix())
	}
	ep.Next() // skip )
	result.Pos = pos
	return &result, nil
}

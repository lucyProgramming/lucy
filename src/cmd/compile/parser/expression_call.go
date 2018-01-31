package parser

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

func (ep *ExpressionParser) parseCallExpression(e *ast.Expression) (*ast.Expression, error) {
	var err error
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
		return nil, fmt.Errorf("%s except ')' ,but %s",
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
	} else if result.Typ == ast.EXPRESSION_TYPE_DOT {
		result.Typ = ast.EXPRESSION_TYPE_METHOD_CALL
		call := &ast.ExpressionMethodCall{}
		index := e.Data.(*ast.ExpressionIndex)
		call.Expression = index.Expression
		call.Args = args
		result.Data = call
		result.Pos = e.Pos
	} else {
		return nil, fmt.Errorf("%s can`t make call on '%s'", ep.parser.errorMsgPrefix())
	}
	ep.Next() // skip )
	return &result, nil
}

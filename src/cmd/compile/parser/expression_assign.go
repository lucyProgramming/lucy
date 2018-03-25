package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//a = 123
func (ep *ExpressionParser) parseAssignExpression() (*ast.Expression, error) {
	left, err := ep.parseLogicalExpression() //
	if err != nil {
		return nil, err
	}
	if left.Typ == ast.EXPRESSION_TYPE_LABLE {
		return left, nil
	}
	for ep.parser.token.Type == lex.TOKEN_COMMA { // read more
		ep.Next()                                 //  skip comma
		left2, err := ep.parseLogicalExpression() //
		if err != nil {
			return nil, err
		}
		if left.Typ == ast.EXPRESSION_TYPE_LIST {
			list := left.Data.([]*ast.Expression)
			left.Data = append(list, left2)
		} else {
			newe := &ast.Expression{}
			newe.Typ = ast.EXPRESSION_TYPE_LIST
			list := []*ast.Expression{left, left2}
			newe.Data = list
			left = newe
		}
	}
	mustBeOneExpression := func() {
		if left.Typ == ast.EXPRESSION_TYPE_LIST {
			es := left.Data.([]*ast.Expression)
			left = es[0]
			if len(es) > 1 {
				ep.parser.errs = append(ep.parser.errs, fmt.Errorf("%s expect one left value", ep.parser.errorMsgPrefix(es[1].Pos)))
			}
		}
	}
	mkBinayExpression := func(typ int, multi bool) (*ast.Expression, error) {
		ep.Next() // skip = :=
		if ep.parser.eof {
			return nil, ep.parser.mkUnexpectedEofErr()
		}
		result := &ast.Expression{}
		result.Typ = typ
		bin := &ast.ExpressionBinary{}
		result.Data = bin
		bin.Left = left
		result.Pos = ep.parser.mkPos()
		if multi {
			es, err := ep.parseExpressions()
			if err != nil {
				return result, err
			}
			bin.Right = &ast.Expression{}
			bin.Right.Typ = ast.EXPRESSION_TYPE_LIST
			bin.Right.Data = es
		} else {
			bin.Right, err = ep.parseExpression()
		}
		return result, err
	}
	// := += -= *= /= %=
	switch ep.parser.token.Type {
	case lex.TOKEN_ASSIGN:
		return mkBinayExpression(ast.EXPRESSION_TYPE_ASSIGN, true)
	case lex.TOKEN_COLON_ASSIGN:
		return mkBinayExpression(ast.EXPRESSION_TYPE_COLON_ASSIGN, true)
	case lex.TOKEN_ADD_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_PLUS_ASSIGN, false)
	case lex.TOKEN_SUB_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_MINUS_ASSIGN, false)
	case lex.TOKEN_MUL_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_MUL_ASSIGN, false)
	case lex.TOKEN_DIV_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_DIV_ASSIGN, false)
	case lex.TOKEN_MOD_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_MOD_ASSIGN, false)
	}
	return left, nil
}

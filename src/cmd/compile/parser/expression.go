package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type ExpressionParser struct {
	parser *Parser
}

func (ep *ExpressionParser) Next() {
	ep.parser.Next()
}

func (ep *ExpressionParser) parseExpressions() ([]*ast.Expression, error) {
	es := []*ast.Expression{}
	for !ep.parser.eof {
		e, err := ep.parseExpression(false)
		if err != nil {
			return es, err
		}
		if e.Typ == ast.EXPRESSION_TYPE_LIST {
			es = append(es, e.Data.([]*ast.Expression)...)
		} else {
			es = append(es, e)
		}
		if ep.parser.token.Type != lex.TOKEN_COMMA {
			break
		}
		// == ,
		ep.Next() // skip ,
	}
	return es, nil
}

//parse assign expression
func (ep *ExpressionParser) parseExpression(statemenLevel bool) (*ast.Expression, error) {
	left, err := ep.parseLogicalExpression(statemenLevel) //
	if err != nil {
		return nil, err
	}
	if left.Typ == ast.EXPRESSION_TYPE_LABLE {
		return left, nil
	}
	for ep.parser.token.Type == lex.TOKEN_COMMA && statemenLevel { // read more
		ep.Next()                                      //  skip comma
		left2, err := ep.parseLogicalExpression(false) //
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
				ep.parser.errs = append(ep.parser.errs, fmt.Errorf("%s expect one expression on left",
					ep.parser.errorMsgPrefix(es[1].Pos)))
			}
		}
	}
	mkBinayExpression := func(typ int, multi bool) (*ast.Expression, error) {
		pos := ep.parser.mkPos()
		ep.Next() // skip = :=
		if ep.parser.eof {
			return nil, ep.parser.mkUnexpectedEofErr()
		}
		result := &ast.Expression{}
		result.Typ = typ
		bin := &ast.ExpressionBinary{}
		result.Data = bin
		bin.Left = left
		result.Pos = pos
		if multi {
			es, err := ep.parseExpressions()
			if err != nil {
				return result, err
			}
			bin.Right = &ast.Expression{}
			bin.Right.Typ = ast.EXPRESSION_TYPE_LIST
			bin.Right.Data = es
		} else {
			bin.Right, err = ep.parseExpression(false)
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
	case lex.TOKEN_LEFT_SHIFT_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_LEFT_SHIFT_ASSIGN, false)
	case lex.TOKEN_RIGHT_SHIFT_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_RIGHT_SHIFT_ASSIGN, false)
	case lex.TOKEN_AND_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_AND_ASSIGN, false)
	case lex.TOKEN_OR_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_OR_ASSIGN, false)
	case lex.TOKEN_XOR_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_XOR_ASSIGN, false)
	}
	return left, nil
}

func (ep *ExpressionParser) parseTypeConvertionExpression() (*ast.Expression, error) {
	t, err := ep.parser.parseType()
	if err != nil {
		return nil, err
	}
	if ep.parser.token.Type != lex.TOKEN_LP {
		return nil, fmt.Errorf("%s not '(' after a type", ep.parser.errorMsgPrefix())
	}
	pos := ep.parser.mkPos()
	ep.Next() // skip (
	e, err := ep.parseExpression(false)
	if err != nil {
		return nil, err
	}
	if ep.parser.token.Type != lex.TOKEN_RP {
		return nil, fmt.Errorf("%s '(' and ')' not match", ep.parser.errorMsgPrefix())
	}
	ep.Next() // skip )
	return &ast.Expression{
		Typ: ast.EXPRESSION_TYPE_CHECK_CAST,
		Data: &ast.ExpressionTypeConvertion{
			Typ:        t,
			Expression: e,
		},
		Pos: pos,
	}, nil
}

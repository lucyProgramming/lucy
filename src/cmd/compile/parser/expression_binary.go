package parser

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//||
func (ep *ExpressionParser) parseLogicalOrExpression() (*ast.Expression, error) {
	e, err := ep.parseLogicalAndExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TOKEN_LOGICAL_OR {
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseLogicalAndExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		newe.Type = ast.EXPRESSION_TYPE_LOGICAL_OR
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// &&
func (ep *ExpressionParser) parseLogicalAndExpression() (*ast.Expression, error) {
	e, err := ep.parseOrExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TOKEN_LOGICAL_AND {
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseOrExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		newe.Type = ast.EXPRESSION_TYPE_LOGICAL_AND
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

//  |
func (ep *ExpressionParser) parseOrExpression() (*ast.Expression, error) {
	e, err := ep.parseXorExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TOKEN_OR {
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseXorExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		newe.Type = ast.EXPRESSION_TYPE_OR
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// ^
func (ep *ExpressionParser) parseXorExpression() (*ast.Expression, error) {
	e, err := ep.parseAndExpression()
	if err != nil {
		return nil, err
	}

	for ep.parser.token.Type == lex.TOKEN_XOR {
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseAndExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		newe.Type = ast.EXPRESSION_TYPE_XOR
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// &
func (ep *ExpressionParser) parseAndExpression() (*ast.Expression, error) {
	e, err := ep.parseEqualExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TOKEN_AND {
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseEqualExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		newe.Type = ast.EXPRESSION_TYPE_AND
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// == and !=
func (ep *ExpressionParser) parseEqualExpression() (*ast.Expression, error) {
	e, err := ep.parseRelationExpression()
	if err != nil {
		return nil, err
	}
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_EQUAL ||
		ep.parser.token.Type == lex.TOKEN_NE) && ep.parser.token.Type != lex.TOKEN_EOF {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseRelationExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_EQUAL {
			newe.Type = ast.EXPRESSION_TYPE_EQ
		} else {
			newe.Type = ast.EXPRESSION_TYPE_NE
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// > < >= <=
func (ep *ExpressionParser) parseRelationExpression() (*ast.Expression, error) {
	e, err := ep.parseShiftExpression()
	if err != nil {
		return nil, err
	}
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_GT || ep.parser.token.Type == lex.TOKEN_GE ||
		ep.parser.token.Type == lex.TOKEN_LT || ep.parser.token.Type == lex.TOKEN_LE) && ep.parser.token.Type != lex.TOKEN_EOF {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseShiftExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_GT {
			newe.Type = ast.EXPRESSION_TYPE_GT
		} else if typ == lex.TOKEN_GE {
			newe.Type = ast.EXPRESSION_TYPE_GE
		} else if typ == lex.TOKEN_LT {
			newe.Type = ast.EXPRESSION_TYPE_LT
		} else {
			newe.Type = ast.EXPRESSION_TYPE_LE
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// << >>
func (ep *ExpressionParser) parseShiftExpression() (*ast.Expression, error) {
	e, err := ep.parseAddExpression()
	if err != nil {
		return nil, err
	}
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_LEFT_SHIFT ||
		ep.parser.token.Type == lex.TOKEN_RIGHT_SHIFT) &&
		ep.parser.token.Type != lex.TOKEN_EOF {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseAddExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_LEFT_SHIFT {
			newe.Type = ast.EXPRESSION_TYPE_LSH
		} else {
			newe.Type = ast.EXPRESSION_TYPE_RSH
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// + -
func (ep *ExpressionParser) parseAddExpression() (*ast.Expression, error) {
	e, err := ep.parseMulExpression()
	if err != nil {
		return nil, err
	}
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_ADD || ep.parser.token.Type == lex.TOKEN_SUB) &&
		ep.parser.token.Type != lex.TOKEN_EOF {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseMulExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_ADD {
			newe.Type = ast.EXPRESSION_TYPE_ADD
		} else {
			newe.Type = ast.EXPRESSION_TYPE_SUB
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// * / %
func (ep *ExpressionParser) parseMulExpression() (*ast.Expression, error) {
	e, err := ep.parseOneExpression(false)
	if err != nil {
		return nil, err
	}
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_MUL ||
		ep.parser.token.Type == lex.TOKEN_DIV ||
		ep.parser.token.Type == lex.TOKEN_MOD) && ep.parser.token.Type != lex.TOKEN_EOF {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseOneExpression(false)
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_MUL {
			newe.Type = ast.EXPRESSION_TYPE_MUL
		} else if typ == lex.TOKEN_DIV {
			newe.Type = ast.EXPRESSION_TYPE_DIV
		} else {
			newe.Type = ast.EXPRESSION_TYPE_MOD
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

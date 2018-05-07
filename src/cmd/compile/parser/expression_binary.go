package parser

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

// && ||
func (ep *ExpressionParser) parseLogicalExpression() (*ast.Expression, error) {
	e, err := ep.parseBitANDORExpression()
	if err != nil {
		return nil, err
	}
	for (ep.parser.token.Type == lex.TOKEN_LOGICAL_AND ||
		ep.parser.token.Type == lex.TOKEN_LOGICAL_OR) && ep.parser.token.Type != lex.TOKEN_EOF {
		typ := ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseBitANDORExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_LOGICAL_AND {
			newe.Typ = ast.EXPRESSION_TYPE_LOGICAL_AND
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_LOGICAL_OR
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// & |
func (ep *ExpressionParser) parseBitANDORExpression() (*ast.Expression, error) {
	e, err := ep.parseEqualExpression()
	if err != nil {
		return nil, err
	}
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_AND || ep.parser.token.Type == lex.TOKEN_OR || ep.parser.token.Type == lex.TOKEN_XOR) &&
		ep.parser.token.Type != lex.TOKEN_EOF {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseEqualExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_AND {
			newe.Typ = ast.EXPRESSION_TYPE_AND
		} else if typ == lex.TOKEN_OR {
			newe.Typ = ast.EXPRESSION_TYPE_OR
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_XOR
		}
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
	for (ep.parser.token.Type == lex.TOKEN_EQUAL || ep.parser.token.Type == lex.TOKEN_NE) && ep.parser.token.Type != lex.TOKEN_EOF {
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
			newe.Typ = ast.EXPRESSION_TYPE_EQ
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_NE
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
			newe.Typ = ast.EXPRESSION_TYPE_GT
		} else if typ == lex.TOKEN_GE {
			newe.Typ = ast.EXPRESSION_TYPE_GE
		} else if typ == lex.TOKEN_LT {
			newe.Typ = ast.EXPRESSION_TYPE_LT
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_LE
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
		ep.parser.token.Type == lex.TOKEN_RIGHT_SHIFT) && ep.parser.token.Type != lex.TOKEN_EOF {
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
			newe.Typ = ast.EXPRESSION_TYPE_LSH
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_RSH
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
	for (ep.parser.token.Type == lex.TOKEN_ADD || ep.parser.token.Type == lex.TOKEN_SUB) && ep.parser.token.Type != lex.TOKEN_EOF {
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
			newe.Typ = ast.EXPRESSION_TYPE_ADD
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_SUB
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// */ %
func (ep *ExpressionParser) parseMulExpression() (*ast.Expression, error) {
	e, err := ep.parseOneExpression()
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
		e2, err := ep.parseOneExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_MUL {
			newe.Typ = ast.EXPRESSION_TYPE_MUL
		} else if typ == lex.TOKEN_DIV {
			newe.Typ = ast.EXPRESSION_TYPE_DIV
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_MOD
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

package parser

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//||
func (expressionParser *ExpressionParser) parseLogicalOrExpression() (*ast.Expression, error) {
	e, err := expressionParser.parseLogicalAndExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TOKEN_LOGICAL_OR {
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		e2, err := expressionParser.parseLogicalAndExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Type = ast.EXPRESSION_TYPE_LOGICAL_OR
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newExpression.Data = binary
		e = newExpression
	}
	return e, nil
}

// &&
func (expressionParser *ExpressionParser) parseLogicalAndExpression() (*ast.Expression, error) {
	e, err := expressionParser.parseOrExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TOKEN_LOGICAL_AND {
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		e2, err := expressionParser.parseOrExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Type = ast.EXPRESSION_TYPE_LOGICAL_AND
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newExpression.Data = binary
		e = newExpression
	}
	return e, nil
}

//  |
func (expressionParser *ExpressionParser) parseOrExpression() (*ast.Expression, error) {
	e, err := expressionParser.parseXorExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TOKEN_OR {
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		e2, err := expressionParser.parseXorExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Type = ast.EXPRESSION_TYPE_OR
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newExpression.Data = binary
		e = newExpression
	}
	return e, nil
}

// ^
func (expressionParser *ExpressionParser) parseXorExpression() (*ast.Expression, error) {
	e, err := expressionParser.parseAndExpression()
	if err != nil {
		return nil, err
	}

	for expressionParser.parser.token.Type == lex.TOKEN_XOR {
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		e2, err := expressionParser.parseAndExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Type = ast.EXPRESSION_TYPE_XOR
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newExpression.Data = binary
		e = newExpression
	}
	return e, nil
}

// &
func (expressionParser *ExpressionParser) parseAndExpression() (*ast.Expression, error) {
	e, err := expressionParser.parseEqualExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TOKEN_AND {
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		e2, err := expressionParser.parseEqualExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Type = ast.EXPRESSION_TYPE_AND
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newExpression.Data = binary
		e = newExpression
	}
	return e, nil
}

// == and !=
func (expressionParser *ExpressionParser) parseEqualExpression() (*ast.Expression, error) {
	e, err := expressionParser.parseRelationExpression()
	if err != nil {
		return nil, err
	}
	var typ int
	for (expressionParser.parser.token.Type == lex.TOKEN_EQUAL ||
		expressionParser.parser.token.Type == lex.TOKEN_NE) && expressionParser.parser.token.Type != lex.TOKEN_EOF {
		typ = expressionParser.parser.token.Type
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		e2, err := expressionParser.parseRelationExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		if typ == lex.TOKEN_EQUAL {
			newExpression.Type = ast.EXPRESSION_TYPE_EQ
		} else {
			newExpression.Type = ast.EXPRESSION_TYPE_NE
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newExpression.Data = binary
		e = newExpression
	}
	return e, nil
}

// > < >= <=
func (expressionParser *ExpressionParser) parseRelationExpression() (*ast.Expression, error) {
	e, err := expressionParser.parseShiftExpression()
	if err != nil {
		return nil, err
	}
	var typ int
	for (expressionParser.parser.token.Type == lex.TOKEN_GT || expressionParser.parser.token.Type == lex.TOKEN_GE ||
		expressionParser.parser.token.Type == lex.TOKEN_LT || expressionParser.parser.token.Type == lex.TOKEN_LE) && expressionParser.parser.token.Type != lex.TOKEN_EOF {
		typ = expressionParser.parser.token.Type
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		e2, err := expressionParser.parseShiftExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		if typ == lex.TOKEN_GT {
			newExpression.Type = ast.EXPRESSION_TYPE_GT
		} else if typ == lex.TOKEN_GE {
			newExpression.Type = ast.EXPRESSION_TYPE_GE
		} else if typ == lex.TOKEN_LT {
			newExpression.Type = ast.EXPRESSION_TYPE_LT
		} else {
			newExpression.Type = ast.EXPRESSION_TYPE_LE
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newExpression.Data = binary
		e = newExpression
	}
	return e, nil
}

// << >>
func (expressionParser *ExpressionParser) parseShiftExpression() (*ast.Expression, error) {
	e, err := expressionParser.parseAddExpression()
	if err != nil {
		return nil, err
	}
	var typ int
	for (expressionParser.parser.token.Type == lex.TOKEN_LEFT_SHIFT ||
		expressionParser.parser.token.Type == lex.TOKEN_RIGHT_SHIFT) &&
		expressionParser.parser.token.Type != lex.TOKEN_EOF {
		typ = expressionParser.parser.token.Type
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		e2, err := expressionParser.parseAddExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		if typ == lex.TOKEN_LEFT_SHIFT {
			newExpression.Type = ast.EXPRESSION_TYPE_LSH
		} else {
			newExpression.Type = ast.EXPRESSION_TYPE_RSH
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newExpression.Data = binary
		e = newExpression
	}
	return e, nil
}

// + -
func (expressionParser *ExpressionParser) parseAddExpression() (*ast.Expression, error) {
	e, err := expressionParser.parseMulExpression()
	if err != nil {
		return nil, err
	}
	var typ int
	for (expressionParser.parser.token.Type == lex.TOKEN_ADD || expressionParser.parser.token.Type == lex.TOKEN_SUB) &&
		expressionParser.parser.token.Type != lex.TOKEN_EOF {
		typ = expressionParser.parser.token.Type
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		e2, err := expressionParser.parseMulExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		if typ == lex.TOKEN_ADD {
			newExpression.Type = ast.EXPRESSION_TYPE_ADD
		} else {
			newExpression.Type = ast.EXPRESSION_TYPE_SUB
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newExpression.Data = binary
		e = newExpression
	}
	return e, nil
}

// * / %
func (expressionParser *ExpressionParser) parseMulExpression() (*ast.Expression, error) {
	e, err := expressionParser.parseOneExpression(false)
	if err != nil {
		return nil, err
	}
	var typ int
	for (expressionParser.parser.token.Type == lex.TOKEN_MUL ||
		expressionParser.parser.token.Type == lex.TOKEN_DIV ||
		expressionParser.parser.token.Type == lex.TOKEN_MOD) && expressionParser.parser.token.Type != lex.TOKEN_EOF {
		typ = expressionParser.parser.token.Type
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		e2, err := expressionParser.parseOneExpression(false)
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		if typ == lex.TOKEN_MUL {
			newExpression.Type = ast.EXPRESSION_TYPE_MUL
		} else if typ == lex.TOKEN_DIV {
			newExpression.Type = ast.EXPRESSION_TYPE_DIV
		} else {
			newExpression.Type = ast.EXPRESSION_TYPE_MOD
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newExpression.Data = binary
		e = newExpression
	}
	return e, nil
}

package parser

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

// ||
func (expressionParser *ExpressionParser) parseLogicalOrExpression() (*ast.Expression, error) {
	e, err := expressionParser.parseLogicalAndExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenLogicalOr {
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(false)
		e2, err := expressionParser.parseLogicalAndExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Type = ast.ExpressionTypeLogicalOr
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
	for expressionParser.parser.token.Type == lex.TokenLogicalAnd {
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(false)
		right, err := expressionParser.parseOrExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Type = ast.ExpressionTypeLogicalAnd
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = right
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
	for expressionParser.parser.token.Type == lex.TokenOr {
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(false)
		e2, err := expressionParser.parseXorExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Type = ast.ExpressionTypeOr
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

	for expressionParser.parser.token.Type == lex.TokenXor {
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(false)
		e2, err := expressionParser.parseAndExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Type = ast.ExpressionTypeXor
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
	for expressionParser.parser.token.Type == lex.TokenAnd {
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(false)
		e2, err := expressionParser.parseEqualExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Type = ast.ExpressionTypeAnd
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
	for expressionParser.parser.token.Type == lex.TokenEqual ||
		expressionParser.parser.token.Type == lex.TokenNe {
		typ = expressionParser.parser.token.Type
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(false)
		e2, err := expressionParser.parseRelationExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		if typ == lex.TokenEqual {
			newExpression.Type = ast.ExpressionTypeEq
		} else {
			newExpression.Type = ast.ExpressionTypeNe
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
	for expressionParser.parser.token.Type == lex.TokenGt ||
		expressionParser.parser.token.Type == lex.TokenGe ||
		expressionParser.parser.token.Type == lex.TokenLt ||
		expressionParser.parser.token.Type == lex.TokenLe {
		typ = expressionParser.parser.token.Type
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(false)
		e2, err := expressionParser.parseShiftExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		if typ == lex.TokenGt {
			newExpression.Type = ast.ExpressionTypeGt
		} else if typ == lex.TokenGe {
			newExpression.Type = ast.ExpressionTypeGe
		} else if typ == lex.TokenLt {
			newExpression.Type = ast.ExpressionTypeLt
		} else {
			newExpression.Type = ast.ExpressionTypeLe
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
	for expressionParser.parser.token.Type == lex.TokenLsh ||
		expressionParser.parser.token.Type == lex.TokenRsh {
		typ = expressionParser.parser.token.Type
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(false)
		e2, err := expressionParser.parseAddExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		if typ == lex.TokenLsh {
			newExpression.Type = ast.ExpressionTypeLsh
		} else {
			newExpression.Type = ast.ExpressionTypeRsh
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
	for expressionParser.parser.token.Type == lex.TokenAdd ||
		expressionParser.parser.token.Type == lex.TokenSub {
		typ = expressionParser.parser.token.Type
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(false)
		e2, err := expressionParser.parseMulExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		if typ == lex.TokenAdd {
			newExpression.Type = ast.ExpressionTypeAdd
		} else {
			newExpression.Type = ast.ExpressionTypeSub
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
	e, err := expressionParser.parseOneExpression()
	if err != nil {
		return nil, err
	}
	var typ int
	for expressionParser.parser.token.Type == lex.TokenMul ||
		expressionParser.parser.token.Type == lex.TokenDiv ||
		expressionParser.parser.token.Type == lex.TokenMod {
		typ = expressionParser.parser.token.Type
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(false)
		e2, err := expressionParser.parseOneExpression()
		if err != nil {
			return nil, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		if typ == lex.TokenMul {
			newExpression.Type = ast.ExpressionTypeMul
		} else if typ == lex.TokenDiv {
			newExpression.Type = ast.ExpressionTypeDiv
		} else {
			newExpression.Type = ast.ExpressionTypeMod
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newExpression.Data = binary
		e = newExpression
	}
	return e, nil
}

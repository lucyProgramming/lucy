package parser

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

// ||
func (expressionParser *ExpressionParser) parseLogicalOrExpression() (*ast.Expression, error) {
	left, err := expressionParser.parseLogicalAndExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenLogicalOr {
		pos := expressionParser.parser.mkPos()
		name := expressionParser.parser.token.Description
		expressionParser.Next(lfNotToken)
		right, err := expressionParser.parseLogicalAndExpression()
		if err != nil {
			return left, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Op = name
		newExpression.Type = ast.ExpressionTypeLogicalOr
		binary := &ast.ExpressionBinary{}
		binary.Left = left
		binary.Right = right
		newExpression.Data = binary
		left = newExpression
	}
	return left, nil
}

// &&
func (expressionParser *ExpressionParser) parseLogicalAndExpression() (*ast.Expression, error) {
	left, err := expressionParser.parseOrExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenLogicalAnd {
		pos := expressionParser.parser.mkPos()
		name := expressionParser.parser.token.Description
		expressionParser.Next(lfNotToken)
		right, err := expressionParser.parseOrExpression()
		if err != nil {
			return left, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Op = name
		newExpression.Type = ast.ExpressionTypeLogicalAnd
		binary := &ast.ExpressionBinary{}
		binary.Left = left
		binary.Right = right
		newExpression.Data = binary
		left = newExpression
	}
	return left, nil
}

//  |
func (expressionParser *ExpressionParser) parseOrExpression() (*ast.Expression, error) {
	left, err := expressionParser.parseXorExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenOr {
		pos := expressionParser.parser.mkPos()
		name := expressionParser.parser.token.Description
		expressionParser.Next(lfNotToken)
		right, err := expressionParser.parseXorExpression()
		if err != nil {
			return left, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Op = name
		newExpression.Type = ast.ExpressionTypeOr
		binary := &ast.ExpressionBinary{}
		binary.Left = left
		binary.Right = right
		newExpression.Data = binary
		left = newExpression
	}
	return left, nil
}

// ^
func (expressionParser *ExpressionParser) parseXorExpression() (*ast.Expression, error) {
	left, err := expressionParser.parseAndExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenXor {
		pos := expressionParser.parser.mkPos()
		name := expressionParser.parser.token.Description
		expressionParser.Next(lfNotToken)
		right, err := expressionParser.parseAndExpression()
		if err != nil {
			return left, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Op = name
		newExpression.Type = ast.ExpressionTypeXor
		binary := &ast.ExpressionBinary{}
		binary.Left = left
		binary.Right = right
		newExpression.Data = binary
		left = newExpression
	}
	return left, nil
}

// &
func (expressionParser *ExpressionParser) parseAndExpression() (*ast.Expression, error) {
	left, err := expressionParser.parseEqualExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenAnd {
		pos := expressionParser.parser.mkPos()
		name := expressionParser.parser.token.Description
		expressionParser.Next(lfNotToken)
		right, err := expressionParser.parseEqualExpression()
		if err != nil {
			return left, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Type = ast.ExpressionTypeAnd
		newExpression.Op = name
		binary := &ast.ExpressionBinary{}
		binary.Left = left
		binary.Right = right
		newExpression.Data = binary
		left = newExpression
	}
	return left, nil
}

// == and !=
func (expressionParser *ExpressionParser) parseEqualExpression() (*ast.Expression, error) {
	left, err := expressionParser.parseRelationExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenEqual ||
		expressionParser.parser.token.Type == lex.TokenNe {
		typ := expressionParser.parser.token.Type
		name := expressionParser.parser.token.Description
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfNotToken)
		right, err := expressionParser.parseRelationExpression()
		if err != nil {
			return left, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Op = name
		if typ == lex.TokenEqual {
			newExpression.Type = ast.ExpressionTypeEq
		} else {
			newExpression.Type = ast.ExpressionTypeNe
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = left
		binary.Right = right
		newExpression.Data = binary
		left = newExpression
	}
	return left, nil
}

// > < >= <=
func (expressionParser *ExpressionParser) parseRelationExpression() (*ast.Expression, error) {
	left, err := expressionParser.parseShiftExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenGt ||
		expressionParser.parser.token.Type == lex.TokenGe ||
		expressionParser.parser.token.Type == lex.TokenLt ||
		expressionParser.parser.token.Type == lex.TokenLe {
		typ := expressionParser.parser.token.Type
		name := expressionParser.parser.token.Description
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfNotToken)
		right, err := expressionParser.parseShiftExpression()
		if err != nil {
			return left, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Op = name
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
		binary.Left = left
		binary.Right = right
		newExpression.Data = binary
		left = newExpression
	}
	return left, nil
}

// << >>
func (expressionParser *ExpressionParser) parseShiftExpression() (*ast.Expression, error) {
	left, err := expressionParser.parseAddExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenLsh ||
		expressionParser.parser.token.Type == lex.TokenRsh {
		typ := expressionParser.parser.token.Type
		name := expressionParser.parser.token.Description
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfNotToken)
		right, err := expressionParser.parseAddExpression()
		if err != nil {
			return left, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Op = name
		if typ == lex.TokenLsh {
			newExpression.Type = ast.ExpressionTypeLsh

		} else {
			newExpression.Type = ast.ExpressionTypeRsh

		}
		binary := &ast.ExpressionBinary{}
		binary.Left = left
		binary.Right = right
		newExpression.Data = binary
		left = newExpression
	}
	return left, nil
}

// + -
func (expressionParser *ExpressionParser) parseAddExpression() (*ast.Expression, error) {
	left, err := expressionParser.parseMulExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenAdd ||
		expressionParser.parser.token.Type == lex.TokenSub {
		typ := expressionParser.parser.token.Type
		name := expressionParser.parser.token.Description
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfNotToken)
		right, err := expressionParser.parseMulExpression()
		if err != nil {
			return left, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Op = name
		if typ == lex.TokenAdd {
			newExpression.Type = ast.ExpressionTypeAdd
		} else {
			newExpression.Type = ast.ExpressionTypeSub
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = left
		binary.Right = right
		newExpression.Data = binary
		left = newExpression
	}
	return left, nil
}

// * / %
func (expressionParser *ExpressionParser) parseMulExpression() (*ast.Expression, error) {
	left, err := expressionParser.parseSuffixExpression()
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenMul ||
		expressionParser.parser.token.Type == lex.TokenDiv ||
		expressionParser.parser.token.Type == lex.TokenMod {
		typ := expressionParser.parser.token.Type
		name := expressionParser.parser.token.Description
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfNotToken)
		right, err := expressionParser.parseSuffixExpression()
		if err != nil {
			return left, err
		}
		newExpression := &ast.Expression{}
		newExpression.Pos = pos
		newExpression.Op = name
		if typ == lex.TokenMul {
			newExpression.Type = ast.ExpressionTypeMul
		} else if typ == lex.TokenDiv {
			newExpression.Type = ast.ExpressionTypeDiv
		} else {
			newExpression.Type = ast.ExpressionTypeMod
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = left
		binary.Right = right
		newExpression.Data = binary
		left = newExpression
	}
	return left, nil
}

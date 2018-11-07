package parser

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

// ||
func (ep *ExpressionParser) parseLogicalOrExpression() (*ast.Expression, error) {
	left, err := ep.parseLogicalAndExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TokenLogicalOr {
		pos := ep.parser.mkPos()
		name := ep.parser.token.Description
		ep.Next(lfNotToken)
		right, err := ep.parseLogicalAndExpression()
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
func (ep *ExpressionParser) parseLogicalAndExpression() (*ast.Expression, error) {
	left, err := ep.parseEqualExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TokenLogicalAnd {
		pos := ep.parser.mkPos()
		name := ep.parser.token.Description
		ep.Next(lfNotToken)
		right, err := ep.parseEqualExpression()
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

// == and !=
func (ep *ExpressionParser) parseEqualExpression() (*ast.Expression, error) {
	left, err := ep.parseRelationExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TokenEqual ||
		ep.parser.token.Type == lex.TokenNe {
		typ := ep.parser.token.Type
		name := ep.parser.token.Description
		pos := ep.parser.mkPos()
		ep.Next(lfNotToken)
		right, err := ep.parseRelationExpression()
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
func (ep *ExpressionParser) parseRelationExpression() (*ast.Expression, error) {
	left, err := ep.parseOrExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TokenGt ||
		ep.parser.token.Type == lex.TokenGe ||
		ep.parser.token.Type == lex.TokenLt ||
		ep.parser.token.Type == lex.TokenLe {
		typ := ep.parser.token.Type
		name := ep.parser.token.Description
		pos := ep.parser.mkPos()
		ep.Next(lfNotToken)
		right, err := ep.parseOrExpression()
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

//  |
func (ep *ExpressionParser) parseOrExpression() (*ast.Expression, error) {
	left, err := ep.parseXorExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TokenOr {
		pos := ep.parser.mkPos()
		name := ep.parser.token.Description
		ep.Next(lfNotToken)
		right, err := ep.parseXorExpression()
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
func (ep *ExpressionParser) parseXorExpression() (*ast.Expression, error) {
	left, err := ep.parseAndExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TokenXor {
		pos := ep.parser.mkPos()
		name := ep.parser.token.Description
		ep.Next(lfNotToken)
		right, err := ep.parseAndExpression()
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
func (ep *ExpressionParser) parseAndExpression() (*ast.Expression, error) {
	left, err := ep.parseShiftExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TokenAnd {
		pos := ep.parser.mkPos()
		name := ep.parser.token.Description
		ep.Next(lfNotToken)
		right, err := ep.parseShiftExpression()
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

// << >>
func (ep *ExpressionParser) parseShiftExpression() (*ast.Expression, error) {
	left, err := ep.parseAddExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TokenLsh ||
		ep.parser.token.Type == lex.TokenRsh {
		typ := ep.parser.token.Type
		name := ep.parser.token.Description
		pos := ep.parser.mkPos()
		ep.Next(lfNotToken)
		right, err := ep.parseAddExpression()
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
func (ep *ExpressionParser) parseAddExpression() (*ast.Expression, error) {
	left, err := ep.parseMulExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TokenAdd ||
		ep.parser.token.Type == lex.TokenSub {
		typ := ep.parser.token.Type
		name := ep.parser.token.Description
		pos := ep.parser.mkPos()
		ep.Next(lfNotToken)
		right, err := ep.parseMulExpression()
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
func (ep *ExpressionParser) parseMulExpression() (*ast.Expression, error) {
	left, err := ep.parseSuffixExpression()
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TokenMul ||
		ep.parser.token.Type == lex.TokenDiv ||
		ep.parser.token.Type == lex.TokenMod {
		typ := ep.parser.token.Type
		name := ep.parser.token.Description
		pos := ep.parser.mkPos()
		ep.Next(lfNotToken)
		right, err := ep.parseSuffixExpression()
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

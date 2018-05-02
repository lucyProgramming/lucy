package parser

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

// && ||
func (ep *ExpressionParser) parseLogicalExpression(statementLevel bool) (*ast.Expression, error) {
	e, err := ep.parseBitANDORExpression(statementLevel)
	if err != nil {
		return nil, err
	}
	if e.Typ == ast.EXPRESSION_TYPE_LABLE {
		return e, nil
	}
	for (ep.parser.token.Type == lex.TOKEN_LOGICAL_AND ||
		ep.parser.token.Type == lex.TOKEN_LOGICAL_OR) && !ep.parser.eof {
		typ := ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseBitANDORExpression(false)
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
func (ep *ExpressionParser) parseBitANDORExpression(statementLevel bool) (*ast.Expression, error) {
	e, err := ep.parseEqualExpression(statementLevel)
	if err != nil {
		return nil, err
	}
	if e.Typ == ast.EXPRESSION_TYPE_LABLE {
		return e, nil
	}
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_AND || ep.parser.token.Type == lex.TOKEN_OR || ep.parser.token.Type == lex.TOKEN_XOR) && !ep.parser.eof {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseEqualExpression(false)
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
func (ep *ExpressionParser) parseEqualExpression(statementLevel bool) (*ast.Expression, error) {
	e, err := ep.parseRelationExpression(statementLevel)
	if err != nil {
		return nil, err
	}
	if e.Typ == ast.EXPRESSION_TYPE_LABLE {
		return e, nil
	}
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_EQUAL || ep.parser.token.Type == lex.TOKEN_NE) && !ep.parser.eof {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseRelationExpression(false)
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
func (ep *ExpressionParser) parseRelationExpression(statementLevel bool) (*ast.Expression, error) {
	e, err := ep.parseShiftExpression(statementLevel)
	if err != nil {
		return nil, err
	}
	if e.Typ == ast.EXPRESSION_TYPE_LABLE {
		return e, nil
	}
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_GT || ep.parser.token.Type == lex.TOKEN_GE ||
		ep.parser.token.Type == lex.TOKEN_LT || ep.parser.token.Type == lex.TOKEN_LE) && !ep.parser.eof {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseShiftExpression(false)
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
func (ep *ExpressionParser) parseShiftExpression(statmentLevel bool) (*ast.Expression, error) {
	e, err := ep.parseAddExpression(statmentLevel)
	if err != nil {
		return nil, err
	}
	if e.Typ == ast.EXPRESSION_TYPE_LABLE {
		return e, nil
	}
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_LEFT_SHIFT ||
		ep.parser.token.Type == lex.TOKEN_RIGHT_SHIFT) && !ep.parser.eof {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseAddExpression(false)
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_LEFT_SHIFT {
			newe.Typ = ast.EXPRESSION_TYPE_LEFT_SHIFT
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_RIGHT_SHIFT
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
func (ep *ExpressionParser) parseAddExpression(statemenLevel bool) (*ast.Expression, error) {
	e, err := ep.parseMulExpression(statemenLevel)
	if err != nil {
		return nil, err
	}
	if e.Typ == ast.EXPRESSION_TYPE_LABLE {
		return e, nil
	}
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_ADD || ep.parser.token.Type == lex.TOKEN_SUB) && !ep.parser.eof {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseMulExpression(false)
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
func (ep *ExpressionParser) parseMulExpression(statementLevel bool) (*ast.Expression, error) {
	e, err := ep.parseOneExpression(statementLevel)
	if err != nil {
		return nil, err
	}
	if e.Typ == ast.EXPRESSION_TYPE_LABLE {
		return e, nil
	}
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_MUL ||
		ep.parser.token.Type == lex.TOKEN_DIV ||
		ep.parser.token.Type == lex.TOKEN_MOD) && !ep.parser.eof {
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

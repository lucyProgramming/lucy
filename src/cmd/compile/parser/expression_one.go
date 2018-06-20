package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (ep *ExpressionParser) parseOneExpression(unary bool) (*ast.Expression, error) {
	var left *ast.Expression
	var err error
	switch ep.parser.token.Type {
	case lex.TOKEN_IDENTIFIER:
		left = &ast.Expression{}
		left.Type = ast.EXPRESSION_TYPE_IDENTIFIER
		identifier := &ast.ExpressionIdentifier{}
		identifier.Name = ep.parser.token.Data.(string)
		left.Data = identifier
		left.Pos = ep.parser.mkPos()
		ep.Next()
	case lex.TOKEN_TRUE:
		left = &ast.Expression{}
		left.Type = ast.EXPRESSION_TYPE_BOOL
		left.Data = true
		left.Pos = ep.parser.mkPos()
		ep.Next()
	case lex.TOKEN_FALSE:
		left = &ast.Expression{}
		left.Type = ast.EXPRESSION_TYPE_BOOL
		left.Data = false
		left.Pos = ep.parser.mkPos()
		ep.Next()
	case lex.TOKEN_LITERAL_BYTE:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_BYTE,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_LITERAL_SHORT:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_SHORT,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_LITERAL_INT:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_INT,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_LITERAL_LONG:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_LONG,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_LITERAL_FLOAT:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_FLOAT,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_LITERAL_DOUBLE:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_DOUBLE,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_LITERAL_STRING:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_STRING,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_NULL:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_NULL,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_LP:
		ep.Next()
		left, err = ep.parseExpression(false)
		if err != nil {
			return nil, err
		}
		if ep.parser.token.Type != lex.TOKEN_RP {
			return nil, fmt.Errorf("%s '(' and ')' not matched, but '%s'",
				ep.parser.errorMsgPrefix(), ep.parser.token.Description)
		}
		ep.Next()
	case lex.TOKEN_INCREMENT:
		pos := ep.parser.mkPos()
		ep.Next() // skip ++
		newE := &ast.Expression{}
		newE.Pos = pos
		left, err = ep.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.EXPRESSION_TYPE_PRE_INCREMENT
		newE.Data = left
		left = newE
	case lex.TOKEN_DECREMENT:
		pos := ep.parser.mkPos()
		ep.Next() // skip --
		newE := &ast.Expression{}
		left, err = ep.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.EXPRESSION_TYPE_PRE_DECREMENT
		newE.Data = left
		newE.Pos = pos
		left = newE
	case lex.TOKEN_NOT:
		pos := ep.parser.mkPos()
		ep.Next()
		newE := &ast.Expression{}
		left, err = ep.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.EXPRESSION_TYPE_NOT
		newE.Data = left
		newE.Pos = pos
		left = newE
	case lex.TOKEN_BITWISE_NOT:
		pos := ep.parser.mkPos()
		ep.Next()
		newE := &ast.Expression{}
		left, err = ep.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.EXPRESSION_TYPE_BITWISE_NOT
		newE.Data = left
		newE.Pos = pos
		left = newE
	case lex.TOKEN_SUB:
		pos := ep.parser.mkPos()
		ep.Next()
		newE := &ast.Expression{}
		left, err = ep.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.EXPRESSION_TYPE_NEGATIVE
		newE.Data = left
		newE.Pos = pos
		left = newE
	// case lex.TOKEN_FUNCTION:
	// 	f, err := ep.parser.Function.parse(false)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return &ast.Expression{
	// 		Typ:  ast.EXPRESSION_TYPE_FUNCTION,
	// 		Data: f,
	// 	}, nil
	case lex.TOKEN_NEW:
		pos := ep.parser.mkPos()
		ep.Next()
		t, err := ep.parser.parseType()
		if err != nil {
			return nil, err
		}
		if ep.parser.token.Type != lex.TOKEN_LP {
			return nil, fmt.Errorf("%s missing '(' after new", ep.parser.errorMsgPrefix())
		}
		ep.Next() // skip (
		var es []*ast.Expression
		if ep.parser.token.Type != lex.TOKEN_RP { //
			es, err = ep.parseExpressions()
			if err != nil {
				return nil, err
			}
		}
		if ep.parser.token.Type != lex.TOKEN_RP {
			return nil, fmt.Errorf("%s ( and ) not match", ep.parser.errorMsgPrefix())
		}
		ep.Next()
		left = &ast.Expression{
			Pos:  pos,
			Type: ast.EXPRESSION_TYPE_NEW,
			Data: &ast.ExpressionNew{
				Args: es,
				Type: t,
			},
		}

	case lex.TOKEN_LB:
		left, err = ep.parseArrayExpression()
		if err != nil {
			return left, err
		}
	// bool(xxx)
	case lex.TOKEN_BOOL:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
		//
	case lex.TOKEN_BYTE:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_SHORT:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_INT:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_LONG:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_FLOAT:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_DOUBLE:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_STRING:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_T:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_RANGE:
		pos := ep.parser.mkPos()
		ep.Next()
		e, err := ep.parseExpression(false)
		if err != nil {
			return nil, err
		}
		left = &ast.Expression{}
		left.Type = ast.EXPRESSION_TYPE_RANGE
		left.Pos = pos
		left.Data = e
		return left, nil
	case lex.TOKEN_MAP:
		left, err = ep.parseMapExpression(true)
		if err != nil {
			return left, err
		}
	case lex.TOKEN_LC:
		left, err = ep.parseMapExpression(false)
		if err != nil {
			return left, err
		}
	default:
		err = fmt.Errorf("%s unkown begining of a expression, token:%s",
			ep.parser.errorMsgPrefix(), ep.parser.token.Description)
		return nil, err
	}

	for ep.parser.token.Type == lex.TOKEN_INCREMENT ||
		ep.parser.token.Type == lex.TOKEN_DECREMENT ||
		ep.parser.token.Type == lex.TOKEN_LP ||
		ep.parser.token.Type == lex.TOKEN_LB ||
		ep.parser.token.Type == lex.TOKEN_DOT {
		// ++ or --
		if ep.parser.token.Type == lex.TOKEN_INCREMENT ||
			ep.parser.token.Type == lex.TOKEN_DECREMENT { //  ++ or --
			if unary {
				return left, nil
			}
			newExpression := &ast.Expression{}
			if ep.parser.token.Type == lex.TOKEN_INCREMENT {
				newExpression.Type = ast.EXPRESSION_TYPE_INCREMENT
			} else {
				newExpression.Type = ast.EXPRESSION_TYPE_DECREMENT
			}
			if left.Type != ast.EXPRESSION_TYPE_LIST {
				newExpression.Data = left
				left = newExpression
			} else {
				list := left.Data.([]*ast.Expression)
				newExpression.Data = list[len(list)-1]
			}
			newExpression.Pos = ep.parser.mkPos()
			ep.Next()
			continue
		}
		// [
		if ep.parser.token.Type == lex.TOKEN_LB {
			pos := ep.parser.mkPos()
			ep.Next()                                    // skip [
			if ep.parser.token.Type == lex.TOKEN_COLON { // a[:]
				ep.Next() // skip :
				var end *ast.Expression
				if ep.parser.token.Type != lex.TOKEN_RB {
					end, err = ep.parseExpression(false)
					if err != nil {
						return nil, err
					}
				}
				if ep.parser.token.Type != lex.TOKEN_RB {
					return nil, fmt.Errorf("%s '[' and ']' not match", ep.parser.errorMsgPrefix())
				}
				ep.Next() // skip ]
				newExpression := &ast.Expression{}
				newExpression.Type = ast.EXPRESSION_TYPE_SLICE
				newExpression.Pos = ep.parser.mkPos()
				slice := &ast.ExpressionSlice{}
				newExpression.Data = slice
				slice.Array = left
				slice.End = end
				left = newExpression
				continue
			}
			e, err := ep.parseExpression(false)
			if err != nil {
				return nil, err
			}
			if ep.parser.token.Type == lex.TOKEN_COLON {
				ep.parser.Next()
				if ep.parser.token.Type == lex.TOKEN_COLON {
					ep.parser.Next() // skip :
				}
				var end *ast.Expression
				if ep.parser.token.Type != lex.TOKEN_RB {
					end, err = ep.parseExpression(false)
					if err != nil {
						return nil, err
					}
				}
				if ep.parser.token.Type != lex.TOKEN_RB {
					return nil, fmt.Errorf("%s '[' and ']' not match", ep.parser.errorMsgPrefix())
				}
				ep.Next() // skip ]
				newExpression := &ast.Expression{}
				newExpression.Type = ast.EXPRESSION_TYPE_SLICE
				newExpression.Pos = ep.parser.mkPos()
				slice := &ast.ExpressionSlice{}
				newExpression.Data = slice
				slice.Start = e
				slice.Array = left
				slice.End = end
				left = newExpression
				continue
			}
			if ep.parser.token.Type != lex.TOKEN_RB {
				return nil, fmt.Errorf("%s '[' and ']' not match", ep.parser.errorMsgPrefix())
			}
			newExpression := &ast.Expression{}
			newExpression.Pos = pos
			newExpression.Type = ast.EXPRESSION_TYPE_INDEX
			index := &ast.ExpressionIndex{}
			index.Expression = left
			index.Index = e
			newExpression.Data = index
			left = newExpression
			ep.Next()
			continue
		}
		// aaa.xxxx
		if ep.parser.token.Type == lex.TOKEN_DOT {
			pos := ep.parser.mkPos()
			ep.parser.Next() // skip .
			if ep.parser.token.Type == lex.TOKEN_IDENTIFIER {
				newExpression := &ast.Expression{}
				newExpression.Pos = pos
				newExpression.Type = ast.EXPRESSION_TYPE_SELECT
				index := &ast.ExpressionSelection{}
				index.Expression = left
				index.Name = ep.parser.token.Data.(string)
				newExpression.Data = index
				left = newExpression
				ep.Next()
			} else if ep.parser.token.Type == lex.TOKEN_LP { //  a.(xxx)
				//
				ep.Next() // skip (
				typ, err := ep.parser.parseType()
				if err != nil {
					return nil, err
				}
				if ep.parser.token.Type != lex.TOKEN_RP {
					return nil, fmt.Errorf("%s '(' and ')' not match", ep.parser.errorMsgPrefix())
				}
				ep.Next() // skip  )
				newExpression := &ast.Expression{}
				newExpression.Pos = pos
				newExpression.Type = ast.EXPRESSION_TYPE_TYPE_ASSERT
				typeAssert := &ast.ExpressionTypeAssert{}
				typeAssert.Type = typ
				typeAssert.Expression = left
				newExpression.Data = typeAssert
				left = newExpression
			} else {
				return nil, fmt.Errorf("%s expect  'identifier' or '(',but '%s'",
					ep.parser.errorMsgPrefix(), ep.parser.token.Description)
			}
			continue
		}
		// aa()
		if ep.parser.token.Type == lex.TOKEN_LP {
			newExpression, err := ep.parseCallExpression(left)
			if err != nil {
				return nil, err
			}
			left = newExpression
			continue
		}
	}
	return left, nil
}

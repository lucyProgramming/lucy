package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (expressionParser *ExpressionParser) parseOneExpression(unary bool) (*ast.Expression, error) {
	var left *ast.Expression
	var err error
	switch expressionParser.parser.token.Type {
	case lex.TOKEN_IDENTIFIER:
		left = &ast.Expression{}
		left.Type = ast.EXPRESSION_TYPE_IDENTIFIER
		identifier := &ast.ExpressionIdentifier{}
		identifier.Name = expressionParser.parser.token.Data.(string)
		left.Data = identifier
		left.Pos = expressionParser.parser.mkPos()
		expressionParser.Next()
	case lex.TOKEN_TRUE:
		left = &ast.Expression{}
		left.Type = ast.EXPRESSION_TYPE_BOOL
		left.Data = true
		left.Pos = expressionParser.parser.mkPos()
		expressionParser.Next()
	case lex.TOKEN_FALSE:
		left = &ast.Expression{}
		left.Type = ast.EXPRESSION_TYPE_BOOL
		left.Data = false
		left.Pos = expressionParser.parser.mkPos()
		expressionParser.Next()
	case lex.TOKEN_LITERAL_BYTE:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_BYTE,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TOKEN_LITERAL_SHORT:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_SHORT,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TOKEN_LITERAL_INT:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_INT,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TOKEN_LITERAL_LONG:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_LONG,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TOKEN_LITERAL_FLOAT:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_FLOAT,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TOKEN_LITERAL_DOUBLE:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_DOUBLE,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TOKEN_LITERAL_STRING:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_STRING,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TOKEN_NULL:
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_NULL,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TOKEN_LP:
		expressionParser.Next()
		left, err = expressionParser.parseExpression(false)
		if err != nil {
			return nil, err
		}
		if expressionParser.parser.token.Type != lex.TOKEN_RP {
			return nil, fmt.Errorf("%s '(' and ')' not matched, but '%s'",
				expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
		}
		expressionParser.Next()
	case lex.TOKEN_INCREMENT:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next() // skip ++
		newE := &ast.Expression{}
		newE.Pos = pos
		left, err = expressionParser.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.EXPRESSION_TYPE_PRE_INCREMENT
		newE.Data = left
		left = newE
	case lex.TOKEN_DECREMENT:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next() // skip --
		newE := &ast.Expression{}
		left, err = expressionParser.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.EXPRESSION_TYPE_PRE_DECREMENT
		newE.Data = left
		newE.Pos = pos
		left = newE
	case lex.TOKEN_NOT:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		newE := &ast.Expression{}
		left, err = expressionParser.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.EXPRESSION_TYPE_NOT
		newE.Data = left
		newE.Pos = pos
		left = newE
	case lex.TOKEN_BITWISE_NOT:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		newE := &ast.Expression{}
		left, err = expressionParser.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.EXPRESSION_TYPE_BIT_NOT
		newE.Data = left
		newE.Pos = pos
		left = newE
	case lex.TOKEN_SUB:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		newE := &ast.Expression{}
		left, err = expressionParser.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.EXPRESSION_TYPE_NEGATIVE
		newE.Data = left
		newE.Pos = pos
		left = newE
	case lex.TOKEN_FUNCTION:
		pos := expressionParser.parser.mkPos()
		f, err := expressionParser.parser.FunctionParser.parse(false)
		if err != nil {
			return nil, err
		}
		left = &ast.Expression{
			Type: ast.EXPRESSION_TYPE_FUNCTION_LITERAL,
			Data: f,
			Pos:  pos,
		}
	case lex.TOKEN_NEW:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		t, err := expressionParser.parser.parseType()
		if err != nil {
			return nil, err
		}
		if expressionParser.parser.token.Type != lex.TOKEN_LP {
			return nil, fmt.Errorf("%s missing '(' after new", expressionParser.parser.errorMsgPrefix())
		}
		expressionParser.Next() // skip (
		var es []*ast.Expression
		if expressionParser.parser.token.Type != lex.TOKEN_RP { //
			es, err = expressionParser.parseExpressions()
			if err != nil {
				return nil, err
			}
		}
		if expressionParser.parser.token.Type != lex.TOKEN_RP {
			return nil, fmt.Errorf("%s ( and ) not match", expressionParser.parser.errorMsgPrefix())
		}
		expressionParser.Next()
		left = &ast.Expression{
			Pos:  pos,
			Type: ast.EXPRESSION_TYPE_NEW,
			Data: &ast.ExpressionNew{
				Args: es,
				Type: t,
			},
		}

	case lex.TOKEN_LB:
		left, err = expressionParser.parseArrayExpression()
		if err != nil {
			return left, err
		}
	// bool(xxx)
	case lex.TOKEN_BOOL:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
		//
	case lex.TOKEN_BYTE:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_SHORT:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_INT:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_LONG:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_FLOAT:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_DOUBLE:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_STRING:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_T:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_RANGE:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		e, err := expressionParser.parseExpression(false)
		if err != nil {
			return nil, err
		}
		left = &ast.Expression{}
		left.Type = ast.EXPRESSION_TYPE_RANGE
		left.Pos = pos
		left.Data = e
		return left, nil
	case lex.TOKEN_MAP:
		left, err = expressionParser.parseMapExpression(true)
		if err != nil {
			return left, err
		}
	case lex.TOKEN_LC:
		left, err = expressionParser.parseMapExpression(false)
		if err != nil {
			return left, err
		}
	default:
		err = fmt.Errorf("%s unkown begining of a expression, token:%s",
			expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
		return nil, err
	}

	for expressionParser.parser.token.Type == lex.TOKEN_INCREMENT ||
		expressionParser.parser.token.Type == lex.TOKEN_DECREMENT ||
		expressionParser.parser.token.Type == lex.TOKEN_LP ||
		expressionParser.parser.token.Type == lex.TOKEN_LB ||
		expressionParser.parser.token.Type == lex.TOKEN_DOT {
		// ++ or --
		if expressionParser.parser.token.Type == lex.TOKEN_INCREMENT ||
			expressionParser.parser.token.Type == lex.TOKEN_DECREMENT { //  ++ or --
			if unary {
				return left, nil
			}
			newExpression := &ast.Expression{}
			if expressionParser.parser.token.Type == lex.TOKEN_INCREMENT {
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
			newExpression.Pos = expressionParser.parser.mkPos()
			expressionParser.Next()
			continue
		}
		// [
		if expressionParser.parser.token.Type == lex.TOKEN_LB {
			pos := expressionParser.parser.mkPos()
			expressionParser.Next()                                    // skip [
			if expressionParser.parser.token.Type == lex.TOKEN_COLON { // a[:]
				expressionParser.Next() // skip :
				var end *ast.Expression
				if expressionParser.parser.token.Type != lex.TOKEN_RB {
					end, err = expressionParser.parseExpression(false)
					if err != nil {
						return nil, err
					}
				}
				if expressionParser.parser.token.Type != lex.TOKEN_RB {
					return nil, fmt.Errorf("%s '[' and ']' not match", expressionParser.parser.errorMsgPrefix())
				}
				expressionParser.Next() // skip ]
				newExpression := &ast.Expression{}
				newExpression.Type = ast.EXPRESSION_TYPE_SLICE
				newExpression.Pos = expressionParser.parser.mkPos()
				slice := &ast.ExpressionSlice{}
				newExpression.Data = slice
				slice.Array = left
				slice.End = end
				left = newExpression
				continue
			}
			e, err := expressionParser.parseExpression(false)
			if err != nil {
				return nil, err
			}
			if expressionParser.parser.token.Type == lex.TOKEN_COLON {
				expressionParser.parser.Next()
				if expressionParser.parser.token.Type == lex.TOKEN_COLON {
					expressionParser.parser.Next() // skip :
				}
				var end *ast.Expression
				if expressionParser.parser.token.Type != lex.TOKEN_RB {
					end, err = expressionParser.parseExpression(false)
					if err != nil {
						return nil, err
					}
				}
				if expressionParser.parser.token.Type != lex.TOKEN_RB {
					return nil, fmt.Errorf("%s '[' and ']' not match", expressionParser.parser.errorMsgPrefix())
				}
				expressionParser.Next() // skip ]
				newExpression := &ast.Expression{}
				newExpression.Type = ast.EXPRESSION_TYPE_SLICE
				newExpression.Pos = expressionParser.parser.mkPos()
				slice := &ast.ExpressionSlice{}
				newExpression.Data = slice
				slice.Start = e
				slice.Array = left
				slice.End = end
				left = newExpression
				continue
			}
			if expressionParser.parser.token.Type != lex.TOKEN_RB {
				return nil, fmt.Errorf("%s '[' and ']' not match", expressionParser.parser.errorMsgPrefix())
			}
			newExpression := &ast.Expression{}
			newExpression.Pos = pos
			newExpression.Type = ast.EXPRESSION_TYPE_INDEX
			index := &ast.ExpressionIndex{}
			index.Expression = left
			index.Index = e
			newExpression.Data = index
			left = newExpression
			expressionParser.Next()
			continue
		}
		// aaa.xxxx
		if expressionParser.parser.token.Type == lex.TOKEN_DOT {
			pos := expressionParser.parser.mkPos()
			expressionParser.parser.Next() // skip .
			if expressionParser.parser.token.Type == lex.TOKEN_IDENTIFIER {
				newExpression := &ast.Expression{}
				newExpression.Pos = pos
				newExpression.Type = ast.EXPRESSION_TYPE_SELECTION
				index := &ast.ExpressionSelection{}
				index.Expression = left
				index.Name = expressionParser.parser.token.Data.(string)
				newExpression.Data = index
				left = newExpression
				expressionParser.Next()
			} else if expressionParser.parser.token.Type == lex.TOKEN_LP { //  a.(xxx)
				//
				expressionParser.Next() // skip (
				typ, err := expressionParser.parser.parseType()
				if err != nil {
					return nil, err
				}
				if expressionParser.parser.token.Type != lex.TOKEN_RP {
					return nil, fmt.Errorf("%s '(' and ')' not match", expressionParser.parser.errorMsgPrefix())
				}
				expressionParser.Next() // skip  )
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
					expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
			}
			continue
		}
		// aa()
		if expressionParser.parser.token.Type == lex.TOKEN_LP {
			newExpression, err := expressionParser.parseCallExpression(left)
			if err != nil {
				return nil, err
			}
			left = newExpression
			continue
		}
	}
	return left, nil
}

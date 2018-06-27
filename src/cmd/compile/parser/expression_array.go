package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

// []int{1,2,3}
func (expressionParser *ExpressionParser) parseArrayExpression() (*ast.Expression, error) {
	pos := expressionParser.parser.mkPos()
	expressionParser.parser.Next() // skip [
	var err error
	if expressionParser.parser.token.Type != lex.TOKEN_RB {
		arr := &ast.ExpressionArray{}
		arr.Expressions, err = expressionParser.parseExpressions()
		if expressionParser.parser.token.Type != lex.TOKEN_RB {
			err = fmt.Errorf("%s '[' and ']' not match", expressionParser.parser.errorMsgPrefix())
		} else {
			expressionParser.Next() // skip ]
		}
		return &ast.Expression{
			Type: ast.EXPRESSION_TYPE_ARRAY,
			Data: arr,
			Pos:  pos,
		}, err
	}
	expressionParser.Next() // skip [
	t, err := expressionParser.parser.parseType()
	if err != nil {
		return nil, err
	}
	if expressionParser.parser.token.Type == lex.TOKEN_LP { // []byte("1111111111")
		expressionParser.Next() // skip (
		e, err := expressionParser.parseExpression(false)
		if err != nil {
			return nil, err
		}
		if expressionParser.parser.token.Type != lex.TOKEN_RP {
			return nil, fmt.Errorf("%s '(' and  ')' not match",
				expressionParser.parser.errorMsgPrefix())
		}
		expressionParser.Next() // skip )
		ret := &ast.Expression{}
		ret.Pos = pos
		ret.Type = ast.EXPRESSION_TYPE_CHECK_CAST
		data := &ast.ExpressionTypeConversion{}
		data.Type = &ast.Type{}
		data.Type.Type = ast.VariableTypeArray
		data.Type.Pos = pos
		data.Type.ArrayType = t
		data.Expression = e
		ret.Data = data
		return ret, nil
	}

	arr := &ast.ExpressionArray{}
	if t != nil {
		arr.Type = &ast.Type{}
		arr.Type.Type = ast.VariableTypeArray
		arr.Type.ArrayType = t
		arr.Type.Pos = pos
	}
	arr.Expressions, err = expressionParser.parseArrayValues()
	return &ast.Expression{
		Type: ast.EXPRESSION_TYPE_ARRAY,
		Data: arr,
		Pos:  pos,
	}, err

}

//{1,2,3}  {{1,2,3},{456}}
func (expressionParser *ExpressionParser) parseArrayValues() ([]*ast.Expression, error) {
	if expressionParser.parser.token.Type != lex.TOKEN_LC {
		return nil, fmt.Errorf("%s expect '{',but '%s'",
			expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
	}
	expressionParser.Next() // skip {
	es := []*ast.Expression{}
	for expressionParser.parser.token.Type != lex.TOKEN_EOF && expressionParser.parser.token.Type != lex.TOKEN_RC {
		if expressionParser.parser.token.Type == lex.TOKEN_LC {
			ees, err := expressionParser.parseArrayValues()
			if err != nil {
				return es, err
			}
			arrayExpression := &ast.Expression{Type: ast.EXPRESSION_TYPE_ARRAY}
			data := ast.ExpressionArray{}
			data.Expressions = ees
			arrayExpression.Data = data
			es = append(es, arrayExpression)
		} else {
			e, err := expressionParser.parseExpression(false)
			if e != nil {
				if e.Type == ast.EXPRESSION_TYPE_LIST {
					es = append(es, e.Data.([]*ast.Expression)...)
				} else {
					es = append(es, e)
				}
			}
			if err != nil {
				return es, err
			}
		}
		if expressionParser.parser.token.Type == lex.TOKEN_COMMA {
			expressionParser.Next() // skip ,
		} else {
			break
		}
	}
	if expressionParser.parser.token.Type != lex.TOKEN_RC {
		return es, fmt.Errorf("%s expect '}',but '%s'", expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
	}
	expressionParser.Next()
	return es, nil
}

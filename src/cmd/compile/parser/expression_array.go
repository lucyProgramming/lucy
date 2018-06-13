package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

// []int{1,2,3}
func (ep *ExpressionParser) parseArrayExpression() (*ast.Expression, error) {
	pos := ep.parser.mkPos()
	ep.parser.Next() // skip [
	var err error
	if ep.parser.token.Type != lex.TOKEN_RB {
		arr := &ast.ExpressionArray{}
		arr.Expressions, err = ep.parseExpressions()
		if ep.parser.token.Type != lex.TOKEN_RB {
			err = fmt.Errorf("%s '[' and ']' not match", ep.parser.errorMsgPrefix())
		} else {
			ep.Next() // skip ]
		}
		return &ast.Expression{
			Type: ast.EXPRESSION_TYPE_ARRAY,
			Data: arr,
			Pos:  pos,
		}, err
	}
	ep.Next() // skip [
	t, err := ep.parser.parseType()
	if err != nil {
		return nil, err
	}
	if ep.parser.token.Type == lex.TOKEN_LP { // []byte("1111111111")
		ep.Next() // skip (
		e, err := ep.parseExpression(false)
		if err != nil {
			return nil, err
		}
		if ep.parser.token.Type != lex.TOKEN_RP {
			return nil, fmt.Errorf("%s '(' and  ')' not match",
				ep.parser.errorMsgPrefix())
		}
		ep.Next() // skip )
		ret := &ast.Expression{}
		ret.Pos = pos
		ret.Type = ast.EXPRESSION_TYPE_CHECK_CAST
		data := &ast.ExpressionTypeConversion{}
		data.Type = &ast.VariableType{}
		data.Type.Type = ast.VARIABLE_TYPE_ARRAY
		data.Type.Pos = pos
		data.Type.ArrayType = t
		data.Expression = e
		ret.Data = data
		return ret, nil
	}

	arr := &ast.ExpressionArray{}
	if t != nil {
		arr.Type = &ast.VariableType{}
		arr.Type.Type = ast.VARIABLE_TYPE_ARRAY
		arr.Type.ArrayType = t
		arr.Type.Pos = pos
	}
	arr.Expressions, err = ep.parseArrayValues()
	return &ast.Expression{
		Type: ast.EXPRESSION_TYPE_ARRAY,
		Data: arr,
		Pos:  pos,
	}, err

}

//{1,2,3}  {{1,2,3},{456}}
func (ep *ExpressionParser) parseArrayValues() ([]*ast.Expression, error) {
	if ep.parser.token.Type != lex.TOKEN_LC {
		return nil, fmt.Errorf("%s expect '{',but '%s'",
			ep.parser.errorMsgPrefix(), ep.parser.token.Description)
	}
	ep.Next() // skip {
	es := []*ast.Expression{}
	for ep.parser.token.Type != lex.TOKEN_EOF && ep.parser.token.Type != lex.TOKEN_RC {
		if ep.parser.token.Type == lex.TOKEN_LC {
			ees, err := ep.parseArrayValues()
			if err != nil {
				return es, err
			}
			arrayExpression := &ast.Expression{Type: ast.EXPRESSION_TYPE_ARRAY}
			data := ast.ExpressionArray{}
			data.Expressions = ees
			arrayExpression.Data = data
			es = append(es, arrayExpression)
		} else {
			e, err := ep.parseExpression(false)
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
		if ep.parser.token.Type == lex.TOKEN_COMMA {
			ep.Next() // skip ,
		} else {
			break
		}
	}
	if ep.parser.token.Type != lex.TOKEN_RC {
		return es, fmt.Errorf("%s expect '}',but '%s'", ep.parser.errorMsgPrefix(), ep.parser.token.Description)
	}
	ep.Next()
	return es, nil
}

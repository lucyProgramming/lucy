package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (ep *ExpressionParser) parseMapExprssion(needType bool) (*ast.Expression, error) {
	var typ *ast.VariableType
	var err error
	if needType {
		typ, err = ep.parser.parseType()
		if err != nil {
			return nil, err
		}
	}
	if ep.parser.token.Type != lex.TOKEN_LC {
		return nil, fmt.Errorf("expect '{',but '%s'", ep.parser.token.Description)
	}
	ep.Next() // skip {
	ret := &ast.Expression{Type: ast.EXPRESSION_TYPE_MAP}
	m := &ast.ExpressionMap{}
	m.Type = typ
	ret.Data = m
	for ep.parser.token.Type != lex.TOKEN_EOF && ep.parser.token.Type != lex.TOKEN_RC {
		// key
		k, err := ep.parseExpression(false)
		if err != nil {
			return ret, err
		}
		// arrow
		if ep.parser.token.Type != lex.TOKEN_ARROW {
			return ret, fmt.Errorf("%s expect '->',but '%s'",
				ep.parser.errorMsgPrefix(), ep.parser.token.Description)
		}
		ep.Next()
		// value
		v, err := ep.parseExpression(false)
		if err != nil {
			return ret, err
		}
		m.KeyValuePairs = append(m.KeyValuePairs, &ast.ExpressionBinary{
			Left:  k,
			Right: v,
		})
		if ep.parser.token.Type == lex.TOKEN_COMMA {
			ep.Next() // read next  key value pair
		} else {
			break
		}
	}
	if ep.parser.token.Type != lex.TOKEN_RC {
		return nil, fmt.Errorf("%s expect '}',but '%s'",
			ep.parser.errorMsgPrefix(), ep.parser.token.Description)
	}
	ep.Next() // skip }
	return ret, nil
}

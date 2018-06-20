package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (expressionParser *ExpressionParser) parseMapExpression(needType bool) (*ast.Expression, error) {
	var typ *ast.Type
	var err error
	if needType {
		typ, err = expressionParser.parser.parseType()
		if err != nil {
			return nil, err
		}
	}
	if expressionParser.parser.token.Type != lex.TOKEN_LC {
		return nil, fmt.Errorf("expect '{',but '%s'", expressionParser.parser.token.Description)
	}
	expressionParser.Next() // skip {
	ret := &ast.Expression{Type: ast.EXPRESSION_TYPE_MAP}
	m := &ast.ExpressionMap{}
	m.Type = typ
	ret.Data = m
	for expressionParser.parser.token.Type != lex.TOKEN_EOF && expressionParser.parser.token.Type != lex.TOKEN_RC {
		// key
		k, err := expressionParser.parseExpression(false)
		if err != nil {
			return ret, err
		}
		// arrow
		if expressionParser.parser.token.Type != lex.TOKEN_ARROW {
			return ret, fmt.Errorf("%s expect '->',but '%s'",
				expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
		}
		expressionParser.Next()
		// value
		v, err := expressionParser.parseExpression(false)
		if err != nil {
			return ret, err
		}
		m.KeyValuePairs = append(m.KeyValuePairs, &ast.ExpressionBinary{
			Left:  k,
			Right: v,
		})
		if expressionParser.parser.token.Type == lex.TOKEN_COMMA {
			expressionParser.Next() // read next  key value pair
		} else {
			break
		}
	}
	if expressionParser.parser.token.Type != lex.TOKEN_RC {
		return nil, fmt.Errorf("%s expect '}',but '%s'",
			expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
	}
	expressionParser.Next() // skip }
	return ret, nil
}

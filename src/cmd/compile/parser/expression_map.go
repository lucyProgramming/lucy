package parser

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
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
		return nil, fmt.Errorf("expect '{',but '%s'", ep.parser.token.Desp)
	}
	ep.Next() // skip {

	ret := &ast.Expression{Typ: ast.EXPRESSION_TYPE_MAP}
	m := &ast.ExpressionMap{}
	m.Typ = typ
	ret.Data = m
	for ep.parser.eof == false && ep.parser.token.Type != lex.TOKEN_RC {
		// key
		k, err := ep.parseExpression()
		if err != nil {
			return ret, err
		}
		// arrow
		if ep.parser.token.Type != lex.TOKEN_ARROW {
			return ret, fmt.Errorf("expect '->',but '%s'", ep.parser.token.Desp)
		}
		ep.Next()
		// value
		v, err := ep.parseExpression()
		if err != nil {
			return ret, err
		}
		m.Values = append(m.Values, &ast.ExpressionBinary{
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
		return nil, fmt.Errorf("expect '}',but '%s'", ep.parser.token.Desp)
	}
	ep.Next() // skip }
	return ret, nil
}

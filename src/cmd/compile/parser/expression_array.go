package parser

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

// []int{1,2,3}
func (ep *ExpressionParser) parseArrayExpression() (*ast.Expression, error) {
	pos := ep.parser.mkPos()
	ep.parser.Next() // skip [
	var t *ast.VariableType
	var err error
	if ep.parser.token.Type != lex.TOKEN_RB {
		arr := &ast.ExpressionArrayLiteral{}
		arr.Expressions, err = ep.parseExpressions()
		if ep.parser.token.Type != lex.TOKEN_RB {
			err = fmt.Errorf("%s '[' and ']' not match", ep.parser.errorMsgPrefix())
		} else {
			ep.Next() // skip ]
		}
		return &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_ARRAY,
			Data: arr,
			Pos:  pos,
		}, err
	} else {
		ep.Next()
		t, err = ep.parser.parseType()
		if err != nil {
			ep.parser.consume(untils_lc)
			ep.parser.Next() //
			return nil, err
		}
	}
	arr := &ast.ExpressionArrayLiteral{}
	if t != nil {
		arr.Typ = &ast.VariableType{}
		arr.Typ.Typ = ast.VARIABLE_TYPE_ARRAY
		arr.Typ.CombinationType = t
		arr.Typ.Pos = pos
	}
	arr.Expressions, err = ep.parseArrayValues()
	return &ast.Expression{
		Typ:  ast.EXPRESSION_TYPE_ARRAY,
		Data: arr,
		Pos:  pos,
	}, err
}

//{1,2,3}  {{1,2,3},{456}}
func (ep *ExpressionParser) parseArrayValues() ([]*ast.Expression, error) {
	if ep.parser.token.Type != lex.TOKEN_LC {
		return nil, fmt.Errorf("%s expect '{',but '%s'", ep.parser.errorMsgPrefix(), ep.parser.token.Desp)
	}
	ep.Next() // skip {
	es := []*ast.Expression{}
	for ep.parser.eof == false && ep.parser.token.Type != lex.TOKEN_RC {
		if ep.parser.token.Type == lex.TOKEN_LC {
			ees, err := ep.parseArrayValues()
			if err != nil {
				return es, err
			}
			arre := &ast.Expression{Typ: ast.EXPRESSION_TYPE_ARRAY}
			data := ast.ExpressionArrayLiteral{}
			data.Expressions = ees
			arre.Data = data
			es = append(es, arre)
		} else {
			e, err := ep.parseExpression()
			if e != nil {
				if e.Typ == ast.EXPRESSION_TYPE_LIST {
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
		return es, fmt.Errorf("%s expect '}',but '%s' ", ep.parser.errorMsgPrefix(), ep.parser.token.Desp)
	}
	ep.Next()
	return es, nil
}

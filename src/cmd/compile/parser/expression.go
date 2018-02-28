package parser

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

type ExpressionParser struct {
	parser *Parser
}

func (ep *ExpressionParser) Next() {
	ep.parser.Next()
}

func (ep *ExpressionParser) parseExpressions() ([]*ast.Expression, error) {
	es := []*ast.Expression{}
	for !ep.parser.eof {
		e, err := ep.parseExpression()
		if err != nil {
			return es, err
		}
		if e.Typ == ast.EXPRESSION_TYPE_LIST {
			es = append(es, e.Data.([]*ast.Expression)...)
		} else {
			es = append(es, e)
		}
		if ep.parser.token.Type != lex.TOKEN_COMMA {
			break
		}
		// == ,
		ep.Next() // skip ,
	}
	return es, nil
}

//parse assign expression
func (ep *ExpressionParser) parseExpression() (*ast.Expression, error) {
	return ep.parseAssignExpression()
}

func (ep *ExpressionParser) parseTypeConvertionExpression() (*ast.Expression, error) {
	t, err := ep.parser.parseType()
	if err != nil {
		return nil, err
	}
	if ep.parser.token.Type != lex.TOKEN_LP {
		return nil, fmt.Errorf("%s not ( after a type", ep.parser.errorMsgPrefix())
	}
	ep.Next() // skip (
	e, err := ep.parseExpression()
	if err != nil {
		return nil, err
	}
	if ep.parser.token.Type != lex.TOKEN_RP {
		return nil, fmt.Errorf("%s ( and ) not match", ep.parser.errorMsgPrefix())
	}
	ep.Next() // skip ) for next process
	return &ast.Expression{
		Typ: ast.EXPRESSION_TYPE_CONVERTION_TYPE,
		Data: &ast.ExpressionTypeConvertion{
			Typ:        t,
			Expression: e,
		},
	}, nil
}

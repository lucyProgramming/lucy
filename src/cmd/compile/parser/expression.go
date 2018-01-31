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

// []int{1,2,3}
func (ep *ExpressionParser) parseArrayExpression() (*ast.Expression, error) {
	fmt.Println(ep.parser.token.Desp)
	t, err := ep.parser.parseType()
	if err != nil {
		return nil, err
	}
	arr := &ast.ExpressionArray{}
	arr.Typ = t
	arr.Expression, err = ep.parseValueGroups()
	return &ast.Expression{
		Typ:  ast.EXPRESSION_TYPE_ARRAY,
		Data: arr,
	}, err
}

//{1,2,3}  {{1,2,3},{456}}
func (ep *ExpressionParser) parseValueGroups() (*ast.Expression, error) {
	if ep.parser.token.Type != lex.TOKEN_LC {
		return nil, fmt.Errorf("%s no { after type definition", ep.parser.errorMsgPrefix())
	}
	ret := &ast.Expression{}
	ret.Typ = ast.EXPRESSION_TYPE_LIST
	es := []*ast.Expression{}
	var e *ast.Expression
	var err error
	defer func() {
		ret.Data = es
	}()
	ep.Next() // skip {
	if ep.parser.token.Type == lex.TOKEN_RC {
		return ret, nil
	}
	if ep.parser.token.Type == lex.TOKEN_LC {
		for !ep.parser.eof && ep.parser.token.Type == lex.TOKEN_LC {
			e, err = ep.parseValueGroups()
			if err != nil {
				return ret, err
			}
			es = append(es, e)
			if ep.parser.token.Type != lex.TOKEN_COMMA {
				break
			}
			ep.Next()
		}
	} else {
		es, err = ep.parseExpressions()
		if err != nil {
			return ret, nil
		}
	}
	if ep.parser.token.Type != lex.TOKEN_RC {
		return ret, fmt.Errorf("%s missing } ", ep.parser.errorMsgPrefix())
	}
	ep.Next()
	return ret, nil

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

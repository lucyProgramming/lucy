package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type ExpressionParser struct {
	parser *Parser
}

func (expressionParser *ExpressionParser) Next() {
	expressionParser.parser.Next()
}

func (expressionParser *ExpressionParser) parseExpressions() ([]*ast.Expression, error) {
	es := []*ast.Expression{}
	for expressionParser.parser.token.Type != lex.TOKEN_EOF {
		e, err := expressionParser.parseExpression(false)
		if err != nil {
			return es, err
		}
		if e.Type == ast.EXPRESSION_TYPE_LIST {
			es = append(es, e.Data.([]*ast.Expression)...)
		} else {
			es = append(es, e)
		}
		if expressionParser.parser.token.Type != lex.TOKEN_COMMA {
			break
		}
		// == ,
		expressionParser.Next() // skip ,
	}
	return es, nil
}

//parse assign expression
func (expressionParser *ExpressionParser) parseExpression(statementLevel bool) (*ast.Expression, error) {
	left, err := expressionParser.parseTernaryExpression() //
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TOKEN_COMMA && statementLevel { // read more
		expressionParser.Next()                                 //  skip comma
		left2, err := expressionParser.parseTernaryExpression() //
		if err != nil {
			return nil, err
		}
		if left.Type == ast.EXPRESSION_TYPE_LIST {
			list := left.Data.([]*ast.Expression)
			left.Data = append(list, left2)
		} else {
			newExpression := &ast.Expression{}
			newExpression.Type = ast.EXPRESSION_TYPE_LIST
			list := []*ast.Expression{left, left2}
			newExpression.Data = list
			left = newExpression
		}
	}
	mustBeOneExpression := func() {
		if left.Type == ast.EXPRESSION_TYPE_LIST {
			es := left.Data.([]*ast.Expression)
			left = es[0]
			if len(es) > 1 {
				expressionParser.parser.errs = append(expressionParser.parser.errs, fmt.Errorf("%s expect one expression on left",
					expressionParser.parser.errorMsgPrefix(es[1].Pos)))
			}
		}
	}
	mkExpression := func(typ int, multi bool) (*ast.Expression, error) {
		pos := expressionParser.parser.mkPos()
		expressionParser.Next() // skip = :=
		result := &ast.Expression{}
		result.Type = typ
		bin := &ast.ExpressionBinary{}
		result.Data = bin
		bin.Left = left
		result.Pos = pos
		if multi {
			es, err := expressionParser.parseExpressions()
			if err != nil {
				return result, err
			}
			bin.Right = &ast.Expression{}
			bin.Right.Type = ast.EXPRESSION_TYPE_LIST
			bin.Right.Data = es
		} else {
			bin.Right, err = expressionParser.parseExpression(false)
		}
		return result, err
	}
	// := += -= *= /= %=
	switch expressionParser.parser.token.Type {
	case lex.TOKEN_ASSIGN:
		return mkExpression(ast.EXPRESSION_TYPE_ASSIGN, true)
	case lex.TOKEN_COLON_ASSIGN:
		return mkExpression(ast.EXPRESSION_TYPE_COLON_ASSIGN, true)
	case lex.TOKEN_ADD_ASSIGN:
		mustBeOneExpression()
		return mkExpression(ast.EXPRESSION_TYPE_PLUS_ASSIGN, false)
	case lex.TOKEN_SUB_ASSIGN:
		mustBeOneExpression()
		return mkExpression(ast.EXPRESSION_TYPE_MINUS_ASSIGN, false)
	case lex.TOKEN_MUL_ASSIGN:
		mustBeOneExpression()
		return mkExpression(ast.EXPRESSION_TYPE_MUL_ASSIGN, false)
	case lex.TOKEN_DIV_ASSIGN:
		mustBeOneExpression()
		return mkExpression(ast.EXPRESSION_TYPE_DIV_ASSIGN, false)
	case lex.TOKEN_MOD_ASSIGN:
		mustBeOneExpression()
		return mkExpression(ast.EXPRESSION_TYPE_MOD_ASSIGN, false)
	case lex.TOKEN_LSH_ASSIGN:
		mustBeOneExpression()
		return mkExpression(ast.EXPRESSION_TYPE_LSH_ASSIGN, false)
	case lex.TOKEN_RSH_ASSIGN:
		mustBeOneExpression()
		return mkExpression(ast.EXPRESSION_TYPE_RSH_ASSIGN, false)
	case lex.TOKEN_AND_ASSIGN:
		mustBeOneExpression()
		return mkExpression(ast.EXPRESSION_TYPE_AND_ASSIGN, false)
	case lex.TOKEN_OR_ASSIGN:
		mustBeOneExpression()
		return mkExpression(ast.EXPRESSION_TYPE_OR_ASSIGN, false)
	case lex.TOKEN_XOR_ASSIGN:
		mustBeOneExpression()
		return mkExpression(ast.EXPRESSION_TYPE_XOR_ASSIGN, false)

	}
	return left, nil
}

func (expressionParser *ExpressionParser) parseTypeConversionExpression() (*ast.Expression, error) {
	t, err := expressionParser.parser.parseType()
	if err != nil {
		return nil, err
	}
	if expressionParser.parser.token.Type != lex.TOKEN_LP {
		return nil, fmt.Errorf("%s not '(' after a type", expressionParser.parser.errorMsgPrefix())
	}
	pos := expressionParser.parser.mkPos()
	expressionParser.Next() // skip (
	e, err := expressionParser.parseExpression(false)
	if err != nil {
		return nil, err
	}
	if expressionParser.parser.token.Type != lex.TOKEN_RP {
		return nil, fmt.Errorf("%s '(' and ')' not match", expressionParser.parser.errorMsgPrefix())
	}
	expressionParser.Next() // skip )
	return &ast.Expression{
		Type: ast.EXPRESSION_TYPE_CHECK_CAST,
		Data: &ast.ExpressionTypeConversion{
			Type:       t,
			Expression: e,
		},
		Pos: pos,
	}, nil
}

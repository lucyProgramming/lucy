package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

// []int{1,2,3}
func (expressionParser *ExpressionParser) parseArrayExpression() (*ast.Expression, error) {
	pos := expressionParser.parser.mkPos()
	expressionParser.parser.Next(lfIsToken) // skip [
	var err error
	if expressionParser.parser.token.Type != lex.TokenRb {
		arr := &ast.ExpressionArray{}
		arr.Expressions, err = expressionParser.parseExpressions()
		if expressionParser.parser.token.Type != lex.TokenRb {
			err = fmt.Errorf("%s '[' and ']' not match", expressionParser.parser.errorMsgPrefix())
			return nil, err
		} else {
			expressionParser.Next(lfNotToken) // skip ]
		}
		return &ast.Expression{
			Type: ast.ExpressionTypeArray,
			Data: arr,
			Pos:  pos,
		}, err
	}
	expressionParser.Next(lfNotToken) // skip ]
	t, err := expressionParser.parser.parseType()
	if err != nil {
		return nil, err
	}
	if expressionParser.parser.token.Type == lex.TokenLp { // []byte("1111111111")
		expressionParser.Next(lfNotToken) // skip (
		e, err := expressionParser.parseExpression(false)
		if err != nil {
			return nil, err
		}
		if expressionParser.parser.token.Type != lex.TokenRp {
			return nil, fmt.Errorf("%s '(' and  ')' not match",
				expressionParser.parser.errorMsgPrefix())
		}
		expressionParser.Next(lfNotToken) // skip )
		ret := &ast.Expression{}
		ret.Pos = pos
		ret.Type = ast.ExpressionTypeCheckCast
		data := &ast.ExpressionTypeConversion{}
		data.Type = &ast.Type{}
		data.Type.Type = ast.VariableTypeArray
		data.Type.Pos = pos
		data.Type.Array = t
		data.Expression = e
		ret.Data = data
		return ret, nil
	}
	arr := &ast.ExpressionArray{}
	if t != nil {
		arr.Type = &ast.Type{}
		arr.Type.Type = ast.VariableTypeArray
		arr.Type.Array = t
		arr.Type.Pos = pos
	}
	arr.Expressions, err = expressionParser.parseArrayValues()
	return &ast.Expression{
		Type: ast.ExpressionTypeArray,
		Data: arr,
		Pos:  pos,
	}, err

}

//{1,2,3}  {{1,2,3},{456}}
func (expressionParser *ExpressionParser) parseArrayValues() ([]*ast.Expression, error) {
	if expressionParser.parser.token.Type != lex.TokenLc {
		return nil, fmt.Errorf("%s expect '{',but '%s'",
			expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
	}
	expressionParser.Next(lfNotToken) // skip {
	es := []*ast.Expression{}
	for expressionParser.parser.token.Type != lex.TokenEof && expressionParser.parser.token.Type != lex.TokenRc {
		if expressionParser.parser.token.Type == lex.TokenLc {
			ees, err := expressionParser.parseArrayValues()
			if err != nil {
				return es, err
			}
			arrayExpression := &ast.Expression{Type: ast.ExpressionTypeArray}
			data := ast.ExpressionArray{}
			data.Expressions = ees
			arrayExpression.Data = data
			es = append(es, arrayExpression)
		} else {
			e, err := expressionParser.parseExpression(false)
			if e != nil {
				es = append(es, e)
			}
			if err != nil {
				return es, err
			}
		}
		if expressionParser.parser.token.Type == lex.TokenComma {
			expressionParser.Next(lfNotToken) // skip ,
		} else {
			break
		}
	}
	expressionParser.parser.ifTokenIsLfSkip()
	if expressionParser.parser.token.Type != lex.TokenRc {
		return es, fmt.Errorf("%s expect '}',but '%s'", expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
	}
	expressionParser.Next(lfNotToken)
	return es, nil
}

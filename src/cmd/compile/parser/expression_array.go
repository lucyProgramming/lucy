package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

// []int{1,2,3}
func (expressionParser *ExpressionParser) parseArrayExpression() (*ast.Expression, error) {
	expressionParser.parser.Next(lfIsToken) // skip [
	expressionParser.parser.unExpectNewLineAndSkip()
	var err error
	if expressionParser.parser.token.Type != lex.TokenRb {
		/*
			[1 ,2]
		*/
		arr := &ast.ExpressionArray{}
		arr.Expressions, err = expressionParser.parseExpressions(lex.TokenRb)
		expressionParser.parser.ifTokenIsLfThenSkip()
		if expressionParser.parser.token.Type != lex.TokenRb {
			err = fmt.Errorf("%s '[' and ']' not match", expressionParser.parser.errorMsgPrefix())
			return nil, err
		}
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken) // skip ]
		return &ast.Expression{
			Type:        ast.ExpressionTypeArray,
			Data:        arr,
			Pos:         pos,
			Description: "arrayLiteral",
		}, err
	}
	expressionParser.Next(lfIsToken) // skip ]
	expressionParser.parser.unExpectNewLineAndSkip()
	array, err := expressionParser.parser.parseType()
	if err != nil {
		return nil, err
	}
	if expressionParser.parser.token.Type == lex.TokenLp {
		/*
			[]byte("1111111111")
		*/
		expressionParser.Next(lfNotToken) // skip (
		e, err := expressionParser.parseExpression(false)
		if err != nil {
			return nil, err
		}
		if expressionParser.parser.token.Type != lex.TokenRp {
			return nil, fmt.Errorf("%s '(' and  ')' not match",
				expressionParser.parser.errorMsgPrefix())
		}
		ret := &ast.Expression{}
		ret.Description = "checkCast"
		pos := expressionParser.parser.mkPos()
		ret.Pos = pos
		ret.Type = ast.ExpressionTypeCheckCast
		data := &ast.ExpressionTypeConversion{}
		data.Type = &ast.Type{}
		data.Type.Type = ast.VariableTypeArray
		data.Type.Pos = pos
		data.Type.Array = array
		data.Expression = e
		ret.Data = data
		expressionParser.Next(lfIsToken) // skip )
		return ret, nil
	}
	expressionParser.parser.unExpectNewLineAndSkip()
	arr := &ast.ExpressionArray{}
	if array != nil {
		arr.Type = &ast.Type{}
		arr.Type.Type = ast.VariableTypeArray
		arr.Type.Array = array
		arr.Type.Pos = array.Pos
	}
	/*
		[]int { 1, 2}
	*/
	arr.Expressions, err = expressionParser.parseArrayValues()
	return &ast.Expression{
		Type:        ast.ExpressionTypeArray,
		Data:        arr,
		Pos:         expressionParser.parser.mkPos(),
		Description: "arrayLiteral",
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
	for expressionParser.parser.token.Type != lex.TokenEof &&
		expressionParser.parser.token.Type != lex.TokenRc {
		if expressionParser.parser.token.Type == lex.TokenComment ||
			expressionParser.parser.token.Type == lex.TokenCommentMultiLine {
			expressionParser.Next(lfIsToken)
			continue
		}
		if expressionParser.parser.token.Type == lex.TokenLc {
			ees, err := expressionParser.parseArrayValues()
			if err != nil {
				return es, err
			}
			arrayExpression := &ast.Expression{
				Type: ast.ExpressionTypeArray,
				Pos:  expressionParser.parser.mkPos(),
			}
			arrayExpression.Description = "arrayLiteral"
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
	expressionParser.parser.ifTokenIsLfThenSkip()
	if expressionParser.parser.token.Type != lex.TokenRc {
		return es, fmt.Errorf("%s expect '}',but '%s'",
			expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
	}
	expressionParser.Next(lfIsToken)
	return es, nil
}

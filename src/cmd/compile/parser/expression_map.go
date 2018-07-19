package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (expressionParser *ExpressionParser) parseMapExpression() (*ast.Expression, error) {
	var typ *ast.Type
	var err error
	pos := expressionParser.parser.mkPos()
	if expressionParser.parser.token.Type == lex.TokenMap {
		typ, err = expressionParser.parser.parseType()
		if err != nil {
			return nil, err
		}
	}
	if expressionParser.parser.token.Type != lex.TokenLc {
		return nil, fmt.Errorf("expect '{',but '%s'", expressionParser.parser.token.Description)
	}
	expressionParser.Next(lfNotToken) // skip {
	ret := &ast.Expression{
		Type: ast.ExpressionTypeMap,
		Pos:  pos,
	}
	m := &ast.ExpressionMap{}
	m.Type = typ
	ret.Data = m
	for expressionParser.parser.token.Type != lex.TokenEof && expressionParser.parser.token.Type != lex.TokenRc {
		// key
		k, err := expressionParser.parseExpression(false)
		if err != nil {
			return ret, err
		}
		// arrow
		if expressionParser.parser.token.Type != lex.TokenArrow {
			return ret, fmt.Errorf("%s expect '->',but '%s'",
				expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
		}
		expressionParser.Next(lfNotToken) // skip ->
		// value
		v, err := expressionParser.parseExpression(false)
		if err != nil {
			return ret, err
		}
		m.KeyValuePairs = append(m.KeyValuePairs, &ast.ExpressionBinary{
			Left:  k,
			Right: v,
		})
		if expressionParser.parser.token.Type == lex.TokenComma {
			// read next  key value pair
			expressionParser.Next(lfNotToken)
		} else {
			break
		}
	}
	expressionParser.parser.ifTokenIsLfSkip()
	if expressionParser.parser.token.Type != lex.TokenRc {
		return nil, fmt.Errorf("%s expect '}',but '%s'",
			expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
	}
	expressionParser.Next(lfNotToken) // skip }
	return ret, nil
}

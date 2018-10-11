package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (expressionParser *ExpressionParser) parseMapExpression() (*ast.Expression, error) {
	var typ *ast.Type
	var err error
	if expressionParser.parser.token.Type == lex.TokenMap {
		typ, err = expressionParser.parser.parseType()
		if err != nil {
			return nil, err
		}
		expressionParser.parser.ifTokenIsLfThenSkip()
	}
	if expressionParser.parser.token.Type != lex.TokenLc {
		err := fmt.Errorf("expect '{',but '%s'", expressionParser.parser.token.Description)
		expressionParser.parser.errs = append(expressionParser.parser.errs, err)
		return nil, err
	}
	expressionParser.Next(lfNotToken) // skip {
	ret := &ast.Expression{
		Type:        ast.ExpressionTypeMap,
		Description: "mapLiteral",
	}
	m := &ast.ExpressionMap{}
	m.Type = typ
	ret.Data = m
	for expressionParser.parser.token.Type != lex.TokenEof &&
		expressionParser.parser.token.Type != lex.TokenRc {
		// key
		k, err := expressionParser.parseExpression(false)
		if err != nil {
			return ret, err
		}
		expressionParser.parser.unExpectNewLineAndSkip()
		// arrow
		if expressionParser.parser.token.Type != lex.TokenArrow {
			err := fmt.Errorf("%s expect '->',but '%s'",
				expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
			expressionParser.parser.errs = append(expressionParser.parser.errs, err)
			return ret, err
		}
		expressionParser.Next(lfNotToken) // skip ->
		// value
		v, err := expressionParser.parseExpression(false)
		if err != nil {
			return ret, err
		}
		m.KeyValuePairs = append(m.KeyValuePairs, &ast.ExpressionKV{
			Key:   k,
			Value: v,
		})
		if expressionParser.parser.token.Type == lex.TokenComma {
			// read next  key value pair
			expressionParser.Next(lfNotToken)
		} else {
			break
		}
	}
	expressionParser.parser.ifTokenIsLfThenSkip()
	if expressionParser.parser.token.Type != lex.TokenRc {
		err := fmt.Errorf("%s expect '}',but '%s'",
			expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
		expressionParser.parser.errs = append(expressionParser.parser.errs, err)
		expressionParser.parser.consume(untilRc)
	}
	ret.Pos = expressionParser.parser.mkPos()
	expressionParser.Next(lfIsToken) // skip }
	return ret, nil
}

package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type ExpressionParser struct {
	parser *Parser
}

// wrapper
func (expressionParser *ExpressionParser) Next(lfIsToken bool) {
	expressionParser.parser.Next(lfIsToken)
}

func (expressionParser *ExpressionParser) parseExpressions(endTokens ...lex.TokenKind) ([]*ast.Expression, error) {
	es := []*ast.Expression{}
	for expressionParser.parser.token.Type != lex.TokenEof {
		if expressionParser.parser.token.Type == lex.TokenComment ||
			expressionParser.parser.token.Type == lex.TokenCommentMultiLine {
			expressionParser.Next(lfIsToken)
			continue
		}
		e, err := expressionParser.parseExpression(false)
		if err != nil {
			return es, err
		}
		es = append(es, e)
		if expressionParser.parser.token.Type != lex.TokenComma {
			if expressionParser.looksLikeExpression() {
				/*
					missing comma
					a(1 2)
				*/
				expressionParser.parser.errs = append(expressionParser.parser.errs, fmt.Errorf("%s missing comma",
					expressionParser.parser.errorMsgPrefix()))
				continue
			}
			break
		}
		// == ,
		expressionParser.Next(lfNotToken) // skip ,
		for expressionParser.parser.token.Type == lex.TokenComma {
			expressionParser.parser.errs = append(expressionParser.parser.errs,
				fmt.Errorf("%s missing expression", expressionParser.parser.errorMsgPrefix()))
			expressionParser.Next(lfNotToken) // skip ,
		}
		for _, v := range endTokens {
			if v == expressionParser.parser.token.Type {
				// found end token
				expressionParser.parser.errs = append(expressionParser.parser.errs,
					fmt.Errorf("%s extra comma", expressionParser.parser.errorMsgPrefix()))
				goto end
			}
		}
	}
end:
	return es, nil
}

/*
	parse assign expression
*/
func (expressionParser *ExpressionParser) parseExpression(statementLevel bool) (*ast.Expression, error) {
	left, err := expressionParser.parseQuestionExpression() //
	if err != nil {
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenComma && statementLevel { // read more
		expressionParser.Next(lfNotToken)                        //  skip comma
		left2, err := expressionParser.parseQuestionExpression() //
		if err != nil {
			return nil, err
		}
		if left.Type == ast.ExpressionTypeList {
			left.Data = append(left.Data.([]*ast.Expression), left2)
		} else {
			newExpression := &ast.Expression{}
			newExpression.Type = ast.ExpressionTypeList
			newExpression.Pos = left.Pos
			list := []*ast.Expression{left, left2}
			newExpression.Data = list
			left = newExpression
		}
	}
	parseRight := func(expressionType ast.ExpressionTypeKind, isMulti bool) (*ast.Expression, error) {
		pos := expressionParser.parser.mkPos()
		name := expressionParser.parser.token.Description
		expressionParser.Next(lfNotToken) // skip = :=
		result := &ast.Expression{}
		result.Type = expressionType
		result.Description = name
		bin := &ast.ExpressionBinary{}
		result.Data = bin
		bin.Left = left
		result.Pos = pos
		if isMulti {
			es, err := expressionParser.parseExpressions(lex.TokenSemicolon)
			if err != nil {
				return nil, err
			}
			bin.Right = &ast.Expression{}
			bin.Right.Type = ast.ExpressionTypeList
			bin.Right.Data = es
		} else {
			bin.Right, err = expressionParser.parseExpression(false)
			if err != nil {
				return nil, err
			}
		}
		return result, err
	}

	switch expressionParser.parser.token.Type {
	case lex.TokenAssign:
		return parseRight(ast.ExpressionTypeAssign, true)
	case lex.TokenVarAssign:
		return parseRight(ast.ExpressionTypeVarAssign, true)
	case lex.TokenAddAssign:
		return parseRight(ast.ExpressionTypePlusAssign, false)
	case lex.TokenSubAssign:
		return parseRight(ast.ExpressionTypeMinusAssign, false)
	case lex.TokenMulAssign:
		return parseRight(ast.ExpressionTypeMulAssign, false)
	case lex.TokenDivAssign:
		return parseRight(ast.ExpressionTypeDivAssign, false)
	case lex.TokenModAssign:
		return parseRight(ast.ExpressionTypeModAssign, false)
	case lex.TokenLshAssign:
		return parseRight(ast.ExpressionTypeLshAssign, false)
	case lex.TokenRshAssign:
		return parseRight(ast.ExpressionTypeRshAssign, false)
	case lex.TokenAndAssign:
		return parseRight(ast.ExpressionTypeAndAssign, false)
	case lex.TokenOrAssign:
		return parseRight(ast.ExpressionTypeOrAssign, false)
	case lex.TokenXorAssign:
		return parseRight(ast.ExpressionTypeXorAssign, false)
	}
	return left, nil
}

func (expressionParser *ExpressionParser) parseTypeConversionExpression() (*ast.Expression, error) {
	to, err := expressionParser.parser.parseType()
	if err != nil {
		return nil, err
	}
	expressionParser.parser.unExpectNewLineAndSkip()
	if expressionParser.parser.token.Type != lex.TokenLp {
		err := fmt.Errorf("%s not '(' after a type",
			expressionParser.parser.errorMsgPrefix())
		expressionParser.parser.errs = append(expressionParser.parser.errs, err)
		return nil, err
	}
	expressionParser.Next(lfNotToken) // skip (
	e, err := expressionParser.parseExpression(false)
	if err != nil {
		return nil, err
	}
	expressionParser.parser.ifTokenIsLfThenSkip()
	if expressionParser.parser.token.Type != lex.TokenRp {
		err := fmt.Errorf("%s '(' and ')' not match", expressionParser.parser.errorMsgPrefix())
		expressionParser.parser.errs = append(expressionParser.parser.errs, err)
		return nil, err
	}
	pos := expressionParser.parser.mkPos()
	expressionParser.Next(lfIsToken) // skip )
	return &ast.Expression{
		Type: ast.ExpressionTypeCheckCast,
		Data: &ast.ExpressionTypeConversion{
			Type:       to,
			Expression: e,
		},
		Pos: pos,
	}, nil
}

func (expressionParser *ExpressionParser) looksLikeExpression() bool {
	return expressionParser.parser.token.Type == lex.TokenIdentifier ||
		expressionParser.parser.token.Type == lex.TokenTrue ||
		expressionParser.parser.token.Type == lex.TokenFalse ||
		expressionParser.parser.token.Type == lex.TokenGlobal ||
		expressionParser.parser.token.Type == lex.TokenLiteralByte ||
		expressionParser.parser.token.Type == lex.TokenLiteralShort ||
		expressionParser.parser.token.Type == lex.TokenLiteralInt ||
		expressionParser.parser.token.Type == lex.TokenLiteralLong ||
		expressionParser.parser.token.Type == lex.TokenLiteralFloat ||
		expressionParser.parser.token.Type == lex.TokenLiteralDouble ||
		expressionParser.parser.token.Type == lex.TokenLiteralString ||
		expressionParser.parser.token.Type == lex.TokenNull ||
		expressionParser.parser.token.Type == lex.TokenLp ||
		expressionParser.parser.token.Type == lex.TokenIncrement ||
		expressionParser.parser.token.Type == lex.TokenDecrement ||
		expressionParser.parser.token.Type == lex.TokenNot ||
		expressionParser.parser.token.Type == lex.TokenBitNot ||
		expressionParser.parser.token.Type == lex.TokenSub ||
		expressionParser.parser.token.Type == lex.TokenFn ||
		expressionParser.parser.token.Type == lex.TokenNew ||
		expressionParser.parser.token.Type == lex.TokenLb ||
		expressionParser.parser.token.Type == lex.TokenSelection
}

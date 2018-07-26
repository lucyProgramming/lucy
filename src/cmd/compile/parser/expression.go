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
		newExpression := &ast.Expression{}
		newExpression.Type = ast.ExpressionTypeList
		newExpression.Pos = left.Pos
		list := []*ast.Expression{left, left2}
		newExpression.Data = list
		left = newExpression
	}
	mustBeOneExpression := func(left *ast.Expression) {
		if left.Type == ast.ExpressionTypeList {
			es := left.Data.([]*ast.Expression)
			left = es[0]
			if len(es) > 1 {
				expressionParser.parser.errs = append(expressionParser.parser.errs,
					fmt.Errorf("%s expect one expression on left",
						expressionParser.parser.errorMsgPrefix(es[1].Pos)))
			}
		}
	}
	mkExpression := func(expressionType ast.ExpressionTypeKind, isMulti bool) (*ast.Expression, error) {
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfNotToken) // skip = :=
		result := &ast.Expression{}
		result.Type = expressionType
		bin := &ast.ExpressionBinary{}
		result.Data = bin
		bin.Left = left
		result.Pos = pos
		if isMulti {
			es, err := expressionParser.parseExpressions(lex.TokenSemicolon)
			if err != nil {
				return result, err
			}
			bin.Right = &ast.Expression{}
			bin.Right.Type = ast.ExpressionTypeList
			bin.Right.Data = es
		} else {
			bin.Right, err = expressionParser.parseExpression(false)
		}
		return result, err
	}
	// := += -= *= /= %=
	switch expressionParser.parser.token.Type {
	case lex.TokenAssign:
		return mkExpression(ast.ExpressionTypeAssign, true)
	case lex.TokenColonAssign:
		return mkExpression(ast.ExpressionTypeColonAssign, true)
	case lex.TokenAddAssign:
		mustBeOneExpression(left)
		return mkExpression(ast.ExpressionTypePlusAssign, false)
	case lex.TokenSubAssign:
		mustBeOneExpression(left)
		return mkExpression(ast.ExpressionTypeMinusAssign, false)
	case lex.TokenMulAssign:
		mustBeOneExpression(left)
		return mkExpression(ast.ExpressionTypeMulAssign, false)
	case lex.TokenDivAssign:
		mustBeOneExpression(left)
		return mkExpression(ast.ExpressionTypeDivAssign, false)
	case lex.TokenModAssign:
		mustBeOneExpression(left)
		return mkExpression(ast.ExpressionTypeModAssign, false)
	case lex.TokenLshAssign:
		mustBeOneExpression(left)
		return mkExpression(ast.ExpressionTypeLshAssign, false)
	case lex.TokenRshAssign:
		mustBeOneExpression(left)
		return mkExpression(ast.ExpressionTypeRshAssign, false)
	case lex.TokenAndAssign:
		mustBeOneExpression(left)
		return mkExpression(ast.ExpressionTypeAndAssign, false)
	case lex.TokenOrAssign:
		mustBeOneExpression(left)
		return mkExpression(ast.ExpressionTypeOrAssign, false)
	case lex.TokenXorAssign:
		mustBeOneExpression(left)
		return mkExpression(ast.ExpressionTypeXorAssign, false)
	}
	return left, nil
}

func (expressionParser *ExpressionParser) parseTypeConversionExpression() (*ast.Expression, error) {
	pos := expressionParser.parser.mkPos()
	to, err := expressionParser.parser.parseType()
	if err != nil {
		return nil, err
	}
	expressionParser.parser.unExpectNewLineAndSkip()
	if expressionParser.parser.token.Type != lex.TokenLp {
		return nil, fmt.Errorf("%s not '(' after a type",
			expressionParser.parser.errorMsgPrefix())
	}
	expressionParser.Next(lfNotToken) // skip (
	e, err := expressionParser.parseExpression(false)
	if err != nil {
		return nil, err
	}
	expressionParser.parser.ifTokenIsLfThenSkip()
	if expressionParser.parser.token.Type != lex.TokenRp {
		return nil, fmt.Errorf("%s '(' and ')' not match", expressionParser.parser.errorMsgPrefix())
	}
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
		expressionParser.parser.token.Type == lex.TokenLb
}

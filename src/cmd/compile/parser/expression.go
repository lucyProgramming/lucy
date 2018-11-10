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
func (ep *ExpressionParser) Next(lfIsToken bool) {
	ep.parser.Next(lfIsToken)
}

func (ep *ExpressionParser) parseExpressions(endTokens ...lex.TokenKind) ([]*ast.Expression, error) {
	es := []*ast.Expression{}
	for ep.parser.token.Type != lex.TokenEof {
		if ep.parser.token.Type == lex.TokenComment ||
			ep.parser.token.Type == lex.TokenMultiLineComment {
			ep.Next(lfIsToken)
			continue
		}
		e, err := ep.parseExpression(false)
		if err != nil {
			return es, err
		}
		es = append(es, e)
		if ep.parser.token.Type != lex.TokenComma {
			if ep.looksLikeExpression() {
				/*
					missing comma
					a(1 2)
				*/
				ep.parser.errs = append(ep.parser.errs, fmt.Errorf("%s missing comma",
					ep.parser.errMsgPrefix()))
				continue
			}
			break
		}
		// == ,
		commnaPos := ep.parser.mkPos()
		ep.Next(lfNotToken) // skip ,
		for ep.parser.token.Type == lex.TokenComma {
			ep.parser.errs = append(ep.parser.errs,
				fmt.Errorf("%s missing expression", ep.parser.errMsgPrefix()))
			ep.Next(lfNotToken) // skip ,
		}
		for _, v := range endTokens {
			if v == ep.parser.token.Type {
				// found end token
				ep.parser.errs = append(ep.parser.errs,
					fmt.Errorf("%s extra comma", ep.parser.errMsgPrefix(commnaPos)))
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
func (ep *ExpressionParser) parseExpression(statementLevel bool) (*ast.Expression, error) {
	left, err := ep.parseQuestionExpression() //
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TokenComma && statementLevel { // read more
		ep.Next(lfNotToken)                        //  skip comma
		left2, err := ep.parseQuestionExpression() //
		if err != nil {
			return nil, err
		}
		if left.Type == ast.ExpressionTypeList {
			left.Data = append(left.Data.([]*ast.Expression), left2)
		} else {
			newExpression := &ast.Expression{}
			newExpression.Type = ast.ExpressionTypeList
			newExpression.Pos = left.Pos
			newExpression.Op = "list"
			list := []*ast.Expression{left, left2}
			newExpression.Data = list
			left = newExpression
		}
	}
	parseRight := func(expressionType ast.ExpressionTypeKind, isMulti bool) (*ast.Expression, error) {
		pos := ep.parser.mkPos()
		opName := ep.parser.token.Description
		ep.Next(lfNotToken) // skip = :=
		result := &ast.Expression{}
		result.Type = expressionType
		result.Op = opName
		bin := &ast.ExpressionBinary{}
		result.Data = bin
		bin.Left = left
		result.Pos = pos
		if isMulti {
			es, err := ep.parseExpressions(lex.TokenSemicolon)
			if err != nil {
				return nil, err
			}
			bin.Right = &ast.Expression{}
			bin.Right.Type = ast.ExpressionTypeList
			bin.Right.Data = es
		} else {
			bin.Right, err = ep.parseExpression(false)
			if err != nil {
				return nil, err
			}
		}
		return result, err
	}

	switch ep.parser.token.Type {
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

func (ep *ExpressionParser) parseTypeConversionExpression() (*ast.Expression, error) {
	to, err := ep.parser.parseType()
	if err != nil {
		return nil, err
	}
	pos := ep.parser.mkPos()
	ep.parser.unExpectNewLineAndSkip()
	if ep.parser.token.Type != lex.TokenLp {
		err := fmt.Errorf("%s not '(' after a type",
			ep.parser.errMsgPrefix())
		ep.parser.errs = append(ep.parser.errs, err)
		return nil, err
	}
	ep.Next(lfNotToken) // skip (
	e, err := ep.parseExpression(false)
	if err != nil {
		return nil, err
	}
	ep.parser.ifTokenIsLfThenSkip()
	if ep.parser.token.Type != lex.TokenRp {
		err := fmt.Errorf("%s '(' and ')' not match", ep.parser.errMsgPrefix())
		ep.parser.errs = append(ep.parser.errs, err)
		return nil, err
	}
	ep.Next(lfIsToken) // skip )
	return &ast.Expression{
		Type: ast.ExpressionTypeCheckCast,
		Data: &ast.ExpressionTypeConversion{
			Type:       to,
			Expression: e,
		},
		Pos: pos,
	}, nil
}

func (ep *ExpressionParser) looksLikeExpression() bool {
	return ep.parser.token.Type == lex.TokenIdentifier ||
		ep.parser.token.Type == lex.TokenTrue ||
		ep.parser.token.Type == lex.TokenFalse ||
		ep.parser.token.Type == lex.TokenGlobal ||
		ep.parser.token.Type == lex.TokenLiteralByte ||
		ep.parser.token.Type == lex.TokenLiteralShort ||
		ep.parser.token.Type == lex.TokenLiteralInt ||
		ep.parser.token.Type == lex.TokenLiteralLong ||
		ep.parser.token.Type == lex.TokenLiteralFloat ||
		ep.parser.token.Type == lex.TokenLiteralDouble ||
		ep.parser.token.Type == lex.TokenLiteralString ||
		ep.parser.token.Type == lex.TokenNull ||
		ep.parser.token.Type == lex.TokenLp ||
		ep.parser.token.Type == lex.TokenIncrement ||
		ep.parser.token.Type == lex.TokenDecrement ||
		ep.parser.token.Type == lex.TokenNot ||
		ep.parser.token.Type == lex.TokenBitNot ||
		ep.parser.token.Type == lex.TokenSub ||
		ep.parser.token.Type == lex.TokenFn ||
		ep.parser.token.Type == lex.TokenNew ||
		ep.parser.token.Type == lex.TokenLb ||
		ep.parser.token.Type == lex.TokenSelection
}

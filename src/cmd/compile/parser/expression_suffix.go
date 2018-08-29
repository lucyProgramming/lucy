package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (expressionParser *ExpressionParser) parseSuffixExpression() (*ast.Expression, error) {
	var suffix *ast.Expression
	var err error
	switch expressionParser.parser.token.Type {
	case lex.TokenIdentifier:
		suffix = &ast.Expression{}
		suffix.Description = "identifer " + expressionParser.parser.token.Data.(string)
		suffix.Type = ast.ExpressionTypeIdentifier
		identifier := &ast.ExpressionIdentifier{}
		identifier.Name = expressionParser.parser.token.Data.(string)
		suffix.Data = identifier
		suffix.Pos = expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
	case lex.TokenTrue:
		suffix = &ast.Expression{}
		suffix.Description = "true"
		suffix.Type = ast.ExpressionTypeBool
		suffix.Data = true
		suffix.Pos = expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
	case lex.TokenFalse:
		suffix = &ast.Expression{}
		suffix.Description = "false"
		suffix.Type = ast.ExpressionTypeBool
		suffix.Data = false
		suffix.Pos = expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
	case lex.TokenGlobal:
		suffix = &ast.Expression{}
		suffix.Description = "global"
		suffix.Type = ast.ExpressionTypeGlobal
		suffix.Pos = expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralByte:
		suffix = &ast.Expression{
			Type:        ast.ExpressionTypeByte,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "byteLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralShort:
		suffix = &ast.Expression{
			Type:        ast.ExpressionTypeShort,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "shortLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralInt:
		suffix = &ast.Expression{
			Type:        ast.ExpressionTypeInt,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "intLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralLong:
		suffix = &ast.Expression{
			Type:        ast.ExpressionTypeLong,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "longLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralFloat:
		suffix = &ast.Expression{
			Type:        ast.ExpressionTypeFloat,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "floatLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralDouble:
		suffix = &ast.Expression{
			Type:        ast.ExpressionTypeDouble,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "doubleLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralString:
		suffix = &ast.Expression{
			Type:        ast.ExpressionTypeString,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "stringLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenNull:
		suffix = &ast.Expression{
			Type:        ast.ExpressionTypeNull,
			Pos:         expressionParser.parser.mkPos(),
			Description: "null",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLp:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfNotToken)
		suffix, err = expressionParser.parseExpression(false)
		if err != nil {
			return nil, err
		}
		expressionParser.parser.ifTokenIsLfThenSkip()
		if expressionParser.parser.token.Type != lex.TokenRp {
			return nil, fmt.Errorf("%s '(' and ')' not matched, but '%s'",
				expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
		}
		newExpression := &ast.Expression{
			Type:        ast.ExpressionTypeParenthesis,
			Pos:         pos,
			Data:        suffix,
			Description: "(" + suffix.Description + ")",
		}
		suffix = newExpression
		expressionParser.Next(lfIsToken)
	case lex.TokenIncrement:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken) // skip ++
		suffix, err = expressionParser.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE := &ast.Expression{}
		newE.Pos = pos
		newE.Description = "++"
		newE.Type = ast.ExpressionTypePrefixIncrement
		newE.Data = suffix
		suffix = newE
	case lex.TokenDecrement:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken) // skip --
		suffix, err = expressionParser.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE := &ast.Expression{}
		newE.Description = "--"
		newE.Type = ast.ExpressionTypePrefixDecrement
		newE.Data = suffix
		newE.Pos = pos
		suffix = newE
	case lex.TokenNot:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
		newE := &ast.Expression{}
		newE.Description = "!"
		suffix, err = expressionParser.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE.Type = ast.ExpressionTypeNot
		newE.Data = suffix
		newE.Pos = pos
		suffix = newE
	case lex.TokenBitNot:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
		suffix, err = expressionParser.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE := &ast.Expression{}
		newE.Description = "~"
		newE.Type = ast.ExpressionTypeBitwiseNot
		newE.Data = suffix
		newE.Pos = pos
		suffix = newE
	case lex.TokenSub:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
		suffix, err = expressionParser.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE := &ast.Expression{}
		newE.Description = "-"
		newE.Type = ast.ExpressionTypeNegative
		newE.Data = suffix
		newE.Pos = pos
		suffix = newE
	case lex.TokenFn:
		pos := expressionParser.parser.mkPos()
		f, err := expressionParser.parser.FunctionParser.parse(false, false)
		if err != nil {
			return nil, err
		}
		suffix = &ast.Expression{
			Type:        ast.ExpressionTypeFunctionLiteral,
			Data:        f,
			Pos:         pos,
			Description: "functionLiteral",
		}
	case lex.TokenNew:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
		expressionParser.parser.unExpectNewLineAndSkip()
		t, err := expressionParser.parser.parseType()
		if err != nil {
			return nil, err
		}
		expressionParser.parser.unExpectNewLineAndSkip()
		if expressionParser.parser.token.Type != lex.TokenLp {
			return nil, fmt.Errorf("%s missing '(' after new", expressionParser.parser.errorMsgPrefix())
		}
		expressionParser.Next(lfNotToken) // skip (
		var es []*ast.Expression
		if expressionParser.parser.token.Type != lex.TokenRp { //
			es, err = expressionParser.parseExpressions(lex.TokenRp)
			if err != nil {
				return nil, err
			}
		}
		expressionParser.parser.ifTokenIsLfThenSkip()
		if expressionParser.parser.token.Type != lex.TokenRp {
			return nil, fmt.Errorf("%s '(' and ')' not match", expressionParser.parser.errorMsgPrefix())
		}
		expressionParser.Next(lfIsToken)
		suffix = &ast.Expression{
			Pos:  pos,
			Type: ast.ExpressionTypeNew,
			Data: &ast.ExpressionNew{
				Args: es,
				Type: t,
			},
			Description: "new",
		}
	case lex.TokenLb:
		suffix, err = expressionParser.parseArrayExpression()
		if err != nil {
			return suffix, err
		}
	// bool(xxx)
	case lex.TokenBool:
		suffix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return suffix, err
		}
		//byte()
	case lex.TokenByte:
		suffix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return suffix, err
		}
		//short()
	case lex.TokenShort:
		suffix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return suffix, err
		}
		//int()
	case lex.TokenInt:
		suffix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return suffix, err
		}
		//long()
	case lex.TokenLong:
		suffix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return suffix, err
		}
		//float()
	case lex.TokenFloat:
		suffix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return suffix, err
		}
		//double
	case lex.TokenDouble:
		suffix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return suffix, err
		}
		//string()
	case lex.TokenString:
		suffix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return suffix, err
		}
		// T()
	case lex.TokenTemplate:
		suffix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return suffix, err
		}
		// range
	case lex.TokenRange:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
		expressionParser.parser.unExpectNewLineAndSkip()
		e, err := expressionParser.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		suffix = &ast.Expression{}
		suffix.Description = "range"
		suffix.Type = ast.ExpressionTypeRange
		suffix.Pos = pos
		suffix.Data = e
		return suffix, nil
	case lex.TokenMap:
		suffix, err = expressionParser.parseMapExpression()
		if err != nil {
			return suffix, err
		}
	case lex.TokenLc:
		suffix, err = expressionParser.parseMapExpression()
		if err != nil {
			return suffix, err
		}
	case lex.TokenLf:
		expressionParser.parser.unExpectNewLineAndSkip()
		return expressionParser.parseSuffixExpression()
	default:
		err = fmt.Errorf("%s unkown begining of a expression, token:'%s'",
			expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenIncrement ||
		expressionParser.parser.token.Type == lex.TokenDecrement ||
		expressionParser.parser.token.Type == lex.TokenLp ||
		expressionParser.parser.token.Type == lex.TokenLb ||
		expressionParser.parser.token.Type == lex.TokenSelection ||
		expressionParser.parser.token.Type == lex.TokenVArgs {
		// ++ or --
		if expressionParser.parser.token.Type == lex.TokenVArgs {
			newExpression := &ast.Expression{}
			newExpression.Description = "..."
			newExpression.Type = ast.ExpressionTypeVArgs
			newExpression.Data = suffix
			newExpression.Pos = expressionParser.parser.mkPos()
			expressionParser.Next(lfIsToken)
			return newExpression, nil
		}
		if expressionParser.parser.token.Type == lex.TokenIncrement ||
			expressionParser.parser.token.Type == lex.TokenDecrement {
			newExpression := &ast.Expression{}
			if expressionParser.parser.token.Type == lex.TokenIncrement {
				newExpression.Type = ast.ExpressionTypeIncrement
				newExpression.Description = "++"
			} else {
				newExpression.Type = ast.ExpressionTypeDecrement
				newExpression.Description = "--"
			}
			newExpression.Data = suffix
			suffix = newExpression
			newExpression.Pos = expressionParser.parser.mkPos()
			expressionParser.Next(lfIsToken)
			continue
		}
		// [
		if expressionParser.parser.token.Type == lex.TokenLb {
			pos := expressionParser.parser.mkPos()
			expressionParser.Next(lfNotToken) // skip [
			if expressionParser.parser.token.Type == lex.TokenColon {
				/*
					a[:]
				*/
				expressionParser.Next(lfNotToken) // skip :
				var end *ast.Expression
				if expressionParser.parser.token.Type != lex.TokenRb {
					end, err = expressionParser.parseExpression(false)
					if err != nil {
						return nil, err
					}
				}
				expressionParser.parser.ifTokenIsLfThenSkip()
				if expressionParser.parser.token.Type != lex.TokenRb {
					return nil, fmt.Errorf("%s '[' and ']' not match", expressionParser.parser.errorMsgPrefix())
				}
				expressionParser.Next(lfIsToken) // skip ]
				newExpression := &ast.Expression{}
				newExpression.Type = ast.ExpressionTypeSlice
				newExpression.Description = "slice"
				newExpression.Pos = expressionParser.parser.mkPos()
				slice := &ast.ExpressionSlice{}
				newExpression.Data = slice
				slice.ExpressionOn = suffix
				slice.End = end
				suffix = newExpression
				continue
			}
			e, err := expressionParser.parseExpression(false)
			if err != nil {
				return nil, err
			}
			if expressionParser.parser.token.Type == lex.TokenColon {
				expressionParser.Next(lfNotToken)
				var end *ast.Expression
				if expressionParser.parser.token.Type != lex.TokenRb {
					end, err = expressionParser.parseExpression(false)
					if err != nil {
						return nil, err
					}
				}
				if expressionParser.parser.token.Type != lex.TokenRb {
					return nil, fmt.Errorf("%s '[' and ']' not match", expressionParser.parser.errorMsgPrefix())
				}
				expressionParser.Next(lfIsToken) // skip ]
				newExpression := &ast.Expression{}
				newExpression.Type = ast.ExpressionTypeSlice
				newExpression.Description = "slice"
				newExpression.Pos = expressionParser.parser.mkPos()
				slice := &ast.ExpressionSlice{}
				newExpression.Data = slice
				slice.Start = e
				slice.ExpressionOn = suffix
				slice.End = end
				suffix = newExpression
				continue
			}
			if expressionParser.parser.token.Type != lex.TokenRb {
				return nil, fmt.Errorf("%s '[' and ']' not match", expressionParser.parser.errorMsgPrefix())
			}
			newExpression := &ast.Expression{}
			newExpression.Pos = pos
			newExpression.Description = "index"
			newExpression.Type = ast.ExpressionTypeIndex
			index := &ast.ExpressionIndex{}
			index.Expression = suffix
			index.Index = e
			newExpression.Data = index
			suffix = newExpression
			expressionParser.Next(lfIsToken)
			continue
		}
		// aaa.xxxx
		if expressionParser.parser.token.Type == lex.TokenSelection {
			pos := expressionParser.parser.mkPos()
			expressionParser.Next(lfNotToken) // skip .
			if expressionParser.parser.token.Type == lex.TokenIdentifier {
				newExpression := &ast.Expression{}
				newExpression.Pos = pos
				newExpression.Description = "selection"
				newExpression.Type = ast.ExpressionTypeSelection
				selection := &ast.ExpressionSelection{}
				selection.Expression = suffix
				selection.Name = expressionParser.parser.token.Data.(string)
				newExpression.Data = selection
				suffix = newExpression
				expressionParser.Next(lfIsToken)
			} else if expressionParser.parser.token.Type == lex.TokenLp { //  a.(xxx)
				//
				expressionParser.Next(lfNotToken) // skip (
				typ, err := expressionParser.parser.parseType()
				if err != nil {
					return nil, err
				}
				expressionParser.parser.ifTokenIsLfThenSkip()
				if expressionParser.parser.token.Type != lex.TokenRp {
					return nil, fmt.Errorf("%s '(' and ')' not match", expressionParser.parser.errorMsgPrefix())
				}
				expressionParser.Next(lfIsToken) // skip  )
				newExpression := &ast.Expression{}
				newExpression.Pos = pos
				newExpression.Description = "assert"
				newExpression.Type = ast.ExpressionTypeTypeAssert
				typeAssert := &ast.ExpressionTypeAssert{}
				typeAssert.Type = typ
				typeAssert.Expression = suffix
				newExpression.Data = typeAssert
				suffix = newExpression
			} else {
				return nil, fmt.Errorf("%s expect  'identifier' or '(',but '%s'",
					expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
			}
			continue
		}
		// aa()
		if expressionParser.parser.token.Type == lex.TokenLp {
			newExpression, err := expressionParser.parseCallExpression(suffix)
			if err != nil {
				return nil, err
			}
			suffix = newExpression
			continue
		}
	}
	return suffix, nil
}

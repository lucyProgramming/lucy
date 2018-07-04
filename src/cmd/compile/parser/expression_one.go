package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (expressionParser *ExpressionParser) parseOneExpression(isPrefixUnary bool) (*ast.Expression, error) {
	var left *ast.Expression
	var err error
	switch expressionParser.parser.token.Type {
	case lex.TokenIdentifier:
		left = &ast.Expression{}
		left.Type = ast.ExpressionTypeIdentifier
		identifier := &ast.ExpressionIdentifier{}
		identifier.Name = expressionParser.parser.token.Data.(string)
		left.Data = identifier
		left.Pos = expressionParser.parser.mkPos()
		expressionParser.Next()
	case lex.TokenTrue:
		left = &ast.Expression{}
		left.Type = ast.ExpressionTypeBool
		left.Data = true
		left.Pos = expressionParser.parser.mkPos()
		expressionParser.Next()
	case lex.TokenFalse:
		left = &ast.Expression{}
		left.Type = ast.ExpressionTypeBool
		left.Data = false
		left.Pos = expressionParser.parser.mkPos()
		expressionParser.Next()
	case lex.TokenGlobal:
		left = &ast.Expression{}
		left.Type = ast.ExpressionTypeGlobal
		left.Pos = expressionParser.parser.mkPos()
		expressionParser.Next()
	case lex.TokenLiteralByte:
		left = &ast.Expression{
			Type: ast.ExpressionTypeByte,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TokenLiteralShort:
		left = &ast.Expression{
			Type: ast.ExpressionTypeShort,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TokenLiteralInt:
		left = &ast.Expression{
			Type: ast.ExpressionTypeInt,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TokenLiteralLong:
		left = &ast.Expression{
			Type: ast.ExpressionTypeLong,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TokenLiteralFloat:
		left = &ast.Expression{
			Type: ast.ExpressionTypeFloat,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TokenLiteralDouble:
		left = &ast.Expression{
			Type: ast.ExpressionTypeDouble,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TokenLiteralString:
		left = &ast.Expression{
			Type: ast.ExpressionTypeString,
			Data: expressionParser.parser.token.Data,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TokenNull:
		left = &ast.Expression{
			Type: ast.ExpressionTypeNull,
			Pos:  expressionParser.parser.mkPos(),
		}
		expressionParser.Next()
	case lex.TokenLp:
		expressionParser.Next()
		left, err = expressionParser.parseExpression(false)
		if err != nil {
			return nil, err
		}
		if expressionParser.parser.token.Type != lex.TokenRp {
			return nil, fmt.Errorf("%s '(' and ')' not matched, but '%s'",
				expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
		}
		expressionParser.Next()
	case lex.TokenIncrement:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next() // skip ++
		newE := &ast.Expression{}
		newE.Pos = pos
		left, err = expressionParser.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.ExpressionTypePrefixIncrement
		newE.Data = left
		left = newE
	case lex.TokenDecrement:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next() // skip --
		newE := &ast.Expression{}
		left, err = expressionParser.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.ExpressionTypePrefixDecrement
		newE.Data = left
		newE.Pos = pos
		left = newE
	case lex.TokenNot:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		newE := &ast.Expression{}
		left, err = expressionParser.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.ExpressionTypeNot
		newE.Data = left
		newE.Pos = pos
		left = newE
	case lex.TokenBitNot:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		newE := &ast.Expression{}
		left, err = expressionParser.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.ExpressionTypeBitwiseNot
		newE.Data = left
		newE.Pos = pos
		left = newE
	case lex.TokenSub:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		newE := &ast.Expression{}
		left, err = expressionParser.parseOneExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Type = ast.ExpressionTypeNegative
		newE.Data = left
		newE.Pos = pos
		left = newE
	case lex.TokenFunction:
		pos := expressionParser.parser.mkPos()
		f, err := expressionParser.parser.FunctionParser.parse(false)
		if err != nil {
			return nil, err
		}
		left = &ast.Expression{
			Type: ast.ExpressionTypeFunctionLiteral,
			Data: f,
			Pos:  pos,
		}
	case lex.TokenNew:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		t, err := expressionParser.parser.parseType()
		if err != nil {
			return nil, err
		}
		if expressionParser.parser.token.Type != lex.TokenLp {
			return nil, fmt.Errorf("%s missing '(' after new", expressionParser.parser.errorMsgPrefix())
		}
		expressionParser.Next() // skip (
		var es []*ast.Expression
		if expressionParser.parser.token.Type != lex.TokenRp { //
			es, err = expressionParser.parseExpressions()
			if err != nil {
				return nil, err
			}
		}
		if expressionParser.parser.token.Type != lex.TokenRp {
			return nil, fmt.Errorf("%s ( and ) not match", expressionParser.parser.errorMsgPrefix())
		}
		expressionParser.Next()
		left = &ast.Expression{
			Pos:  pos,
			Type: ast.ExpressionTypeNew,
			Data: &ast.ExpressionNew{
				Args: es,
				Type: t,
			},
		}
	case lex.TokenLb:
		left, err = expressionParser.parseArrayExpression()
		if err != nil {
			return left, err
		}
	// bool(xxx)
	case lex.TokenBool:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
		//byte()
	case lex.TokenByte:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
		//short()
	case lex.TokenShort:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
		//int()
	case lex.TokenInt:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
		//long()
	case lex.TokenLong:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
		//float()
	case lex.TokenFloat:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
		//double
	case lex.TokenDouble:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
		//string()
	case lex.TokenString:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
		// T()
	case lex.TokenTemplate:
		left, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return left, err
		}
		// range
	case lex.TokenRange:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next()
		e, err := expressionParser.parseOneExpression(false)
		if err != nil {
			return nil, err
		}
		left = &ast.Expression{}
		left.Type = ast.ExpressionTypeRange
		left.Pos = pos
		left.Data = e
		return left, nil
	case lex.TokenMap:
		left, err = expressionParser.parseMapExpression()
		if err != nil {
			return left, err
		}
	case lex.TokenLc:
		left, err = expressionParser.parseMapExpression()
		if err != nil {
			return left, err
		}
	default:
		err = fmt.Errorf("%s unkown begining of a expression, token:'%s'",
			expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenIncrement ||
		expressionParser.parser.token.Type == lex.TokenDecrement ||
		expressionParser.parser.token.Type == lex.TokenLp ||
		expressionParser.parser.token.Type == lex.TokenLb ||
		expressionParser.parser.token.Type == lex.TokenSelection {
		// ++ or --
		if expressionParser.parser.token.Type == lex.TokenIncrement ||
			expressionParser.parser.token.Type == lex.TokenDecrement { //  ++ or --
			if isPrefixUnary {
				return left, nil
			}
			newExpression := &ast.Expression{}
			if expressionParser.parser.token.Type == lex.TokenIncrement {
				newExpression.Type = ast.ExpressionTypeIncrement
			} else {
				newExpression.Type = ast.ExpressionTypeDecrement
			}
			newExpression.Data = left
			left = newExpression
			newExpression.Pos = expressionParser.parser.mkPos()
			expressionParser.Next()
			continue
		}
		// [
		if expressionParser.parser.token.Type == lex.TokenLb {
			pos := expressionParser.parser.mkPos()
			expressionParser.Next() // skip [
			if expressionParser.parser.token.Type == lex.TokenColon {
				/*
					a[:]
				*/
				expressionParser.Next() // skip :
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
				expressionParser.Next() // skip ]
				newExpression := &ast.Expression{}
				newExpression.Type = ast.ExpressionTypeSlice
				newExpression.Pos = expressionParser.parser.mkPos()
				slice := &ast.ExpressionSlice{}
				newExpression.Data = slice
				slice.Expression = left
				slice.End = end
				left = newExpression
				continue
			}
			e, err := expressionParser.parseExpression(false)
			if err != nil {
				return nil, err
			}
			if expressionParser.parser.token.Type == lex.TokenColon {
				expressionParser.parser.Next()
				if expressionParser.parser.token.Type == lex.TokenColon {
					expressionParser.parser.Next() // skip :
				}
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
				expressionParser.Next() // skip ]
				newExpression := &ast.Expression{}
				newExpression.Type = ast.ExpressionTypeSlice
				newExpression.Pos = expressionParser.parser.mkPos()
				slice := &ast.ExpressionSlice{}
				newExpression.Data = slice
				slice.Start = e
				slice.Expression = left
				slice.End = end
				left = newExpression
				continue
			}
			if expressionParser.parser.token.Type != lex.TokenRb {
				return nil, fmt.Errorf("%s '[' and ']' not match", expressionParser.parser.errorMsgPrefix())
			}
			newExpression := &ast.Expression{}
			newExpression.Pos = pos
			newExpression.Type = ast.ExpressionTypeIndex
			index := &ast.ExpressionIndex{}
			index.Expression = left
			index.Index = e
			newExpression.Data = index
			left = newExpression
			expressionParser.Next()
			continue
		}
		// aaa.xxxx
		if expressionParser.parser.token.Type == lex.TokenSelection {
			pos := expressionParser.parser.mkPos()
			expressionParser.parser.Next() // skip .
			if expressionParser.parser.token.Type == lex.TokenIdentifier {
				newExpression := &ast.Expression{}
				newExpression.Pos = pos
				newExpression.Type = ast.ExpressionTypeSelection
				index := &ast.ExpressionSelection{}
				index.Expression = left
				index.Name = expressionParser.parser.token.Data.(string)
				newExpression.Data = index
				left = newExpression
				expressionParser.Next()
			} else if expressionParser.parser.token.Type == lex.TokenLp { //  a.(xxx)
				//
				expressionParser.Next() // skip (
				typ, err := expressionParser.parser.parseType()
				if err != nil {
					return nil, err
				}
				if expressionParser.parser.token.Type != lex.TokenRp {
					return nil, fmt.Errorf("%s '(' and ')' not match", expressionParser.parser.errorMsgPrefix())
				}
				expressionParser.Next() // skip  )
				newExpression := &ast.Expression{}
				newExpression.Pos = pos
				newExpression.Type = ast.ExpressionTypeTypeAssert
				typeAssert := &ast.ExpressionTypeAssert{}
				typeAssert.Type = typ
				typeAssert.Expression = left
				newExpression.Data = typeAssert
				left = newExpression
			} else {
				return nil, fmt.Errorf("%s expect  'identifier' or '(',but '%s'",
					expressionParser.parser.errorMsgPrefix(), expressionParser.parser.token.Description)
			}
			continue
		}
		// aa()
		if expressionParser.parser.token.Type == lex.TokenLp {
			newExpression, err := expressionParser.parseCallExpression(left)
			if err != nil {
				return nil, err
			}
			left = newExpression
			continue
		}
	}
	return left, nil
}

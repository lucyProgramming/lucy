package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (expressionParser *ExpressionParser) parseSuffixExpression() (*ast.Expression, error) {
	var prefix *ast.Expression
	var err error
	switch expressionParser.parser.token.Type {
	case lex.TokenIdentifier:
		prefix = &ast.Expression{}
		prefix.Description = expressionParser.parser.token.Data.(string)
		prefix.Type = ast.ExpressionTypeIdentifier
		identifier := &ast.ExpressionIdentifier{}
		identifier.Name = expressionParser.parser.token.Data.(string)
		prefix.Data = identifier
		prefix.Pos = expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
	case lex.TokenTrue:
		prefix = &ast.Expression{}
		prefix.Description = "true"
		prefix.Type = ast.ExpressionTypeBool
		prefix.Data = true
		prefix.Pos = expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
	case lex.TokenFalse:
		prefix = &ast.Expression{}
		prefix.Description = "false"
		prefix.Type = ast.ExpressionTypeBool
		prefix.Data = false
		prefix.Pos = expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
	case lex.TokenSelection:
		prefix = &ast.Expression{}
		prefix.Description = "."
		prefix.Type = ast.ExpressionTypeDot
		prefix.Data = true
		prefix.Pos = expressionParser.parser.mkPos()
		//special case , no next
	case lex.TokenGlobal:
		prefix = &ast.Expression{}
		prefix.Description = "global"
		prefix.Type = ast.ExpressionTypeGlobal
		prefix.Pos = expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralByte:
		prefix = &ast.Expression{
			Type:        ast.ExpressionTypeByte,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "byteLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralShort:
		prefix = &ast.Expression{
			Type:        ast.ExpressionTypeShort,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "shortLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralChar:
		prefix = &ast.Expression{
			Type:        ast.ExpressionTypeChar,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "shortLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralInt:
		prefix = &ast.Expression{
			Type:        ast.ExpressionTypeInt,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "intLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralLong:
		prefix = &ast.Expression{
			Type:        ast.ExpressionTypeLong,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "longLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralFloat:
		prefix = &ast.Expression{
			Type:        ast.ExpressionTypeFloat,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "floatLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralDouble:
		prefix = &ast.Expression{
			Type:        ast.ExpressionTypeDouble,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "doubleLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLiteralString:
		prefix = &ast.Expression{
			Type:        ast.ExpressionTypeString,
			Data:        expressionParser.parser.token.Data,
			Pos:         expressionParser.parser.mkPos(),
			Description: "stringLiteral",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenNull:
		prefix = &ast.Expression{
			Type:        ast.ExpressionTypeNull,
			Pos:         expressionParser.parser.mkPos(),
			Description: "null",
		}
		expressionParser.Next(lfIsToken)
	case lex.TokenLp:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfNotToken)
		prefix, err = expressionParser.parseExpression(false)
		if err != nil {
			return nil, err
		}
		expressionParser.parser.ifTokenIsLfThenSkip()
		if expressionParser.parser.token.Type != lex.TokenRp {
			err := fmt.Errorf("%s '(' and ')' not matched, but '%s'",
				expressionParser.parser.errMsgPrefix(), expressionParser.parser.token.Description)
			expressionParser.parser.errs = append(expressionParser.parser.errs, err)
			return nil, err
		}
		newExpression := &ast.Expression{
			Type:        ast.ExpressionTypeParenthesis,
			Pos:         pos,
			Data:        prefix,
			Description: "(" + prefix.Description + ")",
		}
		prefix = newExpression
		expressionParser.Next(lfIsToken)
	case lex.TokenIncrement:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken) // skip ++
		prefix, err = expressionParser.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE := &ast.Expression{}
		newE.Pos = pos
		newE.Description = expressionParser.parser.token.Description
		newE.Type = ast.ExpressionTypePrefixIncrement
		newE.Data = prefix
		prefix = newE
	case lex.TokenDecrement:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken) // skip --
		prefix, err = expressionParser.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE := &ast.Expression{}
		newE.Description = expressionParser.parser.token.Description
		newE.Type = ast.ExpressionTypePrefixDecrement
		newE.Data = prefix
		newE.Pos = pos
		prefix = newE
	case lex.TokenNot:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
		newE := &ast.Expression{}
		newE.Description = expressionParser.parser.token.Description
		prefix, err = expressionParser.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE.Type = ast.ExpressionTypeNot
		newE.Data = prefix
		newE.Pos = pos
		prefix = newE
	case lex.TokenBitNot:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
		prefix, err = expressionParser.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE := &ast.Expression{}
		newE.Description = expressionParser.parser.token.Description
		newE.Type = ast.ExpressionTypeBitwiseNot
		newE.Data = prefix
		newE.Pos = pos
		prefix = newE
	case lex.TokenSub:
		pos := expressionParser.parser.mkPos()
		expressionParser.Next(lfIsToken)
		prefix, err = expressionParser.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE := &ast.Expression{}
		newE.Description = expressionParser.parser.token.Description
		newE.Type = ast.ExpressionTypeNegative
		newE.Data = prefix
		newE.Pos = pos
		prefix = newE
	case lex.TokenFn:
		pos := expressionParser.parser.mkPos()
		f, err := expressionParser.parser.FunctionParser.parse(false, false)
		if err != nil {
			return nil, err
		}
		prefix = &ast.Expression{
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
			err := fmt.Errorf("%s missing '(' after new", expressionParser.parser.errMsgPrefix())
			expressionParser.parser.errs = append(expressionParser.parser.errs, err)
			return nil, err
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
			err := fmt.Errorf("%s '(' and ')' not match", expressionParser.parser.errMsgPrefix())
			expressionParser.parser.errs = append(expressionParser.parser.errs, err)
			return nil, err
		}
		expressionParser.Next(lfIsToken)
		prefix = &ast.Expression{
			Pos:  pos,
			Type: ast.ExpressionTypeNew,
			Data: &ast.ExpressionNew{
				Args: es,
				Type: t,
			},
			Description: "new",
		}
	case lex.TokenLb:
		prefix, err = expressionParser.parseArrayExpression()
		if err != nil {
			return prefix, err
		}
	// bool(xxx)
	case lex.TokenBool:
		prefix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//byte()
	case lex.TokenByte:
		prefix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//short()
	case lex.TokenShort:
		prefix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//char()
	case lex.TokenChar:
		prefix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//int()
	case lex.TokenInt:
		prefix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//long()
	case lex.TokenLong:
		prefix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//float()
	case lex.TokenFloat:
		prefix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//double
	case lex.TokenDouble:
		prefix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//string()
	case lex.TokenString:
		prefix, err = expressionParser.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
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
		prefix = &ast.Expression{}
		prefix.Description = "range"
		prefix.Type = ast.ExpressionTypeRange
		prefix.Pos = pos
		prefix.Data = e
		return prefix, nil
	case lex.TokenMap:
		prefix, err = expressionParser.parseMapExpression()
		if err != nil {
			return prefix, err
		}
	case lex.TokenLc:
		prefix, err = expressionParser.parseMapExpression()
		if err != nil {
			return prefix, err
		}

	case lex.TokenLf:
		expressionParser.parser.unExpectNewLineAndSkip()
		return expressionParser.parseSuffixExpression()
	default:
		err = fmt.Errorf("%s unkown begining of a expression, token:'%s'",
			expressionParser.parser.errMsgPrefix(), expressionParser.parser.token.Description)
		expressionParser.parser.errs = append(expressionParser.parser.errs, err)
		return nil, err
	}
	for expressionParser.parser.token.Type == lex.TokenIncrement ||
		expressionParser.parser.token.Type == lex.TokenDecrement ||
		expressionParser.parser.token.Type == lex.TokenLp ||
		expressionParser.parser.token.Type == lex.TokenLb ||
		expressionParser.parser.token.Type == lex.TokenSelection ||
		expressionParser.parser.token.Type == lex.TokenVArgs ||
		expressionParser.parser.token.Type == lex.Token2Colon {
		switch expressionParser.parser.token.Type {
		case lex.TokenVArgs:
			newExpression := &ast.Expression{}
			newExpression.Description = "..."
			newExpression.Type = ast.ExpressionTypeVArgs
			newExpression.Data = prefix
			newExpression.Pos = expressionParser.parser.mkPos()
			expressionParser.Next(lfIsToken)
			return newExpression, nil
		case lex.TokenIncrement,
			lex.TokenDecrement:
			newExpression := &ast.Expression{}
			newExpression.Description = expressionParser.parser.token.Description
			if expressionParser.parser.token.Type == lex.TokenIncrement {
				newExpression.Type = ast.ExpressionTypeIncrement
			} else {
				newExpression.Type = ast.ExpressionTypeDecrement
			}
			newExpression.Data = prefix
			prefix = newExpression
			newExpression.Pos = expressionParser.parser.mkPos()
			expressionParser.Next(lfIsToken)
		case lex.TokenLb:
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
					err := fmt.Errorf("%s '[' and ']' not match", expressionParser.parser.errMsgPrefix())
					expressionParser.parser.errs = append(expressionParser.parser.errs, err)
					return nil, err
				}
				newExpression := &ast.Expression{}
				newExpression.Type = ast.ExpressionTypeSlice
				newExpression.Description = "slice"
				newExpression.Pos = expressionParser.parser.mkPos()
				slice := &ast.ExpressionSlice{}
				newExpression.Data = slice
				slice.ExpressionOn = prefix
				slice.End = end
				prefix = newExpression
				expressionParser.Next(lfIsToken) // skip ]
			} else {
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
						err := fmt.Errorf("%s '[' and ']' not match", expressionParser.parser.errMsgPrefix())
						expressionParser.parser.errs = append(expressionParser.parser.errs, err)
						return nil, err
					}
					newExpression := &ast.Expression{}
					newExpression.Type = ast.ExpressionTypeSlice
					newExpression.Description = "slice"
					newExpression.Pos = expressionParser.parser.mkPos()
					slice := &ast.ExpressionSlice{}
					newExpression.Data = slice
					slice.Start = e
					slice.ExpressionOn = prefix
					slice.End = end
					prefix = newExpression
					expressionParser.Next(lfIsToken) // skip ]
				} else {
					if expressionParser.parser.token.Type != lex.TokenRb {
						err := fmt.Errorf("%s '[' and ']' not match", expressionParser.parser.errMsgPrefix())
						expressionParser.parser.errs = append(expressionParser.parser.errs, err)
						return nil, err
					}
					newExpression := &ast.Expression{}
					newExpression.Pos = pos
					newExpression.Description = "index"
					newExpression.Type = ast.ExpressionTypeIndex
					index := &ast.ExpressionIndex{}
					index.Expression = prefix
					index.Index = e
					newExpression.Data = index
					prefix = newExpression
					expressionParser.Next(lfIsToken)
				}
			}
		case lex.Token2Colon:
			pos := expressionParser.parser.mkPos()
			expressionParser.Next(lfNotToken) // skip ::
			var constName string
			if expressionParser.parser.token.Type != lex.TokenIdentifier {
				expressionParser.parser.errs = append(expressionParser.parser.errs,
					fmt.Errorf("%s expect idnetifier , but '%s'",
						expressionParser.parser.errMsgPrefix(), expressionParser.parser.token.Description))
				constName = compileAutoName()
			} else {
				constName = expressionParser.parser.token.Data.(string)
				expressionParser.Next(lfIsToken)
			}
			newExpression := &ast.Expression{}
			newExpression.Pos = pos
			newExpression.Description = "selectConst"
			newExpression.Type = ast.ExpressionTypeSelectionConst
			selection := &ast.ExpressionSelection{}
			selection.Expression = prefix
			selection.Name = constName
			newExpression.Data = selection
			prefix = newExpression
		case lex.TokenSelection:
			expressionParser.Next(lfNotToken) // skip .
			if expressionParser.parser.token.Type == lex.TokenIdentifier {
				newExpression := &ast.Expression{}
				newExpression.Pos = expressionParser.parser.mkPos()
				newExpression.Description = "selection"
				newExpression.Type = ast.ExpressionTypeSelection
				selection := &ast.ExpressionSelection{}
				selection.Expression = prefix
				selection.Name = expressionParser.parser.token.Data.(string)
				newExpression.Data = selection
				prefix = newExpression
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
					err := fmt.Errorf("%s '(' and ')' not match", expressionParser.parser.errMsgPrefix())
					expressionParser.parser.errs = append(expressionParser.parser.errs, err)
					return nil, err
				}
				newExpression := &ast.Expression{}
				newExpression.Pos = expressionParser.parser.mkPos()
				newExpression.Description = "assert"
				newExpression.Type = ast.ExpressionTypeTypeAssert
				typeAssert := &ast.ExpressionTypeAssert{}
				typeAssert.Type = typ
				typeAssert.Expression = prefix
				newExpression.Data = typeAssert
				prefix = newExpression
				expressionParser.Next(lfIsToken) // skip  )
			} else {
				err := fmt.Errorf("%s expect  'identifier' or '(',but '%s'",
					expressionParser.parser.errMsgPrefix(), expressionParser.parser.token.Description)
				expressionParser.parser.errs = append(expressionParser.parser.errs, err)
				return nil, err
			}
		case lex.TokenLp:
			newExpression, err := expressionParser.parseCallExpression(prefix)
			if err != nil {
				return nil, err
			}
			prefix = newExpression
		}
	}

	return prefix, nil
}

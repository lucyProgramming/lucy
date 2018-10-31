package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (ep *ExpressionParser) parseSuffixExpression() (*ast.Expression, error) {
	var prefix *ast.Expression
	var err error
	switch ep.parser.token.Type {
	case lex.TokenIdentifier:
		prefix = &ast.Expression{}
		prefix.Op = ep.parser.token.Data.(string)
		prefix.Type = ast.ExpressionTypeIdentifier
		identifier := &ast.ExpressionIdentifier{}
		identifier.Name = ep.parser.token.Data.(string)
		prefix.Data = identifier
		prefix.Pos = ep.parser.mkPos()
		ep.Next(lfIsToken)
	case lex.TokenTrue:
		prefix = &ast.Expression{}
		prefix.Op = "true"
		prefix.Type = ast.ExpressionTypeBool
		prefix.Data = true
		prefix.Pos = ep.parser.mkPos()
		ep.Next(lfIsToken)
	case lex.TokenFalse:
		prefix = &ast.Expression{}
		prefix.Op = "false"
		prefix.Type = ast.ExpressionTypeBool
		prefix.Data = false
		prefix.Pos = ep.parser.mkPos()
		ep.Next(lfIsToken)
	case lex.TokenSelection:
		prefix = &ast.Expression{}
		prefix.Op = "."
		prefix.Type = ast.ExpressionTypeDot
		prefix.Data = true
		prefix.Pos = ep.parser.mkPos()
		//special case , no next
	case lex.TokenGlobal:
		prefix = &ast.Expression{}
		prefix.Op = "global"
		prefix.Type = ast.ExpressionTypeGlobal
		prefix.Pos = ep.parser.mkPos()
		ep.Next(lfIsToken)
	case lex.TokenLiteralByte:
		prefix = &ast.Expression{
			Type: ast.ExpressionTypeByte,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
			Op:   "byteLiteral",
		}
		ep.Next(lfIsToken)
	case lex.TokenLiteralShort:
		prefix = &ast.Expression{
			Type: ast.ExpressionTypeShort,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
			Op:   "shortLiteral",
		}
		ep.Next(lfIsToken)
	case lex.TokenLiteralChar:
		prefix = &ast.Expression{
			Type: ast.ExpressionTypeChar,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
			Op:   "charLiteral",
		}
		ep.Next(lfIsToken)
	case lex.TokenLiteralInt:
		prefix = &ast.Expression{
			Type: ast.ExpressionTypeInt,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
			Op:   "intLiteral",
		}
		ep.Next(lfIsToken)
	case lex.TokenLiteralLong:
		prefix = &ast.Expression{
			Type: ast.ExpressionTypeLong,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
			Op:   "longLiteral",
		}
		ep.Next(lfIsToken)
	case lex.TokenLiteralFloat:
		prefix = &ast.Expression{
			Type: ast.ExpressionTypeFloat,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
			Op:   "floatLiteral",
		}
		ep.Next(lfIsToken)
	case lex.TokenLiteralDouble:
		prefix = &ast.Expression{
			Type: ast.ExpressionTypeDouble,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
			Op:   "doubleLiteral",
		}
		ep.Next(lfIsToken)
	case lex.TokenLiteralString:
		prefix = &ast.Expression{
			Type: ast.ExpressionTypeString,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
			Op:   "stringLiteral",
		}
		ep.Next(lfIsToken)
	case lex.TokenNull:
		prefix = &ast.Expression{
			Type: ast.ExpressionTypeNull,
			Pos:  ep.parser.mkPos(),
			Op:   "null",
		}
		ep.Next(lfIsToken)
	case lex.TokenLp:
		pos := ep.parser.mkPos()
		ep.Next(lfNotToken)
		prefix, err = ep.parseExpression(false)
		if err != nil {
			return nil, err
		}
		ep.parser.ifTokenIsLfThenSkip()
		if ep.parser.token.Type != lex.TokenRp {
			err := fmt.Errorf("%s '(' and ')' not matched, but '%s'",
				ep.parser.errMsgPrefix(), ep.parser.token.Description)
			ep.parser.errs = append(ep.parser.errs, err)
			return nil, err
		}
		newExpression := &ast.Expression{
			Type: ast.ExpressionTypeParenthesis,
			Pos:  pos,
			Data: prefix,
			Op:   "(" + prefix.Op + ")",
		}
		prefix = newExpression
		ep.Next(lfIsToken)
	case lex.TokenIncrement:
		pos := ep.parser.mkPos()
		ep.Next(lfIsToken) // skip ++
		prefix, err = ep.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE := &ast.Expression{}
		newE.Pos = pos
		newE.Op = "++()"
		newE.Type = ast.ExpressionTypePrefixIncrement
		newE.Data = prefix
		prefix = newE
	case lex.TokenDecrement:
		pos := ep.parser.mkPos()
		ep.Next(lfIsToken) // skip --
		prefix, err = ep.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE := &ast.Expression{}
		newE.Op = "--()"
		newE.Type = ast.ExpressionTypePrefixDecrement
		newE.Data = prefix
		newE.Pos = pos
		prefix = newE
	case lex.TokenNot:
		op := ep.parser.token.Description
		pos := ep.parser.mkPos()
		ep.Next(lfIsToken)
		newE := &ast.Expression{}
		newE.Op = op
		prefix, err = ep.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE.Type = ast.ExpressionTypeNot
		newE.Data = prefix
		newE.Pos = pos
		prefix = newE
	case lex.TokenBitNot:
		op := ep.parser.token.Description
		pos := ep.parser.mkPos()
		ep.Next(lfIsToken)
		prefix, err = ep.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE := &ast.Expression{}
		newE.Op = op
		newE.Type = ast.ExpressionTypeBitwiseNot
		newE.Data = prefix
		newE.Pos = pos
		prefix = newE
	case lex.TokenSub:
		op := ep.parser.token.Description
		pos := ep.parser.mkPos()
		ep.Next(lfIsToken)
		prefix, err = ep.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		newE := &ast.Expression{}
		newE.Op = op
		newE.Type = ast.ExpressionTypeNegative
		newE.Data = prefix
		newE.Pos = pos
		prefix.NegativeExpression = newE
		prefix = newE
	case lex.TokenFn:
		pos := ep.parser.mkPos()
		f, err := ep.parser.FunctionParser.parse(false, false)
		if err != nil {
			return nil, err
		}
		prefix = &ast.Expression{
			Type: ast.ExpressionTypeFunctionLiteral,
			Data: f,
			Pos:  pos,
			Op:   "functionLiteral",
		}
	case lex.TokenNew:
		pos := ep.parser.mkPos()
		ep.Next(lfIsToken)
		ep.parser.unExpectNewLineAndSkip()
		t, err := ep.parser.parseType()
		if err != nil {
			return nil, err
		}
		ep.parser.unExpectNewLineAndSkip()
		if ep.parser.token.Type != lex.TokenLp {
			err := fmt.Errorf("%s missing '(' after new", ep.parser.errMsgPrefix())
			ep.parser.errs = append(ep.parser.errs, err)
			return nil, err
		}
		ep.Next(lfNotToken) // skip (
		var es []*ast.Expression
		if ep.parser.token.Type != lex.TokenRp { //
			es, err = ep.parseExpressions(lex.TokenRp)
			if err != nil {
				return nil, err
			}
		}
		ep.parser.ifTokenIsLfThenSkip()
		if ep.parser.token.Type != lex.TokenRp {
			err := fmt.Errorf("%s '(' and ')' not match", ep.parser.errMsgPrefix())
			ep.parser.errs = append(ep.parser.errs, err)
			return nil, err
		}
		ep.Next(lfIsToken)
		prefix = &ast.Expression{
			Pos:  pos,
			Type: ast.ExpressionTypeNew,
			Data: &ast.ExpressionNew{
				Args: es,
				Type: t,
			},
			Op: "new",
		}
	case lex.TokenLb:
		prefix, err = ep.parseArrayExpression()
		if err != nil {
			return prefix, err
		}
	// bool(xxx)
	case lex.TokenBool:
		prefix, err = ep.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//byte()
	case lex.TokenByte:
		prefix, err = ep.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//short()
	case lex.TokenShort:
		prefix, err = ep.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//char()
	case lex.TokenChar:
		prefix, err = ep.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//int()
	case lex.TokenInt:
		prefix, err = ep.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//long()
	case lex.TokenLong:
		prefix, err = ep.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//float()
	case lex.TokenFloat:
		prefix, err = ep.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//double
	case lex.TokenDouble:
		prefix, err = ep.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		//string()
	case lex.TokenString:
		prefix, err = ep.parseTypeConversionExpression()
		if err != nil {
			return prefix, err
		}
		// range
	case lex.TokenRange:
		pos := ep.parser.mkPos()
		ep.Next(lfIsToken)
		ep.parser.unExpectNewLineAndSkip()
		e, err := ep.parseSuffixExpression()
		if err != nil {
			return nil, err
		}
		prefix = &ast.Expression{}
		prefix.Op = "range"
		prefix.Type = ast.ExpressionTypeRange
		prefix.Pos = pos
		prefix.Data = e
		return prefix, nil
	case lex.TokenMap:
		prefix, err = ep.parseMapExpression()
		if err != nil {
			return prefix, err
		}
	case lex.TokenLc:
		prefix, err = ep.parseMapExpression()
		if err != nil {
			return prefix, err
		}

	case lex.TokenLf:
		ep.parser.unExpectNewLineAndSkip()
		return ep.parseSuffixExpression()
	default:
		err = fmt.Errorf("%s unkown begining of a expression, token:'%s'",
			ep.parser.errMsgPrefix(), ep.parser.token.Description)
		ep.parser.errs = append(ep.parser.errs, err)
		return nil, err
	}
	for ep.parser.token.Type == lex.TokenIncrement ||
		ep.parser.token.Type == lex.TokenDecrement ||
		ep.parser.token.Type == lex.TokenLp ||
		ep.parser.token.Type == lex.TokenLb ||
		ep.parser.token.Type == lex.TokenSelection ||
		ep.parser.token.Type == lex.TokenVArgs ||
		ep.parser.token.Type == lex.TokenSelectConst {
		switch ep.parser.token.Type {
		case lex.TokenVArgs:
			newExpression := &ast.Expression{}
			newExpression.Op = "..."
			newExpression.Type = ast.ExpressionTypeVArgs
			newExpression.Data = prefix
			newExpression.Pos = ep.parser.mkPos()
			ep.Next(lfIsToken)
			return newExpression, nil
		case lex.TokenIncrement,
			lex.TokenDecrement:
			newExpression := &ast.Expression{}
			if ep.parser.token.Type == lex.TokenIncrement {
				newExpression.Op = "()++"
				newExpression.Type = ast.ExpressionTypeIncrement
			} else {
				newExpression.Op = "()--"
				newExpression.Type = ast.ExpressionTypeDecrement
			}
			newExpression.Data = prefix
			prefix = newExpression
			newExpression.Pos = ep.parser.mkPos()
			ep.Next(lfIsToken)
		case lex.TokenLb:
			ep.Next(lfNotToken) // skip [
			if ep.parser.token.Type == lex.TokenColon {
				/*
					a[:]
				*/
				ep.Next(lfNotToken) // skip :
				var end *ast.Expression
				if ep.parser.token.Type != lex.TokenRb {
					end, err = ep.parseExpression(false)
					if err != nil {
						return nil, err
					}
				}
				ep.parser.ifTokenIsLfThenSkip()
				if ep.parser.token.Type != lex.TokenRb {
					err := fmt.Errorf("%s '[' and ']' not match", ep.parser.errMsgPrefix())
					ep.parser.errs = append(ep.parser.errs, err)
					return nil, err
				}
				newExpression := &ast.Expression{}
				newExpression.Type = ast.ExpressionTypeSlice
				newExpression.Op = "slice"
				newExpression.Pos = ep.parser.mkPos()
				slice := &ast.ExpressionSlice{}
				newExpression.Data = slice
				slice.ExpressionOn = prefix
				slice.End = end
				prefix = newExpression
				ep.Next(lfIsToken) // skip ]
			} else {
				e, err := ep.parseExpression(false)
				if err != nil {
					return nil, err
				}
				if ep.parser.token.Type == lex.TokenColon {
					ep.Next(lfNotToken)
					var end *ast.Expression
					if ep.parser.token.Type != lex.TokenRb {
						end, err = ep.parseExpression(false)
						if err != nil {
							return nil, err
						}
					}
					if ep.parser.token.Type != lex.TokenRb {
						err := fmt.Errorf("%s '[' and ']' not match", ep.parser.errMsgPrefix())
						ep.parser.errs = append(ep.parser.errs, err)
						return nil, err
					}
					newExpression := &ast.Expression{}
					newExpression.Type = ast.ExpressionTypeSlice
					newExpression.Op = "slice"
					newExpression.Pos = ep.parser.mkPos()
					slice := &ast.ExpressionSlice{}
					newExpression.Data = slice
					slice.Start = e
					slice.ExpressionOn = prefix
					slice.End = end
					prefix = newExpression
					ep.Next(lfIsToken) // skip ]
				} else {
					if ep.parser.token.Type != lex.TokenRb {
						err := fmt.Errorf("%s '[' and ']' not match", ep.parser.errMsgPrefix())
						ep.parser.errs = append(ep.parser.errs, err)
						return nil, err
					}
					newExpression := &ast.Expression{}
					newExpression.Pos = ep.parser.mkPos()
					newExpression.Op = "index"
					newExpression.Type = ast.ExpressionTypeIndex
					index := &ast.ExpressionIndex{}
					index.Expression = prefix
					index.Index = e
					newExpression.Data = index
					prefix = newExpression
					ep.Next(lfIsToken)
				}
			}
		case lex.TokenSelectConst:
			pos := ep.parser.mkPos()
			ep.Next(lfNotToken) // skip ::
			var constName string
			if ep.parser.token.Type != lex.TokenIdentifier {
				ep.parser.errs = append(ep.parser.errs,
					fmt.Errorf("%s expect idnetifier , but '%s'",
						ep.parser.errMsgPrefix(), ep.parser.token.Description))
				constName = compileAutoName()
			} else {
				constName = ep.parser.token.Data.(string)
				ep.Next(lfIsToken)
			}
			newExpression := &ast.Expression{}
			newExpression.Pos = pos
			newExpression.Op = "selectConst"
			newExpression.Type = ast.ExpressionTypeSelectionConst
			selection := &ast.ExpressionSelection{}
			selection.Expression = prefix
			selection.Name = constName
			newExpression.Data = selection
			prefix = newExpression
		case lex.TokenSelection:
			pos := ep.parser.mkPos()
			ep.Next(lfNotToken) // skip .
			if ep.parser.token.Type == lex.TokenIdentifier {
				newExpression := &ast.Expression{}
				newExpression.Pos = pos
				newExpression.Op = "selection"
				newExpression.Type = ast.ExpressionTypeSelection
				selection := &ast.ExpressionSelection{}
				selection.Expression = prefix
				selection.Name = ep.parser.token.Data.(string)
				newExpression.Data = selection
				prefix = newExpression
				ep.Next(lfIsToken)
			} else if ep.parser.token.Type == lex.TokenLp { //  a.(xxx)
				//
				ep.Next(lfNotToken) // skip (
				typ, err := ep.parser.parseType()
				if err != nil {
					return nil, err
				}
				ep.parser.ifTokenIsLfThenSkip()
				if ep.parser.token.Type != lex.TokenRp {
					err := fmt.Errorf("%s '(' and ')' not match", ep.parser.errMsgPrefix())
					ep.parser.errs = append(ep.parser.errs, err)
					return nil, err
				}
				newExpression := &ast.Expression{}
				newExpression.Pos = pos
				newExpression.Op = "assert"
				newExpression.Type = ast.ExpressionTypeTypeAssert
				typeAssert := &ast.ExpressionTypeAssert{}
				typeAssert.Type = typ
				typeAssert.Expression = prefix
				newExpression.Data = typeAssert
				prefix = newExpression
				ep.Next(lfIsToken) // skip  )
			} else {
				err := fmt.Errorf("%s expect  'identifier' or '(',but '%s'",
					ep.parser.errMsgPrefix(), ep.parser.token.Description)
				ep.parser.errs = append(ep.parser.errs, err)
				return nil, err
			}
		case lex.TokenLp:
			newExpression, err := ep.parseCallExpression(prefix)
			if err != nil {
				return nil, err
			}
			prefix = newExpression
		}
	}

	return prefix, nil
}

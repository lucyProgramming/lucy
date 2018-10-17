package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (parser *Parser) parseType() (*ast.Type, error) {
	var err error
	var ret *ast.Type
	pos := parser.mkPos()
	switch parser.token.Type {
	case lex.TokenLb:
		parser.Next(lfIsToken)
		parser.unExpectNewLineAndSkip()
		if parser.token.Type != lex.TokenRb {
			// [ and ] not match
			err = fmt.Errorf("%s '[' and ']' not match", parser.errorMsgPrefix())
			parser.errs = append(parser.errs, err)
			return nil, err
		}
		//lookahead
		parser.Next(lfIsToken) //skip ]
		parser.unExpectNewLineAndSkip()
		array, err := parser.parseType()
		if err != nil {
			return nil, err
		}
		ret = &ast.Type{}
		ret.Pos = pos
		ret.Type = ast.VariableTypeArray
		ret.Array = array
	case lex.TokenBool:
		parser.Next(lfIsToken)
		ret = &ast.Type{
			Type: ast.VariableTypeBool,
			Pos:  pos,
		}
	case lex.TokenByte:
		parser.Next(lfIsToken)
		ret = &ast.Type{
			Type: ast.VariableTypeByte,
			Pos:  pos,
		}
	case lex.TokenShort:
		parser.Next(lfIsToken)
		ret = &ast.Type{
			Type: ast.VariableTypeShort,
			Pos:  pos,
		}
	case lex.TokenChar:
		parser.Next(lfIsToken)
		ret = &ast.Type{
			Type: ast.VariableTypeChar,
			Pos:  pos,
		}
	case lex.TokenInt:
		parser.Next(lfIsToken)
		ret = &ast.Type{
			Type: ast.VariableTypeInt,
			Pos:  pos,
		}
	case lex.TokenFloat:
		parser.Next(lfIsToken)
		ret = &ast.Type{
			Type: ast.VariableTypeFloat,
			Pos:  pos,
		}
	case lex.TokenDouble:
		parser.Next(lfIsToken)
		ret = &ast.Type{
			Type: ast.VariableTypeDouble,
			Pos:  pos,
		}
	case lex.TokenLong:
		parser.Next(lfIsToken)
		ret = &ast.Type{
			Type: ast.VariableTypeLong,
			Pos:  pos,
		}
	case lex.TokenString:
		parser.Next(lfIsToken)
		ret = &ast.Type{
			Type: ast.VariableTypeString,
			Pos:  pos,
		}
	case lex.TokenIdentifier:
		ret, err = parser.parseIdentifierType()
	case lex.TokenMap:
		parser.Next(lfNotToken) // skip map key word
		if parser.token.Type != lex.TokenLc {
			return nil, fmt.Errorf("%s expect '{',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description)
		}
		parser.Next(lfNotToken) // skip {
		var k, v *ast.Type
		k, err = parser.parseType()
		if err != nil {
			return nil, err
		}
		parser.ifTokenIsLfThenSkip()
		if parser.token.Type != lex.TokenArrow {
			return nil, fmt.Errorf("%s expect '->',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description)
		}
		parser.Next(lfNotToken) // skip ->
		v, err := parser.parseType()
		if err != nil {
			return nil, err
		}
		parser.ifTokenIsLfThenSkip()
		if parser.token.Type != lex.TokenRc {
			return nil, fmt.Errorf("%s expect '}',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description)
		}
		parser.Next(lfIsToken)
		m := &ast.Map{
			K: k,
			V: v,
		}
		ret = &ast.Type{
			Type: ast.VariableTypeMap,
			Map:  m,
			Pos:  pos,
		}
	case lex.TokenTemplate:
		ret = &ast.Type{
			Type: ast.VariableTypeTemplate,
			Pos:  pos,
			Name: parser.token.Data.(string),
		}
		parser.Next(lfIsToken)
	case lex.TokenFn:
		parser.Next(lfIsToken)
		ft, err := parser.parseFunctionType()
		if err != nil {
			return nil, err
		}
		ret = &ast.Type{
			Type:         ast.VariableTypeFunction,
			Pos:          pos,
			FunctionType: &ft,
		}
	case lex.TokenGlobal:
		parser.Next(lfIsToken)
		parser.unExpectNewLineAndSkip()
		if parser.token.Type != lex.TokenSelection {
			return nil, fmt.Errorf("%s expect '.' , but '%s'",
				parser.errorMsgPrefix(), parser.token.Description)
		}
		parser.Next(lfNotToken)
		if parser.token.Type != lex.TokenIdentifier {
			parser.errs = append(parser.errs, fmt.Errorf("%s expect identifier , but '%s'",
				parser.errorMsgPrefix(), parser.token.Description))
		} else {
			ret = &ast.Type{
				Type: ast.VariableTypeGlobal,
				Pos:  pos,
				Name: parser.token.Data.(string),
			}
			parser.Next(lfIsToken)
		}
	default:
		return nil, fmt.Errorf("%s unkown begining '%s' token for a type",
			parser.errorMsgPrefix(), parser.token.Description)
	}
	if err != nil {
		parser.errs = append(parser.errs, err)
		return nil, err
	}
	if parser.token.Type == lex.TokenVArgs {
		parser.Next(lfIsToken) // skip ...
		newRet := &ast.Type{
			Pos:            pos,
			Type:           ast.VariableTypeJavaArray,
			Array:          ret,
			IsVariableArgs: true,
		}
		ret = newRet
		return ret, nil
	}
	for parser.token.Type == lex.TokenLb { // int [
		pos := parser.mkPos()
		parser.Next(lfIsToken) // skip [
		parser.unExpectNewLineAndSkip()
		if parser.token.Type != lex.TokenRb {
			err = fmt.Errorf("%s '[' and ']' not match", parser.errorMsgPrefix())
			parser.errs = append(parser.errs, err)
			return ret, err
		}
		parser.Next(lfIsToken) // skip ]
		newRet := &ast.Type{
			Pos:   pos,
			Type:  ast.VariableTypeJavaArray,
			Array: ret,
		}
		ret = newRet
	}
	return ret, err
}

/*
	valid begin token of a type
*/
func (parser *Parser) isValidTypeBegin() bool {
	return parser.token.Type == lex.TokenLb ||
		parser.token.Type == lex.TokenBool ||
		parser.token.Type == lex.TokenByte ||
		parser.token.Type == lex.TokenShort ||
		parser.token.Type == lex.TokenChar ||
		parser.token.Type == lex.TokenInt ||
		parser.token.Type == lex.TokenFloat ||
		parser.token.Type == lex.TokenDouble ||
		parser.token.Type == lex.TokenLong ||
		parser.token.Type == lex.TokenString ||
		parser.token.Type == lex.TokenMap ||
		parser.token.Type == lex.TokenIdentifier ||
		parser.token.Type == lex.TokenTemplate ||
		parser.token.Type == lex.TokenFn
}

func (parser *Parser) parseIdentifierType() (*ast.Type, error) {
	name := parser.token.Data.(string)
	ret := &ast.Type{
		Pos:  parser.mkPos(),
		Type: ast.VariableTypeName,
	}
	parser.Next(lfIsToken) // skip name identifier
	for parser.token.Type == lex.TokenSelection {
		parser.Next(lfNotToken) // skip .
		if parser.token.Type != lex.TokenIdentifier {
			return nil, fmt.Errorf("%s not a identifier after dot",
				parser.errorMsgPrefix())
		}
		name += "." + parser.token.Data.(string)
		parser.Next(lfIsToken) // skip identifier
	}
	ret.Name = name
	return ret, nil
}

func (parser *Parser) parseTypes(endTokens ...lex.TokenKind) ([]*ast.Type, error) {
	ret := []*ast.Type{}
	for parser.token.Type != lex.TokenEof {
		t, err := parser.parseType()
		if err != nil {
			return ret, err
		}
		ret = append(ret, t)
		if parser.token.Type != lex.TokenComma {
			if parser.isValidTypeBegin() {
				parser.errs = append(parser.errs, fmt.Errorf("%s missing comma",
					parser.errorMsgPrefix()))
				continue
			}
			break
		}
		parser.Next(lfNotToken) // skip ,
		for _, v := range endTokens {
			if v == parser.token.Type {
				parser.errs = append(parser.errs, fmt.Errorf("%s extra comma", parser.errorMsgPrefix()))
				goto end
			}
		}
	}
end:
	return ret, nil
}

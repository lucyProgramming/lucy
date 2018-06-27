package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (parser *Parser) parseType() (*ast.Type, error) {
	var err error
	var ret *ast.Type
	switch parser.token.Type {
	case lex.TokenLb:

		//if len(pre) > 0 {
		//	parser.Next()
		//	if parser.token.Type != lex.TOKEN_RB {
		//		// [ and ] not match
		//		err = fmt.Errorf("%s '[' and ']' not match", parser.errorMsgPrefix())
		//		parser.errs = append(parser.errs, err)
		//		return nil, err
		//	}
		//	parser.Next() //skip ]
		//	ret = &ast.Type{
		//		Pos:       pos,
		//		Type:      ast.VARIABLE_TYPE_JAVA_ARRAY,
		//		ArrayType: pre[0],
		//	}
		//	break
		//}
		pos := parser.mkPos()
		parser.Next()
		if parser.token.Type != lex.TokenRb {
			// [ and ] not match
			err = fmt.Errorf("%s '[' and ']' not match", parser.errorMsgPrefix())
			parser.errs = append(parser.errs, err)
			return nil, err
		}
		//lookahead
		parser.Next() //skip ]
		t, err := parser.parseType()
		if err != nil {
			return nil, err
		}
		ret = &ast.Type{}
		ret.Pos = pos
		ret.Type = ast.VariableTypeArray
		ret.ArrayType = t
	case lex.TokenBool:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VariableTypeBool,
			Pos:  pos,
		}
	case lex.TokenByte:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VariableTypeByte,
			Pos:  pos,
		}
	case lex.TokenShort:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VariableTypeShort,
			Pos:  pos,
		}
	case lex.TokenInt:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VariableTypeInt,
			Pos:  pos,
		}
	case lex.TokenFloat:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VariableTypeFloat,
			Pos:  pos,
		}
	case lex.TokenDouble:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VariableTypeDouble,
			Pos:  pos,
		}
	case lex.TokenLong:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VariableTypeLong,
			Pos:  pos,
		}
	case lex.TokenString:
		pos := parser.mkPos()
		parser.Next()
		ret = &ast.Type{
			Type: ast.VariableTypeString,
			Pos:  pos,
		}
	case lex.TokenIdentifier:
		ret, err = parser.parseIdentifierType()

	case lex.TokenMap:
		pos := parser.mkPos()
		parser.Next() // skip map key word
		if parser.token.Type != lex.TokenLc {
			return nil, fmt.Errorf("%s expect '{',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description)
		}
		parser.Next() // skip {
		var k, v *ast.Type
		k, err = parser.parseType()
		if err != nil {
			return nil, err
		}
		if parser.token.Type != lex.TokenArrow {
			return nil, fmt.Errorf("%s expect '->',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description)
		}
		parser.Next() // skip ->
		v, err := parser.parseType()
		if err != nil {
			return nil, err
		}
		if parser.token.Type != lex.TokenRc {
			return nil, fmt.Errorf("%s expect '}',but '%s'",
				parser.errorMsgPrefix(), parser.token.Description)
		}
		parser.Next()
		m := &ast.Map{
			Key:   k,
			Value: v,
		}
		ret = &ast.Type{
			Type: ast.VariableTypeMap,
			Map:  m,
			Pos:  pos,
		}
	case lex.TokenTemplate:
		pos := parser.mkPos()
		ret = &ast.Type{
			Type: ast.VariableTypeTemplate,
			Pos:  pos,
			Name: parser.token.Data.(string),
		}
		parser.Next()
	case lex.TokenFunction:
		pos := parser.mkPos()
		parser.Next()
		ft, err := parser.parseFunctionType()
		if err != nil {
			return nil, err
		}
		ret = &ast.Type{
			Type:         ast.VariableTypeFunction,
			Pos:          pos,
			FunctionType: &ft,
		}
	default:
		err = fmt.Errorf("%s unkown type,begining token is '%s'",
			parser.errorMsgPrefix(), parser.token.Description)
	}
	if err != nil {
		parser.errs = append(parser.errs, err)
		return nil, err
	}
	for parser.token.Type == lex.TokenLb {
		pos := parser.mkPos()
		parser.Next() // skip [
		if parser.token.Type != lex.TokenRb {
			err = fmt.Errorf("%s '[' and ']' not match", parser.errorMsgPrefix())
			parser.errs = append(parser.errs, err)
			return ret, err
		}
		parser.Next() // skip ]
		newRet := &ast.Type{
			Pos:       pos,
			Type:      ast.VariableTypeJavaArray,
			ArrayType: ret,
		}
		ret = newRet
	}
	return ret, err

}

func (parser *Parser) isValidTypeBegin() bool {
	return parser.token.Type == lex.TokenLb ||
		parser.token.Type == lex.TokenBool ||
		parser.token.Type == lex.TokenByte ||
		parser.token.Type == lex.TokenShort ||
		parser.token.Type == lex.TokenInt ||
		parser.token.Type == lex.TokenFloat ||
		parser.token.Type == lex.TokenDouble ||
		parser.token.Type == lex.TokenLong ||
		parser.token.Type == lex.TokenString ||
		parser.token.Type == lex.TokenMap ||
		parser.token.Type == lex.TokenIdentifier ||
		parser.token.Type == lex.TokenTemplate

}
func (parser *Parser) parseIdentifierType() (*ast.Type, error) {
	name := parser.token.Data.(string)
	ret := &ast.Type{
		Pos:  parser.mkPos(),
		Type: ast.VariableTypeName,
	}
	parser.Next() // skip name identifier
	for parser.token.Type == lex.TokenSelection {
		parser.Next() // skip .
		if parser.token.Type != lex.TokenIdentifier {
			return nil, fmt.Errorf("%s not a identifier after dot",
				parser.errorMsgPrefix())
		}
		name += "." + parser.token.Data.(string)
		parser.Next() // if
	}
	ret.Name = name
	return ret, nil
}

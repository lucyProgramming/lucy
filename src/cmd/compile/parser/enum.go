package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (parser *Parser) parseEnum() (e *ast.Enum, err error) {
	enumName := &ast.NameWithPos{
		Pos: parser.mkPos(),
	}
	parser.Next(lfIsToken) // skip enum
	parser.unExpectNewLineAndSkip()
	if parser.token.Type != lex.TokenIdentifier {
		err = fmt.Errorf("%s expect 'identifier' for enum name, but '%s'",
			parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		enumName.Name = compileAutoName()
		parser.consume(untilLc)
	} else {
		enumName.Name = parser.token.Data.(string)
		parser.Next(lfNotToken) // skip enum name
	}
	if parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'",
			parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return nil, err
	}
	parser.Next(lfNotToken) // skip {
	e = &ast.Enum{}
	e.Name = enumName.Name
	e.Pos = enumName.Pos
	//first name
	if parser.token.Type != lex.TokenIdentifier {
		err = fmt.Errorf("%s expect 'identifier',but '%s'",
			parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return nil, err
	}
	names := []*ast.NameWithPos{
		&ast.NameWithPos{
			Name: parser.token.Data.(string),
			Pos:  parser.mkPos(),
		},
	}
	parser.Next(lfIsToken) // skip first name
	var initExpression *ast.Expression
	if parser.token.Type == lex.TokenAssign { // first value defined here
		parser.Next(lfNotToken) // skip assign
		initExpression, err = parser.ExpressionParser.parseExpression(false)
		if err != nil {
			parser.errs = append(parser.errs, err)
			return nil, err
		}
	}
	if parser.token.Type == lex.TokenComma {
		parser.Next(lfNotToken) // skip ,should be a identifier after  comma
		ns, err := parser.parseNameList()
		if err != nil {
			parser.consume(untilRc)
		}
		if ns != nil {
			names = append(names, ns...)
		}
	}
	parser.ifTokenIsLfThenSkip()
	if parser.token.Type != lex.TokenRc {
		err = fmt.Errorf("%s expect '}',but '%s'", parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		parser.consume(untilRc)
	}
	parser.Next(lfNotToken) // skip }
	e.Init = initExpression
	for _, v := range names {
		t := &ast.EnumName{}
		t.Name = v.Name
		t.Pos = v.Pos
		t.Enum = e
		e.Enums = append(e.Enums, t)
	}
	return e, err
}

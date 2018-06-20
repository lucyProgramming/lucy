package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (parser *Parser) parseEnum(isPublic bool) (e *ast.Enum, err error) {
	parser.Next() // skip enum

	if parser.token.Type != lex.TOKEN_IDENTIFIER {
		err = fmt.Errorf("%s expect 'identifier', but '%s'",
			parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return nil, err
	}
	enumName := &ast.NameWithPos{
		Name: parser.token.Data.(string),
		Pos:  parser.mkPos(),
	}
	parser.Next() // skip enum name
	if parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s expect '{',but '%s'",
			parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return nil, err
	}
	parser.Next() // skip {
	e = &ast.Enum{}
	e.Name = enumName.Name
	e.Pos = enumName.Pos
	//first name
	if parser.token.Type != lex.TOKEN_IDENTIFIER {
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
	parser.Next()
	var initExpression *ast.Expression
	if parser.token.Type == lex.TOKEN_ASSIGN { // first value defined here
		parser.Next() // skip assign
		initExpression, err = parser.ExpressionParser.parseExpression(false)
		if err != nil {
			parser.errs = append(parser.errs, err)
			return nil, err
		}
	}
	if parser.token.Type == lex.TOKEN_COMMA {
		parser.Next() // skip ,should be a identifier after  comma
		ns, err := parser.parseNameList()
		if err != nil {
			return nil, err
		}
		names = append(names, ns...)
	}
	if parser.token.Type != lex.TOKEN_RC {
		err = fmt.Errorf("%s expect '}',but '%s'", parser.token.Description, parser.token.Description)
		parser.errs = append(parser.errs, err)
		return nil, err
	}
	parser.Next() // skip }
	e.Init = initExpression
	for _, v := range names {
		t := &ast.EnumName{}
		t.Name = v.Name
		t.Pos = v.Pos
		t.Enum = e
		e.Enums = append(e.Enums, t)
	}
	e.AccessFlags = 0
	if isPublic {
		e.AccessFlags |= cg.ACC_CLASS_PUBLIC
	}
	return e, err
}

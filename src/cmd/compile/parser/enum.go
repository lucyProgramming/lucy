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
	e = &ast.Enum{}
	e.Name = enumName.Name
	e.Pos = enumName.Pos
	comment := &CommentParser{
		parser: parser,
	}
	reset := func() {
		comment.reset()
	}
	if parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'", parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		parser.consume(untilLc)
	}
	parser.Next(lfNotToken)
	for parser.token.Type != lex.TokenRc &&
		parser.token.Type != lex.TokenEof {
		switch parser.token.Type {
		case lex.TokenLf:
			parser.Next(lfNotToken)
		case lex.TokenCommentMultiLine,
			lex.TokenComment:
			comment.read()
		case lex.TokenIdentifier:
			name := parser.token.Data.(string)
			pos := parser.mkPos()
			var value *ast.Expression
			var err error
			parser.Next(lfIsToken)
			if parser.token.Type == lex.TokenAssign {
				parser.Next(lfNotToken)
				value, err = parser.ExpressionParser.parseExpression(false)
				if err != nil {
					parser.errs = append(parser.errs, err)
					parser.consume(untilSemicolonOrLf)
				}
			}
			enumComment := comment.Comment
			if parser.token.Type == lex.TokenComment {
				enumComment = parser.token.Data.(string)
				parser.Next(lfIsToken)
			}
			if len(e.Enums) == 0 {
				e.Init = value
				value = nil
			}
			enumName := &ast.EnumName{
				Name:    name,
				Pos:     pos,
				NoNeed:  value,
				Enum:    e,
				Comment: enumComment,
			}
			e.Enums = append(e.Enums, enumName)
			reset()
		case lex.TokenComma:
			parser.Next(lfNotToken)
		default:
			parser.errs = append(parser.errs, fmt.Errorf("%s token '%s' is not except",
				parser.errorMsgPrefix(), parser.token.Description))
			parser.consume(untilSemicolonOrLf)
			parser.Next(lfNotToken)
			reset()
		}
	}
	if len(e.Enums) == 0 {
		enumName := &ast.EnumName{
			Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", //
			Pos:  parser.mkPos(),
			Enum: e,
		}
		e.Enums = []*ast.EnumName{
			enumName,
		}
		parser.errs = append(parser.errs, fmt.Errorf("%s enum expect at least 1 enumName",
			parser.errorMsgPrefix()))
	}
	parser.ifTokenIsLfThenSkip()
	if parser.token.Type != lex.TokenRc {
		err = fmt.Errorf("%s expect '}',but '%s'", parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		parser.consume(untilRc)
	}
	parser.Next(lfNotToken)
	return e, err
}

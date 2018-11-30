package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (this *Parser) parseEnum() (e *ast.Enum, err error) {
	var enumName string
	this.Next(lfIsToken) // skip enum
	this.unExpectNewLineAndSkip()
	if this.token.Type != lex.TokenIdentifier {
		err = fmt.Errorf("%s expect 'identifier' for enum name, but '%s'",
			this.errMsgPrefix(), this.token.Description)
		this.errs = append(this.errs, err)
		enumName = compileAutoName()
		this.consume(untilLc)

	} else {
		enumName = this.token.Data.(string)
		this.Next(lfNotToken) // skip enum name
	}
	e = &ast.Enum{}
	e.Name = enumName
	e.Pos = this.mkPos()
	comment := &CommentParser{
		parser: this,
	}
	reset := func() {
		comment.reset()
	}
	if this.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'", this.errMsgPrefix(), this.token.Description)
		this.errs = append(this.errs, err)
		this.consume(untilLc)
	}
	this.Next(lfNotToken)
	for this.token.Type != lex.TokenRc &&
		this.token.Type != lex.TokenEof {
		switch this.token.Type {
		case lex.TokenLf:
			this.Next(lfNotToken)
		case lex.TokenMultiLineComment,
			lex.TokenComment:
			
			comment.read()
		case lex.TokenIdentifier:
			name := this.token.Data.(string)
			pos := this.mkPos()
			var value *ast.Expression
			var err error
			this.Next(lfIsToken)
			if this.token.Type == lex.TokenAssign {
				this.Next(lfNotToken)
				value, err = this.ExpressionParser.parseExpression(false)
				if err != nil {
					this.consume(untilSemicolonOrLf)
				}
			}
			enumComment := comment.Comment
			if e.Init == nil && value != nil {
				e.Init = value
				e.FirstValueIndex = len(e.Enums)
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
			this.Next(lfNotToken)
		default:
			this.errs = append(this.errs, fmt.Errorf("%s token '%s' is not except",
				this.errMsgPrefix(), this.token.Description))
			this.Next(lfNotToken)
			reset()
		}
	}
	if len(e.Enums) == 0 {
		enumName := &ast.EnumName{
			Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", //easter egg
			Pos:  this.mkPos(),
			Enum: e,
		}
		e.Enums = []*ast.EnumName{
			enumName,
		}
		this.errs = append(this.errs, fmt.Errorf("%s enum expect at least 1 enumName",
			this.errMsgPrefix()))
	}
	this.ifTokenIsLfThenSkip()
	if this.token.Type != lex.TokenRc {
		err = fmt.Errorf("%s expect '}',but '%s'", this.errMsgPrefix(), this.token.Description)
		this.errs = append(this.errs, err)
		this.consume(untilRc)
	}
	this.Next(lfNotToken)
	return e, err
}

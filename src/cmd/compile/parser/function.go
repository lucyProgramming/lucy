package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type FunctionParser struct {
	parser *Parser
}

func (this *FunctionParser) Next(lfIsToken bool) {
	this.parser.Next(lfIsToken)
}

func (this *FunctionParser) consume(until map[lex.TokenKind]bool) {
	this.parser.consume(until)
}

/*
	when canBeAbstract is true , means can have no body
*/
func (this *FunctionParser) parse(needName bool, isAbstract bool) (f *ast.Function, err error) {
	f = &ast.Function{}
	offset := this.parser.token.Offset
	this.Next(lfIsToken) // skip fn key word
	this.parser.unExpectNewLineAndSkip()
	if needName && this.parser.token.Type != lex.TokenIdentifier {
		err := fmt.Errorf("%s expect function name,but '%s'",
			this.parser.errMsgPrefix(), this.parser.token.Description)
		this.parser.errs = append(this.parser.errs, err)
		this.consume(untilLp)
	}
	f.Pos = this.parser.mkPos()
	if this.parser.token.Type == lex.TokenIdentifier {
		f.Name = this.parser.token.Data.(string)
		this.Next(lfNotToken)
	}
	f.Type, err = this.parser.parseFunctionType()
	if err != nil {
		if isAbstract {
			this.consume(untilSemicolonOrLf)
		} else {
			this.consume(untilLc)
		}
	}
	if isAbstract {
		return f, nil
	}
	this.parser.ifTokenIsLfThenSkip()
	if this.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s except '{' but '%s'",
			this.parser.errMsgPrefix(), this.parser.token.Description)
		this.parser.errs = append(this.parser.errs, err)
		this.consume(untilLc)
	}
	f.Block.IsFunctionBlock = true
	f.Block.Fn = f
	this.Next(lfNotToken) // skip {
	this.parser.BlockParser.parseStatementList(&f.Block, false)
	if this.parser.token.Type != lex.TokenRc {
		err = fmt.Errorf("%s expect '}', but '%s'",
			this.parser.errMsgPrefix(), this.parser.token.Description)
		this.parser.errs = append(this.parser.errs, err)
		this.consume(untilRc)
	} else {
		f.SourceCode =
			this.parser.
				bs[offset : this.parser.token.Offset+1]
	}
	this.Next(lfIsToken)
	return f, err
}

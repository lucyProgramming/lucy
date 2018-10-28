package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type FunctionParser struct {
	parser *Parser
}

func (functionParser *FunctionParser) Next(lfIsToken bool) {
	functionParser.parser.Next(lfIsToken)
}

func (functionParser *FunctionParser) consume(until map[lex.TokenKind]bool) {
	functionParser.parser.consume(until)
}

/*
	when canBeAbstract is true , means can have no body
*/
func (functionParser *FunctionParser) parse(needName bool, isAbstract bool) (f *ast.Function, err error) {
	f = &ast.Function{}
	offset := functionParser.parser.token.Offset
	functionParser.Next(lfIsToken) // skip fn key word
	functionParser.parser.unExpectNewLineAndSkip()
	if needName && functionParser.parser.token.Type != lex.TokenIdentifier {
		err := fmt.Errorf("%s expect function name,but '%s'",
			functionParser.parser.errMsgPrefix(), functionParser.parser.token.Description)
		functionParser.parser.errs = append(functionParser.parser.errs, err)
		functionParser.consume(untilLp)
	}
	f.Pos = functionParser.parser.mkPos()
	if functionParser.parser.token.Type == lex.TokenIdentifier {
		f.Name = functionParser.parser.token.Data.(string)
		functionParser.Next(lfNotToken)
	}
	f.Type, err = functionParser.parser.parseFunctionType()
	if err != nil {
		if isAbstract {
			functionParser.consume(untilSemicolonOrLf)
		} else {
			functionParser.consume(untilLc)
		}
	}
	if isAbstract {
		return f, nil
	}
	functionParser.parser.ifTokenIsLfThenSkip()
	if functionParser.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s except '{' but '%s'",
			functionParser.parser.errMsgPrefix(), functionParser.parser.token.Description)
		functionParser.parser.errs = append(functionParser.parser.errs, err)
		functionParser.consume(untilLc)
	}
	f.Block.IsFunctionBlock = true
	f.Block.Fn = f
	functionParser.Next(lfNotToken) // skip {
	functionParser.parser.BlockParser.parseStatementList(&f.Block, false)
	if functionParser.parser.token.Type != lex.TokenRc {
		err = fmt.Errorf("%s expect '}', but '%s'",
			functionParser.parser.errMsgPrefix(), functionParser.parser.token.Description)
		functionParser.parser.errs = append(functionParser.parser.errs, err)
		functionParser.consume(untilRc)
	} else {
		f.SourceCode =
			functionParser.parser.
				bs[offset : functionParser.parser.token.Offset+1]
	}
	functionParser.Next(lfIsToken)
	return f, err
}

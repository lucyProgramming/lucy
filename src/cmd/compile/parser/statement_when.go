package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (bp *BlockParser) parseWhen() (*ast.StatementWhen, error) {
	bp.parser.Next(lfIsToken)
	bp.parser.unExpectNewLineAndSkip()
	condition, err := bp.parser.parseType()
	if err != nil {
		bp.parser.errs = append(bp.parser.errs, err)
		bp.consume(untilLc)
	}
	bp.parser.ifTokenIsLfThenSkip()
	if bp.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'",
			bp.parser.errMsgPrefix(), bp.parser.token.Description)
		bp.parser.errs = append(bp.parser.errs, err)
		bp.consume(untilLc)
	}
	bp.Next(lfNotToken) // skip {  , must be case
	if bp.parser.token.Type != lex.TokenCase {
		err = fmt.Errorf("%s expect 'case',but '%s'",
			bp.parser.errMsgPrefix(), bp.parser.token.Description)
		bp.parser.errs = append(bp.parser.errs, err)
		return nil, err
	}
	when := &ast.StatementWhen{}
	when.Condition = condition
	for bp.parser.token.Type == lex.TokenCase {
		bp.Next(lfIsToken) // skip case
		bp.parser.unExpectNewLineAndSkip()
		ts, err := bp.parser.parseTypes(lex.TokenColon)
		if err != nil {
			bp.parser.errs = append(bp.parser.errs, err)
			return when, err
		}
		bp.parser.unExpectNewLineAndSkip()
		if bp.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s expect ':',but '%s'",
				bp.parser.errMsgPrefix(), bp.parser.token.Description)
			bp.parser.errs = append(bp.parser.errs, err)
			return when, err
		}
		bp.Next(lfIsToken) // skip :
		bp.parser.expectNewLineAndSkip()
		var block *ast.Block
		if bp.parser.token.Type != lex.TokenCase &&
			bp.parser.token.Type != lex.TokenDefault &&
			bp.parser.token.Type != lex.TokenRc {
			block = &ast.Block{}
			block.IsSwitchBlock = true
			bp.parseStatementList(block, false)
		}
		when.Cases =
			append(when.Cases, &ast.StatementWhenCase{
				Matches: ts,
				Block:   block,
			})
	}
	//default value
	if bp.parser.token.Type == lex.TokenDefault {
		bp.Next(lfIsToken) // skip default key word
		bp.parser.unExpectNewLineAndSkip()
		if bp.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s missing colon after default",
				bp.parser.errMsgPrefix())
			bp.parser.errs = append(bp.parser.errs, err)
		} else {
			bp.Next(lfIsToken)
			bp.parser.expectNewLineAndSkip()
		}
		if bp.parser.token.Type != lex.TokenRc {
			block := ast.Block{}
			block.IsSwitchBlock = true
			bp.parseStatementList(&block, false)
			when.Default = &block
		}
	}
	if bp.parser.token.Type != lex.TokenRc {
		err = fmt.Errorf("%s expect '}',but '%s'",
			bp.parser.errMsgPrefix(), bp.parser.token.Description)
		bp.parser.errs = append(bp.parser.errs, err)
		return when, err
	}
	bp.Next(lfNotToken) //  skip }
	return when, nil
}

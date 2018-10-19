package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//TODO:: missing format
func (blockParser *BlockParser) parseWhen() (*ast.StatementWhen, error) {
	blockParser.parser.Next(lfIsToken)
	blockParser.parser.unExpectNewLineAndSkip()
	condition, err := blockParser.parser.parseType()
	if err != nil {
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		blockParser.consume(untilLc)
	}
	if blockParser.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'",
			blockParser.parser.errMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		blockParser.consume(untilLc)
	}
	blockParser.Next(lfNotToken) // skip {  , must be case
	if blockParser.parser.token.Type != lex.TokenCase {
		err = fmt.Errorf("%s expect 'case',but '%s'",
			blockParser.parser.errMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return nil, err
	}
	when := &ast.StatementWhen{}
	when.Condition = condition
	for blockParser.parser.token.Type == lex.TokenCase {
		blockParser.Next(lfNotToken) // skip case
		ts, err := blockParser.parser.parseTypes(lex.TokenColon)
		if err != nil {
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			return when, err
		}
		if blockParser.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s expect ':',but '%s'",
				blockParser.parser.errMsgPrefix(), blockParser.parser.token.Description)
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			return when, err
		}
		blockParser.Next(lfNotToken) // skip :
		var block *ast.Block
		if blockParser.parser.token.Type != lex.TokenCase &&
			blockParser.parser.token.Type != lex.TokenDefault &&
			blockParser.parser.token.Type != lex.TokenRc {
			block = &ast.Block{}
			block.IsSwitchBlock = true
			blockParser.parseStatementList(block, false)
		}
		when.Cases =
			append(when.Cases, &ast.StatementWhenIs{
				Matches: ts,
				Block:   block,
			})
	}
	//default value
	if blockParser.parser.token.Type == lex.TokenDefault {
		blockParser.Next(lfIsToken) // skip default key word
		if blockParser.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s missing colon after default",
				blockParser.parser.errMsgPrefix())
			blockParser.parser.errs = append(blockParser.parser.errs, err)
		} else {
			blockParser.Next(lfNotToken)
		}
		if blockParser.parser.token.Type != lex.TokenRc {
			block := ast.Block{}
			block.IsSwitchBlock = true
			blockParser.parseStatementList(&block, false)
			when.Default = &block
		}
	}
	if blockParser.parser.token.Type != lex.TokenRc {
		err = fmt.Errorf("%s expect '}',but '%s'",
			blockParser.parser.errMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return when, err
	}
	blockParser.Next(lfNotToken) //  skip }
	return when, nil
}
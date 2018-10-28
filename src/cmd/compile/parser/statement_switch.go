package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (blockParser *BlockParser) parseSwitch() (*ast.StatementSwitch, error) {
	statementSwitch := &ast.StatementSwitch{
		Pos: blockParser.parser.mkPos(),
	}
	blockParser.Next(lfIsToken) // skip switch key word
	blockParser.parser.unExpectNewLineAndSkip()
	statementSwitch.EndPos = blockParser.parser.mkPos()
	var err error
	statementSwitch.Condition, err = blockParser.parser.ExpressionParser.parseExpression(false)
	if err != nil {
		blockParser.consume(untilLc)
	}
	blockParser.parser.ifTokenIsLfThenSkip()
	for blockParser.parser.token.Type == lex.TokenSemicolon {
		if statementSwitch.Condition != nil {
			statementSwitch.PrefixExpressions = append(statementSwitch.PrefixExpressions, statementSwitch.Condition)
			statementSwitch.Condition = nil
		}
		blockParser.parser.Next(lfNotToken)
		statementSwitch.Condition, err = blockParser.parser.ExpressionParser.parseExpression(false)
		if err != nil {
			blockParser.consume(untilLc)
		}
		blockParser.parser.ifTokenIsLfThenSkip()
	}
	if blockParser.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'",
			blockParser.parser.errMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		blockParser.consume(untilLc)
	}
	blockParser.Next(lfIsToken) // skip {  , must be case
	blockParser.parser.expectNewLineAndSkip()
	if blockParser.parser.token.Type != lex.TokenCase {
		err = fmt.Errorf("%s expect 'case',but '%s'",
			blockParser.parser.errMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return nil, err
	}
	for blockParser.parser.token.Type == lex.TokenCase {
		blockParser.Next(lfIsToken) // skip case
		blockParser.parser.unExpectNewLineAndSkip()
		es, err := blockParser.parser.ExpressionParser.parseExpressions(lex.TokenColon)
		if err != nil {
			return statementSwitch, err
		}
		if blockParser.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s expect ':',but '%s'",
				blockParser.parser.errMsgPrefix(), blockParser.parser.token.Description)
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			return statementSwitch, err
		}
		blockParser.Next(lfIsToken) // skip :
		blockParser.parser.expectNewLineAndSkip()
		var block *ast.Block
		if blockParser.parser.token.Type != lex.TokenCase &&
			blockParser.parser.token.Type != lex.TokenDefault &&
			blockParser.parser.token.Type != lex.TokenRc {
			block = &ast.Block{}
			block.IsSwitchBlock = true
			blockParser.parseStatementList(block, false)
		}
		statementSwitch.StatementSwitchCases = append(statementSwitch.StatementSwitchCases, &ast.StatementSwitchCase{
			Matches: es,
			Block:   block,
		})
	}
	//default value
	if blockParser.parser.token.Type == lex.TokenDefault {
		blockParser.Next(lfIsToken) // skip default key word
		blockParser.parser.unExpectNewLineAndSkip()
		if blockParser.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s missing colon after 'default'",
				blockParser.parser.errMsgPrefix())
			blockParser.parser.errs = append(blockParser.parser.errs, err)
		} else {
			blockParser.Next(lfIsToken)
		}
		blockParser.parser.expectNewLineAndSkip()
		if blockParser.parser.token.Type != lex.TokenRc {
			block := ast.Block{}
			block.IsSwitchBlock = true
			blockParser.parseStatementList(&block, false)
			statementSwitch.Default = &block
		}
	}
	if blockParser.parser.token.Type != lex.TokenRc {
		err = fmt.Errorf("%s expect '}',but '%s'",
			blockParser.parser.errMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return statementSwitch, err
	}
	statementSwitch.EndPos = blockParser.parser.mkEndPos()
	blockParser.Next(lfNotToken) //  skip }
	return statementSwitch, nil
}

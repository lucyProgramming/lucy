package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (bp *BlockParser) parseSwitch() (*ast.StatementSwitch, error) {
	statementSwitch := &ast.StatementSwitch{
		Pos: bp.parser.mkPos(),
	}
	bp.Next(lfIsToken) // skip switch key word
	bp.parser.unExpectNewLineAndSkip()
	statementSwitch.EndPos = bp.parser.mkPos()
	var err error
	statementSwitch.Condition, err = bp.parser.ExpressionParser.parseExpression(false)
	if err != nil {
		bp.consume(untilLc)
	}
	bp.parser.ifTokenIsLfThenSkip()
	for bp.parser.token.Type == lex.TokenSemicolon {
		if statementSwitch.Condition != nil {
			statementSwitch.PrefixExpressions = append(statementSwitch.PrefixExpressions, statementSwitch.Condition)
			statementSwitch.Condition = nil
		}
		bp.parser.Next(lfNotToken)
		statementSwitch.Condition, err = bp.parser.ExpressionParser.parseExpression(false)
		if err != nil {
			bp.consume(untilLc)
		}
		bp.parser.ifTokenIsLfThenSkip()
	}
	if bp.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'",
			bp.parser.errMsgPrefix(), bp.parser.token.Description)
		bp.parser.errs = append(bp.parser.errs, err)
		bp.consume(untilLc)
	}
	bp.Next(lfIsToken) // skip {  , must be case
	bp.parser.expectNewLineAndSkip()
	if bp.parser.token.Type != lex.TokenCase {
		err = fmt.Errorf("%s expect 'case',but '%s'",
			bp.parser.errMsgPrefix(), bp.parser.token.Description)
		bp.parser.errs = append(bp.parser.errs, err)
		return nil, err
	}
	for bp.parser.token.Type == lex.TokenCase {
		bp.Next(lfIsToken) // skip case
		bp.parser.unExpectNewLineAndSkip()
		es, err := bp.parser.ExpressionParser.parseExpressions(lex.TokenColon)
		if err != nil {
			return statementSwitch, err
		}
		if bp.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s expect ':',but '%s'",
				bp.parser.errMsgPrefix(), bp.parser.token.Description)
			bp.parser.errs = append(bp.parser.errs, err)
			return statementSwitch, err
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
		statementSwitch.StatementSwitchCases = append(statementSwitch.StatementSwitchCases, &ast.StatementSwitchCase{
			Matches: es,
			Block:   block,
		})
	}
	//default value
	if bp.parser.token.Type == lex.TokenDefault {
		bp.Next(lfIsToken) // skip default key word
		bp.parser.unExpectNewLineAndSkip()
		if bp.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s missing colon after 'default'",
				bp.parser.errMsgPrefix())
			bp.parser.errs = append(bp.parser.errs, err)
		} else {
			bp.Next(lfIsToken)
		}
		bp.parser.expectNewLineAndSkip()
		if bp.parser.token.Type != lex.TokenRc {
			block := ast.Block{}
			block.IsSwitchBlock = true
			bp.parseStatementList(&block, false)
			statementSwitch.Default = &block
		}
	}
	if bp.parser.token.Type != lex.TokenRc {
		err = fmt.Errorf("%s expect '}',but '%s'",
			bp.parser.errMsgPrefix(), bp.parser.token.Description)
		bp.parser.errs = append(bp.parser.errs, err)
		return statementSwitch, err
	}
	statementSwitch.EndPos = bp.parser.mkEndPos()
	bp.Next(lfNotToken) //  skip }
	return statementSwitch, nil
}

package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (this *BlockParser) parseSwitch() (*ast.StatementSwitch, error) {
	statementSwitch := &ast.StatementSwitch{
		Pos: this.parser.mkPos(),
	}
	this.Next(lfIsToken) // skip switch key word
	this.parser.unExpectNewLineAndSkip()
	statementSwitch.EndPos = this.parser.mkPos()
	var err error
	statementSwitch.Condition, err = this.parser.ExpressionParser.parseExpression(false)
	if err != nil {
		this.consume(untilLc)
	}
	this.parser.ifTokenIsLfThenSkip()
	for this.parser.token.Type == lex.TokenSemicolon {
		if statementSwitch.Condition != nil {
			statementSwitch.PrefixExpressions = append(statementSwitch.PrefixExpressions, statementSwitch.Condition)
			statementSwitch.Condition = nil
		}
		this.parser.Next(lfNotToken)
		statementSwitch.Condition, err = this.parser.ExpressionParser.parseExpression(false)
		if err != nil {
			this.consume(untilLc)
		}
		this.parser.ifTokenIsLfThenSkip()
	}
	if this.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'",
			this.parser.errMsgPrefix(), this.parser.token.Description)
		this.parser.errs = append(this.parser.errs, err)
		this.consume(untilLc)
	}
	this.Next(lfIsToken) // skip {  , must be case
	this.parser.expectNewLineAndSkip()
	if this.parser.token.Type != lex.TokenCase {
		err = fmt.Errorf("%s expect 'case',but '%s'",
			this.parser.errMsgPrefix(), this.parser.token.Description)
		this.parser.errs = append(this.parser.errs, err)
		return nil, err
	}
	for this.parser.token.Type == lex.TokenCase {
		this.Next(lfIsToken) // skip case
		this.parser.unExpectNewLineAndSkip()
		es, err := this.parser.ExpressionParser.parseExpressions(lex.TokenColon)
		if err != nil {
			return statementSwitch, err
		}
		if this.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s expect ':',but '%s'",
				this.parser.errMsgPrefix(), this.parser.token.Description)
			this.parser.errs = append(this.parser.errs, err)
			return statementSwitch, err
		}
		this.Next(lfIsToken) // skip :
		this.parser.expectNewLineAndSkip()
		var block *ast.Block
		if this.parser.token.Type != lex.TokenCase &&
			this.parser.token.Type != lex.TokenDefault &&
			this.parser.token.Type != lex.TokenRc {
			block = &ast.Block{}
			block.IsSwitchBlock = true
			this.parseStatementList(block, false)
		}
		statementSwitch.StatementSwitchCases = append(statementSwitch.StatementSwitchCases, &ast.StatementSwitchCase{
			Matches: es,
			Block:   block,
		})
	}
	//default value
	if this.parser.token.Type == lex.TokenDefault {
		this.Next(lfIsToken) // skip default key word
		this.parser.unExpectNewLineAndSkip()
		if this.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s missing colon after 'default'",
				this.parser.errMsgPrefix())
			this.parser.errs = append(this.parser.errs, err)
		} else {
			this.Next(lfIsToken)
		}
		this.parser.expectNewLineAndSkip()
		block := &ast.Block{}
		block.IsSwitchBlock = true
		statementSwitch.Default = block
		if this.parser.token.Type != lex.TokenRc {
			this.parseStatementList(block, false)
		}
	}
	if this.parser.token.Type != lex.TokenRc {
		err = fmt.Errorf("%s expect '}',but '%s'",
			this.parser.errMsgPrefix(), this.parser.token.Description)
		this.parser.errs = append(this.parser.errs, err)
		return statementSwitch, err
	}
	statementSwitch.EndPos = this.parser.mkEndPos()
	this.Next(lfNotToken) //  skip }
	return statementSwitch, nil
}

package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (this *BlockParser) parseFor() (statementFor *ast.StatementFor, err error) {
	statementFor = &ast.StatementFor{}
	statementFor.Pos = this.parser.mkPos()
	statementFor.Block = &ast.Block{}
	this.Next(lfIsToken) // skip for
	this.parser.unExpectNewLineAndSkip()
	if this.parser.token.Type != lex.TokenLc &&
		this.parser.token.Type != lex.TokenSemicolon { // not '{' and not ';'
		statementFor.Condition, err = this.parser.ExpressionParser.parseExpression(true)
		if err != nil {
			this.consume(untilLc)
			goto parseBlock
		}
	}
	if this.parser.token.Type == lex.TokenSemicolon {
		this.Next(lfNotToken) // skip ;
		statementFor.Init = statementFor.Condition
		statementFor.Condition = nil // mk nil
		//condition
		var err error
		if this.parser.token.Type != lex.TokenSemicolon {
			statementFor.Condition, err = this.parser.ExpressionParser.parseExpression(false)
			if err != nil {
				this.consume(untilLc)
				goto parseBlock
			}
			if this.parser.token.Type != lex.TokenSemicolon {
				this.parser.errs = append(this.parser.errs, fmt.Errorf("%s missing semicolon after expression",
					this.parser.errMsgPrefix()))
				this.consume(untilLc)
				goto parseBlock
			}
		}
		this.Next(lfIsToken)
		this.parser.unExpectNewLineAndSkip()
		if this.parser.token.Type != lex.TokenLc {
			statementFor.Increment, err = this.parser.ExpressionParser.parseExpression(true)
			if err != nil {
				this.consume(untilLc)
				goto parseBlock
			}
		}
	}
parseBlock:
	this.parser.ifTokenIsLfThenSkip()
	if this.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'",
			this.parser.errMsgPrefix(), this.parser.token.Description)
		this.parser.errs = append(this.parser.errs, err)
		return
	}
	this.Next(lfNotToken) // skip {
	this.parseStatementList(statementFor.Block, false)
	if this.parser.token.Type != lex.TokenRc {
		this.parser.errs = append(this.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
			this.parser.errMsgPrefix(), this.parser.token.Description))
		this.consume(untilRc)
	}
	this.Next(lfNotToken) // skip }
	return statementFor, nil
}

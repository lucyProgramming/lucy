package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (this *BlockParser) parseIf() (statementIf *ast.StatementIf, err error) {
	statementIf = &ast.StatementIf{
		Pos: this.parser.mkPos(),
	}
	this.Next(lfIsToken) // skip if
	var condition *ast.Expression
	this.parser.unExpectNewLineAndSkip()
	condition, err = this.parser.ExpressionParser.parseExpression(true)
	if err != nil {
		this.consume(untilLc)
		this.Next(lfNotToken)
	}
	statementIf.Condition = condition
	this.parser.ifTokenIsLfThenSkip()
	for this.parser.token.Type == lex.TokenSemicolon {
		if statementIf.Condition != nil {
			statementIf.PrefixExpressions = append(statementIf.PrefixExpressions, statementIf.Condition)
		}
		this.Next(lfNotToken) // skip ;
		statementIf.Condition, err = this.parser.ExpressionParser.parseExpression(true)
		if err != nil {
			this.consume(untilLc)
			this.Next(lfNotToken)
		}
	}
	this.parser.ifTokenIsLfThenSkip()
	if this.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s missing '{' after condtion,but '%s'",
			this.parser.errMsgPrefix(), this.parser.token.Description)
		this.parser.errs = append(this.parser.errs, err)
		this.consume(untilLc)
	}
	this.Next(lfNotToken) //skip {
	this.parseStatementList(&statementIf.Block, false)
	if this.parser.token.Type != lex.TokenRc {
		this.parser.errs = append(this.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
			this.parser.errMsgPrefix(), this.parser.token.Description))
		this.consume(untilRc)
	}
	this.Next(lfIsToken) // skip }
	if this.parser.token.Type == lex.TokenLf {
		pos := this.parser.mkPos()
		this.Next(lfNotToken)
		if this.parser.token.Type == lex.TokenElseif ||
			this.parser.token.Type == lex.TokenElse {
			this.parser.errs = append(this.parser.errs, fmt.Errorf("%s unexpected new line",
				this.parser.errMsgPrefix(pos)))
		}
	}
	if this.parser.token.Type == lex.TokenElseif {
		statementIf.ElseIfList, err = this.parseElseIfList()
		if err != nil {
			return statementIf, err
		}
	}
	if this.parser.token.Type == lex.TokenLf {
		pos := this.parser.mkPos()
		this.Next(lfNotToken)
		if this.parser.token.Type == lex.TokenElse {
			this.parser.errs = append(this.parser.errs, fmt.Errorf("%s unexpected new line",
				this.parser.errMsgPrefix(pos)))
		}
	}
	if this.parser.token.Type == lex.TokenElse {
		this.Next(lfNotToken)
		if this.parser.token.Type != lex.TokenLc {
			err = fmt.Errorf("%s missing '{' after else", this.parser.errMsgPrefix())
			this.parser.errs = append(this.parser.errs, err)
			this.consume(untilLc)
		}
		this.Next(lfNotToken) // skip {
		statementIf.Else = &ast.Block{}
		this.parseStatementList(statementIf.Else, false)
		if this.parser.token.Type != lex.TokenRc {
			err = fmt.Errorf("%s expect '}', but '%s'",
				this.parser.errMsgPrefix(), this.parser.token.Description)
			this.parser.errs = append(this.parser.errs, err)
			this.consume(untilRc)
		}
		this.Next(lfNotToken) // skip }
	}
	return statementIf, err
}

func (this *BlockParser) parseElseIfList() (elseIfList []*ast.StatementElseIf, err error) {
	elseIfList = []*ast.StatementElseIf{}
	var condition *ast.Expression
	for this.parser.token.Type == lex.TokenElseif {
		this.Next(lfIsToken) // skip elseif token
		this.parser.unExpectNewLineAndSkip()
		condition, err = this.parser.ExpressionParser.parseExpression(false)
		if err != nil {
			this.consume(untilLc)
		}
		if this.parser.token.Type != lex.TokenLc {
			err = fmt.Errorf("%s not '{' after a expression,but '%s'",
				this.parser.errMsgPrefix(), this.parser.token.Description)
			this.parser.errs = append(this.parser.errs)
			this.consume(untilLc)
		}
		this.Next(lfNotToken) // skip {
		block := &ast.Block{}
		this.parseStatementList(block, false)
		elseIfList = append(elseIfList, &ast.StatementElseIf{
			Condition: condition,
			Block:     block,
		})
		if this.parser.token.Type != lex.TokenRc {
			err = fmt.Errorf("%s expect '}', but '%s'",
				this.parser.errMsgPrefix(), this.parser.token.Description)
			this.parser.errs = append(this.parser.errs)
			this.consume(untilRc)
		}
		this.Next(lfIsToken) // skip }
	}
	return elseIfList, err
}

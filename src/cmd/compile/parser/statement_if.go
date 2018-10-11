package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (blockParser *BlockParser) parseIf() (statementIf *ast.StatementIf, err error) {
	blockParser.Next(lfIsToken) // skip if
	var condition *ast.Expression
	blockParser.parser.unExpectNewLineAndSkip()
	condition, err = blockParser.parser.ExpressionParser.parseExpression(true)
	if err != nil {
		blockParser.consume(untilLc)
		blockParser.Next(lfNotToken)
	}
	statementIf = &ast.StatementIf{}
	statementIf.Condition = condition
	blockParser.parser.ifTokenIsLfThenSkip()
	for blockParser.parser.token.Type == lex.TokenSemicolon {
		if statementIf.Condition != nil {
			statementIf.PrefixExpressions = append(statementIf.PrefixExpressions, statementIf.Condition)
		}
		blockParser.Next(lfNotToken) // skip ;
		statementIf.Condition, err = blockParser.parser.ExpressionParser.parseExpression(true)
		if err != nil {
			blockParser.consume(untilLc)
			blockParser.Next(lfNotToken)
		}
	}
	blockParser.parser.ifTokenIsLfThenSkip()
	if blockParser.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s missing '{' after condtion,but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		blockParser.consume(untilLc)
	}
	blockParser.Next(lfNotToken) //skip {
	blockParser.parseStatementList(&statementIf.TrueBlock, false)
	if blockParser.parser.token.Type != lex.TokenRc {
		blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description))
		blockParser.consume(untilRc)
	}
	blockParser.Next(lfIsToken) // skip }
	if blockParser.parser.token.Type == lex.TokenLf {
		pos := blockParser.parser.mkPos()
		blockParser.Next(lfNotToken)
		if blockParser.parser.token.Type == lex.TokenElseif ||
			blockParser.parser.token.Type == lex.TokenElse {
			blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s unexpected new line",
				blockParser.parser.errorMsgPrefix(pos)))
		}
	}
	if blockParser.parser.token.Type == lex.TokenElseif {
		statementIf.ElseIfList, err = blockParser.parseElseIfList()
		if err != nil {
			return statementIf, err
		}
	}
	if blockParser.parser.token.Type == lex.TokenLf {
		pos := blockParser.parser.mkPos()
		blockParser.Next(lfNotToken)
		if blockParser.parser.token.Type == lex.TokenElse {
			blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s unexpected new line",
				blockParser.parser.errorMsgPrefix(pos)))
		}
	}
	if blockParser.parser.token.Type == lex.TokenElse {
		blockParser.Next(lfNotToken)
		if blockParser.parser.token.Type != lex.TokenLc {
			err = fmt.Errorf("%s missing '{' after else", blockParser.parser.errorMsgPrefix())
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			blockParser.consume(untilLc)
		}
		blockParser.Next(lfNotToken) // skip {
		statementIf.ElseBlock = &ast.Block{}
		blockParser.parseStatementList(statementIf.ElseBlock, false)
		if blockParser.parser.token.Type != lex.TokenRc {
			err = fmt.Errorf("%s expect '}', but '%s'",
				blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			blockParser.consume(untilRc)
		}
		blockParser.Next(lfNotToken) // skip }
	}
	return statementIf, err
}

func (blockParser *BlockParser) parseElseIfList() (elseIfList []*ast.StatementElseIf, err error) {
	elseIfList = []*ast.StatementElseIf{}
	var condition *ast.Expression
	for blockParser.parser.token.Type == lex.TokenElseif {
		blockParser.Next(lfIsToken) // skip elseif token
		blockParser.parser.unExpectNewLineAndSkip()
		condition, err = blockParser.parser.ExpressionParser.parseExpression(false)
		if err != nil {
			blockParser.consume(untilLc)
		}
		if blockParser.parser.token.Type != lex.TokenLc {
			err = fmt.Errorf("%s not '{' after a expression,but '%s'",
				blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
			blockParser.parser.errs = append(blockParser.parser.errs)
			blockParser.consume(untilLc)
		}
		blockParser.Next(lfNotToken) // skip {
		block := &ast.Block{}
		blockParser.parseStatementList(block, false)
		elseIfList = append(elseIfList, &ast.StatementElseIf{
			Condition: condition,
			Block:     block,
		})
		if blockParser.parser.token.Type != lex.TokenRc {
			err = fmt.Errorf("%s expect '}', but '%s'",
				blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
			blockParser.parser.errs = append(blockParser.parser.errs)
			blockParser.consume(untilRc)
		}
		blockParser.Next(lfIsToken) // skip }
	}
	return elseIfList, err
}

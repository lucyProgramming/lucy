package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (bp *BlockParser) parseIf() (statementIf *ast.StatementIf, err error) {
	statementIf = &ast.StatementIf{
		Pos: bp.parser.mkPos(),
	}
	bp.Next(lfIsToken) // skip if
	var condition *ast.Expression
	bp.parser.unExpectNewLineAndSkip()
	condition, err = bp.parser.ExpressionParser.parseExpression(true)
	if err != nil {
		bp.consume(untilLc)
		bp.Next(lfNotToken)
	}
	statementIf.Condition = condition
	bp.parser.ifTokenIsLfThenSkip()
	for bp.parser.token.Type == lex.TokenSemicolon {
		if statementIf.Condition != nil {
			statementIf.PrefixExpressions = append(statementIf.PrefixExpressions, statementIf.Condition)
		}
		bp.Next(lfNotToken) // skip ;
		statementIf.Condition, err = bp.parser.ExpressionParser.parseExpression(true)
		if err != nil {
			bp.consume(untilLc)
			bp.Next(lfNotToken)
		}
	}
	bp.parser.ifTokenIsLfThenSkip()
	if bp.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s missing '{' after condtion,but '%s'",
			bp.parser.errMsgPrefix(), bp.parser.token.Description)
		bp.parser.errs = append(bp.parser.errs, err)
		bp.consume(untilLc)
	}
	bp.Next(lfNotToken) //skip {
	bp.parseStatementList(&statementIf.Block, false)
	if bp.parser.token.Type != lex.TokenRc {
		bp.parser.errs = append(bp.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
			bp.parser.errMsgPrefix(), bp.parser.token.Description))
		bp.consume(untilRc)
	}
	bp.Next(lfIsToken) // skip }
	if bp.parser.token.Type == lex.TokenLf {
		pos := bp.parser.mkPos()
		bp.Next(lfNotToken)
		if bp.parser.token.Type == lex.TokenElseif ||
			bp.parser.token.Type == lex.TokenElse {
			bp.parser.errs = append(bp.parser.errs, fmt.Errorf("%s unexpected new line",
				bp.parser.errMsgPrefix(pos)))
		}
	}
	if bp.parser.token.Type == lex.TokenElseif {
		statementIf.ElseIfList, err = bp.parseElseIfList()
		if err != nil {
			return statementIf, err
		}
	}
	if bp.parser.token.Type == lex.TokenLf {
		pos := bp.parser.mkPos()
		bp.Next(lfNotToken)
		if bp.parser.token.Type == lex.TokenElse {
			bp.parser.errs = append(bp.parser.errs, fmt.Errorf("%s unexpected new line",
				bp.parser.errMsgPrefix(pos)))
		}
	}
	if bp.parser.token.Type == lex.TokenElse {
		bp.Next(lfNotToken)
		if bp.parser.token.Type != lex.TokenLc {
			err = fmt.Errorf("%s missing '{' after else", bp.parser.errMsgPrefix())
			bp.parser.errs = append(bp.parser.errs, err)
			bp.consume(untilLc)
		}
		bp.Next(lfNotToken) // skip {
		statementIf.Else = &ast.Block{}
		bp.parseStatementList(statementIf.Else, false)
		if bp.parser.token.Type != lex.TokenRc {
			err = fmt.Errorf("%s expect '}', but '%s'",
				bp.parser.errMsgPrefix(), bp.parser.token.Description)
			bp.parser.errs = append(bp.parser.errs, err)
			bp.consume(untilRc)
		}
		bp.Next(lfNotToken) // skip }
	}
	return statementIf, err
}

func (bp *BlockParser) parseElseIfList() (elseIfList []*ast.StatementElseIf, err error) {
	elseIfList = []*ast.StatementElseIf{}
	var condition *ast.Expression
	for bp.parser.token.Type == lex.TokenElseif {
		bp.Next(lfIsToken) // skip elseif token
		bp.parser.unExpectNewLineAndSkip()
		condition, err = bp.parser.ExpressionParser.parseExpression(false)
		if err != nil {
			bp.consume(untilLc)
		}
		if bp.parser.token.Type != lex.TokenLc {
			err = fmt.Errorf("%s not '{' after a expression,but '%s'",
				bp.parser.errMsgPrefix(), bp.parser.token.Description)
			bp.parser.errs = append(bp.parser.errs)
			bp.consume(untilLc)
		}
		bp.Next(lfNotToken) // skip {
		block := &ast.Block{}
		bp.parseStatementList(block, false)
		elseIfList = append(elseIfList, &ast.StatementElseIf{
			Condition: condition,
			Block:     block,
		})
		if bp.parser.token.Type != lex.TokenRc {
			err = fmt.Errorf("%s expect '}', but '%s'",
				bp.parser.errMsgPrefix(), bp.parser.token.Description)
			bp.parser.errs = append(bp.parser.errs)
			bp.consume(untilRc)
		}
		bp.Next(lfIsToken) // skip }
	}
	return elseIfList, err
}

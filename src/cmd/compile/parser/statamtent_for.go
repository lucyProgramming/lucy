package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (bp *BlockParser) parseFor() (statementFor *ast.StatementFor, err error) {
	statementFor = &ast.StatementFor{}
	statementFor.Pos = bp.parser.mkPos()
	statementFor.Block = &ast.Block{}
	bp.Next(lfIsToken) // skip for
	bp.parser.unExpectNewLineAndSkip()
	if bp.parser.token.Type != lex.TokenLc &&
		bp.parser.token.Type != lex.TokenSemicolon { // not '{' and not ';'
		statementFor.Condition, err = bp.parser.ExpressionParser.parseExpression(true)
		if err != nil {
			bp.consume(untilLc)
			goto parseBlock
		}
	}
	if bp.parser.token.Type == lex.TokenSemicolon {
		bp.Next(lfNotToken) // skip ;
		statementFor.Init = statementFor.Condition
		statementFor.Condition = nil // mk nil
		//condition
		var err error
		if bp.parser.token.Type != lex.TokenSemicolon {
			statementFor.Condition, err = bp.parser.ExpressionParser.parseExpression(false)
			if err != nil {
				bp.consume(untilLc)
				goto parseBlock
			}
			if bp.parser.token.Type != lex.TokenSemicolon {
				bp.parser.errs = append(bp.parser.errs, fmt.Errorf("%s missing semicolon after expression",
					bp.parser.errMsgPrefix()))
				bp.consume(untilLc)
				goto parseBlock
			}
		}
		bp.Next(lfIsToken)
		bp.parser.unExpectNewLineAndSkip()
		if bp.parser.token.Type != lex.TokenLc {
			statementFor.Increment, err = bp.parser.ExpressionParser.parseExpression(true)
			if err != nil {
				bp.consume(untilLc)
				goto parseBlock
			}
		}
	}
parseBlock:
	bp.parser.ifTokenIsLfThenSkip()
	if bp.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'",
			bp.parser.errMsgPrefix(), bp.parser.token.Description)
		bp.parser.errs = append(bp.parser.errs, err)
		return
	}
	bp.Next(lfNotToken) // skip {
	bp.parseStatementList(statementFor.Block, false)
	if bp.parser.token.Type != lex.TokenRc {
		bp.parser.errs = append(bp.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
			bp.parser.errMsgPrefix(), bp.parser.token.Description))
		bp.consume(untilRc)
	}
	bp.Next(lfNotToken) // skip }
	return statementFor, nil
}

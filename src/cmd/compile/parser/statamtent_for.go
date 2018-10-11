package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (blockParser *BlockParser) parseFor() (statementFor *ast.StatementFor, err error) {
	statementFor = &ast.StatementFor{}
	statementFor.Pos = blockParser.parser.mkPos()
	statementFor.Block = &ast.Block{}
	blockParser.Next(lfIsToken) // skip for
	blockParser.parser.unExpectNewLineAndSkip()
	if blockParser.parser.token.Type != lex.TokenLc &&
		blockParser.parser.token.Type != lex.TokenSemicolon { // not '{' and not ';'
		statementFor.Condition, err = blockParser.parser.ExpressionParser.parseExpression(true)
		if err != nil {
			blockParser.consume(untilLc)
			goto parseBlock
		}
	}
	if blockParser.parser.token.Type == lex.TokenSemicolon {
		blockParser.Next(lfNotToken) // skip ;
		statementFor.Init = statementFor.Condition
		statementFor.Condition = nil // mk nil
		//condition
		var err error
		if blockParser.parser.token.Type != lex.TokenSemicolon {
			statementFor.Condition, err = blockParser.parser.ExpressionParser.parseExpression(false)
			if err != nil {
				blockParser.consume(untilLc)
				goto parseBlock
			}
			if blockParser.parser.token.Type != lex.TokenSemicolon {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s missing semicolon after expression",
					blockParser.parser.errorMsgPrefix()))
				blockParser.consume(untilLc)
				goto parseBlock
			}
		}
		blockParser.Next(lfIsToken)
		blockParser.parser.unExpectNewLineAndSkip()
		if blockParser.parser.token.Type != lex.TokenLc {
			statementFor.Increment, err = blockParser.parser.ExpressionParser.parseExpression(true)
			if err != nil {
				blockParser.consume(untilLc)
				goto parseBlock
			}
		}
	}
parseBlock:
	blockParser.parser.ifTokenIsLfThenSkip()
	if blockParser.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return
	}
	blockParser.Next(lfNotToken) // skip {
	blockParser.parseStatementList(statementFor.Block, false)
	if blockParser.parser.token.Type != lex.TokenRc {
		blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description))
		blockParser.consume(untilRc)
	}
	blockParser.Next(lfNotToken) // skip }
	return statementFor, nil
}

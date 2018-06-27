package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (blockParser *BlockParser) parseFor() (f *ast.StatementFor, err error) {
	f = &ast.StatementFor{}
	f.Pos = blockParser.parser.mkPos()
	f.Block = &ast.Block{}
	blockParser.Next()                                                                                       // skip for
	if blockParser.parser.token.Type != lex.TokenLc && blockParser.parser.token.Type != lex.TokenSemicolon { // not {
		e, err := blockParser.parser.ExpressionParser.parseExpression(true)
		if err != nil {
			blockParser.parser.errs = append(blockParser.parser.errs, err)
		} else {
			f.Condition = e
		}
	}
	if blockParser.parser.token.Type == lex.TokenSemicolon {
		blockParser.Next() // skip ;
		f.Init = f.Condition
		f.Condition = nil // mk nil
		//condition
		if blockParser.parser.token.Type != lex.TokenSemicolon {
			e, err := blockParser.parser.ExpressionParser.parseExpression(false)
			if err != nil {
				blockParser.parser.errs = append(blockParser.parser.errs, err)
				blockParser.consume(untilSemicolon)
			} else {
				f.Condition = e
			}
			if blockParser.parser.token.Type != lex.TokenSemicolon {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s missing semicolon after expression",
					blockParser.parser.errorMsgPrefix()))
				blockParser.consume(untilLc)
			}
		}
		blockParser.Next()
		if blockParser.parser.token.Type != lex.TokenLc {
			e, err := blockParser.parser.ExpressionParser.parseExpression(true)
			if err != nil {
				blockParser.parser.errs = append(blockParser.parser.errs, err)
			}
			f.After = e
		}

	}
	if blockParser.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return
	}
	blockParser.Next() // skip {
	blockParser.parseStatementList(f.Block, false)
	if blockParser.parser.token.Type != lex.TokenRc {
		blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description))
		blockParser.consume(untilRc)
	}
	blockParser.Next() // }
	return f, nil
}

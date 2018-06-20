package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (b *BlockParser) parseFor() (f *ast.StatementFor, err error) {
	f = &ast.StatementFor{}
	f.Pos = b.parser.mkPos()
	f.Block = &ast.Block{}
	b.Next()                                                                               // skip for
	if b.parser.token.Type != lex.TOKEN_LC && b.parser.token.Type != lex.TOKEN_SEMICOLON { // not {
		e, err := b.parser.ExpressionParser.parseExpression(true)
		if err != nil {
			b.parser.errs = append(b.parser.errs, err)
		} else {
			f.Condition = e
		}
	}
	if b.parser.token.Type == lex.TOKEN_SEMICOLON {
		b.Next() // skip ;
		f.Init = f.Condition
		f.Condition = nil // mk nil
		//condition
		if b.parser.token.Type != lex.TOKEN_SEMICOLON {
			e, err := b.parser.ExpressionParser.parseExpression(false)
			if err != nil {
				b.parser.errs = append(b.parser.errs, err)
				b.consume(untilSemicolon)
			} else {
				f.Condition = e
			}
			if b.parser.token.Type != lex.TOKEN_SEMICOLON {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s missing semicolon after expression",
					b.parser.errorMsgPrefix()))
				b.consume(untilLc)
			}
		}
		b.Next()
		if b.parser.token.Type != lex.TOKEN_LC {
			e, err := b.parser.ExpressionParser.parseExpression(true)
			if err != nil {
				b.parser.errs = append(b.parser.errs, err)
			}
			f.After = e
		}

	}
	if b.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s expect '{',but '%s'",
			b.parser.errorMsgPrefix(), b.parser.token.Description)
		b.parser.errs = append(b.parser.errs, err)
		return
	}
	b.Next() // skip {
	b.parseStatementList(f.Block, false)
	if b.parser.token.Type != lex.TOKEN_RC {
		b.parser.errs = append(b.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
			b.parser.errorMsgPrefix(), b.parser.token.Description))
		b.consume(untilRc)
	}
	b.Next() // }
	return f, nil
}

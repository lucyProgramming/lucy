package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (b *Block) parseFor() (f *ast.StatementFor, err error) {
	f = &ast.StatementFor{}
	f.Pos = b.parser.mkPos()
	f.Block = &ast.Block{}
	b.Next()                                                                               // skip for
	if b.parser.token.Type != lex.TOKEN_LC && b.parser.token.Type != lex.TOKEN_SEMICOLON { // not {
		e, err := b.parser.Expression.parseExpression(true)
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
			e, err := b.parser.Expression.parseExpression(false)
			if err != nil {
				b.parser.errs = append(b.parser.errs, err)
				b.consume(untils_semicolon)
			} else {
				f.Condition = e
			}
			if b.parser.token.Type != lex.TOKEN_SEMICOLON {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s missing semicolon after expression", b.parser.errorMsgPrefix()))
				b.consume(untils_lc)
			}
		}
		b.Next()
		if b.parser.token.Type != lex.TOKEN_LC {
			e, err := b.parser.Expression.parseExpression(true)
			if err != nil {
				b.parser.errs = append(b.parser.errs, err)
			}
			f.Post = e
		}

	}
	if b.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s expect '{',but '%s'",
			b.parser.errorMsgPrefix(), b.parser.token.Desp)
		b.parser.errs = append(b.parser.errs, err)
		return
	}
	b.Next() // skip {
	b.parseStatementList(f.Block, false)
	if b.parser.token.Type != lex.TOKEN_RC {
		b.parser.errs = append(b.parser.errs, fmt.Errorf("%s expect '}', but '%s'"))
		b.consume(untils_rc)
	}
	b.Next() // }
	return f, nil
}

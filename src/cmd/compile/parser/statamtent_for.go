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
		err = fmt.Errorf("%s not { after for", b.parser.errorMsgPrefix())
		b.parser.errs = append(b.parser.errs, err)
		return
	}
	b.Next()
	err = b.parse(f.Block, false, lex.TOKEN_RC)
	if err != nil {
		return nil, err
	}
	return f, nil
}

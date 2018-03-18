package parser

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

func (b *Block) parseIf() (i *ast.StatementIF, err error) {
	b.Next() // skip if
	if b.parser.eof {
		err = b.parser.mkUnexpectedEofErr()
		b.parser.errs = append(b.parser.errs, err)
		return nil, err
	}
	var e *ast.Expression
	e, err = b.parser.ExpressionParser.parseExpression()
	if err != nil {
		b.parser.errs = append(b.parser.errs, err)
		b.consume(untils_lc)
		b.Next()
	}
	if b.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s missing '{' after a expression,but '%s'", b.parser.errorMsgPrefix(), b.parser.token.Desp)
		b.parser.errs = append(b.parser.errs)
		b.consume(untils_lc) // consume and next
		b.Next()
	}
	i = &ast.StatementIF{}
	i.Condition = e
	i.Block = &ast.Block{}
	b.Next()
	err = b.parse(i.Block, false, lex.TOKEN_RC)
	if b.parser.token.Type == lex.TOKEN_ELSEIF {
		es, err := b.parseElseIfList()
		if err != nil {
			return i, err
		}
		i.ElseIfList = es
	}
	if b.parser.token.Type == lex.TOKEN_ELSE {
		b.Next()
		if b.parser.token.Type != lex.TOKEN_LC {
			err = fmt.Errorf("%s missing { after else", b.parser.errorMsgPrefix())
			return i, err
		}
		i.ElseBlock = &ast.Block{}
		b.Next()
		err = b.parse(i.ElseBlock, false, lex.TOKEN_RC)
	}
	return i, err
}

func (b *Block) parseElseIfList() (es []*ast.StatementElseIf, err error) {
	es = []*ast.StatementElseIf{}
	var e *ast.Expression
	for (b.parser.token.Type == lex.TOKEN_ELSEIF) && !b.parser.eof {
		b.Next() // skip elseif token
		e, err = b.parser.ExpressionParser.parseExpression()
		if err != nil {
			b.parser.errs = append(b.parser.errs, err)
			return es, err
		}
		if b.parser.token.Type != lex.TOKEN_LC {
			err = fmt.Errorf("%s not { after a expression,but %s", b.parser.errorMsgPrefix(), b.parser.token.Desp)
			b.parser.errs = append(b.parser.errs)
			return es, err
		}
		block := &ast.Block{}
		b.Next()
		err = b.parse(block, false, lex.TOKEN_RC)
		if err != nil {
			b.consume(untils_rc)
			b.Next()
			continue
		}
		es = append(es, &ast.StatementElseIf{
			Condition: e,
			Block:     block,
		})
	}
	return es, err
}

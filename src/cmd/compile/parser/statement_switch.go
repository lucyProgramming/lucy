package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (b *Block) parseSwitch() (*ast.StatementSwitch, error) {
	pos := b.parser.mkPos()
	b.Next() // skip switch key word
	condition, err := b.parser.ExpressionParser.parseExpression()
	if err != nil {
		b.parser.errs = append(b.parser.errs, err)
		return nil, err
	}
	if b.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s expect '{',but '%s'", b.parser.errorMsgPrefix(), b.parser.token.Desp)
		b.parser.errs = append(b.parser.errs, err)
		return nil, err
	}
	b.Next() // skip {  , must be case
	if b.parser.token.Type != lex.TOKEN_CASE {
		err = fmt.Errorf("%s expect 'case',but '%s'", b.parser.errorMsgPrefix(), b.parser.token.Desp)
		b.parser.errs = append(b.parser.errs, err)
		return nil, err
	}
	s := &ast.StatementSwitch{}
	s.Pos = pos
	s.Condition = condition
	for b.parser.eof == false && b.parser.token.Type == lex.TOKEN_CASE {
		b.Next() // skip case
		es, err := b.parser.ExpressionParser.parseExpressions()
		if err != nil {
			b.parser.errs = append(b.parser.errs, err)
			return s, err
		}
		if b.parser.token.Type != lex.TOKEN_COLON {
			err = fmt.Errorf("%s expect ':',but '%s'", b.parser.errorMsgPrefix(), b.parser.token.Desp)
			b.parser.errs = append(b.parser.errs, err)
			return s, err
		}
		b.Next() // skip :
		var block *ast.Block
		if b.parser.token.Type != lex.TOKEN_CASE && b.parser.token.Type != lex.TOKEN_DEFAULT {
			block = &ast.Block{}
			err = b.parse(block, true, lex.TOKEN_CASE, lex.TOKEN_RC, lex.TOKEN_DEFAULT)
			if err != nil {
				return s, nil
			}
		}
		s.StatmentSwitchCases = append(s.StatmentSwitchCases, &ast.StatmentSwitchCase{
			Matches: es,
			Block:   block,
		})
	}
	//default value
	if b.parser.token.Type == lex.TOKEN_DEFAULT {
		b.Next() // skip default key word
		if b.parser.token.Type != lex.TOKEN_COLON {
			err = fmt.Errorf("%s missing clon after default", b.parser.errorMsgPrefix())
			b.parser.errs = append(b.parser.errs, err)
		} else {
			b.Next()
		}
		block := ast.Block{}
		err = b.parse(&block, true, lex.TOKEN_RC)
		if err != nil {
			return s, err
		}
		s.Default = &block
	}
	if b.parser.token.Type != lex.TOKEN_RC {
		err = fmt.Errorf("%s expect '}',but '%s'", b.parser.errorMsgPrefix(), b.parser.token.Desp)
		b.parser.errs = append(b.parser.errs, err)
		return s, err
	}
	b.Next() //  skip }
	return s, nil
}

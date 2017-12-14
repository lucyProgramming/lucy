package parser

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

type Block struct {
	parser *Parser
}

func (b *Block) Next() {
	b.parser.Next()
}

func (b *Block) consume(c map[int]bool) {
	b.parser.consume(c)
}

func (b *Block) parse(block *ast.Block) (err error) {
	b.Next() // skip {
	block.Statements = []*ast.Statement{}
	for !b.parser.eof {
		switch b.parser.token.Type {
		case lex.TOKEN_SEMICOLON:
			b.Next() // look up next
			continue
		case lex.TOKEN_RC: // end
			b.Next()
			return
		case lex.TOKEN_IDENTIFIER:
			e, err := b.parser.ExpressionParser.parseExpression()
			if err != nil {
				b.parser.errs = append(b.parser.errs, err)
				b.parser.consume(untils_semicolon)
				b.Next()
				continue
			}
			if b.parser.token.Type != lex.TOKEN_SEMICOLON {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s missing semicolon", b.parser.errorMsgPrefix(e.Pos)))
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:        ast.STATEMENT_TYPE_EXPRESSION,
				Expression: e,
			})
			continue
		case lex.TOKEN_VAR:
			vs, err := b.parser.parseVarDefinition()
			if err != nil {
				b.consume(untils_semicolon)
				b.Next()
				continue
			}
			s := &ast.Statement{
				Typ: ast.STATEMENT_TYPE_EXPRESSION,
				Expression: &ast.Expression{
					Typ:  ast.EXPRESSION_TYPE_VAR,
					Data: &ast.ExpressionDeclareVariable{vs},
				},
			}
			block.Statements = append(block.Statements, s)
		case lex.TOKEN_IF:
			i, err := b.parseIf()
			if err != nil {
				b.consume(untils_rc)
				b.Next()
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:         ast.STATEMENT_TYPE_IF,
				StatementIf: i,
			})
		case lex.TOKEN_FOR:
			f, err := b.parseFor()
			if err != nil {
				b.consume(untils_rc)
				b.Next()
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:          ast.STATEMENT_TYPE_FOR,
				StatementFor: f,
			})
		case lex.TOKEN_SWITCH:
			s, err := b.parseSwitch()
			if err != nil {
				b.consume(untils_rc)
				b.Next()
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:             ast.STATEMENT_TYPE_SWITCH,
				StatementSwitch: s,
			})
		case lex.TOKEN_CONST:
			b.Next()
			if b.parser.token.Type != lex.TOKEN_IDENTIFIER {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s missing identifier after const,but ％s", b.parser.errorMsgPrefix(), b.parser.token.Desp))
				b.consume(untils_semicolon)
				b.Next()
				continue
			}
			vs, typ, err := b.parser.parseAssignedNames()
			if err != nil {
				b.consume(untils_rc_semicolon)
				b.Next()
				continue
			}
			if typ != lex.TOKEN_ASSIGN {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s declare const should use ‘=’ instead of ‘:=’", b.parser.errorMsgPrefix(vs[0].Pos)))
			}
			if b.parser.token.Type != lex.TOKEN_SEMICOLON {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s missing semicolon after const declaration", b.parser.errorMsgPrefix()))
				b.consume(untils_rc_semicolon)
			}
			cs := make([]*ast.Const, len(vs))
			for k, v := range vs {
				c := &ast.Const{}
				c.VariableDefinition = *v
				cs[k] = c
			}
			r := &ast.Statement{}
			r.Typ = ast.STATEMENT_TYPE_EXPRESSION
			r.Expression = &ast.Expression{
				Typ: ast.EXPRESSION_TYPE_CONST,
				Data: &ast.ExpressionDeclareConsts{
					Cs: cs,
				},
				Pos: b.parser.mkPos(),
			}
			block.Statements = append(block.Statements, r)
			b.Next()
		case lex.TOKEN_RETURN:
			b.Next()
			r := &ast.StatementReturn{}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:             ast.STATEMENT_TYPE_RETURN,
				StatementReturn: r,
			})
			if b.parser.token.Type == lex.TOKEN_SEMICOLON {
				b.Next()
				continue
			}
			var es []*ast.Expression
			es, err = b.parser.ExpressionParser.parseExpressions()
			if err != nil {
				b.parser.errs = append(b.parser.errs, err)
				b.consume(untils_semicolon)
				b.Next()
			}
			r.Expressions = es
			if b.parser.token.Type != lex.TOKEN_SEMICOLON {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s  no ‘;’after return statement, but %s", b.parser.errorMsgPrefix(), b.parser.token.Desp))
				continue
			}
			b.Next()
		case lex.TOKEN_LC:
			newblock := &ast.Block{}
			err = b.parse(newblock)
			if err != nil {
				b.consume(untils_rc)
				b.Next()
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:   ast.STATEMENT_TYPE_BLOCK,
				Block: newblock,
			})
		case lex.TOKEN_SKIP:
			b.Next()
			if b.parser.token.Type != lex.TOKEN_SEMICOLON {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s  missing semicolon after skip", b.parser.errorMsgPrefix(), b.parser.token.Desp))
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ: ast.STATEMENT_TYPE_SKIP,
			})
		default:
			b.parser.errs = append(b.parser.errs, fmt.Errorf("%s unkown begining of a statement, but %s", b.parser.errorMsgPrefix(), b.parser.token.Desp))
			b.consume(untils_rc_semicolon)
			b.Next()
		}
	}
	return
}

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
		return nil, err
	}
	if b.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s missing { after a expression,but %s", b.parser.errorMsgPrefix(), b.parser.token.Desp)
		b.parser.errs = append(b.parser.errs)
		b.consume(untils_lc) // consume and next
		b.Next()
	}
	i = &ast.StatementIF{}
	i.Condition = e
	i.Block = &ast.Block{}
	err = b.parse(i.Block)
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
		err = b.parse(i.ElseBlock)
	}
	return i, err
}

func (b *Block) parseElseIfList() (es []*ast.StatementElseIf, err error) {
	es = []*ast.StatementElseIf{}
	var e *ast.Expression
	for (b.parser.token.Type == lex.TOKEN_ELSEIF) && !b.parser.eof {
		b.Next()
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
		err = b.parse(block)
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

func (b *Block) parseFor() (f *ast.StatementFor, err error) {
	f = &ast.StatementFor{}
	f.Pos = b.parser.mkPos()
	f.Block = &ast.Block{}
	b.Next()                                 // skip for
	if b.parser.token.Type != lex.TOKEN_LC { // not {
		e, err := b.parser.ExpressionParser.parseExpression()
		if err != nil {
			b.parser.errs = append(b.parser.errs, err)
		} else {
			f.Condition = e
		}
	}
	if b.parser.token.Type == lex.TOKEN_SEMICOLON {
		b.Next() // skip ;
		e, err := b.parser.ExpressionParser.parseExpression()
		if err != nil {
			b.parser.errs = append(b.parser.errs, err)
			b.consume(untils_semicolon)
		} else {
			f.Init = f.Condition
			f.Condition = e
		}
		if b.parser.token.Type != lex.TOKEN_SEMICOLON {
			b.parser.errs = append(b.parser.errs, fmt.Errorf("%s missing semicolon after expression", b.parser.errorMsgPrefix()))
			b.consume(untils_lc)
		} else {
			b.Next()
			e, err = b.parser.ExpressionParser.parseExpression()
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
	err = b.parse(f.Block)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (b *Block) parseSwitch() (*ast.StatementSwitch, error) {
	return nil, nil
}

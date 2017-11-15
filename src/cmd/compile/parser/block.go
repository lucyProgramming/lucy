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
	block.Statements = []*ast.Statement{}
	for !b.parser.eof {
		switch b.parser.token.Type {
		case lex.TOKEN_RC: // end
			return
		case lex.TOKEN_IDENTIFIER:
			e, err := b.parser.ExpressionParser.parseExpression()
			if err != nil {
				b.parser.errs = append(b.parser.errs, err)
				b.parser.consume(untils_statement)
				b.Next()
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:        ast.STATEMENT_TYPE_EXPRESSION,
				Expression: e,
			})
			b.Next()
		case lex.TOKEN_VAR:
			vs, err := b.parser.parseVarDefinition()
			if err != nil {
				b.consume(untils_statement)
				b.Next()
				continue
			}
			for _, v := range vs {
				block.SymbolicTable.Insert(v.Name, v)
			}
			b.Next()
		case lex.TOKEN_IF:
			i, err := b.parseIf()
			if err != nil {
				b.consume(untils_block)
				b.Next()
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:         ast.STATEMENT_TYPE_IF,
				StatementIf: i,
			})
			b.Next()
		case lex.TOKEN_FOR:
			f, err := b.parseFor()
			if err != nil {
				b.consume(untils_block)
				b.Next()
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:          ast.STATEMENT_TYPE_FOR,
				StatementFor: f,
			})
			b.Next()
		case lex.TOKEN_SWITCH:
			s, err := b.parseSwitch()
			if err != nil {
				b.consume(untils_block)
				b.Next()
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:             ast.STATEMENT_TYPE_SWITCH,
				StatementSwitch: s,
			})
			b.Next()
		case lex.TOKEN_CONST:
			b.Next()
			if b.parser.token.Type != lex.TOKEN_IDENTIFIER {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s not identifier after const,but ％s", b.parser.errorMsgPrefix(), b.parser.token.Desp))
				b.consume(untils_statement)
				b.Next()
				continue
			}
			b.parser.parseAssignedNames()
		case lex.TOKEN_RETURN:
			b.Next()
			var es []*ast.Expression
			r := &ast.StatementReturn{}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:             ast.STATEMENT_TYPE_RETURN,
				StatementReturn: r,
			})
			fmt.Println("##############################", b.parser.token.Desp)
			if b.parser.ExpressionParser.looksLikeAExprssion() {
				es, err = b.parser.ExpressionParser.parseExpressions()
				if err != nil {
					b.parser.errs = append(b.parser.errs, fmt.Errorf("%s not identifier after const,but ％s", b.parser.errorMsgPrefix(), b.parser.token.Desp))
					b.consume(untils_statement)
					b.Next()
				}
				r.Expressions = es
			}
			if b.parser.token.Type != lex.TOKEN_SEMICOLON {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s  no ; after return statement, but %s", b.parser.errorMsgPrefix(), b.parser.token.Desp))
				continue
			}
			b.Next()
		case lex.TOKEN_LC:
			newblock := &ast.Block{}
			err = b.parse(newblock)
			if err != nil {
				b.consume(untils_block)
				b.Next()
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:   ast.STATEMENT_TYPE_BLOCK,
				Block: newblock,
			})
		default:
			b.parser.errs = append(b.parser.errs, fmt.Errorf("%s unkown begining of a statement, but %s", b.parser.errorMsgPrefix()))
			b.consume(untils_statement)
			b.Next()
		}
	}
	return
}

func (b *Block) parseIf() (i *ast.StatementIF, err error) {
	b.Next()
	if b.parser.eof {
		err = b.parser.mkUnexpectedEofErr()
		b.parser.errs = append(b.parser.errs)
		return nil, err
	}
	var e *ast.Expression
	e, err = b.parser.ExpressionParser.parseExpression()
	if err != nil {
		return nil, err
	}
	if b.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s not { after a expression,but %s", b.parser.errorMsgPrefix(), b.parser.token.Desp)
		b.parser.errs = append(b.parser.errs)
		return nil, err
	}
	block := &ast.Block{}
	err = b.parse(block)
	if err != nil {
		return
	}
	i = &ast.StatementIF{}
	i.Condition = e
	i.Block = block
	// skip }
	b.Next()

	if b.parser.token.Type == lex.TOKEN_ELSEIF {
		es, err := b.parseElseIfList()
		if err != nil {
			b.consume(untils_block)
		}
		i.ElseIfList = es
	}

	if b.parser.token.Type == lex.TOKEN_ELSE {

	}

	return nil, nil
}

func (b *Block) parseElseIfList() (es []*ast.StatementElseIf, err error) {
	es = []*ast.StatementElseIf{}
	var e *ast.Expression
	for (b.parser.token.Type == lex.TOKEN_ELSEIF) && !b.parser.eof {
		b.Next()
		e, err = b.parser.ExpressionParser.parseExpression()
		if err != nil {
			return es, err
		}
		if b.parser.token.Type != lex.TOKEN_LC {
			err = fmt.Errorf("%s not { after a expression,but %s", b.parser.errorMsgPrefix(), b.parser.token.Desp)
			b.parser.errs = append(b.parser.errs)
			return es, err
		}
		b.Next()
		block := &ast.Block{}
		err = b.parse(block)
		if err != nil {
			b.consume(untils_block)
			b.Next()
			continue
		}
		es = append(es, &ast.StatementElseIf{
			Condition: e,
			Block:     block,
		})

		b.Next()
	}
	return es, nil
}

func (b *Block) parseFor() (f *ast.StatementFor, err error) {
	b.Next()
	var e *ast.Expression
	e, err = b.parser.ExpressionParser.parseExpression()
	if err != nil {
		b.parser.errs = append(b.parser.errs, err)
		return nil, err
	}

	f = &ast.StatementFor{}
	f.Condition = e
	f.Block = &ast.Block{}
	err = b.parse(f.Block)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (b *Block) parseSwitch() (*ast.StatementSwitch, error) {
	return nil, nil
}

func (b *Block) insertVariableIntoBlock(block *ast.Block, vars []*ast.VariableDefinition) (errs []error) {
	errs = []error{}
	if vars == nil || len(vars) == 0 {
		return
	}
	if block.SymbolicTable.ItemsMap == nil {
		block.SymbolicTable.ItemsMap = make(map[string]*ast.SymbolicItem)
	}
	var err error
	for _, v := range vars {
		if v.Name == "" {
			continue
		}
		err = block.SymbolicTable.Insert(v.Name, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return
}

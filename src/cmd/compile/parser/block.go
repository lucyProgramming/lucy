package parser

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

type Block struct {
	parser *Parser
	token  *lex.Token
}

func (b *Block) Next() {
	b.parser.Next()
	b.token = b.parser.token
}
func (b *Block) consume(t ...int) {
	b.parser.consume(t...)
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
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:        ast.STATEMENT_TYPE_EXPRESSION,
				Expression: e,
			})
		case lex.TOKEN_VAR:
			vs := b.parser.parseVarDefinition()
			for _, v := range vs {
				block.SymbolicTable.Insert(v.Name, v)
			}
		case lex.TOKEN_IF:
		case lex.TOKEN_FOR:
		case lex.TOKEN_SWITCH:
		default:
			b.parser.errs = append(b.parser.errs, fmt.Errorf("%s unkown begining of a statement", b.parser.errorMsgPrefix()))
			b.consume(lex.TOKEN_SEMICOLON)
			b.Next()
		}
	}
	return
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

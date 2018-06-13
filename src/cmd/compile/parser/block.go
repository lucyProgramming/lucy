package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
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

func (b *Block) parseStatementList(block *ast.Block, isGlobal bool) {
	block.Pos = b.parser.mkPos()
	isDefer := false
	reset := func() {
		isDefer = false
	}
	validAfterDefer := func() bool {
		return b.parser.token.Type == lex.TOKEN_IDENTIFIER ||
			b.parser.token.Type == lex.TOKEN_LP ||
			b.parser.token.Type == lex.TOKEN_LC
	}
	var err error
	block.Statements = []*ast.Statement{}
	for lex.TOKEN_EOF != b.parser.token.Type {
		if len(b.parser.errs) > b.parser.nerr {
			block.EndPos = b.parser.mkPos()
			break
		}
		switch b.parser.token.Type {
		case lex.TOKEN_SEMICOLON:
			reset()
			b.Next() // look up next
			continue
		case lex.TOKEN_DEFER:
			isDefer = true
			b.Next()
			if validAfterDefer() == false {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s not a valid token '%s' after defer",
					b.parser.errorMsgPrefix(), b.parser.token.Description))
				reset()
			}
		case lex.TOKEN_IDENTIFIER:
			b.parseExpressionStatement(block, isDefer)
			reset()
		case lex.TOKEN_LP:
			b.parseExpressionStatement(block, isDefer)
			reset()
		case lex.TOKEN_FUNCTION:
			pos := b.parser.mkPos()
			f, err := b.parser.Function.parse(true)
			if err != nil {
				b.parser.consume(untils_rc_semicolon)
			}
			f.AccessFlags |= cg.ACC_METHOD_PRIVATE
			s := &ast.Statement{}
			s.Pos = pos
			s.Typ = ast.STATEMENT_TYPE_EXPRESSION
			s.Expression = &ast.Expression{}
			s.Expression.Typ = ast.EXPRESSION_TYPE_FUNCTION
			s.Expression.Data = f
			block.Statements = append(block.Statements, s)
		case lex.TOKEN_VAR:
			pos := b.parser.mkPos()
			b.Next() // skip var key word
			vs, es, typ, err := b.parser.parseConstDefinition(true)
			if err != nil {
				b.consume(untils_semicolon)
				b.Next()
				continue
			}
			if typ != nil && typ.Type != lex.TOKEN_ASSIGN {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s use '=' to initialize value",
						b.parser.errorMsgPrefix()))
			}
			s := &ast.Statement{
				Typ: ast.STATEMENT_TYPE_EXPRESSION,
				Expression: &ast.Expression{
					Typ:  ast.EXPRESSION_TYPE_VAR,
					Data: &ast.ExpressionDeclareVariable{Variables: vs, Values: es},
					Pos:  pos,
				},
				Pos: pos,
			}
			block.Statements = append(block.Statements, s)
			if isDefer {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s defer mixup with expression var not allow",
						b.parser.errorMsgPrefix()))
			}
			reset()
		case lex.TOKEN_IF:
			pos := b.parser.mkPos()
			i, err := b.parseIf()
			if err != nil {
				b.consume(untils_rc)
				b.Next()
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:         ast.STATEMENT_TYPE_IF,
				StatementIf: i,
				Pos:         pos,
			})
			if isDefer {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s defer mixup with  statment if not allow",
						b.parser.errorMsgPrefix()))
			}
			reset()
		case lex.TOKEN_FOR:
			pos := b.parser.mkPos()
			f, err := b.parseFor()
			if err != nil {
				b.consume(untils_rc)
				b.Next()
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:          ast.STATEMENT_TYPE_FOR,
				StatementFor: f,
				Pos:          pos,
			})
		case lex.TOKEN_SWITCH:
			pos := b.parser.mkPos()
			s, err := b.parseSwitch()
			if err != nil {
				b.consume(untils_rc)
				b.Next()
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:             ast.STATEMENT_TYPE_SWITCH,
				StatementSwitch: s,
				Pos:             pos,
			})
		case lex.TOKEN_CONST:
			if isDefer {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s defer mixup with const definition not allow",
						b.parser.errorMsgPrefix()))
				reset()
			}
			pos := b.parser.mkPos()
			b.Next()
			if b.parser.token.Type != lex.TOKEN_IDENTIFIER {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s missing identifier after const,but '%s'",
						b.parser.errorMsgPrefix(), b.parser.token.Description))
				b.consume(untils_semicolon)
				b.Next()
				continue
			}
			vs, es, typ, err := b.parser.parseConstDefinition(false)
			if err != nil {
				b.consume(untils_rc_semicolon)
				b.Next()
				continue
			}
			if typ != nil && typ.Type != lex.TOKEN_ASSIGN {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s declare const should use ‘=’ instead of ‘:=’",
						b.parser.errorMsgPrefix(vs[0].Pos)))
			}
			if b.parser.token.Type != lex.TOKEN_SEMICOLON {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s missing semicolon after const declaration",
						b.parser.errorMsgPrefix()))
				b.consume(untils_rc_semicolon)
			}
			if len(vs) != len(es) {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s cannot assign '%d' values to '%d' destination",
						b.parser.errorMsgPrefix(vs[0].Pos), len(es), len(vs)))
			}
			cs := make([]*ast.Constant, len(vs))
			for k, v := range vs {
				c := &ast.Constant{}
				c.VariableDefinition = *v
				cs[k] = c
				if k < len(es) {
					cs[k].Expression = es[k] // assignment
				}
			}
			r := &ast.Statement{}
			r.Typ = ast.STATEMENT_TYPE_EXPRESSION
			r.Pos = pos
			r.Expression = &ast.Expression{
				Typ:  ast.EXPRESSION_TYPE_CONST,
				Data: cs,
				Pos:  pos,
			}
			block.Statements = append(block.Statements, r)
			b.Next()
		case lex.TOKEN_RETURN:
			pos := b.parser.mkPos()
			if isDefer {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s defer mixup with statement return not allow",
						b.parser.errorMsgPrefix()))
				reset()
			}
			if isGlobal {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s 'return' cannot used in global block",
						b.parser.errorMsgPrefix()))
			}
			b.Next()
			r := &ast.StatementReturn{}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:             ast.STATEMENT_TYPE_RETURN,
				StatementReturn: r,
				Pos:             pos,
			})
			if b.parser.token.Type == lex.TOKEN_SEMICOLON {
				b.Next()
				continue
			}
			var es []*ast.Expression
			es, err = b.parser.Expression.parseExpressions()
			if err != nil {
				b.parser.errs = append(b.parser.errs, err)
				b.consume(untils_semicolon)
				b.Next()
			}
			r.Expressions = es
			if b.parser.token.Type != lex.TOKEN_SEMICOLON {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s  no semicolon after return statement, but %s",
						b.parser.errorMsgPrefix(), b.parser.token.Description))
				continue
			}
			b.Next()
		case lex.TOKEN_LC:
			pos := b.parser.mkPos()
			newblock := ast.Block{}
			b.Next() // skip {
			b.parseStatementList(&newblock, false)
			if b.parser.token.Type != lex.TOKEN_RC {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
					b.parser.errorMsgPrefix(), b.parser.token.Description))
				b.consume(untils_rc)
			}
			b.Next()
			if isDefer {
				d := &ast.Defer{
					Block: newblock,
				}
				block.Statements = append(block.Statements, &ast.Statement{
					Typ:   ast.STATEMENT_TYPE_DEFER,
					Defer: d,
					Pos:   pos,
				})
			} else {
				block.Statements = append(block.Statements, &ast.Statement{
					Typ:   ast.STATEMENT_TYPE_BLOCK,
					Block: &newblock,
					Pos:   pos,
				})
			}
			reset()
		case lex.TOKEN_PASS:
			pos := b.parser.mkPos()
			if isDefer {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s defer mixup with statement not allow",
						b.parser.errorMsgPrefix()))
				reset()
			}
			if isGlobal == false {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s 'pass' can only be used in global blocks",
						b.parser.errorMsgPrefix()))
			}
			b.Next()
			if b.parser.token.Type != lex.TOKEN_SEMICOLON {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s  missing semicolon after 'skip'",
					b.parser.errorMsgPrefix()))
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:             ast.STATEMENT_TYPE_RETURN,
				Pos:             pos,
				StatementReturn: &ast.StatementReturn{},
			})
		case lex.TOKEN_CONTINUE:
			pos := b.parser.mkPos()
			if isDefer {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s defer mixup with statement not allow",
					b.parser.errorMsgPrefix()))
				reset()
			}
			b.Next()
			if b.parser.token.Type != lex.TOKEN_SEMICOLON {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s  missing semicolon after 'continue'",
					b.parser.errorMsgPrefix()))
			} else {
				b.Next()
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:               ast.STATEMENT_TYPE_CONTINUE,
				StatementContinue: &ast.StatementContinue{},
				Pos:               pos,
			})
		case lex.TOKEN_BREAK:
			pos := b.parser.mkPos()
			if isDefer {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s defer mixup with statement 'break' not allow",
					b.parser.errorMsgPrefix()))
				reset()
			}
			b.Next()
			if b.parser.token.Type != lex.TOKEN_SEMICOLON {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s  missing semicolon after 'break'",
					b.parser.errorMsgPrefix()))
			} else {
				b.Next()
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:            ast.STATEMENT_TYPE_BREAK,
				StatementBreak: &ast.StatementBreak{},
				Pos:            pos,
			})
		case lex.TOKEN_GOTO:
			pos := b.parser.mkPos()
			if isDefer {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s defer mixup with statement 'goto' not allow",
					b.parser.errorMsgPrefix()))
				reset()
			}
			b.Next() // skip goto key word
			if b.parser.token.Type != lex.TOKEN_IDENTIFIER {
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s  missing identifier after goto statement, but '%s'",
						b.parser.errorMsgPrefix(), b.parser.token.Description))
				b.consume(untils_semicolon)
				b.Next()
				continue
			}
			s := &ast.StatementGoto{}
			s.Name = b.parser.token.Data.(string)
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:           ast.STATEMENT_TYPE_GOTO,
				StatementGoto: s,
				Pos:           pos,
			})
			b.Next()
			if b.parser.token.Type != lex.TOKEN_SEMICOLON { // incase forget
				b.parser.errs = append(b.parser.errs,
					fmt.Errorf("%s  missing semicolon after goto statement,but '%s'",
						b.parser.errorMsgPrefix(), b.parser.token.Description))
			}
			b.Next()
		case lex.TOKEN_TYPE:
			pos := b.parser.mkPos()
			if isDefer {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s defer mixup with statement 'type' not allow",
					b.parser.errorMsgPrefix()))
				reset()
			}
			alias, err := b.parser.parseTypeaAlias()
			if err != nil {
				b.consume(untils_semicolon)
				b.Next()
				continue
			}
			if b.parser.token.Type != lex.TOKEN_SEMICOLON {
				b.parser.errs = append(b.parser.errs, fmt.Errorf("%s  missing semicolon", b.parser.errorMsgPrefix()))
			}
			s := &ast.Statement{}
			s.Pos = pos
			s.Typ = ast.STATEMENT_TYPE_EXPRESSION
			s.Expression = &ast.Expression{}
			s.Expression.Typ = ast.EXPRESSION_TYPE_TYPE_ALIAS
			s.Expression.Data = alias
			block.Statements = append(block.Statements, s)
			b.Next()
		case lex.TOKEN_CLASS, lex.TOKEN_INTERFACE:
			pos := b.parser.mkPos()
			var class *ast.Class
			var err error
			if b.parser.token.Type == lex.TOKEN_CLASS {
				class, err = b.parser.Class.parse()
			} else {
				class, err = b.parser.Interface.parse()
			}
			if err != nil {
				b.consume(untils_rc)
				b.Next()
				continue
			}
			s := &ast.Statement{}
			s.Pos = pos
			s.Typ = ast.STATEMENT_TYPE_CLASS
			s.Class = class
			block.Statements = append(block.Statements, s)
		case lex.TOKEN_ENUM:
			pos := b.parser.mkPos()
			e, err := b.parser.parseEnum(false)
			if err != nil {
				b.consume(untils_rc)
				b.Next()
				continue
			}
			s := &ast.Statement{}
			s.Pos = pos
			s.Typ = ast.STATEMENT_TYPE_ENUM
			s.Enum = e
			block.Statements = append(block.Statements, s)
		default:
			return
		}
	}
	return
}

func (b *Block) parseExpressionStatement(block *ast.Block, isDefer bool) {
	pos := b.parser.mkPos()
	e, err := b.parser.Expression.parseExpression(true)
	if err != nil {
		b.parser.errs = append(b.parser.errs, err)
		b.parser.consume(untils_semicolon)
		b.Next()
		return
	}
	if e.Typ == ast.EXPRESSION_TYPE_IDENTIFIER && b.parser.token.Type == lex.TOKEN_COLON {
		if isDefer {
			b.parser.errs = append(b.parser.errs, fmt.Errorf("%s defer mixup with statement lable has no meaning",
				b.parser.errorMsgPrefix()))
		}
		b.Next() // skip :
		s := &ast.Statement{}
		s.Pos = pos
		s.Typ = ast.STATEMENT_TYPE_LABLE
		lable := &ast.StatementLabel{}
		s.StatementLabel = lable
		lable.Statement = s
		lable.Name = e.Data.(*ast.ExpressionIdentifier).Name
		block.Statements = append(block.Statements, s)
		lable.Block = block
		block.Insert(lable.Name, e.Pos, lable) // insert first,so this label can be found before it is checked
	} else {
		if b.parser.token.Type != lex.TOKEN_SEMICOLON {
			b.parser.errs = append(b.parser.errs, fmt.Errorf("%s missing semicolon afete a statement expression",
				b.parser.errorMsgPrefix(e.Pos)))
		}
		if isDefer {
			d := &ast.Defer{}
			d.Block.Statements = []*ast.Statement{&ast.Statement{
				Typ:        ast.STATEMENT_TYPE_EXPRESSION,
				Expression: e,
				Pos:        pos,
			}}
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:   ast.STATEMENT_TYPE_DEFER,
				Defer: d,
			})
		} else {
			block.Statements = append(block.Statements, &ast.Statement{
				Typ:        ast.STATEMENT_TYPE_EXPRESSION,
				Expression: e,
				Pos:        pos,
			})
		}
	}
}

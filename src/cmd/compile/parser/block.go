package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type BlockParser struct {
	parser *Parser
}

func (blockParser *BlockParser) Next() {
	blockParser.parser.Next()
}

func (blockParser *BlockParser) consume(c map[int]bool) {
	blockParser.parser.consume(c)
}

func (blockParser *BlockParser) parseStatementList(block *ast.Block, isGlobal bool) {
	block.Pos = blockParser.parser.mkPos()
	isDefer := false
	reset := func() {
		isDefer = false
	}
	validAfterDefer := func() bool {
		return blockParser.parser.token.Type == lex.TOKEN_IDENTIFIER ||
			blockParser.parser.token.Type == lex.TOKEN_LP ||
			blockParser.parser.token.Type == lex.TOKEN_LC
	}
	var err error
	block.Statements = []*ast.Statement{}
	for lex.TOKEN_EOF != blockParser.parser.token.Type {
		if len(blockParser.parser.errs) > blockParser.parser.nErrors2Stop {
			block.EndPos = blockParser.parser.mkPos()
			break
		}
		switch blockParser.parser.token.Type {
		case lex.TOKEN_SEMICOLON:
			reset()
			blockParser.Next() // look up next
			continue
		case lex.TOKEN_DEFER:
			isDefer = true
			blockParser.Next()
			if validAfterDefer() == false {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s not a valid token '%s' after defer",
					blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description))
				reset()
			}
		case lex.TOKEN_IDENTIFIER, lex.TOKEN_LP, lex.TOKEN_FUNCTION:
			blockParser.parseExpressionStatement(block, isDefer)
			reset()
		case lex.TOKEN_VAR:
			pos := blockParser.parser.mkPos()
			blockParser.Next() // skip var key word
			vs, es, typ, err := blockParser.parser.parseConstDefinition(true)
			if err != nil {
				blockParser.consume(untilSemicolon)
				blockParser.Next()
				continue
			}
			if typ != nil && typ.Type != lex.TOKEN_ASSIGN {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s use '=' to initialize value",
						blockParser.parser.errorMsgPrefix()))
			}
			s := &ast.Statement{
				Type: ast.StatementTypeExpression,
				Expression: &ast.Expression{
					Type: ast.EXPRESSION_TYPE_VAR,
					Data: &ast.ExpressionDeclareVariable{Variables: vs, InitValues: es},
					Pos:  pos,
				},
				Pos: pos,
			}
			block.Statements = append(block.Statements, s)
			if isDefer {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s defer mixup with expression var not allow",
						blockParser.parser.errorMsgPrefix()))
			}
			reset()
		case lex.TOKEN_IF:
			pos := blockParser.parser.mkPos()
			i, err := blockParser.parseIf()
			if err != nil {
				blockParser.consume(untilRc)
				blockParser.Next()
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:        ast.StatementTypeIf,
				StatementIf: i,
				Pos:         pos,
			})
			if isDefer {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s defer mixup with  statment if not allow",
						blockParser.parser.errorMsgPrefix()))
			}
			reset()
		case lex.TOKEN_FOR:
			pos := blockParser.parser.mkPos()
			f, err := blockParser.parseFor()
			if err != nil {
				blockParser.consume(untilRc)
				blockParser.Next()
				continue
			}
			f.Block.IsForBlock = true
			block.Statements = append(block.Statements, &ast.Statement{
				Type:         ast.StatementTypeFor,
				StatementFor: f,
				Pos:          pos,
			})
		case lex.TOKEN_SWITCH:
			pos := blockParser.parser.mkPos()
			s, err := blockParser.parseSwitch()
			if err != nil {
				blockParser.consume(untilRc)
				blockParser.Next()
				continue
			}
			if _, ok := s.(*ast.StatementSwitch); ok {
				block.Statements = append(block.Statements, &ast.Statement{
					Type:            ast.StatementTypeSwitch,
					StatementSwitch: s.(*ast.StatementSwitch),
					Pos:             pos,
				})
			} else {
				block.Statements = append(block.Statements, &ast.Statement{
					Type: ast.StatementTypeSwitchTemplate,
					StatementSwitchTemplate: s.(*ast.StatementSwitchTemplate),
					Pos: pos,
				})
			}

		case lex.TOKEN_CONST:
			if isDefer {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s defer mixup with const definition not allow",
						blockParser.parser.errorMsgPrefix()))
				reset()
			}
			pos := blockParser.parser.mkPos()
			blockParser.Next()
			if blockParser.parser.token.Type != lex.TOKEN_IDENTIFIER {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s missing identifier after const,but '%s'",
						blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description))
				blockParser.consume(untilSemicolon)
				blockParser.Next()
				continue
			}
			vs, es, typ, err := blockParser.parser.parseConstDefinition(false)
			if err != nil {
				blockParser.consume(untilRcAndSemicolon)
				blockParser.Next()
				continue
			}
			if typ != nil && typ.Type != lex.TOKEN_ASSIGN {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s declare const should use ‘=’ instead of ‘:=’",
						blockParser.parser.errorMsgPrefix(vs[0].Pos)))
			}
			if blockParser.parser.token.Type != lex.TOKEN_SEMICOLON {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s missing semicolon after const declaration",
						blockParser.parser.errorMsgPrefix()))
				blockParser.consume(untilRcAndSemicolon)
			}
			if len(vs) != len(es) {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s cannot assign '%d' values to '%d' destination",
						blockParser.parser.errorMsgPrefix(vs[0].Pos), len(es), len(vs)))
			}
			cs := make([]*ast.Constant, len(vs))
			for k, v := range vs {
				c := &ast.Constant{}
				c.Variable = *v
				cs[k] = c
				if k < len(es) {
					cs[k].Expression = es[k] // assignment
				}
			}
			r := &ast.Statement{}
			r.Type = ast.StatementTypeExpression
			r.Pos = pos
			r.Expression = &ast.Expression{
				Type: ast.EXPRESSION_TYPE_CONST,
				Data: cs,
				Pos:  pos,
			}
			block.Statements = append(block.Statements, r)
			blockParser.Next()
		case lex.TOKEN_RETURN:
			pos := blockParser.parser.mkPos()
			if isDefer {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s defer mixup with statement return not allow",
						blockParser.parser.errorMsgPrefix()))
				reset()
			}
			if isGlobal {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s 'return' cannot used in global block",
						blockParser.parser.errorMsgPrefix()))
			}
			blockParser.Next()
			r := &ast.StatementReturn{}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:            ast.StatementTypeReturn,
				StatementReturn: r,
				Pos:             pos,
			})
			if blockParser.parser.token.Type == lex.TOKEN_SEMICOLON {
				blockParser.Next()
				continue
			}
			var es []*ast.Expression
			es, err = blockParser.parser.ExpressionParser.parseExpressions()
			if err != nil {
				blockParser.parser.errs = append(blockParser.parser.errs, err)
				blockParser.consume(untilSemicolon)
				blockParser.Next()
			}
			r.Expressions = es
			if blockParser.parser.token.Type != lex.TOKEN_SEMICOLON {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s  no semicolon after return statement, but %s",
						blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description))
				continue
			}
			blockParser.Next()
		case lex.TOKEN_LC:
			pos := blockParser.parser.mkPos()
			newBlock := ast.Block{}
			blockParser.Next() // skip {
			blockParser.parseStatementList(&newBlock, false)
			if blockParser.parser.token.Type != lex.TOKEN_RC {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
					blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description))
				blockParser.consume(untilRc)
			}
			blockParser.Next()
			if isDefer {
				d := &ast.StatementDefer{
					Block: newBlock,
				}
				block.Statements = append(block.Statements, &ast.Statement{
					Type:  ast.StatementTypeDefer,
					Defer: d,
					Pos:   pos,
				})
			} else {
				block.Statements = append(block.Statements, &ast.Statement{
					Type:  ast.StatementTypeBlock,
					Block: &newBlock,
					Pos:   pos,
				})
			}
			reset()
		case lex.TOKEN_PASS:
			pos := blockParser.parser.mkPos()
			if isDefer {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s defer mixup with statement not allow",
						blockParser.parser.errorMsgPrefix()))
				reset()
			}
			if isGlobal == false {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s 'pass' can only be used in global blocks",
						blockParser.parser.errorMsgPrefix()))
			}
			blockParser.Next()
			if blockParser.parser.token.Type != lex.TOKEN_SEMICOLON {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s  missing semicolon after 'skip'",
					blockParser.parser.errorMsgPrefix()))
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:            ast.StatementTypeReturn,
				Pos:             pos,
				StatementReturn: &ast.StatementReturn{},
			})
		case lex.TOKEN_CONTINUE:
			pos := blockParser.parser.mkPos()
			if isDefer {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s defer mixup with statement not allow",
					blockParser.parser.errorMsgPrefix()))
				reset()
			}
			blockParser.Next()
			if blockParser.parser.token.Type != lex.TOKEN_SEMICOLON {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s  missing semicolon after 'continue'",
					blockParser.parser.errorMsgPrefix()))
			} else {
				blockParser.Next()
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:              ast.StatementTypeContinue,
				StatementContinue: &ast.StatementContinue{},
				Pos:               pos,
			})
		case lex.TOKEN_BREAK:
			pos := blockParser.parser.mkPos()
			if isDefer {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s defer mixup with statement 'break' not allow",
					blockParser.parser.errorMsgPrefix()))
				reset()
			}
			blockParser.Next()
			if blockParser.parser.token.Type != lex.TOKEN_SEMICOLON {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s  missing semicolon after 'break'",
					blockParser.parser.errorMsgPrefix()))
			} else {
				blockParser.Next()
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:           ast.StatementTypeBreak,
				StatementBreak: &ast.StatementBreak{},
				Pos:            pos,
			})
		case lex.TOKEN_GOTO:
			pos := blockParser.parser.mkPos()
			if isDefer {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s defer mixup with statement 'goto' not allow",
					blockParser.parser.errorMsgPrefix()))
				reset()
			}
			blockParser.Next() // skip goto key word
			if blockParser.parser.token.Type != lex.TOKEN_IDENTIFIER {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s  missing identifier after goto statement, but '%s'",
						blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description))
				blockParser.consume(untilSemicolon)
				blockParser.Next()
				continue
			}
			s := &ast.StatementGoTo{}
			s.LabelName = blockParser.parser.token.Data.(string)
			block.Statements = append(block.Statements, &ast.Statement{
				Type:          ast.StatementTypeGoto,
				StatementGoTo: s,
				Pos:           pos,
			})
			blockParser.Next()
			if blockParser.parser.token.Type != lex.TOKEN_SEMICOLON { // in case forget
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s  missing semicolon after goto statement,but '%s'",
						blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description))
			}
			blockParser.Next()
		case lex.TOKEN_TYPE:
			pos := blockParser.parser.mkPos()
			if isDefer {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s defer mixup with statement 'type' not allow",
					blockParser.parser.errorMsgPrefix()))
				reset()
			}
			alias, err := blockParser.parser.parseTypeAlias()
			if err != nil {
				blockParser.consume(untilSemicolon)
				blockParser.Next()
				continue
			}
			if blockParser.parser.token.Type != lex.TOKEN_SEMICOLON {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s  missing semicolon", blockParser.parser.errorMsgPrefix()))
			}
			s := &ast.Statement{}
			s.Pos = pos
			s.Type = ast.StatementTypeExpression
			s.Expression = &ast.Expression{}
			s.Expression.Type = ast.EXPRESSION_TYPE_TYPE_ALIAS
			s.Expression.Data = alias
			block.Statements = append(block.Statements, s)
			blockParser.Next()
		case lex.TOKEN_CLASS, lex.TOKEN_INTERFACE:
			pos := blockParser.parser.mkPos()
			var class *ast.Class
			var err error
			if blockParser.parser.token.Type == lex.TOKEN_CLASS {
				class, err = blockParser.parser.ClassParser.parse()
			} else {
				class, err = blockParser.parser.InterfaceParser.parse()
			}
			if err != nil {
				blockParser.consume(untilRc)
				blockParser.Next()
				continue
			}
			s := &ast.Statement{}
			s.Pos = pos
			s.Type = ast.StatementTypeClass
			s.Class = class
			block.Statements = append(block.Statements, s)
		case lex.TOKEN_ENUM:
			pos := blockParser.parser.mkPos()
			e, err := blockParser.parser.parseEnum(false)
			if err != nil {
				blockParser.consume(untilRc)
				blockParser.Next()
				continue
			}
			s := &ast.Statement{}
			s.Pos = pos
			s.Type = ast.StatementTypeEnum
			s.Enum = e
			block.Statements = append(block.Statements, s)
		default:
			return
		}
	}
	return
}

func (blockParser *BlockParser) parseExpressionStatement(block *ast.Block, isDefer bool) {
	pos := blockParser.parser.mkPos()
	e, err := blockParser.parser.ExpressionParser.parseExpression(true)
	if err != nil {
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		blockParser.parser.consume(untilSemicolon)
		blockParser.Next()
		return
	}
	if e.Type == ast.EXPRESSION_TYPE_IDENTIFIER && blockParser.parser.token.Type == lex.TOKEN_COLON {
		if isDefer {
			blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s defer mixup with statement lable has no meaning",
				blockParser.parser.errorMsgPrefix()))
		}
		blockParser.Next() // skip :
		s := &ast.Statement{}
		s.Pos = pos
		s.Type = ast.StatementTypeLabel
		label := &ast.StatementLabel{}
		s.StatementLabel = label
		label.Statement = s
		label.Name = e.Data.(*ast.ExpressionIdentifier).Name
		block.Statements = append(block.Statements, s)
		label.Block = block
		block.Insert(label.Name, e.Pos, label) // insert first,so this label can be found before it is checked
	} else {
		if blockParser.parser.token.Type != lex.TOKEN_SEMICOLON {
			if blockParser.parser.lastToken != nil && blockParser.parser.lastToken.Type == lex.TOKEN_RC {
			} else {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s missing semicolon afete a statement expression",
					blockParser.parser.errorMsgPrefix(e.Pos)))
			}
		}
		if isDefer {
			d := &ast.StatementDefer{}
			d.Block.Statements = []*ast.Statement{&ast.Statement{
				Type:       ast.StatementTypeExpression,
				Expression: e,
				Pos:        pos,
			}}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:  ast.StatementTypeDefer,
				Defer: d,
			})
		} else {
			block.Statements = append(block.Statements, &ast.Statement{
				Type:       ast.StatementTypeExpression,
				Expression: e,
				Pos:        pos,
			})
		}
	}
}

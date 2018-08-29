package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type BlockParser struct {
	parser *Parser
}

func (blockParser *BlockParser) Next(lfIsToken bool) {
	blockParser.parser.Next(lfIsToken)
}

func (blockParser *BlockParser) consume(c map[lex.TokenKind]bool) {
	blockParser.parser.consume(c)
}

func (blockParser *BlockParser) parseStatementList(block *ast.Block, isGlobal bool) {
	block.Pos = blockParser.parser.mkPos()
	defer func() {
		block.EndPos = blockParser.parser.mkPos()
	}()
	isDefer := false
	resetDefer := func() {
		isDefer = false
	}
	validAfterDefer := func() error {
		if blockParser.parser.ExpressionParser.looksLikeExpression() ||
			blockParser.parser.token.Type == lex.TokenLc {
			return nil
		}
		return fmt.Errorf("%s not valid token '%s' after defer",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
	}
	var err error
	for lex.TokenEof != blockParser.parser.token.Type {
		if len(blockParser.parser.errs) > blockParser.parser.nErrors2Stop {
			break
		}
		if blockParser.parser.ExpressionParser.looksLikeExpression() {
			blockParser.parseExpressionStatement(block, isDefer)
			resetDefer()
			continue
		}
		switch blockParser.parser.token.Type {
		case lex.TokenSemicolon, lex.TokenLf: // may be empty statement
			resetDefer()
			blockParser.Next(lfNotToken) // look up next
			continue
		case lex.TokenDefer:
			blockParser.Next(lfIsToken)
			if err := validAfterDefer(); err != nil {
				blockParser.parser.errs = append(blockParser.parser.errs, err)
			} else {
				isDefer = true
			}
		case lex.TokenVar:
			pos := blockParser.parser.mkPos()
			blockParser.Next(lfIsToken) // skip var key word
			vs, err := blockParser.parser.parseVar()
			if err != nil {
				blockParser.consume(untilSemicolonOrLf)
				blockParser.Next(lfNotToken)
				continue
			}
			statement := &ast.Statement{
				Type: ast.StatementTypeExpression,
				Expression: &ast.Expression{
					Type: ast.ExpressionTypeVar,
					Data: vs,
					Pos:  pos,
				},
				Pos: pos,
			}
			block.Statements = append(block.Statements, statement)
			blockParser.parser.validStatementEnding()

		case lex.TokenIf:
			pos := blockParser.parser.mkPos()
			statement, err := blockParser.parseIf()
			if err != nil {
				blockParser.consume(untilRc)
				blockParser.Next(lfNotToken)
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:        ast.StatementTypeIf,
				StatementIf: statement,
				Pos:         pos,
			})
		case lex.TokenFor:
			pos := blockParser.parser.mkPos()
			statement, err := blockParser.parseFor()
			if err != nil {
				blockParser.consume(untilRc)
				blockParser.Next(lfNotToken)
				continue
			}
			statement.Block.IsForBlock = true
			block.Statements = append(block.Statements, &ast.Statement{
				Type:         ast.StatementTypeFor,
				StatementFor: statement,
				Pos:          pos,
			})
		case lex.TokenSwitch:
			pos := blockParser.parser.mkPos()
			statement, err := blockParser.parseSwitch()
			if err != nil {
				blockParser.consume(untilRc)
				blockParser.Next(lfNotToken)
				continue
			}
			if _, ok := statement.(*ast.StatementSwitch); ok {
				block.Statements = append(block.Statements, &ast.Statement{
					Type:            ast.StatementTypeSwitch,
					StatementSwitch: statement.(*ast.StatementSwitch),
					Pos:             pos,
				})
			} else {
				block.Statements = append(block.Statements, &ast.Statement{
					Type: ast.StatementTypeSwitchTemplate,
					StatementSwitchTemplate: statement.(*ast.StatementSwitchTemplate),
					Pos: pos,
				})
			}
		case lex.TokenConst:
			pos := blockParser.parser.mkPos()
			blockParser.Next(lfIsToken)
			cs, err := blockParser.parser.parseConst()
			if err != nil {
				blockParser.consume(untilSemicolonOrLf)
				blockParser.Next(lfNotToken)
				continue
			}
			statement := &ast.Statement{}
			statement.Type = ast.StatementTypeExpression
			statement.Pos = pos
			statement.Expression = &ast.Expression{
				Type: ast.ExpressionTypeConst,
				Data: cs,
				Pos:  pos,
			}
			block.Statements = append(block.Statements, statement)
			blockParser.parser.validStatementEnding()
			if blockParser.parser.token.Type == lex.TokenSemicolon {
				blockParser.Next(lfNotToken)
			}
		case lex.TokenReturn:
			pos := blockParser.parser.mkPos()
			if isGlobal {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s 'return' cannot used in packge init block",
						blockParser.parser.errorMsgPrefix()))
			}
			blockParser.Next(lfIsToken)
			r := &ast.StatementReturn{}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:            ast.StatementTypeReturn,
				StatementReturn: r,
				Pos:             pos,
			})
			if blockParser.parser.token.Type == lex.TokenRc {
				continue
			}
			if blockParser.parser.token.Type == lex.TokenLf ||
				blockParser.parser.token.Type == lex.TokenSemicolon ||
				blockParser.parser.token.Type == lex.TokenRc {
				blockParser.Next(lfNotToken)
				continue
			}
			var es []*ast.Expression
			es, err = blockParser.parser.ExpressionParser.parseExpressions(lex.TokenSemicolon)
			if err != nil {
				blockParser.parser.errs = append(blockParser.parser.errs, err)
				blockParser.consume(untilSemicolonOrLf)
				blockParser.Next(lfNotToken)
				continue
			}
			r.Expressions = es
			blockParser.parser.validStatementEnding()
			blockParser.Next(lfNotToken)
		case lex.TokenLc:
			pos := blockParser.parser.mkPos()
			newBlock := ast.Block{}
			blockParser.Next(lfNotToken) // skip {
			blockParser.parseStatementList(&newBlock, false)
			blockParser.parser.ifTokenIsLfThenSkip()
			if blockParser.parser.token.Type != lex.TokenRc {
				blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
					blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description))
				blockParser.consume(untilRc)
			}
			blockParser.Next(lfNotToken)
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
			resetDefer()
		case lex.TokenPass:
			pos := blockParser.parser.mkPos()
			if isGlobal == false {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s 'pass' can only be used in package init block",
						blockParser.parser.errorMsgPrefix()))
			}
			blockParser.Next(lfIsToken)
			blockParser.parser.validStatementEnding()
			block.Statements = append(block.Statements, &ast.Statement{
				Type:            ast.StatementTypeReturn,
				Pos:             pos,
				StatementReturn: &ast.StatementReturn{},
			})
		case lex.TokenContinue:
			pos := blockParser.parser.mkPos()
			blockParser.Next(lfIsToken)
			blockParser.parser.validStatementEnding()
			block.Statements = append(block.Statements, &ast.Statement{
				Type:              ast.StatementTypeContinue,
				StatementContinue: &ast.StatementContinue{},
				Pos:               pos,
			})
		case lex.TokenBreak:
			pos := blockParser.parser.mkPos()
			blockParser.Next(lfIsToken)
			blockParser.parser.validStatementEnding()
			block.Statements = append(block.Statements, &ast.Statement{
				Type:           ast.StatementTypeBreak,
				StatementBreak: &ast.StatementBreak{},
				Pos:            pos,
			})
		case lex.TokenGoto:
			pos := blockParser.parser.mkPos()
			blockParser.Next(lfIsToken) // skip goto key word
			if blockParser.parser.token.Type != lex.TokenIdentifier {
				blockParser.parser.errs = append(blockParser.parser.errs,
					fmt.Errorf("%s  missing identifier after goto statement, but '%s'",
						blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description))
				blockParser.consume(untilSemicolonOrLf)
				blockParser.Next(lfNotToken)
				continue
			}
			statementGoto := &ast.StatementGoTo{}
			statementGoto.LabelName = blockParser.parser.token.Data.(string)
			block.Statements = append(block.Statements, &ast.Statement{
				Type:          ast.StatementTypeGoTo,
				StatementGoTo: statementGoto,
				Pos:           pos,
			})
			blockParser.Next(lfIsToken)
			blockParser.parser.validStatementEnding()
			blockParser.Next(lfNotToken)
		case lex.TokenType:
			pos := blockParser.parser.mkPos()
			alias, err := blockParser.parser.parseTypeAlias()
			if err != nil {
				blockParser.consume(untilSemicolonOrLf)
				blockParser.Next(lfNotToken)
				continue
			}
			blockParser.parser.validStatementEnding()
			statement := &ast.Statement{}
			statement.Pos = pos
			statement.Type = ast.StatementTypeTypeAlias
			statement.TypeAlias = alias
			block.Statements = append(block.Statements, statement)
			blockParser.Next(lfNotToken)
		case lex.TokenClass, lex.TokenInterface:
			pos := blockParser.parser.mkPos()
			class, err := blockParser.parser.ClassParser.parse()
			if err != nil {
				continue
			}
			statement := &ast.Statement{}
			statement.Pos = pos
			statement.Type = ast.StatementTypeClass
			statement.Class = class
			block.Statements = append(block.Statements, statement)
		case lex.TokenEnum:
			pos := blockParser.parser.mkPos()
			e, err := blockParser.parser.parseEnum()
			if err != nil {
				continue
			}
			s := &ast.Statement{}
			s.Pos = pos
			s.Type = ast.StatementTypeEnum
			s.Enum = e
			block.Statements = append(block.Statements, s)
		case lex.TokenImport:
			ims := blockParser.parser.parseImports()
			for _, t := range ims {
				s := &ast.Statement{
					Type:   ast.StatementTypeImport,
					Import: t,
					Pos:    t.Pos,
				}
				block.Statements = append(block.Statements, s)
			}
		default:
			// something I cannot handle
			return
		}
	}
	return
}

func (blockParser *BlockParser) parseExpressionStatement(block *ast.Block, isDefer bool) (isLabel bool) {
	pos := blockParser.parser.mkPos()
	e, err := blockParser.parser.ExpressionParser.parseExpression(true)
	if err != nil {
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		blockParser.consume(untilSemicolonOrLf)
		blockParser.Next(lfNotToken)
		return
	}
	if e.Type == ast.ExpressionTypeIdentifier && blockParser.parser.token.Type == lex.TokenColon {
		//lable found , good...
		if isDefer {
			blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s defer mixup with statement lable has no meaning",
				blockParser.parser.errorMsgPrefix()))
		}
		isLabel = true
		blockParser.Next(lfIsToken) // skip :
		if blockParser.parser.token.Type != lex.TokenLf {
			blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s expect new line",
				blockParser.parser.errorMsgPrefix()))
		}
		statement := &ast.Statement{}
		statement.Pos = pos
		statement.Type = ast.StatementTypeLabel
		label := &ast.StatementLabel{}
		label.CodeOffset = -1
		statement.StatementLabel = label
		label.Statement = statement
		label.Name = e.Data.(*ast.ExpressionIdentifier).Name
		block.Statements = append(block.Statements, statement)
		label.Block = block
		err = block.Insert(label.Name, e.Pos, label) // insert first,so this label can be found before it is checked
		if err != nil {
			blockParser.parser.errs = append(blockParser.parser.errs, err)
		}
	} else {
		blockParser.parser.validStatementEnding()
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
	return
}

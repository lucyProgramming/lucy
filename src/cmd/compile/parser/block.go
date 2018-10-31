package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type BlockParser struct {
	parser *Parser
}

func (bp *BlockParser) Next(lfIsToken bool) {
	bp.parser.Next(lfIsToken)
}

func (bp *BlockParser) consume(c map[lex.TokenKind]bool) {
	bp.parser.consume(c)
}

func (bp *BlockParser) parseStatementList(block *ast.Block, isGlobal bool) {
	block.Pos = bp.parser.mkPos()
	defer func() {
		block.EndPos = bp.parser.mkPos()
	}()
	isDefer := false
	var deferPos *ast.Pos
	isAbstract := false
	isFinal := false
	var finalPos *ast.Pos
	comment := &CommentParser{
		parser: bp.parser,
	}
	resetPrefix := func() {
		isFinal = false
		isDefer = false
		deferPos = nil
		isAbstract = false
		comment.reset()
	}
	validAfterDefer := func() error {
		if bp.parser.ExpressionParser.looksLikeExpression() ||
			bp.parser.token.Type == lex.TokenLc {
			return nil
		}
		return fmt.Errorf("%s not valid token '%s' after defer",
			bp.parser.errMsgPrefix(), bp.parser.token.Description)
	}
	var err error
	for lex.TokenEof != bp.parser.token.Type {
		if len(bp.parser.errs) > bp.parser.nErrors2Stop {
			break
		}
		if bp.parser.ExpressionParser.looksLikeExpression() {
			bp.parseExpressionStatement(block, isDefer, deferPos)
			resetPrefix()
			continue
		}
		switch bp.parser.token.Type {
		case lex.TokenComment, lex.TokenMultiLineComment:
			comment.read()
		case lex.TokenSemicolon, lex.TokenLf: // may be empty statement
			resetPrefix()
			bp.Next(lfNotToken) // look up next
			continue
		case lex.TokenFinal:
			pos := bp.parser.mkPos()
			bp.parser.Next(lfIsToken)
			bp.parser.unExpectNewLineAndSkip()
			if err := bp.parser.validAfterFinal(); err != nil {
				isFinal = false
			} else {
				isFinal = true
				finalPos = pos
			}
			continue
		case lex.TokenDefer:
			pos := bp.parser.mkPos()
			bp.Next(lfIsToken)
			if err := validAfterDefer(); err != nil {
				bp.parser.errs = append(bp.parser.errs, err)
			} else {
				isDefer = true
				deferPos = pos
			}
		case lex.TokenVar:
			pos := bp.parser.mkPos()
			bp.Next(lfIsToken) // skip var key word
			vs, err := bp.parser.parseVar()
			if err != nil {
				bp.consume(untilSemicolonOrLf)
				bp.Next(lfNotToken)
				continue
			}
			statement := &ast.Statement{
				Type: ast.StatementTypeExpression,
				Expression: &ast.Expression{
					Type: ast.ExpressionTypeVar,
					Data: vs,
					Pos:  pos,
					Op:   "var",
				},
				Pos: pos,
			}
			block.Statements = append(block.Statements, statement)
			bp.parser.validStatementEnding()

		case lex.TokenIf:
			pos := bp.parser.mkPos()
			statement, err := bp.parseIf()
			if err != nil {
				bp.consume(untilRc)
				bp.Next(lfNotToken)
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:        ast.StatementTypeIf,
				StatementIf: statement,
				Pos:         pos,
			})
		case lex.TokenFor:
			pos := bp.parser.mkPos()
			statement, err := bp.parseFor()
			if err != nil {
				bp.consume(untilRc)
				bp.Next(lfNotToken)
				continue
			}
			statement.Block.IsForBlock = true
			block.Statements = append(block.Statements, &ast.Statement{
				Type:         ast.StatementTypeFor,
				StatementFor: statement,
				Pos:          pos,
			})
		case lex.TokenAbstract:
			bp.parser.Next(lfIsToken)
			bp.parser.unExpectNewLineAndSkip()
			if err := bp.parser.validAfterAbstract(); err == nil {
				isAbstract = true
			}
		case lex.TokenSwitch:
			pos := bp.parser.mkPos()
			statement, err := bp.parseSwitch()
			if err != nil {
				bp.consume(untilRc)
				bp.Next(lfNotToken)
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:            ast.StatementTypeSwitch,
				StatementSwitch: statement,
				Pos:             pos,
			})
		case lex.TokenWhen:
			pos := bp.parser.mkPos()
			statement, err := bp.parseWhen()
			if err != nil {
				bp.consume(untilRc)
				bp.Next(lfNotToken)
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:          ast.StatementTypeWhen,
				StatementWhen: statement,
				Pos:           pos,
			})
		case lex.TokenConst:
			pos := bp.parser.mkPos()
			bp.Next(lfIsToken)
			cs, err := bp.parser.parseConst()
			if err != nil {
				bp.consume(untilSemicolonOrLf)
				bp.Next(lfNotToken)
				continue
			}
			statement := &ast.Statement{}
			statement.Type = ast.StatementTypeExpression
			statement.Pos = pos
			statement.Expression = &ast.Expression{
				Type: ast.ExpressionTypeConst,
				Data: cs,
				Pos:  pos,
				Op:   "const",
			}
			block.Statements = append(block.Statements, statement)
			bp.parser.validStatementEnding()
			if bp.parser.token.Type == lex.TokenSemicolon {
				bp.Next(lfNotToken)
			}
		case lex.TokenReturn:
			if isGlobal {
				bp.parser.errs = append(bp.parser.errs,
					fmt.Errorf("%s 'return' cannot used in packge init block",
						bp.parser.errMsgPrefix()))
			}
			st := &ast.StatementReturn{
				Pos: bp.parser.mkPos(),
			}
			bp.Next(lfIsToken)
			block.Statements = append(block.Statements, &ast.Statement{
				Type:            ast.StatementTypeReturn,
				StatementReturn: st,
				Pos:             st.Pos,
			})
			if bp.parser.token.Type == lex.TokenRc {
				continue
			}
			if bp.parser.token.Type == lex.TokenRc ||
				bp.parser.token.Type == lex.TokenSemicolon ||
				bp.parser.token.Type == lex.TokenLf ||
				bp.parser.token.Type == lex.TokenComma ||
				bp.parser.token.Type == lex.TokenMultiLineComment {
				bp.Next(lfNotToken)
				continue
			}
			var es []*ast.Expression
			es, err = bp.parser.ExpressionParser.parseExpressions(lex.TokenSemicolon)
			if err != nil {
				bp.consume(untilSemicolonOrLf)
				bp.Next(lfNotToken)
				continue
			}
			st.Expressions = es
			bp.parser.validStatementEnding()
			bp.Next(lfNotToken)
		case lex.TokenLc:
			pos := bp.parser.mkPos()
			newBlock := ast.Block{}
			bp.Next(lfNotToken) // skip {
			bp.parseStatementList(&newBlock, false)
			bp.parser.ifTokenIsLfThenSkip()
			if bp.parser.token.Type != lex.TokenRc {
				bp.parser.errs = append(bp.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
					bp.parser.errMsgPrefix(), bp.parser.token.Description))
				bp.consume(untilRc)
			}
			bp.Next(lfNotToken)
			if isDefer {
				d := &ast.StatementDefer{
					Block: newBlock,
					Pos:   deferPos,
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
			resetPrefix()
		case lex.TokenPass:
			if isGlobal == false {
				bp.parser.errs = append(bp.parser.errs,
					fmt.Errorf("%s 'pass' can only be used in package init block",
						bp.parser.errMsgPrefix()))
			}
			pos := bp.parser.mkPos()
			bp.Next(lfIsToken)
			bp.parser.validStatementEnding()
			block.Statements = append(block.Statements, &ast.Statement{
				Type: ast.StatementTypeReturn,
				Pos:  pos,
				StatementReturn: &ast.StatementReturn{
					Pos: pos,
				},
			})
		case lex.TokenContinue:
			pos := bp.parser.mkPos()
			bp.Next(lfIsToken)
			bp.parser.validStatementEnding()
			block.Statements = append(block.Statements, &ast.Statement{
				Type: ast.StatementTypeContinue,
				StatementContinue: &ast.StatementContinue{
					Pos: pos,
				},
				Pos: pos,
			})
		case lex.TokenBreak:
			pos := bp.parser.mkPos()
			bp.Next(lfIsToken)
			bp.parser.validStatementEnding()
			block.Statements = append(block.Statements, &ast.Statement{
				Type: ast.StatementTypeBreak,
				StatementBreak: &ast.StatementBreak{
					Pos: pos,
				},
				Pos: pos,
			})
		case lex.TokenGoto:
			pos := bp.parser.mkPos()
			bp.Next(lfIsToken) // skip goto key word
			if bp.parser.token.Type != lex.TokenIdentifier {
				bp.parser.errs = append(bp.parser.errs,
					fmt.Errorf("%s  missing identifier after goto statement, but '%s'",
						bp.parser.errMsgPrefix(), bp.parser.token.Description))
				bp.consume(untilSemicolonOrLf)
				bp.Next(lfNotToken)
				continue
			}
			statementGoto := &ast.StatementGoTo{
				Pos: pos,
			}
			statementGoto.LabelName = bp.parser.token.Data.(string)
			block.Statements = append(block.Statements, &ast.Statement{
				Type:          ast.StatementTypeGoTo,
				StatementGoTo: statementGoto,
				Pos:           pos,
			})
			bp.Next(lfIsToken)
			bp.parser.validStatementEnding()
			bp.Next(lfNotToken)
		case lex.TokenTypeAlias:
			pos := bp.parser.mkPos()
			alias, err := bp.parser.parseTypeAlias(comment)
			if err != nil {
				bp.consume(untilSemicolonOrLf)
				bp.Next(lfNotToken)
				continue
			}
			bp.parser.validStatementEnding()
			statement := &ast.Statement{}
			statement.Pos = pos
			statement.Type = ast.StatementTypeTypeAlias
			statement.TypeAlias = alias
			block.Statements = append(block.Statements, statement)
			bp.Next(lfNotToken)
		case lex.TokenClass, lex.TokenInterface:
			pos := bp.parser.mkPos()
			class, _ := bp.parser.ClassParser.parse(isAbstract)
			if class != nil {
				statement := &ast.Statement{}
				statement.Pos = pos
				class.FinalPos = finalPos
				if isFinal {
					class.AccessFlags |= cg.ACC_CLASS_FINAL
				}
				statement.Type = ast.StatementTypeClass
				statement.Class = class
				block.Statements = append(block.Statements, statement)
			}

		case lex.TokenEnum:
			pos := bp.parser.mkPos()
			e, _ := bp.parser.parseEnum()
			if e != nil {
				s := &ast.Statement{}
				s.Pos = pos
				s.Type = ast.StatementTypeEnum
				s.Enum = e
				block.Statements = append(block.Statements, s)
			}
		case lex.TokenImport:
			pos := bp.parser.mkPos()
			ims := bp.parser.parseImports()
			for _, t := range ims {
				s := &ast.Statement{
					Type:   ast.StatementTypeImport,
					Import: t,
					Pos:    pos,
				}
				block.Statements = append(block.Statements, s)
			}
		case lex.TokenElse, lex.TokenElseif:
			bp.parser.errs = append(bp.parser.errs,
				fmt.Errorf("%s unexpected '%s'", bp.parser.errMsgPrefix(), bp.parser.token.Description))
			bp.Next(lfIsToken)

		default:
			// something I cannot handle
			return
		}
	}
	return
}

func (bp *BlockParser) parseExpressionStatement(block *ast.Block, isDefer bool, deferPos *ast.Pos) (isLabel bool) {
	pos := bp.parser.mkPos()
	e, err := bp.parser.ExpressionParser.parseExpression(true)
	if err != nil {
		bp.consume(untilSemicolonOrLf)
		bp.Next(lfNotToken)
		return
	}
	if e.Type == ast.ExpressionTypeIdentifier &&
		bp.parser.token.Type == lex.TokenColon {
		//lable found , good...
		if isDefer {
			bp.parser.errs = append(bp.parser.errs,
				fmt.Errorf("%s defer mixup with statement lable has no meaning",
					bp.parser.errMsgPrefix()))
		}
		isLabel = true
		pos := bp.parser.mkPos()
		bp.Next(lfIsToken) // skip :
		if bp.parser.token.Type != lex.TokenLf {
			bp.parser.errs = append(bp.parser.errs,
				fmt.Errorf("%s expect new line",
					bp.parser.errMsgPrefix()))
		}
		statement := &ast.Statement{}
		statement.Pos = pos
		statement.Type = ast.StatementTypeLabel
		label := &ast.StatementLabel{}
		label.Pos = pos
		label.CodeOffset = -1
		statement.StatementLabel = label
		label.Statement = statement
		label.Name = e.Data.(*ast.ExpressionIdentifier).Name
		block.Statements = append(block.Statements, statement)
		label.Block = block
		err = block.Insert(label.Name, e.Pos, label) // insert first,so this label can be found before it is checked
		if err != nil {
			bp.parser.errs = append(bp.parser.errs, err)
		}
	} else {
		bp.parser.validStatementEnding()
		if isDefer {
			d := &ast.StatementDefer{
				Pos: deferPos,
			}
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

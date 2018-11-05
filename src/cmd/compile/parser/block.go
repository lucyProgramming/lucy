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

func (this *BlockParser) Next(lfIsToken bool) {
	this.parser.Next(lfIsToken)
}

func (this *BlockParser) consume(c map[lex.TokenKind]bool) {
	this.parser.consume(c)
}

func (this *BlockParser) parseStatementList(block *ast.Block, isGlobal bool) {
	block.Pos = this.parser.mkPos()
	defer func() {
		block.EndPos = this.parser.mkPos()
	}()
	isDefer := false
	var deferPos *ast.Pos
	isAbstract := false
	isFinal := false
	var finalPos *ast.Pos
	comment := &CommentParser{
		parser: this.parser,
	}
	resetPrefix := func() {
		isFinal = false
		isDefer = false
		deferPos = nil
		isAbstract = false
		comment.reset()
	}
	validAfterDefer := func() error {
		if this.parser.ExpressionParser.looksLikeExpression() ||
			this.parser.token.Type == lex.TokenLc {
			return nil
		}
		return fmt.Errorf("%s not valid token '%s' after defer",
			this.parser.errMsgPrefix(), this.parser.token.Description)
	}
	var err error
	for lex.TokenEof != this.parser.token.Type {
		if len(this.parser.errs) > this.parser.nErrors2Stop {
			break
		}
		if this.parser.ExpressionParser.looksLikeExpression() {
			this.parseExpressionStatement(block, isDefer, deferPos)
			resetPrefix()
			continue
		}
		switch this.parser.token.Type {
		case lex.TokenComment, lex.TokenMultiLineComment:
			comment.read()
		case lex.TokenSemicolon, lex.TokenLf: // may be empty statement
			resetPrefix()
			this.Next(lfNotToken) // look up next
			continue
		case lex.TokenFinal:
			pos := this.parser.mkPos()
			this.parser.Next(lfIsToken)
			this.parser.unExpectNewLineAndSkip()
			if err := this.parser.validAfterFinal(); err != nil {
				isFinal = false
			} else {
				isFinal = true
				finalPos = pos
			}
			continue
		case lex.TokenDefer:
			pos := this.parser.mkPos()
			this.Next(lfIsToken)
			if err := validAfterDefer(); err != nil {
				this.parser.errs = append(this.parser.errs, err)
			} else {
				isDefer = true
				deferPos = pos
			}
		case lex.TokenVar:
			pos := this.parser.mkPos()
			this.Next(lfIsToken) // skip var key word
			vs, err := this.parser.parseVar()
			if err != nil {
				this.consume(untilSemicolonOrLf)
				this.Next(lfNotToken)
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
			this.parser.validStatementEnding()

		case lex.TokenIf:
			pos := this.parser.mkPos()
			statement, err := this.parseIf()
			if err != nil {
				this.consume(untilRc)
				this.Next(lfNotToken)
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:        ast.StatementTypeIf,
				StatementIf: statement,
				Pos:         pos,
			})
		case lex.TokenFor:
			pos := this.parser.mkPos()
			statement, err := this.parseFor()
			if err != nil {
				this.consume(untilRc)
				this.Next(lfNotToken)
				continue
			}
			statement.Block.IsForBlock = true
			block.Statements = append(block.Statements, &ast.Statement{
				Type:         ast.StatementTypeFor,
				StatementFor: statement,
				Pos:          pos,
			})
		case lex.TokenAbstract:
			this.parser.Next(lfIsToken)
			this.parser.unExpectNewLineAndSkip()
			if err := this.parser.validAfterAbstract(); err == nil {
				isAbstract = true
			}
		case lex.TokenSwitch:
			pos := this.parser.mkPos()
			statement, err := this.parseSwitch()
			if err != nil {
				this.consume(untilRc)
				this.Next(lfNotToken)
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:            ast.StatementTypeSwitch,
				StatementSwitch: statement,
				Pos:             pos,
			})
		case lex.TokenWhen:
			pos := this.parser.mkPos()
			statement, err := this.parseWhen()
			if err != nil {
				this.consume(untilRc)
				this.Next(lfNotToken)
				continue
			}
			block.Statements = append(block.Statements, &ast.Statement{
				Type:          ast.StatementTypeWhen,
				StatementWhen: statement,
				Pos:           pos,
			})
		case lex.TokenConst:
			pos := this.parser.mkPos()
			this.Next(lfIsToken)
			cs, err := this.parser.parseConst()
			if err != nil {
				this.consume(untilSemicolonOrLf)
				this.Next(lfNotToken)
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
			this.parser.validStatementEnding()
			if this.parser.token.Type == lex.TokenSemicolon {
				this.Next(lfNotToken)
			}
		case lex.TokenReturn:
			if isGlobal {
				this.parser.errs = append(this.parser.errs,
					fmt.Errorf("%s 'return' cannot used in packge init block",
						this.parser.errMsgPrefix()))
			}
			st := &ast.StatementReturn{
				Pos: this.parser.mkPos(),
			}
			this.Next(lfIsToken)
			block.Statements = append(block.Statements, &ast.Statement{
				Type:            ast.StatementTypeReturn,
				StatementReturn: st,
				Pos:             st.Pos,
			})
			if this.parser.token.Type == lex.TokenRc {
				continue
			}
			if this.parser.token.Type == lex.TokenRc ||
				this.parser.token.Type == lex.TokenSemicolon ||
				this.parser.token.Type == lex.TokenLf ||
				this.parser.token.Type == lex.TokenComma ||
				this.parser.token.Type == lex.TokenMultiLineComment {
				this.Next(lfNotToken)
				continue
			}
			var es []*ast.Expression
			es, err = this.parser.ExpressionParser.parseExpressions(lex.TokenSemicolon)
			if err != nil {
				this.consume(untilSemicolonOrLf)
				this.Next(lfNotToken)
				continue
			}
			st.Expressions = es
			this.parser.validStatementEnding()
			this.Next(lfNotToken)
		case lex.TokenLc:
			pos := this.parser.mkPos()
			newBlock := ast.Block{}
			this.Next(lfNotToken) // skip {
			this.parseStatementList(&newBlock, false)
			this.parser.ifTokenIsLfThenSkip()
			if this.parser.token.Type != lex.TokenRc {
				this.parser.errs = append(this.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
					this.parser.errMsgPrefix(), this.parser.token.Description))
				this.consume(untilRc)
			}
			this.Next(lfNotToken)
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
				this.parser.errs = append(this.parser.errs,
					fmt.Errorf("%s 'pass' can only be used in package init block",
						this.parser.errMsgPrefix()))
			}
			pos := this.parser.mkPos()
			this.Next(lfIsToken)
			this.parser.validStatementEnding()
			block.Statements = append(block.Statements, &ast.Statement{
				Type: ast.StatementTypeReturn,
				Pos:  pos,
				StatementReturn: &ast.StatementReturn{
					Pos: pos,
				},
			})
		case lex.TokenContinue:
			pos := this.parser.mkPos()
			this.Next(lfIsToken)
			this.parser.validStatementEnding()
			block.Statements = append(block.Statements, &ast.Statement{
				Type: ast.StatementTypeContinue,
				StatementContinue: &ast.StatementContinue{
					Pos: pos,
				},
				Pos: pos,
			})
		case lex.TokenBreak:
			pos := this.parser.mkPos()
			this.Next(lfIsToken)
			this.parser.validStatementEnding()
			block.Statements = append(block.Statements, &ast.Statement{
				Type: ast.StatementTypeBreak,
				StatementBreak: &ast.StatementBreak{
					Pos: pos,
				},
				Pos: pos,
			})
		case lex.TokenGoto:
			pos := this.parser.mkPos()
			this.Next(lfIsToken) // skip goto key word
			if this.parser.token.Type != lex.TokenIdentifier {
				this.parser.errs = append(this.parser.errs,
					fmt.Errorf("%s  missing identifier after goto statement, but '%s'",
						this.parser.errMsgPrefix(), this.parser.token.Description))
				this.consume(untilSemicolonOrLf)
				this.Next(lfNotToken)
				continue
			}
			statementGoto := &ast.StatementGoTo{
				Pos: pos,
			}
			statementGoto.LabelName = this.parser.token.Data.(string)
			block.Statements = append(block.Statements, &ast.Statement{
				Type:          ast.StatementTypeGoTo,
				StatementGoTo: statementGoto,
				Pos:           pos,
			})
			this.Next(lfIsToken)
			this.parser.validStatementEnding()
			this.Next(lfNotToken)
		case lex.TokenTypeAlias:
			pos := this.parser.mkPos()
			alias, err := this.parser.parseTypeAlias(comment)
			if err != nil {
				this.consume(untilSemicolonOrLf)
				this.Next(lfNotToken)
				continue
			}
			this.parser.validStatementEnding()
			statement := &ast.Statement{}
			statement.Pos = pos
			statement.Type = ast.StatementTypeTypeAlias
			statement.TypeAlias = alias
			block.Statements = append(block.Statements, statement)
			this.Next(lfNotToken)
		case lex.TokenClass, lex.TokenInterface:
			pos := this.parser.mkPos()
			class, _ := this.parser.ClassParser.parse(isAbstract)
			if class != nil {
				statement := &ast.Statement{}
				statement.Pos = pos
				class.FinalPos = finalPos
				if isFinal {
					class.AccessFlags |= cg.AccClassFinal
				}
				statement.Type = ast.StatementTypeClass
				statement.Class = class
				block.Statements = append(block.Statements, statement)
			}

		case lex.TokenEnum:
			pos := this.parser.mkPos()
			e, _ := this.parser.parseEnum()
			if e != nil {
				s := &ast.Statement{}
				s.Pos = pos
				s.Type = ast.StatementTypeEnum
				s.Enum = e
				block.Statements = append(block.Statements, s)
			}
		case lex.TokenImport:
			pos := this.parser.mkPos()
			ims := this.parser.parseImports()
			for _, t := range ims {
				s := &ast.Statement{
					Type:   ast.StatementTypeImport,
					Import: t,
					Pos:    pos,
				}
				block.Statements = append(block.Statements, s)
			}
		case lex.TokenElse, lex.TokenElseif:
			this.parser.errs = append(this.parser.errs,
				fmt.Errorf("%s unexpected '%s'", this.parser.errMsgPrefix(), this.parser.token.Description))
			this.Next(lfIsToken)

		default:
			// something I cannot handle
			return
		}
	}
	return
}

func (this *BlockParser) parseExpressionStatement(block *ast.Block, isDefer bool, deferPos *ast.Pos) (isLabel bool) {
	pos := this.parser.mkPos()
	e, err := this.parser.ExpressionParser.parseExpression(true)
	if err != nil {
		this.consume(untilSemicolonOrLf)
		this.Next(lfNotToken)
		return
	}
	if e.Type == ast.ExpressionTypeIdentifier &&
		this.parser.token.Type == lex.TokenColon {
		//lable found , good...
		if isDefer {
			this.parser.errs = append(this.parser.errs,
				fmt.Errorf("%s defer mixup with statement lable has no meaning",
					this.parser.errMsgPrefix()))
		}
		isLabel = true
		pos := this.parser.mkPos()
		this.Next(lfIsToken) // skip :
		if this.parser.token.Type != lex.TokenLf {
			this.parser.errs = append(this.parser.errs,
				fmt.Errorf("%s expect new line",
					this.parser.errMsgPrefix()))
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
			this.parser.errs = append(this.parser.errs, err)
		}
	} else {
		this.parser.validStatementEnding()
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

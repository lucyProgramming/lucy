package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func Parse(tops *[]*ast.TopNode, filename string, bs []byte, onlyParseImport bool, nErrors2Stop int) []error {
	p := &Parser{
		bs:              bs,
		tops:            tops,
		filename:        filename,
		onlyParseImport: onlyParseImport,
		nErrors2Stop:    nErrors2Stop,
	}
	return p.Parse()
}

type Parser struct {
	onlyParseImport        bool
	bs                     []byte
	tops                   *[]*ast.TopNode
	lexer                  *lex.Lexer
	filename               string
	lastToken              *lex.Token
	token                  *lex.Token
	errs                   []error
	importsByAccessName    map[string]*ast.Import
	importsByResourceName  map[string]*ast.Import
	nErrors2Stop           int
	consumeFoundValidToken bool
	ExpressionParser       *ExpressionParser
	FunctionParser         *FunctionParser
	ClassParser            *ClassParser
	BlockParser            *BlockParser
}

/*
	call before parse source file
*/
func (parser *Parser) initParser() {
	parser.ExpressionParser = &ExpressionParser{parser}
	parser.FunctionParser = &FunctionParser{}
	parser.FunctionParser.parser = parser
	parser.ClassParser = &ClassParser{}
	parser.ClassParser.parser = parser
	parser.BlockParser = &BlockParser{}
	parser.BlockParser.parser = parser
}

func (parser *Parser) Parse() []error {
	parser.initParser()
	parser.lexer = lex.New(parser.bs, 1, 1)
	parser.Next(lfNotToken) //
	if parser.token.Type == lex.TokenEof {
		//TODO::empty source file , should forbidden???
		return nil
	}
	for _, t := range parser.parseImports() {
		parser.insertImports(t)
	}
	if parser.onlyParseImport { // only parse imports
		return parser.errs
	}
	var accessControlToken *lex.Token
	isFinal := false
	isAbstract := false
	var finalPos *ast.Pos
	resetProperty := func() {
		accessControlToken = nil
		isFinal = false
		isAbstract = false
		finalPos = nil
	}
	isPublic := func() bool {
		return accessControlToken != nil && accessControlToken.Type == lex.TokenPublic
	}
	for parser.token.Type != lex.TokenEof {
		if len(parser.errs) > parser.nErrors2Stop {
			break
		}
		switch parser.token.Type {
		case lex.TokenSemicolon, lex.TokenLf: // empty statement, no big deal
			parser.Next(lfNotToken)
			continue
		case lex.TokenPublic:
			accessControlToken = parser.token
			parser.Next(lfIsToken)
			parser.unExpectNewLineAndSkip()
			if err := parser.validAfterPublic(); err != nil {
				accessControlToken = nil
			}
			continue
		case lex.TokenAbstract:
			parser.Next(lfIsToken)
			parser.unExpectNewLineAndSkip()
			if err := parser.validAfterAbstract(); err == nil {
				isAbstract = true
			}
		case lex.TokenFinal:
			pos := parser.mkPos()
			parser.Next(lfIsToken)
			parser.unExpectNewLineAndSkip()
			if err := parser.validAfterFinal(); err != nil {
				isFinal = false
			} else {
				isFinal = true
				finalPos = pos
			}
			continue
		case lex.TokenVar:
			pos := parser.mkPos()
			parser.Next(lfIsToken) // skip var key word
			vs, err := parser.parseVar()
			if err != nil {
				parser.consume(untilSemicolonOrLf)
				parser.Next(lfNotToken)
				continue
			}
			isPublic := isPublic()
			e := &ast.Expression{
				Type:        ast.ExpressionTypeVar,
				Data:        vs,
				Pos:         pos,
				IsPublic:    isPublic,
				Description: "var",
			}
			*parser.tops = append(*parser.tops, &ast.TopNode{
				Data: e,
			})
			resetProperty()
		case lex.TokenIdentifier:
			e, err := parser.ExpressionParser.parseExpression(true)
			if err != nil {
				parser.errs = append(parser.errs, err)
				parser.consume(untilSemicolonOrLf)
				parser.Next(lfNotToken)
				continue
			}
			e.IsPublic = isPublic()
			parser.validStatementEnding()
			if e.Type == ast.ExpressionTypeVarAssign {
				*parser.tops = append(*parser.tops, &ast.TopNode{
					Data: e,
				})
			} else {
				parser.errs = append(parser.errs, fmt.Errorf("%s cannot have expression '%s' in top",
					parser.errorMsgPrefix(e.Pos), e.Description))
			}
			resetProperty()
		case lex.TokenEnum:
			e, err := parser.parseEnum()
			if err != nil {
				resetProperty()
				continue
			}
			isPublic := isPublic()
			if isPublic {
				e.AccessFlags |= cg.ACC_CLASS_PUBLIC
			}
			if e != nil {
				*parser.tops = append(*parser.tops, &ast.TopNode{
					Data: e,
				})
			}
			resetProperty()
		case lex.TokenFn:
			f, err := parser.FunctionParser.parse(true, false)
			if err != nil {
				parser.Next(lfNotToken)
				continue
			}
			isPublic := isPublic()
			if isPublic {
				f.AccessFlags |= cg.ACC_METHOD_PUBLIC
			}
			*parser.tops = append(*parser.tops, &ast.TopNode{
				Data: f,
			})
			resetProperty()
		case lex.TokenLc:
			b := &ast.Block{}
			parser.Next(lfNotToken) // skip {
			parser.BlockParser.parseStatementList(b, true)
			if parser.token.Type != lex.TokenRc {
				parser.errs = append(parser.errs, fmt.Errorf("%s expect '}', but '%s'",
					parser.errorMsgPrefix(), parser.token.Description))
				parser.consume(untilRc)
			}
			parser.Next(lfNotToken) // skip }
			*parser.tops = append(*parser.tops, &ast.TopNode{
				Data: b,
			})

		case lex.TokenClass, lex.TokenInterface:
			c, err := parser.ClassParser.parse(isAbstract)
			if err != nil {
				resetProperty()
				continue
			}
			*parser.tops = append(*parser.tops, &ast.TopNode{
				Data: c,
			})
			isPublic := isPublic()
			if isPublic {
				c.AccessFlags |= cg.ACC_CLASS_PUBLIC
			}
			if isAbstract {
				c.AccessFlags |= cg.ACC_CLASS_ABSTRACT
			}
			if isFinal {
				c.AccessFlags |= cg.ACC_CLASS_FINAL
				c.FinalPos = finalPos
			}
			resetProperty()
		case lex.TokenConst:
			parser.Next(lfIsToken) // skip const key word
			cs, err := parser.parseConst()
			if err != nil {
				parser.consume(untilSemicolonOrLf)
				parser.Next(lfNotToken)
				resetProperty()
				continue
			}
			isPublic := isPublic()
			for _, v := range cs {
				if isPublic {
					v.AccessFlags |= cg.ACC_FIELD_PUBLIC
				}
				*parser.tops = append(*parser.tops, &ast.TopNode{
					Data: v,
				})

			}
			resetProperty()
			continue
		case lex.TokenType:
			a, err := parser.parseTypeAlias()
			if err != nil {
				parser.consume(untilSemicolonOrLf)
				parser.Next(lfNotToken)
				resetProperty()
				continue
			}
			*parser.tops = append(*parser.tops, &ast.TopNode{
				Data: a,
			})
		case lex.TokenImport:
			pos := parser.mkPos()
			parser.parseImports()
			parser.errs = append(parser.errs, fmt.Errorf("%s cannot have import at this scope",
				parser.errorMsgPrefix(pos)))
		case lex.TokenEof:
			break
		default:
			if parser.ExpressionParser.looksLikeExpression() {
				e, err := parser.ExpressionParser.parseExpression(true)
				if err != nil {
					parser.errs = append(parser.errs, err)
					continue
				}
				if e.Type == ast.ExpressionTypeVarAssign {
					*parser.tops = append(*parser.tops, &ast.TopNode{
						Data: e,
					})
				} else {
					parser.errs = append(parser.errs, fmt.Errorf("%s cannot have expression '%s' in top",
						parser.errorMsgPrefix(e.Pos), e.Description))
				}
				continue
			}
			parser.errs = append(parser.errs, fmt.Errorf("%s token '%s' is not except",
				parser.errorMsgPrefix(), parser.token.Description))
			parser.consume(untilSemicolonOrLf)
			resetProperty()
		}
	}
	return parser.errs
}

func (parser *Parser) validAfterPublic() error {
	if parser.token.Type == lex.TokenFn ||
		parser.token.Type == lex.TokenClass ||
		parser.token.Type == lex.TokenEnum ||
		parser.token.Type == lex.TokenIdentifier ||
		parser.token.Type == lex.TokenInterface ||
		parser.token.Type == lex.TokenConst ||
		parser.token.Type == lex.TokenVar ||
		parser.token.Type == lex.TokenFinal ||
		parser.token.Type == lex.TokenAbstract {
		return nil
	}
	err := fmt.Errorf("%s cannot have token '%s' after 'public'",
		parser.errorMsgPrefix(), parser.token.Description)
	parser.errs = append(parser.errs, err)
	return err
}
func (parser *Parser) validAfterAbstract() error {
	if parser.token.Type == lex.TokenClass {
		return nil
	}
	err := fmt.Errorf("%s cannot have token '%s' after 'abstract'",
		parser.errorMsgPrefix(), parser.token.Description)
	parser.errs = append(parser.errs, err)
	return err
}
func (parser *Parser) validAfterFinal() error {
	if parser.token.Type == lex.TokenClass ||
		parser.token.Type == lex.TokenInterface {
		return nil
	}
	err := fmt.Errorf("%s cannot have token '%s' after 'final'",
		parser.errorMsgPrefix(), parser.token.Description)
	parser.errs = append(parser.errs, err)
	return err
}

/*
	statement ending
*/
func (parser *Parser) isStatementEnding() bool {
	return parser.token.Type == lex.TokenSemicolon ||
		parser.token.Type == lex.TokenLf ||
		parser.token.Type == lex.TokenRc
}
func (parser *Parser) validStatementEnding() error {
	if parser.isStatementEnding() {
		return nil
	}
	token := parser.token
	err := fmt.Errorf("%s expect semicolon or new line", parser.errorMsgPrefix(&ast.Pos{
		Filename:    parser.filename,
		StartLine:   token.StartLine,
		StartColumn: token.StartColumn,
	}))
	parser.errs = append(parser.errs, err)
	return nil
}

func (parser *Parser) mkPos() *ast.Pos {
	return &ast.Pos{
		Filename:    parser.filename,
		StartLine:   parser.token.StartLine,
		StartColumn: parser.token.StartColumn,
		Offset:      parser.lexer.GetOffSet(),
	}
}

// str := "hello world"   a,b = 123 or a b ;
func (parser *Parser) parseConst() (constants []*ast.Constant, err error) {
	names, err := parser.parseNameList()
	if err != nil {
		return
	}
	constants = make([]*ast.Constant, len(names))
	for k, v := range names {
		vd := &ast.Constant{}
		vd.Name = v.Name
		vd.Pos = v.Pos
		constants[k] = vd
	}
	var variableType *ast.Type
	if parser.isValidTypeBegin() {
		variableType, err = parser.parseType()
		if err != nil {
			return
		}
	}
	if variableType != nil {
		for _, c := range constants {
			c.Type = variableType.Clone()
		}
	}
	if parser.token.Type != lex.TokenAssign {
		err = fmt.Errorf("%s missing assign", parser.errorMsgPrefix())
		parser.errs = append(parser.errs, err)
		return
	}
	parser.Next(lfNotToken) // skip =
	es, err := parser.ExpressionParser.parseExpressions(lex.TokenSemicolon)
	if err != nil {
		return
	}
	if len(es) != len(constants) {
		err = fmt.Errorf("%s cannot assign %d value to %d constant",
			parser.errorMsgPrefix(), len(es), len(constants))
		parser.errs = append(parser.errs, err)
	}
	for k, _ := range constants {
		if k < len(es) {
			constants[k].Expression = es[k]
		}
	}
	return
}

// str := "hello world"   a,b = 123 or a b ;
func (parser *Parser) parseVar() (ret *ast.ExpressionVar, err error) {
	names, err := parser.parseNameList()
	if err != nil {
		return
	}
	ret = &ast.ExpressionVar{}
	ret.Variables = make([]*ast.Variable, len(names))
	for k, v := range names {
		vd := &ast.Variable{}
		vd.Name = v.Name
		vd.Pos = v.Pos
		ret.Variables[k] = vd
	}
	if parser.token.Type != lex.TokenAssign {
		ret.Type, err = parser.parseType()
		if err != nil {
			return
		}
	}
	if parser.token.Type == lex.TokenAssign {
		parser.Next(lfNotToken) // skip = or :=
		ret.InitValues, err = parser.ExpressionParser.parseExpressions(lex.TokenSemicolon)
		if err != nil {
			parser.errs = append(parser.errs, err)
			return
		}
	}
	return
}

func (parser *Parser) Next(lfIsToken bool) {
	if parser.consumeFoundValidToken {
		parser.consumeFoundValidToken = false
		return
	}
	var err error
	var tok *lex.Token
	parser.lastToken = parser.token
	for {
		tok, err = parser.lexer.Next()
		if err != nil {
			parser.errs = append(parser.errs,
				fmt.Errorf("%s %s", parser.errorMsgPrefix(), err.Error()))
		}
		if tok == nil {
			continue
		}
		parser.token = tok
		if lfIsToken {
			break
		}
		if tok.Type != lex.TokenLf {
			break
		}
	}
	return
}

/*
	errorMsgPrefix(pos) only receive one argument
*/
func (parser *Parser) errorMsgPrefix(pos ...*ast.Pos) string {
	var line, column int
	if len(pos) > 0 {
		line = pos[0].StartLine
		column = pos[0].StartColumn
	} else {
		line, column = parser.token.StartLine, parser.token.StartColumn
	}
	return fmt.Sprintf("%s:%d:%d", parser.filename, line, column)
}

func (parser *Parser) consume(until map[lex.TokenKind]bool) {
	if len(until) == 0 {
		panic("no token to consume")
	}
	for parser.token.Type != lex.TokenEof {
		if parser.token.Type == lex.TokenPublic ||
			parser.token.Type == lex.TokenProtected ||
			parser.token.Type == lex.TokenPrivate ||
			parser.token.Type == lex.TokenClass ||
			parser.token.Type == lex.TokenInterface ||
			parser.token.Type == lex.TokenFn ||
			parser.token.Type == lex.TokenFor ||
			parser.token.Type == lex.TokenIf ||
			parser.token.Type == lex.TokenSwitch ||
			parser.token.Type == lex.TokenEnum ||
			parser.token.Type == lex.TokenConst ||
			parser.token.Type == lex.TokenVar ||
			parser.token.Type == lex.TokenImport ||
			parser.token.Type == lex.TokenType ||
			parser.token.Type == lex.TokenGoto ||
			parser.token.Type == lex.TokenBreak ||
			parser.token.Type == lex.TokenContinue ||
			parser.token.Type == lex.TokenDefer ||
			parser.token.Type == lex.TokenReturn ||
			parser.token.Type == lex.TokenPass ||
			parser.token.Type == lex.TokenExtends ||
			parser.token.Type == lex.TokenImplements ||
			parser.token.Type == lex.TokenGlobal ||
			parser.token.Type == lex.TokenCase ||
			parser.token.Type == lex.TokenDefault {
			parser.consumeFoundValidToken = true
			return
		}
		if parser.token.Type == lex.TokenLc {
			if _, ok := until[lex.TokenLc]; ok == false {
				parser.consumeFoundValidToken = true
				return
			}
		}
		if parser.token.Type == lex.TokenRc {
			if _, ok := until[lex.TokenRc]; ok == false {
				parser.consumeFoundValidToken = true
				return
			}
		}
		if _, ok := until[parser.token.Type]; ok {
			return
		}
		parser.Next(lfIsToken)
	}
}

func (parser *Parser) ifTokenIsLfThenSkip() {
	if parser.token.Type == lex.TokenLf {
		parser.Next(lfNotToken)
	}
}

func (parser *Parser) unExpectNewLineAndSkip() {
	if err := parser.unExpectNewLine(); err != nil {
		parser.Next(lfNotToken)
	}
}
func (parser *Parser) unExpectNewLine() error {
	var err error
	if parser.token.Type == lex.TokenLf {
		err = fmt.Errorf("%s unexpected new line",
			parser.errorMsgPrefix(parser.mkPos()))
		parser.errs = append(parser.errs, err)
	}
	return err
}
func (parser *Parser) expectNewLineAndSkip() {
	if err := parser.expectNewLine(); err == nil {
		parser.Next(lfNotToken)
	}
}
func (parser *Parser) expectNewLine() error {
	var err error
	if parser.token.Type != lex.TokenLf {
		err = fmt.Errorf("%s expect new line , but '%s'",
			parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
	}
	return err
}

func (parser *Parser) parseTypeAlias() (*ast.TypeAlias, error) {
	parser.Next(lfIsToken) // skip type key word
	parser.unExpectNewLineAndSkip()
	if parser.token.Type != lex.TokenIdentifier {
		err := fmt.Errorf("%s expect identifer,but '%s'", parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return nil, err
	}
	ret := &ast.TypeAlias{}
	ret.Pos = parser.mkPos()
	ret.Name = parser.token.Data.(string)
	parser.Next(lfIsToken) // skip identifier
	if parser.token.Type != lex.TokenAssign {
		err := fmt.Errorf("%s expect '=',but '%s'", parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return nil, err
	}
	parser.Next(lfNotToken) // skip =
	var err error
	ret.Type, err = parser.parseType()
	if err != nil {
		return nil, err
	}
	return ret, err
}

/*
	a int
	int
*/
func (parser *Parser) parseTypedName() (vs []*ast.Variable, err error) {
	if parser.token.Type != lex.TokenIdentifier {
		/*
			not identifier begin
			must be type
			// int
		*/
		pos := parser.mkPos()
		t, err := parser.parseType()
		if err != nil {
			return nil, err
		}
		v := &ast.Variable{}
		v.Type = t
		v.Pos = pos
		return []*ast.Variable{v}, nil
	}
	names, err := parser.parseNameList()
	if err != nil {
		return nil, err
	}
	if parser.isValidTypeBegin() {
		/*
			a , b int
		*/
		t, err := parser.parseType()
		if err != nil {
			return nil, err
		}
		vs = make([]*ast.Variable, len(names))
		for k, v := range names {
			vd := &ast.Variable{}
			vs[k] = vd
			vd.Name = v.Name
			vd.Pos = v.Pos
			vd.Type = t.Clone()
			vd.Type.Pos = v.Pos // override pos
		}
		return vs, nil
	} else {
		/*
			syntax a,b
			not valid type begins, "a" and b must indicate types not double
		*/

		vs = make([]*ast.Variable, len(names))
		for k, v := range names {
			vd := &ast.Variable{}
			vs[k] = vd
			vd.Pos = v.Pos
			vd.Type = &ast.Type{
				Type: ast.VariableTypeName,
				Pos:  v.Pos,
				Name: v.Name,
			}
			vd.Type.Pos = v.Pos // override pos
		}
		return vs, nil
	}
}

// a,b int or int,bool  c xxx
func (parser *Parser) parseTypedNames() (vs []*ast.Variable, err error) {
	vs = []*ast.Variable{}
	for parser.token.Type != lex.TokenEof {
		ns, err := parser.parseNameList()
		if err != nil {
			return vs, err
		}
		t, err := parser.parseType()
		if err != nil {
			return vs, err
		}
		for _, v := range ns {
			vd := &ast.Variable{}
			vd.Name = v.Name
			vd.Pos = v.Pos
			vd.Type = t.Clone()
			vs = append(vs, vd)
		}
		if parser.token.Type != lex.TokenComma { // not a comma
			break
		} else {
			parser.Next(lfNotToken)
		}
	}
	return vs, nil
}

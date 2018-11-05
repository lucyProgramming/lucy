package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func Parse(
	tops *[]*ast.TopNode,
	filename string,
	bs []byte,
	onlyParseImport bool,
	nErrors2Stop int) []error {
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
func (this *Parser) initParser() {
	this.ExpressionParser = &ExpressionParser{this}
	this.FunctionParser = &FunctionParser{
		parser: this,
	}
	this.ClassParser = &ClassParser{
		parser: this,
	}
	this.BlockParser = &BlockParser{
		parser: this,
	}
}

func (this *Parser) Parse() []error {
	this.initParser()
	this.lexer = lex.New(this.bs, 1, 1)
	this.Next(lfNotToken) //
	if this.token.Type == lex.TokenEof {
		return nil
	}
	for _, t := range this.parseImports() {
		*this.tops = append(*this.tops, &ast.TopNode{
			Node: t,
		})
	}
	if this.onlyParseImport { // only parse imports
		return this.errs
	}
	var accessControlToken *lex.Token
	isFinal := false
	isAbstract := false
	var finalPos *ast.Pos
	comment := &CommentParser{
		parser: this,
	}
	resetProperty := func() {
		accessControlToken = nil
		isFinal = false
		isAbstract = false
		finalPos = nil
		comment.reset()
	}
	isPublic := func() bool {
		return accessControlToken != nil && accessControlToken.Type == lex.TokenPublic
	}
	for this.token.Type != lex.TokenEof {
		if len(this.errs) > this.nErrors2Stop {
			break
		}

		switch this.token.Type {
		case lex.TokenComment, lex.TokenMultiLineComment:
			comment.read()
		case lex.TokenSemicolon, lex.TokenLf: // empty statement, no big deal
			this.Next(lfNotToken)
			continue
		case lex.TokenPublic:
			accessControlToken = this.token
			this.Next(lfIsToken)
			this.unExpectNewLineAndSkip()
			if err := this.validAfterPublic(); err != nil {
				accessControlToken = nil
			}
			continue
		case lex.TokenAbstract:
			this.Next(lfIsToken)
			this.unExpectNewLineAndSkip()
			if err := this.validAfterAbstract(); err == nil {
				isAbstract = true
			}
		case lex.TokenFinal:
			pos := this.mkPos()
			this.Next(lfIsToken)
			this.unExpectNewLineAndSkip()
			if err := this.validAfterFinal(); err != nil {
				isFinal = false
			} else {
				isFinal = true
				finalPos = pos
			}
			continue
		case lex.TokenVar:
			pos := this.mkPos()
			this.Next(lfIsToken) // skip var key word
			vs, err := this.parseVar()
			if err != nil {
				this.consume(untilSemicolonOrLf)
				this.Next(lfNotToken)
				continue
			}
			isPublic := isPublic()
			e := &ast.Expression{
				Type:     ast.ExpressionTypeVar,
				Data:     vs,
				Pos:      pos,
				IsPublic: isPublic,
				IsGlobal: true,
				Op:       "var",
			}
			*this.tops = append(*this.tops, &ast.TopNode{
				Node: e,
			})
			resetProperty()
		case lex.TokenEnum:
			e, err := this.parseEnum()
			if err != nil {
				resetProperty()
				continue
			}
			e.Comment = comment.Comment
			isPublic := isPublic()
			if isPublic {
				e.AccessFlags |= cg.AccClassPublic
			}
			if e != nil {
				*this.tops = append(*this.tops, &ast.TopNode{
					Node: e,
				})
			}
			resetProperty()
		case lex.TokenFn:
			f, err := this.FunctionParser.parse(true, false)
			if err != nil {
				this.Next(lfNotToken)
				continue
			}
			f.Comment = comment.Comment
			isPublic := isPublic()
			if isPublic {
				f.AccessFlags |= cg.AccMethodPublic
			}
			*this.tops = append(*this.tops, &ast.TopNode{
				Node: f,
			})
			resetProperty()
		case lex.TokenLc:
			b := &ast.Block{}
			this.Next(lfNotToken) // skip {
			this.BlockParser.parseStatementList(b, true)
			if this.token.Type != lex.TokenRc {
				this.errs = append(this.errs, fmt.Errorf("%s expect '}', but '%s'",
					this.errMsgPrefix(), this.token.Description))
				this.consume(untilRc)
			}
			this.Next(lfNotToken) // skip }
			*this.tops = append(*this.tops, &ast.TopNode{
				Node: b,
			})

		case lex.TokenClass, lex.TokenInterface:
			c, err := this.ClassParser.parse(isAbstract)
			if err != nil {
				resetProperty()
				continue
			}
			c.Comment = comment.Comment
			*this.tops = append(*this.tops, &ast.TopNode{
				Node: c,
			})
			isPublic := isPublic()
			if isPublic {
				c.AccessFlags |= cg.AccClassPublic
			}
			if isAbstract {
				c.AccessFlags |= cg.AccClassAbstract
			}
			if isFinal {
				c.AccessFlags |= cg.AccClassFinal
				c.FinalPos = finalPos
			}
			resetProperty()
		case lex.TokenConst:
			this.Next(lfIsToken) // skip const key word
			cs, err := this.parseConst()
			if err != nil {
				this.consume(untilSemicolonOrLf)
				this.Next(lfNotToken)
				resetProperty()
				continue
			}
			isPublic := isPublic()
			for _, v := range cs {
				if isPublic {
					v.AccessFlags |= cg.AccFieldPublic
				}
				*this.tops = append(*this.tops, &ast.TopNode{
					Node: v,
				})
			}
			resetProperty()
			continue
		case lex.TokenTypeAlias:
			a, err := this.parseTypeAlias(comment)
			if err != nil {
				this.consume(untilSemicolonOrLf)
				this.Next(lfNotToken)
				resetProperty()
				continue
			}
			*this.tops = append(*this.tops, &ast.TopNode{
				Node: a,
			})
		case lex.TokenImport:
			pos := this.mkPos()
			this.parseImports()
			this.errs = append(this.errs, fmt.Errorf("%s cannot have import at this scope",
				this.errMsgPrefix(pos)))
		case lex.TokenEof:
			break
		default:
			if this.ExpressionParser.looksLikeExpression() {
				e, err := this.ExpressionParser.parseExpression(true)
				if err != nil {
					continue
				}
				e.IsPublic = isPublic()
				e.IsGlobal = true
				*this.tops = append(*this.tops, &ast.TopNode{
					Node: e,
				})
			} else {
				this.errs = append(this.errs, fmt.Errorf("%s token '%s' is not except",
					this.errMsgPrefix(), this.token.Description))
				this.Next(lfNotToken)
			}
		}
	}
	return this.errs
}

func (this *Parser) validAfterPublic() error {
	if this.token.Type == lex.TokenFn ||
		this.token.Type == lex.TokenClass ||
		this.token.Type == lex.TokenEnum ||
		this.token.Type == lex.TokenIdentifier ||
		this.token.Type == lex.TokenInterface ||
		this.token.Type == lex.TokenConst ||
		this.token.Type == lex.TokenVar ||
		this.token.Type == lex.TokenFinal ||
		this.token.Type == lex.TokenAbstract {
		return nil
	}
	err := fmt.Errorf("%s cannot have token '%s' after 'public'",
		this.errMsgPrefix(), this.token.Description)
	this.errs = append(this.errs, err)
	return err
}
func (this *Parser) validAfterAbstract() error {
	if this.token.Type == lex.TokenClass {
		return nil
	}
	err := fmt.Errorf("%s cannot have token '%s' after 'abstract'",
		this.errMsgPrefix(), this.token.Description)
	this.errs = append(this.errs, err)
	return err
}
func (this *Parser) validAfterFinal() error {
	if this.token.Type == lex.TokenClass ||
		this.token.Type == lex.TokenInterface {
		return nil
	}
	err := fmt.Errorf("%s cannot have token '%s' after 'final'",
		this.errMsgPrefix(), this.token.Description)
	this.errs = append(this.errs, err)
	return err
}

/*
	statement ending
*/
func (this *Parser) isStatementEnding() bool {
	return this.token.Type == lex.TokenSemicolon ||
		this.token.Type == lex.TokenLf ||
		this.token.Type == lex.TokenRc ||
		this.token.Type == lex.TokenComment ||
		this.token.Type == lex.TokenMultiLineComment
}
func (this *Parser) validStatementEnding() error {
	if this.isStatementEnding() {
		return nil
	}
	token := this.token
	err := fmt.Errorf("%s expect semicolon or new line", this.errMsgPrefix(&ast.Pos{
		Filename: this.filename,
		Line:     token.StartLine,
		Column:   token.StartColumn,
	}))
	this.errs = append(this.errs, err)
	return nil
}

func (this *Parser) mkPos() *ast.Pos {
	if this.token != nil {
		return &ast.Pos{
			Filename: this.filename,
			Line:     this.token.EndLine,
			Column:   this.token.EndColumn,
			Offset:   this.lexer.GetOffSet(),
		}
	} else {
		line, column := this.lexer.GetLineAndColumn()
		pos := &ast.Pos{
			Filename: this.filename,
			Line:     line,
			Column:   column,
		}
		return pos
	}
}

func (this *Parser) mkEndPos() *ast.Pos {
	if this.lastToken == nil {
		return &ast.Pos{
			Filename: this.filename,
			Line:     this.token.EndLine,
			Column:   this.token.EndColumn,
			Offset:   this.lexer.GetOffSet(),
		}
	} else {
		return &ast.Pos{
			Filename: this.filename,
			Line:     this.lastToken.EndLine,
			Column:   this.lastToken.EndColumn,
		}
	}
}

// str := "hello world"   a,b = 123 or a b ;
func (this *Parser) parseConst() (constants []*ast.Constant, err error) {
	names, err := this.parseNameList()
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
	if this.isValidTypeBegin() {
		variableType, err = this.parseType()
		if err != nil {
			return
		}
	}
	if variableType != nil {
		for _, c := range constants {
			c.Type = variableType.Clone()
		}
	}
	if this.token.Type != lex.TokenAssign {
		err = fmt.Errorf("%s missing assign", this.errMsgPrefix())
		this.errs = append(this.errs, err)
		return
	}
	this.Next(lfNotToken) // skip =
	es, err := this.ExpressionParser.parseExpressions(lex.TokenSemicolon)
	if err != nil {
		return
	}
	if len(es) != len(constants) {
		err = fmt.Errorf("%s cannot assign %d value to %d constant",
			this.errMsgPrefix(), len(es), len(constants))
		this.errs = append(this.errs, err)
	}
	for k, _ := range constants {
		if k < len(es) {
			constants[k].DefaultValueExpression = es[k]
		}
	}
	return
}

// str := "hello world"   a,b = 123 or a b ;
func (this *Parser) parseVar() (ret *ast.ExpressionVar, err error) {
	names, err := this.parseNameList()
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
	if this.token.Type != lex.TokenAssign {
		ret.Type, err = this.parseType()
		if err != nil {
			return
		}
	}
	if this.token.Type == lex.TokenAssign {
		this.Next(lfNotToken) // skip = or :=
		ret.InitValues, err = this.ExpressionParser.parseExpressions(lex.TokenSemicolon)
		if err != nil {
			return
		}
	}
	return
}

func (this *Parser) Next(lfIsToken bool) {
	if this.consumeFoundValidToken {
		this.consumeFoundValidToken = false
		return
	}
	var err error
	var tok *lex.Token
	this.lastToken = this.token
	defer func() {
		if this.lastToken == nil {
			this.lastToken = this.token
		}
	}()
	for {
		tok, err = this.lexer.Next()
		if tok != nil {
			this.token = tok
		}
		if err != nil {
			this.errs = append(this.errs,
				fmt.Errorf("%s %s", this.errMsgPrefix(), err.Error()))
		}
		if tok == nil {
			continue
		}
		this.token = tok
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
	pos.ErrMsgPrefix() only receive one argument
*/
func (this *Parser) errMsgPrefix(pos ...*ast.Pos) string {
	if len(pos) > 0 {
		return pos[0].ErrMsgPrefix()
	}
	return this.mkPos().ErrMsgPrefix()
}

func (this *Parser) consume(until map[lex.TokenKind]bool) {
	if len(until) == 0 {
		panic("no token to consume")
	}
	for this.token.Type != lex.TokenEof {
		if this.token.Type == lex.TokenPublic ||
			this.token.Type == lex.TokenProtected ||
			this.token.Type == lex.TokenPrivate ||
			this.token.Type == lex.TokenClass ||
			this.token.Type == lex.TokenInterface ||
			this.token.Type == lex.TokenFn ||
			this.token.Type == lex.TokenFor ||
			this.token.Type == lex.TokenIf ||
			this.token.Type == lex.TokenSwitch ||
			this.token.Type == lex.TokenEnum ||
			this.token.Type == lex.TokenConst ||
			this.token.Type == lex.TokenVar ||
			this.token.Type == lex.TokenImport ||
			this.token.Type == lex.TokenTypeAlias ||
			this.token.Type == lex.TokenGoto ||
			this.token.Type == lex.TokenBreak ||
			this.token.Type == lex.TokenContinue ||
			this.token.Type == lex.TokenDefer ||
			this.token.Type == lex.TokenReturn ||
			this.token.Type == lex.TokenPass ||
			this.token.Type == lex.TokenExtends ||
			this.token.Type == lex.TokenImplements ||
			this.token.Type == lex.TokenGlobal ||
			this.token.Type == lex.TokenCase ||
			this.token.Type == lex.TokenDefault {
			if _, ok := until[this.token.Type]; ok == false {
				this.consumeFoundValidToken = true
				return
			}
		}
		if this.token.Type == lex.TokenLc {
			if _, ok := until[lex.TokenLc]; ok == false {
				this.consumeFoundValidToken = true
				return
			}
		}
		if this.token.Type == lex.TokenRc {
			if _, ok := until[lex.TokenRc]; ok == false {
				this.consumeFoundValidToken = true
				return
			}
		}
		if _, ok := until[this.token.Type]; ok {
			return
		}
		this.Next(lfIsToken)
	}
}

func (this *Parser) ifTokenIsLfThenSkip() {
	if this.token.Type == lex.TokenLf {
		this.Next(lfNotToken)
	}
}

func (this *Parser) unExpectNewLineAndSkip() {
	if err := this.unExpectNewLine(); err != nil {
		this.Next(lfNotToken)
	}
}
func (this *Parser) unExpectNewLine() error {
	var err error
	if this.token.Type == lex.TokenLf {
		err = fmt.Errorf("%s unexpected new line",
			this.errMsgPrefix(this.mkPos()))
		this.errs = append(this.errs, err)
	}
	return err
}
func (this *Parser) expectNewLineAndSkip() {
	if err := this.expectNewLine(); err == nil {
		this.Next(lfNotToken)
	}
}
func (this *Parser) expectNewLine() error {
	var err error
	if this.token.Type != lex.TokenLf &&
		this.token.Type != lex.TokenComment {
		err = fmt.Errorf("%s expect new line , but '%s'",
			this.errMsgPrefix(), this.token.Description)
		this.errs = append(this.errs, err)
	}
	return err
}

func (this *Parser) parseTypeAlias(comment *CommentParser) (*ast.TypeAlias, error) {
	this.Next(lfIsToken) // skip type key word
	this.unExpectNewLineAndSkip()
	if this.token.Type != lex.TokenIdentifier {
		err := fmt.Errorf("%s expect identifer,but '%s'", this.errMsgPrefix(), this.token.Description)
		this.errs = append(this.errs, err)
		return nil, err
	}
	ret := &ast.TypeAlias{}
	ret.Pos = this.mkPos()
	ret.Name = this.token.Data.(string)
	this.Next(lfIsToken) // skip identifier
	if this.token.Type != lex.TokenAssign {
		err := fmt.Errorf("%s expect '=',but '%s'", this.errMsgPrefix(), this.token.Description)
		this.errs = append(this.errs, err)
		return nil, err
	}
	this.Next(lfNotToken) // skip =
	var err error
	ret.Type, err = this.parseType()
	if err != nil {
		return nil, err
	}
	ret.Comment = comment.Comment
	if this.token.Type == lex.TokenComment {
		this.Next(lfIsToken)
	}
	return ret, err
}

/*
	a int
	int
*/
func (this *Parser) parseTypedName() (vs []*ast.Variable, err error) {
	if this.token.Type != lex.TokenIdentifier {
		/*
			not identifier begin
			must be type
			// int
		*/
		t, err := this.parseType()
		if err != nil {
			return nil, err
		}
		v := &ast.Variable{}
		v.Type = t
		v.Pos = this.mkPos()
		return []*ast.Variable{v}, nil
	}
	names, err := this.parseNameList()
	if err != nil {
		return nil, err
	}
	if this.isValidTypeBegin() {
		/*
			a , b int
		*/
		t, err := this.parseType()
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
			not valid type after name list, "a" and "b" must indicate types
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
func (this *Parser) parseTypedNames() (vs []*ast.Variable, err error) {
	vs = []*ast.Variable{}
	for this.token.Type != lex.TokenEof {
		ns, err := this.parseNameList()
		if err != nil {
			return vs, err
		}
		t, err := this.parseType()
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
		if this.token.Type != lex.TokenComma { // not a comma
			break
		} else {
			this.Next(lfNotToken)
		}
	}
	return vs, nil
}

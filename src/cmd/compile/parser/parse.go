package parser

import (
	"bytes"
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func Parse(tops *[]*ast.Top, filename string, bs []byte, onlyImport bool, nerr int) []error {
	p := &Parser{
		bs:           bs,
		tops:         tops,
		filename:     filename,
		onlyImport:   onlyImport,
		nErrors2Stop: nerr,
	}
	return p.Parse()
}

type Parser struct {
	onlyImport   bool
	bs           []byte
	lines        [][]byte
	tops         *[]*ast.Top
	scanner      *lex.Lexer
	filename     string
	lastToken    *lex.Token
	token        *lex.Token
	errs         []error
	imports      map[string]*ast.Import
	nErrors2Stop int
	// parsers
	ExpressionParser *ExpressionParser
	FunctionParser   *FunctionParser
	ClassParser      *ClassParser
	BlockParser      *BlockParser
	InterfaceParser  *InterfaceParser
}

func (p *Parser) Parse() []error {
	p.ExpressionParser = &ExpressionParser{p}
	p.FunctionParser = &FunctionParser{}
	p.FunctionParser.parser = p
	p.ClassParser = &ClassParser{}
	p.ClassParser.parser = p
	p.InterfaceParser = &InterfaceParser{}
	p.InterfaceParser.parser = p
	p.BlockParser = &BlockParser{}
	p.BlockParser.parser = p
	p.errs = []error{}
	p.scanner = lex.New(p.bs, 1, 1)
	p.lines = bytes.Split(p.bs, []byte("\n"))
	p.Next()
	if p.token.Type == lex.TOKEN_EOF {
		return nil
	}
	p.parseImports() // next is called
	if p.token.Type == lex.TOKEN_EOF {
		return p.errs
	}
	if p.onlyImport { // only parse imports
		return p.errs
	}
	isPublic := false
	resetProperty := func() {
		isPublic = false
	}
	for p.token.Type != lex.TOKEN_EOF {
		if len(p.errs) > p.nErrors2Stop {
			break
		}
		switch p.token.Type {
		case lex.TOKEN_SEMICOLON: // empty statement, no big deal
			p.Next()
			continue
		case lex.TOKEN_VAR:
			pos := p.mkPos()
			p.Next() // skip var key word
			vs, es, typ, err := p.parseConstDefinition(true)
			if err != nil {
				p.consume(untilSemicolon)
				p.Next()
				continue
			}
			if typ != nil && typ.Type != lex.TOKEN_ASSIGN {
				p.errs = append(p.errs,
					fmt.Errorf("%s use '=' to initialize value",
						p.errorMsgPrefix()))
			}
			d := &ast.ExpressionDeclareVariable{Variables: vs, InitValues: es}
			e := &ast.Expression{
				Type:     ast.EXPRESSION_TYPE_VAR,
				Data:     d,
				Pos:      pos,
				IsPublic: isPublic,
			}
			*p.tops = append(*p.tops, &ast.Top{
				Data: e,
			})
			resetProperty()
		case lex.TOKEN_IDENTIFIER:
			e, err := p.ExpressionParser.parseExpression(true)
			if err != nil {
				p.consume(untilSemicolon)
				p.Next()
				continue
			}
			e.IsPublic = isPublic
			p.validStatementEnding(e.Pos)
			*p.tops = append(*p.tops, &ast.Top{
				Data: e,
			})
			resetProperty()
		case lex.TOKEN_ENUM:
			e, err := p.parseEnum(isPublic)
			if err != nil {
				p.consume(untilRc)
				p.Next()
				resetProperty()
				continue
			}
			if e != nil {
				*p.tops = append(*p.tops, &ast.Top{
					Data: e,
				})
			}
			resetProperty()
		case lex.TOKEN_FUNCTION:
			f, err := p.FunctionParser.parse(true)
			if err != nil {
				p.consume(untilRc)
				p.Next()
				continue
			}
			if isPublic {
				f.AccessFlags |= cg.ACC_METHOD_PUBLIC
			} else {
				f.AccessFlags |= cg.ACC_METHOD_PRIVATE
			}
			*p.tops = append(*p.tops, &ast.Top{
				Data: f,
			})
			resetProperty()
		case lex.TOKEN_LC:
			b := &ast.Block{}
			p.Next()                                  // skip {
			p.BlockParser.parseStatementList(b, true) // this function will lookup next
			if p.token.Type != lex.TOKEN_RC {
				p.errs = append(p.errs, fmt.Errorf("%s expect '}', but '%s'",
					p.errorMsgPrefix(), p.token.Description))
				p.consume(untilRc)
			}
			p.Next() // skip }
			*p.tops = append(*p.tops, &ast.Top{
				Data: b,
			})
			resetProperty()
		case lex.TOKEN_CLASS:
			c, err := p.ClassParser.parse()
			if err != nil {
				p.errs = append(p.errs, err)
				p.consume(untilRc)
				p.Next()
				resetProperty()
				continue
			}
			*p.tops = append(*p.tops, &ast.Top{
				Data: c,
			})
			if isPublic {
				c.AccessFlags |= cg.ACC_CLASS_PUBLIC
			}
			resetProperty()
		case lex.TOKEN_INTERFACE:
			c, err := p.InterfaceParser.parse()
			if err != nil {
				p.errs = append(p.errs, err)
				p.consume(untilRc)
				p.Next()
				resetProperty()
				continue
			}
			*p.tops = append(*p.tops, &ast.Top{
				Data: c,
			})
			if isPublic {
				c.AccessFlags |= cg.ACC_CLASS_PUBLIC
			}
			resetProperty()
		case lex.TOKEN_PUBLIC:
			isPublic = true
			p.Next()
			p.validAfterPublic(isPublic)
			continue
		case lex.TOKEN_CONST:
			p.Next() // skip const key word
			vs, es, typ, err := p.parseConstDefinition(false)
			if err != nil {
				p.consume(untilSemicolon)
				p.Next()
				resetProperty()
				continue
			}
			if p.validStatementEnding() == false { //assume missing ; not big deal
				p.Next()
				p.consume(untilSemicolon)
				resetProperty()
				continue
			}
			// const a := 1 is wrong,
			if typ != nil && typ.Type != lex.TOKEN_ASSIGN {
				p.errs = append(p.errs, fmt.Errorf("%s use '=' instead of ':=' for const definition",
					p.errorMsgPrefix()))
				resetProperty()
				continue
			}
			if len(vs) != len(es) {
				p.errs = append(p.errs,
					fmt.Errorf("%s cannot assign %d values to %d destinations",
						p.errorMsgPrefix(p.mkPos()), len(es), len(vs)))
			}
			for k, v := range vs {
				if k < len(es) {
					c := &ast.Constant{}
					c.Variable = *v
					c.Expression = es[k]
					if isPublic {
						c.AccessFlags |= cg.ACC_FIELD_PUBLIC
					} else {
						c.AccessFlags |= cg.ACC_FIELD_PRIVATE
					}
					*p.tops = append(*p.tops, &ast.Top{
						Data: c,
					})
				}
			}
			resetProperty()
			continue
		case lex.TOKEN_PRIVATE: //is a default attribute
			isPublic = false
			p.Next()
			p.validAfterPublic(isPublic)
			continue
		case lex.TOKEN_TYPE:
			a, err := p.parseTypeAlias()
			if err != nil {
				p.consume(untilSemicolon)
				p.Next()
				resetProperty()
				continue
			}
			*p.tops = append(*p.tops, &ast.Top{
				Data: a,
			})

		case lex.TOKEN_EOF:
			break
		default:
			p.errs = append(p.errs, fmt.Errorf("%s token(%s) is not except",
				p.errorMsgPrefix(), p.token.Description))
			p.consume(untilSemicolon)
			resetProperty()
		}
	}
	return p.errs
}

func (p *Parser) parseTypes() ([]*ast.Type, error) {
	ret := []*ast.Type{}
	for p.token.Type != lex.TOKEN_EOF {
		t, err := p.parseType()
		if err != nil {
			return ret, err
		}
		ret = append(ret, t)
		if p.token.Type != lex.TOKEN_COMMA {
			break
		}
		p.Next() // skip ,
	}
	return ret, nil
}

func (p *Parser) validAfterPublic(isPublic bool) {
	if p.token.Type == lex.TOKEN_FUNCTION ||
		p.token.Type == lex.TOKEN_CLASS ||
		p.token.Type == lex.TOKEN_ENUM ||
		p.token.Type == lex.TOKEN_IDENTIFIER ||
		p.token.Type == lex.TOKEN_INTERFACE ||
		p.token.Type == lex.TOKEN_CONST ||
		p.token.Type == lex.TOKEN_VAR {
		return
	}
	var err error
	token := "public"
	if isPublic == false {
		token = "private"
	}
	if p.token.Description != "" {
		err = fmt.Errorf("%s cannot have token:%s after '%s'",
			p.errorMsgPrefix(), p.token.Description, token)
	} else {
		err = fmt.Errorf("%s cannot have token:%s after '%s'",
			p.errorMsgPrefix(), p.token.Description, token)
	}
	p.errs = append(p.errs, err)
}
func (p *Parser) validStatementEnding(pos ...*ast.Position) bool {
	if p.token.Type == lex.TOKEN_SEMICOLON ||
		(p.lastToken != nil && p.lastToken.Type == lex.TOKEN_RC) {
		return true
	}
	if len(pos) > 0 {
		p.errs = append(p.errs, fmt.Errorf("%s missing semicolon", p.errorMsgPrefix(pos[0])))
	} else {
		p.errs = append(p.errs, fmt.Errorf("%s missing semicolon", p.errorMsgPrefix()))
	}
	return false
}

func (p *Parser) mkPos() *ast.Position {
	return &ast.Position{
		Filename:    p.filename,
		StartLine:   p.token.StartLine,
		StartColumn: p.token.StartColumn,
		Offset:      p.scanner.GetOffSet(),
	}
}

// str := "hello world"   a,b = 123 or a b ;
func (p *Parser) parseConstDefinition(needType bool) ([]*ast.Variable, []*ast.Expression, *lex.Token, error) {
	names, err := p.parseNameList()
	if err != nil {
		return nil, nil, nil, err
	}
	var variableType *ast.Type
	//trying to parse type
	if p.isValidTypeBegin() || needType {
		variableType, err = p.parseType()
		if err != nil {
			p.errs = append(p.errs, err)
			return nil, nil, nil, err
		}
	}
	f := func() []*ast.Variable {
		vs := make([]*ast.Variable, len(names))
		for k, v := range names {
			vd := &ast.Variable{}
			vd.Name = v.Name
			vd.Pos = v.Pos
			if variableType != nil {
				vd.Type = variableType.Clone()
			}
			vs[k] = vd
		}
		return vs
	}
	if p.token.Type != lex.TOKEN_ASSIGN &&
		p.token.Type != lex.TOKEN_COLON_ASSIGN {
		return f(), nil, nil, err
	}
	typ := p.token
	p.Next() // skip = or :=
	es, err := p.ExpressionParser.parseExpressions()
	if err != nil {
		return nil, nil, typ, err
	}
	return f(), es, typ, nil
}

func (p *Parser) Next() {
	var err error
	var tok *lex.Token
	p.lastToken = p.token
	for {
		tok, err = p.scanner.Next()
		if err != nil {
			p.errs = append(p.errs, fmt.Errorf("%s %s", p.errorMsgPrefix(), err.Error()))
		}
		if tok == nil {
			continue
		}
		p.token = tok
		if tok.Type != lex.TOKEN_LF {
			if p.token.Description != "" {
				//	fmt.Println("#########", p.token.Type, p.token.Desp)
			} else {
				//fmt.Println("#########", p.token.Type, p.token.Data)
			}
			break
		}
	}
	return
}

/*
	errorMsgPrefix(pos) only receive one argument
*/
func (p *Parser) errorMsgPrefix(pos ...*ast.Position) string {
	if len(pos) > 0 {
		return fmt.Sprintf("%s:%d:%d", pos[0].Filename, pos[0].StartLine, pos[0].StartColumn)
	}
	line, column := p.scanner.GetPos()
	return fmt.Sprintf("%s:%d:%d", p.filename, line, column)
}

func (p *Parser) consume(until map[int]bool) {
	if len(until) == 0 {
		panic("no token to consume")
	}
	var ok bool
	for p.token.Type != lex.TOKEN_EOF {
		if _, ok = until[p.token.Type]; ok {
			return
		}
		p.Next()
	}
}

func (p *Parser) lexPos2AstPos(t *lex.Token, pos *ast.Position) {
	pos.Filename = p.filename
	pos.StartLine = t.StartLine
	pos.StartColumn = t.StartColumn
}

func (p *Parser) parseTypeAlias() (*ast.ExpressionTypeAlias, error) {
	p.Next() // skip type key word
	if p.token.Type != lex.TOKEN_IDENTIFIER {
		err := fmt.Errorf("%s expect identifer,but %s", p.errorMsgPrefix(), p.token.Description)
		p.errs = append(p.errs, err)
		return nil, err
	}
	ret := &ast.ExpressionTypeAlias{}
	ret.Pos = p.mkPos()
	ret.Name = p.token.Data.(string)
	p.Next() // skip identifier
	if p.token.Type != lex.TOKEN_ASSIGN {
		err := fmt.Errorf("%s expect '=',but %s", p.errorMsgPrefix(), p.token.Description)
		p.errs = append(p.errs, err)
		return nil, err
	}
	p.Next() // skip =
	var err error
	ret.Type, err = p.parseType()
	return ret, err
}

func (p *Parser) parseTypedName() (vs []*ast.Variable, err error) {
	names, err := p.parseNameList()
	if err != nil {
		return nil, err
	}
	t, err := p.parseType()
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
	}
	return vs, nil
}

// a,b int or int,bool  c xxx
func (p *Parser) parseTypedNames() (vs []*ast.Variable, err error) {
	vs = []*ast.Variable{}
	for p.token.Type != lex.TOKEN_EOF {
		ns, err := p.parseNameList()
		if err != nil {
			return vs, err
		}
		t, err := p.parseType()
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
		if p.token.Type != lex.TOKEN_COMMA { // not a comma
			break
		} else {
			p.Next()
		}
	}
	return vs, nil
}

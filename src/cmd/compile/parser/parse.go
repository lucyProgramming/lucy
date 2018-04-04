package parser

import (
	"bytes"
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func Parse(tops *[]*ast.Node, filename string, bs []byte, onlyimport bool, nerr int) []error {
	return (&Parser{bs: bs, tops: tops, filename: filename, onlyimport: onlyimport, nerr: nerr}).Parse()
}

type Parser struct {
	onlyimport       bool
	bs               []byte
	lines            [][]byte
	tops             *[]*ast.Node
	ExpressionParser *ExpressionParser
	Function         *Function
	Class            *Class
	Block            *Block
	scanner          *lex.LucyLexer
	filename         string
	lastToken        *lex.Token
	token            *lex.Token
	eof              bool
	errs             []error
	imports          map[string]*ast.Import
	nerr             int
}

func (p *Parser) Parse() []error {
	p.ExpressionParser = &ExpressionParser{p}
	p.Function = &Function{}
	p.Function.parser = p
	p.Class = &Class{}
	p.Class.parser = p
	p.Block = &Block{}
	p.Block.parser = p
	p.errs = []error{}
	p.scanner = lex.New(p.bs)
	p.lines = bytes.Split(p.bs, []byte("\n"))
	p.Next()
	if p.eof {
		return nil
	}
	p.parseImports() // next is called
	if p.eof {
		return p.errs
	}
	if p.onlyimport { // only parse imports
		return p.errs
	}
	ispublic := false
	resetProperty := func() {
		ispublic = false
	}
	var err error
	for !p.eof {
		if len(p.errs) > p.nerr {
			break
		}
		switch p.token.Type {
		case lex.TOKEN_SEMICOLON: // empty statement, no big deal
			p.Next()
			continue
		case lex.TOKEN_VAR:
			pos := p.mkPos()
			p.Next() // skip var key word
			vs, es, _, err := p.parseConstDefinition()
			if err != nil {
				p.consume(untils_semicolon)
				p.Next()
				continue
			}
			d := &ast.ExpressionDeclareVariable{Vs: vs, Values: es}
			e := &ast.Expression{
				Typ:  ast.EXPRESSION_TYPE_VAR,
				Data: d,
				Pos:  pos,
			}
			*p.tops = append(*p.tops, &ast.Node{
				Data: e,
			})
			resetProperty()
		case lex.TOKEN_IDENTIFIER:
			e, err := p.ExpressionParser.parseExpression()
			if err != nil {
				p.consume(untils_semicolon)
				p.Next()
				continue
			}
			p.validStatementEnding(e.Pos)
			*p.tops = append(*p.tops, &ast.Node{
				Data: e,
			})
			resetProperty()
		case lex.TOKEN_ENUM:
			e, err := p.parseEnum(ispublic)
			if err != nil {
				p.consume(untils_rc)
				p.Next()
				resetProperty()
				continue
			}
			if e != nil {
				*p.tops = append(*p.tops, &ast.Node{
					Data: e,
				})
			}
			resetProperty()
		case lex.TOKEN_FUNCTION:
			f, err := p.Function.parse(ispublic)
			if err != nil {
				p.errs = append(p.errs, err)
				p.consume(untils_rc)
				p.Next()
				continue
			}
			*p.tops = append(*p.tops, &ast.Node{
				Data: f,
			})
			resetProperty()
		case lex.TOKEN_LC:
			b := &ast.Block{}
			p.Next()
			err = p.Block.parse(b, false, lex.TOKEN_RC) // this function will lookup next
			if err != nil {
				p.consume(untils_rc)
				p.Next()
			}
			*p.tops = append(*p.tops, &ast.Node{
				Data: b,
			})
			resetProperty()
		case lex.TOKEN_CLASS:
			c, err := p.Class.parse()
			if err != nil {
				p.errs = append(p.errs, err)
				p.consume(untils_rc)
				p.Next()
				resetProperty()
				continue
			}
			*p.tops = append(*p.tops, &ast.Node{
				Data: c,
			})
			if ispublic {
				c.Access |= cg.ACC_FIELD_PUBLIC
			} else {
				c.Access |= cg.ACC_FIELD_PRIVATE
			}
			resetProperty()
		case lex.TOKEN_PUBLIC:
			ispublic = true
			p.Next()
			p.validAfterPublic()
			continue
		case lex.TOKEN_CONST:
			p.Next() // skip const key word
			vs, es, typ, err := p.parseConstDefinition()
			if err != nil {
				p.consume(untils_semicolon)
				p.Next()
				resetProperty()
				continue
			}
			if p.token.Type != lex.TOKEN_SEMICOLON && (p.lastToken != nil && p.lastToken.Type != lex.TOKEN_RC) { //assume missing ; not big deal
				p.errs = append(p.errs, fmt.Errorf("%s not ; after variable or const definition,but %s", p.errorMsgPrefix(), p.token.Desp))
				p.Next()
				p.consume(untils_semicolon)
				resetProperty()
				continue
			}
			// const a := 1 is wrong,
			if typ == lex.TOKEN_COLON_ASSIGN {
				p.errs = append(p.errs, fmt.Errorf("%s use = instead of := for const definition", p.errorMsgPrefix()))
				resetProperty()
				continue
			}
			for k, v := range vs {
				if k < len(es) {
					c := &ast.Const{}
					c.VariableDefinition = *v
					c.Expression = es[k]
					if ispublic {
						c.AccessFlags |= cg.ACC_FIELD_PUBLIC
					} else {
						c.AccessFlags |= cg.ACC_FIELD_PRIVATE
					}
					*p.tops = append(*p.tops, &ast.Node{
						Data: c,
					})
				}
			}
			resetProperty()
			continue
		case lex.TOKEN_PRIVATE: //is a default attribute
			ispublic = false
			p.Next()
			p.validAfterPublic()
			continue
		case lex.TOKEN_TYPE:
			a, err := p.parseTypeaAlias()
			if err != nil {
				p.consume(untils_semicolon)
				p.Next()
				resetProperty()
				continue
			}
			*p.tops = append(*p.tops, &ast.Node{
				Data: a,
			})
		default:
			p.errs = append(p.errs, fmt.Errorf("%s token(%s) is not except", p.errorMsgPrefix(), p.token.Desp))
			p.consume(untils_semicolon)
			resetProperty()
		}
	}
	return p.errs
}

func (p *Parser) validAfterPublic() {
	if p.token.Type == lex.TOKEN_FUNCTION || p.token.Type == lex.TOKEN_CLASS || p.token.Type == lex.TOKEN_ENUM || p.token.Type == lex.TOKEN_IDENTIFIER {
		return
	}
	var err error
	if p.token.Desp != "" {
		err = fmt.Errorf("%s cannot have token:'%s' after 'public' or 'private'", p.errorMsgPrefix(), p.token.Desp)
	} else {
		err = fmt.Errorf("%s cannot have token:'%v' after 'public' or 'private'", p.errorMsgPrefix(), p.token.Data)
	}
	p.errs = append(p.errs, err)
}
func (p *Parser) validStatementEnding(pos ...*ast.Pos) {
	if p.token.Type == lex.TOKEN_SEMICOLON || p.lastToken != nil && p.lastToken.Type == lex.TOKEN_RC {
		return
	}
	if len(pos) > 0 {
		p.errs = append(p.errs, fmt.Errorf("%s missing semicolon", p.errorMsgPrefix(pos[0])))
	} else {
		p.errs = append(p.errs, fmt.Errorf("%s missing semicolon", p.errorMsgPrefix()))
	}

}

func (p *Parser) mkPos() *ast.Pos {
	return &ast.Pos{
		Filename:    p.filename,
		StartLine:   p.token.StartLine,
		StartColumn: p.token.StartColumn,
	}
}

// str := "hello world"   a,b = 123 or a b ;
func (p *Parser) parseConstDefinition() ([]*ast.VariableDefinition, []*ast.Expression, int, error) {
	names, err := p.parseNameList()
	if err != nil {
		return nil, nil, 0, err
	}
	var variableType *ast.VariableType
	//trying to parse type
	if p.isValidTypeBegin() {
		variableType, err = p.parseType()
		if err != nil {
			p.errs = append(p.errs, err)
			return nil, nil, 0, err
		}
	}
	f := func() []*ast.VariableDefinition {
		vs := make([]*ast.VariableDefinition, len(names))
		for k, v := range names {
			vd := &ast.VariableDefinition{}
			vd.Name = v.Name
			vd.Pos = v.Pos
			if variableType != nil {
				vd.Typ = variableType.Clone()
			}
			vs[k] = vd
		}
		return vs
	}
	if p.token.Type != lex.TOKEN_ASSIGN && p.token.Type != lex.TOKEN_COLON_ASSIGN {
		return f(), nil, 0, err
	}
	typ := p.token.Type
	p.Next() // skip = or :=
	if p.eof {
		err = p.mkUnexpectedEofErr()
		p.errs = append(p.errs, err)
		return nil, nil, typ, err
	}
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
	for !p.eof {
		tok, p.eof, err = p.scanner.Next()
		if err != nil {
			p.errs = append(p.errs, fmt.Errorf("%s %s", p.errorMsgPrefix(), err.Error()))
		}
		if p.eof {
			break
		}
		if tok == nil {
			continue
		}
		if tok.Type != lex.TOKEN_CRLF {
			p.token = tok
			if p.token.Desp != "" {
				//fmt.Println("#########", p.token.Desp)
			} else {
				//fmt.Println("#########", p.token.Data)
			}
			break
		}
	}
	return
}

/*
	errorMsgPrefix(pos) only receive one argument
*/
func (p *Parser) errorMsgPrefix(pos ...*ast.Pos) string {
	if len(pos) > 0 {
		return fmt.Sprintf("%s:%d:%d ", pos[0].Filename, pos[0].StartLine, pos[0].StartColumn)
	}
	line, column := p.scanner.Pos()
	return fmt.Sprintf("%s:%d:%d ", p.filename, line, column)
}

func (p *Parser) consume(untils map[int]bool) {
	if len(untils) == 0 {
		panic("no token to consume")
	}
	var ok bool
	for p.eof == false {
		if _, ok = untils[p.token.Type]; ok {
			return
		}
		p.Next()
	}
}

func (p *Parser) lexPos2AstPos(t *lex.Token, pos *ast.Pos) {
	pos.Filename = p.filename
	pos.StartLine = t.StartLine
	pos.StartColumn = t.StartColumn
}

func (p *Parser) parseTypeaAlias() (*ast.ExpressionTypeAlias, error) {
	p.Next() // skip type key word
	if p.token.Type != lex.TOKEN_IDENTIFIER {
		err := fmt.Errorf("%s expect identifer,but %s", p.errorMsgPrefix(), p.token.Desp)
		p.errs = append(p.errs, err)
		return nil, err
	}
	ret := &ast.ExpressionTypeAlias{}
	ret.Pos = p.mkPos()
	ret.Name = p.token.Data.(string)
	p.Next() // skip identifier
	if p.token.Type != lex.TOKEN_ASSIGN {
		err := fmt.Errorf("%s expect '=',but %s", p.errorMsgPrefix(), p.token.Desp)
		p.errs = append(p.errs, err)
		return nil, err
	}
	p.Next() // skip =
	var err error
	ret.Typ, err = p.parseType()
	return ret, err
}

func (p *Parser) parseTypedName() (vs []*ast.VariableDefinition, err error) {
	names, err := p.parseNameList()
	if err != nil {
		return nil, err
	}
	t, err := p.parseType()
	if err != nil {
		return nil, err
	}
	vs = make([]*ast.VariableDefinition, len(names))
	for k, v := range names {
		vd := &ast.VariableDefinition{}
		vs[k] = vd
		vd.Name = v.Name
		vd.Pos = v.Pos
		vd.Typ = t.Clone()
	}
	return vs, nil
}

// a,b int or int,bool  c xxx
func (p *Parser) parseTypedNames() (vs []*ast.VariableDefinition, err error) {
	vs = []*ast.VariableDefinition{}
	for !p.eof {
		ns, err := p.parseNameList()
		if err != nil {
			return vs, err
		}
		t, err := p.parseType()
		if err != nil {
			return vs, err
		}
		for _, v := range ns {
			vd := &ast.VariableDefinition{}
			vd.Name = v.Name
			vd.Pos = v.Pos
			vd.Typ = t.Clone()
			vs = append(vs, vd)
		}
		if p.token.Type != lex.TOKEN_COMMA { // not a commna
			break
		} else {
			p.Next()
		}
	}
	return vs, nil
}

////var a,b,c int,char,bool  | var a,b,c int = 123;
//func (p *Parser) parseVarDefinition(ispublic ...bool) (vs []*ast.VariableDefinition, expressions []*ast.Expression, err error) {
//	p.Next()
//	if p.eof {
//		err = p.mkUnexpectedEofErr()
//		p.errs = append(p.errs, err)
//		return
//	}
//	names, err := p.parseNameList()
//	if err != nil {
//		return nil, nil, err
//	}
//	if p.eof {
//		err = p.mkUnexpectedEofErr()
//		p.errs = append(p.errs, err)
//		return
//	}
//	t, err := p.parseType()
//	if t == nil {
//		err = fmt.Errorf("%s no variable type found or defined wrong", p.errorMsgPrefix())
//		p.errs = append(p.errs, err)
//		return nil, nil, err
//	}
//	//value , no default value definition
//	if lex.TOKEN_ASSIGN == p.token.Type {
//		//assign
//		p.Next() // skip =
//		expressions, err = p.ExpressionParser.parseExpressions()
//		if err != nil {
//			p.errs = append(p.errs, err)
//		}
//	}
//	if p.token.Type != lex.TOKEN_SEMICOLON {
//		err = fmt.Errorf("%s not a \";\" after a variable declare ", p.errorMsgPrefix())
//		p.errs = append(p.errs, err)
//		return
//	}
//	p.Next() // look next
//	vs = make([]*ast.VariableDefinition, len(names))
//	for k, v := range names {
//		vd := &ast.VariableDefinition{}
//		vd.Name = v.Name
//		vd.Typ = t.Clone()
//		if len(ispublic) > 0 && ispublic[0] {
//			vd.AccessFlags |= cg.ACC_FIELD_PUBLIC
//		} else {
//			vd.AccessFlags |= cg.ACC_FIELD_PRIVATE
//		}
//		vd.Pos = v.Pos
//		vs[k] = vd
//	}
//	return
//}

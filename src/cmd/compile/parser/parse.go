package parser

import (
	"bytes"
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
	"github.com/timtadh/lexmachine"
)

func Parse(tops *[]*ast.Node, filename string, bs []byte, onlyimport bool) []error {
	return (&Parser{bs: bs, tops: tops, filename: filename, onlyimport: onlyimport}).Parse()
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
	scanner          *lexmachine.Scanner
	filename         string
	lastToken        *lex.Token
	token            *lex.Token
	eof              bool
	errs             []error
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
	var err error
	p.scanner, err = lex.Lexer.Scanner(p.bs)
	p.lines = bytes.Split(p.bs, []byte("\n"))
	if err != nil {
		p.errs = append(p.errs, err)
		return p.errs
	}
	p.Next()
	//package name definition
	if p.eof {
		p.errs = append(p.errs, fmt.Errorf("no package name definition found"))
		return p.errs
	}
	if p.token.Type != lex.TOKEN_PACKAGE {
		p.errs = append(p.errs, fmt.Errorf("first token must be a  package name definition"))
		return p.errs
	}
	p.Next()
	if p.eof {
		p.errs = append(p.errs, fmt.Errorf("no package name definition found(no name after)"))
		return p.errs
	}
	if p.token.Type != lex.TOKEN_IDENTIFIER {
		p.errs = append(p.errs, fmt.Errorf("no package name definition found(no name after)"))
		return p.errs
	}
	pd := &ast.PackageNameDeclare{
		Name: p.token.Data.(string),
	}
	p.lexPos2AstPos(p.token, &pd.Pos)
	*p.tops = append(*p.tops, &ast.Node{
		Data: pd,
	})
	p.parseImports() // next is called
	if p.eof {
		return p.errs
	}
	ispublic := false
	isconst := false
	resetProperty := func() {
		ispublic = false
		isconst = false
	}
	for !p.eof {
		switch p.token.Type {
		case lex.TOKEN_SEMICOLON: // empty statement, no big deal
			p.Next()
			continue
		case lex.TOKEN_VAR:
			vs, err := p.parseVarDefinition(ispublic)
			if err != nil {
				p.consume(untils_semicolon)
				p.Next()
				continue
			}
			if vs != nil && len(vs) > 0 {
				for _, v := range vs {
					*p.tops = append(*p.tops, &ast.Node{
						Data: v,
					})
				}
			}
			resetProperty()
		case lex.TOKEN_IDENTIFIER:
			names, typ, es, err := p.parseAssignedNames()
			if err != nil {
				p.errs = append(p.errs, err)
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
			if typ != lex.TOKEN_COLON_ASSIGN {
				p.errs = append(p.errs, fmt.Errorf("%s cannot have statement at top,possible do you mean a := 1 to create a global variable", p.errorMsgPrefix()))
				p.Next()
				p.consume(untils_semicolon)
				resetProperty()
				continue
			}
			for k, v := range names {
				c := &ast.Const{}
				c.Name = v.Name
				c.Expression = es[k]
				if ispublic {
					c.Access = ast.ACCESS_PUBLIC
				} else {
					c.Access = ast.ACCESS_PRIVATE
				}
				*p.tops = append(*p.tops, &ast.Node{
					Data: c,
				})
			}
			resetProperty()
		case lex.TOKEN_ENUM:
			if isconst {
				p.errs = append(p.errs, fmt.Errorf("%s cannot use const for a enum", p.errorMsgPrefix()))
			}
			e := p.parseEnum(ispublic)
			if e != nil {
				*p.tops = append(*p.tops, &ast.Node{
					Data: e,
				})
			}
			resetProperty()
		case lex.TOKEN_FUNCTION:
			if isconst {
				p.errs = append(p.errs, fmt.Errorf("%s cannot use const for a enum", p.errorMsgPrefix()))
			}
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
			if isconst {
				p.errs = append(p.errs, fmt.Errorf("%s cannot use const for a block", p.errorMsgPrefix()))
			}
			if ispublic {
				p.errs = append(p.errs, fmt.Errorf("%s cannot use public for a block", p.errorMsgPrefix()))
			}
			b := &ast.Block{}
			err = p.Block.parse(b)
			if err != nil {
				p.consume(untils_rc)
			}
			*p.tops = append(*p.tops, &ast.Node{
				Data: b,
			})
			p.Next()
			resetProperty()
		case lex.TOKEN_CLASS:
			if isconst {
				p.errs = append(p.errs, fmt.Errorf("%s can`t use const for a class", p.errorMsgPrefix()))
			}
			c, err := p.Class.parse(ispublic)
			if err != nil {
				p.errs = append(p.errs, err)
				continue
			}
			*p.tops = append(*p.tops, &ast.Node{
				Data: c,
			})
			resetProperty()
		case lex.TOKEN_PUBLIC:
			ispublic = true
			p.Next()
			continue
		case lex.TOKEN_CONST:
			isconst = true
			p.Next()
			continue
		case lex.TOKEN_PRIVATE: //is a default attribute
			ispublic = false
			p.Next()
			continue
		default:
			p.errs = append(p.errs, fmt.Errorf("%s %d:%d token(%s) is not except", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn, p.token.Desp))
			p.consume(untils_semicolon)
			p.Next()
			resetProperty()
		}
	}
	return p.errs
}

func (p *Parser) mkPos() *ast.Pos {
	return &ast.Pos{
		Filename:    p.filename,
		StartLine:   p.token.Match.StartLine,
		StartColumn: p.token.Match.StartColumn,
	}
}

// str := "hello world"   a,b = 123 or a b ;
func (p *Parser) parseAssignedNames() ([]*ast.NameWithPos, int, []*ast.Expression, error) {
	names, err := p.parseNameList()
	if err != nil {
		return nil, 0, nil, err
	}
	if p.token.Type != lex.TOKEN_ASSIGN && p.token.Type != lex.TOKEN_COLON_ASSIGN {
		return nil, 0, nil, err
	}
	typ := p.token.Type
	p.Next()
	if p.eof {
		return nil, 0, nil, fmt.Errorf("%s %d:%d eof after const", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn)
	}
	es, err := p.ExpressionParser.parseExpressions()
	if err != nil {
		return nil, 0, nil, nil
	}
	if len(es) != len(names) {
		return nil, 0, nil, fmt.Errorf("%s %d:%d mame and value not match", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn)
	}
	return names, typ, es, nil
}

func (p *Parser) Next() {
	var err error
	var tok interface{}
	p.lastToken = p.token
	for p.eof == false {
		tok, err, p.eof = p.scanner.Next()
		if err != nil {
			p.eof = true
			return
		}
		if tok != nil && tok.(*lex.Token).Type != lex.TOKEN_CRLF {
			p.token = tok.(*lex.Token)
			fmt.Println("#########", p.token.Desp)
			break
		}
	}
	return
}

func (p *Parser) unexpectedErr() {
	p.errs = append(p.errs, p.mkUnexpectedEofErr())
}
func (p *Parser) mkUnexpectedEofErr() error {
	return fmt.Errorf("%s %d:%d unexpected EOF", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn)
}

func (p *Parser) errorMsgPrefix(pos ...*ast.Pos) string {
	if len(pos) > 0 {
		return fmt.Sprintf("%s %d:%d '%s'", pos[0].Filename, pos[0].StartLine, pos[0].StartColumn, string(p.lines[pos[0].StartLine-1]))
	}
	return fmt.Sprintf("%s %d:%d '%s'", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn, string(p.lines[p.token.Match.StartLine-1]))
}

//var a,b,c int,char,bool  | var a,b,c int = 123;
func (p *Parser) parseVarDefinition(ispublic ...bool) (vs []*ast.VariableDefinition, err error) {
	p.Next()
	if p.eof {
		p.unexpectedErr()
		return
	}
	names, err := p.parseNameList()
	if err != nil {
		return nil, err
	}
	if p.eof {
		p.unexpectedErr()
		return
	}
	if len(names) == 0 {
		err = fmt.Errorf("%s %d:%d no variable name defined", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn)
		p.errs = append(p.errs, err)
		return nil, err
	}
	t, err := p.parseType()
	if t == nil {
		err = fmt.Errorf("%s %d:%d no variable type found or defined wrong", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn)
		p.errs = append(p.errs, err)
		return nil, err
	}
	var expressions []*ast.Expression
	//value , no default value definition
	noDefaultValue := false
	if p.token.Type == lex.TOKEN_SEMICOLON {
		noDefaultValue = true
	} else if lex.TOKEN_ASSIGN == p.token.Type { //assign
		p.Next()
		expressions, err = p.ExpressionParser.parseExpressions()
		if err != nil {
			p.errs = append(p.errs, err)
		}
	} else {
		p.errs = append(p.errs, fmt.Errorf("%s %d:%d not a ; after type definition", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn))
		p.consume(untils_semicolon)
		p.Next()
		return
	}
	if p.token.Type != lex.TOKEN_SEMICOLON {
		err = fmt.Errorf("%s not a \";\" after a variable declare ", p.errorMsgPrefix())
		p.errs = append(p.errs, err)
		return
	}
	p.Next() // look next
	if !noDefaultValue && len(names) != len(expressions) {
		p.errs = append(p.errs, fmt.Errorf("%s name list and value list has no same length", p.errorMsgPrefix()))
		return
	}
	vs = []*ast.VariableDefinition{}
	for _, v := range names {
		vd := &ast.VariableDefinition{}
		vd.Name = v.Name
		vt := &ast.VariableType{}
		*vt = *t
		vd.Typ = vt
		if len(ispublic) > 0 && ispublic[0] {
			vd.AccessProperty.Access = ast.ACCESS_PUBLIC
		} else {
			vd.AccessProperty.Access = ast.ACCESS_PRIVATE
		}
		vd.Pos = v.Pos
		vs = append(vs, vd)
	}
	return vs, nil
}

//func (p *Parser) parseTypes() ([]*ast.VariableType, error) {
//	ret := []*ast.VariableType{}
//	var t *ast.VariableType
//	for {
//		t = p.parseType()
//		if t == nil { // not a type
//			return ret, nil
//		}
//		ret = append(ret, t)
//		p.Next()
//		if p.token.Type != lex.TOKEN_COMMA {
//			break
//		}
//		t = p.parseType()
//		if t == nil {
//			return ret, fmt.Errorf("%s %d:%d is not type", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn)
//		} else {
//			ret = append(ret, t)
//		}
//	}
//	return ret, nil
//}

func (p *Parser) parseType() (*ast.VariableType, error) {
	switch p.token.Type {
	case lex.TOKEN_LB:
		p.Next()
		if p.token.Type != lex.TOKEN_RB {
			// [ and ] not match
			return nil, fmt.Errorf("%s [ and ] not match", p.errorMsgPrefix())
		}
		//lookahead
		p.Next() //skip ]
		t, err := p.parseType()
		if err != nil {
			return nil, err
		}
		tt := &ast.VariableType{}
		tt.CombinationType = &ast.VariableType{}
		tt.CombinationType.Typ = ast.VARIABLE_TYPE_ARRAY
		tt.CombinationType = t
		return tt, nil
	case lex.TOKEN_BOOL:
		p.Next()
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BOOL,
		}, nil
	case lex.TOKEN_BYTE:
		p.Next()
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BYTE,
		}, nil
	case lex.TOKEN_INT:
		p.Next()
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_INT,
		}, nil
	case lex.TOKEN_FLOAT:
		p.Next()
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_FLOAT,
		}, nil
	case lex.TOKEN_STRING:
		p.Next()
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_STRING,
		}, nil
	case lex.TOKEN_IDENTIFIER:
		return p.parseIdentiferType()
	}
	return nil, fmt.Errorf("%s unkown type", p.errorMsgPrefix())
}

func (p *Parser) parseIdentiferType() (*ast.VariableType, error) {
	name := p.token.Data.(string)
	p.Next()
	if p.token.Type != lex.TOKEN_DOT {
		return &ast.VariableType{
			Typ:   ast.VARIABLE_TYPE_NAME,
			Lname: name,
		}, nil
	}
	// a.b
	p.Next() // skip dot
	if p.token.Type != lex.TOKEN_IDENTIFIER {
		return nil, fmt.Errorf("%s not a identifier after dot", p.errorMsgPrefix())
	}
	name2 := p.token.Data.(string)
	p.Next()
	return &ast.VariableType{
		Typ:   ast.VARIABLE_TYPE_DOT,
		Lname: name,
		Rname: name2,
	}, nil
}

//at least one name
func (p *Parser) parseNameList() (names []*ast.NameWithPos, err error) {
	if p.token.Type != lex.TOKEN_IDENTIFIER {
		return nil, fmt.Errorf("%s is not name,but %s", p.errorMsgPrefix(), p.token.Desp)
	}
	names = []*ast.NameWithPos{}
	for p.token.Type == lex.TOKEN_IDENTIFIER && !p.eof {
		names = append(names, &ast.NameWithPos{
			Name: p.token.Data.(string),
			Pos:  p.mkPos(),
		})
		p.Next()
		if p.token.Type != lex.TOKEN_COMMA {
			// not a ,
			break
		}
		p.Next()
		if p.token.Type != lex.TOKEN_IDENTIFIER {
			err = fmt.Errorf("%s %d:%d not identifier after a comma,but %s ", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn, p.token.Desp)
		}
	}
	return
}

func (p *Parser) consume(untils map[int]bool) {
	if len(untils) == 0 {
		panic("no token to consume")
	}
	var ok bool
	for !p.eof {
		if _, ok = untils[p.token.Type]; ok {
			return
		}
		p.Next()
	}
}

//imports,alway call next
func (p *Parser) parseImports() {
	p.Next()
	if p.eof {
		return
	}
	if p.token.Type != lex.TOKEN_IMPORT {
		// not a import
		return
	}
	syntaxErr := func() error {
		return fmt.Errorf("%s %d:%d import should be like this import \"github.com/xxx/yyy\" [as xyz] ; ", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn)
	}
	// p.token.Type == lex.TOKEN_IMPORT
	p.Next()
	if p.token.Type != lex.TOKEN_LITERAL_STRING {
		p.consume(untils_semicolon)
		p.errs = append(p.errs, syntaxErr())
		p.parseImports()
		return
	}
	packagename := p.token.Data.(string)
	p.Next()
	if p.token.Type == lex.TOKEN_AS {
		i := &ast.Imports{}
		p.lexPos2AstPos(p.token, &i.Pos)
		i.Name = packagename
		p.Next()
		if p.token.Type != lex.TOKEN_IDENTIFIER {
			p.consume(untils_semicolon)
			p.errs = append(p.errs, syntaxErr())
			p.parseImports()
			return
		}
		i.Alias = p.token.Data.(string)
		p.Next()
		if p.token.Type != lex.TOKEN_SEMICOLON {
			p.consume(untils_semicolon)
			p.errs = append(p.errs, syntaxErr())
			p.parseImports()
			return
		}
		*p.tops = append(*p.tops)
		p.parseImports()
		return
	} else if p.token.Type == lex.TOKEN_SEMICOLON {
		i := &ast.Imports{}
		i.Name = packagename
		p.lexPos2AstPos(p.token, &i.Pos)
		*p.tops = append(*p.tops, &ast.Node{
			Data: i,
		})
		p.parseImports()
		return
	} else {
		p.consume(untils_rc)
		p.errs = append(p.errs, syntaxErr())
		p.parseImports()
		return
	}
}

func (p *Parser) lexPos2AstPos(t *lex.Token, pos *ast.Pos) {
	pos.Filename = p.filename
	pos.StartLine = t.Match.StartLine
	pos.StartColumn = t.Match.StartColumn
}

func (p *Parser) parseTypeOrTypedNames() (names []*ast.VariableDefinition, err error) {

	return nil, nil
}

// a,b int or int,bool  c xxx
func (p *Parser) parseTypedNames() (names []*ast.VariableDefinition, err error) {
	names = []*ast.VariableDefinition{}
	for !p.eof {
		ns, err := p.parseNameList()
		if err != nil {
			return names, err
		}
		t, err := p.parseType()
		if err != nil {
			return names, err
		}
		for _, v := range ns {
			vd := &ast.VariableDefinition{}
			vd.Name = v.Name
			vd.Pos = v.Pos
			vd.Typ = &ast.VariableType{}
			*vd.Typ = *t
			names = append(names, vd)
		}
		if p.token.Type != lex.TOKEN_COMMA { // not a commna
			break
		}
	}
	return nil, nil
}

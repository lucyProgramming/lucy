package parser

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
	"github.com/timtadh/lexmachine"
)

func Parse(tops *[]*ast.Node, filename string, bs []byte) []error {
	return (&Parser{bs: bs, tops: tops, filename: filename}).parse()
}

type Parser struct {
	bs               []byte
	tops             *[]*ast.Node
	ExpressionParser *ExpressionParser
	Function         *Function
	Class            *Class
	scanner          *lexmachine.Scanner
	filename         string
	token            *lex.Token
	eof              bool
	errs             []error
}

func (p *Parser) parse() []error {
	p.ExpressionParser = &ExpressionParser{p}
	p.Function = &Function{p}
	p.Class = &Class{}
	p.Class.parser = p
	p.errs = []error{}
	var err error
	p.scanner, err = lex.Lexer.Scanner(p.bs)
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
		//end of file
		return p.errs
	}
	ispublic := false
	resetProperty := func() {
		ispublic = false
	}
	for !p.eof {
		fmt.Println("!!!!!!!!!!!!!!!!!!!!", p.token.Desp)
		switch p.token.Type {
		case lex.TOKEN_SEMICOLON:
			p.Next()
			continue
		case lex.TOKEN_VAR:
			vs := p.parseVarDefinition(ispublic)
			if vs != nil && len(vs) > 0 {
				for _, v := range vs {
					*p.tops = append(*p.tops, &ast.Node{
						Data: v,
					})
				}
			}
			resetProperty()
		//case lex.TOKEN_IDENTIFIER:

		case lex.TOKEN_ENUM:
			e := p.parseEnum(ispublic)
			if e != nil {
				*p.tops = append(*p.tops, &ast.Node{
					Data: e,
				})
			}
		case lex.TOKEN_FUNCTION:
			f, err := p.Function.parse(ispublic)
			if err != nil {
				p.errs = append(p.errs, err)
				continue
			}
			*p.tops = append(*p.tops, &ast.Node{
				Data: f,
			})
		//case lex.TOKEN_LC:
		case lex.TOKEN_CLASS:
			c, err := p.Class.parse(ispublic)
			if err != nil {
				p.errs = append(p.errs, err)
				continue
			}
			*p.tops = append(*p.tops, &ast.Node{
				Data: c,
			})
		case lex.TOKEN_PUBLIC:
			ispublic = true
			p.Next()
		case lex.TOKEN_CONST:
			p.Next()
			if p.eof {
				p.errs = append(p.errs, fmt.Errorf("%s %d:%d eof after const", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn))
				return p.errs
			}
			names, err := p.parseNameList()
			if err != nil {
				p.errs = append(p.errs, err)
				p.consume(lex.TOKEN_SEMICOLON)
				p.Next()
				continue
			}
			if p.token.Type != lex.TOKEN_ASSIGN {
				p.errs = append(p.errs, fmt.Errorf("%s except = ,but %s", p.errorMsgPrefix(), p.token.Desp))
				p.consume(lex.TOKEN_SEMICOLON)
				p.Next()
				continue
			}
			p.Next()
			if p.eof {
				p.errs = append(p.errs, fmt.Errorf("%s %d:%d eof after const", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn))
				continue
			}
			es, err := p.ExpressionParser.parseExpressions()
			if err != nil {
				p.errs = append(p.errs, err)
				p.consume(lex.TOKEN_SEMICOLON)
				p.Next()
				continue
			}
			if len(es) != len(names) {
				p.errs = append(p.errs, fmt.Errorf("%s %d:%d mame and value not match", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn))
				p.consume(lex.TOKEN_SEMICOLON)
				p.Next()
				continue
			}
			if p.token.Type != lex.TOKEN_SEMICOLON {
				p.errs = append(p.errs, fmt.Errorf("%s %d:%d not ; after const definition", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn))
				p.consume(lex.TOKEN_SEMICOLON)
				p.Next()
				continue
			}
			for k, v := range names {
				c := &ast.Const{}
				c.Name = v.Name
				c.Init = es[k]
				if ispublic {
					c.Access = ast.ACCESS_PUBLIC
				} else {
					c.Access = ast.ACCESS_PRIVATE
				}
				*p.tops = append(*p.tops, &ast.Node{
					Data: c,
				})
			}
			p.Next()
			continue
		case lex.TOKEN_PRIVATE: //is a default attribute
			ispublic = false
			p.Next()
		default:
			p.errs = append(p.errs, fmt.Errorf("%s %d:%d token(%s) is not except", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn, p.token.Desp))
			p.consume(lex.TOKEN_SEMICOLON)
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

func (p *Parser) Next() {
	var err error
	var tok interface{}
	for p.eof == false {
		tok, err, p.eof = p.scanner.Next()
		if err != nil {
			p.eof = true
			return
		}
		if tok != nil && tok.(*lex.Token).Type != lex.TOKEN_CRLF {
			p.token = tok.(*lex.Token)
			fmt.Println("################", p.token.Desp)
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
		return fmt.Sprintf("%s %d:%d", pos[0].Filename, pos[0].StartLine, pos[0].StartColumn)
	}
	return fmt.Sprintf("%s %d:%d", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn)
}

//var a,b,c int,char,bool  | var a,b,c int = 123;
func (p *Parser) parseVarDefinition(ispublic ...bool) (vs []*ast.VariableDefinition) {
	p.Next()
	if p.eof {
		p.unexpectedErr()
		return
	}
	names, err := p.parseNameList()
	if err != nil {
		p.errs = append(p.errs, err)
		p.consume(lex.TOKEN_SEMICOLON)
		p.Next()
		return
	}
	if p.eof {
		p.unexpectedErr()
		return
	}
	if len(names) == 0 {
		p.errs = append(p.errs, fmt.Errorf("%s %d:%d no variable name defined", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn))
		p.consume(lex.TOKEN_SEMICOLON)
		p.Next()
		return
	}
	t, err := p.parseType()
	if t == nil {
		p.errs = append(p.errs, fmt.Errorf("%s %d:%d no variable type found or defined wrong", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn))
		p.consume(lex.TOKEN_SEMICOLON)
		p.Next()
		return
	}

	var expressions []*ast.Expression
	//value , no default value definition
	if p.token.Type == lex.TOKEN_SEMICOLON {
		p.Next()
	} else if lex.TOKEN_ASSIGN == p.token.Type { //assign
		p.Next()
		expressions, err = p.ExpressionParser.parseExpressions()
		if err != nil {
			p.errs = append(p.errs, err)
		}
		if p.token.Type != lex.TOKEN_SEMICOLON {
			p.errs = append(p.errs, fmt.Errorf("%s not a \";\" after a expression list ", p.errorMsgPrefix()))
			p.consume(lex.TOKEN_SEMICOLON)
			p.Next()
			return
		}
	} else {
		p.errs = append(p.errs, fmt.Errorf("%s %d:%d not a ; after type definition", p.filename, p.token.Match.StartLine, p.token.Match.StartColumn))
		p.consume(lex.TOKEN_SEMICOLON)
		p.Next()
		return
	}
	if len(names) != len(expressions) {
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
	return vs
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
		p.Next()
		t, err := p.parseType()
		if err != nil {
			return nil, err
		}
		tt := &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_COMBINATION,
		}
		tt.CombinationType.Typ = ast.COMBINATION_TYPE_ARRAY
		tt.CombinationType.Combination = *t
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
			Typ: ast.VARIABLE_TYPE_BYTE,
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
	return nil, nil
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

func (p *Parser) consume(untils ...int) {
	fmt.Println("consume.........")
	if len(untils) == 0 {
		panic("no token to consume")
	}
	m := make(map[int]bool)
	for _, v := range untils {
		m[v] = true
	}
	var ok bool
	for !p.eof {
		if _, ok = m[p.token.Type]; ok {
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
		p.consume(lex.TOKEN_SEMICOLON)
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
			p.consume(lex.TOKEN_SEMICOLON)
			p.errs = append(p.errs, syntaxErr())
			p.parseImports()
			return
		}
		i.Alias = p.token.Data.(string)
		p.Next()
		if p.token.Type != lex.TOKEN_SEMICOLON {
			p.consume(lex.TOKEN_SEMICOLON)
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
		p.consume(lex.TOKEN_SEMICOLON)
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

// a,b int or int,bool  c xxx
func (p *Parser) parseTypedNames() (names []*ast.VariableDefinition, err error) {
	//names = []*ast.VariableDefinition{}
	//for {
	//	if p.token.Type == lex.TOKEN_IDENTIFIER {
	//		ns, err := p.parseNameList()
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//	}
	//}
	return nil, nil
}

func (p *Parser) insertVariableIntoBlock(b *ast.Block, vars []*ast.VariableDefinition) (errs []error) {
	errs = []error{}
	if vars == nil || len(vars) == 0 {
		return
	}
	if b.SymbolicTable.ItemsMap == nil {
		b.SymbolicTable.ItemsMap = make(map[string]*ast.SymbolicItem)
	}
	var err error
	for _, v := range vars {
		if v.Name == "" {
			continue
		}
		err = b.SymbolicTable.Insert(v.Name, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return
}

package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//imports,alway call next
func (p *Parser) parseImports() {
	if p.token.Type != lex.TOKEN_IMPORT {
		// not a import
		return
	}
	// p.token.Type == lex.TOKEN_IMPORT
	p.Next()
	if p.token.Type != lex.TOKEN_LITERAL_STRING {
		p.consume(untils_semicolon)
		p.errs = append(p.errs, fmt.Errorf("%s expect string literal after import", p.errorMsgPrefix()))
		p.parseImports()
		return
	}
	packagename := p.token.Data.(string)
	p.Next()
	if p.token.Type == lex.TOKEN_AS {
		i := &ast.Import{}
		i.Pos = &ast.Pos{}
		p.lexPos2AstPos(p.token, i.Pos)
		i.Name = packagename
		p.Next()
		if p.token.Type != lex.TOKEN_IDENTIFIER {
			p.consume(untils_semicolon)
			p.Next()
			p.errs = append(p.errs, fmt.Errorf("%s expect identifier after as", p.errorMsgPrefix()))
			p.parseImports()
			return
		}
		i.AccessName = p.token.Data.(string)
		p.Next()
		if p.token.Type != lex.TOKEN_SEMICOLON {
			p.consume(untils_semicolon)
			p.Next()
			p.errs = append(p.errs, fmt.Errorf("%s  semicolon after import statement", p.errorMsgPrefix()))
			p.parseImports()
			return
		}
		p.Next()
		*p.tops = append(*p.tops, &ast.Node{
			Data: i,
		})
		p.insertImports(i)
		p.parseImports()
		return
	} else if p.token.Type == lex.TOKEN_SEMICOLON {
		i := &ast.Import{}
		i.Name = packagename
		i.Pos = &ast.Pos{}
		p.lexPos2AstPos(p.token, i.Pos)
		*p.tops = append(*p.tops, &ast.Node{
			Data: i,
		})
		p.Next()
		p.insertImports(i)
		p.parseImports()
		return
	} else {
		p.consume(untils_semicolon)
		p.Next()
		p.errs = append(p.errs, fmt.Errorf("%s expect semicolon after", p.errorMsgPrefix()))
		p.parseImports()
		return
	}
}

func (p *Parser) insertImports(im *ast.Import) {
	if p.imports == nil {
		p.imports = make(map[string]*ast.Import)
	}
	access, err := im.GetAccessName()
	if err != nil {
		p.errs = append(p.errs, fmt.Errorf("%s %v", p.errorMsgPrefix(im.Pos), err))
		return
	}
	if p.imports[access] != nil {
		p.errs = append(p.errs, fmt.Errorf("%s package %s reimported", p.errorMsgPrefix(im.Pos), access))
		return
	}
	p.imports[access] = im
}

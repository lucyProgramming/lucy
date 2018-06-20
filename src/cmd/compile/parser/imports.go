package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

//imports,always call next
func (p *Parser) parseImports() {
	if p.token.Type != lex.TOKEN_IMPORT {
		// not a import
		return
	}
	p.Next()
	if p.token.Type != lex.TOKEN_LITERAL_STRING {
		p.consume(untilSemicolon)
		p.errs = append(p.errs, fmt.Errorf("%s expect 'string_literal' after import,but '%s'",
			p.errorMsgPrefix(), p.token.Description))
		p.parseImports()
		return
	}
	packageName := p.token.Data.(string)
	p.Next()
	i := &ast.Import{}
	i.ImportName = packageName
	i.Pos = p.mkPos()
	if p.token.Type == lex.TOKEN_AS { // import "xxxxxxxxxxx" as yyy ;
		p.Next() // skip as
		if p.token.Type != lex.TOKEN_IDENTIFIER {
			p.errs = append(p.errs, fmt.Errorf("%s expect 'identifier' after 'as',but '%s'",
				p.errorMsgPrefix(), p.token.Description))
			p.consume(untilSemicolon)
			p.Next()
			p.insertImports(i)
			p.parseImports()
			return
		} else {
			i.AccessName = p.token.Data.(string)
			p.Next() // skip identifier
		}
	}
	if p.token.Type != lex.TOKEN_SEMICOLON {
		p.errs = append(p.errs, fmt.Errorf("%s expect semicolon, but '%s'",
			p.errorMsgPrefix(), p.token.Description))
		p.consume(untilSemicolon)
		p.consume(untilSemicolon)
		p.Next()
		p.insertImports(i)
		p.parseImports()
		return
	}
	p.insertImports(i)
	p.Next() // skip ;
	p.parseImports()
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
		p.errs = append(p.errs, fmt.Errorf("%s package '%s' reimported",
			p.errorMsgPrefix(im.Pos), access))
		return
	}
	p.imports[access] = im
	*p.tops = append(*p.tops, &ast.Top{
		Data: im,
	})
}

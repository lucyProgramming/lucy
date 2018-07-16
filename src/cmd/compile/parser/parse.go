package parser

import (
	"fmt"

	"bytes"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func Parse(tops *[]*ast.Top, filename string, bs []byte, onlyParseImport bool, nErrors2Stop int) []error {
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
	onlyParseImport bool
	bs              []byte
	lines           [][]byte
	tops            *[]*ast.Top
	scanner         *lex.Lexer
	filename        string
	expectLf        bool
	lastToken       *lex.Token
	token           *lex.Token
	errs            []error
	imports         map[string]*ast.Import
	nErrors2Stop    int
	// parsers
	ExpressionParser *ExpressionParser
	FunctionParser   *FunctionParser
	ClassParser      *ClassParser
	BlockParser      *BlockParser
	InterfaceParser  *InterfaceParser
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
	parser.InterfaceParser = &InterfaceParser{}
	parser.InterfaceParser.parser = parser
	parser.BlockParser = &BlockParser{}
	parser.BlockParser.parser = parser
}

func (parser *Parser) Parse() []error {
	parser.initParser()
	parser.scanner = lex.New(parser.bs, 1, 1)
	parser.Next() //
	parser.lines = bytes.Split(parser.bs, []byte("\n"))
	if parser.token.Type == lex.TokenEof {
		//TODO::empty source file , should forbidden???
		return nil
	}
	parser.parseImports()
	if parser.onlyParseImport { // only parse imports
		return parser.errs
	}
	var accessControlToken *lex.Token
	isFinal := false
	resetSomeProperty := func() {
		accessControlToken = nil
		isFinal = false
	}
	isPublic := func() bool {
		return accessControlToken != nil && accessControlToken.Type == lex.TokenPublic
	}
	for parser.token.Type != lex.TokenEof {
		if len(parser.errs) > parser.nErrors2Stop {
			break
		}
		switch parser.token.Type {
		case lex.TokenSemicolon: // empty statement, no big deal
			parser.Next()
			continue
		case lex.TokenPublic:
			accessControlToken = parser.token
			parser.Next()
			if err := parser.validAfterPublic(); err != nil {
				accessControlToken = nil
			}
			continue
		case lex.TokenFinal:
			isFinal = true
			parser.Next()
			if err := parser.validAfterFinal(); err != nil {
				isFinal = false
			}
			continue
		case lex.TokenVar:
			pos := parser.mkPos()
			parser.Next() // skip var key word
			vs, es, err := parser.parseConstDefinition(true)
			if err != nil {
				parser.consume(untilSemicolon)
				parser.Next()
				continue
			}
			d := &ast.ExpressionDeclareVariable{Variables: vs, InitValues: es}
			isPublic := isPublic()
			e := &ast.Expression{
				Type:     ast.ExpressionTypeVar,
				Data:     d,
				Pos:      pos,
				IsPublic: isPublic,
			}
			*parser.tops = append(*parser.tops, &ast.Top{
				Data: e,
			})
			resetSomeProperty()
		case lex.TokenIdentifier:
			e, err := parser.ExpressionParser.parseExpression(true)
			if err != nil {
				parser.errs = append(parser.errs, err)
				parser.consume(untilSemicolon)
				parser.Next()
				continue
			}
			e.IsPublic = isPublic()
			parser.validStatementEnding()
			*parser.tops = append(*parser.tops, &ast.Top{
				Data: e,
			})
			resetSomeProperty()
		case lex.TokenEnum:
			e, err := parser.parseEnum()
			if err != nil {
				parser.consume(untilRc)
				parser.Next()
				resetSomeProperty()
				continue
			}
			isPublic := isPublic()
			if isPublic {
				e.AccessFlags |= cg.ACC_CLASS_PUBLIC
			}
			if e != nil {
				*parser.tops = append(*parser.tops, &ast.Top{
					Data: e,
				})
			}
			resetSomeProperty()
		case lex.TokenFunction:
			f, err := parser.FunctionParser.parse(true)
			if err != nil {
				parser.consume(untilRc)
				parser.Next()
				continue
			}
			isPublic := isPublic()
			if isPublic {
				f.AccessFlags |= cg.ACC_METHOD_PUBLIC
			}
			*parser.tops = append(*parser.tops, &ast.Top{
				Data: f,
			})
			resetSomeProperty()
		case lex.TokenLc:
			b := &ast.Block{}
			parser.Next() // skip {
			parser.BlockParser.parseStatementList(b, true)
			if parser.token.Type != lex.TokenRc {
				parser.errs = append(parser.errs, fmt.Errorf("%s expect '}', but '%s'",
					parser.errorMsgPrefix(), parser.token.Description))
				parser.consume(untilRc)
			}
			parser.Next() // skip }
			*parser.tops = append(*parser.tops, &ast.Top{
				Data: b,
			})
		case lex.TokenClass, lex.TokenInterface:
			var c *ast.Class
			var err error
			if parser.token.Type == lex.TokenClass {
				c, err = parser.ClassParser.parse()
			} else {
				c, err = parser.InterfaceParser.parse()
			}
			if err != nil {
				parser.errs = append(parser.errs, err)
				parser.consume(untilRc)
				parser.Next()
				resetSomeProperty()
				continue
			}
			if c == nil && err == nil {
				panic(1)
			}
			*parser.tops = append(*parser.tops, &ast.Top{
				Data: c,
			})
			if isPublic() {
				c.AccessFlags |= cg.ACC_CLASS_PUBLIC
			}
			if isFinal {
				c.AccessFlags |= cg.ACC_CLASS_FINAL
			}
			resetSomeProperty()
		case lex.TokenConst:
			parser.Next() // skip const key word
			vs, es, err := parser.parseConstDefinition(false)
			if err != nil {
				parser.consume(untilSemicolon)
				parser.Next()
				resetSomeProperty()
				continue
			}
			if len(vs) != len(es) {
				parser.errs = append(parser.errs,
					fmt.Errorf("%s cannot assign %d values to %d destinations",
						parser.errorMsgPrefix(parser.mkPos()), len(es), len(vs)))
			}
			isPublic := isPublic()
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
					*parser.tops = append(*parser.tops, &ast.Top{
						Data: c,
					})
				}
			}
			resetSomeProperty()
			continue
		case lex.TokenType:
			a, err := parser.parseTypeAlias()
			if err != nil {
				parser.consume(untilSemicolon)
				parser.Next()
				resetSomeProperty()
				continue
			}
			*parser.tops = append(*parser.tops, &ast.Top{
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
			parser.errs = append(parser.errs, fmt.Errorf("%s token '%s' is not except",
				parser.errorMsgPrefix(), parser.token.Description))
			parser.consume(untilSemicolon)
			resetSomeProperty()
		}
	}
	return parser.errs
}

func (parser *Parser) parseTypes() ([]*ast.Type, error) {
	ret := []*ast.Type{}
	for parser.token.Type != lex.TokenEof {
		t, err := parser.parseType()
		if err != nil {
			return ret, err
		}
		ret = append(ret, t)
		if parser.token.Type != lex.TokenComma {
			break
		}
		parser.Next() // skip ,
	}
	return ret, nil
}

func (parser *Parser) validAfterPublic() error {
	if parser.token.Type == lex.TokenFunction ||
		parser.token.Type == lex.TokenClass ||
		parser.token.Type == lex.TokenEnum ||
		parser.token.Type == lex.TokenIdentifier ||
		parser.token.Type == lex.TokenInterface ||
		parser.token.Type == lex.TokenConst ||
		parser.token.Type == lex.TokenVar ||
		parser.token.Type == lex.TokenFinal {
		return nil
	}
	err := fmt.Errorf("%s cannot have token '%s' after 'public'",
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

func (parser *Parser) validStatementEnding() {
	if parser.token.Type == lex.TokenSemicolon ||
		(parser.lastToken != nil && parser.lastToken.Type == lex.TokenRc) {
		return
	}
	var token *lex.Token
	if nil != parser.lastToken {
		token = parser.lastToken
	}
	if token == nil {
		token = parser.token
	}
	parser.errs = append(parser.errs, fmt.Errorf("%s missing semicolon", parser.errorMsgPrefix(&ast.Position{
		Filename:    parser.filename,
		StartLine:   token.StartLine,
		StartColumn: token.StartColumn,
	})))
}

func (parser *Parser) mkPos() *ast.Position {
	return &ast.Position{
		Filename:    parser.filename,
		StartLine:   parser.token.StartLine,
		StartColumn: parser.token.StartColumn,
		Offset:      parser.scanner.GetOffSet(),
	}
}

// str := "hello world"   a,b = 123 or a b ;
func (parser *Parser) parseConstDefinition(needType bool) ([]*ast.Variable, []*ast.Expression, error) {
	names, err := parser.parseNameList()
	if err != nil {
		return nil, nil, err
	}
	var variableType *ast.Type
	//trying to parse type
	if parser.isValidTypeBegin() || needType {
		variableType, err = parser.parseType()
		if err != nil {
			parser.errs = append(parser.errs, err)
			return nil, nil, err
		}
	}
	mkResult := func() []*ast.Variable {
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
	if parser.token.Type != lex.TokenAssign &&
		parser.token.Type != lex.TokenColonAssign {
		return mkResult(), nil, err
	}
	if parser.token.Type != lex.TokenAssign {
		parser.errs = append(parser.errs, fmt.Errorf("%s use '=' instead of ':='", parser.errorMsgPrefix()))
	}
	parser.Next() // skip = or :=
	es, err := parser.ExpressionParser.parseExpressions()
	if err != nil {
		return nil, nil, err
	}
	return mkResult(), es, nil
}

func (parser *Parser) Next() {
	var err error
	var tok *lex.Token
	parser.lastToken = parser.token
	for {
		tok, err = parser.scanner.Next()
		if err != nil {
			parser.errs = append(parser.errs, fmt.Errorf("%s %s", parser.errorMsgPrefix(), err.Error()))
		}
		if tok == nil {
			continue
		}
		if parser.expectLf {
			if tok.Type != lex.TokenLf {
				parser.errs = append(parser.errs, fmt.Errorf("%s expect new line", parser.errorMsgPrefix()))
			}
			parser.expectLf = false
		}
		parser.token = tok
		if tok.Type != lex.TokenLf {
			//if parser.token.Description != "" {
			//	fmt.Println("#########", parser.token.Type, parser.token.Description)
			//} else {
			//	fmt.Println("#########", parser.token.Type, parser.token.Data)
			//}
			break
		}
	}
	return
}

/*
	errorMsgPrefix(pos) only receive one argument
*/
func (parser *Parser) errorMsgPrefix(pos ...*ast.Position) string {
	if len(pos) > 0 {
		return fmt.Sprintf("%s:%d:%d", pos[0].Filename, pos[0].StartLine, pos[0].StartColumn)
	}
	line, column := parser.scanner.GetLineAndColumn()
	return fmt.Sprintf("%s:%d:%d", parser.filename, line, column)
}

func (parser *Parser) consume(until map[int]bool) {
	if len(until) == 0 {
		panic("no token to consume")
	}
	var ok bool
	for parser.token.Type != lex.TokenEof {
		if _, ok = until[parser.token.Type]; ok {
			return
		}
		parser.Next()
	}
}

func (parser *Parser) parseTypeAlias() (*ast.ExpressionTypeAlias, error) {
	parser.Next() // skip type key word
	if parser.token.Type != lex.TokenIdentifier {
		err := fmt.Errorf("%s expect identifer,but %s", parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return nil, err
	}
	ret := &ast.ExpressionTypeAlias{}
	ret.Pos = parser.mkPos()
	ret.Name = parser.token.Data.(string)
	parser.Next() // skip identifier
	if parser.token.Type != lex.TokenAssign {
		err := fmt.Errorf("%s expect '=',but %s", parser.errorMsgPrefix(), parser.token.Description)
		parser.errs = append(parser.errs, err)
		return nil, err
	}
	parser.Next() // skip =
	var err error
	ret.Type, err = parser.parseType()
	if err != nil {
		return nil, err
	}
	return ret, err
}

func (parser *Parser) parseTypedName() (vs []*ast.Variable, err error) {
	names, err := parser.parseNameList()
	if err != nil {
		return nil, err
	}
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
			parser.Next()
		}
	}
	return vs, nil
}

//func (parser *Parser) lexPos2AstPos(t *lex.Token, pos *ast.Position) {
//	pos.Filename = parser.filename
//	pos.StartLine = t.StartLine
//	pos.StartColumn = t.StartColumn
//}

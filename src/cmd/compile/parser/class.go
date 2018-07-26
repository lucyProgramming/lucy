package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type ClassParser struct {
	parser             *Parser
	ret                *ast.Class
	isStatic           bool
	isVolatile         bool
	isSynchronized     bool
	isFinal            bool
	accessControlToken *lex.Token
}

func (classParser *ClassParser) resetProperty() {
	classParser.isStatic = false
	classParser.isVolatile = false
	classParser.isSynchronized = false
	classParser.isFinal = false
	classParser.accessControlToken = nil
}

func (classParser *ClassParser) Next(lfIsToken bool) {
	classParser.parser.Next(lfIsToken)
}

func (classParser *ClassParser) consume(m map[lex.TokenKind]bool) {
	classParser.parser.consume(m)
}

func (classParser *ClassParser) parseClassName() (string, error) {
	if classParser.parser.token.Type != lex.TokenIdentifier {
		err := fmt.Errorf("%s expect identifier for class`s name,but '%s'",
			classParser.parser.errorMsgPrefix(), classParser.parser.token.Description)
		classParser.parser.errs = append(classParser.parser.errs, err)
		return "", err
	}
	name := classParser.parser.token.Data.(string)
	classParser.Next(lfNotToken)
	if classParser.parser.token.Type == lex.TokenSelection {
		classParser.Next(lfNotToken) // skip .
		if classParser.parser.token.Type != lex.TokenIdentifier {
			err := fmt.Errorf("%s expect identifer for class`s name,but '%s'", classParser.parser.errorMsgPrefix(),
				classParser.parser.token.Description)
			classParser.parser.errs = append(classParser.parser.errs, err)
		}
		name += "." + classParser.parser.token.Data.(string)
		classParser.Next(lfNotToken) // skip name identifier
	}
	return name, nil
}

func (classParser *ClassParser) parseImplementsInterfaces() ([]*ast.NameWithPos, error) {
	ret := []*ast.NameWithPos{}
	for classParser.parser.token.Type != lex.TokenEof {
		pos := classParser.parser.mkPos()
		name, err := classParser.parseClassName()
		if err != nil {
			return nil, err
		}
		ret = append(ret, &ast.NameWithPos{
			Name: name,
			Pos:  pos,
		})
		if classParser.parser.token.Type == lex.TokenComma {
			classParser.Next(lfNotToken)
		} else {
			break
		}
	}
	return ret, nil
}

func (classParser *ClassParser) parse() (classDefinition *ast.Class, err error) {
	classParser.resetProperty()
	classDefinition = &ast.Class{}
	classParser.ret = classDefinition
	classParser.ret.Pos = classParser.parser.mkPos()
	classParser.Next(lfIsToken) // skip class key word
	classParser.parser.unExpectNewLineAndSkip()
	classParser.ret.Name, err = classParser.parseClassName()
	classParser.ret.Block.IsClassBlock = true
	if err != nil {
		if classParser.ret.Name == "" {
			compileAutoName()
		}
		classParser.consume(untilLc)
	}
	if classParser.parser.token.Type == lex.TokenExtends { // parse father expression
		//classParser.ret.HaveExtends = true
		classParser.Next(lfNotToken) // skip extends
		classParser.ret.Pos = classParser.parser.mkPos()
		t, err := classParser.parseClassName()
		classParser.ret.SuperClassName = t
		if err != nil {
			classParser.parser.errs = append(classParser.parser.errs, err)
			classParser.consume(untilLc)
		}
	}
	if classParser.parser.token.Type == lex.TokenImplements {
		classParser.Next(lfNotToken) // skip key word
		classParser.ret.InterfaceNames, err = classParser.parseImplementsInterfaces()
		if err != nil {
			classParser.parser.errs = append(classParser.parser.errs, err)
			classParser.consume(untilLc)
		}
	}
	classParser.parser.ifTokenIsLfThenSkip()
	if classParser.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{' but '%s'", classParser.parser.errorMsgPrefix(), classParser.parser.token.Description)
		classParser.parser.errs = append(classParser.parser.errs, err)
		return nil, err
	}
	validAfterPublic := func(keyWord string, token *lex.Token) error {
		if token.Type == lex.TokenIdentifier ||
			token.Type == lex.TokenFn ||
			token.Type == lex.TokenStatic ||
			token.Type == lex.TokenSynchronized ||
			token.Type == lex.TokenFinal {
			return nil
		}
		return fmt.Errorf("%s not a valid token after '%s'", classParser.parser.errorMsgPrefix(), keyWord)
	}
	validAfterVolatile := func(token *lex.Token) error {
		if token.Type == lex.TokenIdentifier {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'volatile'", classParser.parser.errorMsgPrefix())
	}
	validAfterSynchronized := func(token *lex.Token) error {
		if token.Type == lex.TokenFn ||
			token.Type == lex.TokenFinal {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'synchronized'", classParser.parser.errorMsgPrefix())
	}
	validAfterStatic := func(token *lex.Token) error {
		if token.Type == lex.TokenIdentifier ||
			token.Type == lex.TokenFn ||
			token.Type == lex.TokenFinal {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'static'", classParser.parser.errorMsgPrefix())
	}
	validAfterFinal := func(token *lex.Token) error {
		if token.Type == lex.TokenFn ||
			token.Type == lex.TokenSynchronized {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'final'", classParser.parser.errorMsgPrefix())
	}
	classParser.Next(lfNotToken) // skip {
	for classParser.parser.token.Type != lex.TokenEof {
		if len(classParser.parser.errs) > classParser.parser.nErrors2Stop {
			break
		}
		switch classParser.parser.token.Type {
		case lex.TokenRc:
			classParser.Next(lfNotToken)
			return
		case lex.TokenSemicolon, lex.TokenLf:
			classParser.Next(lfNotToken)
			continue
		case lex.TokenStatic:
			classParser.isStatic = true
			classParser.Next(lfIsToken)
			if classParser.parser.token.Type == lex.TokenLc {
				classParser.Next(lfNotToken) // skip {
				block := &ast.Block{}
				classParser.parser.BlockParser.parseStatementList(block, false)
				if classParser.parser.token.Type != lex.TokenRc {
					classParser.parser.errs = append(classParser.parser.errs,
						fmt.Errorf("%s expect '}' , but '%s'", classParser.parser.errorMsgPrefix(),
							classParser.parser.token.Description))
				} else {
					classParser.Next(lfNotToken) // skip }
					classParser.ret.StaticBlocks = append(classParser.ret.StaticBlocks, block)
				}
				continue
			}
			err := validAfterStatic(classParser.parser.token)
			if err != nil {
				classParser.parser.errs = append(classParser.parser.errs, err)
				classParser.isStatic = false
			}
		//access private
		case lex.TokenPublic, lex.TokenProtected, lex.TokenPrivate:
			classParser.accessControlToken = classParser.parser.token
			classParser.Next(lfIsToken)
			err = validAfterPublic(classParser.accessControlToken.Description, classParser.parser.token)
			if err != nil {
				classParser.parser.errs = append(classParser.parser.errs, err)
				classParser.accessControlToken = nil // set to nil
			}
		case lex.TokenVolatile:
			classParser.isVolatile = true
			classParser.Next(lfIsToken)
			if err := validAfterVolatile(classParser.parser.token); err != nil {
				classParser.parser.errs = append(classParser.parser.errs, err)
				classParser.isVolatile = false
			}
		case lex.TokenFinal:
			classParser.isFinal = true
			classParser.Next(lfIsToken)
			if err := validAfterFinal(classParser.parser.token); err != nil {
				classParser.parser.errs = append(classParser.parser.errs, err)
				classParser.isFinal = false
			}
		case lex.TokenIdentifier:
			err = classParser.parseField(&classParser.parser.errs)
			if err != nil {
				classParser.consume(untilSemicolonOrLf)
				classParser.Next(lfNotToken)
			}
			classParser.resetProperty()
		case lex.TokenConst: // const is for local use
			classParser.Next(lfIsToken)
			err := classParser.parseConst()
			if err != nil {
				classParser.consume(untilSemicolonOrLf)
				classParser.Next(lfNotToken)
				continue
			}
		case lex.TokenSynchronized:
			classParser.isSynchronized = true
			classParser.Next(lfIsToken)
			if err := validAfterSynchronized(classParser.parser.token); err != nil {
				classParser.parser.errs = append(classParser.parser.errs, err)
				classParser.isSynchronized = false
			}
		case lex.TokenFn:
			f, err := classParser.parser.FunctionParser.parse(true)
			if err != nil {
				classParser.resetProperty()
				classParser.Next(lfNotToken)
				continue
			}
			if classParser.ret.Methods == nil {
				classParser.ret.Methods = make(map[string][]*ast.ClassMethod)
			}
			if f.Name == "" {
				// some error
				f.Name = compileAutoName()
			}
			m := &ast.ClassMethod{}
			m.Function = f
			f.AccessFlags = 0
			if classParser.accessControlToken != nil {
				switch classParser.accessControlToken.Type {
				case lex.TokenPrivate:
					m.Function.AccessFlags |= cg.ACC_METHOD_PRIVATE
				case lex.TokenProtected:
					m.Function.AccessFlags |= cg.ACC_METHOD_PROTECTED
				case lex.TokenPublic:
					m.Function.AccessFlags |= cg.ACC_METHOD_PUBLIC
				}
			}
			if classParser.isSynchronized {
				m.Function.AccessFlags |= cg.ACC_METHOD_SYNCHRONIZED
			}
			if classParser.isStatic {
				f.AccessFlags |= cg.ACC_METHOD_STATIC
			}
			if classParser.isFinal {
				f.AccessFlags |= cg.ACC_METHOD_FINAL
			}
			if f.Name == classParser.ret.Name {
				f.Name = ast.SpecialMethodInit
			}
			classParser.ret.Methods[f.Name] = append(classParser.ret.Methods[f.Name], m)
			classParser.resetProperty()
		case lex.TokenImport:
			pos := classParser.parser.mkPos()
			classParser.parser.parseImports()
			classParser.parser.errs = append(classParser.parser.errs, fmt.Errorf("%s cannot have import at this scope",
				classParser.parser.errorMsgPrefix(pos)))
		default:
			classParser.parser.errs = append(classParser.parser.errs, fmt.Errorf("%s unexpect '%s'",
				classParser.parser.errorMsgPrefix(), classParser.parser.token.Description))
			classParser.Next(lfNotToken)
		}
	}
	return
}

func (classParser *ClassParser) parseConst() error {
	pos := classParser.parser.mkPos()
	vs, es, err := classParser.parser.parseConstDefinition(false)
	if err != nil {
		return err
	}
	if len(vs) != len(es) {
		classParser.parser.errs = append(classParser.parser.errs,
			fmt.Errorf("%s cannot assign %d values to %d destinations",
				classParser.parser.errorMsgPrefix(pos), len(es), len(vs)))
	}

	if classParser.ret.Block.Constants == nil {
		classParser.ret.Block.Constants = make(map[string]*ast.Constant)
	}
	for k, v := range vs {
		if _, ok := classParser.ret.Block.Constants[v.Name]; ok {
			classParser.parser.errs = append(classParser.parser.errs, fmt.Errorf("%s const %s alreay declared",
				classParser.parser.errorMsgPrefix(), v.Name))
			continue
		}
		if k < len(es) && es[k] != nil {
			t := &ast.Constant{}
			t.Variable = *v
			t.Expression = es[k]
			classParser.ret.Block.Constants[v.Name] = t
		}
	}
	return nil
}

func (classParser *ClassParser) parseField(errs *[]error) error {
	variables, es, err := classParser.parser.parseConstDefinition(true)
	if err == nil {
		classParser.parser.validStatementEnding()
	}
	if len(variables) == 0 {
		return err
	}
	if classParser.ret.Fields == nil {
		classParser.ret.Fields = make(map[string]*ast.ClassField)
	}
	if es != nil && len(es) != len(variables) {
		err := fmt.Errorf("%s cannot assign %d values to %d destinations",
			classParser.parser.errorMsgPrefix(), len(es), len(variables))
		*errs = append(*errs, err)
	}
	for k, v := range variables {
		if _, ok := classParser.ret.Fields[v.Name]; ok {
			classParser.parser.errs = append(classParser.parser.errs,
				fmt.Errorf("%s field %s is alreay declared",
					classParser.parser.errorMsgPrefix(), v.Name))
			continue
		}
		f := &ast.ClassField{}
		f.Name = v.Name
		f.Pos = v.Pos
		f.Type = &ast.Type{}
		*f.Type = *v.Type
		if k < len(es) && es[k] != nil {
			f.Expression = es[k]
		}
		f.AccessFlags = 0
		if classParser.isStatic {
			f.AccessFlags |= cg.ACC_FIELD_STATIC
		}
		if classParser.accessControlToken != nil {
			switch classParser.accessControlToken.Type {
			case lex.TokenPublic:
				f.AccessFlags |= cg.ACC_FIELD_PUBLIC
			case lex.TokenProtected:
				f.AccessFlags |= cg.ACC_FIELD_PROTECTED
			default: // private
				f.AccessFlags |= cg.ACC_FIELD_PRIVATE
			}
		}
		if classParser.isVolatile {
			f.AccessFlags |= cg.ACC_FIELD_VOLATILE
		}
		classParser.ret.Fields[v.Name] = f
	}
	return nil
}

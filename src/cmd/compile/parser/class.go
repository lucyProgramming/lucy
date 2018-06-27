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
	accessControlToken *lex.Token
}

func (classParser *ClassParser) Next() {
	classParser.parser.Next()
}

func (classParser *ClassParser) consume(m map[int]bool) {
	classParser.parser.consume(m)
}

func (classParser *ClassParser) parseClassName() (string, error) {
	if classParser.parser.token.Type != lex.TokenIdentifier {
		err := fmt.Errorf("%s expect class`s name,but '%s'",
			classParser.parser.errorMsgPrefix(), classParser.parser.token.Description)
		classParser.parser.errs = append(classParser.parser.errs, err)
		return "", err
	}
	name := classParser.parser.token.Data.(string)
	classParser.Next()
	if classParser.parser.token.Type == lex.TokenSelection {
		classParser.Next()
		if classParser.parser.token.Type != lex.TokenIdentifier {
			err := fmt.Errorf("%s expect identifer,but '%s'", classParser.parser.errorMsgPrefix(),
				classParser.parser.token.Description)
			classParser.parser.errs = append(classParser.parser.errs, err)
		}
		name += "." + classParser.parser.token.Data.(string)
		classParser.Next() // skip name identifier
	}
	return name, nil
}

func (classParser *ClassParser) parseInterfaces() ([]*ast.NameWithPos, error) {
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
			classParser.Next()
		} else {
			break
		}
	}
	return ret, nil
}

func (classParser *ClassParser) parse() (classDefinition *ast.Class, err error) {
	defer classParser.resetProperty()
	classDefinition = &ast.Class{}
	classParser.ret = classDefinition
	classParser.Next() // skip class key word
	classParser.ret.Pos = classParser.parser.mkPos()
	classParser.ret.Name, err = classParser.parseClassName()
	classParser.ret.Block.IsClassBlock = true
	if err != nil {
		return nil, err
	}
	if classParser.parser.token.Type == lex.TokenExtends { // parse father expression
		//classParser.ret.HaveExtends = true
		classParser.Next() // skip extends
		classParser.ret.Pos = classParser.parser.mkPos()
		if classParser.parser.token.Type != lex.TokenIdentifier {
			err = fmt.Errorf("%s class`s father must be a identifier", classParser.parser.errorMsgPrefix())
			classParser.parser.errs = append(classParser.parser.errs, err)
			classParser.consume(untilLc) //
		} else {
			t, err := classParser.parseClassName()
			classParser.ret.SuperClassName = t
			if err != nil {
				classParser.parser.errs = append(classParser.parser.errs, err)
				return nil, err
			}
		}
	}
	if classParser.parser.token.Type == lex.TokenImplements {
		classParser.Next() // skip key word
		classParser.ret.InterfaceNames, err = classParser.parseInterfaces()
		if err != nil {
			classParser.consume(untilLc)
		}
	}
	if classParser.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{' but '%s'", classParser.parser.errorMsgPrefix(), classParser.parser.token.Description)
		classParser.parser.errs = append(classParser.parser.errs, err)
		return nil, err
	}
	classParser.Next() // skip {
	classParser.resetProperty()
	validAfterPublic := func(keyWord string, token *lex.Token) error {
		if token.Type == lex.TokenIdentifier ||
			token.Type == lex.TokenFunction ||
			token.Type == lex.TokenStatic ||
			token.Type == lex.TokenSynchronized {
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
		if token.Type == lex.TokenFunction {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'synchronized'", classParser.parser.errorMsgPrefix())
	}
	validAfterStatic := func(token *lex.Token) error {
		if token.Type == lex.TokenIdentifier ||
			token.Type == lex.TokenFunction || token.Type == lex.TokenLc {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'static'", classParser.parser.errorMsgPrefix())
	}
	for classParser.parser.token.Type != lex.TokenEof {
		if len(classParser.parser.errs) > classParser.parser.nErrors2Stop {
			break
		}
		switch classParser.parser.token.Type {
		case lex.TokenRc:
			classParser.Next()
			return
		case lex.TokenSemicolon:
			classParser.Next()
			continue
		case lex.TokenStatic:
			classParser.isStatic = true
			classParser.Next()
			if classParser.parser.token.Type == lex.TokenLc {
				classParser.Next() // skip {
				block := &ast.Block{}
				classParser.parser.BlockParser.parseStatementList(block, false)
				if classParser.parser.token.Type != lex.TokenRc {
					classParser.parser.errs = append(classParser.parser.errs,
						fmt.Errorf("%s expect '}' , but '%s'", classParser.parser.errorMsgPrefix(), classParser.parser.token.Description))
				} else {
					classParser.Next() // skip }
					classParser.ret.StaticBlocks = append(classParser.ret.StaticBlocks, block)
				}
				continue
			}
			err := validAfterStatic(classParser.parser.token)
			if err != nil {
				classParser.parser.errs = append(classParser.parser.errs, err)
			}
		//access private
		case lex.TokenPublic:
			classParser.accessControlToken = classParser.parser.token
			classParser.Next()
			err = validAfterPublic("public", classParser.parser.token)
			if err != nil {
				classParser.parser.errs = append(classParser.parser.errs, err)
			}
		case lex.TokenProtected:
			classParser.accessControlToken = classParser.parser.token
			classParser.Next()
			err = validAfterPublic("protected", classParser.parser.token)
			if err != nil {
				classParser.parser.errs = append(classParser.parser.errs, err)
			}
		case lex.TokenPrivate:
			classParser.accessControlToken = classParser.parser.token
			classParser.Next() // skip private
			err = validAfterPublic("private", classParser.parser.token)
			if err != nil {
				classParser.parser.errs = append(classParser.parser.errs, err)
			}
		case lex.TokenVolatile:
			classParser.isVolatile = true
			classParser.Next()
			if err := validAfterVolatile(classParser.parser.token); err != nil {
				classParser.parser.errs = append(classParser.parser.errs, err)
			}
		case lex.TokenIdentifier:
			err = classParser.parseField(&classParser.parser.errs)
			if err != nil {
				classParser.consume(untilSemicolon)
				classParser.Next()
			}
			classParser.resetProperty()
		case lex.TokenConst: // const is for local use
			classParser.Next()
			err := classParser.parseConst()
			if err != nil {
				classParser.consume(untilSemicolon)
				classParser.Next()
				continue
			}
			classParser.resetProperty()
		case lex.TokenSynchronized:
			classParser.isSynchronized = true
			classParser.Next()
			if err := validAfterSynchronized(classParser.parser.token); err != nil {
				classParser.parser.errs = append(classParser.parser.errs, err)
			}
		case lex.TokenFunction:
			f, err := classParser.parser.FunctionParser.parse(true)
			if err != nil {
				classParser.consume(untilRc)
				classParser.Next()
				classParser.resetProperty()
				continue
			}
			if classParser.ret.Methods == nil {
				classParser.ret.Methods = make(map[string][]*ast.ClassMethod)
			}
			if f.Name == "" {
				classParser.parser.errs = append(classParser.parser.errs, fmt.Errorf("%s method has no name", classParser.parser.errorMsgPrefix(f.Pos)))
				classParser.resetProperty()
				continue
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
			} else {
				m.Function.AccessFlags |= cg.ACC_METHOD_PRIVATE
			}
			if classParser.isSynchronized {
				m.Function.AccessFlags |= cg.ACC_METHOD_SYNCHRONIZED
			}
			if classParser.isStatic {
				f.AccessFlags |= cg.ACC_METHOD_STATIC
			}
			if classParser.ret.Methods == nil {
				classParser.ret.Methods = make(map[string][]*ast.ClassMethod)
			}
			if f.Name == classParser.ret.Name {
				f.Name = ast.ConstructionMethodName
			}
			classParser.ret.Methods[f.Name] = append(classParser.ret.Methods[f.Name], m)
			classParser.resetProperty()
		default:
			classParser.parser.errs = append(classParser.parser.errs, fmt.Errorf("%s unexpect '%s'",
				classParser.parser.errorMsgPrefix(), classParser.parser.token.Description))
			classParser.Next()
		}
	}
	return
}

func (classParser *ClassParser) resetProperty() {
	classParser.isStatic = false
	classParser.isVolatile = false
	classParser.isSynchronized = false
	classParser.accessControlToken = nil
}

func (classParser *ClassParser) parseConst() error {
	pos := classParser.parser.mkPos()
	vs, es, typ, err := classParser.parser.parseConstDefinition(false)
	if err != nil {
		return err
	}
	if typ != nil && typ.Type != lex.TokenAssign {
		classParser.parser.errs = append(classParser.parser.errs,
			fmt.Errorf("%s declare const should use ‘=’ instead of ‘:=’",
				classParser.parser.errorMsgPrefix(vs[0].Pos)))
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
		if k < len(es) {
			t := &ast.Constant{}
			t.Variable = *v
			t.Expression = es[k]
			classParser.ret.Block.Constants[v.Name] = t
		}
	}
	return nil
}

func (classParser *ClassParser) parseField(errs *[]error) error {
	names, err := classParser.parser.parseNameList()
	if err != nil {
		return err
	}
	t, err := classParser.parser.parseType()
	if err != nil {
		return err
	}
	var es []*ast.Expression
	if classParser.parser.token.Type == lex.TokenAssign {
		classParser.Next() // skip =
		es, err = classParser.parser.ExpressionParser.parseExpressions()
		if err != nil {
			*errs = append(*errs, err)
		}
	}

	if classParser.parser.token.Type != lex.TokenSemicolon {
		*errs = append(*errs, fmt.Errorf("%s missing senicolon after field definition",
			classParser.parser.errorMsgPrefix()))
	}
	if classParser.ret.Fields == nil {
		classParser.ret.Fields = make(map[string]*ast.ClassField)
	}
	if es != nil && len(es) != len(names) {
		err := fmt.Errorf("%s cannot assign %d values to %d destinations",
			classParser.parser.errorMsgPrefix(), len(es), len(names))
		*errs = append(*errs, err)
	}
	for k, v := range names {
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
		if t == nil {
			panic(11)
		}
		*f.Type = *t
		f.AccessFlags = 0
		if k < len(es) && es[k] != nil {
			f.Expression = es[k]
		}
		if classParser.isStatic {
			f.AccessFlags |= cg.ACC_FIELD_STATIC
		}
		if classParser.accessControlToken == nil {
			f.AccessFlags |= cg.ACC_FIELD_PRIVATE
		} else {
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

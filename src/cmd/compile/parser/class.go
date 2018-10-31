package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type ClassParser struct {
	parser *Parser
}

func (cp *ClassParser) Next(lfIsToken bool) {
	cp.parser.Next(lfIsToken)
}

func (cp *ClassParser) consume(m map[lex.TokenKind]bool) {
	cp.parser.consume(m)
}

func (cp *ClassParser) parseClassName() (*ast.NameWithPos, error) {
	if cp.parser.token.Type != lex.TokenIdentifier {
		err := fmt.Errorf("%s expect identifier for class`s name,but '%s'",
			cp.parser.errMsgPrefix(), cp.parser.token.Description)
		cp.parser.errs = append(cp.parser.errs, err)
		return nil, err
	}
	name := cp.parser.token.Data.(string)
	pos := cp.parser.mkPos()
	cp.Next(lfNotToken)
	if cp.parser.token.Type == lex.TokenSelection {
		cp.Next(lfNotToken) // skip .
		if cp.parser.token.Type != lex.TokenIdentifier {
			err := fmt.Errorf("%s expect identifer for class`s name,but '%s'",
				cp.parser.errMsgPrefix(),
				cp.parser.token.Description)
			cp.parser.errs = append(cp.parser.errs, err)
		} else {
			name += "." + cp.parser.token.Data.(string)
			cp.Next(lfNotToken) // skip name identifier
		}
	}
	return &ast.NameWithPos{
		Name: name,
		Pos:  pos,
	}, nil
}

func (cp *ClassParser) parseImplementsInterfaces() ([]*ast.NameWithPos, error) {
	ret := []*ast.NameWithPos{}
	for cp.parser.token.Type != lex.TokenEof {
		name, err := cp.parseClassName()
		if err != nil {
			return nil, err
		}
		ret = append(ret, &ast.NameWithPos{
			Name: name.Name,
			Pos:  name.Pos,
		})
		if cp.parser.token.Type == lex.TokenComma {
			cp.Next(lfNotToken)
		} else {
			break
		}
	}
	return ret, nil
}

func (cp *ClassParser) parse(isAbstract bool) (classDefinition *ast.Class, err error) {
	isInterface := cp.parser.token.Type == lex.TokenInterface
	classDefinition = &ast.Class{}
	if isInterface {
		classDefinition.AccessFlags |= cg.ACC_CLASS_INTERFACE
		classDefinition.AccessFlags |= cg.ACC_CLASS_ABSTRACT
	}
	if isAbstract {
		classDefinition.AccessFlags |= cg.ACC_CLASS_ABSTRACT
	}
	cp.Next(lfIsToken) // skip class key word
	cp.parser.unExpectNewLineAndSkip()
	t, err := cp.parseClassName()
	if t != nil {
		classDefinition.Name = t.Name
	}
	classDefinition.Block.IsClassBlock = true
	classDefinition.Block.Class = classDefinition
	if err != nil {
		if classDefinition.Name == "" {
			compileAutoName()
		}
		cp.consume(untilLc)
	}
	classDefinition.Pos = cp.parser.mkPos()
	if cp.parser.token.Type == lex.TokenExtends { // parse father expression
		cp.Next(lfNotToken) // skip extends
		var err error
		classDefinition.SuperClassName, err = cp.parseClassName()
		if err != nil {
			cp.parser.errs = append(cp.parser.errs, err)
			cp.consume(untilLc)
		}
	}
	if cp.parser.token.Type == lex.TokenImplements {
		cp.Next(lfNotToken) // skip key word
		classDefinition.InterfaceNames, err = cp.parseImplementsInterfaces()
		if err != nil {
			cp.parser.errs = append(cp.parser.errs, err)
			cp.consume(untilLc)
		}
	}
	cp.parser.ifTokenIsLfThenSkip()
	if cp.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{' but '%s'",
			cp.parser.errMsgPrefix(), cp.parser.token.Description)
		cp.parser.errs = append(cp.parser.errs, err)
		return nil, err
	}
	validAfterAccessControlToken := func(keyWord string) error {
		if cp.parser.token.Type == lex.TokenIdentifier ||
			cp.parser.token.Type == lex.TokenFn ||
			cp.parser.token.Type == lex.TokenStatic ||
			cp.parser.token.Type == lex.TokenSynchronized ||
			cp.parser.token.Type == lex.TokenFinal ||
			cp.parser.token.Type == lex.TokenAbstract {
			return nil
		}
		return fmt.Errorf("%s not a valid token after '%s'",
			cp.parser.errMsgPrefix(), keyWord)
	}
	validAfterVolatile := func(token *lex.Token) error {
		if token.Type == lex.TokenIdentifier {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'volatile'",
			cp.parser.errMsgPrefix())
	}
	validAfterAbstract := func() error {
		if cp.parser.token.Type == lex.TokenFn {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'abstract'",
			cp.parser.errMsgPrefix())
	}
	validAfterSynchronized := func() error {
		if cp.parser.token.Type == lex.TokenFn ||
			cp.parser.token.Type == lex.TokenFinal {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'synchronized'",
			cp.parser.errMsgPrefix())
	}
	validAfterStatic := func() error {
		if cp.parser.token.Type == lex.TokenIdentifier ||
			cp.parser.token.Type == lex.TokenFn ||
			cp.parser.token.Type == lex.TokenFinal {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'static'",
			cp.parser.errMsgPrefix())
	}
	validAfterFinal := func() error {
		if cp.parser.token.Type == lex.TokenFn ||
			cp.parser.token.Type == lex.TokenSynchronized {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'final'",
			cp.parser.errMsgPrefix())
	}
	cp.Next(lfNotToken) // skip {
	comment := &CommentParser{
		parser: cp.parser,
	}
	var (
		isStatic           bool
		isVolatile         bool
		isSynchronized     bool
		isFinal            bool
		accessControlToken *lex.Token
	)
	resetProperty := func() {
		isStatic = false
		isVolatile = false
		isSynchronized = false
		isFinal = false
		isAbstract = false
		accessControlToken = nil
	}
	for cp.parser.token.Type != lex.TokenEof {
		if len(cp.parser.errs) > cp.parser.nErrors2Stop {
			break
		}
		switch cp.parser.token.Type {
		case lex.TokenComment, lex.TokenMultiLineComment:
			comment.read()
		case lex.TokenRc:
			cp.Next(lfNotToken)
			return
		case lex.TokenSemicolon, lex.TokenLf:
			cp.Next(lfNotToken)
			continue
		case lex.TokenStatic:
			isStatic = true
			cp.Next(lfIsToken)
			cp.parser.unExpectNewLineAndSkip()
			if cp.parser.token.Type == lex.TokenLc {
				cp.Next(lfNotToken) // skip {
				block := &ast.Block{}
				cp.parser.BlockParser.parseStatementList(block, false)
				if cp.parser.token.Type != lex.TokenRc {
					cp.parser.errs = append(cp.parser.errs,
						fmt.Errorf("%s expect '}' , but '%s'", cp.parser.errMsgPrefix(),
							cp.parser.token.Description))
				} else {
					cp.Next(lfNotToken) // skip }
					classDefinition.StaticBlocks = append(classDefinition.StaticBlocks, block)
				}
				continue
			}
			err := validAfterStatic()
			if err != nil {
				cp.parser.errs = append(cp.parser.errs, err)
				isStatic = false
			}
		//access private
		case lex.TokenPublic, lex.TokenProtected, lex.TokenPrivate:
			accessControlToken = cp.parser.token
			cp.Next(lfIsToken)
			cp.parser.unExpectNewLineAndSkip()
			err = validAfterAccessControlToken(accessControlToken.Description)
			if err != nil {
				cp.parser.errs = append(cp.parser.errs, err)
				accessControlToken = nil // set to nil
			}
		case lex.TokenAbstract:
			cp.Next(lfIsToken)
			cp.parser.unExpectNewLineAndSkip()
			err = validAfterAbstract()
			if err != nil {
				cp.parser.errs = append(cp.parser.errs, err)
				accessControlToken = nil // set to nil
			} else {
				isAbstract = true
			}
		case lex.TokenVolatile:
			isVolatile = true
			cp.Next(lfIsToken)
			if err := validAfterVolatile(cp.parser.token); err != nil {
				cp.parser.errs = append(cp.parser.errs, err)
				isVolatile = false
			}
		case lex.TokenFinal:
			isFinal = true
			cp.Next(lfIsToken)
			if err := validAfterFinal(); err != nil {
				cp.parser.errs = append(cp.parser.errs, err)
				isFinal = false
			}
		case lex.TokenIdentifier:
			err = cp.parseField(classDefinition, &cp.parser.errs, isStatic, isVolatile, accessControlToken, comment)
			if err != nil {
				cp.consume(untilSemicolonOrLf)
				cp.Next(lfNotToken)
			}
			resetProperty()
		case lex.TokenConst: // const is for local use
			cp.Next(lfIsToken)
			err := cp.parseConst(classDefinition, comment)
			if err != nil {
				cp.consume(untilSemicolonOrLf)
				cp.Next(lfNotToken)
				continue
			}
		case lex.TokenSynchronized:
			isSynchronized = true
			cp.Next(lfIsToken)
			if err := validAfterSynchronized(); err != nil {
				cp.parser.errs = append(cp.parser.errs, err)
				isSynchronized = false
			}
		case lex.TokenFn:
			if isAbstract &&
				(classDefinition.IsAbstract() == false && classDefinition.IsInterface() == false) {
				cp.parser.errs = append(cp.parser.errs,
					fmt.Errorf("%s cannot  abstact method is non-abstract class",
						cp.parser.errMsgPrefix()))
			}
			isAbstract := isAbstract || isInterface
			f, err := cp.parser.FunctionParser.parse(true, isAbstract)
			if err != nil {
				resetProperty()
				cp.Next(lfNotToken)
				continue
			}
			f.Comment = comment.Comment
			if classDefinition.Methods == nil {
				classDefinition.Methods = make(map[string][]*ast.ClassMethod)
			}
			if f.Name == "" {
				f.Name = compileAutoName()
			}
			m := &ast.ClassMethod{}
			m.Function = f
			if accessControlToken != nil {
				switch accessControlToken.Type {
				case lex.TokenPrivate:
					m.Function.AccessFlags |= cg.ACC_METHOD_PRIVATE
				case lex.TokenProtected:
					m.Function.AccessFlags |= cg.ACC_METHOD_PROTECTED
				case lex.TokenPublic:
					m.Function.AccessFlags |= cg.ACC_METHOD_PUBLIC
				}
			}
			if isSynchronized {
				m.Function.AccessFlags |= cg.ACC_METHOD_SYNCHRONIZED
			}
			if isStatic {
				f.AccessFlags |= cg.ACC_METHOD_STATIC
			}
			if isAbstract {
				f.AccessFlags |= cg.ACC_METHOD_ABSTRACT
			}
			if isFinal {
				f.AccessFlags |= cg.ACC_METHOD_FINAL
			}
			if f.Name == classDefinition.Name && isInterface == false {
				f.Name = ast.SpecialMethodInit
			}
			classDefinition.Methods[f.Name] = append(classDefinition.Methods[f.Name], m)
			resetProperty()
		case lex.TokenImport:
			pos := cp.parser.mkPos()
			cp.parser.parseImports()
			cp.parser.errs = append(cp.parser.errs,
				fmt.Errorf("%s cannot have import at this scope",
					cp.parser.errMsgPrefix(pos)))
		default:
			cp.parser.errs = append(cp.parser.errs,
				fmt.Errorf("%s unexpected '%s'",
					cp.parser.errMsgPrefix(), cp.parser.token.Description))
			cp.Next(lfNotToken)
		}
	}
	return
}

func (cp *ClassParser) parseConst(classDefinition *ast.Class, comment *CommentParser) error {
	cs, err := cp.parser.parseConst()
	if err != nil {
		return err
	}
	constComment := comment.Comment
	if cp.parser.token.Type == lex.TokenComment {
		cp.Next(lfIsToken)
	} else {
		cp.parser.validStatementEnding()
	}
	if classDefinition.Block.Constants == nil {
		classDefinition.Block.Constants = make(map[string]*ast.Constant)
	}
	for _, v := range cs {
		if _, ok := classDefinition.Block.Constants[v.Name]; ok {
			cp.parser.errs = append(cp.parser.errs,
				fmt.Errorf("%s const %s alreay declared",
					cp.parser.errMsgPrefix(), v.Name))
			continue
		}
		classDefinition.Block.Constants[v.Name] = v
		v.Comment = constComment
	}
	return nil
}

func (cp *ClassParser) parseField(
	classDefinition *ast.Class,
	errs *[]error,
	isStatic bool,
	isVolatile bool,
	accessControlToken *lex.Token,
	comment *CommentParser) error {
	names, err := cp.parser.parseNameList()
	if err != nil {
		return err
	}
	t, err := cp.parser.parseType()
	if err != nil {
		return err
	}
	var initValues []*ast.Expression
	if cp.parser.token.Type == lex.TokenAssign {
		cp.parser.Next(lfNotToken) // skip = or :=
		initValues, err = cp.parser.ExpressionParser.parseExpressions(lex.TokenSemicolon)
		if err != nil {
			cp.consume(untilSemicolonOrLf)
		}
	}
	fieldComment := comment.Comment
	if cp.parser.token.Type == lex.TokenComment {
		cp.Next(lfIsToken)
	} else {
		cp.parser.validStatementEnding()
	}
	if classDefinition.Fields == nil {
		classDefinition.Fields = make(map[string]*ast.ClassField)
	}
	for k, v := range names {
		if _, ok := classDefinition.Fields[v.Name]; ok {
			cp.parser.errs = append(cp.parser.errs,
				fmt.Errorf("%s field %s is alreay declared",
					cp.parser.errMsgPrefix(), v.Name))
			continue
		}
		f := &ast.ClassField{}
		f.Name = v.Name
		f.Pos = v.Pos
		f.Type = t.Clone()
		f.AccessFlags = 0
		if k < len(initValues) {
			f.DefaultValueExpression = initValues[k]
		}
		f.Comment = fieldComment
		if isStatic {
			f.AccessFlags |= cg.ACC_FIELD_STATIC
		}
		if accessControlToken != nil {
			switch accessControlToken.Type {
			case lex.TokenPublic:
				f.AccessFlags |= cg.ACC_FIELD_PUBLIC
			case lex.TokenProtected:
				f.AccessFlags |= cg.ACC_FIELD_PROTECTED
			default: // private
				f.AccessFlags |= cg.ACC_FIELD_PRIVATE
			}
		}
		if isVolatile {
			f.AccessFlags |= cg.ACC_FIELD_VOLATILE
		}
		classDefinition.Fields[v.Name] = f
	}
	return nil
}

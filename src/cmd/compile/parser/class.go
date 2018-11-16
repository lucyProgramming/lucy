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

func (this *ClassParser) Next(lfIsToken bool) {
	this.parser.Next(lfIsToken)
}

func (this *ClassParser) consume(m map[lex.TokenKind]bool) {
	this.parser.consume(m)
}

func (this *ClassParser) parseClassName() (*ast.NameWithPos, error) {
	if this.parser.token.Type != lex.TokenIdentifier {
		err := fmt.Errorf("%s expect identifier for class`s name,but '%s'",
			this.parser.errMsgPrefix(), this.parser.token.Description)
		this.parser.errs = append(this.parser.errs, err)
		return nil, err
	}
	name := this.parser.token.Data.(string)
	pos := this.parser.mkPos()
	this.Next(lfNotToken)
	if this.parser.token.Type == lex.TokenSelection {
		this.Next(lfNotToken) // skip .
		if this.parser.token.Type != lex.TokenIdentifier {
			err := fmt.Errorf("%s expect identifer for class`s name,but '%s'",
				this.parser.errMsgPrefix(),
				this.parser.token.Description)
			this.parser.errs = append(this.parser.errs, err)
		} else {
			name += "." + this.parser.token.Data.(string)
			this.Next(lfNotToken) // skip name identifier
		}
	}
	return &ast.NameWithPos{
		Name: name,
		Pos:  pos,
	}, nil
}

func (this *ClassParser) parseImplementsInterfaces() ([]*ast.NameWithPos, error) {
	ret := []*ast.NameWithPos{}
	for this.parser.token.Type != lex.TokenEof {
		name, err := this.parseClassName()
		if err != nil {
			return nil, err
		}
		ret = append(ret, &ast.NameWithPos{
			Name: name.Name,
			Pos:  name.Pos,
		})
		if this.parser.token.Type == lex.TokenComma {
			this.Next(lfNotToken)
		} else {
			break
		}
	}
	return ret, nil
}

func (this *ClassParser) parse(isAbstract bool) (classDefinition *ast.Class, err error) {
	isInterface := this.parser.token.Type == lex.TokenInterface
	classDefinition = &ast.Class{}
	if isInterface {
		classDefinition.AccessFlags |= cg.AccClassInterface
		classDefinition.AccessFlags |= cg.AccClassAbstract
	}
	if isAbstract {
		classDefinition.AccessFlags |= cg.AccClassAbstract
	}
	this.Next(lfIsToken) // skip class key word
	this.parser.unExpectNewLineAndSkip()
	t, err := this.parseClassName()
	if t != nil {
		classDefinition.Name = t.Name
	}
	classDefinition.Block.IsClassBlock = true
	classDefinition.Block.Class = classDefinition
	if err != nil {
		if classDefinition.Name == "" {
			compileAutoName()
		}
		this.consume(untilLc)
	}
	classDefinition.Pos = this.parser.mkPos()
	if this.parser.token.Type == lex.TokenExtends { // parse father expression
		this.Next(lfNotToken) // skip extends
		var err error
		classDefinition.SuperClassName, err = this.parseClassName()
		if err != nil {
			this.parser.errs = append(this.parser.errs, err)
			this.consume(untilLc)
		}
	}
	if this.parser.token.Type == lex.TokenImplements {
		this.Next(lfNotToken) // skip key word
		classDefinition.InterfaceNames, err = this.parseImplementsInterfaces()
		if err != nil {
			this.parser.errs = append(this.parser.errs, err)
			this.consume(untilLc)
		}
	}
	this.parser.ifTokenIsLfThenSkip()
	if this.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{' but '%s'",
			this.parser.errMsgPrefix(), this.parser.token.Description)
		this.parser.errs = append(this.parser.errs, err)
		return nil, err
	}
	validAfterAccessControlToken := func(keyWord string) error {
		if this.parser.token.Type == lex.TokenIdentifier ||
			this.parser.token.Type == lex.TokenFn ||
			this.parser.token.Type == lex.TokenStatic ||
			this.parser.token.Type == lex.TokenSynchronized ||
			this.parser.token.Type == lex.TokenFinal ||
			this.parser.token.Type == lex.TokenAbstract {
			return nil
		}
		return fmt.Errorf("%s not a valid token after '%s'",
			this.parser.errMsgPrefix(), keyWord)
	}
	validAfterVolatile := func(token *lex.Token) error {
		if token.Type == lex.TokenIdentifier {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'volatile'",
			this.parser.errMsgPrefix())
	}
	validAfterAbstract := func() error {
		if this.parser.token.Type == lex.TokenFn {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'abstract'",
			this.parser.errMsgPrefix())
	}
	validAfterSynchronized := func() error {
		if this.parser.token.Type == lex.TokenFn ||
			this.parser.token.Type == lex.TokenFinal {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'synchronized'",
			this.parser.errMsgPrefix())
	}
	validAfterStatic := func() error {
		if this.parser.token.Type == lex.TokenIdentifier ||
			this.parser.token.Type == lex.TokenFn ||
			this.parser.token.Type == lex.TokenFinal {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'static'",
			this.parser.errMsgPrefix())
	}
	validAfterFinal := func() error {
		if this.parser.token.Type == lex.TokenFn ||
			this.parser.token.Type == lex.TokenSynchronized {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'final'",
			this.parser.errMsgPrefix())
	}
	this.Next(lfNotToken) // skip {
	comment := &CommentParser{
		parser: this.parser,
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
	for this.parser.token.Type != lex.TokenEof {
		if len(this.parser.errs) > this.parser.nErrors2Stop {
			break
		}
		switch this.parser.token.Type {
		case lex.TokenComment, lex.TokenMultiLineComment:
			comment.read()
		case lex.TokenRc:
			this.Next(lfNotToken)
			return
		case lex.TokenSemicolon, lex.TokenLf:
			this.Next(lfNotToken)
			continue
		case lex.TokenStatic:
			isStatic = true
			this.Next(lfIsToken)
			this.parser.unExpectNewLineAndSkip()
			if this.parser.token.Type == lex.TokenLc {
				this.Next(lfNotToken) // skip {
				block := &ast.Block{}
				this.parser.BlockParser.parseStatementList(block, false)
				if this.parser.token.Type != lex.TokenRc {
					this.parser.errs = append(this.parser.errs,
						fmt.Errorf("%s expect '}' , but '%s'", this.parser.errMsgPrefix(),
							this.parser.token.Description))
				} else {
					this.Next(lfNotToken) // skip }
					classDefinition.StaticBlocks = append(classDefinition.StaticBlocks, block)
				}
				continue
			}
			err := validAfterStatic()
			if err != nil {
				this.parser.errs = append(this.parser.errs, err)
				isStatic = false
			}
		//access private
		case lex.TokenPublic, lex.TokenProtected, lex.TokenPrivate:
			accessControlToken = this.parser.token
			this.Next(lfIsToken)
			this.parser.unExpectNewLineAndSkip()
			err = validAfterAccessControlToken(accessControlToken.Description)
			if err != nil {
				this.parser.errs = append(this.parser.errs, err)
				accessControlToken = nil // set to nil
			}
		case lex.TokenAbstract:
			this.Next(lfIsToken)
			this.parser.unExpectNewLineAndSkip()
			err = validAfterAbstract()
			if err != nil {
				this.parser.errs = append(this.parser.errs, err)
				accessControlToken = nil // set to nil
			} else {
				isAbstract = true
			}
		case lex.TokenVolatile:
			isVolatile = true
			this.Next(lfIsToken)
			if err := validAfterVolatile(this.parser.token); err != nil {
				this.parser.errs = append(this.parser.errs, err)
				isVolatile = false
			}
		case lex.TokenFinal:
			isFinal = true
			this.Next(lfIsToken)
			if err := validAfterFinal(); err != nil {
				this.parser.errs = append(this.parser.errs, err)
				isFinal = false
			}
		case lex.TokenIdentifier:
			err = this.parseField(classDefinition, &this.parser.errs, isStatic, isVolatile, accessControlToken, comment)
			if err != nil {
				this.consume(untilSemicolonOrLf)
				this.Next(lfNotToken)
			}
			resetProperty()
		case lex.TokenConst: // const is for local use
			this.Next(lfIsToken)
			err := this.parseConst(classDefinition, comment)
			if err != nil {
				this.consume(untilSemicolonOrLf)
				this.Next(lfNotToken)
				continue
			}
		case lex.TokenSynchronized:
			isSynchronized = true
			this.Next(lfIsToken)
			if err := validAfterSynchronized(); err != nil {
				this.parser.errs = append(this.parser.errs, err)
				isSynchronized = false
			}
		case lex.TokenFn:
			if isAbstract &&
				(classDefinition.IsAbstract() == false && classDefinition.IsInterface() == false) {
				this.parser.errs = append(this.parser.errs,
					fmt.Errorf("%s cannot  abstact method is non-abstract class",
						this.parser.errMsgPrefix()))
			}
			isAbstract := isAbstract || isInterface
			f, err := this.parser.FunctionParser.parse(true, isAbstract, true)
			if err != nil {
				resetProperty()
				this.Next(lfNotToken)
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
					m.Function.AccessFlags |= cg.AccMethodPrivate
				case lex.TokenProtected:
					m.Function.AccessFlags |= cg.AccMethodProtected
				case lex.TokenPublic:
					m.Function.AccessFlags |= cg.AccMethodPublic
				}
			}
			if isSynchronized {
				m.Function.AccessFlags |= cg.AccMethodSynchronized
			}
			if isStatic {
				f.AccessFlags |= cg.AccMethodStatic
			}
			if isAbstract {
				f.AccessFlags |= cg.AccMethodAbstract
			}
			if isFinal {
				f.AccessFlags |= cg.AccMethodFinal
			}
			if f.Name == classDefinition.Name && isInterface == false {
				f.Name = ast.SpecialMethodInit
			}
			classDefinition.Methods[f.Name] = append(classDefinition.Methods[f.Name], m)
			resetProperty()
		case lex.TokenImport:
			pos := this.parser.mkPos()
			this.parser.parseImports()
			this.parser.errs = append(this.parser.errs,
				fmt.Errorf("%s cannot have import at this scope",
					this.parser.errMsgPrefix(pos)))
		default:
			this.parser.errs = append(this.parser.errs,
				fmt.Errorf("%s unexpected '%s'",
					this.parser.errMsgPrefix(), this.parser.token.Description))
			this.Next(lfNotToken)
		}
	}
	return
}

func (this *ClassParser) parseConst(classDefinition *ast.Class, comment *CommentParser) error {
	cs, err := this.parser.parseConst()
	if err != nil {
		return err
	}
	constComment := comment.Comment
	if this.parser.token.Type == lex.TokenComment {
		this.Next(lfIsToken)
	} else {
		this.parser.validStatementEnding()
	}
	if classDefinition.Block.Constants == nil {
		classDefinition.Block.Constants = make(map[string]*ast.Constant)
	}
	for _, v := range cs {
		if _, ok := classDefinition.Block.Constants[v.Name]; ok {
			this.parser.errs = append(this.parser.errs,
				fmt.Errorf("%s const %s alreay declared",
					this.parser.errMsgPrefix(), v.Name))
			continue
		}
		classDefinition.Block.Constants[v.Name] = v
		v.Comment = constComment
	}
	return nil
}

func (this *ClassParser) parseField(
	classDefinition *ast.Class,
	errs *[]error,
	isStatic bool,
	isVolatile bool,
	accessControlToken *lex.Token,
	comment *CommentParser) error {
	names, err := this.parser.parseNameList()
	if err != nil {
		return err
	}
	t, err := this.parser.parseType()
	if err != nil {
		return err
	}
	var initValues []*ast.Expression
	if this.parser.token.Type == lex.TokenAssign {
		this.parser.Next(lfNotToken) // skip = or :=
		initValues, err = this.parser.ExpressionParser.parseExpressions(lex.TokenSemicolon)
		if err != nil {
			this.consume(untilSemicolonOrLf)
		}
	}
	fieldComment := comment.Comment
	if this.parser.token.Type == lex.TokenComment {
		this.Next(lfIsToken)
	} else {
		this.parser.validStatementEnding()
	}
	if classDefinition.Fields == nil {
		classDefinition.Fields = make(map[string]*ast.ClassField)
	}
	for k, v := range names {
		if _, ok := classDefinition.Fields[v.Name]; ok {
			this.parser.errs = append(this.parser.errs,
				fmt.Errorf("%s field %s is alreay declared",
					this.parser.errMsgPrefix(), v.Name))
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
			f.AccessFlags |= cg.AccFieldStatic
		}
		if accessControlToken != nil {
			switch accessControlToken.Type {
			case lex.TokenPublic:
				f.AccessFlags |= cg.AccFieldPublic
			case lex.TokenProtected:
				f.AccessFlags |= cg.AccFieldProtected
			default: // private
				f.AccessFlags |= cg.AccFieldPrivate
			}
		}
		if isVolatile {
			f.AccessFlags |= cg.AccFieldVolatile
		}
		classDefinition.Fields[v.Name] = f
	}
	return nil
}

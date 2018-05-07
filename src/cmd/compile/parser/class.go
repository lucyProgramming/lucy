package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type Class struct {
	parser             *Parser
	classDefinition    *ast.Class
	isStatic           bool
	accessControlToken *lex.Token
}

func (c *Class) Next() {
	c.parser.Next()
}

func (c *Class) consume(m map[int]bool) {
	c.parser.consume(m)
}

func (c *Class) parseClassName() (string, error) {
	if c.parser.token.Type != lex.TOKEN_IDENTIFIER {
		err := fmt.Errorf("%s expect class`s name,but '%s'",
			c.parser.errorMsgPrefix(), c.parser.token.Desp)
		c.parser.errs = append(c.parser.errs, err)
		return "", err
	}
	name := c.parser.token.Data.(string)
	c.Next()
	if c.parser.token.Type == lex.TOKEN_DOT {
		c.Next()
		if c.parser.token.Type != lex.TOKEN_IDENTIFIER {
			err := fmt.Errorf("%s expect identifer,but '%s'", c.parser.errorMsgPrefix(),
				c.parser.token.Desp)
			c.parser.errs = append(c.parser.errs, err)
		}
		name += "." + c.parser.token.Data.(string)
		c.Next() // skip name identifier
	}
	return name, nil
}

func (c *Class) parseInterfaces() ([]*ast.NameWithPos, error) {
	ret := []*ast.NameWithPos{}
	for c.parser.token.Type != lex.TOKEN_EOF {
		pos := c.parser.mkPos()
		name, err := c.parseClassName()
		if err != nil {
			return nil, err
		}
		ret = append(ret, &ast.NameWithPos{
			Name: name,
			Pos:  pos,
		})
		if c.parser.token.Type == lex.TOKEN_COMMA {
			c.Next()
		} else {
			break
		}
	}
	return ret, nil
}

func (c *Class) parse() (classDefinition *ast.Class, err error) {
	defer c.resetProperty()
	classDefinition = &ast.Class{}
	c.classDefinition = classDefinition
	c.Next() // skip class key word
	c.classDefinition.Pos = c.parser.mkPos()
	c.classDefinition.Name, err = c.parseClassName()
	c.classDefinition.Block.IsClassBlock = true
	if err != nil {
		return nil, err
	}
	if c.parser.token.Type == lex.TOKEN_EXTENDS { // parse father expression
		c.Next() // skip extends
		c.classDefinition.Pos = c.parser.mkPos()
		if c.parser.token.Type != lex.TOKEN_IDENTIFIER {
			err = fmt.Errorf("%s class`s father must be a identifier", c.parser.errorMsgPrefix())
			c.parser.errs = append(c.parser.errs, err)
			c.consume(untils_lc) //
		} else {
			t, err := c.parseClassName()
			c.classDefinition.SuperClassName = t
			if err != nil {
				c.parser.errs = append(c.parser.errs, err)
				return nil, err
			}
		}
	}
	if c.parser.token.Type == lex.TOKEN_IMPLEMENTS {
		c.Next() // skip key word
		c.classDefinition.InterfaceNames, err = c.parseInterfaces()
		if err != nil {
			c.consume(untils_lc)
		}
	}
	if c.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s expect '{' but '%s'", c.parser.errorMsgPrefix(), c.parser.token.Desp)
		c.parser.errs = append(c.parser.errs, err)
		return nil, err
	}
	c.Next() // skip {
	c.resetProperty()
	validAfterPublic := func(keyWord string, token *lex.Token) error {
		if token.Type == lex.TOKEN_IDENTIFIER ||
			token.Type == lex.TOKEN_FUNCTION ||
			token.Type == lex.TOKEN_STATIC {
			return nil
		}
		return fmt.Errorf("%s not a valid token after '%s'", c.parser.errorMsgPrefix(), keyWord)
	}
	validAfterStatic := func(token *lex.Token) error {
		if token.Type == lex.TOKEN_IDENTIFIER ||
			token.Type == lex.TOKEN_FUNCTION {
			return nil
		}
		return fmt.Errorf("%s not a valid token after 'static'", c.parser.errorMsgPrefix())
	}
	for c.parser.token.Type != lex.TOKEN_EOF {
		if len(c.parser.errs) > c.parser.nerr {
			break
		}
		switch c.parser.token.Type {
		case lex.TOKEN_RC:
			c.Next()
			return
		case lex.TOKEN_SEMICOLON:
			c.Next()
			continue
		case lex.TOKEN_STATIC:
			c.isStatic = true
			c.Next()
			err := validAfterStatic(c.parser.token)
			if err != nil {
				c.parser.errs = append(c.parser.errs, err)
			}
		//access private
		case lex.TOKEN_PUBLIC:
			c.accessControlToken = c.parser.token
			c.Next()
			err = validAfterPublic("public", c.parser.token)
			if err != nil {
				c.parser.errs = append(c.parser.errs, err)
			}
		case lex.TOKEN_PROTECTED:
			c.accessControlToken = c.parser.token
			c.Next()
			err = validAfterPublic("protected", c.parser.token)
			if err != nil {
				c.parser.errs = append(c.parser.errs, err)
			}
		case lex.TOKEN_PRIVATE:
			c.accessControlToken = c.parser.token
			c.Next() // skip private
			err = validAfterPublic("private", c.parser.token)
			if err != nil {
				c.parser.errs = append(c.parser.errs, err)
			}
		case lex.TOKEN_IDENTIFIER:
			err = c.parseField(&c.parser.errs)
			if err != nil {
				c.consume(untils_semicolon)
				c.Next()
			}
			c.resetProperty()
		case lex.TOKEN_CONST: // const is for local use
			c.Next()
			err := c.parseConst()
			if err != nil {
				c.consume(untils_semicolon)
				c.Next()
				continue
			}
			c.resetProperty()
		case lex.TOKEN_FUNCTION:
			f, err := c.parser.Function.parse(true)
			if err != nil {
				c.consume(untils_rc)
				c.Next()
				c.resetProperty()
				continue
			}
			if c.classDefinition.Methods == nil {
				c.classDefinition.Methods = make(map[string][]*ast.ClassMethod)
			}
			if f.Name == "" {
				c.parser.errs = append(c.parser.errs, fmt.Errorf("%s method has no name", c.parser.errorMsgPrefix(f.Pos)))
				c.resetProperty()
				continue
			}
			m := &ast.ClassMethod{}
			m.Func = f
			f.AccessFlags = 0
			if c.accessControlToken != nil {
				switch c.accessControlToken.Type {
				case lex.TOKEN_PRIVATE:
					m.Func.AccessFlags |= cg.ACC_METHOD_PRIVATE
				case lex.TOKEN_PROTECTED:
					m.Func.AccessFlags |= cg.ACC_METHOD_PROTECTED
				case lex.TOKEN_PUBLIC:
					m.Func.AccessFlags |= cg.ACC_METHOD_PUBLIC
				}
			} else {
				m.Func.AccessFlags |= cg.ACC_METHOD_PRIVATE
			}
			if c.isStatic {
				f.AccessFlags |= cg.ACC_METHOD_STATIC
			}

			if c.classDefinition.Methods == nil {
				c.classDefinition.Methods = make(map[string][]*ast.ClassMethod)
			}
			if f.Name == c.classDefinition.Name {
				f.Name = ast.CONSTRUCTION_METHOD_NAME
			}
			c.classDefinition.Methods[f.Name] = append(c.classDefinition.Methods[f.Name], m)
			c.resetProperty()
		default:
			c.parser.errs = append(c.parser.errs, fmt.Errorf("%s unexpect token:%s", c.parser.errorMsgPrefix(), c.parser.token.Desp))
			c.Next()
		}
	}
	return
}

func (c *Class) resetProperty() {
	c.isStatic = false
	c.accessControlToken = nil
}

func (c *Class) parseConst() error {
	pos := c.parser.mkPos()
	vs, es, typ, err := c.parser.parseConstDefinition(false)
	if err != nil {
		return err
	}
	if typ != lex.TOKEN_ASSIGN {
		c.parser.errs = append(c.parser.errs,
			fmt.Errorf("%s declare const should use ‘=’ instead of ‘:=’",
				c.parser.errorMsgPrefix(vs[0].Pos)))
	}
	if len(vs) != len(es) {
		c.parser.errs = append(c.parser.errs,
			fmt.Errorf("%s cannot assign %d values to %d destinations",
				c.parser.errorMsgPrefix(pos), len(es), len(vs)))
	}

	if c.classDefinition.Block.Consts == nil {
		c.classDefinition.Block.Consts = make(map[string]*ast.Const)
	}
	for k, v := range vs {
		if _, ok := c.classDefinition.Block.Consts[v.Name]; ok {
			c.parser.errs = append(c.parser.errs, fmt.Errorf("%s const %s alreay declared", v.Name))
			continue
		}
		if k < len(es) {
			t := &ast.Const{}
			t.VariableDefinition = *v
			t.Expression = es[k]
			c.classDefinition.Block.Consts[v.Name] = t
		}
	}
	return nil
}

func (c *Class) parseField(errs *[]error) error {
	names, err := c.parser.parseNameList()
	if err != nil {
		return err
	}
	t, err := c.parser.parseType()
	if err != nil {
		return err
	}
	var es []*ast.Expression
	if c.parser.token.Type == lex.TOKEN_ASSIGN {
		c.Next() // skip =
		es, err = c.parser.ExpressionParser.parseExpressions()
		if err != nil {
			*errs = append(*errs, err)
		}
	}

	if c.parser.token.Type != lex.TOKEN_SEMICOLON {
		*errs = append(*errs, fmt.Errorf("%s missing senicolon after field definition",
			c.parser.errorMsgPrefix()))
	}
	if c.classDefinition.Fields == nil {
		c.classDefinition.Fields = make(map[string]*ast.ClassField)
	}
	if es != nil && len(es) != len(names) {
		err := fmt.Errorf("%s cannot assign %d values to %d destinations",
			c.parser.errorMsgPrefix(), len(es), len(names))
		*errs = append(*errs, err)
	}
	for k, v := range names {
		if _, ok := c.classDefinition.Fields[v.Name]; ok {
			c.parser.errs = append(c.parser.errs,
				fmt.Errorf("%s field %s is alreay declared",
					c.parser.errorMsgPrefix(), v.Name))
			continue
		}
		f := &ast.ClassField{}
		f.Name = v.Name
		f.Pos = v.Pos
		f.Typ = &ast.VariableType{}
		*f.Typ = *t
		f.AccessFlags = 0
		if k < len(es) && es[k] != nil {
			f.Expression = es[k]
		}
		if c.isStatic {
			f.AccessFlags |= cg.ACC_FIELD_STATIC
		}
		if c.accessControlToken == nil {
			f.AccessFlags |= cg.ACC_FIELD_PRIVATE
		} else {
			switch c.accessControlToken.Type {
			case lex.TOKEN_PUBLIC:
				f.AccessFlags |= cg.ACC_FIELD_PUBLIC
			case lex.TOKEN_PROTECTED:
				f.AccessFlags |= cg.ACC_FIELD_PROTECTED
			default: // private
				f.AccessFlags |= cg.ACC_FIELD_PRIVATE
			}
		}
		c.classDefinition.Fields[v.Name] = f
	}
	return nil
}

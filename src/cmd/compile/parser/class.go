package parser

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
	"github.com/756445638/lucy/src/cmd/compile/lex"
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
		err := fmt.Errorf("%s on name after class,but %s", c.parser.errorMsgPrefix(), c.parser.token.Desp)
		c.parser.errs = append(c.parser.errs, err)
		return "", err
	}
	name := c.parser.token.Data.(string)
	c.Next()
	return name, nil

}

func (c *Class) parse() (classDefinition *ast.Class, err error) {
	c.resetProperty()
	defer c.resetProperty()
	classDefinition = &ast.Class{}
	c.classDefinition = classDefinition
	c.Next() // skip class key word
	c.classDefinition.Name, err = c.parseClassName()
	if err != nil {
		return nil, err
	}
	c.classDefinition.Pos = c.parser.mkPos()
	c.classDefinition.Block.IsClassBlock = true
	if c.parser.eof {
		err = c.parser.mkUnexpectedEofErr()
		c.parser.errs = append(c.parser.errs, err)
		return nil, err
	}
	if c.parser.token.Type == lex.TOKEN_EXTENDS { // parse father expression
		c.Next() // skip extends
		if c.parser.token.Type != lex.TOKEN_IDENTIFIER {
			err = fmt.Errorf("%s class`s father must be a identifier", c.parser.errorMsgPrefix())
			c.parser.errs = append(c.parser.errs, err)
			c.consume(untils_lc) //
		} else {
			t, err := c.parser.parseType()
			c.classDefinition.Name = t.Name
			if err != nil {
				c.parser.errs = append(c.parser.errs, err)
				return nil, err
			}
		}
	}
	if c.parser.token.Type == lex.TOKEN_IMPLEMENTS {
		c.Next() // skip implements
		for c.parser.token.Type != lex.TOKEN_LC {
			c.Next() // skip all until  {
		}
	}
	if c.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s expect '{' but '%s'", c.parser.errorMsgPrefix(), c.parser.token.Desp)
		c.parser.errs = append(c.parser.errs, err)
		return nil, err
	}
	c.Next() // skip {
	c.resetProperty()
	validAfterPublic := func(token *lex.Token) error {
		if token.Type == lex.TOKEN_IDENTIFIER ||
			token.Type == lex.TOKEN_FUNCTION ||
			token.Type == lex.TOKEN_STATIC {
			return nil
		}
		return fmt.Errorf("%s not a valid token after (public|private|protected)", c.parser.errorMsgPrefix())
	}
	for !c.parser.eof {
		if len(c.parser.errs) > c.parser.nerr {
			break
		}
		switch c.parser.token.Type {
		case lex.TOKEN_SEMICOLON:
			c.Next()
			continue
		case lex.TOKEN_STATIC:
			c.isStatic = true
			c.Next()
			err := validAfterPublic(c.parser.token)
			if err != nil {
				c.parser.errs = append(c.parser.errs, err)
				c.Next()
			}
		//access private
		case lex.TOKEN_PUBLIC:
			c.accessControlToken = c.parser.token
			c.Next()
			if err != nil {
				c.parser.errs = append(c.parser.errs, err)
				c.Next()
			}
		case lex.TOKEN_PROTECTED:
			c.accessControlToken = c.parser.token
			c.Next()
			if err != nil {
				c.parser.errs = append(c.parser.errs, err)
				c.Next()
			}
		case lex.TOKEN_PRIVATE:
			c.accessControlToken = c.parser.token
			c.Next()
		case lex.TOKEN_IDENTIFIER:
			err = c.parseFiled()
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
		case lex.TOKEN_FUNCTION:
			f, err := c.parser.Function.parse(false)
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
			if c.isStatic {
				m.Func.AccessFlags |= cg.ACC_METHOD_STATIC
			}
			if c.accessControlToken == nil {
				m.Func.AccessFlags |= cg.ACC_METHOD_PRIVATE
			} else {
				switch c.accessControlToken.Type {
				case lex.TOKEN_PRIVATE:
					m.Func.AccessFlags |= cg.ACC_METHOD_PRIVATE
				case lex.TOKEN_PUBLIC:
					m.Func.AccessFlags |= cg.ACC_METHOD_PUBLIC
				case lex.TOKEN_PROTECTED:
					m.Func.AccessFlags |= cg.ACC_METHOD_PROTECTED
				}
			}
			if f.Name == c.classDefinition.Name {
				if c.classDefinition.Constructors == nil {
					c.classDefinition.Constructors = []*ast.ClassMethod{m}
				} else {
					c.classDefinition.Constructors = append(c.classDefinition.Constructors, m)
				}
			} else {
				if c.classDefinition.Methods == nil {
					c.classDefinition.Methods = make(map[string][]*ast.ClassMethod)
				}
				c.classDefinition.Methods[f.Name] = append(c.classDefinition.Methods[f.Name], m)
			}
			c.resetProperty()
		case lex.TOKEN_RC:
			c.Next()
			return
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
	vs, es, typ, err := c.parser.parseConstDefinition()
	if err != nil {
		return err
	}
	if typ != lex.TOKEN_ASSIGN {
		c.parser.errs = append(c.parser.errs,
			fmt.Errorf("%s declare const should use ‘=’ instead of ‘:=’", c.parser.errorMsgPrefix(vs[0].Pos)))
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

func (c *Class) parseFiled() error {
	names, err := c.parser.parseNameList()
	if err != nil {
		return err
	}
	t, err := c.parser.parseType()
	if err != nil {
		return err
	}
	if c.classDefinition.Fields == nil {
		c.classDefinition.Fields = make(map[string]*ast.ClassField)
	}
	for _, v := range names {
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
		if c.isStatic {
			f.AccessFlags |= cg.ACC_FIELD_STATIC
		}
		if c.accessControlToken == nil {
			f.AccessFlags |= cg.ACC_FIELD_PRIVATE
		} else {
			switch c.accessControlToken.Type {
			case lex.TOKEN_PUBLIC:
				f.AccessFlags |= cg.ACC_FIELD_PUBLIC
			case lex.TOKEN_PRIVATE:
				f.AccessFlags |= cg.ACC_FIELD_PRIVATE
			case lex.TOKEN_PROTECTED:
				f.AccessFlags |= cg.ACC_FIELD_PROTECTED
			}
		}
		c.classDefinition.Fields[v.Name] = f
	}
	return nil
}

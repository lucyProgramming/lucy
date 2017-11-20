package parser

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

type Class struct {
	parser          *Parser
	classDefinition *ast.Class
	access          int
	isStatic        bool
	isConst         bool
}

func (c *Class) Next() {
	c.parser.Next()
}

func (c *Class) consume(m map[int]bool) {
	c.parser.consume(m)
}

func (c *Class) parse() (classDefinition *ast.Class, err error) {
	classDefinition = &ast.Class{}
	c.classDefinition = classDefinition
	c.Next() // skip class key work
	if c.parser.token.Type != lex.TOKEN_IDENTIFIER {
		err = fmt.Errorf("%s on name after class,but %s", c.parser.errorMsgPrefix(), c.parser.token.Desp)
		c.parser.errs = append(c.parser.errs, err)
		return nil, err
	}
	c.classDefinition.Name = c.parser.token.Data.(string)
	c.classDefinition.Pos = c.parser.mkPos()
	c.Next() // skip class name
	if c.parser.eof {
		err = c.parser.mkUnexpectedEofErr()
		c.parser.errs = append(c.parser.errs, err)
		return nil, err
	}
	if c.parser.token.Type == lex.TOKEN_COLON { // parse father expression
		c.Next() // skip :
		if c.parser.token.Type != lex.TOKEN_IDENTIFIER {
			err = fmt.Errorf("%s class`s father must be a identifier", c.parser.errorMsgPrefix())
			c.parser.errs = append(c.parser.errs, err)
			c.consume(untils_lc) //
		} else {
			c.classDefinition.Father, err = c.parser.ExpressionParser.parseIdentifierExpression()
			if err != nil {
				c.parser.errs = append(c.parser.errs, err)
				return nil, err
			}
		}
	}
	if c.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s expect { but %s", c.parser.errorMsgPrefix(), c.parser.token.Desp)
		c.parser.errs = append(c.parser.errs, err)
		return nil, err
	}
	c.Next() // skip {
	c.access = ast.ACCESS_PRIVATE
	c.isStatic = false
	for !c.parser.eof {
		switch c.parser.token.Type {
		case lex.TOKEN_SEMICOLON:
			c.Next()
			continue
		case lex.TOKEN_STATIC:
			c.isStatic = true
			c.Next()
		//access private
		case lex.TOKEN_PUBLIC:
			c.access = ast.ACCESS_PUBLIC
			c.Next()
		case lex.TOKEN_PROTECTED:
			c.access = ast.ACCESS_PROTECTED
			c.Next()
		case lex.TOKEN_PRIVATE:
			c.access = ast.ACCESS_PRIVATE
			c.Next()
		case lex.TOKEN_IDENTIFIER:
			err = c.parseFiled()
			if err != nil {
				c.consume(untils_semicolon)
				c.Next()
			}
			c.resetProperty()
		case lex.TOKEN_CONST: // const is for local use
			if c.access == ast.ACCESS_PUBLIC {
				c.parser.errs = append(c.parser.errs, fmt.Errorf("%s const declared in class,can only be used this block"))
			}
			c.isConst = true
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
				c.classDefinition.Methods = make(map[string]*ast.ClassMethod)
			}
			if f.Name == "" {
				c.parser.errs = append(c.parser.errs, fmt.Errorf("%s method has no name", c.parser.errorMsgPrefix(f.Pos)))
				c.resetProperty()
				continue
			}
			if _, ok := c.classDefinition.Methods[f.Name]; ok || (f.Name == c.classDefinition.Name && c.classDefinition.Constructor != nil) {
				c.parser.errs = append(c.parser.errs, fmt.Errorf("%s method　%s already declared", c.parser.errorMsgPrefix(f.Pos), f.Name))
				c.resetProperty()
				continue
			}
			m := &ast.ClassMethod{}
			m.ClassFieldProperty.Access = c.access
			m.IsStatic = c.isStatic
			m.Func = f
			c.resetProperty()
			if f.Name == c.classDefinition.Name {
				c.classDefinition.Constructor = m
			} else {
				c.classDefinition.Methods[f.Name] = m
			}
		case lex.TOKEN_RC:
			c.Next()
			break
		default:
			c.parser.errs = append(c.parser.errs, fmt.Errorf("%s unexcept token:%s", c.parser.errorMsgPrefix(), c.parser.token.Desp))
			c.Next()
		}
	}
	return
}

func (c *Class) resetProperty() {
	c.access = ast.ACCESS_PRIVATE
	c.isStatic = false
	c.isConst = false
}

func (c *Class) parseConst() error {
	names, _, es, err := c.parser.parseAssignedNames()
	if err != nil {
		return err
	}
	if c.classDefinition.Consts == nil {
		c.classDefinition.Consts = make(map[string]*ast.Const)
	}
	for k, v := range names {
		if _, ok := c.classDefinition.Consts[v.Name]; ok {
			c.parser.errs = append(c.parser.errs, fmt.Errorf("%s const %s alreay declared", v.Name))
			continue
		}
		c.classDefinition.Consts[v.Name] = &ast.Const{}
		c.classDefinition.Consts[v.Name].Pos = v.Pos
		c.classDefinition.Consts[v.Name].Expression = es[k]
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
		f.ClassFieldProperty.Access = c.access
		f.IsStatic = c.isStatic
		c.classDefinition.Fields[v.Name] = f
	}
	return nil
}

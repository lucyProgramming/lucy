package parser

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

type Class struct {
	parser          *Parser
	token           *lex.Token
	classDefinition *ast.Class
	access          int
	isStatic        bool
	isConst         bool
}

func (c *Class) Next() {
	c.parser.Next()
	c.token = c.parser.token
}

func (c *Class) consume(untils ...int) {
	c.parser.consume(untils...)
}

func (c *Class) parse(ispublic bool) (classDefinition *ast.Class, err error) {
	classDefinition = &ast.Class{}
	c.Next()
	if c.parser.eof {
		return nil, c.parser.mkUnexpectedEofErr()
	}
	if c.token.Type != lex.TOKEN_IDENTIFIER {
		c.consume(lex.TOKEN_SEMICOLON, lex.TOKEN_RC)
		c.Next()
		return nil, fmt.Errorf("%s on name after class", c.parser.errorMsgPrefix())
	}
	c.classDefinition.Name = c.token.Data.(string)
	c.classDefinition.Pos = c.parser.mkPos()
	c.Next()
	if c.parser.eof {
		return nil, c.parser.mkUnexpectedEofErr()
	}
	var father *ast.Expression
	if c.token.Type == lex.TOKEN_COLON {
		c.Next()
		if c.parser.eof {
			return nil, c.parser.mkUnexpectedEofErr()
		}
		if c.token.Type != lex.TOKEN_IDENTIFIER {
			c.consume(lex.TOKEN_SEMICOLON, lex.TOKEN_RC)
			c.Next()
			return nil, fmt.Errorf("%s class`s father must be a identifier", c.parser.errorMsgPrefix())
		}
		father, err = c.parser.ExpressionParser.parseIdentifierExpression()
		if err != nil {
			c.consume(lex.TOKEN_SEMICOLON, lex.TOKEN_RC)
			c.Next()
			return
		}
		c.Next()
	}
	if c.token.Type != lex.TOKEN_LC {
		return nil, fmt.Errorf("%s except } but %s", c.parser.errorMsgPrefix(), c.token.Desp)
	}
	c.access = ast.ACCESS_PRIVATE
	c.isStatic = false
	for !c.parser.eof {
		switch c.token.Type {
		case lex.TOKEN_STATIC:
			c.isStatic = true
		//access private
		case lex.TOKEN_PUBLIC:
			c.access = ast.ACCESS_PUBLIC
		case lex.TOKEN_PROTECTED:
			c.access = ast.ACCESS_PROTECTED
		case lex.TOKEN_PRIVATE:
			c.access = ast.ACCESS_PRIVATE
		case lex.TOKEN_IDENTIFIER:
			err = c.parseFiled()
			if err != nil {
				c.parser.errs = append(c.parser.errs, err)
				c.consume(lex.TOKEN_SEMICOLON)
				c.Next()
			}
			c.resetProperty()
		case lex.TOKEN_CONST:
			c.isConst = true
			c.Next()
			err := c.parseConst()
			if err != nil {
				c.consume(lex.TOKEN_SEMICOLON)
				continue
			}
		case lex.TOKEN_FUNCTION:
			c.Next()
			f, err := c.parser.Function.parse(false)
			if err != nil {
				c.consume(lex.TOKEN_RC)
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
			m.Access = c.access
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
			c.parser.errs = append(c.parser.errs, fmt.Errorf("%s unexcept token(%s)", c.parser.errorMsgPrefix(), c.token.Desp))
		}
	}
	c.classDefinition.Father = father
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
		c.classDefinition.Consts[v.Name].Init = es[k]
	}
	return nil
}

func (c *Class) parseFiled() error {
	if c.parser.eof {
		return c.parser.mkUnexpectedEofErr()
	}
	names, err := c.parser.parseNameList()
	if err != nil {
		return err
	}
	t, err := c.parser.parseType()
	if t != nil {
		return err
	}
	for _, v := range names {
		if c.classDefinition.Fields == nil {
			c.classDefinition.Fields = make(map[string]*ast.ClassField)
			if _, ok := c.classDefinition.Fields[v.Name]; ok {
				c.parser.errs = append(c.parser.errs,
					fmt.Errorf("%s field %s is alreay declared",
						c.parser.errorMsgPrefix(), v.Name))
			}
		}
	}
	return nil
}

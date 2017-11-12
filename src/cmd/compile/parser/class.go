package parser

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

type Class struct {
	parser          *Parser
	token           lex.Token
	classDefinition *ast.Class
	access          int
	isstatic        bool
}

func (c *Class) Next() {
	c.parser.Next()
	c.token = c.parser.token
}

func (c *Class) consume(untils ...int) {
	c.parser.consume()
}

func (c *Class) parse(ispublic bool) (classDefinition *ast.Class, err error) {
	classDefinition = &ast.Class{}
	c.Next()
	if c.parser.eof {
		return nil, c.parser.mkUnexpectedErr()
	}
	if c.token.Type != lex.TOKEN_IDENTIFIER {
		c.consume(lex.TOKEN_SEMICOLON, lex.TOKEN_RC)
		c.Next()
		return nil, fmt.Errorf("%s on name after class", c.parser.errorMsgPrefix())
	}
	name := &ast.NameWithPos{
		Name: c.token.Data.(string),
		Pos:  c.parser.mkPos(),
	}
	c.Next()
	if c.parser.eof {
		return nil, c.parser.mkUnexpectedErr()
	}
	var father *ast.Expression
	if c.token.Type == lex.TOKEN_COLON {
		c.Next()
		if c.parser.eof {
			return nil, c.parser.mkUnexpectedErr()
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
	c.isstatic = false
	var err error
	for !c.parser.eof {
		switch c.token.Type {
		case lex.TOKEN_STATIC:
			c.isstatic = true
		//access private
		case lex.TOKEN_PUBLIC:
			access = ast.ACCESS_PUBLIC
		case lex.TOKEN_PROTECTED:
			access = ast.ACCESS_PROTECTED
		case lex.TOKEN_PRIVATE:
			access = ast.ACCESS_PRIVATE
		case lex.TOKEN_IDENTIFIER:
			err = c.parseFiled(access)
			c.resetProperty()
			if err != nil {
				c.consume()
			}
		case lex.TOKEN_FUNCTION:

		}
	}
	return
}

func (c *Class) resetProperty() {
	c.access = ast.ACCESS_PRIVATE
	c.isstatic = false
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
	for _, v := range names {
		if c.classDefinition.Fields == nil {
			c.classDefinition.Fields = make(map[string]*ClassField)
			if _, ok := c.classDefinition.Fields[v.Name]; ok {
				c.parser.errs = append(c.parser.errs,
					fmt.Errorf("%s field %s is alreay declared",
						c.parser.errorMsgPrefix(), v.Name))
			}
		}
	}
	return nil
}

func (c *Class) parseMethod() error {
	names, err := c.parser.parseNameList()
	if err != nil {
		return err
	}
	c.parser.parseType()
	return err
}

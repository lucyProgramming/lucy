package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type Interface struct {
	parser             *Parser
	classDefinition    *ast.Class
	isStatic           bool
	accessControlToken *lex.Token
}

func (c *Interface) Next() {
	c.parser.Next()
}

func (c *Interface) consume(m map[int]bool) {
	c.parser.consume(m)
}

func (c *Interface) parse() (classDefinition *ast.Class, err error) {
	c.Next() // skip interface key word
	classDefinition = &ast.Class{}
	c.classDefinition = classDefinition
	c.classDefinition.Pos = c.parser.mkPos()
	c.classDefinition.Block.IsClassBlock = true
	c.classDefinition.AccessFlags |= cg.ACC_CLASS_INTERFACE // interface
	c.classDefinition.AccessFlags |= cg.ACC_CLASS_ABSTRACT
	c.classDefinition.Name, err = c.parser.Class.parseClassName()
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
			t, err := c.parser.Class.parseClassName()
			c.classDefinition.SuperClassName = t
			if err != nil {
				c.parser.errs = append(c.parser.errs, err)
				return nil, err
			}
		}
	}
	if c.parser.token.Type == lex.TOKEN_IMPLEMENTS {
		c.Next() // skip key word
		c.classDefinition.InterfaceNames, err = c.parser.Class.parseInterfaces()
		if err != nil {
			c.consume(untils_lc)
		}
	}
	if c.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s expect '{' but '%s'", c.parser.errorMsgPrefix(), c.parser.token.Desp)
		c.parser.errs = append(c.parser.errs, err)
		return nil, err
	}
	c.Next()
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
		case lex.TOKEN_FUNCTION:
			c.Next() /// skip key word
			var name string
			if c.parser.token.Type != lex.TOKEN_IDENTIFIER {
				c.parser.errs = append(c.parser.errs, fmt.Errorf("%s expect function name,but '%s'", c.parser.errorMsgPrefix(), c.parser.token.Desp))
				c.consume(untils_rc)
				c.Next()
				continue
			}
			name = c.parser.token.Data.(string)
			c.Next() // skip name
			functionType, err := c.parser.parseFunctionType()
			if err != nil {
				c.consume(untils_rc)
				c.Next()
				continue
			}
			if c.classDefinition.Methods == nil {
				c.classDefinition.Methods = make(map[string][]*ast.ClassMethod)
			}
			m := &ast.ClassMethod{}
			m.Func = &ast.Function{}
			m.Func.Name = name
			m.Func.Typ = functionType
			m.Func.AccessFlags |= cg.ACC_METHOD_PUBLIC
			if c.classDefinition.Methods == nil {
				c.classDefinition.Methods = make(map[string][]*ast.ClassMethod)
			}
			c.classDefinition.Methods[m.Func.Name] = append(c.classDefinition.Methods[m.Func.Name], m)
		default:
			c.parser.errs = append(c.parser.errs, fmt.Errorf("%s unexpect token:%s", c.parser.errorMsgPrefix(), c.parser.token.Desp))
			c.Next()
		}
	}
	return
}

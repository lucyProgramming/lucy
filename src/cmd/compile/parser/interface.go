package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

type InterfaceParser struct {
	parser             *Parser
	classDefinition    *ast.Class
	isStatic           bool
	accessControlToken *lex.Token
}

func (interfaceParser *InterfaceParser) Next() {
	interfaceParser.parser.Next()
}

func (interfaceParser *InterfaceParser) consume(m map[int]bool) {
	interfaceParser.parser.consume(m)
}

func (interfaceParser *InterfaceParser) parse() (classDefinition *ast.Class, err error) {
	interfaceParser.Next() // skip interface key word
	classDefinition = &ast.Class{}
	interfaceParser.classDefinition = classDefinition
	interfaceParser.classDefinition.Pos = interfaceParser.parser.mkPos()
	interfaceParser.classDefinition.Block.IsClassBlock = true
	interfaceParser.classDefinition.AccessFlags |= cg.ACC_CLASS_INTERFACE // interface
	interfaceParser.classDefinition.AccessFlags |= cg.ACC_CLASS_ABSTRACT
	interfaceParser.classDefinition.Name, err = interfaceParser.parser.ClassParser.parseClassName()
	if err != nil {
		return nil, err
	}
	if interfaceParser.parser.token.Type == lex.TOKEN_EXTENDS { // parse father expression
		interfaceParser.Next() // skip extends
		interfaceParser.classDefinition.Pos = interfaceParser.parser.mkPos()
		if interfaceParser.parser.token.Type != lex.TOKEN_IDENTIFIER {
			err = fmt.Errorf("%s class`s father must be a identifier", interfaceParser.parser.errorMsgPrefix())
			interfaceParser.parser.errs = append(interfaceParser.parser.errs, err)
			interfaceParser.consume(untilLc) //
		} else {
			t, err := interfaceParser.parser.ClassParser.parseClassName()
			interfaceParser.classDefinition.SuperClassName = t
			if err != nil {
				interfaceParser.parser.errs = append(interfaceParser.parser.errs, err)
				return nil, err
			}
		}
	}
	if interfaceParser.parser.token.Type == lex.TOKEN_IMPLEMENTS {
		interfaceParser.Next() // skip key word
		interfaceParser.classDefinition.InterfaceNames, err = interfaceParser.parser.ClassParser.parseInterfaces()
		if err != nil {
			interfaceParser.consume(untilLc)
		}
	}
	if interfaceParser.parser.token.Type != lex.TOKEN_LC {
		err = fmt.Errorf("%s expect '{' but '%s'", interfaceParser.parser.errorMsgPrefix(), interfaceParser.parser.token.Description)
		interfaceParser.parser.errs = append(interfaceParser.parser.errs, err)
		interfaceParser.consume(untilLc)
	}
	interfaceParser.Next()
	for interfaceParser.parser.token.Type != lex.TOKEN_EOF {
		if len(interfaceParser.parser.errs) > interfaceParser.parser.nErrors2Stop {
			break
		}
		switch interfaceParser.parser.token.Type {
		case lex.TOKEN_RC:
			interfaceParser.Next()
			return
		case lex.TOKEN_SEMICOLON:
			interfaceParser.Next()
			continue
		case lex.TOKEN_FUNCTION:
			interfaceParser.Next() /// skip key word
			var name string
			if interfaceParser.parser.token.Type != lex.TOKEN_IDENTIFIER {
				interfaceParser.parser.errs = append(interfaceParser.parser.errs, fmt.Errorf("%s expect function name,but '%s'",
					interfaceParser.parser.errorMsgPrefix(), interfaceParser.parser.token.Description))
				interfaceParser.consume(untilRc)
				interfaceParser.Next()
				continue
			}
			name = interfaceParser.parser.token.Data.(string)
			interfaceParser.Next() // skip name
			functionType, err := interfaceParser.parser.parseFunctionType()
			if err != nil {
				interfaceParser.consume(untilRc)
				interfaceParser.Next()
				continue
			}
			if interfaceParser.classDefinition.Methods == nil {
				interfaceParser.classDefinition.Methods = make(map[string][]*ast.ClassMethod)
			}
			m := &ast.ClassMethod{}
			m.Function = &ast.Function{}
			m.Function.Name = name
			m.Function.Type = functionType
			m.Function.AccessFlags |= cg.ACC_METHOD_PUBLIC
			if interfaceParser.classDefinition.Methods == nil {
				interfaceParser.classDefinition.Methods = make(map[string][]*ast.ClassMethod)
			}
			interfaceParser.classDefinition.Methods[m.Function.Name] = append(interfaceParser.classDefinition.Methods[m.Function.Name], m)
		default:
			interfaceParser.parser.errs = append(interfaceParser.parser.errs, fmt.Errorf("%s unexpect token:%s", interfaceParser.parser.errorMsgPrefix(),
				interfaceParser.parser.token.Description))
			interfaceParser.Next()
		}
	}
	return
}

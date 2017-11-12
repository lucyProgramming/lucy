package parser

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

type ExpressionParser struct {
	parser *Parser
}

func (ep *ExpressionParser) Next() {
	ep.parser.Next()
}

func (ep *ExpressionParser) parseExpressions() ([]*ast.Expression, error) {
	es := []*ast.Expression{}
	for {
		e, err := ep.parseExpression(false)
		if err != nil {
			return nil, err
		}
		es = append(es, e)
		if ep.parser.token.Type != lex.TOKEN_COMMA {
			break
		}
		// == ,
		ep.Next()
	}
	return es, nil
}

//parse equal expression
func (ep *ExpressionParser) parseExpression(one bool) (*ast.Expression, error) {
	return ep.parseEqualExpression(one)
}

//begin with identifier
func (ep *ExpressionParser) parseIdentifierExpression() (*ast.Expression, error) {
	if ep.parser.token.Type != lex.TOKEN_IDENTIFIER {
		return nil, fmt.Errorf("it is not a identifier expression")
	}
	result := &ast.Expression{}
	result.Typ = ast.EXPRESSION_TYPE_IDENTIFIER
	result.Data = ep.parser.token.Data.(string)
	result.Pos = ep.parser.mkPos()
	ep.Next()          //look next token
	if ep.parser.eof { // end of file
		return result, nil
	}
	var err error
	for {
		if ep.parser.token.Type == lex.TOKEN_DOT { // a.b.c.e.f
			ep.Next()
			if ep.parser.eof {
				return nil, ep.parser.mkUnexpectedErr()
			}
			if ep.parser.token.Type != lex.TOKEN_IDENTIFIER {
				return nil, fmt.Errorf("%s %d%d excpet identifier afeter \".\",but %s",
					ep.parser.filename,
					ep.parser.token.Match.StartLine,
					ep.parser.token.Match.StartColumn,
					ep.parser.token.Desp)
			}
			newresult := &ast.Expression{}
			newresult.Pos = ep.parser.mkPos()
			newresult.Typ = ast.EXPRESSION_TYPE_DOT
			binary := &ast.ExpressionBinary{}
			newresult.Data = binary
			binary.Left = result
			binary.Right = &ast.Expression{
				Typ:  ast.EXPRESSION_TYPE_IDENTIFIER,
				Data: ep.parser.token.Data.(string),
			}
			result = newresult // reassignment
			ep.Next()
		} else if ep.parser.token.Type == lex.TOKEN_LB { // a["b"]
			ep.Next()
			if ep.parser.eof {
				return nil, ep.parser.mkUnexpectedErr()
			}
			newresult := &ast.Expression{}
			newresult.Pos = ep.parser.mkPos()
			newresult.Typ = ast.EXPRESSION_TYPE_INDEX
			binary := &ast.ExpressionBinary{}
			newresult.Data = binary
			binary.Left = result
			binary.Right, err = ep.parseExpression(false)
			if err != nil {
				return nil, err
			}
			ep.Next()
			if ep.parser.eof {
				return nil, ep.parser.mkUnexpectedErr()
			}
			if ep.parser.token.Type != lex.TOKEN_RB {
				err = fmt.Errorf("%s %d:%d [ and ] not match", ep.parser.filename, ep.parser.token.Match.StartLine, ep.parser.token.Match.StartColumn)
				break
			}
			ep.Next()
			result = newresult
		} else if ep.parser.token.Type == lex.TOKEN_LP { // a() or a.say() a["call"]()
			ep.Next()
			if ep.parser.eof {
				return nil, ep.parser.mkUnexpectedErr()
			}
			args := []*ast.Expression{}
			if ep.parser.token.Type != lex.TOKEN_RP { //a(123)
				args, err = ep.parseExpressions()
				if err != nil {
					break
				}
			} else { //ep.parser.token.Type == lex.TOKEN_RP
			}
			if result.Typ == ast.EXPRESSION_TYPE_IDENTIFIER || result.Typ == ast.EXPRESSION_TYPE_INDEX {
				newresult := &ast.Expression{
					Typ: ast.EXPRESSION_TYPE_FUNCTION_CALL,
				}
				call := &ast.ExpressionFunctionCall{}
				call.Pos = ep.parser.mkPos()
				call.Expression = result
				call.Args = args
				newresult.Data = call
				result = newresult
			} else if result.Typ == ast.EXPRESSION_TYPE_DOT {
				newresult := &ast.Expression{
					Typ: ast.EXPRESSION_TYPE_METHOD_CALL,
				}
				newresult.Pos = ep.parser.mkPos()
				call := &ast.ExpressionMethodCall{}
				binary := result.Data.(*ast.ExpressionBinary)
				call.Expression = binary.Left
				call.Name = binary.Right.Data.(string)
				call.Args = args
				result = newresult
			} else {
				err = fmt.Errorf("%s %d%d can`t make call on that situation", ep.parser.filename, ep.parser.token.Match.StartLine, ep.parser.token.Match.StartColumn)
				break
			}
			if ep.parser.token.Type != lex.TOKEN_RP {
				err = fmt.Errorf("%s %d%d except \")\" ,but %s",
					ep.parser.filename,
					ep.parser.token.Match.StartLine,
					ep.parser.token.Match.StartColumn,
					ep.parser.token.Desp)
				break
			}
			ep.Next() //loop next token
		} else {
			// something i can`t handle
			break
		}
	}
	return result, err
}

// a = 123
func (ep *ExpressionParser) parseEqualExpression(one bool) (*ast.Expression, error) {
	var left *ast.Expression
	var err error
	switch ep.parser.token.Type {
	case lex.TOKEN_IDENTIFIER:
		left, err = ep.parseIdentifierExpression()
		if err != nil {
			return nil, err
		}
	case lex.TOKEN_LITERAL_BOOL:
		left = &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_BOOL,
			Data: ep.parser.token.Data,
		}
	case lex.TOKEN_LITERAL_BYTE:
		left = &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_BYTE,
			Data: ep.parser.token.Data,
		}
	case lex.TOKEN_LITERAL_INT:
		left = &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_INT,
			Data: ep.parser.token.Data,
		}
	case lex.TOKEN_LITERAL_STRING:
		left = &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_STRING,
			Data: ep.parser.token.Data,
		}
	case lex.TOKEN_LITERAL_FLOAT:
		left = &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_STRING,
			Data: ep.parser.token.Data,
		}
	case lex.TOKEN_LP:
		left, err = ep.parseExpression(false)
		if err != nil {
			return nil, err
		}
		if ep.parser.token.Type != lex.TOKEN_RP {
			return nil, fmt.Errorf("%s ( and ) not matched", ep.parser.errorMsgPrefix())
		}
	case lex.TOKEN_INCREMENT:
		newE := &ast.Expression{}
		newE.Pos = ep.parser.mkPos()
		left, err = ep.parseExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Typ = ast.EXPRESSION_TYPE_PRE_INCREMENT
		newE.Data = left
		left = newE
	case lex.TOKEN_DECREMENT:
		newE := &ast.Expression{}
		left, err = ep.parseExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Typ = ast.EXPRESSION_TYPE_PRE_DECREMENT
		newE.Data = left
		left = newE
	case lex.TOKEN_NOT:
		newE := &ast.Expression{}
		left, err = ep.parseExpression(true)
		if err != nil {
			return nil, err
		}
		newE.Typ = ast.EXPRESSION_TYPE_NOT
		newE.Data = left
		left = newE
	default:
		return nil, fmt.Errorf("%s unkown begining of a expression", ep.parser.errorMsgPrefix())
	}
	ep.Next() // look next
	if one {
		return left, nil
	}
	if ep.parser.eof {
		return left, nil
	}
	mkBinayExpression := func(typ int) (*ast.Expression, error) {
		ep.Next()
		if ep.parser.eof {
			return nil, ep.parser.mkUnexpectedErr()
		}
		result := &ast.Expression{}
		result.Typ = typ
		binary := &ast.ExpressionBinary{}
		result.Data = binary
		binary.Left = left
		binary.Right, err = ep.parseExpression(false)
		return result, err
	}

	// := += -= *= /= %=
	switch ep.parser.token.Type {
	case lex.TOKEN_COLON_ASSIGN:
		return mkBinayExpression(ast.EXPRESSION_TYPE_COLON_ASSIGN)
	case lex.TOKEN_PLUS_ASSIGN:
		return mkBinayExpression(ast.EXPRESSION_TYPE_PLUS_ASSIGN)
	case lex.TOKEN_MINUS_ASSIGN:
		return mkBinayExpression(ast.EXPRESSION_TYPE_MINUS_ASSIGN)
	case lex.TOKEN_MUL_ASSIGN:
		return mkBinayExpression(ast.EXPRESSION_TYPE_MUL_ASSIGN)
	case lex.TOKEN_DIV_ASSIGN:
		return mkBinayExpression(ast.EXPRESSION_TYPE_DIV_ASSIGN)
	case lex.TOKEN_MOD_ASSIGN:
		return mkBinayExpression(ast.EXPRESSION_TYPE_MOD_ASSIGN)
	}
	// && ||
	switch ep.parser.token.Type {
	case lex.TOKEN_LOGICAL_AND:
		return mkBinayExpression(ast.EXPRESSION_TYPE_LOGICAL_AND)
	case lex.TOKEN_LOGICAL_OR:
		return mkBinayExpression(ast.EXPRESSION_TYPE_LOGICAL_OR)
	}

	// & |

	switch ep.parser.token.Type {
	case lex.TOKEN_AND:
		return mkBinayExpression(ast.EXPRESSION_TYPE_AND)
	case lex.TOKEN_OR:
		return mkBinayExpression(ast.EXPRESSION_TYPE_OR)
	}

	// == !=

	switch ep.parser.token.Type {
	case lex.TOKEN_EQUAL:
		return mkBinayExpression(ast.EXPRESSION_TYPE_EQ)
	case lex.TOKEN_NE:
		return mkBinayExpression(ast.EXPRESSION_TYPE_NE)
	}
	// > < >= <=
	switch ep.parser.token.Type {
	case lex.TOKEN_GE:
		return mkBinayExpression(ast.EXPRESSION_TYPE_GE)
	case lex.TOKEN_GT:
		return mkBinayExpression(ast.EXPRESSION_TYPE_GT)
	case lex.TOKEN_LT:
		return mkBinayExpression(ast.EXPRESSION_TYPE_LT)
	case lex.TOKEN_LE:
		return mkBinayExpression(ast.EXPRESSION_TYPE_LE)
	}

	// << >>
	switch ep.parser.token.Type {
	case lex.TOKEN_LEFT_SHIFT:
		return mkBinayExpression(ast.EXPRESSION_TYPE_LEFT_SHIFT)
	case lex.TOKEN_RIGHT_SHIFT:
		return mkBinayExpression(ast.EXPRESSION_TYPE_RIGHT_SHIFT)
	}

	// + - * / %
	switch ep.parser.token.Type {
	case lex.TOKEN_ADD:
		return mkBinayExpression(ast.EXPRESSION_TYPE_ADD)
	case lex.TOKEN_SUB:
		return mkBinayExpression(ast.EXPRESSION_TYPE_SUB)
	case lex.TOKEN_MUL:
		return mkBinayExpression(ast.EXPRESSION_TYPE_MUL)
	case lex.TOKEN_DIV:
		return mkBinayExpression(ast.EXPRESSION_TYPE_DIV)
	case lex.TOKEN_MOD:
		return mkBinayExpression(ast.EXPRESSION_TYPE_MOD)
	}
	// ++ --
	switch ep.parser.token.Type {
	case lex.TOKEN_INCREMENT:
		return &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_INCREMENT,
			Data: left,
		}, nil
	case lex.TOKEN_DECREMENT:
		return &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_DECREMENT,
			Data: left,
		}, nil
	}
	return left, nil // no further token can be matched
}

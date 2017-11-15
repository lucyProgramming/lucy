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
	for !ep.parser.eof {
		e, err := ep.parseExpression()
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

//parse assign expression
func (ep *ExpressionParser) parseExpression() (*ast.Expression, error) {
	return ep.parseAssignExpression()
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
				return nil, ep.parser.mkUnexpectedEofErr()
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
				return nil, ep.parser.mkUnexpectedEofErr()
			}
			newresult := &ast.Expression{}
			newresult.Pos = ep.parser.mkPos()
			newresult.Typ = ast.EXPRESSION_TYPE_INDEX
			binary := &ast.ExpressionBinary{}
			newresult.Data = binary
			binary.Left = result
			binary.Right, err = ep.parseExpression()
			if err != nil {
				return nil, err
			}
			ep.Next()
			if ep.parser.eof {
				return nil, ep.parser.mkUnexpectedEofErr()
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
				return nil, ep.parser.mkUnexpectedEofErr()
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

func (ep *ExpressionParser) parseIdentifierExpressions() (ret []*ast.Expression, err error) {
	ret = []*ast.Expression{}
	var e *ast.Expression
	for {
		e, err = ep.parseIdentifierExpression()
		if err != nil {
			return
		}
		ret = append(ret, e)
		ep.Next()
		if ep.parser.token.Type != lex.TOKEN_COMMA {
			return
		}
	}
	return
}

func (ep *ExpressionParser) parseOneExpression() (*ast.Expression, error) {
	if ep.parser.eof {
		return nil, ep.parser.mkUnexpectedEofErr()
	}
	var left *ast.Expression
	var err error
	switch ep.parser.token.Type {
	case lex.TOKEN_IDENTIFIER:
		lefts, err := ep.parseIdentifierExpressions()
		if err != nil {
			return nil, err
		}
		if len(lefts) == 1 {
			left = lefts[0]
		} else {
			left = &ast.Expression{
				Typ:  ast.EXPRESSION_TYPE_NAME_LIST,
				Data: lefts,
			}
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
		ep.parser.Next()
		if ep.parser.eof {
			return nil, ep.parser.mkUnexpectedEofErr()
		}
		left, err = ep.parseExpression()
		if err != nil {
			return nil, err
		}
		if ep.parser.token.Type != lex.TOKEN_RP {
			return nil, fmt.Errorf("%s ( and ) not matched, but %s", ep.parser.errorMsgPrefix(), ep.parser.token.Desp)
		}
	case lex.TOKEN_INCREMENT:
		ep.parser.Next()
		newE := &ast.Expression{}
		newE.Pos = ep.parser.mkPos()
		left, err = ep.parseExpression()
		if err != nil {
			return nil, err
		}
		newE.Typ = ast.EXPRESSION_TYPE_PRE_INCREMENT
		newE.Data = left
		left = newE
	case lex.TOKEN_DECREMENT:
		ep.parser.Next()
		newE := &ast.Expression{}
		left, err = ep.parseExpression()
		if err != nil {
			return nil, err
		}
		newE.Typ = ast.EXPRESSION_TYPE_PRE_DECREMENT
		newE.Data = left
		left = newE
	case lex.TOKEN_NOT:
		ep.parser.Next()
		newE := &ast.Expression{}
		left, err = ep.parseExpression()
		if err != nil {
			return nil, err
		}
		newE.Typ = ast.EXPRESSION_TYPE_NOT
		newE.Data = left
		left = newE
	case lex.TOKEN_SUB:
		ep.parser.Next()
		newE := &ast.Expression{}
		left, err = ep.parseExpression()
		if err != nil {
			return nil, err
		}
		newE.Typ = ast.EXPRESSION_TYPE_NEGATIVE
		newE.Data = left
		left = newE
	default:
		return nil, fmt.Errorf("%s unkown begining of a expression or forget to write a expression, token:%s", ep.parser.errorMsgPrefix(), ep.parser.token.Desp)
	}
	ep.Next() // look next
	// a++ b-- for a++++
	for ep.parser.token.Type == lex.TOKEN_INCREMENT || ep.parser.token.Type == lex.TOKEN_DECREMENT {
		newe := &ast.Expression{}
		if ep.parser.token.Type == lex.TOKEN_INCREMENT {
			newe.Typ = ast.EXPRESSION_TYPE_INCREMENT
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_DECREMENT
		}
		newe.Data = left
		left = newe
		ep.Next()
	}
	return left, nil
}

// && ||
func (ep *ExpressionParser) parseLogicalExpression() (*ast.Expression, error) {
	e, err := ep.parseAndExpression()
	if err != nil {
		return nil, err
	}
	for (ep.parser.token.Type == lex.TOKEN_LOGICAL_AND || ep.parser.token.Type == lex.TOKEN_LOGICAL_OR) && !ep.parser.eof {
		ep.Next()
		e2, err := ep.parseAndExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		if ep.parser.token.Type == lex.TOKEN_LOGICAL_AND {
			newe.Typ = ast.EXPRESSION_TYPE_LOGICAL_AND
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_LOGICAL_OR
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// & |
func (ep *ExpressionParser) parseAndExpression() (*ast.Expression, error) {
	e, err := ep.parseEqualExpression()
	if err != nil {
		return nil, err
	}
	for (ep.parser.token.Type == lex.TOKEN_AND || ep.parser.token.Type == lex.TOKEN_OR) && !ep.parser.eof {
		ep.Next()
		e2, err := ep.parseEqualExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		if ep.parser.token.Type == lex.TOKEN_AND {
			newe.Typ = ast.EXPRESSION_TYPE_AND
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_OR
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// == and !=
func (ep *ExpressionParser) parseEqualExpression() (*ast.Expression, error) {
	e, err := ep.parseRelationExpression()
	if err != nil {
		return nil, err
	}
	for (ep.parser.token.Type == lex.TOKEN_EQUAL || ep.parser.token.Type == lex.TOKEN_NE) && !ep.parser.eof {
		ep.Next()
		e2, err := ep.parseRelationExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		if ep.parser.token.Type == lex.TOKEN_EQUAL {
			newe.Typ = ast.EXPRESSION_TYPE_EQ
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_NE
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// > < >= <=
func (ep *ExpressionParser) parseRelationExpression() (*ast.Expression, error) {
	e, err := ep.parseShiftExpression()
	if err != nil {
		return nil, err
	}
	for (ep.parser.token.Type == lex.TOKEN_GT || ep.parser.token.Type == lex.TOKEN_GE ||
		ep.parser.token.Type == lex.TOKEN_LT || ep.parser.token.Type == lex.TOKEN_LE) && !ep.parser.eof {
		ep.Next()
		e2, err := ep.parseShiftExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		if ep.parser.token.Type == lex.TOKEN_GT {
			newe.Typ = ast.EXPRESSION_TYPE_GT
		} else if ep.parser.token.Type == lex.TOKEN_GE {
			newe.Typ = ast.EXPRESSION_TYPE_GE
		} else if ep.parser.token.Type == lex.TOKEN_LT {
			newe.Typ = ast.EXPRESSION_TYPE_LT
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_LE
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe

	}
	return e, nil
}

// << >>
func (ep *ExpressionParser) parseShiftExpression() (*ast.Expression, error) {
	e, err := ep.parseAddExpression()
	if err != nil {
		return nil, err
	}
	for (ep.parser.token.Type == lex.TOKEN_LEFT_SHIFT || ep.parser.token.Type == lex.TOKEN_RIGHT_SHIFT) && !ep.parser.eof {
		ep.Next()
		e2, err := ep.parseAddExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		if ep.parser.token.Type == lex.TOKEN_LEFT_SHIFT {
			newe.Typ = ast.EXPRESSION_TYPE_LEFT_SHIFT
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_RIGHT_SHIFT
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe

	}
	return e, nil
}

// + -
func (ep *ExpressionParser) parseAddExpression() (*ast.Expression, error) {
	e, err := ep.parseMulExpression()
	if err != nil {
		return nil, err
	}
	for (ep.parser.token.Type == lex.TOKEN_ADD || ep.parser.token.Type == lex.TOKEN_SUB) && !ep.parser.eof {
		ep.Next()
		e2, err := ep.parseMulExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		if ep.parser.token.Type == lex.TOKEN_ADD {
			newe.Typ = ast.EXPRESSION_TYPE_ADD
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_SUB
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

// */ %
func (ep *ExpressionParser) parseMulExpression() (*ast.Expression, error) {
	e, err := ep.parseOneExpression()
	if err != nil {
		return nil, err
	}
	for (ep.parser.token.Type == lex.TOKEN_MUL || ep.parser.token.Type == lex.TOKEN_DIV || ep.parser.token.Type == lex.TOKEN_MOD) && !ep.parser.eof {
		ep.Next()
		e2, err := ep.parseOneExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		if ep.parser.token.Type == lex.TOKEN_MUL {
			newe.Typ = ast.EXPRESSION_TYPE_MUL
		} else if ep.parser.token.Type == lex.TOKEN_DIV {
			newe.Typ = ast.EXPRESSION_TYPE_DIV
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_MOD
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
		ep.Next()
	}
	return e, nil
}

//a = 123
func (ep *ExpressionParser) parseAssignExpression() (*ast.Expression, error) {
	left, err := ep.parseLogicalExpression()
	if err != nil {
		return nil, err
	}
	mkBinayExpression := func(typ int) (*ast.Expression, error) {
		ep.Next()
		if ep.parser.eof {
			return nil, ep.parser.mkUnexpectedEofErr()
		}
		result := &ast.Expression{}
		result.Typ = typ
		binary := &ast.ExpressionBinary{}
		result.Data = binary
		binary.Left = left
		binary.Right, err = ep.parseLogicalExpression()
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
	return left, nil
}
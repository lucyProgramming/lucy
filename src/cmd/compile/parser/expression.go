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
			return es, err
		}
		if e.Typ == ast.EXPRESSION_TYPE_LIST {
			es = append(es, e.Data.([]*ast.Expression)...)
		} else {
			es = append(es, e)
		}
		if ep.parser.token.Type != lex.TOKEN_COMMA {
			break
		}
		// == ,
		ep.Next() // skip ,
	}
	return es, nil
}

//parse assign expression
func (ep *ExpressionParser) parseExpression() (*ast.Expression, error) {
	return ep.parseAssignExpression()
}

func (ep *ExpressionParser) parseCallExpression(e *ast.Expression) (*ast.Expression, error) {
	var err error
	ep.Next() // skip (
	if ep.parser.eof {
		return nil, ep.parser.mkUnexpectedEofErr()
	}
	args := []*ast.Expression{}
	if ep.parser.token.Type != lex.TOKEN_RP { //a(123)
		args, err = ep.parseExpressions()
		if err != nil {
			return nil, err
		}
	}

	if ep.parser.token.Type != lex.TOKEN_RP {
		return nil, fmt.Errorf("%s except ')' ,but %s",
			ep.parser.errorMsgPrefix(),
			ep.parser.token.Desp)
	}
	var result ast.Expression
	if e.Typ == ast.EXPRESSION_TYPE_IDENTIFIER {
		result.Typ = ast.EXPRESSION_TYPE_FUNCTION_CALL
		call := &ast.ExpressionFunctionCall{}
		call.Expression = e
		call.Args = args
		result.Data = call
		result.Pos = ep.parser.mkPos()
	} else if result.Typ == ast.EXPRESSION_TYPE_DOT {
		result.Typ = ast.EXPRESSION_TYPE_METHOD_CALL
		call := &ast.ExpressionMethodCall{}
		index := e.Data.(*ast.ExpressionIndex)
		call.Expression = index.Expression
		call.Args = args
		result.Data = call
		result.Pos = e.Pos
	} else {
		return nil, fmt.Errorf("%s can`t make call on '%s'", ep.parser.errorMsgPrefix())
	}
	ep.Next() // skip )
	return &result, nil
}

// []int{1,2,3}
func (ep *ExpressionParser) parseArrayExpression() (*ast.Expression, error) {
	fmt.Println(ep.parser.token.Desp)
	t, err := ep.parser.parseType()
	if err != nil {
		return nil, err
	}
	arr := &ast.ExpressionArray{}
	arr.Typ = t
	arr.Expression, err = ep.parseValueGroups()
	return &ast.Expression{
		Typ:  ast.EXPRESSION_TYPE_ARRAY,
		Data: arr,
	}, err
}

//{1,2,3}  {{1,2,3},{456}}
func (ep *ExpressionParser) parseValueGroups() (*ast.Expression, error) {
	if ep.parser.token.Type != lex.TOKEN_LC {
		return nil, fmt.Errorf("%s no { after type definition", ep.parser.errorMsgPrefix())
	}
	ret := &ast.Expression{}
	ret.Typ = ast.EXPRESSION_TYPE_LIST
	es := []*ast.Expression{}
	var e *ast.Expression
	var err error
	defer func() {
		ret.Data = es
	}()
	ep.Next() // skip {
	if ep.parser.token.Type == lex.TOKEN_RC {
		return ret, nil
	}
	if ep.parser.token.Type == lex.TOKEN_LC {
		for !ep.parser.eof && ep.parser.token.Type == lex.TOKEN_LC {
			e, err = ep.parseValueGroups()
			if err != nil {
				return ret, err
			}
			es = append(es, e)
			if ep.parser.token.Type != lex.TOKEN_COMMA {
				break
			}
			ep.Next()
		}
	} else {
		es, err = ep.parseExpressions()
		if err != nil {
			return ret, nil
		}
	}
	if ep.parser.token.Type != lex.TOKEN_RC {
		return ret, fmt.Errorf("%s missing } ", ep.parser.errorMsgPrefix())
	}
	ep.Next()
	return ret, nil

}

func (ep *ExpressionParser) parseTypeConvertionExpression() (*ast.Expression, error) {
	t, err := ep.parser.parseType()
	if err != nil {
		return nil, err
	}
	if ep.parser.token.Type != lex.TOKEN_LP {
		return nil, fmt.Errorf("%s not ( after a type", ep.parser.errorMsgPrefix())
	}
	ep.Next() // skip (
	e, err := ep.parseExpression()
	if err != nil {
		return nil, err
	}
	if ep.parser.token.Type != lex.TOKEN_RP {
		return nil, fmt.Errorf("%s ( and ) not match", ep.parser.errorMsgPrefix())
	}
	ep.Next() // skip ) for next process
	return &ast.Expression{
		Typ: ast.EXPRESSION_TYPE_CONVERTION_TYPE,
		Data: &ast.ExpressionTypeConvertion{
			Typ:        t,
			Expression: e,
		},
	}, nil
}

func (ep *ExpressionParser) parseOneExpression() (*ast.Expression, error) {
	if ep.parser.eof {
		return nil, ep.parser.mkUnexpectedEofErr()
	}
	var left *ast.Expression
	var err error
	switch ep.parser.token.Type {
	case lex.TOKEN_IDENTIFIER:
		left = &ast.Expression{}
		left.Typ = ast.EXPRESSION_TYPE_IDENTIFIER
		identifier := &ast.ExpressionIdentifer{}
		identifier.Name = ep.parser.token.Data.(string)
		left.Data = identifier
		ep.Next()
	case lex.TOKEN_TRUE:
		left = &ast.Expression{}
		left.Typ = ast.EXPRESSION_TYPE_BOOL
		left.Data = true
		left.Pos = ep.parser.mkPos()
		ep.Next()
	case lex.TOKEN_FALSE:
		left = &ast.Expression{}
		left.Typ = ast.EXPRESSION_TYPE_BOOL
		left.Data = false
		left.Pos = ep.parser.mkPos()
		ep.Next()
	case lex.TOKEN_LITERAL_BYTE:
		left = &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_BYTE,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_LITERAL_INT:
		left = &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_INT,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_LITERAL_STRING:
		left = &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_STRING,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_LITERAL_FLOAT:
		left = &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_STRING,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_LP:
		ep.Next()
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
		ep.Next()
	case lex.TOKEN_INCREMENT:
		ep.Next() // skip ++
		newE := &ast.Expression{}
		newE.Pos = ep.parser.mkPos()
		left, err = ep.parseOneExpression()
		if err != nil {
			return nil, err
		}
		newE.Typ = ast.EXPRESSION_TYPE_PRE_INCREMENT
		newE.Data = left
		left = newE
	case lex.TOKEN_DECREMENT:
		ep.Next() // skip --
		newE := &ast.Expression{}
		left, err = ep.parseOneExpression()
		if err != nil {
			return nil, err
		}
		newE.Typ = ast.EXPRESSION_TYPE_PRE_DECREMENT
		newE.Data = left
		newE.Pos = ep.parser.mkPos()
		left = newE
	case lex.TOKEN_NOT:
		ep.Next()
		newE := &ast.Expression{}
		left, err = ep.parseOneExpression()
		if err != nil {
			return nil, err
		}
		newE.Typ = ast.EXPRESSION_TYPE_NOT
		newE.Data = left
		newE.Pos = ep.parser.mkPos()
		left = newE
	case lex.TOKEN_SUB:
		ep.Next()
		newE := &ast.Expression{}
		left, err = ep.parseOneExpression()
		if err != nil {
			return nil, err
		}
		newE.Typ = ast.EXPRESSION_TYPE_NEGATIVE
		newE.Data = left
		newE.Pos = ep.parser.mkPos()
		left = newE
	case lex.TOKEN_FUNCTION:
		f, err := ep.parser.Function.parse(false)
		if err != nil {
			return nil, err
		}
		return &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_FUNCTION,
			Data: f,
		}, nil
	case lex.TOKEN_NEW:
		ep.Next()
		t, err := ep.parser.parseIdentifierType()
		if err != nil {
			return nil, err
		}
		if ep.parser.token.Type != lex.TOKEN_LP {
			return nil, fmt.Errorf("%s missing ( after new", ep.parser.errorMsgPrefix())
		}
		ep.Next()
		var es []*ast.Expression
		if ep.parser.token.Type != lex.TOKEN_RP { //
			es, err = ep.parseExpressions()
			if err != nil {
				return nil, err
			}
		}
		if ep.parser.token.Type != lex.TOKEN_RP {
			return nil, fmt.Errorf("%s ( and ) not match", ep.parser.errorMsgPrefix())
		}
		ep.Next()
		left = &ast.Expression{
			Pos: ep.parser.mkPos(),
			Typ: ast.EXPRESSION_TYPE_NEW,
			Data: &ast.ExpressionNew{
				Args: es,
				Typ:  t,
			},
		}
	case lex.TOKEN_LB:
		left, err = ep.parseArrayExpression()
		if err != nil {
			return left, err
		}
	// bool(xxx)
	case lex.TOKEN_BOOL:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
		//
	case lex.TOKEN_BYTE:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_SHORT:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_INT:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_FLOAT:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_DOUBLE:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	case lex.TOKEN_STRING:
		left, err = ep.parseTypeConvertionExpression()
		if err != nil {
			return left, err
		}
	default:
		err = fmt.Errorf("%s unkown begining of a expression, token:%s", ep.parser.errorMsgPrefix(), ep.parser.token.Desp)
		return nil, err
	}
	for ep.parser.token.Type == lex.TOKEN_INCREMENT || ep.parser.token.Type == lex.TOKEN_DECREMENT ||
		ep.parser.token.Type == lex.TOKEN_LP || ep.parser.token.Type == lex.TOKEN_LB ||
		ep.parser.token.Type == lex.TOKEN_DOT {
		if ep.parser.token.Type == lex.TOKEN_INCREMENT || ep.parser.token.Type == lex.TOKEN_DECREMENT { //  ++ or --
			newe := &ast.Expression{}
			if ep.parser.token.Type == lex.TOKEN_INCREMENT {
				newe.Typ = ast.EXPRESSION_TYPE_INCREMENT
			} else {
				newe.Typ = ast.EXPRESSION_TYPE_DECREMENT
			}
			if left.Typ != ast.EXPRESSION_TYPE_LIST {
				newe.Data = left
				left = newe
			} else {
				list := left.Data.([]*ast.Expression)
				newe.Data = list[len(list)-1]
			}
			ep.Next()
			continue
		}
		if ep.parser.token.Type == lex.TOKEN_LB {
			ep.Next() // skip [
			e, err := ep.parseExpression()
			if err != nil {
				return nil, err
			}
			if ep.parser.token.Type != lex.TOKEN_RB {
				return nil, fmt.Errorf("%s not ] after a expression", ep.parser.errorMsgPrefix())
			}
			newe := &ast.Expression{}
			newe.Typ = ast.EXPRESSION_TYPE_INDEX
			index := &ast.ExpressionIndex{}
			index.Expression = left
			index.Index = e
			index.Name = ep.parser.token.Data.(string)
			newe.Data = index
			left = newe
			ep.Next()
			continue
		}
		if ep.parser.token.Type == lex.TOKEN_DOT {
			ep.parser.Next() // skip .
			if ep.parser.token.Type != lex.TOKEN_IDENTIFIER {
				return nil, fmt.Errorf("%s not identifier after dot", ep.parser.errorMsgPrefix())
			}
			newe := &ast.Expression{}
			newe.Typ = ast.EXPRESSION_TYPE_DOT
			index := &ast.ExpressionIndex{}
			index.Expression = left
			index.Name = ep.parser.token.Data.(string)
			newe.Data = index
			left = newe
			ep.Next()
			continue
		}
		if ep.parser.token.Type == lex.TOKEN_LP {
			newe, err := ep.parseCallExpression(left)
			if err != nil {
				return nil, err
			}
			left = newe
			continue
		}
	}
	return left, nil
}

// && ||
func (ep *ExpressionParser) parseLogicalExpression() (*ast.Expression, error) {
	e, err := ep.parseBitANDORExpression()
	if err != nil {
		return nil, err
	}
	for (ep.parser.token.Type == lex.TOKEN_LOGICAL_AND || ep.parser.token.Type == lex.TOKEN_LOGICAL_OR) && !ep.parser.eof {
		typ := ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseBitANDORExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_LOGICAL_AND {
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
func (ep *ExpressionParser) parseBitANDORExpression() (*ast.Expression, error) {
	e, err := ep.parseEqualExpression()
	if err != nil {
		return nil, err
	}
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_AND || ep.parser.token.Type == lex.TOKEN_OR) && !ep.parser.eof {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseEqualExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_AND {
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
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_EQUAL || ep.parser.token.Type == lex.TOKEN_NE) && !ep.parser.eof {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseRelationExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_EQUAL {
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
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_GT || ep.parser.token.Type == lex.TOKEN_GE ||
		ep.parser.token.Type == lex.TOKEN_LT || ep.parser.token.Type == lex.TOKEN_LE) && !ep.parser.eof {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseShiftExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_GT {
			newe.Typ = ast.EXPRESSION_TYPE_GT
		} else if typ == lex.TOKEN_GE {
			newe.Typ = ast.EXPRESSION_TYPE_GE
		} else if typ == lex.TOKEN_LT {
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
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_LEFT_SHIFT || ep.parser.token.Type == lex.TOKEN_RIGHT_SHIFT) && !ep.parser.eof {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseAddExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_LEFT_SHIFT {
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
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_ADD || ep.parser.token.Type == lex.TOKEN_SUB) && !ep.parser.eof {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseMulExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_ADD {
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
	var typ int
	for (ep.parser.token.Type == lex.TOKEN_MUL || ep.parser.token.Type == lex.TOKEN_DIV || ep.parser.token.Type == lex.TOKEN_MOD) && !ep.parser.eof {
		typ = ep.parser.token.Type
		pos := ep.parser.mkPos()
		ep.Next()
		e2, err := ep.parseOneExpression()
		if err != nil {
			return nil, err
		}
		newe := &ast.Expression{}
		newe.Pos = pos
		if typ == lex.TOKEN_MUL {
			newe.Typ = ast.EXPRESSION_TYPE_MUL
		} else if typ == lex.TOKEN_DIV {
			newe.Typ = ast.EXPRESSION_TYPE_DIV
		} else {
			newe.Typ = ast.EXPRESSION_TYPE_MOD
		}
		binary := &ast.ExpressionBinary{}
		binary.Left = e
		binary.Right = e2
		newe.Data = binary
		e = newe
	}
	return e, nil
}

//a = 123
func (ep *ExpressionParser) parseAssignExpression() (*ast.Expression, error) {
	left, err := ep.parseLogicalExpression() //
	if err != nil {
		return nil, err
	}
	for ep.parser.token.Type == lex.TOKEN_COMMA { // read more
		ep.Next()                                 //  skip comma
		left2, err := ep.parseLogicalExpression() //
		if err != nil {
			return nil, err
		}
		if left.Typ == ast.EXPRESSION_TYPE_LIST {
			list := left.Data.([]*ast.Expression)
			left.Data = append(list, left2)
		} else {
			newe := &ast.Expression{}
			newe.Typ = ast.EXPRESSION_TYPE_LIST
			list := []*ast.Expression{left, left2}
			newe.Data = list
			left = newe
		}
	}
	mustBeOneExpression := func() {
		if left.Typ == ast.EXPRESSION_TYPE_LIST {
			es := left.Data.([]*ast.Expression)
			left = es[0]
			if len(es) > 1 {
				ep.parser.errs = append(ep.parser.errs, fmt.Errorf("%s expect one left value", ep.parser.errorMsgPrefix(es[1].Pos)))
			}
		}
	}
	mkBinayExpression := func(typ int, multi bool) (*ast.Expression, error) {
		ep.Next() // skip = :=
		if ep.parser.eof {
			return nil, ep.parser.mkUnexpectedEofErr()
		}
		result := &ast.Expression{}
		result.Typ = typ
		binary := &ast.ExpressionBinary{}
		result.Data = binary
		binary.Left = left
		result.Pos = ep.parser.mkPos()
		if multi {
			es, err := ep.parseExpressions()
			if err != nil {
				return result, err
			}
			binary.Right = &ast.Expression{}
			binary.Right.Typ = ast.EXPRESSION_TYPE_LIST
			binary.Right.Data = es
		} else {
			binary.Right, err = ep.parseExpression()
		}
		return result, err
	}
	// := += -= *= /= %=
	switch ep.parser.token.Type {
	case lex.TOKEN_ASSIGN:
		return mkBinayExpression(ast.EXPRESSION_TYPE_ASSIGN, true)
	case lex.TOKEN_COLON_ASSIGN:
		return mkBinayExpression(ast.EXPRESSION_TYPE_COLON_ASSIGN, true)
	case lex.TOKEN_ADD_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_PLUS_ASSIGN, false)
	case lex.TOKEN_SUB_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_MINUS_ASSIGN, false)
	case lex.TOKEN_MUL_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_MUL_ASSIGN, false)
	case lex.TOKEN_DIV_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_DIV_ASSIGN, false)
	case lex.TOKEN_MOD_ASSIGN:
		mustBeOneExpression()
		return mkBinayExpression(ast.EXPRESSION_TYPE_MOD_ASSIGN, false)
	}
	return left, nil
}

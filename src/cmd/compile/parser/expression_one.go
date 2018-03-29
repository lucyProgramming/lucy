package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

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
		left.Pos = ep.parser.mkPos()
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
	case lex.TOKEN_LITERAL_LONG:
		left = &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_LONG,
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
			Typ:  ast.EXPRESSION_TYPE_FLOAT,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_LITERAL_DOUBLE:
		left = &ast.Expression{
			Typ:  ast.EXPRESSION_TYPE_DOUBLE,
			Data: ep.parser.token.Data,
			Pos:  ep.parser.mkPos(),
		}
		ep.Next()
	case lex.TOKEN_NULL:
		left = &ast.Expression{
			Typ: ast.EXPRESSION_TYPE_NULL,
			Pos: ep.parser.mkPos(),
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
	// case lex.TOKEN_FUNCTION:
	// 	f, err := ep.parser.Function.parse(false)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return &ast.Expression{
	// 		Typ:  ast.EXPRESSION_TYPE_FUNCTION,
	// 		Data: f,
	// 	}, nil
	case lex.TOKEN_NEW:
		ep.Next()
		t, err := ep.parser.parseType()
		if err != nil {
			return nil, err
		}
		if ep.parser.token.Type != lex.TOKEN_LP {
			return nil, fmt.Errorf("%s missing '(' after new", ep.parser.errorMsgPrefix())
		}
		ep.Next() // skip (
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
	case lex.TOKEN_RANGE:
		pos := ep.parser.mkPos()
		ep.Next()
		e, err := ep.parseExpression()
		if err != nil {
			return nil, err
		}
		left = &ast.Expression{}
		left.Typ = ast.EXPRESSION_TYPE_RANGE
		left.Pos = pos
		left.Data = e
		return left, nil
	case lex.TOKEN_MAP:
		left, err = ep.parseMapExprssion(true)
		if err != nil {
			return left, err
		}
	case lex.TOKEN_LC:
		left, err = ep.parseMapExprssion(false)
		if err != nil {
			return left, err
		}
	default:
		err = fmt.Errorf("%s unkown begining of a expression, token:%s", ep.parser.errorMsgPrefix(), ep.parser.token.Desp)
		return nil, err
	}
	if ep.parser.token.Type == lex.TOKEN_COLON && left.Typ == ast.EXPRESSION_TYPE_IDENTIFIER {
		left.Typ = ast.EXPRESSION_TYPE_LABLE
		ep.Next()
		return left, nil // lable here
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
			newe.Pos = ep.parser.mkPos()
			ep.Next()
			continue
		}
		if ep.parser.token.Type == lex.TOKEN_LB {
			pos := ep.parser.mkPos()
			ep.Next()                                    // skip [
			if ep.parser.token.Type == lex.TOKEN_COLON { // a[:]
				ep.Next() // skip :
				var end *ast.Expression
				if ep.parser.token.Type != lex.TOKEN_RB {
					end, err = ep.parseExpression()
					if err != nil {
						return nil, err
					}
				}
				if ep.parser.token.Type != lex.TOKEN_RB {
					return nil, fmt.Errorf("%s '[' and ']' not match", ep.parser.errorMsgPrefix())
				}
				ep.Next() // skip ]
				newe := &ast.Expression{}
				newe.Typ = ast.EXPRESSION_TYPE_SLICE
				newe.Pos = ep.parser.mkPos()
				slice := &ast.ExpressionSlice{}
				newe.Data = slice
				slice.Expression = left
				slice.End = end
				left = newe
				continue
			}
			e, err := ep.parseExpression()
			if err != nil {
				return nil, err
			}
			if e.Typ == ast.EXPRESSION_TYPE_LABLE || ep.parser.token.Type == lex.TOKEN_COLON {
				if e.Typ == ast.EXPRESSION_TYPE_LABLE && ep.parser.token.Type == lex.TOKEN_COLON {
					//
					ep.parser.Next()
				}
				if ep.parser.token.Type == lex.TOKEN_COLON {
					ep.parser.Next() // skip :
				}
				if e.Typ == ast.EXPRESSION_TYPE_LABLE {
					e.Typ = ast.EXPRESSION_TYPE_IDENTIFIER // corrent to identifier
				}
				var end *ast.Expression
				if ep.parser.token.Type != lex.TOKEN_RB {
					end, err = ep.parseExpression()
					if err != nil {
						return nil, err
					}
				}
				if ep.parser.token.Type != lex.TOKEN_RB {
					return nil, fmt.Errorf("%s '[' and ']' not match", ep.parser.errorMsgPrefix())
				}
				ep.Next() // skip ]
				newe := &ast.Expression{}
				newe.Typ = ast.EXPRESSION_TYPE_SLICE
				newe.Pos = ep.parser.mkPos()
				slice := &ast.ExpressionSlice{}
				newe.Data = slice
				slice.Start = e
				slice.Expression = left
				slice.End = end
				left = newe
				continue
			}
			if ep.parser.token.Type != lex.TOKEN_RB {
				return nil, fmt.Errorf("%s '[' and ']' not match", ep.parser.errorMsgPrefix())
			}
			newe := &ast.Expression{}
			newe.Pos = pos
			newe.Typ = ast.EXPRESSION_TYPE_INDEX
			index := &ast.ExpressionIndex{}
			index.Expression = left
			index.Index = e
			newe.Data = index
			left = newe
			ep.Next()
			continue
		}
		if ep.parser.token.Type == lex.TOKEN_DOT {
			pos := ep.parser.mkPos()
			ep.parser.Next() // skip .
			if ep.parser.token.Type != lex.TOKEN_IDENTIFIER {
				return nil, fmt.Errorf("%s not identifier after dot", ep.parser.errorMsgPrefix())
			}
			newe := &ast.Expression{}
			newe.Pos = pos
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

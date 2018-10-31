package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (ep *ExpressionParser) parseCallExpression(on *ast.Expression) (*ast.Expression, error) {
	var err error
	ep.Next(lfNotToken) // skip (
	args := []*ast.Expression{}
	if ep.parser.token.Type != lex.TokenRp { //a(123)
		args, err = ep.parseExpressions(lex.TokenRp)
		if err != nil {
			return nil, err
		}
	}
	ep.parser.ifTokenIsLfThenSkip()
	if ep.parser.token.Type != lex.TokenRp {
		err := fmt.Errorf("%s except ')' ,but '%s'",
			ep.parser.errMsgPrefix(),
			ep.parser.token.Description)
		ep.parser.errs = append(ep.parser.errs, err)
		return nil, err
	}
	pos := ep.parser.mkPos()
	ep.Next(lfIsToken) // skip )
	result := &ast.Expression{}
	if on.Type == ast.ExpressionTypeSelection {
		/*
			x.x()
		*/
		result.Type = ast.ExpressionTypeMethodCall
		result.Op = "methodCall"
		call := &ast.ExpressionMethodCall{}
		index := on.Data.(*ast.ExpressionSelection)
		call.Expression = index.Expression
		call.Args = args
		call.Name = index.Name
		result.Pos = on.Pos
		result.Data = call
	} else {
		result.Type = ast.ExpressionTypeFunctionCall
		result.Op = "functionCall"
		call := &ast.ExpressionFunctionCall{}
		call.Expression = on
		call.Args = args
		result.Data = call
		result.Pos = pos
	}

	if ep.parser.token.Type == lex.TokenLt { // <
		/*
			template function call return type binds
			fn a ()->(r T) {

			}
			a<int , ... >
		*/
		ep.Next(lfNotToken) // skip <
		ts, err := ep.parser.parseTypes(lex.TokenGt)
		if err != nil {
			return result, err
		}
		if ep.parser.token.Type != lex.TokenGt {
			ep.parser.errs = append(ep.parser.errs,
				fmt.Errorf("%s '<' and '>' not match",
					ep.parser.errMsgPrefix()))
			ep.parser.consume(untilGt)
		}
		ep.Next(lfIsToken)
		if result.Type == ast.ExpressionTypeFunctionCall {
			result.Data.(*ast.ExpressionFunctionCall).ParameterTypes = ts
		} else {
			result.Data.(*ast.ExpressionMethodCall).ParameterTypes = ts
		}
	}
	return result, nil
}

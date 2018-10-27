package parser

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (expressionParser *ExpressionParser) parseCallExpression(on *ast.Expression) (*ast.Expression, error) {
	var err error
	expressionParser.Next(lfNotToken) // skip (
	args := []*ast.Expression{}
	if expressionParser.parser.token.Type != lex.TokenRp { //a(123)
		args, err = expressionParser.parseExpressions(lex.TokenRp)
		if err != nil {
			return nil, err
		}
	}
	expressionParser.parser.ifTokenIsLfThenSkip()
	if expressionParser.parser.token.Type != lex.TokenRp {
		err := fmt.Errorf("%s except ')' ,but '%s'",
			expressionParser.parser.errMsgPrefix(),
			expressionParser.parser.token.Description)
		expressionParser.parser.errs = append(expressionParser.parser.errs, err)
		return nil, err
	}
	pos := expressionParser.parser.mkPos()
	expressionParser.Next(lfIsToken) // skip )
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

	if expressionParser.parser.token.Type == lex.TokenLt { // <
		/*
			template function call return type binds
			fn a ()->(r T) {

			}
			a<int , ... >
		*/
		expressionParser.Next(lfNotToken) // skip <
		ts, err := expressionParser.parser.parseTypes(lex.TokenGt)
		if err != nil {
			return result, err
		}
		if expressionParser.parser.token.Type != lex.TokenGt {
			expressionParser.parser.errs = append(expressionParser.parser.errs,
				fmt.Errorf("%s '<' and '>' not match",
					expressionParser.parser.errMsgPrefix()))
			expressionParser.parser.consume(untilGt)
		}
		expressionParser.Next(lfIsToken)
		if result.Type == ast.ExpressionTypeFunctionCall {
			result.Data.(*ast.ExpressionFunctionCall).ParameterTypes = ts
		} else {
			result.Data.(*ast.ExpressionMethodCall).ParameterTypes = ts
		}
	}
	return result, nil
}

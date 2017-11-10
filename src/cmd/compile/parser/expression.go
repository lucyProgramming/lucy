package parser

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/lex"
)

type ExpressionParser struct {
	parser *Parser
}

func (ep *ExpressionParser) parseExpressions() ([]*ast.Expression, error) {

}

//parse equal expression
func (ep *ExpressionParser) parseExpression() (*ast.Expression, error) {
	if ep.parser.token.Type == lex.TOKEN_IDENTIFIER {
		return ep.parseIdentifierExpression()
	}
}

//begin with identifier
func (ep *ExpressionParser) parseIdentifierExpression() {
	identifer := ep.parser.token.Data.(string)
}

//func (ep *ExpressionParser) parseEqualExpression() (*ast.Expression, error) {
//	ep.parser.Next()
//	if ep.parser.eof {
//		return nil, ep.parser.unexpectedErr()
//	}
//	switch ep.parser.token.Type {
//	case lex.TOKEN_COLON_ASSIGN:
//	case lex.TOKEN_PLUS_ASSIGN:
//	case lex.TOKEN_ADD_ASSIGN:

//	}

//}

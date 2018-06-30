package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (blockParser *BlockParser) parseIf() (i *ast.StatementIf, err error) {
	blockParser.Next() // skip if
	var e *ast.Expression
	e, err = blockParser.parser.ExpressionParser.parseExpression(false)
	if err != nil {
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		blockParser.consume(untilLc)
		blockParser.Next()
	}
	i = &ast.StatementIf{}
	i.Condition = e
	for blockParser.parser.token.Type == lex.TokenSemicolon {
		if i.Condition != nil {
			i.PrefixExpressions = append(i.PrefixExpressions, i.Condition)
		}
		blockParser.Next() // skip ;
		i.Condition, err = blockParser.parser.ExpressionParser.parseExpression(false)
		if err != nil {
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			blockParser.consume(untilLc)
			blockParser.Next()
		}
	}

	if blockParser.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s missing '{' after a expression,but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		blockParser.consume(untilLc) // consume and next
		blockParser.Next()
	}
	blockParser.Next() //skip {
	blockParser.parseStatementList(&i.Block, false)
	if blockParser.parser.token.Type != lex.TokenRc {
		blockParser.parser.errs = append(blockParser.parser.errs, fmt.Errorf("%s expect '}', but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description))
		blockParser.consume(untilRc)
	}
	blockParser.Next() // skip }
	//fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	if blockParser.parser.token.Type == lex.TokenElseif {
		es, err := blockParser.parseElseIfList()
		if err != nil {
			return i, err
		}
		i.ElseIfList = es
	}
	//fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	if blockParser.parser.token.Type == lex.TokenElse {
		blockParser.Next()
		if blockParser.parser.token.Type != lex.TokenLc {
			err = fmt.Errorf("%s missing '{' after else", blockParser.parser.errorMsgPrefix())
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			return i, err
		}
		i.ElseBlock = &ast.Block{}
		blockParser.Next()
		blockParser.parseStatementList(i.ElseBlock, false)
		if blockParser.parser.token.Type != lex.TokenRc {
			err = fmt.Errorf("%s expect '}', but '%s'",
				blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
			blockParser.consume(untilRc)
		}
		blockParser.Next()

	}
	return i, err
}

func (blockParser *BlockParser) parseElseIfList() (es []*ast.StatementElseIf, err error) {
	es = []*ast.StatementElseIf{}
	var e *ast.Expression
	for blockParser.parser.token.Type == lex.TokenElseif {
		blockParser.Next() // skip elseif token
		e, err = blockParser.parser.ExpressionParser.parseExpression(false)
		if err != nil {
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			return es, err
		}
		if blockParser.parser.token.Type != lex.TokenLc {
			err = fmt.Errorf("%s not '{' after a expression,but '%s'",
				blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
			blockParser.parser.errs = append(blockParser.parser.errs)
			return es, err
		}
		block := &ast.Block{}
		blockParser.Next()
		blockParser.parseStatementList(block, false)
		es = append(es, &ast.StatementElseIf{
			Condition: e,
			Block:     block,
		})
		if blockParser.parser.token.Type != lex.TokenRc {
			err = fmt.Errorf("%s expect '}', but '%s'",
				blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
			blockParser.consume(untilRc)
		}
		blockParser.Next() // skip }
	}
	return es, err
}

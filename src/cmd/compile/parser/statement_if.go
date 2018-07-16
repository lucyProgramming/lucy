package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (blockParser *BlockParser) parseIf() (i *ast.StatementIf, err error) {
	blockParser.Next() // skip if
	var condition *ast.Expression
	condition, err = blockParser.parser.ExpressionParser.parseExpression(false)
	if err != nil {
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		blockParser.consume(untilLc)
		blockParser.Next()
	}
	i = &ast.StatementIf{}
	i.Condition = condition
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
	if blockParser.parser.token.Type == lex.TokenElseif {
		i.ElseIfList, err = blockParser.parseElseIfList()
		if err != nil {
			return i, err
		}
	}
	if blockParser.parser.token.Type == lex.TokenElse {
		blockParser.Next()
		if blockParser.parser.token.Type != lex.TokenLc {
			err = fmt.Errorf("%s missing '{' after else", blockParser.parser.errorMsgPrefix())
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			blockParser.consume(untilLc)
		}
		blockParser.Next() // skip {
		i.ElseBlock = &ast.Block{}
		blockParser.parseStatementList(i.ElseBlock, false)
		if blockParser.parser.token.Type != lex.TokenRc {
			err = fmt.Errorf("%s expect '}', but '%s'",
				blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			blockParser.consume(untilRc)
		}
		blockParser.Next() // skip }
	}
	return i, err
}

func (blockParser *BlockParser) parseElseIfList() (elseIfList []*ast.StatementElseIf, err error) {
	elseIfList = []*ast.StatementElseIf{}
	var condition *ast.Expression
	for blockParser.parser.token.Type == lex.TokenElseif {
		blockParser.Next() // skip elseif token
		condition, err = blockParser.parser.ExpressionParser.parseExpression(false)
		if err != nil {
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			blockParser.consume(untilLc)
		}
		if blockParser.parser.token.Type != lex.TokenLc {
			err = fmt.Errorf("%s not '{' after a expression,but '%s'",
				blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
			blockParser.parser.errs = append(blockParser.parser.errs)
			blockParser.consume(untilLc)
		}
		blockParser.Next() // skip {
		block := &ast.Block{}
		blockParser.parseStatementList(block, false)
		elseIfList = append(elseIfList, &ast.StatementElseIf{
			Condition: condition,
			Block:     block,
		})
		if blockParser.parser.token.Type != lex.TokenRc {
			err = fmt.Errorf("%s expect '}', but '%s'",
				blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
			blockParser.parser.errs = append(blockParser.parser.errs)
			blockParser.consume(untilRc)
		}
		blockParser.Next() // skip }
	}
	return elseIfList, err
}

package parser

import (
	"fmt"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"
)

func (blockParser *BlockParser) parseSwitchTemplate(pos *ast.Position) (*ast.StatementSwitchTemplate, error) {
	condition, err := blockParser.parser.parseType()
	if err != nil {
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return nil, err
	}
	if blockParser.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return nil, err
	}
	blockParser.Next() // skip {  , must be case
	if blockParser.parser.token.Type != lex.TokenCase {
		err = fmt.Errorf("%s expect 'case',but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return nil, err
	}
	s := &ast.StatementSwitchTemplate{}
	s.Pos = pos
	s.Condition = condition
	for blockParser.parser.token.Type == lex.TokenCase {
		blockParser.Next() // skip case
		ts, err := blockParser.parser.parseTypes()
		if err != nil {
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			return s, err
		}
		if blockParser.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s expect ':',but '%s'",
				blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			return s, err
		}
		blockParser.Next() // skip :
		var block *ast.Block
		if blockParser.parser.token.Type != lex.TokenCase &&
			blockParser.parser.token.Type != lex.TokenDefault &&
			blockParser.parser.token.Type != lex.TokenRc {
			block = &ast.Block{}
			block.IsSwitchStatementTopBlock = true
			blockParser.parseStatementList(block, false)

		}
		s.StatementSwitchCases = append(s.StatementSwitchCases, &ast.StatementSwitchTemplateCase{
			Matches: ts,
			Block:   block,
		})
	}
	//default value
	if blockParser.parser.token.Type == lex.TokenDefault {
		blockParser.Next() // skip default key word
		if blockParser.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s missing clon after default",
				blockParser.parser.errorMsgPrefix())
			blockParser.parser.errs = append(blockParser.parser.errs, err)
		} else {
			blockParser.Next()
		}
		if blockParser.parser.token.Type != lex.TokenRc {
			block := ast.Block{}
			block.IsSwitchStatementTopBlock = true
			blockParser.parseStatementList(&block, false)
			s.Default = &block
		}
	}
	if blockParser.parser.token.Type != lex.TokenRc {
		err = fmt.Errorf("%s expect '}',but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return s, err
	}
	blockParser.Next() //  skip }
	return s, nil
}

func (blockParser *BlockParser) parseSwitch() (interface{}, error) {
	pos := blockParser.parser.mkPos()
	blockParser.Next() // skip switch key word
	if blockParser.parser.token.Type == lex.TokenTemplate {
		return blockParser.parseSwitchTemplate(pos)
	}
	condition, err := blockParser.parser.ExpressionParser.parseExpression(false)
	if err != nil {
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return nil, err
	}
	if blockParser.parser.token.Type != lex.TokenLc {
		err = fmt.Errorf("%s expect '{',but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return nil, err
	}
	blockParser.Next() // skip {  , must be case
	if blockParser.parser.token.Type != lex.TokenCase {
		err = fmt.Errorf("%s expect 'case',but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return nil, err
	}
	s := &ast.StatementSwitch{}
	s.Pos = pos
	s.Condition = condition
	for blockParser.parser.token.Type == lex.TokenCase {
		blockParser.Next() // skip case
		es, err := blockParser.parser.ExpressionParser.parseExpressions()
		if err != nil {
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			return s, err
		}
		if blockParser.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s expect ':',but '%s'",
				blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
			blockParser.parser.errs = append(blockParser.parser.errs, err)
			return s, err
		}
		blockParser.Next() // skip :
		var block *ast.Block
		if blockParser.parser.token.Type != lex.TokenCase &&
			blockParser.parser.token.Type != lex.TokenDefault &&
			blockParser.parser.token.Type != lex.TokenRc {
			block = &ast.Block{}
			block.IsSwitchStatementTopBlock = true
			blockParser.parseStatementList(block, false)

		}
		s.StatementSwitchCases = append(s.StatementSwitchCases, &ast.StatementSwitchCase{
			Matches: es,
			Block:   block,
		})
	}
	//default value
	if blockParser.parser.token.Type == lex.TokenDefault {
		blockParser.Next() // skip default key word
		if blockParser.parser.token.Type != lex.TokenColon {
			err = fmt.Errorf("%s missing clon after default",
				blockParser.parser.errorMsgPrefix())
			blockParser.parser.errs = append(blockParser.parser.errs, err)
		} else {
			blockParser.Next()
		}
		if blockParser.parser.token.Type != lex.TokenRc {
			block := ast.Block{}
			block.IsSwitchStatementTopBlock = true
			blockParser.parseStatementList(&block, false)
			s.Default = &block
		}
	}
	if blockParser.parser.token.Type != lex.TokenRc {
		err = fmt.Errorf("%s expect '}',but '%s'",
			blockParser.parser.errorMsgPrefix(), blockParser.parser.token.Description)
		blockParser.parser.errs = append(blockParser.parser.errs, err)
		return s, err
	}
	blockParser.Next() //  skip }
	return s, nil
}

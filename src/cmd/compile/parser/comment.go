package parser

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/lex"

type CommentParser struct {
	Comment string
	parser  *Parser
}

func (c *CommentParser) reset() {
	c.Comment = ""
}

func (c *CommentParser) read() {
	c.reset()
	if c.parser.token.Type == lex.TokenComment {
		for c.parser.token.Type == lex.TokenComment {
			c.Comment += c.parser.token.Data.(string)
			c.parser.Next(lfIsToken)
		}
	} else {
		c.Comment = c.parser.token.Data.(string)
		c.parser.Next(lfIsToken)
	}
	if c.parser.token.Type != lex.TokenLf {
		return
	}
	c.parser.Next(lfIsToken)
	if c.parser.token.Type == lex.TokenLf {
		c.reset()
		return
	}
}

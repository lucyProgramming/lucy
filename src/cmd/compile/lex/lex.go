package lex

import (
	"fmt"
	"math"
)

func New(bs []byte) *LucyLexer {
	lex := &LucyLexer{bs: bs}
	lex.end = len(bs)
	lex.line = 1
	lex.lastline = 1
	lex.column = 1
	lex.lastcolumn = 1
	return lex
}

type LucyLexer struct {
	bs                   []byte
	lastline, lastcolumn int
	line, column         int
	offset, end          int
}

func (lex *LucyLexer) Pos() (int, int) {
	return lex.line, lex.column
}

func (lex *LucyLexer) getchar() (c byte, eof bool) {
	if lex.offset == lex.end {
		eof = true
		return
	}
	offset := lex.offset
	lex.offset++
	c = lex.bs[offset]
	lex.lastline = lex.line
	lex.lastcolumn = lex.column
	if c == '\n' {
		lex.line++
		lex.column = 1
	} else {
		if c == '\t' {
			lex.column += 4 // TODO:: 4 OR 8
		} else {
			lex.column++
		}
	}
	return

}

func (lex *LucyLexer) ungetchar() {
	lex.offset--
	lex.line, lex.column = lex.lastline, lex.lastcolumn
}

func (lex *LucyLexer) isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z')
}
func (lex *LucyLexer) isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
func (lex *LucyLexer) isOctal(c byte) bool {
	return '0' <= c && c <= '7'
}
func (lex *LucyLexer) isHex(c byte) bool {
	return '0' <= c && c <= '9' ||
		('a' <= c && c <= 'f') ||
		('A' <= c && c <= 'F')
}

func (lex *LucyLexer) hexbyte2Byte(c byte) byte {
	if 'a' <= c && c <= 'f' {
		return c - 'a' + 10
	}
	if 'A' <= c && c <= 'F' {
		return c - 'A' + 10
	}
	return c - '0' //also valid for digit
}

func (lex *LucyLexer) parseInt(bs []byte) int64 {
	base := int64(10)
	if bs[0] == '0' {
		base = 8
	}
	if len(bs) >= 2 && bs[0] == '0' && (bs[1] == 'X' || bs[1] == 'x') { // correct base to hex
		base = 16
		bs = bs[2:]
	}
	var result int64 = 0
	for _, v := range bs {
		result = result*base + int64(lex.hexbyte2Byte(v))
	}
	return result
}

func (lex *LucyLexer) lexNumber(token *Token, c byte) (eof bool, err error) {
	integerpart := []byte{c}
	ishex := false
	isOctal := false
	if c == '0' { // enter when first char is '0'
		c, eof = lex.getchar()
		if c == 'x' || c == 'X' {
			ishex = true
			integerpart = append(integerpart, 'X')
		} else {
			isOctal = true
			lex.ungetchar()
		}
	}
	c, eof = lex.getchar() //get next char
	for eof == false {
		ok := false
		if ishex {
			ok = lex.isHex(c)
		} else if isOctal {
			if lex.isDigit(c) == true && lex.isOctal(c) == false { // integer but not octal
				err = fmt.Errorf("octal number cannot be '8' and '9'")
			}
			ok = lex.isDigit(c)
		} else {
			ok = lex.isDigit(c)
		}
		if ok {
			integerpart = append(integerpart, c)
			c, eof = lex.getchar() // get next char
		} else { // something that I cannot handle
			lex.ungetchar()
			break
		}
	}
	c, eof = lex.getchar()
	floatpart := []byte{}
	isfloat := false // float or double
	if c == '.' {    // float numbers
		isfloat = true
		c, eof = lex.getchar()
		for eof == false {
			if lex.isDigit(c) {
				floatpart = append(floatpart, c)
				c, eof = lex.getchar()
			} else {
				lex.ungetchar()
				break
			}
		}
	} else {
		lex.ungetchar()
	}
	if ishex && isfloat {
		token.Type = TOKEN_LITERAL_INT
		token.Data = 0
		err = fmt.Errorf("mix up float and hex")
		return
	}
	isdouble := false
	islong := false
	isShort := false
	isByte := false
	c, eof = lex.getchar()
	if c == 'l' || c == 'L' {
		islong = true
	} else if c == 'f' || c == 'F' {
		isfloat = true
	} else if c == 's' || c == 'S' {
		isShort = true
	} else if c == 'd' || c == 'D' {
		isdouble = true
	} else if c == 'b' || c == 'B' {
		isByte = true
	} else {
		lex.ungetchar()
	}
	isScientificNotation := false
	power := []byte{}
	powerPositive := true
	c, eof = lex.getchar()
	if (c == 'e' || c == 'E') && eof == false {
		isScientificNotation = true
		c, eof = lex.getchar()
		if eof {
			err = fmt.Errorf("unexpect EOF")
		}
		if c == '-' {
			powerPositive = false
			c, eof = lex.getchar()
		} else if lex.isDigit(c) { // nothing to do

		} else if c == '+' { // default is true
			c, eof = lex.getchar()
		} else {
			err = fmt.Errorf("wrong format scientific notation")
		}
		if lex.isDigit(c) == false {
			lex.ungetchar() //
			err = fmt.Errorf("wrong format scientific notation")
		} else {
			power = append(power, c)
			c, eof = lex.getchar()
			for eof == false && lex.isDigit(c) {
				power = append(power, c)
				c, eof = lex.getchar()
			}
			lex.ungetchar()
		}
	} else {
		lex.ungetchar()
	}
	if ishex && isScientificNotation {
		token.Type = TOKEN_LITERAL_INT
		token.Data = 0
		token.Desp = "0"
		err = fmt.Errorf("mix up hex and seientific notation")
		return
	}
	/*
		parse float part
	*/
	parseFloat := func(bs []byte) float64 {
		index := len(bs) - 1
		var fp float64
		for index >= 0 {
			fp = fp*0.1 + (float64(lex.hexbyte2Byte(bs[index])) / 10.0)
			index--
		}
		return fp
	}
	token.EndLine = lex.line
	token.EndColumn = lex.column
	if isScientificNotation == false {
		if isfloat {
			value := parseFloat(floatpart) + float64(lex.parseInt(integerpart))
			if isdouble {
				token.Type = TOKEN_LITERAL_DOUBLE
				token.Data = value
			} else {
				token.Type = TOKEN_LITERAL_FLOAT
				token.Data = float32(value)
			}
		} else {
			value := lex.parseInt(integerpart)
			if islong {
				token.Type = TOKEN_LITERAL_LONG
				token.Data = value
			} else if isByte {
				token.Type = TOKEN_LITERAL_BYTE
				token.Data = byte(value)
				if int32(value) > math.MaxUint8 {
					err = fmt.Errorf("max byte is %v", math.MaxUint8)
				}
			} else if isShort {
				token.Type = TOKEN_LITERAL_SHORT
				token.Data = int32(value)
				if int32(value) > math.MaxInt16 {
					err = fmt.Errorf("max short is %v", math.MaxInt16)
				}
			} else {
				token.Type = TOKEN_LITERAL_INT
				token.Data = int32(value)
			}
		}
		return
	}
	//scientific notation
	if t := lex.parseInt(integerpart); t > 10 || t < 1 {
		err = fmt.Errorf("wrong format scientific notation")
		token.Type = TOKEN_LITERAL_INT
		token.Data = 0
		return
	}
	p := int(lex.parseInt(power))
	if powerPositive {
		if p >= len(floatpart) { // int
			integerpart = append(integerpart, floatpart...)
			b := make([]byte, p-len(floatpart))
			for k, _ := range b {
				b[k] = '0'
			}
			integerpart = append(integerpart, b...)
			value := lex.parseInt(integerpart)
			token.Type = TOKEN_LITERAL_INT
			token.Data = int32(value)
		} else { // float
			integerpart = append(integerpart, floatpart[0:p]...)
			fmt.Println(floatpart[p:], parseFloat(floatpart[p:]))
			value := float64(lex.parseInt(integerpart)) + parseFloat(floatpart[p:])
			token.Type = TOKEN_LITERAL_FLOAT
			token.Data = value
		}
	} else { // power is negative,must be float number
		b := make([]byte, p-len(integerpart))
		for k, _ := range b {
			b[k] = '0'
		}
		b = append(b, integerpart...)
		b = append(b, floatpart...)
		value := parseFloat(b)
		token.Type = TOKEN_LITERAL_FLOAT
		token.Data = value
	}
	return
}
func (lex *LucyLexer) lexIdentifier(c byte) (token *Token, err error) {
	token = &Token{}
	token.StartLine = lex.line
	token.StartColumn = lex.column - 1 // c is readed
	bs := []byte{c}
	c, eof := lex.getchar()
	for eof == false {
		if lex.isLetter(c) || c == '_' || lex.isDigit(c) || c == '$' {
			bs = append(bs, c)
			c, eof = lex.getchar()
		} else {
			lex.ungetchar()
			break
		}
	}
	token.EndLine = lex.line
	token.EndColumn = lex.column
	identifier := string(bs)
	if t, ok := keywordMap[identifier]; ok {
		token.Type = t
		token.Desp = identifier
		if token.Type == TOKEN_ELSE {
			is := lex.tryLexElseIf()
			if is {
				token.Type = TOKEN_ELSEIF
				token.Desp = "else if"
			}
		}
	} else {
		token.Type = TOKEN_IDENTIFIER
		token.Data = identifier
		token.Desp = "identifer_" + identifier
	}
	token.EndLine = lex.line
	token.EndColumn = lex.column
	return
}

func (lex *LucyLexer) tryLexElseIf() (is bool) {
	c, eof := lex.getchar()
	for (c == ' ' || c == '\t' || c == '\r') && eof == false {
		c, eof = lex.getchar()
	}
	if eof {
		return
	}
	if c != 'i' {
		lex.ungetchar()
		return
	}
	c, eof = lex.getchar()
	if c != 'f' {
		lex.ungetchar()
		lex.ungetchar()
		return
	}
	c, eof = lex.getchar()
	if c != ' ' && c != '\t' && c != '\r' { // white list
		lex.ungetchar()
		lex.ungetchar()
		lex.ungetchar()
		return
	}
	is = true
	return
}

func (lex *LucyLexer) lexString(endc byte) (token *Token, err error) {
	token = &Token{}
	token.StartLine = lex.line
	token.StartColumn = lex.column
	token.Type = TOKEN_LITERAL_STRING
	bs := []byte{}
	var c byte
	c, eof := lex.getchar()
	for c != endc && c != '\n' && eof == false {
		if c != '\\' {
			bs = append(bs, c)
			c, eof = lex.getchar()
			continue
		}
		c, eof = lex.getchar() // get next char
		if eof {
			err = fmt.Errorf("unexpected EOF")
			break
		}
		switch c {
		case 'a':
			bs = append(bs, '\a')
			c, eof = lex.getchar()
		case 'b':
			bs = append(bs, '\b')
			c, eof = lex.getchar()
		case 'f':
			bs = append(bs, '\f')
			c, eof = lex.getchar()
		case 'n':
			bs = append(bs, '\n')
			c, eof = lex.getchar()
		case 'r':
			bs = append(bs, '\r')
			c, eof = lex.getchar()
		case 't':
			bs = append(bs, '\t')
			c, eof = lex.getchar()
		case 'v':
			bs = append(bs, '\v')
			c, eof = lex.getchar()
		case '\\':
			bs = append(bs, '\\')
			c, eof = lex.getchar()
		case '\'':
			bs = append(bs, '\'')
			c, eof = lex.getchar()
		case '"':
			bs = append(bs, '"')
			c, eof = lex.getchar()
		case 'x', 'X':
			var c1, c2 byte
			c1, eof = lex.getchar() //skip 'x' or 'X'
			if eof {
				err = fmt.Errorf("unexpect EOF")
				continue
			}
			if !lex.isHex(c) {
				err = fmt.Errorf("unknown escape sequence")
				continue
			}
			b := lex.hexbyte2Byte(c1)
			c2, eof = lex.getchar()
			if lex.isHex(c2) {
				if t := b*16 + lex.hexbyte2Byte(c2); t < 127 { // only support standard ascii
					b = t
				} else {
					lex.ungetchar()
				}
			} else { //not hex
				lex.ungetchar()
			}
			bs = append(bs, b)
			c, eof = lex.getchar()
		case '0', '1', '2', '3', '4', '5', '7':
			// first char must be octal
			b := byte(0)
			for i := 0; i < 3; i++ {
				if eof {
					break
				}
				if lex.isOctal(c) == false {
					lex.ungetchar()
					break
				}
				if t := b*8 + lex.hexbyte2Byte(c); t > 127 { // only support standard ascii
					lex.ungetchar()
					break
				} else {
					b = t
				}
				c, eof = lex.getchar()
			}
			bs = append(bs, b)
			c, eof = lex.getchar()
		default:
			err = fmt.Errorf("unknown escape sequence")
		}
	}
	if c == '\n' {
		err = fmt.Errorf("string literal start new line")
	}
	token.EndLine = lex.line
	token.EndColumn = lex.column
	token.Data = string(bs)
	token.Desp = string(bs)
	return
}

func (lex *LucyLexer) lexMultiLineComment() {
redo:
	c, eof := lex.getchar()
	if eof {
		return
	}
	for c != '*' && eof == false {
		c, eof = lex.getchar()
	}
	if eof {
		return
	}
	c, eof = lex.getchar()
	if eof || c == '/' {
		return
	}
	goto redo
}

func (lex *LucyLexer) Next() (token *Token, err error) {
redo:
	token = &Token{}
	var c byte
	c, eof := lex.getchar()
	if eof {
		token.StartLine = lex.line
		token.StartColumn = lex.column
		token.Type = TOKEN_EOF
		token.Desp = "EOF"
		return
	}
	for c == ' ' || c == '\t' || c == '\r' {
		c, eof = lex.getchar()
	}
	if eof {
		token.StartLine = lex.line
		token.StartColumn = lex.column
		token.Type = TOKEN_EOF
		token.Desp = "EOF"
		return
	}
	if lex.isLetter(c) || c == '_' || c == '$' {
		return lex.lexIdentifier(c)
	}
	if lex.isDigit(c) {
		eof, err = lex.lexNumber(token, c)
		return
	}
	token.StartLine = lex.line
	token.StartColumn = lex.column
	switch c {
	case '(':
		token.Type = TOKEN_LP
		token.Desp = "("
	case ')':
		token.Type = TOKEN_RP
		token.Desp = ")"
	case '{':
		token.Type = TOKEN_LC
		token.Desp = "{"
	case '}':
		token.Type = TOKEN_RC
		token.Desp = "}"
	case '[':
		token.Type = TOKEN_LB
		token.Desp = "["
	case ']':
		token.Type = TOKEN_RB
		token.Desp = "]"
	case ';':
		token.Type = TOKEN_SEMICOLON
		token.Desp = ";"
	case ',':
		token.Type = TOKEN_COMMA
		token.Desp = ","
	case '&':
		c, eof = lex.getchar()
		if c == '&' {
			token.Type = TOKEN_LOGICAL_AND
			token.Desp = "&&"
		} else if c == '=' {
			token.Type = TOKEN_AND_ASSIGN
			token.Desp = "&="
		} else {
			lex.ungetchar()
			token.Type = TOKEN_AND
			token.Desp = "&"
		}
	case '|':
		c, eof = lex.getchar()
		if c == '|' {
			token.Type = TOKEN_LOGICAL_OR
			token.Desp = "||"
		} else if c == '=' {
			token.Type = TOKEN_OR_ASSIGN
			token.Desp = "|="
		} else {
			lex.ungetchar()
			token.Type = TOKEN_OR
			token.Desp = "|"
		}
	case '=':
		c, eof = lex.getchar()
		if c == '=' {
			token.Type = TOKEN_EQUAL
			token.Desp = "=="
		} else {
			lex.ungetchar()
			token.Type = TOKEN_ASSIGN
			token.Desp = "="
		}
	case '!':
		c, eof = lex.getchar()
		if c == '=' {
			token.Type = TOKEN_NE
			token.Desp = "!="
		} else {
			lex.ungetchar()
			token.Type = TOKEN_NOT
			token.Desp = "!"
		}
	case '>':
		c, eof = lex.getchar()
		if c == '=' {
			token.Type = TOKEN_GE
			token.Desp = ">="
		} else if c == '>' {
			c, eof = lex.getchar()
			if c == '=' {
				token.Type = TOKEN_RSH_ASSIGN
				token.Desp = ">>="
			} else {
				lex.ungetchar()
				token.Type = TOKEN_RIGHT_SHIFT
				token.Desp = ">>"
			}
		} else {
			lex.ungetchar()
			token.Type = TOKEN_GT
			token.Desp = ">"
		}
	case '<':
		c, eof = lex.getchar()
		if c == '=' {
			token.Type = TOKEN_LE
			token.Desp = "<="
		} else if c == '<' {
			c, eof = lex.getchar()
			if c == '=' {
				token.Type = TOKEN_LSH_ASSIGN
				token.Desp = "<<="
			} else {
				lex.ungetchar()
				token.Type = TOKEN_LEFT_SHIFT
				token.Desp = "<<"
			}
		} else {
			lex.ungetchar()
			token.Type = TOKEN_LT
			token.Desp = "<"
		}
	case '^':
		c, eof = lex.getchar()
		if c == '=' {
			token.Type = TOKEN_XOR_ASSIGN
			token.Desp = "^="
		} else {
			lex.ungetchar()
			token.Type = TOKEN_XOR
			token.Desp = "^"
		}
	case '~':
		token.Type = TOKEN_BITWISE_COMPLEMENT
		token.Desp = "~"
	case '+':
		c, eof = lex.getchar()
		if c == '+' {
			token.Type = TOKEN_INCREMENT
			token.Desp = "++"
		} else if c == '=' {
			token.Type = TOKEN_ADD_ASSIGN
			token.Desp = "+="
		} else {
			lex.ungetchar()
			token.Type = TOKEN_ADD
			token.Desp = "+"
		}
	case '-':
		c, eof = lex.getchar()
		if c == '-' {
			token.Type = TOKEN_DECREMENT
			token.Desp = "--"
		} else if c == '=' {
			token.Type = TOKEN_SUB_ASSIGN
			token.Desp = "-="
		} else if c == '>' {
			token.Type = TOKEN_ARROW
			token.Desp = "->"
		} else {
			lex.ungetchar()
			token.Type = TOKEN_SUB
			token.Desp = "-"
		}
	case '*':
		c, eof = lex.getchar()
		if c == '=' {
			token.Type = TOKEN_MUL_ASSIGN
			token.Desp = "*="
		} else {
			lex.ungetchar()
			token.Type = TOKEN_MUL
			token.Desp = "*"
		}
	case '%':
		c, eof = lex.getchar()
		if c == '=' {
			token.Type = TOKEN_MOD_ASSIGN
			token.Desp = "%="
		} else {
			lex.ungetchar()
			token.Type = TOKEN_MOD
			token.Desp = "%"
		}
	case '/':
		c, eof = lex.getchar()
		if c == '=' {
			token.Type = TOKEN_DIV_ASSIGN
			token.Desp = "/="
		} else if c == '/' {
			for c != '\n' && eof == false {
				c, eof = lex.getchar()
			}
			goto redo
		} else if c == '*' {
			lex.lexMultiLineComment()
			goto redo
		} else {
			lex.ungetchar()
			token.Type = TOKEN_DIV
			token.Desp = "/"
		}
	case '\n':
		token.Type = TOKEN_CRLF
		token.Desp = "\n"
	case '.':
		token.Type = TOKEN_DOT
		token.Desp = "."
	case '`':
		bs := []byte{}
		c, eof = lex.getchar()
		for c != '`' && eof == false {
			bs = append(bs, c)
			c, eof = lex.getchar()
		}
		token.Type = TOKEN_LITERAL_STRING
		token.Data = string(bs)
		token.Desp = string(bs)
	case '"':
		return lex.lexString('"')
	case '\'':
		token, err = lex.lexString('\'')
		if err == nil {
			if t := []byte(token.Data.(string)); len(t) != 1 {
				err = fmt.Errorf("expect one char")
			} else {
				token.Type = TOKEN_LITERAL_BYTE
				token.Data = byte([]byte(t)[0])
			}
		}
		return
	case ':':
		c, eof = lex.getchar()
		if c == '=' {
			token.Type = TOKEN_COLON_ASSIGN
			token.Desp = ":= "
		} else {
			token.Type = TOKEN_COLON
			token.Desp = ":"
			lex.ungetchar()
		}

	default:
		err = fmt.Errorf("unkown beginning of token:%d", c)
		return
	}
	token.EndLine = lex.line
	token.EndColumn = lex.column
	return
}

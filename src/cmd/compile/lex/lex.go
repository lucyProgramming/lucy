package lex

import (
	"fmt"
	"math"
)

type LucyLexer struct {
	bs                   []byte
	lastLine, lastColumn int
	line, column         int
	offset, end          int
}

func (lex *LucyLexer) Pos() (int, int) {
	return lex.line, lex.column
}

func (lex *LucyLexer) GetOffSet() int {
	return lex.offset
}

func (lex *LucyLexer) getChar() (c byte, eof bool) {
	if lex.offset == lex.end {
		eof = true
		return
	}
	offset := lex.offset
	lex.offset++
	c = lex.bs[offset]
	lex.lastLine = lex.line
	lex.lastColumn = lex.column
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

func (lex *LucyLexer) unGetChar() {
	lex.offset--
	lex.line, lex.column = lex.lastLine, lex.lastColumn
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

func (lex *LucyLexer) hexByte2Byte(c byte) byte {
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
		result = result*base + int64(lex.hexByte2Byte(v))
	}
	return result
}

func (lex *LucyLexer) lexNumber(token *Token, c byte) (eof bool, err error) {
	integerPart := []byte{c}
	isHex := false
	isOctal := false
	if c == '0' { // enter when first char is '0'
		c, eof = lex.getChar()
		if c == 'x' || c == 'X' {
			isHex = true
			integerPart = append(integerPart, 'X')
		} else {
			isOctal = true
			lex.unGetChar()
		}
	}
	c, eof = lex.getChar() //get next char
	for eof == false {
		ok := false
		if isHex {
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
			integerPart = append(integerPart, c)
			c, eof = lex.getChar() // get next char
		} else { // something that I cannot handle
			lex.unGetChar()
			break
		}
	}
	c, eof = lex.getChar()
	floatPart := []byte{}
	isFloat := false // float or double
	if c == '.' {    // float numbers
		isFloat = true
		c, eof = lex.getChar()
		for eof == false {
			if lex.isDigit(c) {
				floatPart = append(floatPart, c)
				c, eof = lex.getChar()
			} else {
				lex.unGetChar()
				break
			}
		}
	} else {
		lex.unGetChar()
	}
	if isHex && isFloat {
		token.Type = TOKEN_LITERAL_INT
		token.Data = 0

		err = fmt.Errorf("mix up float and hex")
		return
	}
	isDouble := false
	isLong := false
	isShort := false
	isByte := false
	c, eof = lex.getChar()
	if c == 'l' || c == 'L' {
		isLong = true
	} else if c == 'f' || c == 'F' {
		isFloat = true
	} else if c == 's' || c == 'S' {
		isShort = true
	} else if c == 'd' || c == 'D' {
		isDouble = true
	} else if c == 'b' || c == 'B' {
		isByte = true
	} else {
		lex.unGetChar()
	}
	isScientificNotation := false
	power := []byte{}
	powerPositive := true
	c, eof = lex.getChar()
	if (c == 'e' || c == 'E') && eof == false {
		isScientificNotation = true
		c, eof = lex.getChar()
		if eof {
			err = fmt.Errorf("unexpect EOF")
		}
		if c == '-' {
			powerPositive = false
			c, eof = lex.getChar()
		} else if lex.isDigit(c) { // nothing to do

		} else if c == '+' { // default is true
			c, eof = lex.getChar()
		} else {
			err = fmt.Errorf("wrong format scientific notation")
		}
		if lex.isDigit(c) == false {
			lex.unGetChar() //
			err = fmt.Errorf("wrong format scientific notation")
		} else {
			power = append(power, c)
			c, eof = lex.getChar()
			for eof == false && lex.isDigit(c) {
				power = append(power, c)
				c, eof = lex.getChar()
			}
			lex.unGetChar()
		}
	} else {
		lex.unGetChar()
	}
	if isHex && isScientificNotation {
		token.Type = TOKEN_LITERAL_INT
		token.Data = 0
		token.Description = "0"
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
			fp = fp*0.1 + (float64(lex.hexByte2Byte(bs[index])) / 10.0)
			index--
		}
		return fp
	}
	token.EndLine = lex.line
	token.EndColumn = lex.column
	if isScientificNotation == false {
		if isFloat {
			value := parseFloat(floatPart) + float64(lex.parseInt(integerPart))
			if isDouble {
				token.Type = TOKEN_LITERAL_DOUBLE
				token.Data = value
			} else {
				token.Type = TOKEN_LITERAL_FLOAT
				token.Data = float32(value)
			}
		} else {
			value := lex.parseInt(integerPart)
			if isLong {
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
	if t := lex.parseInt(integerPart); t > 10 || t < 1 {
		err = fmt.Errorf("wrong format scientific notation")
		token.Type = TOKEN_LITERAL_INT
		token.Data = 0
		return
	}
	p := int(lex.parseInt(power))
	if powerPositive {
		if p >= len(floatPart) { // int
			integerPart = append(integerPart, floatPart...)
			b := make([]byte, p-len(floatPart))
			for k, _ := range b {
				b[k] = '0'
			}
			integerPart = append(integerPart, b...)
			value := lex.parseInt(integerPart)
			token.Type = TOKEN_LITERAL_INT
			token.Data = int32(value)
		} else { // float
			integerPart = append(integerPart, floatPart[0:p]...)
			fmt.Println(floatPart[p:], parseFloat(floatPart[p:]))
			value := float64(lex.parseInt(integerPart)) + parseFloat(floatPart[p:])
			token.Type = TOKEN_LITERAL_FLOAT
			token.Data = value
		}
	} else { // power is negative,must be float number
		b := make([]byte, p-len(integerPart))
		for k, _ := range b {
			b[k] = '0'
		}
		b = append(b, integerPart...)
		b = append(b, floatPart...)
		value := parseFloat(b)
		token.Type = TOKEN_LITERAL_FLOAT
		token.Data = value
	}
	return
}
func (lex *LucyLexer) looksLikeT(bs []byte) bool {
	if len(bs) == 0 {
		return false
	}
	if bs[0] != 'T' {
		return false
	}
	bs = bs[1:]
	for _, v := range bs {
		if !(v >= '0' && v <= '9') {
			return false
		}
	}
	return true
}

func (lex *LucyLexer) lexIdentifier(c byte) (token *Token, err error) {
	token = &Token{}
	token.StartLine = lex.line
	token.StartColumn = lex.column - 1 // c is readed
	bs := []byte{c}
	token.Offset = lex.offset - 1 // readed
	c, eof := lex.getChar()
	for eof == false {
		if lex.isLetter(c) || c == '_' || lex.isDigit(c) || c == '$' {
			bs = append(bs, c)
			c, eof = lex.getChar()
		} else {
			lex.unGetChar()
			break
		}
	}
	token.EndLine = lex.line
	token.EndColumn = lex.column
	identifier := string(bs)
	if t, ok := keywordsMap[identifier]; ok {
		token.Type = t
		token.Description = identifier
		if token.Type == TOKEN_ELSE {
			is := lex.tryLexElseIf()
			if is {
				token.Type = TOKEN_ELSEIF
				token.Description = "else if"
			}
		}
	} else {
		if lex.looksLikeT(bs) {
			token.Type = TOKEN_T
			token.Data = identifier
			token.Description = identifier
		} else {
			token.Type = TOKEN_IDENTIFIER
			token.Data = identifier
			token.Description = "identifier_" + identifier
		}
	}
	token.EndLine = lex.line
	token.EndColumn = lex.column
	return
}

func (lex *LucyLexer) tryLexElseIf() (is bool) {
	c, eof := lex.getChar()
	for (c == ' ' || c == '\t' || c == '\r') && eof == false {
		c, eof = lex.getChar()
	}
	if eof {
		return
	}
	if c != 'i' {
		lex.unGetChar()
		return
	}
	c, eof = lex.getChar()
	if c != 'f' {
		lex.unGetChar()
		lex.unGetChar()
		return
	}
	c, eof = lex.getChar()
	if c != ' ' && c != '\t' && c != '\r' { // white list
		lex.unGetChar()
		lex.unGetChar()
		lex.unGetChar()
		return
	}
	is = true
	return
}

func (lex *LucyLexer) lexString(endChar byte) (token *Token, err error) {
	token = &Token{}
	token.StartLine = lex.line
	token.StartColumn = lex.column
	token.Type = TOKEN_LITERAL_STRING
	bs := []byte{}
	var c byte
	c, eof := lex.getChar()
	for c != endChar && c != '\n' && eof == false {
		if c != '\\' {
			bs = append(bs, c)
			c, eof = lex.getChar()
			continue
		}
		c, eof = lex.getChar() // get next char
		if eof {
			err = fmt.Errorf("unexpected EOF")
			break
		}
		switch c {
		case 'a':
			bs = append(bs, '\a')
			c, eof = lex.getChar()
		case 'b':
			bs = append(bs, '\b')
			c, eof = lex.getChar()
		case 'f':
			bs = append(bs, '\f')
			c, eof = lex.getChar()
		case 'n':
			bs = append(bs, '\n')
			c, eof = lex.getChar()
		case 'r':
			bs = append(bs, '\r')
			c, eof = lex.getChar()
		case 't':
			bs = append(bs, '\t')
			c, eof = lex.getChar()
		case 'v':
			bs = append(bs, '\v')
			c, eof = lex.getChar()
		case '\\':
			bs = append(bs, '\\')
			c, eof = lex.getChar()
		case '\'':
			bs = append(bs, '\'')
			c, eof = lex.getChar()
		case '"':
			bs = append(bs, '"')
			c, eof = lex.getChar()
		case 'x':
			var c1, c2 byte
			c1, eof = lex.getChar() //skip 'x'
			if eof {
				err = fmt.Errorf("unexpect EOF")
				continue
			}
			if !lex.isHex(c) {
				err = fmt.Errorf("unknown escape sequence")
				continue
			}
			b := lex.hexByte2Byte(c1)
			c2, eof = lex.getChar()
			if lex.isHex(c2) {
				if t := b*16 + lex.hexByte2Byte(c2); t < 127 { // only support standard ascii
					b = t
				} else {
					lex.unGetChar()
				}
			} else { //not hex
				lex.unGetChar()
			}
			bs = append(bs, b)
			c, eof = lex.getChar()
		case '0', '1', '2', '3', '4', '5', '7':
			// first char must be octal
			b := byte(0)
			for i := 0; i < 3; i++ {
				if eof {
					break
				}
				if lex.isOctal(c) == false {
					lex.unGetChar()
					break
				}
				if t := b*8 + lex.hexByte2Byte(c); t > 127 { // only support standard ascii
					lex.unGetChar()
					break
				} else {
					b = t
				}
				c, eof = lex.getChar()
			}
			bs = append(bs, b)
			c, eof = lex.getChar()
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
	token.Description = string(bs)
	return
}

func (lex *LucyLexer) lexMultiLineComment() {
redo:
	c, eof := lex.getChar()
	if eof {
		return
	}
	for c != '*' && eof == false {
		c, eof = lex.getChar()
	}
	if eof {
		return
	}
	c, eof = lex.getChar()
	if eof || c == '/' {
		return
	}
	goto redo
}

func (lex *LucyLexer) Next() (token *Token, err error) {
redo:
	token = &Token{}
	var c byte
	token.StartLine = lex.line
	token.StartColumn = lex.column
	c, eof := lex.getChar()
	if eof {
		token.Type = TOKEN_EOF
		token.Description = "EOF"
		return
	}
	for c == ' ' || c == '\t' || c == '\r' {
		token.StartLine = lex.line
		token.StartColumn = lex.column
		c, eof = lex.getChar()
	}
	if eof {
		token.Type = TOKEN_EOF
		token.Description = "EOF"
		return
	}
	if lex.isLetter(c) || c == '_' || c == '$' {
		return lex.lexIdentifier(c)
	}
	if lex.isDigit(c) {
		eof, err = lex.lexNumber(token, c)
		return
	}
	token.Offset = lex.offset
	switch c {
	case '?':
		token.Type = TOKEN_QUESTION
		token.Description = "?"
	case '(':
		token.Type = TOKEN_LP
		token.Description = "("
	case ')':
		token.Type = TOKEN_RP
		token.Description = ")"
	case '{':
		token.Type = TOKEN_LC
		token.Description = "{"
	case '}':
		token.Type = TOKEN_RC
		token.Description = "}"
	case '[':
		token.Type = TOKEN_LB
		token.Description = "["
	case ']':
		token.Type = TOKEN_RB
		token.Description = "]"
	case ';':
		token.Type = TOKEN_SEMICOLON
		token.Description = ";"
	case ',':
		token.Type = TOKEN_COMMA
		token.Description = ","
	case '&':
		c, eof = lex.getChar()
		if c == '&' {
			token.Type = TOKEN_LOGICAL_AND
			token.Description = "&&"
		} else if c == '=' {
			token.Type = TOKEN_AND_ASSIGN
			token.Description = "&="
		} else {
			lex.unGetChar()
			token.Type = TOKEN_AND
			token.Description = "&"
		}
	case '|':
		c, eof = lex.getChar()
		if c == '|' {
			token.Type = TOKEN_LOGICAL_OR
			token.Description = "||"
		} else if c == '=' {
			token.Type = TOKEN_OR_ASSIGN
			token.Description = "|="
		} else {
			lex.unGetChar()
			token.Type = TOKEN_OR
			token.Description = "|"
		}
	case '=':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TOKEN_EQUAL
			token.Description = "=="
		} else {
			lex.unGetChar()
			token.Type = TOKEN_ASSIGN
			token.Description = "="
		}
	case '!':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TOKEN_NE
			token.Description = "!="
		} else {
			lex.unGetChar()
			token.Type = TOKEN_NOT
			token.Description = "!"
		}
	case '>':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TOKEN_GE
			token.Description = ">="
		} else if c == '>' {
			c, eof = lex.getChar()
			if c == '=' {
				token.Type = TOKEN_RSH_ASSIGN
				token.Description = ">>="
			} else {
				lex.unGetChar()
				token.Type = TOKEN_RIGHT_SHIFT
				token.Description = ">>"
			}
		} else {
			lex.unGetChar()
			token.Type = TOKEN_GT
			token.Description = ">"
		}
	case '<':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TOKEN_LE
			token.Description = "<="
		} else if c == '<' {
			c, eof = lex.getChar()
			if c == '=' {
				token.Type = TOKEN_LSH_ASSIGN
				token.Description = "<<="
			} else {
				lex.unGetChar()
				token.Type = TOKEN_LEFT_SHIFT
				token.Description = "<<"
			}
		} else {
			lex.unGetChar()
			token.Type = TOKEN_LT
			token.Description = "<"
		}
	case '^':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TOKEN_XOR_ASSIGN
			token.Description = "^="
		} else {
			lex.unGetChar()
			token.Type = TOKEN_XOR
			token.Description = "^"
		}
	case '~':
		token.Type = TOKEN_BITWISE_NOT
		token.Description = "~"
	case '+':
		c, eof = lex.getChar()
		if c == '+' {
			token.Type = TOKEN_INCREMENT
			token.Description = "++"
		} else if c == '=' {
			token.Type = TOKEN_ADD_ASSIGN
			token.Description = "+="
		} else {
			lex.unGetChar()
			token.Type = TOKEN_ADD
			token.Description = "+"
		}
	case '-':
		c, eof = lex.getChar()
		if c == '-' {
			token.Type = TOKEN_DECREMENT
			token.Description = "--"
		} else if c == '=' {
			token.Type = TOKEN_SUB_ASSIGN
			token.Description = "-="
		} else if c == '>' {
			token.Type = TOKEN_ARROW
			token.Description = "->"
		} else {
			lex.unGetChar()
			token.Type = TOKEN_SUB
			token.Description = "-"
		}
	case '*':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TOKEN_MUL_ASSIGN
			token.Description = "*="
		} else {
			lex.unGetChar()
			token.Type = TOKEN_MUL
			token.Description = "*"
		}
	case '%':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TOKEN_MOD_ASSIGN
			token.Description = "%="
		} else {
			lex.unGetChar()
			token.Type = TOKEN_MOD
			token.Description = "%"
		}
	case '/':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TOKEN_DIV_ASSIGN
			token.Description = "/="
		} else if c == '/' {
			for c != '\n' && eof == false {
				c, eof = lex.getChar()
			}
			goto redo
		} else if c == '*' {
			lex.lexMultiLineComment()
			goto redo
		} else {
			lex.unGetChar()
			token.Type = TOKEN_DIV
			token.Description = "/"
		}
	case '\n':
		token.Type = TOKEN_CRLF
		token.Description = "\n"
	case '.':
		token.Type = TOKEN_DOT
		token.Description = "."
	case '`':
		bs := []byte{}
		c, eof = lex.getChar()
		for c != '`' && eof == false {
			bs = append(bs, c)
			c, eof = lex.getChar()
		}
		token.Type = TOKEN_LITERAL_STRING
		token.Data = string(bs)
		token.Description = string(bs)
	case '"':
		return lex.lexString('"')
	case '\'':
		token, err = lex.lexString('\'')
		if err == nil {
			if t := []byte(token.Data.(string)); len(t) != 1 {
				err = fmt.Errorf("expect one char")
			} else { // correct token
				token.Type = TOKEN_LITERAL_BYTE
				token.Data = byte([]byte(t)[0])
			}
		}
		return
	case ':':
		c, eof = lex.getChar()
		if c == '=' {
			token.Type = TOKEN_COLON_ASSIGN
			token.Description = ":= "
		} else {
			token.Type = TOKEN_COLON
			token.Description = ":"
			lex.unGetChar()
		}

	default:
		err = fmt.Errorf("unkown beginning of token:%d", c)
		return
	}
	token.EndLine = lex.line
	token.EndColumn = lex.column
	return
}

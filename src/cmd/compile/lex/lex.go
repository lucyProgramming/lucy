package lex

import (
	"fmt"
	"strconv"
)

type Lexer struct {
	bs                   []byte
	lastLine, lastColumn int
	line, column         int
	offset, end          int
}

func (this *Lexer) GetLineAndColumn() (int, int) {
	return this.line, this.column
}
func (this *Lexer) GetOffSet() int {
	return this.offset
}

func (this *Lexer) getChar() (c byte, eof bool) {
	if this.offset == this.end {
		eof = true
		return
	}
	offset := this.offset
	this.offset++ // next
	c = this.bs[offset]
	this.lastLine = this.line
	this.lastColumn = this.column
	if c == '\n' {
		this.line++
		this.column = 1
	} else {
		if c == '\t' {
			this.column += 4 // TODO:: 4 OR 8
		} else {
			this.column++
		}
	}
	return
}

func (this *Lexer) unGetChar() {
	this.offset--
	this.line, this.column = this.lastLine, this.lastColumn
}

func (this *Lexer) unGetChar2(offset int) {
	this.offset -= offset
	this.column -= offset
}

func (this *Lexer) isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z')
}
func (this *Lexer) isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
func (this *Lexer) isOctal(c byte) bool {
	return '0' <= c && c <= '7'
}
func (this *Lexer) isHex(c byte) bool {

	return '0' <= c && c <= '9' ||
		('a' <= c && c <= 'f') ||
		('A' <= c && c <= 'F')
}

func (this *Lexer) hexByte2ByteValue(c byte) byte {
	if 'a' <= c && c <= 'f' {
		return c - 'a' + 10
	}
	if 'A' <= c && c <= 'F' {
		return c - 'A' + 10
	}
	return c - '0' //also valid for digit
}

func (this *Lexer) parseInt64(bs []byte) (int64, error) {
	base := int64(10)
	if bs[0] == '0' {
		base = 8
	}
	if len(bs) >= 2 &&
		bs[0] == '0' &&
		(bs[1] == 'X' || bs[1] == 'x') { // correct base to hex
		base = 16
		bs = bs[2:]
	}
	var result int64 = 0
	bit63is1 := false
	for _, v := range bs {
		result = result*base + int64(this.hexByte2ByteValue(v))
		if false == bit63is1 {
			if (result >> 63) != 0 {
				bit63is1 = true
				continue
			}
		}
		if bit63is1 {
			if (result >> 63) == 0 {
				bit63is1 = true
			}
			return result, fmt.Errorf("exceed max int64")
		}
	}
	return result, nil
}

func (this *Lexer) lexNumber(token *Token, c byte) (eof bool, err error) {
	integerPart := []byte{c}
	isHex := false
	isOctal := false
	if c == '0' { // enter when first char is '0'
		c, eof = this.getChar()
		if c == 'x' || c == 'X' {
			isHex = true
			integerPart = append(integerPart, 'X')
		} else {
			isOctal = true
			this.unGetChar()
		}
	}
	c, eof = this.getChar() //get next char
	for eof == false {
		ok := false
		if isHex {
			ok = this.isHex(c)
		} else if isOctal {
			if this.isDigit(c) == true && this.isOctal(c) == false { // integer but not octal
				err = fmt.Errorf("octal number cannot be '8' and '9'")
			}
			ok = this.isDigit(c)
		} else {
			ok = this.isDigit(c)
		}
		if ok {
			integerPart = append(integerPart, c)
			c, eof = this.getChar() // get next char
			continue
		} else { // something that I cannot handle
			this.unGetChar()
			break
		}
	}
	c, eof = this.getChar()
	floatPart := []byte{}
	haveFloatPart := false // float or double
	if c == '.' {          // float numbers
		haveFloatPart = true
		c, eof = this.getChar()
		for eof == false {
			if this.isDigit(c) {
				floatPart = append(floatPart, c)
				c, eof = this.getChar()
			} else {
				this.unGetChar()
				break
			}
		}
	} else {
		this.unGetChar()
	}
	if isHex && haveFloatPart {
		token.Type = TokenLiteralInt
		token.Data = int64(0)
		err = fmt.Errorf("mix up float and hex")
		return
	}

	isDouble := false
	isLong := false
	isShort := false
	isByte := false
	isFloat := false
	c, eof = this.getChar()
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
		this.unGetChar()
	}
	token.EndLine = this.line
	token.EndColumn = this.column
	if haveFloatPart {
		integerPart = append(integerPart, '.')
		floatValue, _ := strconv.ParseFloat(string(append(integerPart, floatPart...)), 64)
		if isDouble {
			token.Type = TokenLiteralDouble
			token.Data = floatValue
		} else {
			token.Type = TokenLiteralFloat
			token.Data = float32(floatValue)
		}
	} else {
		int64Value, e := this.parseInt64(integerPart)
		if err == nil && e != nil {
			err = e
		}
		if isDouble {
			token.Type = TokenLiteralDouble
			token.Data = float64(int64Value)
		} else if isFloat {
			token.Type = TokenLiteralFloat
			token.Data = float32(int64Value)
		} else if isLong {
			token.Type = TokenLiteralLong
			token.Data = int64(int64Value)
		} else if isByte {
			token.Type = TokenLiteralByte
			token.Data = int64(int64Value)
		} else if isShort {
			token.Type = TokenLiteralShort
			token.Data = int64(int64Value)
		} else {
			token.Type = TokenLiteralInt
			token.Data = int64(int64Value)
		}
	}
	return
}

func (this *Lexer) lexIdentifier(c byte) (token *Token, err error) {
	token = &Token{}
	token.StartLine = this.line
	token.StartColumn = this.column - 1 // c is read
	token.Offset = this.offset - 1      // c is read
	bs := []byte{c}
	c, eof := this.getChar()
	for eof == false {
		if this.isLetter(c) || c == '_' || this.isDigit(c) || c == '$' {
			bs = append(bs, c)
			c, eof = this.getChar()
		} else {
			this.unGetChar()
			break
		}
	}
	identifier := string(bs)
	if t, ok := keywordsMap[identifier]; ok {
		token.Type = t
		token.Description = identifier
		if token.Type == TokenElse {
			is := this.tryLexElseIf()
			if is {
				token.Type = TokenElseif
				token.Description = "else if"
			}
		}
	} else {
		token.Type = TokenIdentifier
		token.Data = identifier
		token.Description = "identifier_" + identifier
	}
	token.EndLine = this.line
	token.EndColumn = this.column
	return
}

func (this *Lexer) tryLexElseIf() (is bool) {
	c, eof := this.getChar()
	for c == ' ' || c == '\t' {
		c, eof = this.getChar()
	}
	if eof {
		return
	}
	if c != 'i' {
		this.unGetChar()
		return
	}
	c, eof = this.getChar()
	if c != 'f' {
		this.unGetChar()
		this.unGetChar2(1)
		return
	}
	c, eof = this.getChar()
	if c != ' ' && c != '\t' { // white list expect ' ' or '\t'
		this.unGetChar()
		this.unGetChar2(2) // un get 'i' and 'f'
		return
	}
	is = true
	return
}

func (this *Lexer) lexString(endChar byte) (token *Token, err error) {
	token = &Token{}
	token.StartLine = this.line
	token.StartColumn = this.column - 1
	token.Type = TokenLiteralString
	bs := []byte{}
	var c byte
	c, eof := this.getChar()
	for c != endChar && c != '\n' && eof == false {
		if c != '\\' {
			bs = append(bs, c)
			c, eof = this.getChar()
			continue
		}
		c, eof = this.getChar() // get next char
		if eof {
			err = fmt.Errorf("unexpected EOF")
			break
		}
		switch c {
		case 'a':
			bs = append(bs, '\a')
			c, eof = this.getChar()
		case 'b':
			bs = append(bs, '\b')
			c, eof = this.getChar()
		case 'f':
			bs = append(bs, '\f')
			c, eof = this.getChar()
		case 'n':
			bs = append(bs, '\n')
			c, eof = this.getChar()
		case 'r':
			bs = append(bs, '\r')
			c, eof = this.getChar()
		case 't':
			bs = append(bs, '\t')
			c, eof = this.getChar()
		case 'v':
			bs = append(bs, '\v')
			c, eof = this.getChar()
		case '\\':
			bs = append(bs, '\\')
			c, eof = this.getChar()
		case '\'':
			bs = append(bs, '\'')
			c, eof = this.getChar()
		case '"':
			bs = append(bs, '"')
			c, eof = this.getChar()
		case 'x':
			var c1, c2 byte
			c1, eof = this.getChar() //skip 'x'
			if eof {
				err = fmt.Errorf("unexpected EOF")
				continue
			}
			if false == this.isHex(c) {
				err = fmt.Errorf("unknown escape sequence")
				continue
			}
			b := this.hexByte2ByteValue(c1)
			c2, eof = this.getChar()
			if this.isHex(c2) {
				if t := b*16 + this.hexByte2ByteValue(c2); t <= 127 { // only support standard ascii
					b = t
				} else {
					this.unGetChar()
				}
			} else { //not hex
				this.unGetChar()
			}
			bs = append(bs, b)
			c, eof = this.getChar()
		case '0', '1', '2', '3', '4', '5', '7':
			// first char must be octal
			b := byte(0)
			for i := 0; i < 3; i++ {
				if eof {
					break
				}
				if this.isOctal(c) == false {
					this.unGetChar()
					break
				}
				if t := b*8 + this.hexByte2ByteValue(c); t > 127 { // only support standard ascii
					this.unGetChar()
					break
				} else {
					b = t
				}
				c, eof = this.getChar()
			}
			bs = append(bs, b)
			c, eof = this.getChar()
		case 'u', 'U':
			var r rune
			n := 4
			if c == 'U' {
				n = 8
			}
			for i := 0; i < n; i++ {
				c, eof = this.getChar()
				if eof {
					err = fmt.Errorf("unexcepted EOF")
					break
				}
				if this.isHex(c) == false {
					err = fmt.Errorf("not enough hex number for unicode, expect '%d' , but '%d'",
						n, i)
					this.unGetChar()
					break
				}
				r = (r << 4) | rune(this.hexByte2ByteValue(c))
			}
			bs = append(bs, []byte(string([]rune{r}))...)
			c, eof = this.getChar()
		default:
			err = fmt.Errorf("unknown escape sequence")
		}
	}
	token.EndLine = this.line
	token.EndColumn = this.column
	if c == '\n' {
		err = fmt.Errorf("string literal start new line")
	}
	token.Data = string(bs)
	token.Description = string(bs)
	return
}

func (this *Lexer) lexMultiLineComment() (string, error) {
	bs := []byte{}
redo:
	c, _ := this.getChar()
	var eof bool
	for c != '*' &&
		eof == false {
		c, eof = this.getChar()
		bs = append(bs, c)
	}
	if eof {
		return string(bs), fmt.Errorf("unexpect EOF")
	}
	c, eof = this.getChar()
	if eof {
		return string(bs), fmt.Errorf("unexpect EOF")
	}
	if eof || c == '/' {
		return string(bs[:len(bs)-1]), // slice out '*'
			nil
	}
	goto redo
}

/*
	one '.' is read
*/
func (this *Lexer) lexVArgs() (is bool) {
	c, _ := this.getChar()
	if c != '.' {
		this.unGetChar()
		return
	}
	// current '..'
	c, _ = this.getChar()
	if c != '.' {
		this.unGetChar()
		this.unGetChar2(1)
		return
	}
	is = true
	return
}

func (this *Lexer) isChar() bool {
	if this.offset+1 >= this.end {
		return false
	}
	if '\\' != this.bs[this.offset] {
		return false
	}

	if 'u' != this.bs[this.offset+1] && 'U' != this.bs[this.offset+1] {
		return false
	}
	return true
}

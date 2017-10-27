package lex

import (
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

var (
	lexer = lexmachine.NewLexer()
)

func init() {
	lexer.Add([]byte("function"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_FUNCTION
		return t, nil
	})
	lexer.Add([]byte("if"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_IF
		return t, nil
	})
	lexer.Add([]byte("elseif"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_ELSEIF
		return t, nil
	})
	lexer.Add([]byte("else"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_ELSE
		return t, nil
	})
	lexer.Add([]byte("for"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_FOR
		return t, nil
	})
	lexer.Add([]byte("return"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_RETURN
		return t, nil
	})
	lexer.Add([]byte("null"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_NULL
		return t, nil
	})
	lexer.Add([]byte("true"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_TRUE
		return t, nil
	})
	lexer.Add([]byte("false"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_FALSE
		return t, nil
	})
	lexer.Add([]byte("("), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_LP
		return t, nil
	})
	lexer.Add([]byte(")"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_RP
		return t, nil
	})
	lexer.Add([]byte("{"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_LC
		return t, nil
	})
	lexer.Add([]byte("}"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_RC
		return t, nil
	})
	lexer.Add([]byte("["), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_LB
		return t, nil
	})
	lexer.Add([]byte("]"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_LB
		return t, nil
	})
	lexer.Add([]byte(";"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_SEMICOLON
		return t, nil
	})
	lexer.Add([]byte(","), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_COMMA
		return t, nil
	})
	lexer.Add([]byte(","), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_COMMA
		return t, nil
	})
	lexer.Add([]byte("&&"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_LOGICAL_AND
		return t, nil
	})
	lexer.Add([]byte("||"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_LOGICAL_OR
		return t, nil
	})
	lexer.Add([]byte("="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_ASSIGN
		return t, nil
	})
	lexer.Add([]byte("=="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_EQUAL
		return t, nil
	})
	lexer.Add([]byte("!="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_NE
		return t, nil
	})
	lexer.Add([]byte(">"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_GT
		return t, nil
	})
	lexer.Add([]byte(">="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_GE
		return t, nil
	})
	lexer.Add([]byte("<"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_LT
		return t, nil
	})
	lexer.Add([]byte("<="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_LT
		return t, nil
	})
	lexer.Add([]byte("+"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_ADD
		return t, nil
	})
	lexer.Add([]byte("-"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_SUB
		return t, nil
	})
	lexer.Add([]byte("*"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_MUL
		return t, nil
	})
	lexer.Add([]byte{'/'}, func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_DIV
		return t, nil
	})
	lexer.Add([]byte("%"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_DIV
		return t, nil
	})
	lexer.Add([]byte("++"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_INCREMENT
		return t, nil
	})
	lexer.Add([]byte("--"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_DECREMENT
		return t, nil
	})
	lexer.Add([]byte("."), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_DOT
		return t, nil
	})
	lexer.Add([]byte("var"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_VAR
		return t, nil
	})
	lexer.Add([]byte("new"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_NEW
		return t, nil
	})
	lexer.Add([]byte(":"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_COLON
		return t, nil
	})
	lexer.Add([]byte("+="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_PLUS_ASSIGN
		return t, nil
	})
	lexer.Add([]byte("-="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_MINUS_ASSIGN
		return t, nil
	})
	lexer.Add([]byte("*="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_MUL_ASSIGN
		return t, nil
	})
	lexer.Add([]byte(`/=`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_DIV_ASSIGN
		return t, nil
	})
	lexer.Add([]byte("%="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_DIV_ASSIGN
		return t, nil
	})
	lexer.Add([]byte("!"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_NOT
		return t, nil
	})

	lexer.Add([]byte("switch"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_SWITCH
		return t, nil
	})
	lexer.Add([]byte("case"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_CASE
		return t, nil
	})
	lexer.Add([]byte("default"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_DEFAULT
		return t, nil
	})
	lexer.Add([]byte("\n"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_CRLF
		return t, nil
	})
	lexer.Add([]byte("( |\t)"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		return nil , nil
	})
	lexer.Add([]byte("//[^\n]*\n"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		return nil , nil
	})
	lexer.Add([]byte("/*[.\n]*/"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		return nil , nil
	})
	lexer.Add([]byte("package"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_PACKAGE
		return t, nil
	})
	lexer.Add([]byte("class"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_CLASS
		return t, nil
	})
	lexer.Add([]byte("int"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_INT
		return t, nil
	})
	lexer.Add([]byte("bool"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_BOOL
		return t, nil
	})
	lexer.Add([]byte("float"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_FLOAT
		return t, nil
	})
	lexer.Add([]byte("string"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_STRING
		return t, nil
	})

	lexer.Add([]byte("([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_)*"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_IDENTIFIER
		return t, nil
	})

	lexer.Add([]byte("(0?[0-9]*)"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_INT
		return t, nil
	})
	lexer.Add([]byte("([0x]?[0-9]+)"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_INT
		return t, nil
	})
	lexer.Add([]byte(`([0-9]+.[0-9]+)`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_LITERAL_FLOAT
		return t, nil
	})
	lexer.Add([]byte("(\".*\")"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Token = TOKEN_STRING
		return t, nil
	})
}



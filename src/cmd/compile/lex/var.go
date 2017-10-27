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
		t.Type = TOKEN_FUNCTION
		t.Desp = "function"
		return t, nil
	})
	lexer.Add([]byte("const"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_CONST
		t.Desp = "const"
		return t, nil
	})
	lexer.Add([]byte("if"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_IF
		t.Desp = "if"
		return t, nil
	})
	lexer.Add([]byte("elseif"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_ELSEIF
		t.Desp = "elseif"
		return t, nil
	})
	lexer.Add([]byte("else"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_ELSE
		t.Desp = "else"
		return t, nil
	})
	lexer.Add([]byte("for"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_FOR
		t.Desp = "for"
		return t, nil
	})
	lexer.Add([]byte("continue"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_CONTINUE
		t.Desp = "continue"
		return t, nil
	})
	lexer.Add([]byte("break"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_BREAK
		t.Desp = "break"
		return t, nil
	})
	lexer.Add([]byte("return"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_RETURN
		t.Desp = "return"
		return t, nil
	})
	lexer.Add([]byte("null"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_NULL
		t.Desp = "null"
		return t, nil
	})
	lexer.Add([]byte("true"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_TRUE
		t.Desp = "true"
		return t, nil
	})
	lexer.Add([]byte("false"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_FALSE
		return t, nil
	})
	lexer.Add([]byte(`([\(])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_LP
		t.Desp = "("
		return t, nil
	})
	lexer.Add([]byte(`([\)])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_RP
		t.Desp = ")"
		return t, nil
	})
	lexer.Add([]byte(`([\{])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_LC
		t.Desp = "{"
		return t, nil
	})
	lexer.Add([]byte(`([\}])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_RC
		t.Desp = "}"
		return t, nil
	})
	lexer.Add([]byte(`([\[])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_LB
		t.Desp = "["
		return t, nil
	})
	lexer.Add([]byte(`([\]])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_RB
		return t, nil
	})
	lexer.Add([]byte(";"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_SEMICOLON
		t.Desp = ";"
		return t, nil
	})
	lexer.Add([]byte(","), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_COMMA
		t.Desp = ","
		return t, nil
	})
	lexer.Add([]byte("&&"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_LOGICAL_AND
		t.Desp = "&&"
		return t, nil
	})
	lexer.Add([]byte("||"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_LOGICAL_OR
		t.Desp = "||"
		return t, nil
	})
	lexer.Add([]byte("="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_ASSIGN
		t.Desp = "="
		return t, nil
	})
	lexer.Add([]byte("=="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_EQUAL
		t.Desp = "=="
		return t, nil
	})
	lexer.Add([]byte("!="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_NE
		t.Desp = "!="
		return t, nil
	})
	lexer.Add([]byte(">"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_GT
		t.Desp = ">"
		return t, nil
	})
	lexer.Add([]byte(">="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_GE
		t.Desp = ">="
		return t, nil
	})
	lexer.Add([]byte("<"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_LT
		t.Desp = "<"
		return t, nil
	})
	lexer.Add([]byte("<="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_LT
		t.Desp = "<="
		return t, nil
	})
	lexer.Add([]byte(`([\+])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_ADD
		t.Desp = "+"
		return t, nil
	})
	lexer.Add([]byte(`([\-])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_SUB
		t.Desp = "-"
		return t, nil
	})
	lexer.Add([]byte(`([\*])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_MUL
		t.Desp = "*"
		return t, nil
	})
	lexer.Add([]byte(`([\/])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_DIV
		t.Desp = "/"
		return t, nil
	})
	lexer.Add([]byte("%"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_MOD
		t.Desp = "%"
		return t, nil
	})
	lexer.Add([]byte(`([\+\+])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_INCREMENT
		t.Desp = "++"
		return t, nil
	})
	lexer.Add([]byte(`(\-\-)`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_DECREMENT
		t.Desp = "--"
		return t, nil
	})
	lexer.Add([]byte(`([\.])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_DOT
		t.Desp = "."
		return t, nil
	})
	lexer.Add([]byte("var"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_VAR
		t.Desp = "var"
		return t, nil
	})
	lexer.Add([]byte("new"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_NEW
		t.Desp = "new"
		return t, nil
	})
	lexer.Add([]byte(":"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_COLON
		t.Desp = ":"
		return t, nil
	})
	lexer.Add([]byte(`([\+=])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_PLUS_ASSIGN
		t.Desp = "+="
		return t, nil
	})
	lexer.Add([]byte(`([\-=])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_MINUS_ASSIGN
		t.Desp = "-="
		return t, nil
	})
	lexer.Add([]byte(`([\*=])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_MUL_ASSIGN
		t.Desp = "*="
		return t, nil
	})
	lexer.Add([]byte(`([\/=])`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_DIV_ASSIGN
		t.Desp = `/=`
		return t, nil
	})
	lexer.Add([]byte("%="), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_DIV_ASSIGN
		t.Desp = "%="
		return t, nil
	})
	lexer.Add([]byte("!"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_NOT
		t.Desp = "!"
		return t, nil
	})

	lexer.Add([]byte("switch"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_SWITCH
		t.Desp = "switch"
		return t, nil
	})
	lexer.Add([]byte("case"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_CASE
		t.Desp = "case"
		return t, nil
	})
	lexer.Add([]byte("default"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_DEFAULT
		t.Desp = "default"
		return t, nil
	})
	lexer.Add([]byte("\n"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_CRLF
		t.Desp = "enter"
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
		t.Type = TOKEN_PACKAGE
		t.Desp = "package"
		return t, nil
	})
	lexer.Add([]byte("class"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_CLASS
		t.Desp = "class"
		return t, nil
	})
	lexer.Add([]byte("int"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_INT
		t.Desp = "int"
		return t, nil
	})
	lexer.Add([]byte("bool"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_BOOL
		t.Desp = "bool"
		return t, nil
	})
	lexer.Add([]byte("float"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_FLOAT
		t.Desp = "float"
		return t, nil
	})
	lexer.Add([]byte("string"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_STRING
		t.Desp = "string"
		return t, nil
	})

	lexer.Add([]byte("([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_)*"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_IDENTIFIER
		t.Desp = "identifer_" + string(match.Bytes)
		return t, nil
	})

	lexer.Add([]byte("(0?[0-9]*)"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_LITERAL_INT
		t.Desp = string(match.Bytes)
		return t, nil
	})
	lexer.Add([]byte("([0x]?[0-9]+)"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_LITERAL_INT
		t.Desp = string(match.Bytes)
		return t, nil
	})
	lexer.Add([]byte(`([0-9]+.[0-9]+)`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_LITERAL_FLOAT
		t.Desp = string(match.Bytes)
		return t, nil
	})
	lexer.Add([]byte("(\".*\")"), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		t := &Token{}
		t.Type = TOKEN_STRING
		t.Desp = string(match.Bytes)
		return t, nil
	})
}



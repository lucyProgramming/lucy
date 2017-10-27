package lex

import (
	"testing"
	"fmt"
)

var(
	code = `package test
	const a = 123;
		function(){
		}

	`
)

func Test_lex(t *testing.T) {
	s, err := lexer.Scanner([]byte(code))
	if err != nil {
		panic(err)
	}

	fmt.Println("\n\n\n\n\n\n\n########################")

	fmt.Println(code)

	fmt.Println("\n\n\n\n\n\n\n########################")
	for {
		t, err, eof := s.Next()
		if eof {
			break
		}
		if err != nil {
			fmt.Println("err:",err)
			continue
		}
		if t == nil {
			continue
		}
		if t.(*Token).Type == TOKEN_CRLF{
			fmt.Println()
		}else{
			fmt.Printf("%s ",t.(*Token).Desp)
		}
	}
	fmt.Println("\n\n\n\n\n\n\n")
}

package lex

import (
	"fmt"
	"strings"
	"testing"
)

var (
	code = ` 
		93434234+++
		\x60
			fdsfds

			fsdfsdf
			fdsfdsds
			fsqtttttttttttt
		\x60
	`
)

func Test_lex(t *testing.T) {
	s := New([]byte(code))
	fmt.Println("\n\n\n\n\n\n\n########################", len(strings.Split(code, "\n")))
	fmt.Println(code)
	fmt.Println("\n\n\n\n\n\n\n########################")
	for {
		token, eof, err := s.Next()
		if eof {
			break
		}
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		if token.Type == TOKEN_CRLF {
			fmt.Println()
			continue
		}
		fmt.Println("token ", "line:", token.StartLine, token.StartColumn, token.Data, token.Desp)
		fmt.Println()
	}

	fmt.Println("\n\n\n\n\n\n\n", s.line)
}

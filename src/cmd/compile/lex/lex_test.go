package lex

import (
	"testing"
)

func Test_lex(t *testing.T) {
	s, err := lexer.Scanner([]byte(`
		//

		package test

		function(){

		}


	`))
	if err != nil {
		panic(err)
	}
	ts := [...]int{
		TOKEN_PACKAGE,
		TOKEN_IDENTIFIER,
		TOKEN_FUNCTION,
		TOKEN_LP,
		TOKEN_RP,
		TOKEN_LC,
		TOKEN_CRLF,
		TOKEN_RC,
	}

	eof := false
	i := 0
	for eof == false {
		t, err, _ := s.Next()
		if err != nil {
			panic(err)
		}
		if t.(*Token).Type != ts[i] {
			panic(11)
		}
		i++
	}

}

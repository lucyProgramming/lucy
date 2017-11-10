package lex

import (
	"fmt"
	"testing"
)

var (
	code = `	package test
	
	import "github.om/xxx/ggg" as vvv;
	const a = 123;
	function(){
		a += b;
		var a int;
		return a;
		if(a > b){
			return true;
		}
		if(a >= b){
			return false;
		}
		if(a && b != 0){

			return null;
		}
		if(a | b != 0){
			return byte;
		}
		a + b;
		a / b;
		a % b;

	}
	class Person{
		int a;
		int b;
		int c;
		public fun Person(){
			for(int i = 0 ;i < 0x100;i++){
				System.out.println(i);
			}
			for(int i = 0 ;i < +0x100;i++){
				System.out.println(i);
			}
		}
		a + 1.00
		b + 0.000000000
		b + +1e5
	}
	{
		try{
			fsfds();
		}catch(e Exception){

		}finnaly{

		}
	}



	`
)

func Test_lex(t *testing.T) {
	s, err := Lexer.Scanner([]byte(code))
	if err != nil {
		panic(err)
	}

	fmt.Println("\n\n\n\n\n\n\n########################")
	fmt.Println(code)

	fmt.Println("\n\n\n\n\n\n\n########################")
	newline := false
	for {
		t, err, eof := s.Next()
		if eof {
			break
		}
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		if t == nil {
			continue
		}
		token := t.(*Token)
		if token.Type == TOKEN_CRLF {
			newline = true
			fmt.Println()
		} else {
			if newline == true {
				fmt.Printf("line:%d\t", token.Match.StartLine)
				newline = false
			}
			if token.Data != nil {
				fmt.Printf("%s(%v) ", token.Desp, token.Data)
			} else {
				fmt.Printf("%s ", token.Desp)
			}

		}
	}
	fmt.Println("\n\n\n\n\n\n\n")

}

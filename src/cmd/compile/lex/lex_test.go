package lex

import (
	"fmt"
	"testing"
)

var (
	code = `
	//11111
	package test
	const a = 123;
	function(){
		a += b;
		var a int;
		return a;
		if(a > b){
			return true
		}
		if(a >= b){
			return false
		}
		if(a && b != 0){

			return null
		}
		if(a | b != 0){
			return byte
		}
		a + b;
		a / b;
		a % b;

	}
	class Person{
		int a
		int b
		int c
		public Person(){
			for(int i = 0 ;i < 0x100;i++){
				System.out.println(i)
			}
			for(int i = 0 ;i < +0x100;i++){
				System.out.println(i)
			}
		}
		a + 1.00
		b + 0.000000000
		b + +1e5
	}
	;;;;100;

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
			fmt.Println("err:", err)
			continue
		}
		if t == nil {
			continue
		}
		if t.(*Token).Type == TOKEN_CRLF {
			fmt.Println()
		} else {
			fmt.Printf("%s ", t.(*Token).Desp)
		}
	}
	fmt.Println("\n\n\n\n\n\n\n")
}

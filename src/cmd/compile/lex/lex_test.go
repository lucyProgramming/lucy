package lex

import (
	"fmt"
	"strings"
	"testing"
)

var (
	code = `import "github.om/xxx/ggg" as vvv;
    const a = 123;
	fn aaa(){
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
	/*
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
	*/
	
	fsdfd fd
	
	fsdfs 
	1134
	;;
	if ske 
	
	a++
	
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

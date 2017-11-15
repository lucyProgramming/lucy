package parser

import (
	"fmt"
	"testing"

	"github.com/756445638/lucy/src/cmd/compile/ast"
)

var str = `

	package main

import "github.com/xxx/yyy" ;



fn add(a,b int)->(c int){
	return a + b ;
}
`

func Test_parse(t *testing.T) {
	nodes := []*ast.Node{}
	errs := Parse(&nodes, "demo.lucy", []byte(str))
	for _, v := range errs {
		fmt.Println(v)
	}
	for _, v := range nodes {
		fmt.Println(v)
	}
}

//var str = `

//	package main

//import "github.com/xxx/yyy" ;

//public enum Day {
//	Monday = 1,
//	TuesDay,
//}

//public const NAME="123"; // public global variable
//const age = 345;  // private global variable

//c := 100;  //private global variable

//{ //block execute at first
//	if c == 100{
//		skip; skip this blok,excute next one
//	}
//	a := 1;
//	if a == 1{
//		print("hello world")
//	}
//}

////function defination
//func (x int) Add(a int,b int){
//    return a +b + c;
//}

////class(key word) Class (class name) this(self name) [colon(:) Persion(fathers name)]
//class Man : Person{
//    public Person(){
//        var abc int;
//        f = fun(){
//            printf(abc);
//        }
//        f();
//    }
//}

//func main(){
//	for i:= 0 ;i < 100;i++{
//		print(i);
//	}
//}

//`

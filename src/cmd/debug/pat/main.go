package main

import (
	"encoding/json"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/parser"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	lucyFiles := make(map[string][]*ast.TopNode)
	var nodes []*ast.TopNode
	var errs []error
	if len(os.Args) == 1 {
		fmt.Println("no file to parse")
		os.Exit(2)
	}
	for _, v := range os.Args[1:] {
		if strings.HasSuffix(v, ".lucy") == false {
			fmt.Printf("'%s' not a lucy file\n", v)
			os.Exit(1)
		}
		bs, err := ioutil.ReadFile(v)
		if err != nil {
			fmt.Printf("read file '%s' failed,err:%v\n", v, err)
			os.Exit(1)
		}
		length := len(nodes)
		es := parser.Parse(&nodes, v, bs, false, 10)
		lucyFiles[v] = nodes[length:len(nodes)]
		errs = append(errs, es...)
	}
	for _, v := range errs {
		fmt.Println(v)
	}
	for _, v := range lucyFiles {
		bs, _ := json.Marshal(v)
		fmt.Println(string(bs))
	}
}

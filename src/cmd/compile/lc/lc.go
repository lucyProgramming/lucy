package lc

import (
	"encoding/json"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/parser"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

func Main(files []string) {
	if len(files) == 0 {
		fmt.Println("no file specfied")
		os.Exit(1)
	}
	if CompileFlags.OnlyImport == false {
		if CompileFlags.PackageName == "" {
			fmt.Println("package name not specfied")
			os.Exit(1)
		}
	}
	compiler.NerrsStopCompile = 10
	compiler.Errs = []error{}
	compiler.Files = files
	compiler.Init()
	compiler.compile()
}

type LucyCompile struct {
	Tops             []*ast.Node
	Files            []string
	Errs             []error
	NerrsStopCompile int
	lucyPath         []string
	ClassPath        []string
	Maker            jvm.MakeClass
}

func (lc *LucyCompile) shouldExit() {
	if len(lc.Errs) > lc.NerrsStopCompile {
		lc.exit()
	}
}

func (lc *LucyCompile) exit() {
	code := 0
	if len(lc.Errs) > 0 {
		code = 2
	}
	for _, v := range lc.Errs {
		fmt.Fprintln(os.Stderr, v)
	}
	os.Exit(code)
}

func (lc *LucyCompile) Init() {
	path := os.Getenv("CLASSPATH")
	if runtime.GOOS == "windows" {
		lc.ClassPath = strings.Split(path, ";") // windows style
	} else {
		lc.ClassPath = strings.Split(path, ":")
	}
	path = os.Getenv("LUCYPATH")
	if runtime.GOOS == "windows" {
		lc.lucyPath = strings.Split(path, ";") // windows style
	} else {
		lc.lucyPath = strings.Split(path, ":")
	}
}

func (lc *LucyCompile) parseImports() {
	if len(lc.Errs) > 0 {
		lc.exit()
	}
	is := make([]string, len(lc.Tops))
	for k, v := range lc.Tops {
		is[k] = v.Data.(*ast.Import).Resource
	}
	bs, _ := json.Marshal(is)
	fmt.Println(string(bs))
	//nameLoader := &RealNameLoader{}
	//packageNames := make(map[string]struct{})
	//for _, v := range is {
	//	p, c, err := nameLoader.LoadName(v)
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//	if p != nil {
	//		packageNames[p.Name] = struct{}{}
	//	}
	//	if pp, ok := c.(*ast.Package); ok && pp != nil {
	//		packageNames[pp.Name] = struct{}{}
	//	}
	//}
	//is = make([]string, len(packageNames))
	//i := 0
	//for name, _ := range packageNames {
	//	is[i] = name
	//	i++
	//}
	//bs, _ := json.Marshal(is)
	//fmt.Println(string(bs))
}

func (lc *LucyCompile) compile() {
	for _, v := range lc.Files {
		bs, err := ioutil.ReadFile(v)
		if err != nil {
			lc.Errs = append(lc.Errs, err)
			continue
		}
		lc.Errs = append(lc.Errs, parser.Parse(&lc.Tops, v, bs,
			CompileFlags.OnlyImport, lc.NerrsStopCompile)...)
		lc.shouldExit()
	}
	// parse import only
	if CompileFlags.OnlyImport {
		lc.parseImports()
		return
	}
	c := ast.ConvertTops2Package{}
	p, rs, errs := c.ConvertTops2Package(lc.Tops)
	p.Name = CompileFlags.PackageName
	lc.Errs = append(lc.Errs, errs...)
	for _, v := range rs {
		lc.Errs = append(lc.Errs, v.Error())
	}
	ast.PackageBeenCompile = p
	lc.shouldExit()
	lc.Errs = append(lc.Errs, p.TypeCheck()...)
	if len(lc.Errs) > 0 {
		lc.exit()
	}
	lc.Maker.Make(p)
	lc.exit()
}

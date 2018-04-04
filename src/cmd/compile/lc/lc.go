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
	code := len(lc.Errs)
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

func (lc *LucyCompile) compile() {

	for _, v := range lc.Files {
		bs, err := ioutil.ReadFile(v)
		if err != nil {
			lc.Errs = append(lc.Errs, err)
			continue
		}
		lc.Errs = append(lc.Errs, parser.Parse(&lc.Tops, v, bs, CompileFlags.OnlyImport, lc.NerrsStopCompile)...)
		lc.shouldExit()
	}

	// parse import only
	if CompileFlags.OnlyImport {
		if len(lc.Errs) > 0 {
			lc.exit()
		}
		is := make([]string, len(lc.Tops))
		for k, v := range lc.Tops {
			is[k] = v.Data.(*ast.Import).Name
		}
		bs, _ := json.Marshal(is)
		fmt.Println(string(bs))
		return
	}

	c := ast.ConvertTops2Package{}
	p, rs, errs := c.ConvertTops2Package(lc.Tops)
	p.FullName = CompileFlags.PackageName
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

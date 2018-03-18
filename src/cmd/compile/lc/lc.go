package lc

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm"
	"github.com/756445638/lucy/src/cmd/compile/parser"
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
	if CompileFlags.PackageName == "" {
		fmt.Println("package name not specfied")
		os.Exit(1)
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
		code = 1
	}
	for _, v := range lc.Errs {
		fmt.Println(v)
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
	//	defer func() {
	//		if err := recover(); err != nil {
	//			debug.PrintStack()
	//		}
	//	}()

	for _, v := range lc.Files {
		bs, err := ioutil.ReadFile(v)
		if err != nil {
			lc.Errs = append(lc.Errs, err)
			continue
		}
		lc.Errs = append(lc.Errs, parser.Parse(&lc.Tops, v, bs, CompileFlags.OnlyImport, lc.NerrsStopCompile)...)
		lc.shouldExit()
	}
	c := ast.ConvertTops2Package{}
	p, rs, errs := c.ConvertTops2Package(lc.Tops)
	p.FullName = CompileFlags.PackageName
	lc.Errs = append(lc.Errs, errs...)
	for _, v := range rs {
		lc.Errs = append(lc.Errs, v.Error())
	}
	lc.shouldExit()
	lc.Errs = append(lc.Errs, p.TypeCheck()...)
	if len(lc.Errs) > 0 {
		lc.exit()
	}
	lc.Maker.Make(p)
	lc.exit()
}

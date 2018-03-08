package lc

import (
	"fmt"
	"io/ioutil"
	"os"
	//	"runtime/debug"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm"
	"github.com/756445638/lucy/src/cmd/compile/parser"
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
	compiler.compile()

}

type LucyCompile struct {
	Tops             []*ast.Node
	Files            []string
	Errs             []error
	NerrsStopCompile int
	lucyPath         []string
	Maker            jvm.MakeClass
}

func (l *LucyCompile) shouldExit() {
	if len(l.Errs) > l.NerrsStopCompile {
		l.exit()
	}
}

func (l *LucyCompile) exit() {
	code := 0
	if len(l.Errs) > 0 {
		code = 1
	}
	for _, v := range l.Errs {
		fmt.Println(v)
	}
	os.Exit(code)
}

func (l *LucyCompile) Init() {
	path := os.Getenv("CLASSPATH")
	l.lucyPath = strings.Split(path, ":")
	if len(l.lucyPath) == 0 {
		fmt.Println("env variable LUCYPATH is not set")
	}
}

func (l *LucyCompile) compile() {
	//	defer func() {
	//		if err := recover(); err != nil {
	//			debug.PrintStack()
	//		}
	//	}()
	l.Init()
	for _, v := range l.Files {
		bs, err := ioutil.ReadFile(v)
		if err != nil {
			l.Errs = append(l.Errs, err)
			continue
		}
		l.Errs = append(l.Errs, parser.Parse(&l.Tops, v, bs, CompileFlags.OnlyImport, l.NerrsStopCompile)...)
		l.shouldExit()
	}
	c := ast.ConvertTops2Package{}
	p, rs, errs := c.ConvertTops2Package(l.Tops)
	p.FullName = CompileFlags.PackageName
	l.Errs = append(l.Errs, errs...)
	for _, v := range rs {
		l.Errs = append(l.Errs, v.Error())
	}
	l.shouldExit()
	l.Errs = append(l.Errs, p.TypeCheck()...)
	if len(l.Errs) > 0 {
		l.exit()
	}
	l.Maker.Make(p)
	l.exit()
}

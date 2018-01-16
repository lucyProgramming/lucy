package lc

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm"
	"github.com/756445638/lucy/src/cmd/compile/parser"
)

func Main(files []string) {
	l.NerrsStopCompile = 10
	l.Nerrs = []error{}
	l.Files = files
	l.compile()
}

type LucyCompile struct {
	Tops             []*ast.Node
	Files            []string
	Nerrs            []error
	NerrsStopCompile int
	lucyPath         []string
	Maker            jvm.MakeClass
}

func (l *LucyCompile) shouldExit() {
	if len(l.Nerrs) > l.NerrsStopCompile {
		l.exit()
	}

}

func (l *LucyCompile) exit() {
	for _, v := range l.Nerrs {
		fmt.Println(v)
	}
}

func (l *LucyCompile) Init() {
	path := os.Getenv("LUCYPATH")
	l.lucyPath = strings.Split(path, ":")
	if len(l.lucyPath) == 0 {
		fmt.Println("env variable LUCYPATH is not set")
	}
}

func (l *LucyCompile) compile() {
	l.Init()
	for _, v := range l.Files {
		bs, err := ioutil.ReadFile(v)
		if err != nil {
			l.Nerrs = append(l.Nerrs, err)
			continue
		}
		l.Nerrs = append(l.Nerrs, parser.Parse(&l.Tops, v, bs, CompileFlags.OnlyImport, l.NerrsStopCompile)...)
		l.shouldExit()
	}
	c := ast.ConvertTops2Package{}
	p, rs, errs := c.ConvertTops2Package(l.Tops)
	l.Nerrs = append(l.Nerrs, errs...)
	for _, v := range rs {
		l.Nerrs = append(l.Nerrs, v.Error())
	}
	l.shouldExit()
	l.Nerrs = append(l.Nerrs, p.TypeCheck()...)
	l.shouldExit()
	l.Maker.Make(p)
	l.exit()
}

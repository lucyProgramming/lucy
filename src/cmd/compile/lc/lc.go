package lc

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
	"github.com/756445638/lucy/src/cmd/compile/parser"
)

func Main(files []string) {
	go cg.Prinf()
	l.NerrsStopCompile = 10
	l.Nerrs = []error{}
	l.compile()
}

type LucyCompile struct {
	Files            []string
	Nerrs            []error
	NerrsStopCompile int
	lucyPath         []string
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
		l.Nerrs = append(l.Nerrs, parser.Parse(&Tops, v, bs, CompileFlags.OnlyImport)...)

		if len(l.Nerrs) > 10 {
			l.exit()
		}
	}

}

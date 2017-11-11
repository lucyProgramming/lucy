package lc

import (
	"github.com/756445638/lucy/src/cmd/compile/parser"
	"io/ioutil"
)

func Main(files []string) {
	l := &LucyCompile{
		Files:            files,
		NerrsStopCompile: 10,
		Nerrs:            []error{},
	}
	l.compile()
}

type LucyCompile struct {
	Files            []string
	Nerrs            []error
	NerrsStopCompile int
}

func (l *LucyCompile) compile() {
	for _, v := range l.Files {
		bs, err := ioutil.ReadFile(v)
		if err != nil {
			l.Nerrs = append(l.Nerrs, err)
			continue
		}
		parser.Parse(&Tops, v, bs)
	}

}

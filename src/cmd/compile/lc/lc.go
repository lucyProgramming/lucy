package lc

import (
	"io/ioutil"

	"github.com/756445638/lucy/src/cmd/compile/yacc"
)

/*
	files in on package
*/
func Main(files []string) {
	l := &LucyCompile{
		Files:            files,
		NerrsStopCompile: 10,
		Nerrs:            []error{},
	}
}

type LucyCompile struct {
	Files            []string
	Nerrs            []error
	NerrsStopCompile int
}

func (l *LucyCompile) compile() {
	for k, v := range l.Files {
		bs, err := ioutil.ReadFile(v)
		if err != nil {
			l.Nerrs = append(l.Nerrs, err)
			continue
		}
	}

}

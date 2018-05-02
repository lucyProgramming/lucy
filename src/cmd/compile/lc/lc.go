package lc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"gitee.com/yuyang-fine/lucy/src/cmd/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/parser"
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
	lucyPaths        []string
	ClassPaths       []string
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
	// sort error
	es := SortErrs(lc.Errs)
	sort.Sort(es)
	for _, v := range es {
		fmt.Fprintln(os.Stderr, v)
	}
	os.Exit(code)
}

func (lc *LucyCompile) Init() {
	lc.ClassPaths = common.GetClassPaths()
	var err error
	lc.lucyPaths, err = common.GetLucyPaths()
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
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
	ast.PackageBeenCompile.Name = CompileFlags.PackageName
	rs, errs := c.ConvertTops2Package(lc.Tops)
	lc.Errs = append(lc.Errs, errs...)
	for _, v := range rs {
		lc.Errs = append(lc.Errs, v.Error())
	}
	lc.shouldExit()
	lc.Errs = append(lc.Errs, ast.PackageBeenCompile.TypeCheck()...)
	if len(lc.Errs) > 0 {
		lc.exit()
	}
	lc.Maker.Make(&ast.PackageBeenCompile)
	lc.exit()
}

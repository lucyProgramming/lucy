package lc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"gitee.com/yuyang-fine/lucy/src/cmd/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	compileCommon "gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm"
	//	"gitee.com/yuyang-fine/lucy/src/cmd/compile/optimizer"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/parser"
)

func Main(files []string) {
	if len(files) == 0 {
		fmt.Println("no file specfied")
		os.Exit(1)
	}
	if compileCommon.CompileFlags.OnlyImport == false {
		if compileCommon.CompileFlags.PackageName == "" {
			fmt.Println("package name not specfied")
			os.Exit(1)
		}
	}
	compiler.NErrsStopCompile = 10
	compiler.Errs = []error{}
	compiler.Files = files
	compiler.Init()
	compiler.compile()
}

type LucyCompile struct {
	Tops             []*ast.Node
	Files            []string
	Errs             []error
	NErrsStopCompile int
	lucyPaths        []string
	ClassPaths       []string
	Maker            jvm.MakeClass
}

func (lc *LucyCompile) shouldExit() {
	if len(lc.Errs) > lc.NErrsStopCompile {
		lc.exit()
	}
}

func (lc *LucyCompile) exit() {
	code := 0
	if len(lc.Errs) > 0 {
		code = 2
	}
	//template function maybe parse twice,maybe same error
	errsM := make(map[string]struct{})
	es := []error{}
	for _, v := range lc.Errs {
		if _, ok := errsM[v.Error()]; ok {
			continue
		}
		es = append(es, v)
		errsM[v.Error()] = struct{}{}
	}
	lc.Errs = es
	// sort errors
	es = SortErrors(lc.Errs)
	sort.Sort(SortErrors(es))
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
			compileCommon.CompileFlags.OnlyImport, lc.NErrsStopCompile)...)
		lc.shouldExit()
	}
	// parse import only
	if compileCommon.CompileFlags.OnlyImport {
		lc.parseImports()
		return
	}
	c := ast.ConvertTops2Package{}
	ast.PackageBeenCompile.Name = compileCommon.CompileFlags.PackageName
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
	//optimizer.Optimize(&ast.PackageBeenCompile)
	lc.Maker.Make(&ast.PackageBeenCompile)
	lc.exit()
}

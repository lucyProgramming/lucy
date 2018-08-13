package lc

import (
	"encoding/json"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	compileCommon "gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/parser"
	"io/ioutil"
	"os"
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

type Compiler struct {
	Tops             []*ast.Top
	Files            []string
	Errs             []error
	NErrsStopCompile int
	lucyPaths        []string
	ClassPaths       []string
	Maker            jvm.BuildPackage
}

func (compiler *Compiler) shouldExit() {
	if len(compiler.Errs) > compiler.NErrsStopCompile {
		compiler.exit()
	}
}

func (compiler *Compiler) exit() {
	code := 0
	if len(compiler.Errs) > 0 {
		code = 2
	}
	for _, v := range compiler.Errs {
		fmt.Fprintln(os.Stderr, v)
	}
	os.Exit(code)
}

func (compiler *Compiler) Init() {
	compiler.ClassPaths = common.GetClassPaths()
	var err error
	compiler.lucyPaths, err = common.GetLucyPaths()
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}

func (compiler *Compiler) dumpImports() {
	if len(compiler.Errs) > 0 {
		compiler.exit()
	}
	is := make([]string, len(compiler.Tops))
	for k, v := range compiler.Tops {
		is[k] = v.Data.(*ast.Import).Import
	}
	bs, _ := json.Marshal(is)
	fmt.Println(string(bs))
}

func (compiler *Compiler) compile() {
	for _, v := range compiler.Files {
		bs, err := ioutil.ReadFile(v)
		if err != nil {
			compiler.Errs = append(compiler.Errs, err)
			continue
		}

		//UTF-16 (BE)
		if len(bs) >= 2 &&
			bs[0] == 0xfe &&
			bs[1] == 0xff {
			fmt.Printf("file:%s looks like UTF-16(BE) file\n", v)
			os.Exit(2)
		}
		//UTF-16 (LE)
		if len(bs) >= 2 &&
			bs[0] == 0xff &&
			bs[1] == 0xfe {
			fmt.Printf("file:%s looks like UTF-16(LE) file\n", v)
			os.Exit(2)
		}
		//UTF-32 (LE)
		if len(bs) >= 4 &&
			bs[0] == 0x0 &&
			bs[1] == 0x0 &&
			bs[2] == 0xfe &&
			bs[3] == 0xff {
			fmt.Printf("file:%s looks like UTF-32(LE) file\n", v)
			os.Exit(2)
		}
		//UTF-32 (BE)
		if len(bs) >= 4 &&
			bs[0] == 0xff &&
			bs[1] == 0xfe &&
			bs[2] == 0x0 &&
			bs[3] == 0x0 {
			fmt.Printf("file:%s looks like UTF-32(BE) file\n", v)
			os.Exit(2)
		}

		if len(bs) >= 3 &&
			bs[0] == 0xef &&
			bs[1] == 0xbb &&
			bs[2] == 0xbf {
			// utf8 bom
			bs = bs[3:]
		}
		compiler.Errs = append(compiler.Errs, parser.Parse(&compiler.Tops, v, bs,
			compileCommon.CompileFlags.OnlyImport, compiler.NErrsStopCompile)...)
		compiler.shouldExit()
	}
	// parse import only
	if compileCommon.CompileFlags.OnlyImport {
		compiler.dumpImports()
		return
	}
	c := ast.ConvertTops2Package{}
	ast.PackageBeenCompile.Name = compileCommon.CompileFlags.PackageName
	rs, errs := c.ConvertTops2Package(compiler.Tops)
	compiler.Errs = append(compiler.Errs, errs...)
	for _, v := range rs {
		compiler.Errs = append(compiler.Errs, v.Error())
	}
	compiler.shouldExit()
	compiler.Errs = append(compiler.Errs, ast.PackageBeenCompile.TypeCheck()...)
	if len(compiler.Errs) > 0 {
		compiler.exit()
	}
	//optimizer.Optimize(&ast.PackageBeenCompile)
	compiler.Maker.Make(&ast.PackageBeenCompile)
	compiler.exit()
}

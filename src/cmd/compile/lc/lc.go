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

type Compiler struct {
	Tops              []*ast.TopNode
	Files             []string
	Errs              []error
	NErrs2StopCompile int
	lucyPaths         []string
	ClassPaths        []string
	Maker             jvm.BuildPackage
}

func Main(files []string) {
	if len(files) == 0 {
		fmt.Println("no file specfied")
		os.Exit(1)
	}
	compiler.NErrs2StopCompile = 10
	compiler.Errs = []error{}
	compiler.Files = files
	compiler.Init()
	compiler.compile()
}

func (this *Compiler) shouldExit() {
	if len(this.Errs) > this.NErrs2StopCompile {
		this.exit()
	}
}

func (this *Compiler) exit() {
	code := 0
	if len(this.Errs) > 0 {
		code = 2
	}
	for _, v := range this.Errs {
		fmt.Fprintln(os.Stderr, v)
	}
	os.Exit(code)
}

func (this *Compiler) Init() {
	this.ClassPaths = common.GetClassPaths()
	var err error
	this.lucyPaths, err = common.GetLucyPaths()
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}

func (this *Compiler) dumpImports() {
	if len(this.Errs) > 0 {
		this.exit()
	}
	is := make([]string, len(this.Tops))
	for k, v := range this.Tops {
		is[k] = v.Node.(*ast.Import).Import
	}
	bs, _ := json.Marshal(is)
	fmt.Println(string(bs))
}

func (this *Compiler) compile() {
	fileNodes := make(map[string][]*ast.TopNode)

	for _, v := range this.Files {

		bs, err := ioutil.ReadFile(v)
		if err != nil {
			this.Errs = append(this.Errs, err)
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
		length := len(this.Tops)
		this.Errs = append(this.Errs, parser.Parse(&this.Tops, v, bs,
			compileCommon.CompileFlags.OnlyImport, this.NErrs2StopCompile)...)
		fileNodes[v] = this.Tops[length:len(this.Tops)]
		this.shouldExit()
	}
	// parse import only
	if compileCommon.CompileFlags.OnlyImport {
		this.dumpImports()
		return
	}

	if compileCommon.CompileFlags.PackageName == "" {
		fmt.Println("package name not specfied")
		os.Exit(1)
	}
	c := ast.ConvertTops2Package{}
	ast.PackageBeenCompile.Name = compileCommon.CompileFlags.PackageName
	rs, errs := c.ConvertTops2Package(this.Tops)
	this.Errs = append(this.Errs, errs...)
	for _, v := range rs {
		this.Errs = append(this.Errs, v.Error())
	}
	this.shouldExit()

	this.Errs = append(this.Errs, ast.PackageBeenCompile.TypeCheck()...)
	if len(this.Errs) > 0 {
		this.exit()
	}
	//optimizer.Optimize(&ast.PackageBeenCompile)
	this.Maker.Make(&ast.PackageBeenCompile)
	this.exit()
}

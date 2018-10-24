package main

import (
	"flag"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/lc"
)

func main() {
	flag.BoolVar(&common.CompileFlags.OnlyImport,
		"only-import", false, "only parse import package")
	flag.BoolVar(&common.CompileFlags.DisableCheckUnUse,
		"disable-check-unuse", false, "disable check un use")
	flag.BoolVar(&common.CompileFlags.DumpParseFile,
		"dump-parse-file", false, "dump parse file")
	flag.StringVar(&common.CompileFlags.PackageName,
		"package-name", "", "package name")
	flag.IntVar(&common.CompileFlags.JvmMajorVersion,
		"jvm-major-version", 52, "jvm major version")
	flag.IntVar(&common.CompileFlags.JvmMinorVersion,
		"jvm-minor-version", 0, "jvm  version")
	flag.Parse()
	lc.Main(flag.Args())
}

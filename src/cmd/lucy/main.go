package main

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/lucy/install_lucy_array"
	"gitee.com/yuyang-fine/lucy/src/cmd/lucy/run"
	"os"
	"runtime"
)

func printUsage() {
	msg := `lucy is a new programing language build on jvm
	version                print version
	build                  build package and don't run					 
	run                    run a lucy package
	clean                  clean compiled files
	pack                   make jar
	test                   test a package`
	fmt.Println(msg)
}

func main() {

	if len(os.Args) == 1 {
		printUsage()
		os.Exit(0)
	}
	switch os.Args[1] {
	case "version":
		fmt.Printf("lucy-%v@(%s/%s)\n", common.VERSION, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	case "build":
		(&run.Run{}).RunCommand("run", append([]string{"-just-build"}, os.Args[2:]...))
	case "run":
		(&run.Run{}).RunCommand(os.Args[1], os.Args[2:])
	case "clean":
		args := []string{"lucy/cmd/langtools/clean"}
		args = append(args, os.Args[2:]...)
		(&run.Run{}).RunCommand("run", args)
	case "test":
		args := []string{"lucy/cmd/langtools/test"}
		args = append(args, os.Args[2:]...)
		(&run.Run{}).RunCommand("run", args)
	case "install_lucy_array":
		(&install_lucy_array.InstallLucyArray{}).RunCommand("install_lucy_array", nil)
	case "pack":
		args := []string{"lucy/cmd/langtools/pack"}
		args = append(args, os.Args[2:]...)
		(&run.Run{}).RunCommand("run", args)
	default:
		printUsage()
		os.Exit(1)
	}
}

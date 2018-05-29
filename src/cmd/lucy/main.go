package main

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/common"

	"gitee.com/yuyang-fine/lucy/src/cmd/lucy/run"
	"os"
	"runtime"
)

func printUsage() {
	msg := `lucy is a new programing language build on jvm
	run                    run a lucy package
	version                print version
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
	case "run":
		(&run.Run{}).RunCommand(os.Args[1], os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}

}

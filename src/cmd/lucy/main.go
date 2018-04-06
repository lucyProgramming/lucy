package main

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/common"
	"gitee.com/yuyang-fine/lucy/src/cmd/lucy/clean"
	"gitee.com/yuyang-fine/lucy/src/cmd/lucy/command"
	"gitee.com/yuyang-fine/lucy/src/cmd/lucy/install_lucy_array"
	"gitee.com/yuyang-fine/lucy/src/cmd/lucy/run"
	"os"
	"runtime"
)

var (
	commands = make(map[string]command.RunCommand)
)

func init() {
	commands["run"] = &run.Run{}
	commands["install_lucy_array"] = &install_lucy_array.InstallLucyArray{}
	commands["clean"] = &clean.Clean{}
}

func printUsage() {
	msg := `lucy is a new programing language build on jvm
	run                    run a lucy package
	version                print version
	clean                  clean compiled files
	pack                   make jar`
	fmt.Println(msg)
}

func main() {
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(0)
	}
	if os.Args[1] == "version" {
		fmt.Printf("lucy-%v@(%s/%s)\n", common.VERSION, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}
	c, ok := commands[os.Args[1]]
	if ok == false {
		printUsage()
		os.Exit(0)
	}

	c.RunCommand(os.Args[1], os.Args[2:])
}

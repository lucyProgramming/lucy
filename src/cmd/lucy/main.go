package main

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/lucy/command"
	"gitee.com/yuyang-fine/lucy/src/cmd/lucy/install_lucy_array"
	"gitee.com/yuyang-fine/lucy/src/cmd/lucy/run"
	"os"
)

var (
	commands = make(map[string]command.RunCommand)
)

func init() {
	commands["run"] = &run.Run{}
	commands["install_lucy_array"] = &install_lucy_array.InstallLucyArray{}
}

func printUsage() {
	msg := "lucy is a new programing language base on jvm\n"
	msg += "\t run		run a lucy package"
	fmt.Println(msg)
}

func main() {
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(0)
	}
	c, ok := commands[os.Args[1]]
	if ok == false {
		printUsage()
		os.Exit(0)
	}
	c.RunCommand(os.Args[1], os.Args[2:])
}

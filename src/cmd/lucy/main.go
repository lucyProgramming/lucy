package main

import (
	"github.com/756445638/lucy/src/cmd/lucy/command"
	"github.com/756445638/lucy/src/cmd/lucy/install_lucy_array"
	"github.com/756445638/lucy/src/cmd/lucy/run"
	"os"
	"fmt"
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

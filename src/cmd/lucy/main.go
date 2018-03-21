package main

import (
	"github.com/756445638/lucy/src/cmd/lucy/command"
	"github.com/756445638/lucy/src/cmd/lucy/run"
	"os"
)

var (
	commands = make(map[string]command.RunCommand)
)

func init() {
	commands["run"] = &run.Run{}
}

func printUsage() {

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
	c.RunCommand(os.Args[2:])
}

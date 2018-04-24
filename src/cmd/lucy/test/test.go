package test

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/common"
	"os"
)

type Test struct {
	LucyRoot    string
	LucyPaths   []string
	In          string
	packageName string
}

func (test *Test) Help(command string) {

}

func (test *Test) RunCommand(command string, args []string) {
	if len(args) == 0 {
		test.Help(command)
		os.Exit(1)
	}
	test.packageName = args[0]
	var err error
	test.LucyRoot, err = common.GetLucyRoot()
	if err != nil {
		test.Help(command)
		os.Exit(1)
	}
	test.LucyPaths, err = common.GetLucyPaths()
	if err != nil {
		test.Help(command)
		os.Exit(2)
	}
	dirs, err := common.FindLucyPackageDirectory(test.packageName, test.LucyPaths)
	if err != nil {
		test.Help(command)
		os.Exit(3)
	}
	if len(dirs) == 0 {
		fmt.Println("package '%s' not found", test.packageName)
		os.Exit(4)
	}
}

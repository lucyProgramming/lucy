package test

import (
	"bufio"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/common"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type Test struct {
	LucyRoot      string
	LucyPaths     []string
	In            string
	packageName   string
	lucyCommandAt string
	LiesIn        string
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
	{
		found := false
		for _, v := range test.LucyPaths {
			if v == test.LucyRoot {
				found = true
				break
			}
		}
		if found == false {
			test.LucyPaths = append(test.LucyPaths, test.LucyRoot)
		}
	}
	dirs := common.FindLucyPackageDirectory(test.packageName, test.LucyPaths)
	if len(dirs) == 0 {
		fmt.Printf("package named '%s' not found\n", test.packageName)
		os.Exit(4)
	}
	if len(dirs) > 1 {
		fmt.Printf("not 1 package named '%s'\n", test.packageName)
		os.Exit(5)
	}
	//check is tests dir is exists
	f, _ := os.Stat(filepath.Join(dirs[0], "test"))
	if f == nil {
		fmt.Printf("sub directory named 'test' not exists in '%s'\n", dirs[0])
		os.Exit(6)
	}
	if f.IsDir() == false {
		fmt.Printf("'test' is not a directory in '%s'\n", dirs[0])
		os.Exit(7)
	}
	test.lucyCommandAt = filepath.Join(test.LucyRoot, "bin", "lucy")
	test.testDir(filepath.Join(dirs[0], common.DIR_FOR_LUCY_SOURCE_FILES), test.packageName+"/"+"test")
}

func (test *Test) testDir(dir string, prefix string) {
	path := filepath.Join(dir, prefix)
	if true == common.SourceFileExist(path) {
		// test this package
		cmd := exec.Command(test.lucyCommandAt, "run", "-forceReBuild", prefix)
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@", test.lucyCommandAt, "run", "-forceReBuild", prefix)
		fmt.Printf(test_package_pre_msg, prefix, path) // output debug infos
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = path
		err := cmd.Start()
		if err != nil {
			fmt.Printf("fatal error start run process failed,err:%v\n", err)
			os.Exit(8)
		}
		err = cmd.Wait()
		if err != nil {
			fmt.Printf("test package '%s'  in '%s' failed\n", prefix, path)
			fmt.Printf("are you wish to continue?[y/n]\n")
			buf := bufio.NewReader(os.Stdin)
			line, _, _ := buf.ReadLine()
			if line != nil && len(line) > 0 && (line[0] == 'n' || line[0] == 'N') {
				os.Exit(0) // exit normally
			}
			fmt.Println("fail")
		} else {
			fmt.Println("ok")
		}
	}
	// test sub directory recursively
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("read sub directoies failed,err:%v\n", err)
		os.Exit(9)
	}
	for _, f := range fis {
		if f.IsDir() == false {
			continue
		}
		test.testDir(dir, prefix+"/"+f.Name())
	}
}

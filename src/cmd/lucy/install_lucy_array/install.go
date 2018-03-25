package install_lucy_array

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type InstallLucyArray struct {
}

type InstallType struct {
	classname    string
	typename     string
	defaultValue string
}

var (
	installs = []*InstallType{}
)

func init() {
	installs = append(installs, &InstallType{
		classname: "ArrayBool",
		typename:  "boolean",
	})
	installs = append(installs, &InstallType{
		classname: "ArrayByte",
		typename:  "byte",
	})
	installs = append(installs, &InstallType{
		classname: "ArrayShort",
		typename:  "short",
	})
	installs = append(installs, &InstallType{
		classname: "ArrayInt",
		typename:  "int",
	})
	installs = append(installs, &InstallType{
		classname: "ArrayLong",
		typename:  "long",
	})
	installs = append(installs, &InstallType{
		classname: "ArrayFloat",
		typename:  "float",
	})
	installs = append(installs, &InstallType{
		classname: "ArrayDouble",
		typename:  "double",
	})
	installs = append(installs, &InstallType{
		classname: "ArrayObject",
		typename:  "Object",
	})
	d := `
		for(int i =0 ;i < this.end;i ++){
			this.elements[i] = "";
		}
	`
	installs = append(installs, &InstallType{
		classname: "ArrayString",
		typename:  "String",
		defaultValue: d ,
		})
}

func (r *InstallLucyArray) RunCommand(command string, args []string) {
	path := os.Getenv("LUCYROOT")
	if path == "" {
		fmt.Println("env variable LUCYPATH is not set")
		os.Exit(1)
	}
	dest := filepath.Join(path, "lib/lucy/deps")
	os.MkdirAll(dest, 0755) // ignore errors
	err := os.Chdir(dest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, v := range installs {
		javafile := v.classname + ".java"
		t := strings.Replace(array_template, "ArrayTTT", v.classname, -1)
		t = strings.Replace(t, "TTT", v.typename, -1)
		t = strings.Replace(t, "DEFAULT_INIT", v.defaultValue , -1)
		err := ioutil.WriteFile(javafile, []byte(t), 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		cmd := exec.Command("javac", javafile)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
			if out != nil && len(out) > 0 {
				os.Stdout.Write(out)
			}
			os.Exit(3)
		}
	}
}

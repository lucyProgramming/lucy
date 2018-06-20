package clean

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gitee.com/yuyang-fine/lucy/src/cmd/common"
)

type Clean struct {
	lucyPaths []string
}

func (c *Clean) Help(command string) {

}

func (c *Clean) RunCommand(command string, args []string) {
	if os.Getenv(common.LUCY_PATH_ENV_KEY) == "" {
		fmt.Printf("env variable  '%s' not set\n", common.LUCY_PATH_ENV_KEY)
		return
	}
	if runtime.GOOS == "windows" {
		for _, v := range strings.Split(os.Getenv(common.LUCY_PATH_ENV_KEY), ";") { // windows style
			if v != "" {
				c.lucyPaths = append(c.lucyPaths, v)
			}
		}
	} else {
		for _, v := range strings.Split(os.Getenv(common.LUCY_PATH_ENV_KEY), ":") { // unix style
			if v != "" {
				c.lucyPaths = append(c.lucyPaths, v)
			}
		}
	}
	for _, v := range c.lucyPaths {
		if v == "" {
			continue
		}
		if filepath.IsAbs(v) == false {
			fmt.Printf("path '%s' is not absolute\n", v)
			return
		}
	}
	if len(c.lucyPaths) == 0 {
		fmt.Printf("env variable  '%s' not set\n", common.LUCY_PATH_ENV_KEY)
		return
	}
	for _, v := range c.lucyPaths {
		fmt.Println("clean:", v)
		c.cleanPath(filepath.Join(v, common.DIR_FOR_COMPILED_CLASS))
	}
}

/*
	don`t delete directory,in case directory have some other files
*/
func (c *Clean) cleanPath(path string) {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range fis {
		if v.IsDir() {
			c.cleanPath(filepath.Join(path, v.Name()))
		}
	}
	bs, err := ioutil.ReadFile(filepath.Join(path, common.LUCY_MAINTAIN_FILE))
	if err != nil {
		return
	}
	var meta common.PackageMeta
	err = json.Unmarshal(bs, &meta)
	if err != nil {
		return
	}
	for _, v := range meta.Classes {
		os.Remove(filepath.Join(path, v))
	}
	os.Remove(filepath.Join(path, common.LUCY_MAINTAIN_FILE))
}

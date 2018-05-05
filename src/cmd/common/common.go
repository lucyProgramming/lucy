package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	VERSION                   = "0.01"
	LUCY_ROOT_ENV_KEY         = "LUCYROOT"
	LUCY_PATH_ENV_KEY         = "LUCYPATH"
	LUCY_MAINTAIN_FILE        = "maintain.json"
	DIR_FOR_COMPILED_CLASS    = "class" // sub directory of $LUCYPATH
	DIR_FOR_LUCY_SOURCE_FILES = "src"
	CORE_PACAKGE              = "lucy/lang"
)

func GetClassPaths() []string {
	lp := os.Getenv("CLASSPATH")
	if runtime.GOOS == "windows" {
		return strings.Split(lp, ";")
	} else {
		return strings.Split(lp, ":")
	}
}

/*
	include lucy root
*/
func GetLucyPaths() ([]string, error) {
	root, err := GetLucyRoot()
	if err != nil {
		return nil, err
	}
	lp := os.Getenv(LUCY_PATH_ENV_KEY)
	if lp == "" {
		return nil, fmt.Errorf("env variable %s not set", LUCY_PATH_ENV_KEY)
	}
	var lps []string
	if runtime.GOOS == "windows" {
		lps = strings.Split(lp, ";")
	} else { // unix style
		lps = strings.Split(lp, ":")
	}
	lucypaths := []string{}
	for _, v := range lps {
		if v == "" {
			continue
		}
		if false == filepath.IsAbs(v) {
			fmt.Printf("env variable %s=%s is not absolute",
				LUCY_PATH_ENV_KEY, LUCY_PATH_ENV_KEY)
			os.Exit(1)
		}
		lucypaths = append(lucypaths, v)
	}
	lucypaths = append(lucypaths, root)
	lucypathMap := make(map[string]struct{})
	for _, v := range lucypaths {
		lucypathMap[v] = struct{}{}
	}
	lucypaths = make([]string, len(lucypathMap))
	i := 0
	for k, _ := range lucypathMap {
		lucypaths[i] = k
		i++
	}
	return lucypaths, nil
}

func GetLucyRoot() (string, error) {
	r := os.Getenv(LUCY_ROOT_ENV_KEY)
	if r == "" {
		return "", fmt.Errorf("env variable %s not set", LUCY_ROOT_ENV_KEY)
	}
	if false == filepath.IsAbs(r) {
		return "", fmt.Errorf("env variable %s=%s is not absolute",
			LUCY_ROOT_ENV_KEY, r)
	}
	return r, nil
}

func FindLucyPackageDirectory(packageName string, paths []string) []string {
	ret := []string{}
	for _, v := range paths {
		f, err := os.Stat(filepath.Join(v, DIR_FOR_LUCY_SOURCE_FILES, packageName))
		if err == nil && f.IsDir() {
			ret = append(ret, v)
		}
	}
	return ret
}

func SourceFileExist(path string) bool {
	f, _ := os.Stat(path)
	if f == nil || f.IsDir() == false {
		return false
	}
	fis, _ := ioutil.ReadDir(path)
	for _, f := range fis {
		if strings.HasSuffix(f.Name(), ".lucy") {
			return true
		}
	}
	return false
}

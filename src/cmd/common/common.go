package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	VERSION                = "0.01"
	LucyRootEnvVariableKey = "LUCYROOT"
	LucyPathEnvVariableKey = "LUCYPATH"
	LucyMaintainFile       = "maintain.json"
	DirForCompiledClass    = "class" // sub directory of compiled class
	DirForLucySourceFile   = "src"   // sub directory of source files
	CorePackage            = "lucy/lang"
)

func GetClassPaths() []string {
	lp := os.Getenv("CLASSPATH")
	if lp == "" {
		return []string{}
	}
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
	lp := os.Getenv(LucyPathEnvVariableKey)
	if lp == "" {
		return nil, fmt.Errorf("env variable %s not set", LucyPathEnvVariableKey)
	}
	var lps []string
	if runtime.GOOS == "windows" {
		lps = strings.Split(lp, ";")
	} else { // unix style
		lps = strings.Split(lp, ":")
	}
	lucyPaths := []string{}
	for _, v := range lps {
		if v == "" {
			continue
		}
		if false == filepath.IsAbs(v) {
			fmt.Printf("env variable %s=%s is not absolute",
				LucyPathEnvVariableKey, LucyPathEnvVariableKey)
			os.Exit(1)
		}
		lucyPaths = append(lucyPaths, v)
	}
	lucyPaths = append(lucyPaths, root)
	lucyPathMap := make(map[string]struct{})
	for _, v := range lucyPaths {
		lucyPathMap[v] = struct{}{}
	}
	lucyPaths = make([]string, len(lucyPathMap))
	i := 0
	for k, _ := range lucyPathMap {
		lucyPaths[i] = k
		i++
	}
	return lucyPaths, nil
}

func GetLucyRoot() (string, error) {
	r := os.Getenv(LucyRootEnvVariableKey)
	if r == "" {
		return "", fmt.Errorf("env variable %s not set", LucyRootEnvVariableKey)
	}
	if false == filepath.IsAbs(r) {
		return "", fmt.Errorf("env variable %s=%s is not absolute",
			LucyRootEnvVariableKey, r)
	}
	return r, nil
}

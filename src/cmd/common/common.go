package common

import (
	"fmt"
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
)

/*
	include
*/
func GetLucyPaths() ([]string, error) {
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

	return nil
}

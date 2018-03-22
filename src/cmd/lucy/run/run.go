package run

import (
	"encoding/json"
	"fmt"
	"github.com/756445638/lucy/src/cmd/common"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Run struct {
	LucyRoot            string
	LucyPaths           []string
	MainPackageLucyPath string
	Package             string
	command             string
	// PackageImported     map[string]struct{}
}

func (r *Run) printUsage() {

}

func (r *Run) RunCommand(command string, args []string) {
	r.command = command
	if len(args) != 1 {
		r.printUsage()
		os.Exit(1)
	}
	r.Package = args[0]
	r.LucyRoot = os.Getenv(common.LUCY_ROOT_ENV_KEY)
	if r.LucyRoot == "" {
		fmt.Printf("env variable %s not set", common.LUCY_ROOT_ENV_KEY)
		os.Exit(1)
	}
	if false == filepath.IsAbs(r.LucyRoot) {
		fmt.Printf("env variable %s=%s is not absolute", common.LUCY_ROOT_ENV_KEY, r.LucyRoot)
		os.Exit(1)
	}
	lp := os.Getenv(common.LUCY_PATH_ENV_KEY)
	if lp == "" {
		fmt.Printf("env variable %s not set", common.LUCY_PATH_ENV_KEY)
		os.Exit(1)
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
			fmt.Printf("env variable %s=%s is not absolute", common.LUCY_PATH_ENV_KEY, r.LucyRoot)
			os.Exit(1)
		}
		lucypaths = append(lucypaths, v)
	}
	r.LucyPaths = lucypaths
	var err error
	r.MainPackageLucyPath, err = r.findPackageIn(r.Package)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	//
	fis, err := ioutil.ReadDir(r.MainPackageLucyPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	lucyFiles := []string{}
	for _, f := range fis {
		if f.IsDir() == false && strings.HasSuffix(f.Name(), ".lucy") {
			lucyFiles = append(lucyFiles, f.Name())
		}
	}
	//
	if len(lucyFiles) == 0 {
		fmt.Printf("no lucy files in %s", filepath.Join(r.MainPackageLucyPath, r.Package))
		os.Exit(1)
	}
	_, err = r.buildPackage(r.MainPackageLucyPath, r.Package)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	//
}

/*
	find package in which directory
*/
func (r *Run) findPackageIn(packageName string) (string, error) {
	pathHavePackage := []string{}
	for _, v := range r.LucyPaths {
		dir := filepath.Join(v, r.Package)
		f, err := os.Stat(dir)
		if err == nil && f.IsDir() {
			fis, _ := ioutil.ReadDir(dir)
			for _, vv := range fis {
				if strings.HasSuffix(vv.Name(), ".lucy") {
					pathHavePackage = append(pathHavePackage, v)
					break
				}
			}
		}
	}
	if len(pathHavePackage) == 0 {
		return "", fmt.Errorf("package %s not found in $%s", r.Package, common.LUCY_PATH_ENV_KEY)
	}
	if len(pathHavePackage) > 1 {
		return "", fmt.Errorf("not 1 package named %s in $%s", r.Package, common.LUCY_PATH_ENV_KEY)
	}
	return pathHavePackage[0], nil
}

/*
	check package if need rebuild
*/
func (r *Run) needCompile(lucypath string, packageName string) (meta *common.PackageMeta, need bool, err error) {
	if lucypath == "" {
		lucypath, err = r.findPackageIn(packageName)
		if err != nil {
			return
		}
	}
	need = true
	_, err = os.Stat(filepath.Join(lucypath, packageName))
	if err != nil {
		err = nil
		return
	}
	bs, err := ioutil.ReadFile(filepath.Join(lucypath, "class", packageName, common.LUCY_MAINTAIN_FILE))
	if err != nil { // maintain file is missing
		err = nil
		return
	}
	meta = &common.PackageMeta{}
	err = json.Unmarshal(bs, meta)
	if err != nil { // this is not happening
		err = nil
		return
	}
	fis, err := ioutil.ReadDir(filepath.Join(lucypath, packageName))
	if err != nil { // shit happens
		return
	}
	fisM := make(map[string]os.FileInfo)
	for _, v := range fis {
		fisM[v.Name()] = v
		if meta.CompiledFrom == nil {
			return
		}
		if _, ok := meta.CompiledFrom[v.Name()]; ok == false { // new file
			return
		}
		if v.ModTime().After(meta.CompiledFrom[v.Name()].LastModify) { // modifyed
			return
		}
	}
	for f, _ := range meta.CompiledFrom {
		_, ok := fisM[f]
		if ok == false { // file deleted,missing file
			return
		}
	}
	need = false
	return
}

func (r *Run) buildPackage(lucypath string, packageName string) (needBuild bool, err error) {

	return false, nil
}

package run

import (
	"encoding/json"
	"fmt"
	"github.com/756445638/lucy/src/cmd/common"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Run struct {
	LucyRoot            string
	LucyPaths           []string
	MainPackageLucyPath string
	Package             string
	command             string
	compilerAt          string
	classPaths          []string
}

func (r *Run) printUsage() {
	fmt.Printf("%s    run a lucy package\n", r.command)
}

func (r *Run) RunCommand(command string, args []string) {
	r.command = command
	if len(args) != 1 {
		r.printUsage()
		os.Exit(0)
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
	if runtime.GOOS == "windows" {
		r.classPaths = strings.Split(os.Getenv("CLASSPATH"), ";")
	} else { // unix style
		r.classPaths = strings.Split(os.Getenv("CLASSPATH"), ":")
	}

	meta, need, lucyFiles, err := r.needCompile(r.MainPackageLucyPath, r.Package)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	//
	if len(lucyFiles) == 0 {
		fmt.Printf("no lucy files in %s", filepath.Join(r.MainPackageLucyPath, r.Package))
		os.Exit(1)
	}
	if need {
		r.compilerAt = filepath.Join(r.LucyRoot, "bin", "compile")
		_, meta, err = r.buildPackage(r.MainPackageLucyPath, r.Package)
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
	}
	//
	if meta.MainClass == "" {
		fmt.Println("has mo main fn,but package is still compiled")
		os.Exit(4)
	}
	cmd := exec.Command("java", meta.MainClass)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	//set CLASSPATHS
	classpath := r.classPaths
	for _, v := range r.LucyPaths {
		classpath = append(classpath, filepath.Join(v, common.DIR_FOR_COMPILED_CLASS))
	}
	if runtime.GOOS == "windows" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("CLASSPATH=%s", strings.Join(classpath, ";")))
	} else { // unix style
		cmd.Env = append(cmd.Env, fmt.Sprintf("CLASSPATH=%s", strings.Join(classpath, ":")))
	}
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		os.Exit(6)
	}
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
func (r *Run) needCompile(lucypath string, packageName string) (meta *common.PackageMeta, need bool, lucyFiles []string, err error) {
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
	fis, err := ioutil.ReadDir(filepath.Join(lucypath, common.DIR_FOR_LUCY_SOURCE_FILES, packageName))
	if err != nil { // shit happens
		return
	}
	lucyFiles = []string{}
	fisM := make(map[string]os.FileInfo)
	for _, v := range fis {
		if strings.HasSuffix(v.Name(), ".lucy") {
			lucyFiles = append(lucyFiles, v.Name())
		}
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
	if len(lucyFiles) == 0 {
		err = fmt.Errorf("no lucy source files in '%s'", filepath.Join(lucypath, common.DIR_FOR_LUCY_SOURCE_FILES, packageName))
		return
	}
	for f := range meta.CompiledFrom {
		_, ok := fisM[f]
		if ok == false { // file deleted,missing file
			return
		}
	}
	need = false
	return
}

func (r *Run) parseImports(files []string) ([]string, error) {
	args := append([]string{"-io"}, files...)
	cmd := exec.Command(r.compilerAt, args...)
	cmd.Stderr = os.Stderr
	bs, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	is := []string{}
	err = json.Unmarshal(bs, &is)
	if err != nil {
		return nil, err
	}
	//only lucy packages
	return r.javaPackageFilter(is)
}

/*
	panic out java package,java package cannot be build by 'lucy'
*/
func (r *Run) javaPackageFilter(is []string) (lucyPackages []string, err error) {
	f := func(name string) (found []string) {
		for _, v := range r.classPaths {
			dir := filepath.Join(v, name)
			f, _ := os.Stat(dir)
			if f != nil && f.IsDir() {
				fis, _ := ioutil.ReadDir(dir)
				for _, ff := range fis {
					if strings.HasSuffix(ff.Name(), ".class") {
						found = append(found, v)
						break
					}
				}
			}
		}
		return
	}
	for _, i := range is {
		found := f(i)
		if len(found) == 0 {
			lucyPackages = append(lucyPackages, i)
		}
		if len(found) > 1 {
			err = fmt.Errorf("not 1 package named '%s' in $CLASSPATH", i)
			return
		}
	}
	return
}

func (r *Run) buildPackage(lucypath string, packageName string) (needBuild bool, meta *common.PackageMeta, err error) {
	if lucypath == "" {
		lucypath, err = r.findPackageIn(packageName)
		if err != nil {
			return
		}
	}
	meta, needBuild, lucyFiles, err := r.needCompile(lucypath, packageName)
	if err != nil {
		return
	}
	if needBuild == false { // current package no need to compile,but I need to check dependies
		need := false
		for _, v := range meta.Imports {
			needBuild, _, err = r.buildPackage("", v)
			if err != nil {
				return
			}
			if needBuild { // means at least one package is rebuild
				need = true
			}
		}
		needBuild = need
	}
	if needBuild == false { // no need actually
		return
	}
	//compile this package really
	is, err := r.parseImports(lucyFiles)
	if err != nil {
		return
	}
	for _, i := range is {
		_, _, err = r.buildPackage("", i) // compile depend
		if err != nil {
			return
		}
	}
	// build this package
	//read  files
	destDir := filepath.Join(lucypath, common.DIR_FOR_COMPILED_CLASS, packageName)
	// mkdir all
	finfo, _ := os.Stat(destDir)
	if finfo == nil {
		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			return
		}
	}
	// cd to destDir
	os.Chdir(destDir)
	args := []string{"-pn", packageName}

	args = append(args, lucyFiles...)
	cmd := exec.Command(r.compilerAt, args...)
	cmd.Stderr = os.Stderr
	bs, err := cmd.Output()
	if err != nil {
		fmt.Println(string(bs))
		return
	}

	// make maitain.json
	maintain := &common.PackageMeta{}
	maintain.CompiledFrom = make(map[string]*common.FileMeta)
	for _, v := range lucyFiles {
		var f os.FileInfo
		f, err = os.Stat(v)
		if err != nil {
			return
		}
		maintain.CompiledFrom[v] = &common.FileMeta{
			LastModify: f.ModTime(),
		}
	}
	maintain.CompileTime = time.Now()
	maintain.Imports = is
	bs, err = json.Marshal(maintain)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filepath.Join(lucypath, common.DIR_FOR_COMPILED_CLASS, packageName, common.LUCY_MAINTAIN_FILE), bs, 0644)
	return
}

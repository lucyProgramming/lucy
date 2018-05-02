package run

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gitee.com/yuyang-fine/lucy/src/cmd/common"
)

type Run struct {
	LucyRoot            string
	LucyPaths           []string
	MainPackageLucyPath string
	Package             string
	command             string
	compilerAt          string
	classPaths          []string
	Flags               Flags
}

func (r *Run) Help(command string) {

}

func (r *Run) parseCmd(args []string) error {
	cmd := flag.NewFlagSet("run", flag.ErrorHandling(1))
	cmd.BoolVar(&r.Flags.forceReBuild, "forceReBuild", false, "force rebuild all package")
	err := cmd.Parse(args)
	if err != nil {
		return err
	}
	args = cmd.Args()
	if len(args) == 0 {
		return fmt.Errorf("no package to run")
	}
	r.Package = args[0]
	return nil
}

func (r *Run) RunCommand(command string, args []string) {
	r.command = command
	err := r.parseCmd(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	r.LucyRoot, err = common.GetLucyRoot()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r.LucyPaths, err = common.GetLucyPaths()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	founds := r.findPackageIn(r.Package)

	err = r.foundError(r.Package, founds)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	r.MainPackageLucyPath = founds[0]

	//
	r.classPaths = common.GetClassPaths()
	r.compilerAt = filepath.Join(r.LucyRoot, "bin", "compile") //compiler at
	{
		t := r.compilerAt
		if runtime.GOOS == "windows" {
			t += ".exe"
		}
		_, e := os.Stat(t)
		if e != nil {
			fmt.Println("compiler not found")
			return
		}
	}

	_, _, err = r.buildPackage(r.MainPackageLucyPath, r.Package)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	//
	cmd := exec.Command("java", r.Package+"/"+"main")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	//set CLASSPATHS
	classpath := make(map[string]struct{}) // if duplicate
	for _, v := range r.classPaths {
		classpath[v] = struct{}{}
	}
	for _, v := range r.LucyPaths {
		classpath[filepath.Join(v, common.DIR_FOR_COMPILED_CLASS)] = struct{}{}
	}
	classpath[filepath.Join(r.LucyRoot, "lib")] = struct{}{}
	classPathArray := make([]string, len(classpath))
	{
		i := 0
		for k := range classpath {
			classPathArray[i] = k
			i++
		}
	}
	if runtime.GOOS == "windows" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("CLASSPATH=%s", strings.Join(classPathArray, ";")))
	} else { // unix style
		cmd.Env = append(cmd.Env, fmt.Sprintf("CLASSPATH=%s", strings.Join(classPathArray, ":")))
	}
	//fmt.Println("CLASSPATH:", classPathArray)
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
func (r *Run) findPackageIn(packageName string) []string {
	pathHavePackage := []string{}
	for _, v := range r.LucyPaths {
		dir := filepath.Join(v, common.DIR_FOR_LUCY_SOURCE_FILES, packageName)
		f, err := os.Stat(dir)
		if err == nil && f.IsDir() {
			pathHavePackage = append(pathHavePackage, v)
		}
	}
	return pathHavePackage
}

/*
	check package if need rebuild
*/
func (r *Run) needCompile(lucypath string, packageName string) (meta *common.PackageMeta, need bool, lucyFiles []string, err error) {
	need = true
	sourceFileDir := filepath.Join(lucypath, common.DIR_FOR_LUCY_SOURCE_FILES, packageName)
	fis, err := ioutil.ReadDir(sourceFileDir)
	if err != nil { // shit happens
		return
	}
	fism := make(map[string]os.FileInfo)
	for _, v := range fis {
		if strings.HasSuffix(v.Name(), ".lucy") {
			lucyFiles = append(lucyFiles, filepath.Join(sourceFileDir, v.Name()))
			fism[v.Name()] = v
		}
	}
	if len(lucyFiles) == 0 {
		err = fmt.Errorf("no lucy source files in '%s'", filepath.Join(lucypath, common.DIR_FOR_LUCY_SOURCE_FILES, packageName))
		return
	}
	destDir := filepath.Join(lucypath, "class", packageName)
	bs, err := ioutil.ReadFile(filepath.Join(destDir, common.LUCY_MAINTAIN_FILE))
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
	// new or add
	for _, v := range fism {
		if meta.CompiledFrom == nil {
			return
		}
		if _, ok := meta.CompiledFrom[v.Name()]; ok == false { // new file
			return
		}
		if v.ModTime().After(meta.CompiledFrom[v.Name()].LastModify) { // modified
			return
		}
	}
	// file missing
	for f := range meta.CompiledFrom {
		_, ok := fism[f]
		if ok == false {
			return
		}
	}
	// if class file is missing
	fis, err = ioutil.ReadDir(destDir)
	fism = make(map[string]os.FileInfo)
	for _, v := range fis {
		if strings.HasSuffix(v.Name(), ".class") {
			fism[v.Name()] = v
		}
	}
	for _, v := range meta.Classes {
		_, ok := fism[v]
		if ok == false {
			return
		}
	}
	need = false
	return
}

func (r *Run) parseImports(files []string) ([]string, error) {
	args := append([]string{"-only-import"}, files...)
	cmd := exec.Command(r.compilerAt, args...)
	cmd.Stderr = os.Stderr
	bs, err := cmd.Output()
	if err != nil {
		fmt.Println(string(bs))
		return nil, err
	}
	is := []string{}
	err = json.Unmarshal(bs, &is)
	if err != nil {
		return nil, err
	}
	//only lucy packages
	//fmt.Println("!!!!!!!!!!!!!!!!", is)
	is, err = r.javaPackageFilter(is)
	//fmt.Println("@@@@@@@@@@@@@", is)
	return is, err
}

func (r *Run) javaPackageFilter(is []string) (lucyPackages []string, err error) {
	existInClassPath := func(name string) (found []string) {
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
	existInLucyPath := func(name string) (found []string) {
		for _, v := range r.LucyPaths {
			dir := filepath.Join(v, common.DIR_FOR_LUCY_SOURCE_FILES, name)
			f, _ := os.Stat(dir)
			if f != nil && f.IsDir() {
				found = append(found, v)
			}
		}
		return
	}
	for _, i := range is {
		found := existInLucyPath(i)
		if len(found) > 1 {
			err = fmt.Errorf("not 1 package named '%s' in $LUCYPATH", i)
			return
		}
		if len(found) == 1 { // perfect found in lucypath
			lucyPackages = append(lucyPackages, i)
			continue
		}
		found = existInClassPath(i)
		if len(found) > 1 {
			err = fmt.Errorf("not 1 package named '%s' in $CLASSPATH", i)
			return
		}
	}
	return
}

func (r *Run) foundError(packageName string, founds []string) error {
	if len(founds) == 0 {
		return fmt.Errorf("package '%s' not found")
	}
	if len(founds) > 1 {

	}
	return nil
}
func (r *Run) buildPackage(lucypath string, packageName string) (needBuild bool, meta *common.PackageMeta, err error) {
	if lucypath == "" {
		founds := r.findPackageIn(packageName)
		err = r.foundError(packageName, founds)
		if err != nil {
			return false, nil, err
		}
		lucypath = founds[0]
	}

	meta, needBuild, lucyFiles, err := r.needCompile(lucypath, packageName)
	if err != nil {
		err = fmt.Errorf("check if not compile,err:%v", err)
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
	fmt.Println("compiling.... ", packageName) // compile this package
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
	//before compile delete old class and maintain.json
	{
		fis, _ := ioutil.ReadDir(destDir)
		for _, f := range fis {
			if strings.HasSuffix(f.Name(), ".class") || f.Name() == common.LUCY_MAINTAIN_FILE {
				file := filepath.Join(destDir, f.Name())
				err := os.Remove(file)
				if err != nil {
					fmt.Printf("delete old compiled file[%s] failed,err:%v\n", file, err)
				}
			}
		}
	}
	// cd to destDir
	os.Chdir(destDir)
	args := []string{"-package-name", packageName}
	args = append(args, lucyFiles...)
	cmd := exec.Command(r.compilerAt, args...)
	cmd.Stderr = os.Stderr
	bs, err := cmd.Output()
	if err != nil {
		fmt.Println(string(bs))
		return
	}
	// make maitain.json
	meta = &common.PackageMeta{}
	meta.CompiledFrom = make(map[string]*common.FileMeta)
	for _, v := range lucyFiles {
		var f os.FileInfo
		f, err = os.Stat(v)
		if err != nil {
			return
		}
		meta.CompiledFrom[filepath.Base(v)] = &common.FileMeta{
			LastModify: f.ModTime(),
		}
	}
	meta.CompileTime = time.Now()
	meta.Imports = is
	fis, err := ioutil.ReadDir(destDir)
	if err != nil {
		return
	}
	for _, v := range fis {
		if strings.HasSuffix(v.Name(), ".class") {
			meta.Classes = append(meta.Classes, v.Name())
		}
	}
	bs, err = json.MarshalIndent(meta, "", "\t")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(
		filepath.Join(lucypath, common.DIR_FOR_COMPILED_CLASS, packageName, common.LUCY_MAINTAIN_FILE),
		bs,
		0644)
	return
}

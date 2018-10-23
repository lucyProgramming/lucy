package run

import (
	"encoding/json"
	"flag"
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/common"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type RunLucyPackage struct {
	lucyRoot            string
	lucyPaths           []string
	mainPackageLucyPath string
	Package             string
	command             string
	compilerExe         string
	classPaths          []string
	flags               Flags
	packagesCompiled    map[string]*PackageCompiled
	lucyProgramArgs     []string // lucy application args
	flagSet             *flag.FlagSet
}

func (runLucyPackage *RunLucyPackage) Help() {
	fmt.Println("run a lucy package")
	var fs flagSet
	runLucyPackage.flagSet.VisitAll(func(f *flag.Flag) {
		t := &flag.Flag{}
		*t = *f
		t.DefValue = fmt.Sprintf(`'%s'`, t.DefValue)
		fs = append(fs, t)
	})
	fs.makeSureLengthIsSame()
	for _, v := range fs {
		fmt.Printf("\t -%s\t default:%s\t%s\n", v.Name, v.DefValue, v.Usage)
	}
}

func (runLucyPackage *RunLucyPackage) parseCmd(args []string) error {
	cmd := flag.NewFlagSet("run", flag.ErrorHandling(1))
	runLucyPackage.flagSet = cmd
	cmd.BoolVar(&runLucyPackage.flags.forceReBuild,
		"forceReBuild", false, "force rebuild all package")
	cmd.StringVar(&runLucyPackage.flags.compilerFlags,
		"cf", "", "compiler flags")
	cmd.BoolVar(&runLucyPackage.flags.build,
		"build", false, "build package and no run")
	cmd.BoolVar(&runLucyPackage.flags.verbose,
		"v", false, "verbose")
	cmd.BoolVar(&runLucyPackage.flags.help, "h", false,
		"print help message")
	var runArgs []string
	var lucyProgramArgs []string
	for k, v := range args {
		if strings.HasPrefix(v, "-") == false {
			runLucyPackage.Package = v
			lucyProgramArgs = args[k+1:]
			break
		}
		runArgs = append(runArgs, v)
	}
	err := cmd.Parse(runArgs)
	if err != nil {
		return err
	}
	runLucyPackage.lucyProgramArgs = lucyProgramArgs
	return nil
}

func (runLucyPackage *RunLucyPackage) setCompiler() error {
	runLucyPackage.compilerExe = filepath.Join(runLucyPackage.lucyRoot, "bin", "compile") //compiler at
	t := runLucyPackage.compilerExe
	if runtime.GOOS == "windows" {
		t += ".exe"
	}
	_, e := os.Stat(t)
	if e != nil {
		return fmt.Errorf("compiler not found")
	}
	return nil
}

func (runLucyPackage *RunLucyPackage) RunCommand(command string, args []string) {
	runLucyPackage.command = command
	err := runLucyPackage.parseCmd(args) // skip run
	if err != nil {
		fmt.Println(err)
		runLucyPackage.Help()
		return
	}
	if runLucyPackage.flags.help {
		runLucyPackage.Help()
		return
	}
	if runLucyPackage.Package == "" {
		fmt.Println("no package to run")
		os.Exit(1)
	}
	runLucyPackage.lucyRoot, err = common.GetLucyRoot()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	runLucyPackage.lucyPaths, err = common.GetLucyPaths()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	runLucyPackage.classPaths = common.GetClassPaths()
	err = runLucyPackage.setCompiler()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	runLucyPackage.packagesCompiled = make(map[string]*PackageCompiled)
	founds := runLucyPackage.findPackageIn(runLucyPackage.Package)
	err = runLucyPackage.foundError(runLucyPackage.Package, founds)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	runLucyPackage.mainPackageLucyPath = founds[0]
	//
	{
		_, _, err = runLucyPackage.buildPackage("", common.CorePackage, &ImportStack{})
		if err != nil {
			fmt.Printf("build  buildin package '%s' failed,err:%v\n", common.CorePackage, err)
			os.Exit(3)
		}
	}
	_, _, err = runLucyPackage.buildPackage(runLucyPackage.mainPackageLucyPath, runLucyPackage.Package, &ImportStack{})
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	if runLucyPackage.flags.build {
		os.Exit(0)
	}
	//
	cmd := exec.Command("java",
		append([]string{runLucyPackage.Package + "/" + "main"}, runLucyPackage.lucyProgramArgs...)...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	//set CLASSPATHS
	classpath := make(map[string]struct{}) // if duplicate
	for _, v := range runLucyPackage.classPaths {
		classpath[v] = struct{}{}
	}
	for _, v := range runLucyPackage.lucyPaths {
		classpath[filepath.Join(v, common.DirForCompiledClass)] = struct{}{}
	}
	classpath[filepath.Join(runLucyPackage.lucyRoot, "lib")] = struct{}{}
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
	{
		envs := os.Environ()
		ts := []string{}
		for _, v := range envs {
			if strings.HasPrefix(v, "CLASSPATH=") {
				continue
			}
			ts = append(ts, v)
		}
		cmd.Env = append(cmd.Env, ts...)
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
func (runLucyPackage *RunLucyPackage) findPackageIn(packageName string) []string {
	pathHavePackage := []string{}
	for _, v := range runLucyPackage.lucyPaths {
		dir := filepath.Join(v, common.DirForLucySourceFile, packageName)
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
func (runLucyPackage *RunLucyPackage) needCompile(lucyPath string, packageName string) (meta *common.PackageMeta,
	need bool, lucyFiles []string, err error) {
	need = true
	sourceFileDir := filepath.Join(lucyPath, common.DirForLucySourceFile, packageName)
	fis, err := ioutil.ReadDir(sourceFileDir)
	if err != nil { // shit happens
		return
	}
	fisM := make(map[string]os.FileInfo)
	for _, v := range fis {
		if strings.HasSuffix(v.Name(), ".lucy") {
			lucyFiles = append(lucyFiles, filepath.Join(sourceFileDir, v.Name()))
			fisM[v.Name()] = v
		}
	}
	if len(lucyFiles) == 0 {
		err = fmt.Errorf("no lucy source files in '%s'",
			filepath.Join(lucyPath, common.DirForLucySourceFile, packageName))
		return
	}
	if p, ok := runLucyPackage.packagesCompiled[packageName]; ok {
		return p.meta, false, nil, nil
	}
	if runLucyPackage.flags.forceReBuild {
		return
	}
	destinationDir := filepath.Join(lucyPath, "class", packageName)
	bs, err := ioutil.ReadFile(filepath.Join(destinationDir, common.LucyMaintainFile))
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
	for _, v := range fisM {
		if meta.CompiledFrom == nil {
			return
		}
		if _, ok := meta.CompiledFrom[v.Name()]; ok == false { // new file
			return
		}
		if v.ModTime().Unix() > (meta.CompiledFrom[v.Name()].LastModify) { // modified
			return
		}
	}
	// file missing
	for f := range meta.CompiledFrom {
		_, ok := fisM[f]
		if ok == false {
			return
		}
	}
	// if class file is missing
	fis, err = ioutil.ReadDir(destinationDir)
	fisM = make(map[string]os.FileInfo)
	for _, v := range fis {
		if strings.HasSuffix(v.Name(), ".class") {
			fisM[v.Name()] = v
		}
	}
	for _, v := range meta.Classes {
		_, ok := fisM[v]
		if ok == false {
			return
		}
	}
	need = false
	return
}

func (runLucyPackage *RunLucyPackage) parseImports(files []string) ([]string, error) {
	args := append([]string{"-only-import"}, files...)
	cmd := exec.Command(runLucyPackage.compilerExe, args...)
	cmd.Stderr = os.Stderr
	bs, err := cmd.Output()
	if err != nil {
		fmt.Println(string(bs))
		return nil, err
	}
	is := []string{}
	err = json.Unmarshal(bs, &is)
	if err != nil {
		return nil, fmt.Errorf("parse import failed,err:%v", err)
	}
	isM := make(map[string]struct{})
	for _, v := range is {
		isM[v] = struct{}{}
	}
	is = []string{}
	for k, _ := range isM {
		is = append(is, k)
	}
	is, err = runLucyPackage.javaPackageFilter(is)
	if err != nil {
		err = fmt.Errorf("parse import failed,err:%v", err)
	}
	return is, err
}

func (runLucyPackage *RunLucyPackage) javaPackageFilter(is []string) (lucyPackages []string, err error) {
	existInClassPath := func(name string) (found []string) {
		for _, v := range runLucyPackage.classPaths {
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
				continue
			}
			dir = filepath.Join(v, name+".class")
			f, _ = os.Stat(dir)
			if f != nil && f.IsDir() == false {
				found = append(found, v)
			}
		}
		return
	}
	existInLucyPath := func(name string) (found []string) {
		for _, v := range runLucyPackage.lucyPaths {
			dir := filepath.Join(v, common.DirForLucySourceFile, name)
			f, _ := os.Stat(dir)
			if f != nil && f.IsDir() {
				found = append(found, v)
			}
		}
		return
	}
	formatPaths := func(paths []string) string {
		var s string
		for _, v := range paths {
			s += "\t" + v + "\n"
		}
		return s
	}
	for _, i := range is {
		found := existInLucyPath(i)
		if len(found) > 1 {
			err = fmt.Errorf("not 1 package named '%s' in $LUCYPATH", i)
			return
		}
		if len(found) == 1 { // perfect found in lucyPath
			if i != common.CorePackage {
				lucyPackages = append(lucyPackages, i)
			}
			continue
		}
		found = existInClassPath(i)
		if len(found) > 1 {
			errMsg := fmt.Sprintf("not 1 package named '%s' in $CLASSPATH,which CLASSPATH are:\n", i)
			errMsg += formatPaths(runLucyPackage.classPaths)
			return nil, fmt.Errorf(errMsg)
		}
		if len(found) == 0 {
			errMsg := fmt.Sprintf("package named '%s' not found in $CLASSPATH,which CLASSPATH are:\n", i)
			errMsg += formatPaths(runLucyPackage.classPaths)
			return nil, fmt.Errorf(errMsg)
		}
	}
	return
}

func (runLucyPackage *RunLucyPackage) foundError(packageName string, founds []string) error {
	if len(founds) == 0 {
		return fmt.Errorf("package '%s' not found", packageName)
	}
	if len(founds) > 1 {

	}
	return nil
}

func (runLucyPackage *RunLucyPackage) buildPackage(lucyPath string, packageName string, importStack *ImportStack) (needBuild bool,
	meta *common.PackageMeta, err error) {
	if p, ok := runLucyPackage.packagesCompiled[packageName]; ok {
		return false, p.meta, nil
	}
	if lucyPath == "" {
		founds := runLucyPackage.findPackageIn(packageName)
		err = runLucyPackage.foundError(packageName, founds)
		if err != nil {
			return false, nil, err
		}
		lucyPath = founds[0]
	}
	meta, needBuild, lucyFiles, err := runLucyPackage.needCompile(lucyPath, packageName)
	if err != nil {
		err = fmt.Errorf("check if need compile,err:%v", err)
		return
	}
	if needBuild == false { //current package no need to compile,but I need to check dependies
		need := false
		for _, v := range meta.Imports {
			if _, ok := runLucyPackage.packagesCompiled[v]; ok {
				continue
			}
			i := (&ImportStack{}).fromLast(importStack)
			err = i.insert(&PackageCompiled{packageName: packageName})
			if err != nil {
				return
			}
			needBuild, _, err = runLucyPackage.buildPackage("", v, i)
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
	is, err := runLucyPackage.parseImports(lucyFiles)
	if err != nil {
		return
	}
	for _, i := range is {
		if _, ok := runLucyPackage.packagesCompiled[i]; ok {
			continue
		}
		im := (&ImportStack{}).fromLast(importStack)
		err = im.insert(&PackageCompiled{packageName: packageName})
		if err != nil {
			return
		}
		_, _, err = runLucyPackage.buildPackage("", i, im) // compile depend
		if err != nil {
			return
		}
	}
	//build this package
	//read  files
	destinationDir := filepath.Join(lucyPath, common.DirForCompiledClass, packageName)
	// mkDir all
	fileInfo, _ := os.Stat(destinationDir)
	if fileInfo == nil {
		err = os.MkdirAll(destinationDir, 0755)
		if err != nil {
			return
		}
	}
	//before compile delete old class and maintain.json
	{
		fis, _ := ioutil.ReadDir(destinationDir)
		for _, f := range fis {
			if strings.HasSuffix(f.Name(), ".class") || f.Name() == common.LucyMaintainFile {
				file := filepath.Join(destinationDir, f.Name())
				err := os.Remove(file)
				if err != nil {
					fmt.Printf("delete old compiled file[%s] failed,err:%v\n", file, err)
				}
			}
		}
	}
	if runLucyPackage.flags.verbose {
		fmt.Printf("# %s\n", packageName) // compile this package
	}
	// cd to destDir
	os.Chdir(destinationDir)
	args := []string{"-package-name", packageName}
	if runLucyPackage.flags.compilerFlags != "" {
		args = append(args, strings.Split(runLucyPackage.flags.compilerFlags, " ")...)
	}
	args = append(args, lucyFiles...)
	cmd := exec.Command(runLucyPackage.compilerExe, args...)
	cmd.Stderr = os.Stderr
	bs, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("compiler err:%v", err)
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
			LastModify: f.ModTime().Unix(),
		}
	}
	meta.CompileTime = time.Now().Unix()
	meta.Imports = is
	fis, err := ioutil.ReadDir(destinationDir)
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
		filepath.Join(lucyPath, common.DirForCompiledClass, packageName, common.LucyMaintainFile),
		bs,
		0644)
	runLucyPackage.packagesCompiled[packageName] = &PackageCompiled{
		meta:        meta,
		packageName: packageName,
	}
	return
}

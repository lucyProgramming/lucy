package lc

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/common"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type RealNameLoader struct {
	Package *ast.Package
}

func (loader *RealNameLoader) LoadName(p *ast.Package, packageName string, name string) (interface{}, error) {
	loader.Package = p
	var realpath []string
	for _, v := range compiler.ClassPath {
		f, err := os.Stat(filepath.Join(v, packageName))
		if err == nil && f.IsDir() { // directory is package
			realpath = append(realpath, filepath.Join(v, packageName))
			loader.Package.Kind = ast.PACKAGE_KIND_JAVA
		}
	}
	for _, v := range compiler.lucyPath {
		f, err := os.Stat(filepath.Join(v, "classes", packageName))
		if err == nil && f.IsDir() { // directory is package
			realpath = append(realpath, filepath.Join(v, packageName))
			loader.Package.Kind = ast.PACKAGE_KIND_LUCY
		}
	}
	if len(realpath) == 0 {
		return nil, fmt.Errorf("package '%v' not found", packageName)
	}
	if len(realpath) > 1 {
		dirs := ""
		for _, v := range realpath {
			dirs += v + " "
		}
		return nil, fmt.Errorf("not 1 package named '%s' present in '%s' ", packageName, dirs)
	}
	var fis []os.FileInfo
	var err error
	if loader.Package.Kind == ast.PACKAGE_KIND_LUCY {
		fis, err = ioutil.ReadDir(realpath[0])
	} else { //java package
		var f os.FileInfo
		f, err = os.Stat(filepath.Join(realpath[0], name+".class"))
		if err == nil {
			fis = append(fis, f)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("read dir '%s' failed,err:%v\n", realpath[0], err)
	}
	classfiles := make(map[string]os.FileInfo)
	for _, v := range fis {
		if strings.HasSuffix(v.Name(), ".class") { // class file is all I need
			classfiles[v.Name()] = v
		}
	}
	if len(classfiles) == 0 {
		return nil, fmt.Errorf("package '%s' has no class files", packageName)
	}
	err = (&RealNameLoader{Package: loader.Package}).load(realpath[0], name, classfiles)
	if err != nil {
		return nil, err
	}
	if loader.Package.Kind == ast.PACKAGE_KIND_LUCY {
		return loader.Package.Block.SearchByName(name), nil
	} else {
		return loader.Package.Block.Classes[name], nil
	}

}

func (loader *RealNameLoader) load(realpath string, name string, files map[string]os.FileInfo) error {
	if loader.Package.Kind == ast.PACKAGE_KIND_LUCY {
		if _, ok := files[common.MAIN_CLASS_NAME]; ok == false {
			return fmt.Errorf("no main class found")
		}
		bs, err := ioutil.ReadFile(filepath.Join(realpath, common.MAIN_CLASS_NAME))
		if err != nil {
			return err
		}
		c, err := (&ClassDecoder{}).decode(bs)
		if err != nil {
			return err
		}
		err = loader.loadLucyMainClass(loader.Package, c)
		if err != nil {
			return err
		}
		delete(files, common.MAIN_CLASS_NAME)
	}

	for _, f := range files {
		bs, err := ioutil.ReadFile(filepath.Join(realpath, f.Name()))
		if err != nil {
			return err
		}
		class, err := (&ClassDecoder{}).decode(bs)
		if err != nil {
			return err
		}
		if loader.Package.Kind == ast.PACKAGE_KIND_JAVA {
			astClass, err := loader.loadAsJava(class)
			if err != nil {
				return err
			}
			if loader.Package.Block.Classes == nil {
				loader.Package.Block.Classes = make(map[string]*ast.Class)
			}
			loader.Package.Block.Classes[name] = astClass
		} else {
			astClass, err := loader.loadAsLucy(class)
			if err != nil {
				return err
			}
			if astClass != nil {
				if loader.Package.Block.Classes == nil {
					loader.Package.Block.Classes = make(map[string]*ast.Class)
				}
				loader.Package.Block.Classes[name] = astClass
			}
		}
	}
	return nil
}

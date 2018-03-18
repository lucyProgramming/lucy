package lc

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type RealNameLoader struct {
	Package *ast.Package
}

func (this *RealNameLoader) LoadName(p *ast.Package, packageName string, name string) (interface{}, error) {
	this.Package = p
	var realpath []string
	for _, v := range compiler.ClassPath {
		f, err := os.Stat(filepath.Join(v, packageName))
		if err == nil && f.IsDir() { // directory is package
			realpath = append(realpath, filepath.Join(v, packageName))
			this.Package.Kind = ast.PACKAGE_KIND_JAVA
		}
	}
	for _, v := range compiler.lucyPath {
		f, err := os.Stat(filepath.Join(v, "classes", packageName))
		if err == nil && f.IsDir() { // directory is package
			realpath = append(realpath, filepath.Join(v, packageName))
			this.Package.Kind = ast.PACKAGE_KIND_LUCY
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
	if this.Package.Kind == ast.PACKAGE_KIND_LUCY {
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
	classfiles := []os.FileInfo{}
	for _, v := range fis {
		if strings.HasSuffix(v.Name(), ".class") { // class file is all I need
			classfiles = append(classfiles, v)
		}
	}
	if len(classfiles) == 0 {
		return nil, fmt.Errorf("package '%s' has no class files")
	}
	err = (&RealNameLoader{Package: this.Package}).load(realpath[0], name, classfiles)
	if err != nil {
		return nil, err
	}
	if this.Package.Kind == ast.PACKAGE_KIND_LUCY {
		return this.Package.Block.SearchByName(name), nil
	} else {
		return this.Package.Block.Classes[name], nil
	}

}

func (p *RealNameLoader) load(realpath string, name string, files []os.FileInfo) error {
	for _, f := range files {
		bs, err := ioutil.ReadFile(filepath.Join(realpath, f.Name()))
		if err != nil {
			return err
		}
		class, err := (&ClassDecoder{}).decode(bs)
		if err != nil {
			return err
		}
		if p.Package.Kind == ast.PACKAGE_KIND_JAVA {
			astClass, err := p.loadAsJava(class)
			if err != nil {
				return err
			}
			if p.Package.Block.Classes == nil {
				p.Package.Block.Classes = make(map[string]*ast.Class)
			}
			p.Package.Block.Classes[name] = astClass
			for k, v := range p.Package.Block.Classes {
				fmt.Println(k, v)
			}

		} else {
			panic("load lucy")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

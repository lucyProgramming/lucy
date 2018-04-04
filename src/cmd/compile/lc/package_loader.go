package lc

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type RealNameLoader struct {
}

const (
	_ = iota
	RESOUCE_KIND_JAVA_CLASS
	RESOUCE_KIND_JAVA_PACKAGE
	RESOUCE_KIND_LUCY_CLASS
	RESOUCE_KIND_LUCY_PACKAGE
)

type Resource struct {
	kind     int
	realpath string
}

func (loader *RealNameLoader) LoadName(resouceName string) (interface{}, error) {
	var realpath []*Resource
	for _, v := range compiler.ClassPath {
		p := filepath.Join(v, resouceName)
		f, err := os.Stat(p)
		if err == nil && f.IsDir() { // directory is package
			realpath = append(realpath, &Resource{
				kind:     RESOUCE_KIND_LUCY_CLASS,
				realpath: p})
		}
		p = filepath.Join(v, resouceName+".class")
		f, err = os.Stat(p)
		if err == nil && f.IsDir() == false { // directory is package
			realpath = append(realpath, &Resource{
				kind:     RESOUCE_KIND_JAVA_CLASS,
				realpath: p})
		}
	}
	for _, v := range compiler.lucyPath {
		p := filepath.Join(v, "class", resouceName)
		f, err := os.Stat(p)
		if err == nil && f.IsDir() { // directory is package
			realpath = append(realpath, &Resource{
				kind:     RESOUCE_KIND_LUCY_PACKAGE,
				realpath: p})
		}
		p = filepath.Join(v, "class", resouceName+".class")
		f, err = os.Stat(p)
		if err == nil && f.IsDir() == false { // directory is package
			realpath = append(realpath, &Resource{
				kind:     RESOUCE_KIND_LUCY_CLASS,
				realpath: p})
		}
	}
	if len(realpath) == 0 {
		return nil, fmt.Errorf("resource '%v' not found", resouceName)
	}
	if len(realpath) > 1 {
		errMsg := "not 1 resource named '" + resouceName + "' present: "
		for _, v := range realpath {
			switch v.kind {
			case RESOUCE_KIND_JAVA_CLASS:
				errMsg += fmt.Sprintf("\t %s is a java class\n", v.realpath)
			case RESOUCE_KIND_JAVA_PACKAGE:
				errMsg += fmt.Sprintf("\t %s is a java package\n", v.realpath)
			case RESOUCE_KIND_LUCY_CLASS:
				errMsg += fmt.Sprintf("\t %s is a lucy class\n", v.realpath)
			case RESOUCE_KIND_LUCY_PACKAGE:
				errMsg += fmt.Sprintf("\t %s is a lucy package\n", v.realpath)
			}
		}
		return nil, fmt.Errorf(errMsg)
	}
	if realpath[0].kind == RESOUCE_KIND_JAVA_CLASS || realpath[0].kind == RESOUCE_KIND_JAVA_CLASS {
		return loader.loadClass(realpath[0])
	} else if realpath[0].kind == RESOUCE_KIND_JAVA_PACKAGE {
		return loader.loadJavaPackage(realpath[0])
	} else {
		return loader.loadLucyPackage(realpath[0])
	}
}
func (loader *RealNameLoader) loadLucyPackage(r *Resource) (*ast.Package, error) {
	return nil, nil
}
func (loader *RealNameLoader) loadJavaPackage(r *Resource) (*ast.Package, error) {
	fis, err := ioutil.ReadDir(r.realpath)
	if err != nil {
		return nil, err
	}
	ret := &ast.Package{}
	ret.Block.Classes = make(map[string]*ast.Class)
	for _, v := range fis {
		var rr Resource
		rr.kind = RESOUCE_KIND_JAVA_CLASS
		if strings.HasSuffix(v.Name(), ".class") {
			rr.realpath = filepath.Join(r.realpath, v.Name())
		}
		class, err := loader.loadClass(&rr)
		if err != nil {
			return nil, err
		}
		ret.Block.Classes[filepath.Base(class.Name)] = class
	}
	return ret, nil
}
func (loader *RealNameLoader) loadClass(r *Resource) (*ast.Class, error) {
	bs, err := ioutil.ReadFile(r.realpath)
	if err != nil {
		return nil, err
	}
	c, err := (&ClassDecoder{}).decode(bs)
	if r.kind == RESOUCE_KIND_LUCY_CLASS {
		return loader.loadAsLucy(c)
	}
	return loader.loadAsJava(c)
}

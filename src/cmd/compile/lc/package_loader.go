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
	name     string
}

func (loader *RealNameLoader) LoadName(resouceName string) (*ast.Package, interface{}, error) {
	var realpaths []*Resource
	for _, v := range compiler.lucyPath {
		p := filepath.Join(v, "class", resouceName)
		f, err := os.Stat(p)
		if err == nil && f.IsDir() { // directory is package
			realpaths = append(realpaths, &Resource{
				kind:     RESOUCE_KIND_LUCY_PACKAGE,
				realpath: p,
				name:     resouceName,
			})
		}
		p = filepath.Join(v, "class", resouceName+".class")
		f, err = os.Stat(p)
		if err == nil && f.IsDir() == false { // directory is package
			realpaths = append(realpaths, &Resource{
				kind:     RESOUCE_KIND_LUCY_CLASS,
				realpath: p,
				name:     resouceName,
			})
		}
	}

	for _, v := range compiler.ClassPath {
		p := filepath.Join(v, resouceName)
		f, err := os.Stat(p)
		if err == nil && f.IsDir() { // directory is package
			realpaths = append(realpaths, &Resource{
				kind:     RESOUCE_KIND_JAVA_PACKAGE,
				realpath: p,
				name:     resouceName,
			})
		}
		p = filepath.Join(v, resouceName+".class")
		f, err = os.Stat(p)
		if err == nil && f.IsDir() == false { // directory is package
			realpaths = append(realpaths, &Resource{
				kind:     RESOUCE_KIND_JAVA_CLASS,
				realpath: p,
				name:     resouceName,
			})
		}
	}
	if len(realpaths) == 0 {
		return nil, nil, fmt.Errorf("resource '%v' not found", resouceName)
	}
	realpathMap := make(map[string]*Resource)
	for _, v := range realpaths {
		realpathMap[v.realpath] = v
	}
	if len(realpathMap) > 1 {
		errMsg := "not 1 resource named '" + resouceName + "' present:\n"
		for _, v := range realpaths {
			switch v.kind {
			case RESOUCE_KIND_JAVA_CLASS:
				errMsg += fmt.Sprintf("\t '%s' is a java class\n", v.realpath)
			case RESOUCE_KIND_JAVA_PACKAGE:
				errMsg += fmt.Sprintf("\t '%s' is a java package\n", v.realpath)
			case RESOUCE_KIND_LUCY_CLASS:
				errMsg += fmt.Sprintf("\t '%s' is a lucy class\n", v.realpath)
			case RESOUCE_KIND_LUCY_PACKAGE:
				errMsg += fmt.Sprintf("\t '%s' is a lucy package\n", v.realpath)
			}
		}
		return nil, nil, fmt.Errorf(errMsg)
	}
	if realpaths[0].kind == RESOUCE_KIND_LUCY_CLASS {
		if filepath.Base(realpaths[0].realpath) == mainClassName {
			return nil, nil, fmt.Errorf("%s is special class for global variable and ...", mainClassName)
		}
	}
	if realpaths[0].kind == RESOUCE_KIND_JAVA_CLASS {
		class, err := loader.loadClass(realpaths[0])
		return nil, class, err
	} else if realpaths[0].kind == RESOUCE_KIND_LUCY_CLASS {
		fmt.Println(realpaths[0].realpath)
		name := filepath.Base(realpaths[0].realpath)
		name = strings.TrimRight(name, ".class")
		realpaths[0].name = filepath.Dir(resouceName)
		realpaths[0].kind = RESOUCE_KIND_LUCY_PACKAGE
		realpaths[0].realpath = filepath.Dir(realpaths[0].realpath)
		p, err := loader.loadLucyPackage(realpaths[0])
		return p, p.Block.SearchByName(name), err
	} else if realpaths[0].kind == RESOUCE_KIND_JAVA_PACKAGE {
		p, err := loader.loadJavaPackage(realpaths[0])
		return nil, p, err
	} else {
		p, err := loader.loadLucyPackage(realpaths[0])
		return nil, p, err
	}
}
func (loader *RealNameLoader) loadLucyPackage(r *Resource) (*ast.Package, error) {
	fis, err := ioutil.ReadDir(r.realpath)
	if err != nil {
		return nil, err
	}
	fisM := make(map[string]os.FileInfo)
	for _, v := range fis {
		if strings.HasSuffix(v.Name(), ".class") {
			fisM[v.Name()] = v
		}
	}

	_, ok := fisM[mainClassName]
	if ok == false {
		return nil, fmt.Errorf("main class not found")
	}
	bs, err := ioutil.ReadFile(filepath.Join(r.realpath, mainClassName))
	if err != nil {
		return nil, fmt.Errorf("read main.class failed,err:%v", err)
	}
	c, err := (&ClassDecoder{}).decode(bs)
	if err != nil {
		return nil, fmt.Errorf("decode main class failed,err:%v", err)
	}
	p := &ast.Package{}
	p.Name = r.name
	err = loader.loadLucyMainClass(p, c)
	if err != nil {
		return nil, fmt.Errorf("parse main class failed,err:%v", err)
	}
	delete(fisM, mainClassName)
	for _, v := range fisM {
		bs, err := ioutil.ReadFile(filepath.Join(r.realpath, v.Name()))
		if err != nil {
			return p, fmt.Errorf("read class failed,err:%v", err)
		}
		c, err := (&ClassDecoder{}).decode(bs)
		if err != nil {
			return nil, fmt.Errorf("decode class failed,err:%v", err)
		}
		class, err := loader.loadAsLucy(c)
		if err != nil {
			return nil, fmt.Errorf("decode class failed,err:%v", err)
		}
		if p.Block.Classes == nil {
			p.Block.Classes = make(map[string]*ast.Class)
		}
		p.Block.Classes[filepath.Base(class.Name)] = class
	}
	return p, nil
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
		if strings.HasSuffix(v.Name(), ".class") == false || strings.Contains(v.Name(), "$") {
			continue
		}
		rr.realpath = filepath.Join(r.realpath, v.Name())
		class, err := loader.loadClass(&rr)
		if err != nil {
			if _, ok := err.(*NotSupportTypeSignatureError); ok == false {
				return nil, err
			}
		}
		if class != nil {
			ret.Block.Classes[filepath.Base(class.Name)] = class
		}
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

package lc

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"
)

type RealNameLoader struct {
	caches map[string]interface{}
}

const (
	_ = iota
	RESOURCE_KIND_JAVA_CLASS
	RESOURCE_KIND_JAVA_PACKAGE
	RESOURCE_KIND_LUCY_CLASS
	RESOURCE_KIND_LUCY_PACKAGE
)

type Resource struct {
	kind     int
	realpath string
	name     string
}

func (loader *RealNameLoader) LoadName(resouceName string) (interface{}, error) {
	if loader.caches != nil && loader.caches[resouceName] != nil {
		return loader.caches[resouceName], nil
	}
	//fmt.Println("load:", resouceName)

	var realpaths []*Resource
	for _, v := range compiler.lucyPaths {
		p := filepath.Join(v, "class", resouceName)
		f, err := os.Stat(p)
		if err == nil && f.IsDir() { // directory is package
			realpaths = append(realpaths, &Resource{
				kind:     RESOURCE_KIND_LUCY_PACKAGE,
				realpath: p,
				name:     resouceName,
			})
		}
		p = filepath.Join(v, "class", resouceName+".class")
		f, err = os.Stat(p)
		if err == nil && f.IsDir() == false { // class file
			realpaths = append(realpaths, &Resource{
				kind:     RESOURCE_KIND_LUCY_CLASS,
				realpath: p,
				name:     resouceName,
			})
		}
	}
	for _, v := range compiler.ClassPaths {
		p := filepath.Join(v, resouceName)
		f, err := os.Stat(p)
		if err == nil && f.IsDir() { // directory is package
			realpaths = append(realpaths, &Resource{
				kind:     RESOURCE_KIND_JAVA_PACKAGE,
				realpath: p,
				name:     resouceName,
			})
		}
		p = filepath.Join(v, resouceName+".class")
		f, err = os.Stat(p)
		if err == nil && f.IsDir() == false { // directory is package
			realpaths = append(realpaths, &Resource{
				kind:     RESOURCE_KIND_JAVA_CLASS,
				realpath: p,
				name:     resouceName,
			})
		}
	}
	if len(realpaths) == 0 {
		return nil, fmt.Errorf("resource '%v' not found", resouceName)
	}
	realpathMap := make(map[string][]*Resource)
	for _, v := range realpaths {
		_, ok := realpathMap[v.realpath]
		if ok {
			realpathMap[v.realpath] = append(realpathMap[v.realpath], v)
		} else {
			realpathMap[v.realpath] = []*Resource{v}
		}
	}
	if len(realpathMap) > 1 {
		errMsg := "not 1 resource named '" + resouceName + "' present:\n"
		for _, v := range realpathMap {
			switch v[0].kind {
			case RESOURCE_KIND_JAVA_CLASS:
				errMsg += fmt.Sprintf("\t in '%s' is a java class\n", v[0].realpath)
			case RESOURCE_KIND_JAVA_PACKAGE:
				errMsg += fmt.Sprintf("\t in '%s' is a java package\n", v[0].realpath)
			case RESOURCE_KIND_LUCY_CLASS:
				errMsg += fmt.Sprintf("\t in '%s' is a lucy class\n", v[0].realpath)
			case RESOURCE_KIND_LUCY_PACKAGE:
				errMsg += fmt.Sprintf("\t in '%s' is a lucy package\n", v[0].realpath)
			}
		}
		return nil, fmt.Errorf(errMsg)
	}
	if realpaths[0].kind == RESOURCE_KIND_LUCY_CLASS {
		if filepath.Base(realpaths[0].realpath) == mainClassName {
			return nil, fmt.Errorf("%s is special class for global variable and other things", mainClassName)
		}
	}
	if realpaths[0].kind == RESOURCE_KIND_JAVA_CLASS {
		class, err := loader.loadClass(realpaths[0])
		if class != nil {
			loader.caches[resouceName] = class
		}
		return class, err
	} else if realpaths[0].kind == RESOURCE_KIND_LUCY_CLASS {
		t, err := loader.loadClass(realpaths[0])
		if t != nil {
			loader.caches[resouceName] = t
		}
		return t, err
	} else if realpaths[0].kind == RESOURCE_KIND_JAVA_PACKAGE {
		p, err := loader.loadJavaPackage(realpaths[0])
		if p != nil {
			loader.caches[resouceName] = p
		}
		return p, err
	} else { // lucy package
		p, err := loader.loadLucyPackage(realpaths[0])
		if p != nil {
			loader.caches[resouceName] = p
		}
		return p, err
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
	mkEnums := func(e *ast.Enum) {
		if p.Block.Enums == nil {
			p.Block.Enums = make(map[string]*ast.Enum)
		}
		if p.Block.EnumNames == nil {
			p.Block.EnumNames = make(map[string]*ast.EnumName)
		}
		p.Block.Enums[filepath.Base(e.Name)] = e
		for _, v := range e.Enums {
			p.Block.EnumNames[v.Name] = v
		}
	}
	for _, v := range fisM {
		bs, err := ioutil.ReadFile(filepath.Join(r.realpath, v.Name()))
		if err != nil {
			return p, fmt.Errorf("read class failed,err:%v", err)
		}
		c, err := (&ClassDecoder{}).decode(bs)
		if err != nil {
			return nil, fmt.Errorf("decode class failed,err:%v", err)
		}
		if len(c.AttributeGroupedByName.GetByName(cg.ATTRIBUTE_NAME_LUCY_ENUM)) > 0 {
			e, err := loader.loadLucyEnum(c)
			if err != nil {
				return nil, err
			}
			mkEnums(e)
			continue
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
		rr.kind = RESOURCE_KIND_JAVA_CLASS
		if strings.HasSuffix(v.Name(), ".class") == false || strings.Contains(v.Name(), "$") {
			continue
		}
		rr.realpath = filepath.Join(r.realpath, v.Name())
		class, err := loader.loadClass(&rr)
		if err != nil {
			return nil, err
		}
		if c, ok := class.(*ast.Class); ok && class != nil {
			ret.Block.Classes[filepath.Base(c.Name)] = c
		}
	}
	return ret, nil
}

func (loader *RealNameLoader) loadClass(r *Resource) (interface{}, error) {
	bs, err := ioutil.ReadFile(r.realpath)
	if err != nil {
		return nil, err
	}
	c, err := (&ClassDecoder{}).decode(bs)
	if r.kind == RESOURCE_KIND_LUCY_CLASS {
		if t := c.AttributeGroupedByName[cg.ATTRIBUTE_NAME_LUCY_ENUM]; len(t) > 0 {
			return loader.loadLucyEnum(c)
		} else {
			return loader.loadAsLucy(c)
		}
	}
	return loader.loadAsJava(c)
}

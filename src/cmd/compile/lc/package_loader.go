package lc

import (
	"encoding/json"
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm"
	"github.com/756445638/lucy/src/cmd/compile/jvm/class_json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type PackageLoader struct {
	P             ast.Package
	name          string
	mainClassName string
}

/*
	类名lucy/lang/Print
在lucy/lang 文件夹中应该有个class为Lang为主类

*/
func (p *PackageLoader) mkMainClassName(name string) {
	n := name
	if strings.Contains(name, "/") {
		t := strings.Split(name, "/")
		n = t[len(t)-1]
	}
	p.mainClassName = strings.Title(n)
}

func (p *PackageLoader) load(realpath string, name string, jsons []os.FileInfo) (*ast.Package, error) {
	p.mkMainClassName(name)
	var bs []byte
	var err error
	for _, v := range jsons {
		bs, err = ioutil.ReadFile(filepath.Join(realpath, v.Name()))
		if err != nil {
			continue
		}
		j := &class_json.ClassJson{}
		json.Unmarshal(bs, j)
		if err != nil { // may be anoth kind of json files
			continue
		}
		if strings.HasSuffix(j.SourceFile, ".java") {
			p.loadAsJava(j)
		} else if strings.HasSuffix(j.SourceFile, ".lucy") {
			p.loadAsLucy(j)
		} else {
			fmt.Printf("unsupported json file,compile from %v", j.SourceFile)
		}
	}
	return nil, nil
}

func (p *PackageLoader) loadAsLucy(j *class_json.ClassJson) {
	shortname := ""
	{
		t := strings.Split(j.ThisClass, "/")
		shortname = t[len(t)-1]
	}
	mainclass := shortname == p.mainClassName // if match
	if mainclass {

	}

}

func (p *PackageLoader) loadAsJava(j *class_json.ClassJson) {
	c := &ast.Class{}
	c.Fields = make(map[string]*ast.ClassField)
	c.Name = j.ThisClass
	c.SuperClassName = j.SuperClass
	for _, v := range j.Fields {
		t, _ := jvm.ParseType(v.Descriptor)
		f := &ast.ClassField{}
		f.ClassFieldProperty.AccessFlags = v.AccessFlags
		f.Name = v.Name
		f.Signature = v.Signature
		f.Typ = t
		c.Fields[v.Name] = f
	}
	shortname := ""
	{
		t := strings.Split(c.Name, "/")
		shortname = t[len(t)-1]
	}
	c.Methods = make(map[string][]*ast.ClassMethod)
	c.Constructors = []*ast.ClassMethod{}
	for _, v := range j.Methods {
		m := &ast.ClassMethod{}
		m.AccessFlags = v.AccessFlags
		m.Signature = v.Signature
		m.Func = &ast.Function{}
		m.Func.Name = v.Name
		m.Func.Typ = &ast.FunctionType{}
		m.Func.Typ.Parameters = make([]*ast.VariableDefinition, len(v.Typ.Parameters))
		for kk, vv := range v.Typ.Parameters {
			t, _ := jvm.ParseType(vv)
			m.Func.Typ.Parameters[kk] = &ast.VariableDefinition{}
			m.Func.Typ.Parameters[kk].Typ = t
		}
		m.Func.Typ.Returns = []*ast.VariableDefinition{&ast.VariableDefinition{}}
		m.Func.Typ.Returns[0].Typ, _ = jvm.ParseType(v.Typ.Return)
		if v.Name == shortname {
			c.Constructors = append(c.Constructors, m)
			continue
		}
		if _, ok := c.Methods[v.Name]; !ok {
			c.Methods[v.Name] = []*ast.ClassMethod{}
		}
		c.Methods[v.Name] = append(c.Methods[v.Name], m)
	}

}

func (*PackageLoader) LoadPackage(name string) (*ast.Package, error) {
	if len(l.lucyPath) == 0 {
		return nil, fmt.Errorf("no env variable LUCYPATH found")
	}
	packagename := name
	if strings.Contains(name, "/") {
		t := strings.Split(name, "/")
		packagename = t[len(t)-1]
		if packagename == "" {
			return nil, fmt.Errorf("no name after separator '/'")
		}
	}
	var realpath string
	for _, v := range l.lucyPath {
		_, err := os.Stat(filepath.Join(v, name))
		if err == nil {
			realpath = filepath.Join(v, name)
			break
		}
	}

	if realpath == "" {
		return nil, fmt.Errorf("package %v not found")
	}
	fis, err := ioutil.ReadDir(realpath)
	if err != nil {
		return nil, fmt.Errorf("read dir %v failed,err:%v\n", err)
	}
	jsonfiles := []os.FileInfo{}
	for _, v := range fis {
		if strings.HasSuffix(v.Name(), ".json") {
			jsonfiles = append(jsonfiles, v)
		}
	}
	if len(jsonfiles) == 0 {
		return nil, fmt.Errorf("package %s has not  been compiled yet")
	}
	return (&PackageLoader{}).load(realpath, packagename, jsonfiles)
}

func init() {
	ast.PackageLoad = &PackageLoader{}
}

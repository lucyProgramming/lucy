package lc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/class_json"
)

type PackageLoader struct {
	//	realpath string
	P    ast.Package
	name string
}

func (p *PackageLoader) load(realpath string, name string, jsons []os.FileInfo) (*ast.Package, error) {
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
			fmt.Println("unsupported json file,compile from %v", j.SourceFile)
		}
	}
	return nil, nil
}

func (p *PackageLoader) loadAsJava(j *class_json.ClassJson) {
	//	c := &ast.Class{}

}
func (p *PackageLoader) loadAsLucy(j *class_json.ClassJson) {

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

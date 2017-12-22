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
		shortname = strings.Title(shortname)
	}
	mainclass := shortname == p.mainClassName // if match
	if mainclass {
		p.loadMainLucy(j)
		return
	}
	// load regular class
	c := &ast.Class{} // rebuild the name,trim prefix
	{
		t := strings.Split(j.ThisClass, "/")
		t[len(t)-1] = strings.TrimLeft(t[len(t)-1], p.mainClassName)
		c.Name = strings.Join(t, "/")
	}
	c.Fields = make(map[string]*ast.ClassField)
	for _, v := range j.Fields {
		f := &ast.ClassField{}
		f.VariableDefinition = *p.loadFieldAsVariableDefination(v)
		c.Fields[v.Name] = f
	}
	c.Methods = make(map[string][]*ast.ClassMethod)
	for _, v := range j.Methods {
		f := p.loadMethod(v)
		m := &ast.ClassMethod{}
		m.Func = f
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

/*
	main class
	比如包名为 lucy/lang/xxx/
	main class wei Xxx
	main class cannot
*/
func (p *PackageLoader) loadMainLucy(j *class_json.ClassJson) {
	if p.P.Block.Vars == nil {
		p.P.Block.Vars = make(map[string]*ast.VariableDefinition)
	}
	for _, v := range j.Fields {
		p.P.Block.Vars[v.Name] = p.loadFieldAsVariableDefination(v)
	}
	if p.P.Block.Funcs == nil {
		p.P.Block.Funcs = make(map[string]*ast.Function)
	}
	for _, v := range j.Methods {
		f := p.loadMethod(v)
		p.P.Block.Funcs[v.Name] = f
	}
}

func (p *PackageLoader) loadFieldAsVariableDefination(field *class_json.Field) *ast.VariableDefinition {
	v := &ast.VariableDefinition{}
	v.Name = field.Name
	v.AccessFlags = field.AccessFlags
	v.Typ, _ = jvm.ParseType(field.Descriptor)
	v.Signature = field.Signature
	return v
}

func (p *PackageLoader) loadMethod(m *class_json.Method) *ast.Function {
	f := &ast.Function{}
	f.Typ = &ast.FunctionType{}
	f.AccessFlags = m.AccessFlags
	f.Signature = m.Signature
	t, _ := jvm.ParseType(m.Typ.Return)
	f.Typ.Returns = []*ast.VariableDefinition{
		{},
	}
	f.Typ.Returns[0].Typ = t
	f.Typ.Parameters = make([]*ast.VariableDefinition, len(f.Typ.Parameters))
	for k, v := range m.Typ.Parameters {
		vd := &ast.VariableDefinition{}
		vd.Typ, _ = jvm.ParseType(v)
		f.Typ.Parameters[k] = vd
	}

	return f
}

func (p *PackageLoader) loadAsJava(j *class_json.ClassJson) {
	c := &ast.Class{}
	c.Fields = make(map[string]*ast.ClassField)
	c.Name = j.ThisClass
	c.SuperClassName = j.SuperClass
	c.SouceFile = j.SourceFile
	c.Signature = j.Signature
	for _, v := range j.Fields {
		t, _ := jvm.ParseType(v.Descriptor)
		f := &ast.ClassField{}
		f.AccessFlags = v.AccessFlags
		f.Name = v.Name
		f.Typ = t
		f.Signature = v.Signature
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
		m.Func = p.loadMethod(v)
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

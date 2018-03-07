package lc

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type PackageLoader struct {
	Package ast.Package
}

func (*PackageLoader) LoadPackage(name string) (*ast.Package, error) {
	if len(compiler.lucyPath) == 0 {
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
	var realpath []string
	for _, v := range compiler.lucyPath {
		_, err := os.Stat(filepath.Join(v, name))
		if err == nil {
			realpath = append(realpath, filepath.Join(v, name))
			break
		}
	}
	if len(realpath) == 0 {
		return nil, fmt.Errorf("package '%v' not found", name)
	}
	if len(realpath) > 1 {
		return nil, fmt.Errorf("not 1 package name '%s' present in $CLASSPATH", name)
	}
	fis, err := ioutil.ReadDir(realpath[0])
	if err != nil {
		return nil, fmt.Errorf(" dir %v failed,err:%v\n", err)
	}

	classfiles := []os.FileInfo{}
	for _, v := range fis {
		if strings.HasSuffix(v.Name(), ".class") {
			classfiles = append(classfiles, v)
		}
	}
	if len(classfiles) == 0 {
		return nil, fmt.Errorf("package '%s' has no class files")
	}
	return (&PackageLoader{}).load(realpath[0], packagename, classfiles)
}

func (p *PackageLoader) load(realpath string, name string, files []os.FileInfo) (*ast.Package, error) {
	return nil, fmt.Errorf(".......")
}

//
//func (p *PackageLoader) loadAsLucy(j *class_json.ClassJson) {
//	shortname := ""
//	{
//		t := strings.Split(j.ThisClass, "/")
//		shortname = t[len(t)-1]
//		shortname = strings.Title(shortname)
//	}
//	mainclass := shortname == p.mainClassName // if match
//	if mainclass {
//		p.loadMainLucy(j)
//		return
//	}
//	//load regular class
//	c := &ast.Class{} // rebuild the name,trim prefix
//	{
//		t := strings.Split(j.ThisClass, "/")
//		t[len(t)-1] = strings.TrimLeft(t[len(t)-1], p.mainClassName)
//		c.ClassNameDefinition.Name = strings.Join(t, "/")
//	}
//	c.Fields = make(map[string]*ast.ClassField)
//	for _, v := range j.Fields {
//		f := &ast.ClassField{}
//		f.VariableDefinition = *p.loadFieldAsVariableDefination(v)
//		c.Fields[v.Name] = f
//	}
//	c.Methods = make(map[string][]*ast.ClassMethod)
//	for _, v := range j.Methods {
//		f := p.loadMethod(v)
//		m := &ast.ClassMethod{}
//		m.Func = f
//		if v.Name == shortname {
//			c.Constructors = append(c.Constructors, m)
//			continue
//		}
//		if _, ok := c.Methods[v.Name]; !ok {
//			c.Methods[v.Name] = []*ast.ClassMethod{}
//		}
//		c.Methods[v.Name] = append(c.Methods[v.Name], m)
//	}
//}

/*
	main class
	比如包名为 lucy/lang/xxx
	main class wei Xxx
	main class cannot
*/
//func (p *PackageLoader) loadMainLucy(j *class_json.ClassJson) {
//	if p.P.Block.Vars == nil {
//		p.P.Block.Vars = make(map[string]*ast.VariableDefinition)
//	}
//	for _, v := range j.Fields {
//		p.P.Block.Vars[v.Name] = p.loadFieldAsVariableDefination(v)
//	}
//	if p.P.Block.Funcs == nil {
//		p.P.Block.Funcs = make(map[string]*ast.Function)
//	}
//	for _, v := range j.Methods {
//		f := p.loadMethod(v)
//		p.P.Block.Funcs[v.Name] = f
//	}
//}

//func (p *PackageLoader) loadFieldAsVariableDefination(field *class_json.Field) *ast.VariableDefinition {
//	v := &ast.VariableDefinition{}
//	v.Name = field.Name
//	v.AccessFlags = field.AccessFlags
//	v.Typ, _ = jvm.ParseType(field.Descriptor)
//	v.Signature = field.Signature
//	return v
//}
//
//func (p *PackageLoader) loadMethod(m *class_json.Method) *ast.Function {
//	f := &ast.Function{}
//	f.Typ = &ast.FunctionType{}
//	f.AccessFlags = m.AccessFlags
//	f.Signature = m.Signature
//	t, _ := jvm.ParseType(m.Typ.Return)
//	f.Typ.ReturnList = []*ast.VariableDefinition{
//		{},
//	}
//	f.Typ.ReturnList[0].Typ = t
//	f.Typ.ParameterList = make([]*ast.VariableDefinition, len(f.Typ.ParameterList))
//	for k, v := range m.Typ.Parameters {
//		vd := &ast.VariableDefinition{}
//		vd.Typ, _ = jvm.ParseType(v)
//		f.Typ.ParameterList[k] = vd
//	}
//	return f
//}
//
//func (p *PackageLoader) loadAsJava(j *class_json.ClassJson) {
//	c := &ast.Class{}
//	c.Fields = make(map[string]*ast.ClassField)
//	c.ClassNameDefinition.Name = j.ThisClass
//	c.SuperClassNameDefinition.BinaryName = j.SuperClass
//	c.SouceFile = j.SourceFile
//	c.Signature = j.Signature
//	for _, v := range j.Fields {
//		t, _ := jvm.ParseType(v.Descriptor)
//		f := &ast.ClassField{}
//		f.AccessFlags = v.AccessFlags
//		f.Name = v.Name
//		f.Typ = t
//		f.Signature = v.Signature
//		c.Fields[v.Name] = f
//	}
//	shortname := ""
//	{
//		t := strings.Split(c.ClassNameDefinition.Name, "/")
//		shortname = t[len(t)-1]
//	}
//	c.Methods = make(map[string][]*ast.ClassMethod)
//	c.Constructors = []*ast.ClassMethod{}
//	for _, v := range j.Methods {
//		m := &ast.ClassMethod{}
//		m.Func = p.loadMethod(v)
//		if v.Name == shortname {
//			c.Constructors = append(c.Constructors, m)
//			continue
//		}
//		if _, ok := c.Methods[v.Name]; !ok {
//			c.Methods[v.Name] = []*ast.ClassMethod{}
//		}
//		c.Methods[v.Name] = append(c.Methods[v.Name], m)
//	}
//}
///*
//	类名lucy/lang/Print
//	在lucy/lang 文件夹中应该有个class为Lang为主类
//*/
//func (p *PackageLoader) mkMainClassName(name string) {
//	n := name
//	if strings.Contains(name, "/") {
//		t := strings.Split(name, "/")
//		n = t[len(t)-1]
//	}
//	p.mainClassName = strings.Title(n)
//}

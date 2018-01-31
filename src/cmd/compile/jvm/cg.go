package jvm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type MakeClass struct {
	p              *ast.Package
	Classes        []*cg.ClassHighLevel
	mainclass      *cg.ClassHighLevel
	MakeExpression MakeExpression
}

func (m *MakeClass) Make(p *ast.Package) {
	m.p = p
	mainclass := &cg.ClassHighLevel{}
	m.mainclass = mainclass
	//	mainclass.AccessFlags |= cg.ACC_CLASS_PUBLIC
	//	mainclass.AccessFlags |= cg.ACC_CLASS_FINAL
	mainclass.SuperClass = ast.JAVA_ROOT_CLASS
	if p.Name == "" {
		p.Name = "test"
	}
	mainclass.Name = p.Name
	mainclass.Fields = make(map[string]*cg.FiledHighLevel)
	m.mkVars()
	m.mkEnums()
	m.mkClass()
	m.mkFuncs()
	m.mkInitFunctions()
	err := m.Dump()
	if err != nil {
		panic(fmt.Sprintf("dump to file failed,err:%v\n", err))
	}
}

func (m *MakeClass) Dump() error {
	//dump main class
	f, err := os.OpenFile("test.class", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if err := cg.FromHighLevel(m.mainclass).OutPut(f); err != nil {
		f.Close()
		return err
	}
	f.Close()
	for _, c := range m.Classes {
		f, err = os.OpenFile(filepath.Join(m.p.DestPath, c.Name+".class"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		if err = cg.FromHighLevel(c).OutPut(f); err != nil {
			f.Close()
			return err
		}
	}
	return nil
}

func (m *MakeClass) mkVars() {
	for k, v := range m.p.Block.Vars {
		f := &cg.FiledHighLevel{}
		f.AccessFlags = v.AccessFlags
		f.Descriptor = v.Typ.Descriptor()
		m.mainclass.Fields[k] = f
	}
}

func (m *MakeClass) mkInitFunctions() {
	//	ms := []*cg.MethodHighLevel{}
	//	for k, v := range m.p.InitFunctions {
	//		method := &cg.MethodHighLevel{}
	//		ms = append(ms, method)
	//		method.AccessFlags |= cg.ACC_METHOD_STATIC
	//		method.AccessFlags |= cg.ACC_METHOD_FINAL
	//		method.AccessFlags |= cg.ACC_METHOD_PRIVATE
	//		method.Name = fmt.Sprintf("block%d", k)
	//		method.Class = m.mainclass
	//		method.Descriptor = "()V"
	//		context := &Context{v, nil}
	//		m.buildFunction(m.mainclass, method, v, context)
	//		fmt.Println(method.Code)
	//	}

	//	method := &cg.MethodHighLevel{}
	//	m.buildEntryMethod(method, ms, true)
	//	m.mainclass.AppendMethod(method)
	//	method2 := &cg.MethodHighLevel{}
	//	m.buildEntryMethod(method2, ms, false)
	//	m.mainclass.AppendMethod(method2)
	//	m.mainclass.AppendMethod(ms...)
}

func (m *MakeClass) mkEnums() {

}
func (m *MakeClass) mkClass() {

}

func (m *MakeClass) mkFuncs() {
	for _, f := range m.p.Block.Funcs {
		if f.Isbuildin { //
			continue
		}
		m.mkFunc(f)
	}
}

func (m *MakeClass) mkClosureFunctionClass() *cg.ClassHighLevel {
	ret := &cg.ClassHighLevel{}
	ret.AccessFlags = cg.ACC_CLASS_FINAL
	return ret
}

func (m *MakeClass) buildBlock(class *cg.ClassHighLevel, code *cg.AttributeCode, b *ast.Block, context *Context) {
	var maxstack uint16
	for _, s := range b.Statements {
		maxstack = m.buildStatement(class, code, s, context)
		if maxstack > code.MaxStack {
			code.MaxStack = maxstack
		}
	}
	return
}

func (m *MakeClass) mkFuncClassMode(f *ast.Function) *cg.MethodHighLevel {
	ret := &cg.MethodHighLevel{}
	return ret
}

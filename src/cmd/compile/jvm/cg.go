package jvm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
)

type MakeClass struct {
	Descriptor     *Descriptor
	p              *ast.Package
	Classes        []*cg.ClassHighLevel
	mainclass      *cg.ClassHighLevel
	MakeExpression MakeExpression
}

func (m *MakeClass) Make(p *ast.Package) {
	m.p = p
	mainclass := &cg.ClassHighLevel{}
	m.mainclass = mainclass
	mainclass.AccessFlags |= cg.ACC_CLASS_FINAL
	mainclass.SuperClass = ast.JAVA_ROOT_CLASS
	if p.Name == "" {
		p.Name = "test"
	}
	mainclass.Name = p.Name
	mainclass.Fields = make(map[string]*cg.FiledHighLevel)
	mkClassDefaultContruction(m.mainclass)
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
		f.AccessFlags = 0
		f.AccessFlags |= cg.ACC_FIELD_STATIC
		f.AccessFlags |= cg.ACC_FIELD_SYNTHETIC
		f.Descriptor = m.Descriptor.typeDescriptor(v.Typ)
		if v.AccessFlags&cg.ACC_FIELD_PUBLIC != 0 {
			f.AccessFlags |= cg.ACC_FIELD_PUBLIC
		}
		f.Name = v.Name
		m.mainclass.Fields[k] = f
	}
}

func (m *MakeClass) mkInitFunctions() {
	ms := []*cg.MethodHighLevel{}
	for _, v := range m.p.InitFunctions {
		method := &cg.MethodHighLevel{}
		ms = append(ms, method)
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		method.AccessFlags |= cg.ACC_METHOD_FINAL
		method.AccessFlags |= cg.ACC_METHOD_PRIVATE
		method.Name = m.mainclass.NewFunctionName("block")
		method.Class = m.mainclass
		method.Descriptor = "()V"
		m.buildFunction(m.mainclass, method, v)
		m.mainclass.AppendMethod(method)
	}
	if len(ms) == 0 {
		return
	}
	method := &cg.MethodHighLevel{}
	method.AccessFlags |= cg.ACC_METHOD_STATIC
	method.Name = "<clinit>"
	method.Descriptor = "()V"
	codes := make([]byte, 65536)
	codelength := uint16(0)
	for _, v := range ms {
		codes[codelength] = cg.OP_invokestatic
		m.mainclass.InsertMethodRefConst(cg.CONSTANT_Methodref_info_high_level{
			Class:      m.mainclass.Name,
			Name:       v.Name,
			Descriptor: "()V",
		}, codes[codelength+1:codelength+3])
		codelength += 3
	}
	codes[codelength] = cg.OP_return
	codelength++
	codes = codes[0:codelength]
	method.Code.Codes = codes
	method.Code.CodeLength = codelength
	m.mainclass.AppendMethod(method)
}

func (m *MakeClass) mkEnums() {

}
func (m *MakeClass) mkClass() {

}

func (m *MakeClass) mkFuncs() {
	ms := make(map[string]*cg.MethodHighLevel)
	for k, f := range m.p.Block.Funcs { // fisrt round
		if f.Isbuildin { //
			continue
		}
		method := &cg.MethodHighLevel{}
		method.Class = m.mainclass
		method.Name = f.Name
		method.Descriptor = m.Descriptor.methodDescriptor(f)
		method.AccessFlags = 0
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		if f.AccessFlags&cg.ACC_METHOD_PUBLIC != 0 || f.Name == ast.MAIN_FUNCTION_NAME {
			method.AccessFlags |= cg.ACC_METHOD_PUBLIC
		}
		ms[k] = method
		f.ClassMethod = method
		m.mainclass.AppendMethod(method)
	}
	for k, f := range m.p.Block.Funcs { // fisrt round
		if f.Isbuildin { //
			continue
		}
		m.buildFunction(ms[k].Class, ms[k], f)
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

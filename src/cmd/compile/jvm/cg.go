package jvm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	mainclass.AccessFlags |= cg.ACC_CLASS_PUBLIC
	mainclass.AccessFlags |= cg.ACC_CLASS_FINAL
	mainclass.AccessFlags |= cg.ACC_CLASS_ABSTRACT
	mainclass.AccessFlags |= cg.ACC_CLASS_SYNTHETIC
	mainclass.SuperClass = ast.LUCY_ROOT_CLASS
	mainclass.Name = strings.Title(p.Name)
	mainclass.Fields = make(map[string]*cg.FiledHighLevel)
	mainclass.Methods = make(map[string][]*cg.MethodHighLevel)
	m.mkVars()
	m.mkEnums()
	m.mkClass()
	m.mkFuncs()

	m.mkInitFunctions()

	err := m.Dump()
	if err != nil {
		panic(fmt.Sprintf("dump to file failed,err:%v\n", err))
	}
	//dump main class

}

func (m *MakeClass) Dump() error {
	//dump main class
	f, err := os.OpenFile(filepath.Join(m.p.DestPath, m.mainclass.Name+".class"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
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
	ms := []*cg.MethodHighLevel{}
	for k, v := range m.p.InitFunctions {
		method := &cg.MethodHighLevel{}
		ms = append(ms, method)
		method.AccessFlags |= cg.ACC_METHOD_STATIC
		method.AccessFlags |= cg.ACC_METHOD_FINAL
		method.AccessFlags |= cg.ACC_METHOD_BRIDGE
		method.Name = fmt.Sprintf("<block%d>", k)
		method.Class = m.mainclass
		method.Descriptor = "()V"
		context := &Context{v, nil}
		m.buildFunction(m.mainclass, method, v, context)
	}
	// mk main function
	method := &cg.MethodHighLevel{}
	m.buildEntryMethod(method, ms, true)
	m.mainclass.AppendMethod(method)
	method2 := &cg.MethodHighLevel{}
	m.buildEntryMethod(method2, ms, false)
	m.mainclass.AppendMethod(method2)
	m.mainclass.AppendMethod(ms...)
}

func (m *MakeClass) buildEntryMethod(method *cg.MethodHighLevel, ms []*cg.MethodHighLevel, ismain bool) {
	if ismain {
		method.Name = "main"
	} else {
		method.Name = "<bloks>"
	}
	method.AccessFlags |= cg.ACC_METHOD_PUBLIC
	method.AccessFlags |= cg.ACC_METHOD_STATIC
	method.AccessFlags |= cg.ACC_METHOD_SYNTHETIC
	method.Descriptor = "()V"
	method.Class = m.mainclass
	method.Code.Codes = make([]byte, 65536)
	defer func() {
		method.Code.Codes = method.Code.Codes[0:method.Code.CodeLength]
	}()
	if ismain {
		method.Code.Codes[method.Code.CodeLength] = cg.OP_iconst_1
		method.Code.Codes[method.Code.CodeLength+1] = cg.OP_putstatic
		m.mainclass.InsertFieldRef(cg.CONSTANT_Fieldref_info_high_level{
			Class:      m.mainclass.Name,
			Name:       "__main__",
			Descriptor: "I",
		}, method.Code.Codes[method.Code.CodeLength+2:method.Code.CodeLength+4])
		method.Code.CodeLength += 4
	}
	for _, v := range ms {
		method.Code.Codes[method.Code.CodeLength] = cg.OP_invokestatic
		m.mainclass.InsertMethodRef(cg.CONSTANT_Methodref_info_high_level{
			Class:      m.mainclass.Name,
			Name:       v.Name,
			Descriptor: "()V",
		}, method.Code.Codes[method.Code.CodeLength+1:method.Code.CodeLength+3])
	}
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
	for _, s := range b.Statements {
		m.buildStatement(class, code, s, context)
	}
	return
}

func (m *MakeClass) mkFuncClassMode(f *ast.Function) *cg.MethodHighLevel {
	ret := &cg.MethodHighLevel{}
	return ret
}

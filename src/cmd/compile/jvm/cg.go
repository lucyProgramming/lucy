package jvm

import (
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
	mainclass.SuperClass = ast.LUCY_ROOT_CLASS
	mainclass.Name = strings.Title(p.Name)
	mainclass.Fields = make(map[string]*cg.FiledHighLevel)
	mainclass.Methods = make(map[string][]*cg.MethodHighLevel)
	m.mkVars()
	m.mkConsts()
	m.mkEnums()
	m.mkClass()
	m.mkFuncs()
	m.mkBlocks()
}

func (m *MakeClass) mkVars() {
	for k, v := range m.p.Block.Vars {
		f := &cg.FiledHighLevel{}
		f.AccessFlags = v.AccessFlags
		f.Descriptor = v.Typ.Descriptor()
		m.mainclass.Fields[k] = f
	}
}
func (m *MakeClass) mkConsts() {

}
func (m *MakeClass) mkBlocks() {

}
func (m *MakeClass) mkEnums() {

}
func (m *MakeClass) mkClass() {

}

func (m *MakeClass) mkFuncs() {
	for _, f := range m.p.Block.Funcs {
		if f.Isbuildin {
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

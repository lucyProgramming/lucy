package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
	"strings"
)

type MakeClass struct {
	p         *ast.Package
	Classes   []*cg.ClassHighLevel
	mainclass *cg.ClassHighLevel
}

func (m *MakeClass) Make(p *ast.Package) {
	m.p = p
	mainclass := &cg.ClassHighLevel{}
	m.mainclass = mainclass
	mainclass.AccessFlags |= cg.ACC_CLASS_PUBLIC
	mainclass.AccessFlags |= cg.ACC_CLASS_FINAL
	mainclass.AccessFlags |= cg.ACC_CLASS_ABSTRACT
	mainclass.SuperClass = ast.JAVA_ROOT_CLASS
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
		m.mkFunc(f, "")
	}
}

func (m *MakeClass) mkFunc(f *ast.Function, path string) {
	if f.IsGlobal || f.Typ.ClosureVars == nil || len(f.Typ.ClosureVars) == 0 {
		m.mkFuncStaticMethodMode(f, path)
	} else {
		m.mkFuncClassMode(f, path)
	}
}

func (m *MakeClass) mkFuncStaticMethodMode(f *ast.Function, path string) *cg.MethodHighLevel {
	ret := &cg.MethodHighLevel{}
	ret.AccessFlags = 0
	if f.AccessFlags&cg.ACC_METHOD_PUBLIC != 0 {
		ret.AccessFlags |= cg.ACC_METHOD_PUBLIC
	}
	ret.AccessFlags |= cg.ACC_METHOD_STATIC
	ret.AccessFlags |= cg.ACC_METHOD_FINAL
	ret.Name = mkPath(path, f.Name)
	ret.Code.MaxLocals, ret.Code.MaxStack = m.buildBlock(&ret.Code, f.Block, ret.Name)
	ret.Descriptor = f.Descriptor
	return ret
}

func (m *MakeClass) buildBlock(code *cg.AttributeCode, b *ast.Block, path string) (locals uint16, stack uint16) {

	return
}

func (m *MakeClass) buildStatement(code *cg.AttributeCode, s *ast.Statement, offset uint16) (locals uint16, stack uint16) {
	switch s.Typ {
	case ast.STATEMENT_TYPE_EXPRESSION:

	case ast.STATEMENT_TYPE_IF:

	case ast.STATEMENT_TYPE_BLOCK:
	case ast.STATEMENT_TYPE_FOR:
	case ast.STATEMENT_TYPE_CONTINUE:
	case ast.STATEMENT_TYPE_RETURN:
	case ast.STATEMENT_TYPE_BREAK:
	case ast.STATEMENT_TYPE_SWITCH:
	case ast.STATEMENT_TYPE_SKIP: // skip this block

	}
	return
}

func (m *MakeClass) mkFuncClassMode(f *ast.Function, path string) *cg.MethodHighLevel {
	ret := &cg.MethodHighLevel{}
	return ret
}

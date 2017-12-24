package jvm

import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"github.com/756445638/lucy/src/cmd/compile/jvm/cg"
	"strings"
)

type MakeClass struct {
	p       *ast.Package
	Classes []*cg.ClassHighLevel
}

func (m *MakeClass) Make(p *ast.Package) []error {
	errs := []error{}
	mainclass := &cg.ClassHighLevel{}
	mainclass.AccessFlags |= cg.ACC_CLASS_PUBLIC
	mainclass.AccessFlags |= cg.ACC_CLASS_FINAL
	mainclass.AccessFlags |= cg.ACC_CLASS_ABSTRACT
	mainclass.SuperClass = ast.JAVA_ROOT_CLASS
	mainclass.Name = strings.Title(p.Name)
	mainclass.Fields = make(map[string]*cg.FiledHighLevel)
	mainclass.Methods = make(map[string][]*cg.MethodHighLevel)
	for k, v := range p.Block.Vars {
		f := &cg.FiledHighLevel{}
		f.AccessFlags = v.AccessFlags
		f.Descriptor = v.Typ.Descriptor()
		mainclass.Fields[k] = f
	}

	return errs
}

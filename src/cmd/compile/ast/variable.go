package ast

import "gitee.com/yuyang-fine/lucy/src/cmd/compile/jvm/cg"

type Variable struct {
	IsBuildIn                bool
	IsGlobal                 bool
	IsFunctionParameter      bool
	IsReturn                 bool
	BeenCapturedAsLeftValue  int
	BeenCapturedAsRightValue int
	Used                     bool   // use as right value
	AccessFlags              uint16 // public private or protected
	Pos                      *Pos
	DefaultValueExpression   *Expression
	Name                     string
	Type                     *Type
	LocalValOffset           uint16 // offset in stack frame
	JvmDescriptor            string // jvm
	Comment                  string
}

func (this *Variable) isPublic() bool {
	return this.AccessFlags&
		cg.AccFieldPublic != 0
}

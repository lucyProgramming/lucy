package optimizer

import (
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

func Optimize(p *ast.Package) {
	(&Optimizer{p: p}).Optimize()
}

type Optimizer struct {
	p *ast.Package
}

func (o *Optimizer) Optimize() {

}

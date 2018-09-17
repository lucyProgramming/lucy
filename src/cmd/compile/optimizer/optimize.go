package optimizer

//import (
//	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
//)
//
//func Optimize(p *ast.Package) {
//	(&Optimizer{p: p}).Optimize()
//}
//
//type Optimizer struct {
//	p *ast.Package
//}
//
//func (o *Optimizer) Optimize() {
//
//}
//
////func (o *Optimizer) optimizeFunction(f *ast.Function) {
////	for _, v := range f.Typ.ReturnList {
////		if v.Expression.IsCompileAuto {
////			continue
////		}
////		(&Expression{}).optimize(&f.Block, v.Expression)
////	}
////}

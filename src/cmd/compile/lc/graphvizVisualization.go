package lc

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type GraphvizVisualization struct {
}

func (g *GraphvizVisualization) VisualNodes(fileNodes map[string][]*ast.TopNode) {
	var graph = `digraph {`
	for k, v := range fileNodes {
		graph += g.VisualOneFile(k, v)
	}
	graph += "}"
	fmt.Println(graph)
}

func (g *GraphvizVisualization) VisualOneFile(filename string, nodes []*ast.TopNode) (graph string) {
	graph = "subgraph{"
	for _, v := range nodes {
		switch v.Data.(type) {
		//case *ast.Block:
		//	graph += g.VisualBlock(v.Data.(*ast.Block)) + ";"
		//case *ast.Function:
		//	return g.VisualFunction(v.Data.(*ast.Function)) + ";"
		//case *ast.Enum:
		//	return g.VisualEnum(v.Data.(*ast.Enum)) + ";"
		//case *ast.Class:
		//	return g.VisualClass(v.Data.(*ast.Class)) + ";"
		//case *ast.Constant:
		//	return g.VisualConst(v.Data.(*ast.Constant)) + ";"
		//case *ast.Import:
		//	return g.VisualImport(v.Data.(*ast.Import)) + ";"
		case *ast.Expression: // a,b = f();
			graph += g.VisualExpression(v.Data.(*ast.Expression)) + ";"
			//case *ast.TypeAlias:
			//	return g.VisualTypeAlias(v.Data.(*ast.TypeAlias)) + ";"
			//default:
			//	panic("tops have unKnow  type")
		}
	}

	return graph + " } ; "
}
func (g *GraphvizVisualization) VisualImport(nodes *ast.Import) (graph string) {
	return
}
func (g *GraphvizVisualization) VisualFunction(nodes *ast.Function) (graph string) {
	return
}
func (g *GraphvizVisualization) VisualClass(nodes *ast.Class) (graph string) {
	return
}

func (g *GraphvizVisualization) VisualConst(c *ast.Constant) (graph string) {
	return
}

func (g *GraphvizVisualization) VisualStatement(s *ast.Statement) (graph string) {
	return
}
func (g *GraphvizVisualization) VisualEnum(e *ast.Enum) (graph string) {
	return
}
func (g *GraphvizVisualization) VisualTypeAlias(alias *ast.TypeAlias) (graph string) {

	return
}

func (g *GraphvizVisualization) VisualBlock(block *ast.Block) (graph string) {
	return
}

func (g *GraphvizVisualization) VisualPackage(p *ast.Package) {

}

func (g *GraphvizVisualization) VisualExpression(e *ast.Expression) (graph string) {
	switch e.Type {
	case ast.ExpressionTypeNull:
		return "null"
	case ast.ExpressionTypeBool:
		return fmt.Sprintf("%v", e.Data.(bool))
	case ast.ExpressionTypeByte:
		return fmt.Sprintf("%v", e.Data.(byte))
	case ast.ExpressionTypeShort:
		fallthrough
	case ast.ExpressionTypeChar:
		fallthrough
	case ast.ExpressionTypeInt:
		return fmt.Sprintf("%v", e.Data.(int32))
	case ast.ExpressionTypeLong:
		return fmt.Sprintf("%v", e.Data.(int64))
	case ast.ExpressionTypeFloat:
		return fmt.Sprintf("%v", e.Data.(float32))
	case ast.ExpressionTypeDouble:
		return fmt.Sprintf("%v", e.Data.(float64))
	case ast.ExpressionTypeString:
		return fmt.Sprintf("\"%v\"", e.Data.(string))
	case ast.ExpressionTypeArray:
		graph = "{"
		for _, v := range e.Data.([]*ast.Expression) {
			graph += "digraph{" + g.VisualExpression(v) + "} ;"
		}
		graph += `node ["op" = "arrayliteral"]`
		graph += "}"
		return graph
	case ast.ExpressionTypeLogicalOr:
		fallthrough
	case ast.ExpressionTypeLogicalAnd:
		fallthrough
	case ast.ExpressionTypeOr:
		fallthrough
	case ast.ExpressionTypeAnd:
		fallthrough
	case ast.ExpressionTypeXor:
		fallthrough
	case ast.ExpressionTypeLsh:
		fallthrough
	case ast.ExpressionTypeRsh:
		fallthrough
	case ast.ExpressionTypeAdd:
		fallthrough
	case ast.ExpressionTypeSub:
		fallthrough
	case ast.ExpressionTypeMul:
		fallthrough
	case ast.ExpressionTypeDiv:
		fallthrough
	case ast.ExpressionTypeMod:
		fallthrough
	case ast.ExpressionTypeAssign:
		fallthrough
	case ast.ExpressionTypeVarAssign:
		fallthrough
	case ast.ExpressionTypePlusAssign:
		fallthrough
	case ast.ExpressionTypeMinusAssign:
		fallthrough
	case ast.ExpressionTypeMulAssign:
		fallthrough
	case ast.ExpressionTypeDivAssign:
		fallthrough
	case ast.ExpressionTypeModAssign:
		fallthrough
	case ast.ExpressionTypeAndAssign:
		fallthrough
	case ast.ExpressionTypeOrAssign:
		fallthrough
	case ast.ExpressionTypeXorAssign:
		fallthrough
	case ast.ExpressionTypeLshAssign:
		fallthrough
	case ast.ExpressionTypeRshAssign:
		fallthrough
	case ast.ExpressionTypeEq:
		fallthrough
	case ast.ExpressionTypeNe:
		fallthrough
	case ast.ExpressionTypeGe:
		fallthrough
	case ast.ExpressionTypeGt:
		fallthrough
	case ast.ExpressionTypeLe:
		fallthrough
	case ast.ExpressionTypeLt:
		bin := e.Data.(*ast.ExpressionBinary)
		graph += "subgraph {" + g.VisualExpression(bin.Left) + "};"
		graph += fmt.Sprintf(`node ["op" = "%s"] ; `, e.Description)
		graph += "subgraph {" + g.VisualExpression(bin.Right) + "} "
		return graph
	case ast.ExpressionTypeIndex:

	case ast.ExpressionTypeSelection:
	case ast.ExpressionTypeSelectionConst:
	case ast.ExpressionTypeMethodCall:
	case ast.ExpressionTypeFunctionCall:

	case ast.ExpressionTypeIncrement:
		graph = "{"
		ee := e.Data.(*ast.Expression)
		graph += "digraph{" + g.VisualExpression(ee) + "}"
		graph += `node[ "op" = "++" ]`
		graph += "}"
		return graph
	case ast.ExpressionTypeDecrement:
		graph = "{"
		ee := e.Data.(*ast.Expression)
		graph += "digraph{" + g.VisualExpression(ee) + "}"
		graph += `node[ "op" = "--" ]`
		graph += "}"
		return graph
	case ast.ExpressionTypePrefixIncrement:
		graph = "{"
		ee := e.Data.(*ast.Expression)
		graph += `node[ "op" = "++" ];`
		graph += "digraph{" + g.VisualExpression(ee) + "}"
		graph += "}"
		return graph
	case ast.ExpressionTypePrefixDecrement:
		graph = "{"
		ee := e.Data.(*ast.Expression)
		graph += `node[ "op" = "--" ];`
		graph += "digraph{" + g.VisualExpression(ee) + "}"
		graph += "}"
		return graph
	case ast.ExpressionTypeNegative:
		graph = "{"
		ee := e.Data.(*ast.Expression)
		graph += `node[ "op" = "-" ];`
		graph += "digraph{" + g.VisualExpression(ee) + "}"
		graph += "}"
		return graph
	case ast.ExpressionTypeNot:
		graph = "{"
		ee := e.Data.(*ast.Expression)
		graph += `node[ "op" = "!" ];`
		graph += "digraph{" + g.VisualExpression(ee) + ""
		graph += "}"
		return graph
	case ast.ExpressionTypeBitwiseNot:
		graph = "{"
		ee := e.Data.(*ast.Expression)
		graph += `node[ "op" = "~" ];`
		graph += "digraph{" + g.VisualExpression(ee) + "}"
		graph += "}"
		return graph
	case ast.ExpressionTypeIdentifier:
		return fmt.Sprintf(`"%s"`, e.Data.(*ast.ExpressionIdentifier).Name)
	case ast.ExpressionTypeNew:

	case ast.ExpressionTypeList:
	case ast.ExpressionTypeFunctionLiteral:
	case ast.ExpressionTypeVar:
	case ast.ExpressionTypeConst:
	case ast.ExpressionTypeCheckCast:
	case ast.ExpressionTypeRange:
	case ast.ExpressionTypeSlice:
	case ast.ExpressionTypeMap:
	case ast.ExpressionTypeTypeAssert:
	case ast.ExpressionTypeQuestion:
	case ast.ExpressionTypeGlobal:
	case ast.ExpressionTypeParenthesis:
	case ast.ExpressionTypeVArgs:
	case ast.ExpressionTypeDot:
	}
	return
}

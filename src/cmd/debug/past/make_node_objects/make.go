package make_node_objects

import (
	"fmt"
	"gitee.com/yuyang-fine/lucy/src/cmd/compile/ast"
)

type MakeNodesObjects struct {
}

func (makeNodes *MakeNodesObjects) Make(files map[string][]*ast.TopNode) interface{} {
	if len(files) == 1 {
		for _, v := range files {
			return makeNodes.MakeNodes(v)
		}
	}
	ret := make(map[string]interface{})
	for k, v := range files {
		ret[k] = makeNodes.MakeNodes(v)
	}
	return ret

}

func (makeNodes *MakeNodesObjects) MakeNodes(nodes []*ast.TopNode) (ret []interface{}) {
	for _, n := range nodes {
		switch n.Node.(type) {
		case *ast.Block:
			ss := makeNodes.makeBlock(n.Node.(*ast.Block))
			ret = append(ret, map[string]interface{}{"block": ss})
		case *ast.Function:
			ret = append(ret, makeNodes.makeFunction(n.Node.(*ast.Function)))
		case *ast.Enum:
			ret = append(ret, makeNodes.makeEnum(n.Node.(*ast.Enum)))
		case *ast.Class:
			ret = append(ret, makeNodes.makeClass(n.Node.(*ast.Class)))
		case *ast.Constant:
			ret = append(ret, makeNodes.makeConst(n.Node.(*ast.Constant)))
		case *ast.Import:
			ret = append(ret, makeNodes.makeImport(n.Node.(*ast.Import)))
		case *ast.Expression:
			ret = append(ret, makeNodes.makeExpression(n.Node.(*ast.Expression)))
		case *ast.TypeAlias:
			ret = append(ret, makeNodes.makeTypeAlias(n.Node.(*ast.TypeAlias)))
		default:
			panic("tops have unKnow  type")
		}
	}
	return ret
}

func (makeNodes *MakeNodesObjects) makeConst(c *ast.Constant) interface{} {
	s := fmt.Sprintf("const%d %s", c.Pos.Line, c.Name)
	if c.Type != nil {
		s += fmt.Sprintf(" base on '%s'", c.Type.TypeString())
	}
	var ret []interface{}
	ret = append(ret, s)
	ret = append(ret, map[string]interface{}{"init": makeNodes.makeExpression(c.DefaultValueExpression)})
	return nil
}
func (makeNodes *MakeNodesObjects) makeImport(i *ast.Import) interface{} {
	s := fmt.Sprintf("import@%d", i.Pos.Line)
	s += fmt.Sprintf(" '%s'", i.Import)
	if i.Alias != "" {
		s += fmt.Sprintf(" as '%s'", i.Alias)
	}
	return s
}
func (makeNodes *MakeNodesObjects) makeFunction(fn *ast.Function) interface{} {
	ret := make(map[string]interface{})
	{
		var t []interface{}
		for _, v := range fn.Type.ParameterList {
			p := map[string]interface{}{}
			p["name"] = v.Name
			p["type"] = v.Type.TypeString()
			if v.DefaultValueExpression != nil {
				p["defaultValue"] = makeNodes.makeExpression(v.DefaultValueExpression)
			}
			t = append(t, p)
		}
		ret["parameterlist"] = t
	}
	if fn.Type.VoidReturn() == false {
		if len(fn.Type.ParameterList) == 1 {
			for _, v := range fn.Type.ReturnList {
				p := map[string]interface{}{}
				p["name"] = v.Name
				p["type"] = v.Type.TypeString()
				if v.DefaultValueExpression != nil {
					p["defaultValue"] = makeNodes.makeExpression(v.DefaultValueExpression)
				}
				ret["returnlist"] = p
			}
		} else {
			var t []interface{}
			for _, v := range fn.Type.ReturnList {
				p := map[string]interface{}{}
				p["name"] = v.Name
				p["type"] = v.Type.TypeString()
				if v.DefaultValueExpression != nil {
					p["defaultValue"] = makeNodes.makeExpression(v.DefaultValueExpression)
				}
				t = append(t, p)
			}
			ret["returnlist"] = t
		}

	} else {
		ret["returnlist"] = "VOID"
	}
	ret["block"] = makeNodes.makeBlock(&fn.Block)
	return map[string]interface{}{
		fmt.Sprintf("fn@%d '%s'", fn.Pos.Line, fn.Name): ret,
	}
}

func (makeNodes *MakeNodesObjects) makeEnum(e *ast.Enum) interface{} {
	ret := make(map[string]interface{})
	{
		var t []interface{}
		for _, v := range e.Enums {
			en := make(map[string]interface{})
			en["name"] = v.Name
			if v.NoNeed != nil {
				en["value"] = makeNodes.makeExpression(v.NoNeed)
			}
			t = append(t, en)
		}
		ret["enums"] = t
	}
	return map[string]interface{}{
		fmt.Sprintf("enum@%d '%s'", e.Pos.Line, e.Name): ret,
	}
}

func (makeNodes *MakeNodesObjects) makeClass(c *ast.Class) interface{} {
	ret := make(map[string]interface{})
	ret["name"] = c.Name
	ret["accessFlags"] = c.AccessFlags
	if c.SuperClassName != nil {
		ret["superClass"] = c.SuperClassName.Name
	}
	if len(c.InterfaceNames) > 0 {
		var t []interface{}
		for _, v := range c.InterfaceNames {
			t = append(t, v.Name)
		}
		ret["interfaces"] = t
	}

	if len(c.Block.Constants) > 0 {
		var t []interface{}
		for _, v := range c.Block.Constants {
			t = append(t, makeNodes.makeConst(v))
		}
		ret["constants"] = t
	}
	if len(c.Fields) > 0 {
		fields := make(map[string]interface{})
		for _, v := range c.Fields {
			f := make(map[string]interface{})
			f["accessFlags"] = v.AccessFlags
			f["name"] = v.Name
			f["type"] = v.Type.TypeString()
			if v.DefaultValueExpression != nil {
				f["defaultValue"] = makeNodes.makeExpression(v.DefaultValueExpression)
			}
			fields[v.Name] = f
		}
		ret["fields"] = fields
	}
	if len(c.StaticBlocks) > 0 {
		var t []interface{}
		for _, v := range c.StaticBlocks {
			t = append(t, makeNodes.makeBlock(v))
		}
		ret["staticBlocks"] = t
	}
	if len(c.Methods) > 0 {
		methods := make(map[string]interface{})
		for name, ms := range c.Methods {
			if len(ms) != 1 {
				var t []interface{}
				for _, method := range ms {
					m := make(map[string]interface{})
					m["accessFlags"] = method.Function.AccessFlags
					m["function"] = makeNodes.makeFunction(method.Function)
					t = append(t, m)
				}
				methods[name] = t
			} else {
				method := ms[0]
				m := make(map[string]interface{})
				m["accessFlags"] = method.Function.AccessFlags
				m["function"] = makeNodes.makeFunction(method.Function)
				methods[name] = m
			}

		}
		ret["methods"] = methods
	}
	return map[string]interface{}{
		fmt.Sprintf("class@%d '%s'", c.Pos.Line, c.Name): ret,
	}
}

func (makeNodes *MakeNodesObjects) makeTypeAlias(a *ast.TypeAlias) interface{} {
	s := fmt.Sprintf("typealias@%d '%s' base on '%s'", a.Pos.Line, a.Name, a.Type.TypeString())
	return s
}
func (makeNodes *MakeNodesObjects) makeBlock(b *ast.Block) interface{} {
	var ret []interface{}
	for _, s := range b.Statements {
		ret = append(ret, makeNodes.makeStatement(s))
	}
	return ret
}
func (makeNodes *MakeNodesObjects) makeStatementIf(s *ast.StatementIf) interface{} {
	ret := make(map[string]interface{})
	if len(s.PrefixExpressions) > 0 {
		var t []interface{}
		for _, v := range s.PrefixExpressions {
			t = append(t, makeNodes.makeExpression(v))
		}
		ret["prefixExpressions"] = t
	}
	if s.Condition != nil {
		ret["condition"] = makeNodes.makeExpression(s.Condition)
	}
	ret["block"] = makeNodes.makeBlock(&s.Block)
	if len(s.ElseIfList) > 0 {
		var t []interface{}
		for _, v := range s.ElseIfList {
			t = append(t, map[string]interface{}{
				"condition": makeNodes.makeExpression(v.Condition),
				"block":     makeNodes.makeBlock(v.Block),
			})
		}
		ret["elseifList"] = t
	}
	if s.Else != nil {
		ret["else"] = makeNodes.makeBlock(s.Else)
	}
	return map[string]interface{}{
		fmt.Sprintf("if@%d", s.Pos.Line): ret,
	}
}
func (makeNodes *MakeNodesObjects) makeStatementFor(s *ast.StatementFor) interface{} {
	ret := make(map[string]interface{})
	if s.Init != nil {
		ret["init"] = makeNodes.makeExpression(s.Init)
	}
	if s.Condition != nil {
		ret["condition"] = makeNodes.makeExpression(s.Condition)
	}
	if s.Increment != nil {
		ret["increment"] = makeNodes.makeExpression(s.Increment)
	}
	ret["block"] = makeNodes.makeBlock(s.Block)
	return map[string]interface{}{
		fmt.Sprintf("for@%d", s.Pos.Line): ret,
	}
}
func (makeNodes *MakeNodesObjects) makeStatementSwitch(s *ast.StatementSwitch) interface{} {
	ret := make(map[string]interface{})
	if len(s.PrefixExpressions) > 0 {
		var t []interface{}
		for _, v := range s.PrefixExpressions {
			t = append(t, makeNodes.makeExpression(v))
		}
		ret["prefixExpressions"] = t
	}
	if s.Condition != nil {
		ret["condition"] = makeNodes.makeExpression(s.Condition)
	}
	if len(s.StatementSwitchCases) > 0 {
		var cases []interface{}
		for _, v := range s.StatementSwitchCases {
			onecase := make(map[string]interface{})
			var conditions []interface{}
			for _, vv := range v.Matches {
				conditions = append(conditions, makeNodes.makeExpression(vv))
			}
			onecase["conditions"] = conditions
			onecase["block"] = makeNodes.makeBlock(v.Block)
			cases = append(cases, onecase)
		}
		ret["cases"] = cases
	}
	if s.Default != nil {
		ret["block"] = makeNodes.makeBlock(s.Default)
	}
	return map[string]interface{}{
		fmt.Sprintf("switch@%d", s.Pos.Line): ret,
	}
}
func (makeNodes *MakeNodesObjects) makeStatementWhen(w *ast.StatementWhen) interface{} {
	ret := make(map[string]interface{})
	if w.Condition != nil {
		ret["condition"] = w.Condition.TypeString()
	}
	if len(w.Cases) > 0 {
		var cases []interface{}
		for _, v := range w.Cases {
			onecase := make(map[string]interface{})
			var conditions []interface{}
			for _, vv := range v.Matches {
				conditions = append(conditions, vv.TypeString())
			}
			onecase["conditions"] = conditions
			onecase["block"] = makeNodes.makeBlock(v.Block)
			cases = append(cases, onecase)
		}
		ret["cases"] = cases
	}
	if w.Default != nil {
		ret["block"] = makeNodes.makeBlock(w.Default)
	}
	return map[string]interface{}{
		fmt.Sprintf("when%d", w.Pos.Line): ret,
	}
}

func (makeNodes *MakeNodesObjects) makeStatementLabel(label *ast.StatementLabel) interface{} {
	return fmt.Sprintf("label@%d '%s'", label.Pos.Line, label.Name)
}
func (makeNodes *MakeNodesObjects) makeStatementGoto(g *ast.StatementGoTo) interface{} {
	return fmt.Sprintf("goto@%d '%s'", g.Pos.Line, g.LabelName)
}

func (makeNodes *MakeNodesObjects) makeStatementDefer(d *ast.StatementDefer) interface{} {
	return map[string]interface{}{
		fmt.Sprintf("defer%d", d.Pos.Line): makeNodes.makeBlock(&d.Block),
	}
}
func (makeNodes *MakeNodesObjects) makeStatement(s *ast.Statement) interface{} {
	switch s.Type {
	case ast.StatementTypeExpression:
		return map[string]interface{}{
			fmt.Sprintf("statementExpression@%d", s.Pos.Line): makeNodes.makeExpression(s.Expression),
		}
	case ast.StatementTypeIf:
		return makeNodes.makeStatementIf(s.StatementIf)
	case ast.StatementTypeBlock:
		return map[string]interface{}{
			"block": makeNodes.makeBlock(s.Block),
		}
	case ast.StatementTypeFor:
		return makeNodes.makeStatementFor(s.StatementFor)
	case ast.StatementTypeContinue:
		return fmt.Sprintf("continue@%d", s.Pos.Line)
	case ast.StatementTypeReturn:
		if len(s.StatementReturn.Expressions) == 0 {
			return fmt.Sprintf("return@%d", s.Pos.Line)
		} else {
			key := fmt.Sprintf("return@%d", s.Pos.Line)
			ret := make(map[string]interface{})
			if len(s.StatementReturn.Expressions) == 1 {
				ret[key] = makeNodes.makeExpression(s.StatementReturn.Expressions[0])
			} else {
				var t []interface{}
				for _, v := range s.StatementReturn.Expressions {
					t = append(t, makeNodes.makeExpression(v))
				}
				ret[key] = t
			}
			return ret
		}
	case ast.StatementTypeBreak:
		return fmt.Sprintf("break@%d", s.Pos.Line)
	case ast.StatementTypeSwitch:
		return makeNodes.makeStatementSwitch(s.StatementSwitch)
	case ast.StatementTypeWhen:
		return makeNodes.makeStatementWhen(s.StatementWhen)
	case ast.StatementTypeLabel:
		return makeNodes.makeStatementLabel(s.StatementLabel)
	case ast.StatementTypeGoTo:
		return makeNodes.makeStatementGoto(s.StatementGoTo)
	case ast.StatementTypeDefer:
		return makeNodes.makeStatementDefer(s.Defer)
	case ast.StatementTypeClass:
		return makeNodes.makeClass(s.Class)
	case ast.StatementTypeEnum:
		return makeNodes.makeEnum(s.Enum)
	case ast.StatementTypeNop:
		return fmt.Sprintf("nop")
	case ast.StatementTypeImport:
		return makeNodes.makeImport(s.Import)
	case ast.StatementTypeTypeAlias:
		return makeNodes.makeTypeAlias(s.TypeAlias)
	}
	return nil
}

func (makeNodes *MakeNodesObjects) makeExpression(e *ast.Expression) interface{} {
	if e == nil {
		return fmt.Errorf("EXPRESSION IS NULL")
	}
	switch e.Type {
	case ast.ExpressionTypeNull:
		return "null"
	case ast.ExpressionTypeBool:
		return e.Data.(bool)
	case ast.ExpressionTypeByte:
		return e.Data.(byte)
	case ast.ExpressionTypeShort:
		return e.Data.(int32)
	case ast.ExpressionTypeChar:
		return e.Data.(int32)
	case ast.ExpressionTypeInt:
		return e.Data.(int32)
	case ast.ExpressionTypeLong:
		return e.Data.(int64)
	case ast.ExpressionTypeFloat:
		return e.Data.(float32)
	case ast.ExpressionTypeDouble:
		return e.Data.(float64)
	case ast.ExpressionTypeString:
		return fmt.Sprintf(`literal string@%d "%s"`, e.Pos.Line, e.Data.(string))
	case ast.ExpressionTypeArray:
		array := e.Data.(*ast.ExpressionArray)
		ret := make(map[string]interface{})
		ret["op"] = e.Op

		if array.Type != nil {
			ret["type"] = array.Type.TypeString()
		}
		var t []interface{}
		for _, v := range array.Expressions {
			t = append(t, makeNodes.makeExpression(v))
		}
		ret["elements"] = t
		return ret
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
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		bin := e.Data.(*ast.ExpressionBinary)
		ret["left"] = makeNodes.makeExpression(bin.Left)
		ret["right"] = makeNodes.makeExpression(bin.Right)
		return ret
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
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		bin := e.Data.(*ast.ExpressionBinary)
		if bin.Left.Type == ast.ExpressionTypeList {
			var t []interface{}
			for _, v := range bin.Left.Data.([]*ast.Expression) {
				t = append(t, makeNodes.makeExpression(v))
			}
			ret["left"] = t
		} else {
			ret["left"] = makeNodes.makeExpression(bin.Left)
		}
		ret["right"] = makeNodes.makeExpression(bin.Right)
		return ret
	case ast.ExpressionTypeIndex:
		index := e.Data.(*ast.ExpressionIndex)
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		ret["on"] = makeNodes.makeExpression(index.Expression)
		ret["index"] = makeNodes.makeExpression(index.Index)
		return ret
	case ast.ExpressionTypeSelection:
		selection := e.Data.(*ast.ExpressionSelection)
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		ret["on"] = makeNodes.makeExpression(selection.Expression)
		ret["index"] = selection.Name
		return ret
	case ast.ExpressionTypeSelectionConst:
		selection := e.Data.(*ast.ExpressionSelection)
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		ret["on"] = makeNodes.makeExpression(selection.Expression)
		ret["index"] = selection.Name
		return ret
	case ast.ExpressionTypeMethodCall:
		call := e.Data.(*ast.ExpressionMethodCall)
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		ret["name"] = call.Name
		ret["on"] = makeNodes.makeExpression(call.Expression)
		if len(call.Args) > 0 {
			var t []interface{}
			for _, v := range call.Args {
				t = append(t, makeNodes.makeExpression(v))
			}
			ret["args"] = t
		}
		if len(call.ParameterTypes) > 0 {
			var t []interface{}
			for _, v := range call.ParameterTypes {
				t = append(t, v.TypeString())
			}
			ret["parameterTypes"] = t
		}
		return ret
	case ast.ExpressionTypeFunctionCall:
		call := e.Data.(*ast.ExpressionFunctionCall)
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		ret["on"] = makeNodes.makeExpression(call.Expression)
		if len(call.Args) > 0 {
			var t []interface{}
			for _, v := range call.Args {
				t = append(t, makeNodes.makeExpression(v))
			}
			ret["args"] = t
		}
		if len(call.ParameterTypes) > 0 {
			var t []interface{}
			for _, v := range call.ParameterTypes {
				t = append(t, v.TypeString())
			}
			ret["parameterTypes"] = t
		}
		return ret
	case ast.ExpressionTypeIncrement:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		ret["on"] = makeNodes.makeExpression(e.Data.(*ast.Expression))
		return ret
	case ast.ExpressionTypeDecrement:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		ret["on"] = makeNodes.makeExpression(e.Data.(*ast.Expression))
		return ret
	case ast.ExpressionTypePrefixIncrement:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		ret["on"] = makeNodes.makeExpression(e.Data.(*ast.Expression))
		return ret
	case ast.ExpressionTypePrefixDecrement:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		ret["on"] = makeNodes.makeExpression(e.Data.(*ast.Expression))
		return ret
	case ast.ExpressionTypeNegative:
		fallthrough
	case ast.ExpressionTypeNot:
		fallthrough
	case ast.ExpressionTypeBitwiseNot:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		ret["on"] = makeNodes.makeExpression(e.Data.(*ast.Expression))
		return ret
	case ast.ExpressionTypeIdentifier:
		return e.Data.(*ast.ExpressionIdentifier).Name
	case ast.ExpressionTypeNew:
		n := e.Data.(*ast.ExpressionNew)
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		ret["type"] = n.Type.TypeString()
		{
			var t []interface{}
			for _, v := range n.Args {
				t = append(t, makeNodes.makeExpression(v))
			}
			ret["args"] = t
		}
		return ret
	case ast.ExpressionTypeList:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		list := e.Data.([]*ast.Expression)
		if len(list) == 1 {
			return makeNodes.makeExpression(list[0])
		}
		var t []interface{}
		for _, v := range list {
			t = append(t, makeNodes.makeExpression(v))
		}
		ret["list"] = t

		return ret
	case ast.ExpressionTypeFunctionLiteral:
		return makeNodes.makeFunction(e.Data.(*ast.Function))
	case ast.ExpressionTypeVar:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		v := e.Data.(*ast.ExpressionVar)
		if v.Type != nil {
			ret["type"] = v.Type.TypeString()
		}
		{
			var t []interface{}
			for _, vv := range v.Variables {
				t = append(t, vv.Name)
			}
			ret["variables"] = t
		}
		if len(v.InitValues) > 0 {
			var t []interface{}
			for _, vv := range v.InitValues {
				t = append(t, makeNodes.makeExpression(vv))
			}
			ret["initValues"] = t
		}
		return ret
	case ast.ExpressionTypeConst:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		{
			var t []interface{}
			for _, v := range e.Data.([]*ast.Constant) {
				m := make(map[string]interface{})
				m["name"] = v.Name
				if v.Type != nil {
					m["type"] = v.Type.TypeString()
				}
				if v.DefaultValueExpression != nil {
					m["defaultValue"] = makeNodes.makeExpression(v.DefaultValueExpression)
				}
				t = append(t, m)
			}
			ret["constants"] = t
		}
		return ret
	case ast.ExpressionTypeCheckCast:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		cast := e.Data.(*ast.ExpressionTypeConversion)
		ret["expression"] = makeNodes.makeExpression(cast.Expression)
		ret["type"] = cast.Type.TypeString()
		return ret
	case ast.ExpressionTypeRange:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		ret["on"] = makeNodes.makeExpression(e.Data.(*ast.Expression))
		return ret
	case ast.ExpressionTypeSlice:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		slice := e.Data.(*ast.ExpressionSlice)
		ret["expression"] = makeNodes.makeExpression(slice.ExpressionOn)
		if slice.Start != nil {
			ret["start"] = makeNodes.makeExpression(slice.Start)
		}
		if slice.End != nil {
			ret["end"] = makeNodes.makeExpression(slice.Start)
		}
		return ret
	case ast.ExpressionTypeMap:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		m := e.Data.(*ast.ExpressionMap)
		if m.Type != nil {
			ret["type"] = m.Type.TypeString()
		}
		{
			var t []interface{}
			for _, v := range m.KeyValuePairs {
				tt := make(map[string]interface{})
				tt["key"] = makeNodes.makeExpression(v.Key)
				tt["value"] = makeNodes.makeExpression(v.Value)
				t = append(t, tt)
			}
			ret["paris"] = t
		}
	case ast.ExpressionTypeTypeAssert:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		cast := e.Data.(*ast.ExpressionTypeAssert)
		ret["expression"] = makeNodes.makeExpression(cast.Expression)
		ret["type"] = cast.Type.TypeString()
		return ret
	case ast.ExpressionTypeQuestion:
		ret := make(map[string]interface{})
		ret["op"] = e.Op
		q := e.Data.(*ast.ExpressionQuestion)
		ret["selection"] = makeNodes.makeExpression(q.Selection)
		ret["true"] = makeNodes.makeExpression(q.True)
		ret["false"] = makeNodes.makeExpression(q.False)
		return ret
	case ast.ExpressionTypeGlobal:
		return e.Op // should be "global"
	case ast.ExpressionTypeParenthesis:
		return map[string]interface{}{
			"()": makeNodes.makeExpression(e.Data.(*ast.Expression)),
		}
	case ast.ExpressionTypeVArgs:
		return map[string]interface{}{
			"...": makeNodes.makeExpression(e.Data.(*ast.Expression)),
		}
	case ast.ExpressionTypeDot:
		return "."
	}
	return nil
}

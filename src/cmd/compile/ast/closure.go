package ast

type ClosureVars struct {
	Vars map[string]*ClosureVar
}

func (c *ClosureVars) ClosureVarsExist(name string, v *VariableDefinition) (times uint8, is bool) {
	if c.Vars == nil {
		is = false
		return
	}
	vv, ok := c.Vars[name]
	if !ok || vv.Var != v {
		is = false
		return
	}
	times, is = vv.Level, true
	return
}
func (c *ClosureVars) NotEmpty() bool {
	return c.Vars != nil && len(c.Vars) > 0
}
func (c *ClosureVars) Insert(v *VariableDefinition) {
	if c.Vars == nil {
		c.Vars = make(map[string]*ClosureVar)
	}
	v.BeenCaptured++
	c.Vars[v.Name] = &ClosureVar{
		Var:   v,
		Level: v.BeenCaptured,
	}
}

func (c *ClosureVars) Search(name string) *VariableDefinition {
	if c.Vars == nil {
		return nil
	}
	if x, ok := c.Vars[name]; ok {
		return x.Var
	} else {
		return nil
	}
}

type ClosureVar struct {
	Var   *VariableDefinition
	Level uint8
}

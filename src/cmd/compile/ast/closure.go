package ast

type ClosureVars struct {
	Vars map[*VariableDefinition]struct{}
}

func (c *ClosureVars) ClosureVarsExist(v *VariableDefinition) bool {
	if c.Vars == nil {
		return false
	}
	_, ok := c.Vars[v]
	return ok
}

func (c *ClosureVars) NotEmpty() bool {
	return c.Vars != nil && len(c.Vars) > 0
}
func (c *ClosureVars) Insert(f *Function, v *VariableDefinition) {
	if c.Vars == nil || len(c.Vars) == 0 {
		for _, v := range f.OffsetDestinations {
			*v += 1
		}
		f.VarOffset++
	}
	if c.Vars == nil {
		c.Vars = make(map[*VariableDefinition]struct{})
	}
	c.Vars[v] = struct{}{}
	v.BeenCaptured = true

}

func (c *ClosureVars) Search(name string) *VariableDefinition {
	if c.Vars == nil {
		return nil
	}
	for v, _ := range c.Vars {
		if v.Name == name {
			return v
		}
	}
	return nil
}

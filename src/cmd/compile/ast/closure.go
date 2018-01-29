package ast

type ClosureVars struct {
	Vars map[*VariableDefinition]*ClosureVar
}

func (c *ClosureVars) ClosureVarsExist(v *VariableDefinition) (level uint8, is bool) {
	if c.Vars == nil {
		is = false
		return
	}
	vv, ok := c.Vars[v]
	if ok == false {
		is = false
		return
	}
	return vv.Level, true
}

func (c *ClosureVars) NotEmpty() bool {
	return c.Vars != nil && len(c.Vars) > 0
}
func (c *ClosureVars) Insert(v *VariableDefinition) {
	if c.Vars == nil {
		c.Vars = make(map[*VariableDefinition]*ClosureVar)
	}
	c.Vars[v] = &ClosureVar{
		Level: v.CaptureLevel,
	}
	v.BeenCaptured = true
	v.CaptureLevel++
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

type ClosureVar struct {
	Level uint8
}

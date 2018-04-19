package ast

type Closure struct {
	Vars  map[*VariableDefinition]struct{}
	Funcs map[*Function]struct{}
}

func (c *Closure) ClosureVariableExist(v *VariableDefinition) bool {
	if c.Vars == nil {
		return false
	}
	_, ok := c.Vars[v]
	return ok
}

func (c *Closure) ClosureFunctionExist(v *Function) bool {
	if c.Funcs == nil {
		return false
	}
	_, ok := c.Funcs[v]
	return ok
}

func (c *Closure) NotEmpty(f *Function) bool {
	fff := func() {
		fs := make(map[*Function]struct{})
		for f, _ := range c.Funcs {
			if f.IsClosureFunction {
				fs[f] = struct{}{}
			}
		}
		c.Funcs = fs
	}
	if c.Vars != nil && len(c.Vars) > 0 {
		f.IsClosureFunction = true // incase capture it self
		fff()
		return true
	}
	if c.Funcs == nil || len(c.Funcs) > 0 {
		return false
	}
	fff()
	for _, t := range f.OffsetDestinations {
		*t += 1
	}
	f.VarOffset++
	return true
}

func (c *Closure) Insert(v *VariableDefinition) {
	if c.Vars == nil {
		c.Vars = make(map[*VariableDefinition]struct{})
	}
	c.Vars[v] = struct{}{}
	v.BeenCaptured = true
}

func (c *Closure) InsertFunction(f *Function) {
	if c.Funcs == nil {
		c.Funcs = make(map[*Function]struct{})
	}
	c.Funcs[f] = struct{}{}
}

func (c *Closure) Search(name string) interface{} {
	for v, _ := range c.Vars {
		if v.Name == name {
			return v
		}
	}
	for v, _ := range c.Funcs {
		if v.Name == name {
			return v
		}
	}
	return nil
}

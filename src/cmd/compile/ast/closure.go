package ast

type Closure struct {
	Vars      map[*VariableDefinition]struct{}
	Functions map[*Function]struct{}
}

func (c *Closure) ClosureVariableExist(v *VariableDefinition) bool {
	if c.Vars == nil {
		return false
	}
	_, ok := c.Vars[v]
	return ok
}

func (c *Closure) ClosureFunctionExist(v *Function) bool {
	if c.Functions == nil {
		return false
	}
	_, ok := c.Functions[v]
	return ok
}

func (c *Closure) NotEmpty(f *Function) bool {
	filterOutNotClosureFunction := func() {
		fs := make(map[*Function]struct{})
		for f, _ := range c.Functions {
			if f.IsClosureFunction {
				fs[f] = struct{}{}
			}
		}
		c.Functions = fs
	}
	if c.Vars != nil && len(c.Vars) > 0 {
		f.IsClosureFunction = true // incase capture it self
		filterOutNotClosureFunction()
		return true
	}
	if c.Functions == nil || len(c.Functions) > 0 {
		return false
	}
	filterOutNotClosureFunction()
	return true
}

func (c *Closure) InsertVar(v *VariableDefinition) {
	if c.Vars == nil {
		c.Vars = make(map[*VariableDefinition]struct{})
	}
	c.Vars[v] = struct{}{}
	v.BeenCaptured = true
}

func (c *Closure) InsertFunction(f *Function) {
	if c.Functions == nil {
		c.Functions = make(map[*Function]struct{})
	}
	c.Functions[f] = struct{}{}
}

func (c *Closure) Search(name string) interface{} {
	for v, _ := range c.Vars {
		if v.Name == name {
			return v
		}
	}
	for v, _ := range c.Functions {
		if v.Name == name {
			return v
		}
	}
	return nil
}

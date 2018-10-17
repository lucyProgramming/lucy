package ast

type Closure struct {
	Variables map[*Variable]struct{}
	Functions map[*Function]struct{}
}

func (c *Closure) ClosureVariableExist(v *Variable) bool {
	if c.Variables == nil {
		return false
	}
	_, ok := c.Variables[v]
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
	keepClosureFunction := func() {
		fs := make(map[*Function]struct{})
		for f, _ := range c.Functions {
			if f.IsClosureFunction {
				fs[f] = struct{}{}
			}
		}
		c.Functions = fs
	}
	if c.Variables != nil && len(c.Variables) > 0 {
		f.IsClosureFunction = true // in case capture it self
		keepClosureFunction()
		return true
	}
	keepClosureFunction() // closure function is function too
	return len(c.Functions) > 0
}

func (c *Closure) InsertVar(v *Variable) {
	if c.Variables == nil {
		c.Variables = make(map[*Variable]struct{})
	}
	c.Variables[v] = struct{}{}

}

func (c *Closure) InsertFunction(f *Function) {
	if c.Functions == nil {
		c.Functions = make(map[*Function]struct{})
	}
	c.Functions[f] = struct{}{}
}

func (c *Closure) Search(name string) interface{} {
	for f, _ := range c.Functions {
		if f.Name == name {
			return f
		}
	}
	return nil
}

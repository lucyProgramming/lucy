package ast

type Closure struct {
	Variables map[*Variable]*ClosureMeta
	Functions map[*Function]*ClosureMeta
}

type ClosureMeta struct {
	pos *Pos
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

func (c *Closure) CaptureCount(f *Function) int {
	sum := len(c.Variables)
	for v, _ := range c.Functions {
		if f == v {
			continue
		}
		if v.IsClosureFunction ||
			v.Closure.CaptureCount(f) > 0 {
			sum++
		}
	}
	return sum
}

func (c *Closure) InsertVar(pos *Pos, v *Variable) {
	if c.Variables == nil {
		c.Variables = make(map[*Variable]*ClosureMeta)
	}
	c.Variables[v] = &ClosureMeta{
		pos: pos,
	}
}

func (c *Closure) InsertFunction(pos *Pos, f *Function) {
	if c.Functions == nil {
		c.Functions = make(map[*Function]*ClosureMeta)
	}
	c.Functions[f] = &ClosureMeta{
		pos: pos,
	}
}

func (c *Closure) Search(name string) interface{} {
	for f, _ := range c.Functions {
		if f.Name == name {
			return f
		}
	}
	return nil
}

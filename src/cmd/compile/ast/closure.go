package ast

type Closure struct {
	Variables map[*Variable]*ClosureMeta
	Functions map[*Function]*ClosureMeta
}

type ClosureMeta struct {
	pos *Pos
}



func (this *Closure) ClosureVariableExist(v *Variable) bool {
	if this.Variables == nil {
		return false
	}
	_, ok := this.Variables[v]
	return ok
}

func (this *Closure) ClosureFunctionExist(v *Function) bool {
	if this.Functions == nil {
		return false
	}
	_, ok := this.Functions[v]
	return ok
}

func (this *Closure) CaptureCount(f *Function) int {
	sum := len(this.Variables)
	for v, _ := range this.Functions {
		if f == v {
			continue
		}
		if v.IsClosureFunction {
			sum++
		}
	}
	return sum
}

func (this *Closure) InsertVar(pos *Pos, v *Variable) {
	if this.Variables == nil {
		this.Variables = make(map[*Variable]*ClosureMeta)
	}
	this.Variables[v] = &ClosureMeta{
		pos: pos,
	}
}

func (this *Closure) InsertFunction(pos *Pos, f *Function) {
	if this == nil {
		panic(".........")
	}
	if this.Functions == nil {
		this.Functions = make(map[*Function]*ClosureMeta)
	}
	this.Functions[f] = &ClosureMeta{
		pos: pos,
	}
}

func (this *Closure) Search(name string) interface{} {
	for f, _ := range this.Functions {
		if f.Name == name {
			return f
		}
	}
	return nil
}

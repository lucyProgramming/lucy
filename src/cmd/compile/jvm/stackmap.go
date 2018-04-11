package jvm

type StackMapStateLocalsNumber struct {
	Locals int
}

func (s *StackMapStateLocalsNumber) FromContext(context *Context) {
	s.Locals = len(context.Locals)
}

package ast

import "fmt"

type Defer struct {
	StartPc        int
	StackMapState  interface{} // *jvm.StackMapState
	Block          Block
	ExceptionClass *Class
}

func (d *Defer) registerExceptionClass(c *Class) error {
	if d.ExceptionClass != nil {
		return fmt.Errorf("exception class already registed as '%s'", d.ExceptionClass.Name)
	}
	d.ExceptionClass = c
	return nil
}

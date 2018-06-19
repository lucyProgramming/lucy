package ast

import "fmt"

type StatementDefer struct {
	allowCatch     bool        // if in function top block
	StartPc        int         // for jvm
	StackMapState  interface{} // *jvm.StackMapState
	Block          Block
	ExceptionClass *Class
}

func (d *StatementDefer) registerExceptionClass(c *Class) error {
	if d.ExceptionClass != nil {
		return fmt.Errorf("exception class already registed as '%s'",
			d.ExceptionClass.Name)
	}
	d.ExceptionClass = c
	return nil
}

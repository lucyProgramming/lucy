package ast

import "fmt"

type Defer struct {
	StartPc        int
	Block          Block
	ExceptionClass *Class
}

func (d *Defer) registerExceptionClass(c *Class) error {
	if d.ExceptionClass != nil {
		return fmt.Errorf("register class already registed as '%d'", d.ExceptionClass.Name)
	}
	d.ExceptionClass = c
	return nil
}

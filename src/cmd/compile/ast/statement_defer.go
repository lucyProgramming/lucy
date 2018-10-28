package ast

import "fmt"

type StatementDefer struct {
	Pos            *Pos
	Labels         []*StatementLabel
	StartPc        int         // for jvm
	StackMapState  interface{} // *jvm.StackMapState
	Block          Block
	ExceptionClass *Class
}

func (d *StatementDefer) ResetLabels() {
	for _, v := range d.Labels {
		v.Reset()
	}
}

func (d *StatementDefer) registerExceptionClass(c *Class) error {
	if d.ExceptionClass != nil {
		return fmt.Errorf("exception class already registed as '%s'",
			d.ExceptionClass.Name)
	}
	d.ExceptionClass = c
	return nil
}

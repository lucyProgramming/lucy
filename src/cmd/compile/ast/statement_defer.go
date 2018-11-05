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

func (this *StatementDefer) ResetLabels() {
	for _, v := range this.Labels {
		v.Reset()
	}
}

func (this *StatementDefer) registerExceptionClass(c *Class) error {
	if this.ExceptionClass != nil {
		return fmt.Errorf("exception class already registed as '%s'",
			this.ExceptionClass.Name)
	}
	this.ExceptionClass = c
	return nil
}

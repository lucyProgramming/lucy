package ast

import "fmt"

type Pos struct {
	Filename string
	Line     int
	Column   int
	/*
		offset at bs , for special use。
		for template function only currently
	*/
	Offset int
}

func (pos *Pos) ErrMsgPrefix() string {
	return fmt.Sprintf("%s:%d:%d:", pos.Filename, pos.Line, pos.Column)
}

type NameWithPos struct {
	Name string
	Pos  *Pos
}

func errMsgPrefix(pos *Pos) string {
	return pos.ErrMsgPrefix()
}

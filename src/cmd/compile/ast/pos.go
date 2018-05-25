package ast

type Pos struct {
	Filename    string
	StartLine   int
	StartColumn int
	Offset      int // offset at bs , for special use
}

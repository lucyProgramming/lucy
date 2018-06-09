package ast

type Pos struct {
	Filename    string
	StartLine   int
	StartColumn int
	/*
		offset at bs , for special use
	*/
	Offset int
}

type NameWithPos struct {
	Name string
	Pos  *Pos
}

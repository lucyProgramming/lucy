package ast

type Position struct {
	Filename    string
	StartLine   int
	StartColumn int
	//EndLint int
	//EndCloumn int
	/*
		offset at bs , for special use
	*/
	Offset int
}

type NameWithPos struct {
	Name string
	Pos  *Position
}

package ast

type Position struct {
	Filename    string
	StartLine   int
	StartColumn int
	//EndLint int
	//EndColumn int
	/*
		offset at bs , for special use
	*/
	Offset int
}

type NameWithPos struct {
	Name string
	Pos  *Position
}

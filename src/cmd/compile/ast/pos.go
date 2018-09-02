package ast

type Pos struct {
	Filename    string
	StartLine   int ``
	StartColumn int
	/*
		offset at bs , for special use。
		for template function only currently
	*/
	Offset int
	//EndLint int
	//EndColumn int
}

type NameWithPos struct {
	Name string
	Pos  *Pos
}

package ast

var (
	small_float          = 0.0001
	negative_small_float = -small_float
	Nodes                *[]*Node //
)

type NameWithPos struct {
	Name string
	Pos  *Pos
}

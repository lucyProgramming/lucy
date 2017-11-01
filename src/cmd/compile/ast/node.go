package ast

//代表语法数的一个节点
type Node struct {
	Pos  Pos
	Data interface{} //class defination or varialbe Defination
}

//type Tops []*Node //语法树顶层结构

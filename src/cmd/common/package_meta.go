package common

import (
	"time"
)

type FileMeta struct {
	LastModify time.Time
}

type PackageMeta struct {
	MainClass    string
	CompiledFrom map[string]*FileMeta // filename -> meta
	Imports      []string             //lucy package that imported
	CompileTime  time.Time
}

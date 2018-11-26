package common

type FileMeta struct {
	LastModify int64 `json:"lastModify"` // unix seconds
	Name string `json:"name"`
}

type PackageMeta struct {
	CompiledFrom []*FileMeta `json:"compiledFrom"` // filename -> meta
	Imports      []string             `json:"imports"`      //lucy package that imported
	CompileTime  int64                `json:"compileTime"`  // unix seconds
	Classes      []string             `json:"classes"`
}



package common

type FileMeta struct {
	LastModify int64 `json:"lastModify"` // unix seconds
}

type PackageMeta struct {
	CompiledFrom map[string]*FileMeta `json:"compiledFrom"` // filename -> meta
	Imports      []string             `json:"imports"`      //lucy package that imported
	CompileTime  int64                `json:"compileTime"`  // unix seconds
	Classes      []string             `json:"classes"`
}

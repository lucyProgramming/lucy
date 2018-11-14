package run

type Flags struct {
	forceReBuild  bool
	build         bool
	verbose       bool
	compilerFlags string
	help          bool
	goCompiler    bool
}

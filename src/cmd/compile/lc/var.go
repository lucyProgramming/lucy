package lc

//	"github.com/756445638/lucy/src/cmd/compile/ast"

var (
	CompileFlags Flags
	l            LucyCompile
)

type Flags struct {
	OnlyImport bool
}

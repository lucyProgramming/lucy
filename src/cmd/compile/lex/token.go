package lex

const (
	_                  = iota
	TokenFunction      // fn
	TokenEnum          // enum
	TokenConst         // const
	TokenIf            // if
	TokenElseif        // elseif
	TokenElse          // else
	TokenFor           // for
	TokenBreak         // break
	TokenContinue      // continue
	TokenReturn        // return
	TokenNull          // null
	TokenBool          // bool
	TokenTrue          // true
	TokenFalse         // false
	TokenLp            // (
	TokenRp            // )
	TokenLc            // {
	TokenRc            // }
	TokenLb            // [
	TokenRb            // ]
	TokenPass          // pass
	TokenSemicolon     // ;
	TokenLf            // "\n"
	TokenComma         // ,
	TokenLogicalAnd    // &&
	TokenLogicalOr     // ||
	TokenAnd           // &
	TokenOr            // |
	TokenLsh           // <<
	TokenRsh           // >>
	TokenXor           // ^
	TokenBitNot        // ~
	TokenAssign        // =
	TokenEqual         // ==
	TokenNe            // !=
	TokenGt            // >
	TokenGe            // >=
	TokenLt            // <
	TokenLe            // <=
	TokenAdd           // +
	TokenSub           // -
	TokenMul           // *
	TokenDiv           // a/c
	TokenMod           // a%b
	TokenIncrement     // a++
	TokenDecrement     // a--
	TokenSelection     // a.do()
	TokenVar           // var a
	TokenNew           // new Object()
	TokenColon         // :
	TokenColonAssign   // :=
	TokenAddAssign     // +=
	TokenSubAssign     // -=
	TokenMulAssign     // *=
	TokenDivAssign     // /=
	TokenModAssign     // %=
	TokenAndAssign     // &=
	TokenOrAssign      // |=
	TokenXorAssign     // ^=
	TokenLshAssign     // <<=
	TokenRshAssign     // >>=
	TokenNot           // !false
	TokenSwitch        // switch
	TokenCase          // case
	TokenDefault       // default
	TokenImport        // import
	TokenAs            // as
	TokenClass         // class
	TokenStatic        // static
	TokenPublic        // public
	TokenProtected     // protected
	TokenPrivate       // private
	TokenInterface     // interface
	TokenByte          // byte
	TokenShort         // short
	TokenInt           // int
	TokenFloat         // float
	TokenDouble        // double
	TokenLong          // long
	TokenString        // string
	TokenIdentifier    // identifier
	TokenLiteralByte   // 'a'
	TokenLiteralShort  // 1s
	TokenLiteralInt    // 123
	TokenLiteralString // ""
	TokenLiteralFloat  // 0.000
	TokenLiteralDouble // 0.0d
	TokenLiteralLong   // 100L
	TokenDefer         // defer
	TokenType          // type
	TokenArrow         // ->
	TokenExtends       // extends
	TokenImplements    // implements
	TokenGoto          // goto
	TokenRange         // range
	TokenMap           // map
	TokenTemplate      // T or T1
	TokenQuestion      // ?
	TokenVolatile      // volatile
	TokenSynchronized  // synchronized
	TokenFinal         // final
	TokenGlobal        // global
	TokenVArgs         // ...
	TokenEof           // end of file
)

var (
	keywordsMap = map[string]int{
		"fn":           TokenFunction,
		"enum":         TokenEnum,
		"const":        TokenConst,
		"if":           TokenIf,
		"elseif":       TokenElseif,
		"else":         TokenElse,
		"for":          TokenFor,
		"break":        TokenBreak,
		"continue":     TokenContinue,
		"return":       TokenReturn,
		"null":         TokenNull,
		"bool":         TokenBool,
		"true":         TokenTrue,
		"false":        TokenFalse,
		"pass":         TokenPass,
		"var":          TokenVar,
		"new":          TokenNew,
		"switch":       TokenSwitch,
		"case":         TokenCase,
		"default":      TokenDefault,
		"import":       TokenImport,
		"as":           TokenAs,
		"class":        TokenClass,
		"static":       TokenStatic,
		"public":       TokenPublic,
		"protected":    TokenProtected,
		"private":      TokenPrivate,
		"interface":    TokenInterface,
		"byte":         TokenByte,
		"short":        TokenShort,
		"int":          TokenInt,
		"float":        TokenFloat,
		"double":       TokenDouble,
		"long":         TokenLong,
		"string":       TokenString,
		"defer":        TokenDefer,
		"type":         TokenType,
		"extends":      TokenExtends,
		"implements":   TokenImplements,
		"goto":         TokenGoto,
		"range":        TokenRange,
		"map":          TokenMap,
		"volatile":     TokenVolatile,
		"synchronized": TokenSynchronized,
		"final":        TokenFinal,
		"global":       TokenGlobal,
	}
)

type Token struct {
	Offset      int // bs offset
	Type        int
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
	Description string
	Data        interface{}
}

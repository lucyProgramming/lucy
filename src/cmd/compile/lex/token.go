package lex

type TokenKind int

const (
	_                     TokenKind = iota
	TokenFn                         // fn
	TokenEnum                       // enum
	TokenConst                      // const
	TokenIf                         // if
	TokenElseif                     // elseif
	TokenElse                       // else
	TokenFor                        // for
	TokenBreak                      // break
	TokenContinue                   // continue
	TokenReturn                     // return
	TokenNull                       // null
	TokenBool                       // bool
	TokenTrue                       // true
	TokenFalse                      // false
	TokenLp                         // (
	TokenRp                         // )
	TokenLc                         // {
	TokenRc                         // }
	TokenLb                         // [
	TokenRb                         // ]
	TokenPass                       // pass
	TokenSemicolon                  // ;
	TokenLf                         // "\n"
	TokenComma                      // ,
	TokenLogicalAnd                 // &&
	TokenLogicalOr                  // ||
	TokenAnd                        // &
	TokenOr                         // |
	TokenLsh                        // <<
	TokenRsh                        // >>
	TokenXor                        // ^
	TokenBitNot                     // ~
	TokenAssign                     // =
	TokenEqual                      // ==
	TokenNe                         // !=
	TokenGt                         // >
	TokenGe                         // >=
	TokenLt                         // <
	TokenLe                         // <=
	TokenAdd                        // +
	TokenSub                        // -
	TokenMul                        // *
	TokenDiv                        // a/c
	TokenMod                        // a%b
	TokenIncrement                  // a++
	TokenDecrement                  // a--
	TokenSelection                  // a.do()
	TokenVar                        // var a
	TokenNew                        // new Object()
	TokenColon                      // :
	TokenSelectConst                // ::
	TokenVarAssign                  // :=
	TokenAddAssign                  // +=
	TokenSubAssign                  // -=
	TokenMulAssign                  // *=
	TokenDivAssign                  // /=
	TokenModAssign                  // %=
	TokenAndAssign                  // &=
	TokenOrAssign                   // |=
	TokenXorAssign                  // ^=
	TokenLshAssign                  // <<=
	TokenRshAssign                  // >>=
	TokenNot                        // !false
	TokenSwitch                     // switch
	TokenCase                       // case
	TokenDefault                    // default
	TokenImport                     // import
	TokenAs                         // as
	TokenClass                      // class
	TokenStatic                     // static
	TokenPublic                     // public
	TokenProtected                  // protected
	TokenPrivate                    // private
	TokenInterface                  // interface
	TokenByte                       // byte
	TokenShort                      // short
	TokenChar                       // char
	TokenInt                        // int
	TokenFloat                      // float
	TokenDouble                     // double
	TokenLong                       // long
	TokenString                     // string
	TokenIdentifier                 // identifier
	TokenLiteralByte                // 'a'
	TokenLiteralChar                // 'a'
	TokenLiteralShort               // 1s
	TokenLiteralInt                 // 123
	TokenLiteralLong                // 100L
	TokenLiteralFloat               // 0.000
	TokenLiteralDouble              // 0.0d
	TokenLiteralString              // ""
	TokenDefer                      // defer
	TokenTypeAlias                  // type
	TokenArrow                      // ->
	TokenExtends                    // extends
	TokenImplements                 // implements
	TokenGoto                       // goto
	TokenRange                      // range
	TokenMap                        // map
	TokenQuestion                   // ?
	TokenVolatile                   // volatile
	TokenSynchronized               // synchronized
	TokenFinal                      // final
	TokenAbstract                   // abstract
	TokenGlobal                     // global
	TokenVArgs                      // ...
	TokenWhen                       // when
	TokenComment                    //
	TokenMultiLineComment           //
	TokenEof                        // end of file
)

var (
	keywordsMap = map[string]TokenKind{
		"fn":           TokenFn,
		"enum":         TokenEnum,
		"const":        TokenConst,
		"if":           TokenIf,
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
		"char":         TokenChar,
		"int":          TokenInt,
		"float":        TokenFloat,
		"double":       TokenDouble,
		"long":         TokenLong,
		"string":       TokenString,
		"defer":        TokenDefer,
		"typealias":    TokenTypeAlias,
		"extends":      TokenExtends,
		"implements":   TokenImplements,
		"goto":         TokenGoto,
		"range":        TokenRange,
		"map":          TokenMap,
		"volatile":     TokenVolatile,
		"synchronized": TokenSynchronized,
		"final":        TokenFinal,
		"global":       TokenGlobal,
		"abstract":     TokenAbstract,
		"when":         TokenWhen,
	}
)

type Token struct {
	Type        TokenKind
	Offset      int // bs offset
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
	Description string
	Data        interface{}
}

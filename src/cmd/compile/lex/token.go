package lex

const (
	_                    = iota
	TOKEN_FUNCTION       // fn
	TOKEN_ENUM           // enum
	TOKEN_CONST          //const
	TOKEN_IF             // if
	TOKEN_ELSEIF         //elseif
	TOKEN_ELSE           // else
	TOKEN_FOR            //for
	TOKEN_BREAK          //break
	TOKEN_CONTINUE       //continue
	TOKEN_RETURN         //return
	TOKEN_NULL           // null
	TOKEN_BOOL           //bool
	TOKEN_TRUE           //true
	TOKEN_FALSE          //false
	TOKEN_LP             //(
	TOKEN_RP             //)
	TOKEN_LC             //{
	TOKEN_RC             //}
	TOKEN_LB             //[
	TOKEN_RB             //]
	TOKEN_SKIP           //skip;
	TOKEN_SEMICOLON      // ;
	TOKEN_CRLF           // enter
	TOKEN_COMMA          //,
	TOKEN_LOGICAL_AND    // &&
	TOKEN_LOGICAL_OR     // ||
	TOKEN_AND            // &
	TOKEN_OR             // |
	TOKEN_LEFT_SHIFT     // <<
	TOKEN_RIGHT_SHIFT    // >>
	TOKEN_XOR            // ^
	TOKEN_ASSIGN         //=
	TOKEN_EQUAL          //== or ==
	TOKEN_NE             // !=
	TOKEN_GT             //>
	TOKEN_GE             //>=
	TOKEN_LT             //<
	TOKEN_LE             //<=
	TOKEN_ADD            //+
	TOKEN_SUB            //-
	TOKEN_MUL            //*
	TOKEN_DIV            // a/c
	TOKEN_MOD            // a%b
	TOKEN_INCREMENT      //a++
	TOKEN_DECREMENT      //a--
	TOKEN_DOT            // a.do()
	TOKEN_VAR            // var a
	TOKEN_NEW            // new Object()
	TOKEN_COLON          // :
	TOKEN_COLON_ASSIGN   // :=
	TOKEN_ADD_ASSIGN     // +=
	TOKEN_SUB_ASSIGN     // -=
	TOKEN_MUL_ASSIGN     // *=
	TOKEN_DIV_ASSIGN     // /=
	TOKEN_MOD_ASSIGN     // %=
	TOKEN_NOT            // !false
	TOKEN_SWITCH         //swtich
	TOKEN_CASE           //case
	TOKEN_DEFAULT        //default
	TOKEN_IMPORT         //import
	TOKEN_AS             //as
	TOKEN_CLASS          //class
	TOKEN_STATIC         //static
	TOKEN_PUBLIC         //public
	TOKEN_PROTECTED      //protected
	TOKEN_PRIVATE        //private
	TOKEN_INTERFACE      //interface
	TOKEN_BYTE           //byte
	TOKEN_SHORT          // short
	TOKEN_INT            //int
	TOKEN_FLOAT          //float
	TOKEN_DOUBLE         //double
	TOKEN_LONG           //long
	TOKEN_STRING         //string
	TOKEN_IDENTIFIER     // identifier
	TOKEN_LITERAL_BOOL   // true or false
	TOKEN_LITERAL_BYTE   //'a'
	TOKEN_LITERAL_SHORT  // 1s
	TOKEN_LITERAL_INT    // 123
	TOKEN_LITERAL_STRING // ""
	TOKEN_LITERAL_FLOAT  // 0.000
	TOKEN_LITERAL_DOUBLE // 0.0d
	TOKEN_LITERAL_LONG   //
	TOKEN_DEFER          // defer
	TOKEN_TYPE           //type
	TOKEN_ARROW          //->
	TOKEN_EXTENDS        //extends
	TOKEN_IMPLEMENTS     // implements
	TOKEN_GOTO           //goto
	TOKEN_RANGE          //range
	TOKEN_MAP            //map
)

var (
	keywordMap = map[string]int{
		"fn":         TOKEN_FUNCTION,
		"enum":       TOKEN_ENUM,
		"const":      TOKEN_CONST,
		"if":         TOKEN_IF,
		"elseif":     TOKEN_ELSEIF,
		"else":       TOKEN_ELSE,
		"for":        TOKEN_FOR,
		"break":      TOKEN_BREAK,
		"continue":   TOKEN_CONTINUE,
		"return":     TOKEN_RETURN,
		"null":       TOKEN_NULL,
		"bool":       TOKEN_BOOL,
		"true":       TOKEN_TRUE,
		"false":      TOKEN_FALSE,
		"skip":       TOKEN_SKIP,
		"var":        TOKEN_VAR,
		"new":        TOKEN_NEW,
		"switch":     TOKEN_SWITCH,
		"case":       TOKEN_CASE,
		"default":    TOKEN_DEFAULT,
		"import":     TOKEN_IMPORT,
		"as":         TOKEN_AS,
		"class":      TOKEN_CLASS,
		"static":     TOKEN_STATIC,
		"public":     TOKEN_PUBLIC,
		"protected":  TOKEN_PROTECTED,
		"private":    TOKEN_PRIVATE,
		"interface":  TOKEN_INTERFACE,
		"byte":       TOKEN_BYTE,
		"short":      TOKEN_SHORT,
		"int":        TOKEN_INT,
		"float":      TOKEN_FLOAT,
		"double":     TOKEN_DOUBLE,
		"long":       TOKEN_LONG,
		"string":     TOKEN_STRING,
		"defer":      TOKEN_DEFER,
		"type":       TOKEN_TYPE,
		"extends":    TOKEN_EXTENDS,
		"implements": TOKEN_IMPLEMENTS,
		"goto":       TOKEN_GOTO,
		"range":      TOKEN_RANGE,
		"map":        TOKEN_MAP,
	}
)

type Token struct {
	Type        int
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
	Desp        string
	Data        interface{}
}

package lex

import (
	"strconv"

	"github.com/timtadh/lexmachine/machines"
)

const (
	TOKEN_FUNCTION       = 1  // function
	TOKEN_ENUM           = 2  // enum
	TOKEN_CONST          = 3  //const
	TOKEN_IF             = 4  // if
	TOKEN_ELSEIF         = 5  //elseif
	TOKEN_ELSE           = 6  // else
	TOKEN_FOR            = 7  //for
	TOKEN_BREAK          = 8  //break
	TOKEN_CONTINUE       = 9  //continue
	TOKEN_RETURN         = 10 //return
	TOKEN_NULL           = 11 // null
	TOKEN_BOOL           = 12 //bool
	TOKEN_LP             = 13 //(
	TOKEN_RP             = 14 //)
	TOKEN_LC             = 15 //{
	TOKEN_RC             = 16 //}
	TOKEN_LB             = 17 //[
	TOKEN_RB             = 18 //]
	TOKEN_SKIP           = 19 //skip;
	TOKEN_SEMICOLON      = 20 // ;
	TOKEN_CRLF           = 21 // enter
	TOKEN_COMMA          = 22 //,
	TOKEN_LOGICAL_AND    = 23 // &&
	TOKEN_LOGICAL_OR     = 24 // ||
	TOKEN_AND            = 25 // &
	TOKEN_OR             = 26 // |
	TOKEN_LEFT_SHIFT     = 27 // <<
	TOKEN_RIGHT_SHIFT    = 28 // >>
	TOKEN_ASSIGN         = 29 //=
	TOKEN_EQUAL          = 30 //== or ===
	TOKEN_NE             = 31 // !=
	TOKEN_GT             = 32 //>
	TOKEN_GE             = 33 //>=
	TOKEN_LT             = 34 //<
	TOKEN_LE             = 35 //<=
	TOKEN_ADD            = 36 //+
	TOKEN_SUB            = 37 //-
	TOKEN_MUL            = 38 //*
	TOKEN_DIV            = 39 // a/c
	TOKEN_MOD            = 40 // a%b
	TOKEN_INCREMENT      = 41 //a++
	TOKEN_DECREMENT      = 42 //a--
	TOKEN_DOT            = 43 // a.do()
	TOKEN_VAR            = 44 // var a
	TOKEN_NEW            = 45 // new Object()
	TOKEN_COLON          = 46 // :
	TOKEN_COLON_ASSIGN   = 47 // :=
	TOKEN_PLUS_ASSIGN    = 48 // +=
	TOKEN_MINUS_ASSIGN   = 49 // -=
	TOKEN_MUL_ASSIGN     = 50 // *=
	TOKEN_DIV_ASSIGN     = 51 // /=
	TOKEN_MOD_ASSIGN     = 52 // %=
	TOKEN_NOT            = 53 // !false
	TOKEN_SWITCH         = 54 //swtich
	TOKEN_CASE           = 55 //case
	TOKEN_DEFAULT        = 56 //default
	TOKEN_PACKAGE        = 57 //package
	TOKEN_IMPORT         = 58 //import
	TOKEN_AS             = 59 //as
	TOKEN_CLASS          = 60 //class
	TOKEN_STATIC         = 61 //static
	TOKEN_PUBLIC         = 62 //public
	TOKEN_PROTECTED      = 63 //protected
	TOKEN_PRIVATE        = 64 //private
	TOKEN_INTERFACE      = 65 //interface
	TOKEN_BYTE           = 66 //byte
	TOKEN_INT            = 67 //int
	TOKEN_FLOAT          = 68 //float
	TOKEN_STRING         = 69 //string
	TOKEN_IDENTIFIER     = 70 // identifier
	TOKEN_LITERAL_BOOL   = 71 // true or false
	TOKEN_LITERAL_BYTE   = 72 //'a'
	TOKEN_LITERAL_INT    = 73 // 123
	TOKEN_LITERAL_STRING = 74 // ""
	TOKEN_LITERAL_FLOAT  = 75 // 0.000
	TOKEN_TRY            = 76 // try
	TOKEN_CATCH          = 77 //catch
	TOKEN_FINALLY        = 78 //finally
	TOKEN_THROW          = 79 //throw
	TOKEN_TYPE           = 80 //type
)

type Token struct {
	Type  int
	Match *machines.Match
	Desp  string
	Data  interface{}
}

func parseInt64(bs []byte) int64 {
	t, _ := strconv.ParseInt(string(bs), 0, 64)
	return t
}

//TODO::解析科学计数法
func parseScientificNotation(bs []byte) (data interface{}, token int) {
	////negative := false
	//if bs[0] == '+' {
	//	bs = bs[1:]
	//}
	//if bs[0] == '-' {
	//	//negative = true
	//	bs = bs[1:]
	//}
	//index := bytes.IndexByte(bs, 'e')
	//pre := bs[0:index]
	//suf := bs[index+1:]
	//power := int(parseInt64(suf))
	////1e5  or 1e-5
	//token = TOKEN_INT
	//base := parseInt64(pre)
	//
	//fmt.Println("@@@@@@@@@@@@@@@@", string(pre), string(suf), string(power))
	return 100, TOKEN_INT

}

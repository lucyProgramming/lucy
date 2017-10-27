package lex


import(
	"github.com/timtadh/lexmachine"
)
const (
	TOKEN_FUNCTION = iota // function
	TOKEN_IF              // if
	TOKEN_ELSEIF          //elseif
	TOKEN_ELSE            // else
	TOKEN_FOR             //for
	TOKEN_RETURN          //return
	TOKEN_NULL,    // null
	TOKEN_TRUE, //true
	TOKEN_FALSE, //false
	TOKEN_LP, //(
	TOKEN_RP, //)
	TOKEN_LC, //{
	TOKEN_RC, //}
	TOKEN_LB, //[
	TOKEN_RB, //]
	TOKEN_SEMICOLON, // ;
	TOKEN_COMMA, //,
	TOKEN_LOGICAL_AND, // &&
	TOKEN_LOGICAL_OR, // ||
	TOKEN_ASSIGN, //=
	TOKEN_EQUAL, //== or ===
	TOKEN_NE, // !=
	TOKEN_GT, //>
	TOKEN_GE, //>=
	TOKEN_LT, //<
	TOKEN_LE, //<=
	TOKEN_ADD, //+
	TOKEN_SUB, //-
	TOKEN_MUL, //*
	TOKEN_DIV, // a/c
	TOKEN_MOD, // a%b
	TOKEN_INCREMENT, //a++
	TOKEN_DECREMENT, //a--
	TOKEN_DOT, // a.do()
	TOKEN_VAR, // var a
	TOKEN_NEW, // new Object()
	TOKEN_COLON, // :
	TOKEN_PLUS_ASSIGN, // +=
	TOKEN_MINUS_ASSIGN, // -=
	TOKEN_MUL_ASSIGN, // *=
	TOKEN_DIV_ASSIGN, // /=
	TOKEN_MOD_ASSIGN, // %=
	TOKEN_NOT, // !false
	TOKEN_SWITCH, //swtich
	TOKEN_CASE, //case
	TOKEN_DEFAULT //default
	TOKEN_CRLF       // enter
	TOKEN_PACKAGE    //package
	TOKEN_CLASS      //class
	TOKEN_INT //int
	TOKEN_BOOL //bool
	TOKEN_FLOAT //float
	TOKEN_STRING //string
	TOKEN_IDENTIFIER // identifier

	TOKEN_LITERAL_INT // 123
	TOKEN_LITERAL_STRING // ""
	TOKEN_LITERAL_FLOAT



)

type Token struct {
	lexmachine.Token
}

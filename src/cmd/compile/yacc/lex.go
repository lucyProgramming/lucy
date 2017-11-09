package yacc

import (
	"fmt"

	"github.com/756445638/lucy/src/cmd/compile/lex"
	"github.com/timtadh/lexmachine"
)

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// The parser uses the type <prefix>Lex as a lexer. It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type LucyLex struct {
	Filename string
	Scanner *lexmachine.Scanner
	//	line    []byte
	//	peek    rune
	Errs []error
}

// The parser calls this method to get each new token. This
// implementation returns operators and NUM.
func (x *LucyLex) Lex(yylval *LucySymType) int {
retry:
	var token *lex.Token
	var t interface{}
	var err error
	var endOfFile bool
	for t == nil { //first time,must be nil
		t, err, endOfFile = x.Scanner.Next()
		if endOfFile {
			return eof
		}
		if err != nil {
			x.Error(err.Error())
			continue
		}
	}
	token = t.(*lex.Token)
	defer func(){
		current_pos.Filename = x.Filename
		current_pos.StartLine =
	}()
	switch token.Type {
	case lex.TOKEN_FUNCTION:
		return TOKEN_FUNCTION
	case lex.TOKEN_ENUM:
		return TOKEN_ENUM
	case lex.TOKEN_CONST:
		return TOKEN_CONST
	case lex.TOKEN_IF:
		return TOKEN_CONST
	case lex.TOKEN_ELSEIF:
		return TOKEN_IF
	case lex.TOKEN_ELSE:
		return TOKEN_ELSE
	case lex.TOKEN_FOR:
		return TOKEN_FOR
	case lex.TOKEN_BREAK:
		return TOKEN_BREAK
	case lex.TOKEN_CONTINUE:
		return TOKEN_CONTINUE
	case lex.TOKEN_RETURN:
		return TOKEN_RETURN
	case lex.TOKEN_NULL:
		return TOKEN_NULL
	case lex.TOKEN_BOOL:
		return TOKEN_BOOL
	case lex.TOKEN_LP:
		return TOKEN_LP
	case lex.TOKEN_RP:
		return TOKEN_RP
	case lex.TOKEN_LC:
		return TOKEN_LC
	case lex.TOKEN_RC:
		return TOKEN_RC
	case lex.TOKEN_LB:
		return TOKEN_LB
	case lex.TOKEN_RB:
		return TOKEN_RB
	case lex.TOKEN_SKIP:
		return TOKEN_SKIP
	case lex.TOKEN_SEMICOLON:
		return TOKEN_SEMICOLON
	case lex.TOKEN_CRLF:
		goto retry
	case lex.TOKEN_COMMA:
		return TOKEN_COMMA
	case lex.TOKEN_LOGICAL_AND:
		return TOKEN_LOGICAL_AND
	case lex.TOKEN_LOGICAL_OR:
		return TOKEN_LOGICAL_OR
	case lex.TOKEN_AND:
		return TOKEN_AND
	case lex.TOKEN_OR:
		return TOKEN_OR
	case lex.TOKEN_LEFT_SHIFT:
		return TOKEN_LEFT_SHIFT
	case lex.TOKEN_RIGHT_SHIFT:
		return TOKEN_RIGHT_SHIFT
	case lex.TOKEN_ASSIGN:
		return TOKEN_ASSIGN
	case lex.TOKEN_EQUAL:
		return TOKEN_EQUAL
	case lex.TOKEN_NE:
		return TOKEN_NE
	case lex.TOKEN_GT:
		return TOKEN_GT
	case lex.TOKEN_GE:
		return TOKEN_GE
	case lex.TOKEN_LT:
		return TOKEN_LT
	case lex.TOKEN_LE:
		return TOKEN_LE
	case lex.TOKEN_ADD:
		return TOKEN_ADD
	case lex.TOKEN_SUB:
		return TOKEN_SUB
	case lex.TOKEN_MUL:
		return TOKEN_MUL
	case lex.TOKEN_DIV:
		return TOKEN_DIV
	case lex.TOKEN_MOD:
		return TOKEN_MOD
	case lex.TOKEN_INCREMENT:
		return TOKEN_INCREMENT
	case lex.TOKEN_DECREMENT:
		return TOKEN_DECREMENT
	case lex.TOKEN_DOT:
		return TOKEN_DOT
	case lex.TOKEN_VAR:
		return TOKEN_VAR
	case lex.TOKEN_NEW:
		return TOKEN_NEW
	case lex.TOKEN_COLON:
		return TOKEN_COLON
	case lex.TOKEN_COLON_ASSIGN:
		return TOKEN_COLON_ASSIGN
	case lex.TOKEN_PLUS_ASSIGN:
		return TOKEN_PLUS_ASSIGN
	case lex.TOKEN_MINUS_ASSIGN:
		return TOKEN_MINUS_ASSIGN
	case lex.TOKEN_MUL_ASSIGN:
		return TOKEN_MUL_ASSIGN
	case lex.TOKEN_DIV_ASSIGN:
		return TOKEN_DIV_ASSIGN
	case lex.TOKEN_MOD_ASSIGN:
		return TOKEN_MOD_ASSIGN
	case lex.TOKEN_NOT:
		return TOKEN_NOT
	case lex.TOKEN_SWITCH:
		return TOKEN_SWITCH
	case lex.TOKEN_CASE:
		return TOKEN_CASE
	case lex.TOKEN_DEFAULT:
		return TOKEN_DEFAULT
	case lex.TOKEN_PACKAGE:
		return TOKEN_PACKAGE
	case lex.TOKEN_IMPORT:
		return TOKEN_IMPORT
	case lex.TOKEN_AS:
		return TOKEN_AS
	case lex.TOKEN_CLASS:
		return TOKEN_CLASS
	case lex.TOKEN_STATIC:
		return TOKEN_STATIC
	case lex.TOKEN_PUBLIC:
		return TOKEN_PUBLIC
	case lex.TOKEN_PROTECTED:
		return TOKEN_PROTECTED
	case lex.TOKEN_PRIVATE:
		return TOKEN_PRIVATE
	case lex.TOKEN_INTERFACE:
		return TOKEN_INTERFACE
	case lex.TOKEN_BYTE:
		return TOKEN_BYTE
	case lex.TOKEN_INT:
		return TOKEN_INT
	case lex.TOKEN_FLOAT:
		return TOKEN_FLOAT
	case lex.TOKEN_STRING:
		return TOKEN_STRING
	case lex.TOKEN_IDENTIFIER:
		return TOKEN_IDENTIFIER
	case lex.TOKEN_LITERAL_BOOL:
		return TOKEN_LITERAL_BOOL
	case lex.TOKEN_LITERAL_BYTE:
		return TOKEN_LITERAL_BYTE
	case lex.TOKEN_LITERAL_INT:
		return TOKEN_LITERAL_INT
	case lex.TOKEN_LITERAL_STRING:
		return TOKEN_LITERAL_STRING
	case lex.TOKEN_LITERAL_FLOAT:
		return TOKEN_LITERAL_FLOAT
	case lex.TOKEN_TRY:
		return TOKEN_TRY
	case lex.TOKEN_CATCH:
		return TOKEN_CATCH
	case lex.TOKEN_FINALLY:
		return TOKEN_FINALLY
	case lex.TOKEN_THROW:
		return TOKEN_THROW
	default:
		panic(fmt.Sprintf("unkown token %d", token.Type))
	}
	return 0
}

// The parser calls this method on a parse error.
func (x *LucyLex) Error(s string) {
	if x.Errs == nil {
		x.Errs = []error{fmt.Errorf(s)}
	} else {
		x.Errs = append(x.Errs, fmt.Errorf(s))
	}
}

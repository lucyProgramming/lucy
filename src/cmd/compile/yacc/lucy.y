%{

package yacc

import (
    "github.com/756445638/lucy/src/cmd/compile/ast"
)

%}


%token TOKEN_FUNCTION   TOKEN_CONST    TOKEN_IF      TOKEN_ELSEIF    TOKEN_ELSE   TOKEN_FOR     TOKEN_BREAK TOKEN_AS  TOKEN_STATIC
    	TOKEN_CONTINUE TOKEN_RETURN TOKEN_NULL  TOKEN_LP TOKEN_RP TOKEN_LC TOKEN_RC TOKEN_LB TOKEN_RB TOKEN_TRY  TOKEN_CATCH TOKEN_FINALLY TOKEN_THROW
    	TOKEN_SEMICOLON TOKEN_COMMA TOKEN_LOGICAL_AND TOKEN_LOGICAL_OR TOKEN_AND TOKEN_OR TOKEN_ASSIGN TOKEN_LEFT_SHIFT TOKEN_RIGHT_SHIFT
    	TOKEN_EQUAL TOKEN_NE TOKEN_GT TOKEN_GE TOKEN_LT TOKEN_LE TOKEN_ADD TOKEN_SUB TOKEN_MUL TOKEN_LITERAL_BYTE
    	TOKEN_DIV TOKEN_MOD TOKEN_INCREMENT TOKEN_DECREMENT TOKEN_DOT TOKEN_VAR TOKEN_NEW TOKEN_COLON
    	TOKEN_PLUS_ASSIGN TOKEN_MINUS_ASSIGN TOKEN_MUL_ASSIGN TOKEN_DIV_ASSIGN TOKEN_MOD_ASSIGN TOKEN_NOT
    	TOKEN_SWITCH TOKEN_CASE TOKEN_DEFAULT TOKEN_PACKAGE TOKEN_CLASS TOKEN_PUBLIC TOKEN_SKIP TOKEN_LITERAL_BOOL
    	TOKEN_PROTECTED TOKEN_PRIVATE TOKEN_BOOL TOKEN_BYTE TOKEN_INT TOKEN_FLOAT TOKEN_STRING TOKEN_ENUM TOKEN_INTERFACE
    	TOKEN_IDENTIFIER TOKEN_LITERAL_INT TOKEN_LITERAL_STRING TOKEN_LITERAL_FLOAT TOKEN_IMPORT TOKEN_COLON_ASSIGN

%token <expression> TOKEN_TRUE TOKEN_FALSE

%token <str> TOKEN_IDENTIFIER TOKEN_LITERAL_STRING
%type <names> namelist
%type <typ> typ
%type <typednames>  typednames typedname typednames_or_nil

%union {
    expression *ast.Expression
    str string
    names []string
    typ *ast.VariableType
	typednames []*ast.TypedName
	top *ast.Node
}


%%

top:
    TOKEN_PACKAGE import mainbody

package_name_definition:
    TOKEN_PACKAGE TOKEN_IDENTIFIER
    {
        packageDefinition($2)
    }

import:
    TOKEN_IMPORT TOKEN_LITERAL_STRING
    {
         importDefinition($2)
    }
    | TOKEN_IMPORT TOKEN_LITERAL_STRING TOKEN_AS TOKEN_IDENTIFIER
    {
        importDefinition($2,$4)
    }
    |
    {

    }


mainbody:
    mainbody function_definition
    {

    }
    | mainbody enum_definition
    {

    }
    |
    {

    }




enum_definition:
    TOKEN_ENUM TOKEN_IDENTIFIER TOKEN_LC  namelist  TOKEN_RC
    {

    }


statementlist:
    statementlist TOKEN_SEMICOLON statement
    {

    }



statement:



block:
    TOKEN_LC statementlist TOKEN_RC


function_definition:
    TOKEN_FUNCTION TOKEN_LP typednames TOKEN_RP  TOKEN_IDENTIFIER TOKEN_LP typednames_or_nil  TOKEN_RP block
	{

	}
	| TOKEN_FUNCTION TOKEN_IDENTIFIER TOKEN_LP typednames_or_nil  TOKEN_RP block
	{

	}

namelist:
    TOKEN_IDENTIFIER
    {
        $$ = []string{$1}
    }
    | namelist  TOKEN_COMMA TOKEN_IDENTIFIER
    {
        $$ = append($1,$3)
    }


typednames:
    typedname
    {

    }
    | typednames TOKEN_COMMA typedname
    {

    }

typednames_or_nil:
	typednames
	{
		$$ = $1
	}
	|
	{
		$$ = nil
	}


typedname:
    namelist typ
    {

    }
	| typ
	{

	}

typ:
    TOKEN_BOOL
    {

    }
    | TOKEN_BYTE
    {

    }
    | TOKEN_INT
    {

    }
    | TOKEN_FLOAT
    {

    }
    | TOKEN_STRING
    {

    }
    | TOKEN_LB TOKEN_RB typ
    {

    }












%%
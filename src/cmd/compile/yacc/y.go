//line lucy.y:2
package yacc

import __yyfmt__ "fmt"

//line lucy.y:3
import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"math/big"
)

//line lucy.y:27
type lucySymType struct {
	yys        int
	expression *ast.Expression
}

const TOKEN_FUNCTION = 57346
const TOKEN_CONST = 57347
const TOKEN_IF = 57348
const TOKEN_ELSEIF = 57349
const TOKEN_ELSE = 57350
const TOKEN_FOR = 57351
const TOKEN_BREAK = 57352
const TOKEN_CONTINUE = 57353
const TOKEN_RETURN = 57354
const TOKEN_NULL = 57355
const TOKEN_LP = 57356
const TOKEN_RP = 57357
const TOKEN_LC = 57358
const TOKEN_RC = 57359
const TOKEN_LB = 57360
const TOKEN_RB = 57361
const TOKEN_SEMICOLON = 57362
const TOKEN_COMMA = 57363
const TOKEN_LOGICAL_AND = 57364
const TOKEN_LOGICAL_OR = 57365
const TOKEN_AND = 57366
const TOKEN_OR = 57367
const TOKEN_ASSIGN = 57368
const TOKEN_EQUAL = 57369
const TOKEN_NE = 57370
const TOKEN_GT = 57371
const TOKEN_GE = 57372
const TOKEN_LT = 57373
const TOKEN_LE = 57374
const TOKEN_ADD = 57375
const TOKEN_SUB = 57376
const TOKEN_MUL = 57377
const TOKEN_DIV = 57378
const TOKEN_MOD = 57379
const TOKEN_INCREMENT = 57380
const TOKEN_DECREMENT = 57381
const TOKEN_DOT = 57382
const TOKEN_VAR = 57383
const TOKEN_NEW = 57384
const TOKEN_COLON = 57385
const TOKEN_PLUS_ASSIGN = 57386
const TOKEN_MINUS_ASSIGN = 57387
const TOKEN_MUL_ASSIGN = 57388
const TOKEN_DIV_ASSIGN = 57389
const TOKEN_MOD_ASSIGN = 57390
const TOKEN_NOT = 57391
const TOKEN_SWITCH = 57392
const TOKEN_CASE = 57393
const TOKEN_DEFAULT = 57394
const TOKEN_PACKAGE = 57395
const TOKEN_CLASS = 57396
const TOKEN_PUBLIC = 57397
const TOKEN_PROTECTED = 57398
const TOKEN_PRIVATE = 57399
const TOKEN_BOOL = 57400
const TOKEN_BYTE = 57401
const TOKEN_INT = 57402
const TOKEN_FLOAT = 57403
const TOKEN_STRING = 57404
const TOKEN_IDENTIFIER = 57405
const TOKEN_LITERAL_INT = 57406
const TOKEN_LITERAL_STRING = 57407
const TOKEN_LITERAL_FLOAT = 57408
const TOKEN_IMPORT = 57409
const TOKEN_COLON_ASSIGN = 57410
const TOKEN_TRUE = 57411
const TOKEN_FALSE = 57412

var lucyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"TOKEN_FUNCTION",
	"TOKEN_CONST",
	"TOKEN_IF",
	"TOKEN_ELSEIF",
	"TOKEN_ELSE",
	"TOKEN_FOR",
	"TOKEN_BREAK",
	"TOKEN_CONTINUE",
	"TOKEN_RETURN",
	"TOKEN_NULL",
	"TOKEN_LP",
	"TOKEN_RP",
	"TOKEN_LC",
	"TOKEN_RC",
	"TOKEN_LB",
	"TOKEN_RB",
	"TOKEN_SEMICOLON",
	"TOKEN_COMMA",
	"TOKEN_LOGICAL_AND",
	"TOKEN_LOGICAL_OR",
	"TOKEN_AND",
	"TOKEN_OR",
	"TOKEN_ASSIGN",
	"TOKEN_EQUAL",
	"TOKEN_NE",
	"TOKEN_GT",
	"TOKEN_GE",
	"TOKEN_LT",
	"TOKEN_LE",
	"TOKEN_ADD",
	"TOKEN_SUB",
	"TOKEN_MUL",
	"TOKEN_DIV",
	"TOKEN_MOD",
	"TOKEN_INCREMENT",
	"TOKEN_DECREMENT",
	"TOKEN_DOT",
	"TOKEN_VAR",
	"TOKEN_NEW",
	"TOKEN_COLON",
	"TOKEN_PLUS_ASSIGN",
	"TOKEN_MINUS_ASSIGN",
	"TOKEN_MUL_ASSIGN",
	"TOKEN_DIV_ASSIGN",
	"TOKEN_MOD_ASSIGN",
	"TOKEN_NOT",
	"TOKEN_SWITCH",
	"TOKEN_CASE",
	"TOKEN_DEFAULT",
	"TOKEN_PACKAGE",
	"TOKEN_CLASS",
	"TOKEN_PUBLIC",
	"TOKEN_PROTECTED",
	"TOKEN_PRIVATE",
	"TOKEN_BOOL",
	"TOKEN_BYTE",
	"TOKEN_INT",
	"TOKEN_FLOAT",
	"TOKEN_STRING",
	"TOKEN_IDENTIFIER",
	"TOKEN_LITERAL_INT",
	"TOKEN_LITERAL_STRING",
	"TOKEN_LITERAL_FLOAT",
	"TOKEN_IMPORT",
	"TOKEN_COLON_ASSIGN",
	"TOKEN_TRUE",
	"TOKEN_FALSE",
}
var lucyStatenames = [...]string{}

const lucyEofCode = 1
const lucyErrCode = 2
const lucyInitialStackSize = 16

//line lucy.y:43
//line yacctab:1
var lucyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const lucyPrivate = 57344

const lucyLast = 1

var lucyAct = [...]int{

	1,
}
var lucyPact = [...]int{

	-1000, -1000,
}
var lucyPgo = [...]int{

	0, 0,
}
var lucyR1 = [...]int{

	0, 1,
}
var lucyR2 = [...]int{

	0, 0,
}
var lucyChk = [...]int{

	-1000, -1,
}
var lucyDef = [...]int{

	1, -2,
}
var lucyTok1 = [...]int{

	1,
}
var lucyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 68, 69, 70,
}
var lucyTok3 = [...]int{
	0,
}

var lucyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	lucyDebug        = 0
	lucyErrorVerbose = false
)

type lucyLexer interface {
	Lex(lval *lucySymType) int
	Error(s string)
}

type lucyParser interface {
	Parse(lucyLexer) int
	Lookahead() int
}

type lucyParserImpl struct {
	lval  lucySymType
	stack [lucyInitialStackSize]lucySymType
	char  int
}

func (p *lucyParserImpl) Lookahead() int {
	return p.char
}

func lucyNewParser() lucyParser {
	return &lucyParserImpl{}
}

const lucyFlag = -1000

func lucyTokname(c int) string {
	if c >= 1 && c-1 < len(lucyToknames) {
		if lucyToknames[c-1] != "" {
			return lucyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func lucyStatname(s int) string {
	if s >= 0 && s < len(lucyStatenames) {
		if lucyStatenames[s] != "" {
			return lucyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func lucyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !lucyErrorVerbose {
		return "syntax error"
	}

	for _, e := range lucyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + lucyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := lucyPact[state]
	for tok := TOKSTART; tok-1 < len(lucyToknames); tok++ {
		if n := base + tok; n >= 0 && n < lucyLast && lucyChk[lucyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if lucyDef[state] == -2 {
		i := 0
		for lucyExca[i] != -1 || lucyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; lucyExca[i] >= 0; i += 2 {
			tok := lucyExca[i]
			if tok < TOKSTART || lucyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if lucyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += lucyTokname(tok)
	}
	return res
}

func lucylex1(lex lucyLexer, lval *lucySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = lucyTok1[0]
		goto out
	}
	if char < len(lucyTok1) {
		token = lucyTok1[char]
		goto out
	}
	if char >= lucyPrivate {
		if char < lucyPrivate+len(lucyTok2) {
			token = lucyTok2[char-lucyPrivate]
			goto out
		}
	}
	for i := 0; i < len(lucyTok3); i += 2 {
		token = lucyTok3[i+0]
		if token == char {
			token = lucyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = lucyTok2[1] /* unknown char */
	}
	if lucyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", lucyTokname(token), uint(char))
	}
	return char, token
}

func lucyParse(lucylex lucyLexer) int {
	return lucyNewParser().Parse(lucylex)
}

func (lucyrcvr *lucyParserImpl) Parse(lucylex lucyLexer) int {
	var lucyn int
	var lucyVAL lucySymType
	var lucyDollar []lucySymType
	_ = lucyDollar // silence set and not used
	lucyS := lucyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	lucystate := 0
	lucyrcvr.char = -1
	lucytoken := -1 // lucyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		lucystate = -1
		lucyrcvr.char = -1
		lucytoken = -1
	}()
	lucyp := -1
	goto lucystack

ret0:
	return 0

ret1:
	return 1

lucystack:
	/* put a state and value onto the stack */
	if lucyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", lucyTokname(lucytoken), lucyStatname(lucystate))
	}

	lucyp++
	if lucyp >= len(lucyS) {
		nyys := make([]lucySymType, len(lucyS)*2)
		copy(nyys, lucyS)
		lucyS = nyys
	}
	lucyS[lucyp] = lucyVAL
	lucyS[lucyp].yys = lucystate

lucynewstate:
	lucyn = lucyPact[lucystate]
	if lucyn <= lucyFlag {
		goto lucydefault /* simple state */
	}
	if lucyrcvr.char < 0 {
		lucyrcvr.char, lucytoken = lucylex1(lucylex, &lucyrcvr.lval)
	}
	lucyn += lucytoken
	if lucyn < 0 || lucyn >= lucyLast {
		goto lucydefault
	}
	lucyn = lucyAct[lucyn]
	if lucyChk[lucyn] == lucytoken { /* valid shift */
		lucyrcvr.char = -1
		lucytoken = -1
		lucyVAL = lucyrcvr.lval
		lucystate = lucyn
		if Errflag > 0 {
			Errflag--
		}
		goto lucystack
	}

lucydefault:
	/* default state action */
	lucyn = lucyDef[lucystate]
	if lucyn == -2 {
		if lucyrcvr.char < 0 {
			lucyrcvr.char, lucytoken = lucylex1(lucylex, &lucyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if lucyExca[xi+0] == -1 && lucyExca[xi+1] == lucystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			lucyn = lucyExca[xi+0]
			if lucyn < 0 || lucyn == lucytoken {
				break
			}
		}
		lucyn = lucyExca[xi+1]
		if lucyn < 0 {
			goto ret0
		}
	}
	if lucyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			lucylex.Error(lucyErrorMessage(lucystate, lucytoken))
			Nerrs++
			if lucyDebug >= 1 {
				__yyfmt__.Printf("%s", lucyStatname(lucystate))
				__yyfmt__.Printf(" saw %s\n", lucyTokname(lucytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for lucyp >= 0 {
				lucyn = lucyPact[lucyS[lucyp].yys] + lucyErrCode
				if lucyn >= 0 && lucyn < lucyLast {
					lucystate = lucyAct[lucyn] /* simulate a shift of "error" */
					if lucyChk[lucystate] == lucyErrCode {
						goto lucystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if lucyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", lucyS[lucyp].yys)
				}
				lucyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if lucyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", lucyTokname(lucytoken))
			}
			if lucytoken == lucyEofCode {
				goto ret1
			}
			lucyrcvr.char = -1
			lucytoken = -1
			goto lucynewstate /* try again in the same state */
		}
	}

	/* reduction by production lucyn */
	if lucyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", lucyn, lucyStatname(lucystate))
	}

	lucynt := lucyn
	lucypt := lucyp
	_ = lucypt // guard against "declared and not used"

	lucyp -= lucyR2[lucyn]
	// lucyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if lucyp+1 >= len(lucyS) {
		nyys := make([]lucySymType, len(lucyS)*2)
		copy(nyys, lucyS)
		lucyS = nyys
	}
	lucyVAL = lucyS[lucyp+1]

	/* consult goto table to find next state */
	lucyn = lucyR1[lucyn]
	lucyg := lucyPgo[lucyn]
	lucyj := lucyg + lucyS[lucyp].yys + 1

	if lucyj >= lucyLast {
		lucystate = lucyAct[lucyg]
	} else {
		lucystate = lucyAct[lucyj]
		if lucyChk[lucystate] != -lucyn {
			lucystate = lucyAct[lucyg]
		}
	}
	// dummy call; replaced with literal code
	switch lucynt {

	}
	goto lucystack /* stack new state and value */
}

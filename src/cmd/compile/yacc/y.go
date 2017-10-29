
//line lucy.y:13

package yacc
import __yyfmt__ "fmt"
//line lucy.y:14
		
import (
	"fmt"
	"math/big"
    "github.com/756445638/lucy/src/cmd/compile/lex"
)


//line lucy.y:37
type lucySymType struct {
	yys int
    boolvalue bool


}

const lex.TOKEN_FUNCTION = 57346
const lex.TOKEN_CONST = 57347
const lex.TOKEN_IF = 57348
const lex.TOKEN_ELSEIF = 57349
const lex.TOKEN_ELSE = 57350
const lex.TOKEN_FOR = 57351
const lex.TOKEN_BREAK = 57352
const lex.TOKEN_CONTINUE = 57353
const ex.TOKEN_RETURN = 57354
const ex.TOKEN_NULL = 57355
const lex.TOKEN_LP = 57356
const lex.TOKEN_RP = 57357
const lex.TOKEN_LC = 57358
const lex.TOKEN_RC = 57359
const lex.TOKEN_LB = 57360
const lex.TOKEN_RB = 57361
const lex.TOKEN_SEMICOLON = 57362
const lex.TOKEN_COMMA = 57363
const lex.TOKEN_LOGICAL_AND = 57364
const lex.TOKEN_LOGICAL_OR = 57365
const lex.TOKEN_AND = 57366
const lex.TOKEN_OR = 57367
const lex.TOKEN_ASSIGN = 57368
const lex.TOKEN_EQUAL = 57369
const lex.TOKEN_NE = 57370
const lex.TOKEN_GT = 57371
const lex.TOKEN_GE = 57372
const lex.TOKEN_LT = 57373
const lex.TOKEN_LE = 57374
const lex.TOKEN_ADD = 57375
const lex.TOKEN_SUB = 57376
const lex.TOKEN_MUL = 57377
const lex.TOKEN_DIV = 57378
const lex.TOKEN_MOD = 57379
const lex.TOKEN_INCREMENT = 57380
const lex.TOKEN_DECREMENT = 57381
const lex.TOKEN_DOT = 57382
const lex.TOKEN_VAR = 57383
const lex.TOKEN_NEW = 57384
const lex.TOKEN_COLON = 57385
const lex.TOKEN_PLUS_ASSIGN = 57386
const lex.TOKEN_MINUS_ASSIGN = 57387
const lex.TOKEN_MUL_ASSIGN = 57388
const lex.TOKEN_DIV_ASSIGN = 57389
const lex.TOKEN_MOD_ASSIGN = 57390
const lex.TOKEN_NOT = 57391
const lex.TOKEN_SWITCH = 57392
const lex.TOKEN_CASE = 57393
const lex.TOKEN_DEFAULT = 57394
const lex.TOKEN_CRLF = 57395
const lex.TOKEN_PACKAGE = 57396
const lex.TOKEN_CLASS = 57397
const lex.TOKEN_PUBLIC = 57398
const lex.TOKEN_PROTECTED = 57399
const lex.TOKEN_PRIVATE = 57400
const lex.TOKEN_BOOL = 57401
const lex.TOKEN_BYTE = 57402
const lex.TOKEN_INT = 57403
const lex.TOKEN_FLOAT = 57404
const lex.TOKEN_STRING = 57405
const lex.TOKEN_IDENTIFIER = 57406
const lex.TOKEN_LITERAL_INT = 57407
const lex.TOKEN_LITERAL_STRING = 57408
const lex.TOKEN_LITERAL_FLOAT = 57409
const lex.TOKEN_TRUE = 57410
const lex.TOKEN_FALSE = 57411

var lucyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"lex.TOKEN_FUNCTION",
	"lex.TOKEN_CONST",
	"lex.TOKEN_IF",
	"lex.TOKEN_ELSEIF",
	"lex.TOKEN_ELSE",
	"lex.TOKEN_FOR",
	"lex.TOKEN_BREAK",
	"lex.TOKEN_CONTINUE",
	"ex.TOKEN_RETURN",
	"ex.TOKEN_NULL",
	"lex.TOKEN_LP",
	"lex.TOKEN_RP",
	"lex.TOKEN_LC",
	"lex.TOKEN_RC",
	"lex.TOKEN_LB",
	"lex.TOKEN_RB",
	"lex.TOKEN_SEMICOLON",
	"lex.TOKEN_COMMA",
	"lex.TOKEN_LOGICAL_AND",
	"lex.TOKEN_LOGICAL_OR",
	"lex.TOKEN_AND",
	"lex.TOKEN_OR",
	"lex.TOKEN_ASSIGN",
	"lex.TOKEN_EQUAL",
	"lex.TOKEN_NE",
	"lex.TOKEN_GT",
	"lex.TOKEN_GE",
	"lex.TOKEN_LT",
	"lex.TOKEN_LE",
	"lex.TOKEN_ADD",
	"lex.TOKEN_SUB",
	"lex.TOKEN_MUL",
	"lex.TOKEN_DIV",
	"lex.TOKEN_MOD",
	"lex.TOKEN_INCREMENT",
	"lex.TOKEN_DECREMENT",
	"lex.TOKEN_DOT",
	"lex.TOKEN_VAR",
	"lex.TOKEN_NEW",
	"lex.TOKEN_COLON",
	"lex.TOKEN_PLUS_ASSIGN",
	"lex.TOKEN_MINUS_ASSIGN",
	"lex.TOKEN_MUL_ASSIGN",
	"lex.TOKEN_DIV_ASSIGN",
	"lex.TOKEN_MOD_ASSIGN",
	"lex.TOKEN_NOT",
	"lex.TOKEN_SWITCH",
	"lex.TOKEN_CASE",
	"lex.TOKEN_DEFAULT",
	"lex.TOKEN_CRLF",
	"lex.TOKEN_PACKAGE",
	"lex.TOKEN_CLASS",
	"lex.TOKEN_PUBLIC",
	"lex.TOKEN_PROTECTED",
	"lex.TOKEN_PRIVATE",
	"lex.TOKEN_BOOL",
	"lex.TOKEN_BYTE",
	"lex.TOKEN_INT",
	"lex.TOKEN_FLOAT",
	"lex.TOKEN_STRING",
	"lex.TOKEN_IDENTIFIER",
	"lex.TOKEN_LITERAL_INT",
	"lex.TOKEN_LITERAL_STRING",
	"lex.TOKEN_LITERAL_FLOAT",
	"lex.TOKEN_TRUE",
	"lex.TOKEN_FALSE",
}
var lucyStatenames = [...]string{}

const lucyEofCode = 1
const lucyErrCode = 2
const lucyInitialStackSize = 16

//line lucy.y:55

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
	62, 63, 64, 65, 66, 67, 68, 69,
}
var lucyTok3 = [...]int{
	0,
}

var lucyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{
}

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

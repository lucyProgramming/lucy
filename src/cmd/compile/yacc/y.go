//line lucy.y:2
package yacc

import __yyfmt__ "fmt"

//line lucy.y:3
import (
	"github.com/756445638/lucy/src/cmd/compile/ast"
)

//line lucy.y:29
type LucySymType struct {
	yys        int
	expression *ast.Expression
	str        string
	names      []string
	typ        *ast.VariableType
	typednames []*ast.TypedName
	top        *ast.Node
}

const TOKEN_FUNCTION = 57346
const TOKEN_CONST = 57347
const TOKEN_IF = 57348
const TOKEN_ELSEIF = 57349
const TOKEN_ELSE = 57350
const TOKEN_FOR = 57351
const TOKEN_BREAK = 57352
const TOKEN_AS = 57353
const TOKEN_STATIC = 57354
const TOKEN_CONTINUE = 57355
const TOKEN_RETURN = 57356
const TOKEN_NULL = 57357
const TOKEN_LP = 57358
const TOKEN_RP = 57359
const TOKEN_LC = 57360
const TOKEN_RC = 57361
const TOKEN_LB = 57362
const TOKEN_RB = 57363
const TOKEN_TRY = 57364
const TOKEN_CATCH = 57365
const TOKEN_FINALLY = 57366
const TOKEN_THROW = 57367
const TOKEN_SEMICOLON = 57368
const TOKEN_COMMA = 57369
const TOKEN_LOGICAL_AND = 57370
const TOKEN_LOGICAL_OR = 57371
const TOKEN_AND = 57372
const TOKEN_OR = 57373
const TOKEN_ASSIGN = 57374
const TOKEN_LEFT_SHIFT = 57375
const TOKEN_RIGHT_SHIFT = 57376
const TOKEN_EQUAL = 57377
const TOKEN_NE = 57378
const TOKEN_GT = 57379
const TOKEN_GE = 57380
const TOKEN_LT = 57381
const TOKEN_LE = 57382
const TOKEN_ADD = 57383
const TOKEN_SUB = 57384
const TOKEN_MUL = 57385
const TOKEN_LITERAL_BYTE = 57386
const TOKEN_DIV = 57387
const TOKEN_MOD = 57388
const TOKEN_INCREMENT = 57389
const TOKEN_DECREMENT = 57390
const TOKEN_DOT = 57391
const TOKEN_VAR = 57392
const TOKEN_NEW = 57393
const TOKEN_COLON = 57394
const TOKEN_PLUS_ASSIGN = 57395
const TOKEN_MINUS_ASSIGN = 57396
const TOKEN_MUL_ASSIGN = 57397
const TOKEN_DIV_ASSIGN = 57398
const TOKEN_MOD_ASSIGN = 57399
const TOKEN_NOT = 57400
const TOKEN_SWITCH = 57401
const TOKEN_CASE = 57402
const TOKEN_DEFAULT = 57403
const TOKEN_PACKAGE = 57404
const TOKEN_CLASS = 57405
const TOKEN_PUBLIC = 57406
const TOKEN_SKIP = 57407
const TOKEN_LITERAL_BOOL = 57408
const TOKEN_PROTECTED = 57409
const TOKEN_PRIVATE = 57410
const TOKEN_BOOL = 57411
const TOKEN_BYTE = 57412
const TOKEN_INT = 57413
const TOKEN_FLOAT = 57414
const TOKEN_STRING = 57415
const TOKEN_ENUM = 57416
const TOKEN_INTERFACE = 57417
const TOKEN_IDENTIFIER = 57418
const TOKEN_LITERAL_INT = 57419
const TOKEN_LITERAL_STRING = 57420
const TOKEN_LITERAL_FLOAT = 57421
const TOKEN_IMPORT = 57422
const TOKEN_COLON_ASSIGN = 57423
const TOKEN_TRUE = 57424
const TOKEN_FALSE = 57425

var LucyToknames = [...]string{
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
	"TOKEN_AS",
	"TOKEN_STATIC",
	"TOKEN_CONTINUE",
	"TOKEN_RETURN",
	"TOKEN_NULL",
	"TOKEN_LP",
	"TOKEN_RP",
	"TOKEN_LC",
	"TOKEN_RC",
	"TOKEN_LB",
	"TOKEN_RB",
	"TOKEN_TRY",
	"TOKEN_CATCH",
	"TOKEN_FINALLY",
	"TOKEN_THROW",
	"TOKEN_SEMICOLON",
	"TOKEN_COMMA",
	"TOKEN_LOGICAL_AND",
	"TOKEN_LOGICAL_OR",
	"TOKEN_AND",
	"TOKEN_OR",
	"TOKEN_ASSIGN",
	"TOKEN_LEFT_SHIFT",
	"TOKEN_RIGHT_SHIFT",
	"TOKEN_EQUAL",
	"TOKEN_NE",
	"TOKEN_GT",
	"TOKEN_GE",
	"TOKEN_LT",
	"TOKEN_LE",
	"TOKEN_ADD",
	"TOKEN_SUB",
	"TOKEN_MUL",
	"TOKEN_LITERAL_BYTE",
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
	"TOKEN_SKIP",
	"TOKEN_LITERAL_BOOL",
	"TOKEN_PROTECTED",
	"TOKEN_PRIVATE",
	"TOKEN_BOOL",
	"TOKEN_BYTE",
	"TOKEN_INT",
	"TOKEN_FLOAT",
	"TOKEN_STRING",
	"TOKEN_ENUM",
	"TOKEN_INTERFACE",
	"TOKEN_IDENTIFIER",
	"TOKEN_LITERAL_INT",
	"TOKEN_LITERAL_STRING",
	"TOKEN_LITERAL_FLOAT",
	"TOKEN_IMPORT",
	"TOKEN_COLON_ASSIGN",
	"TOKEN_TRUE",
	"TOKEN_FALSE",
}
var LucyStatenames = [...]string{}

const LucyEofCode = 1
const LucyErrCode = 2
const LucyInitialStackSize = 16

//line lucy.y:149
//line yacctab:1
var LucyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const LucyPrivate = 57344

const LucyLast = 64

var LucyAct = [...]int{

	5, 20, 10, 11, 32, 30, 22, 6, 3, 20,
	20, 28, 36, 17, 23, 24, 25, 15, 38, 27,
	34, 35, 21, 29, 24, 13, 9, 8, 7, 4,
	26, 2, 1, 16, 0, 14, 0, 0, 0, 0,
	0, 33, 31, 0, 0, 0, 0, 37, 0, 0,
	19, 39, 0, 0, 0, 0, 0, 18, 19, 19,
	0, 0, 0, 12,
}
var LucyPact = [...]int{

	-54, -1000, -80, -69, 22, -76, -1000, -1000, -1000, -13,
	14, -19, 6, -70, -3, -1000, -11, -1000, -1000, -1000,
	-2, -19, -1000, -71, -19, -72, -1000, -10, 3, -12,
	5, -1000, -1000, -1000, -1000, -19, -1000, 1, -1000, -1000,
}
var LucyPgo = [...]int{

	0, 33, 13, 23, 17, 11, 32, 31, 29, 28,
	27, 12,
}
var LucyR1 = [...]int{

	0, 6, 7, 8, 8, 8, 9, 11, 10, 10,
	1, 1, 3, 3, 5, 5, 4, 4, 2, 2,
}
var LucyR2 = [...]int{

	0, 3, 2, 2, 4, 0, 1, 0, 9, 6,
	1, 3, 1, 3, 1, 0, 2, 1, 1, 3,
}
var LucyChk = [...]int{

	-1000, -6, -7, 62, -8, 80, 76, -9, -10, 4,
	78, 16, 76, 11, -3, -4, -1, -2, 76, 69,
	20, 16, 76, 17, 27, 27, -2, 21, -5, -3,
	76, -4, 76, -2, 17, 16, -11, -5, 17, -11,
}
var LucyDef = [...]int{

	0, -2, 5, 0, 0, 0, 2, 1, 6, 0,
	3, 0, 0, 0, 0, 12, 0, 17, 10, 18,
	0, 15, 4, 0, 0, 0, 16, 0, 0, 14,
	0, 13, 11, 19, 7, 15, 9, 0, 7, 8,
}
var LucyTok1 = [...]int{

	1,
}
var LucyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 68, 69, 70, 71,
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81,
	82, 83,
}
var LucyTok3 = [...]int{
	0,
}

var LucyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	LucyDebug        = 0
	LucyErrorVerbose = false
)

type LucyLexer interface {
	Lex(lval *LucySymType) int
	Error(s string)
}

type LucyParser interface {
	Parse(LucyLexer) int
	Lookahead() int
}

type LucyParserImpl struct {
	lval  LucySymType
	stack [LucyInitialStackSize]LucySymType
	char  int
}

func (p *LucyParserImpl) Lookahead() int {
	return p.char
}

func LucyNewParser() LucyParser {
	return &LucyParserImpl{}
}

const LucyFlag = -1000

func LucyTokname(c int) string {
	if c >= 1 && c-1 < len(LucyToknames) {
		if LucyToknames[c-1] != "" {
			return LucyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func LucyStatname(s int) string {
	if s >= 0 && s < len(LucyStatenames) {
		if LucyStatenames[s] != "" {
			return LucyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func LucyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !LucyErrorVerbose {
		return "syntax error"
	}

	for _, e := range LucyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + LucyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := LucyPact[state]
	for tok := TOKSTART; tok-1 < len(LucyToknames); tok++ {
		if n := base + tok; n >= 0 && n < LucyLast && LucyChk[LucyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if LucyDef[state] == -2 {
		i := 0
		for LucyExca[i] != -1 || LucyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; LucyExca[i] >= 0; i += 2 {
			tok := LucyExca[i]
			if tok < TOKSTART || LucyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if LucyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += LucyTokname(tok)
	}
	return res
}

func Lucylex1(lex LucyLexer, lval *LucySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = LucyTok1[0]
		goto out
	}
	if char < len(LucyTok1) {
		token = LucyTok1[char]
		goto out
	}
	if char >= LucyPrivate {
		if char < LucyPrivate+len(LucyTok2) {
			token = LucyTok2[char-LucyPrivate]
			goto out
		}
	}
	for i := 0; i < len(LucyTok3); i += 2 {
		token = LucyTok3[i+0]
		if token == char {
			token = LucyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = LucyTok2[1] /* unknown char */
	}
	if LucyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", LucyTokname(token), uint(char))
	}
	return char, token
}

func LucyParse(Lucylex LucyLexer) int {
	return LucyNewParser().Parse(Lucylex)
}

func (Lucyrcvr *LucyParserImpl) Parse(Lucylex LucyLexer) int {
	var Lucyn int
	var LucyVAL LucySymType
	var LucyDollar []LucySymType
	_ = LucyDollar // silence set and not used
	LucyS := Lucyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	Lucystate := 0
	Lucyrcvr.char = -1
	Lucytoken := -1 // Lucyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		Lucystate = -1
		Lucyrcvr.char = -1
		Lucytoken = -1
	}()
	Lucyp := -1
	goto Lucystack

ret0:
	return 0

ret1:
	return 1

Lucystack:
	/* put a state and value onto the stack */
	if LucyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", LucyTokname(Lucytoken), LucyStatname(Lucystate))
	}

	Lucyp++
	if Lucyp >= len(LucyS) {
		nyys := make([]LucySymType, len(LucyS)*2)
		copy(nyys, LucyS)
		LucyS = nyys
	}
	LucyS[Lucyp] = LucyVAL
	LucyS[Lucyp].yys = Lucystate

Lucynewstate:
	Lucyn = LucyPact[Lucystate]
	if Lucyn <= LucyFlag {
		goto Lucydefault /* simple state */
	}
	if Lucyrcvr.char < 0 {
		Lucyrcvr.char, Lucytoken = Lucylex1(Lucylex, &Lucyrcvr.lval)
	}
	Lucyn += Lucytoken
	if Lucyn < 0 || Lucyn >= LucyLast {
		goto Lucydefault
	}
	Lucyn = LucyAct[Lucyn]
	if LucyChk[Lucyn] == Lucytoken { /* valid shift */
		Lucyrcvr.char = -1
		Lucytoken = -1
		LucyVAL = Lucyrcvr.lval
		Lucystate = Lucyn
		if Errflag > 0 {
			Errflag--
		}
		goto Lucystack
	}

Lucydefault:
	/* default state action */
	Lucyn = LucyDef[Lucystate]
	if Lucyn == -2 {
		if Lucyrcvr.char < 0 {
			Lucyrcvr.char, Lucytoken = Lucylex1(Lucylex, &Lucyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if LucyExca[xi+0] == -1 && LucyExca[xi+1] == Lucystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			Lucyn = LucyExca[xi+0]
			if Lucyn < 0 || Lucyn == Lucytoken {
				break
			}
		}
		Lucyn = LucyExca[xi+1]
		if Lucyn < 0 {
			goto ret0
		}
	}
	if Lucyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			Lucylex.Error(LucyErrorMessage(Lucystate, Lucytoken))
			Nerrs++
			if LucyDebug >= 1 {
				__yyfmt__.Printf("%s", LucyStatname(Lucystate))
				__yyfmt__.Printf(" saw %s\n", LucyTokname(Lucytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for Lucyp >= 0 {
				Lucyn = LucyPact[LucyS[Lucyp].yys] + LucyErrCode
				if Lucyn >= 0 && Lucyn < LucyLast {
					Lucystate = LucyAct[Lucyn] /* simulate a shift of "error" */
					if LucyChk[Lucystate] == LucyErrCode {
						goto Lucystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if LucyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", LucyS[Lucyp].yys)
				}
				Lucyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if LucyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", LucyTokname(Lucytoken))
			}
			if Lucytoken == LucyEofCode {
				goto ret1
			}
			Lucyrcvr.char = -1
			Lucytoken = -1
			goto Lucynewstate /* try again in the same state */
		}
	}

	/* reduction by production Lucyn */
	if LucyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", Lucyn, LucyStatname(Lucystate))
	}

	Lucynt := Lucyn
	Lucypt := Lucyp
	_ = Lucypt // guard against "declared and not used"

	Lucyp -= LucyR2[Lucyn]
	// Lucyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if Lucyp+1 >= len(LucyS) {
		nyys := make([]LucySymType, len(LucyS)*2)
		copy(nyys, LucyS)
		LucyS = nyys
	}
	LucyVAL = LucyS[Lucyp+1]

	/* consult goto table to find next state */
	Lucyn = LucyR1[Lucyn]
	Lucyg := LucyPgo[Lucyn]
	Lucyj := Lucyg + LucyS[Lucyp].yys + 1

	if Lucyj >= LucyLast {
		Lucystate = LucyAct[Lucyg]
	} else {
		Lucystate = LucyAct[Lucyj]
		if LucyChk[Lucystate] != -Lucyn {
			Lucystate = LucyAct[Lucyg]
		}
	}
	// dummy call; replaced with literal code
	switch Lucynt {

	case 2:
		LucyDollar = LucyS[Lucypt-2 : Lucypt+1]
		//line lucy.y:46
		{
			packageDefination(LucyDollar[2].str)
		}
	case 3:
		LucyDollar = LucyS[Lucypt-2 : Lucypt+1]
		//line lucy.y:52
		{
			importDefination(LucyDollar[2].str)
		}
	case 4:
		LucyDollar = LucyS[Lucypt-4 : Lucypt+1]
		//line lucy.y:56
		{
			importDefination(LucyDollar[2].str, LucyDollar[4].str)
		}
	case 5:
		LucyDollar = LucyS[Lucypt-0 : Lucypt+1]
		//line lucy.y:60
		{

		}
	case 6:
		LucyDollar = LucyS[Lucypt-1 : Lucypt+1]
		//line lucy.y:68
		{

		}
	case 8:
		LucyDollar = LucyS[Lucypt-9 : Lucypt+1]
		//line lucy.y:77
		{

		}
	case 9:
		LucyDollar = LucyS[Lucypt-6 : Lucypt+1]
		//line lucy.y:81
		{

		}
	case 10:
		LucyDollar = LucyS[Lucypt-1 : Lucypt+1]
		//line lucy.y:87
		{
			LucyVAL.names = []string{LucyDollar[1].str}
		}
	case 11:
		LucyDollar = LucyS[Lucypt-3 : Lucypt+1]
		//line lucy.y:91
		{
			LucyVAL.names = append(LucyDollar[1].names, LucyDollar[3].str)
		}
	case 12:
		LucyDollar = LucyS[Lucypt-1 : Lucypt+1]
		//line lucy.y:98
		{

		}
	case 13:
		LucyDollar = LucyS[Lucypt-3 : Lucypt+1]
		//line lucy.y:102
		{

		}
	case 14:
		LucyDollar = LucyS[Lucypt-1 : Lucypt+1]
		//line lucy.y:108
		{
			LucyVAL.typednames = LucyDollar[1].typednames
		}
	case 15:
		LucyDollar = LucyS[Lucypt-0 : Lucypt+1]
		//line lucy.y:112
		{
			LucyVAL.typednames = nil
		}
	case 16:
		LucyDollar = LucyS[Lucypt-2 : Lucypt+1]
		//line lucy.y:119
		{

		}
	case 17:
		LucyDollar = LucyS[Lucypt-1 : Lucypt+1]
		//line lucy.y:123
		{

		}
	case 18:
		LucyDollar = LucyS[Lucypt-1 : Lucypt+1]
		//line lucy.y:129
		{

		}
	case 19:
		LucyDollar = LucyS[Lucypt-3 : Lucypt+1]
		//line lucy.y:133
		{

		}
	}
	goto Lucystack /* stack new state and value */
}

//line src/github.com/sayotte/iscdhcp/parse.y:2
package iscdhcp

import __yyfmt__ "fmt"

//line src/github.com/sayotte/iscdhcp/parse.y:2
import (
	"fmt"
	"net"
	"strings"
)

//line src/github.com/sayotte/iscdhcp/parse.y:43
type yySymType struct {
	yys             int
	num             int
	str             string
	strList         []string
	ipList          []net.IP
	ip              net.IP
	statement       Statement
	statementList   []Statement
	dataTerm        fmt.Stringer
	boolExpr        BooleanExpression
	subConditionals []ConditionalStatement
}

const openBrace = 57346
const closeBrace = 57347
const quote = 57348
const semicolon = 57349
const comma = 57350
const ConditionIf = 57351
const ConditionElsif = 57352
const ConditionElse = 57353
const BoolAnd = 57354
const BoolOr = 57355
const BoolNot = 57356
const BoolEqual = 57357
const BoolInequal = 57358
const BoolRegexMatch = 57359
const BoolRegexIMatch = 57360
const BoolExists = 57361
const BoolKnown = 57362
const BoolStatic = 57363
const number = 57364
const ipAddr = 57365
const cidr = 57366
const stringConst = 57367
const macAddr = 57368
const stateTok = 57369
const authoritativeTok = 57370
const groupTok = 57371
const hostTok = 57372
const subnetTok = 57373
const netmaskTok = 57374
const optionTok = 57375
const includeTok = 57376
const hardwareTok = 57377
const ethernetTok = 57378
const fixedAddrTok = 57379
const useHostDeclNamesTok = 57380
const optDomainNameServersTok = 57381
const word = 57382

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"openBrace",
	"closeBrace",
	"quote",
	"semicolon",
	"comma",
	"ConditionIf",
	"ConditionElsif",
	"ConditionElse",
	"BoolAnd",
	"BoolOr",
	"BoolNot",
	"BoolEqual",
	"BoolInequal",
	"BoolRegexMatch",
	"BoolRegexIMatch",
	"BoolExists",
	"BoolKnown",
	"BoolStatic",
	"number",
	"ipAddr",
	"cidr",
	"stringConst",
	"macAddr",
	"stateTok",
	"authoritativeTok",
	"groupTok",
	"hostTok",
	"subnetTok",
	"netmaskTok",
	"optionTok",
	"includeTok",
	"hardwareTok",
	"ethernetTok",
	"fixedAddrTok",
	"useHostDeclNamesTok",
	"optDomainNameServersTok",
	"word",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line src/github.com/sayotte/iscdhcp/parse.y:318

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 125

var yyAct = [...]int{

	26, 36, 3, 31, 42, 25, 62, 28, 46, 41,
	52, 32, 37, 39, 2, 69, 35, 34, 33, 18,
	38, 47, 37, 29, 19, 64, 79, 70, 43, 50,
	38, 30, 53, 54, 55, 78, 56, 57, 20, 14,
	15, 17, 49, 23, 16, 21, 27, 22, 24, 82,
	83, 67, 25, 58, 59, 60, 61, 27, 72, 73,
	74, 75, 76, 77, 48, 54, 55, 68, 18, 80,
	65, 81, 63, 19, 66, 65, 51, 40, 45, 44,
	71, 13, 12, 11, 85, 86, 84, 20, 14, 15,
	17, 10, 23, 16, 21, 18, 22, 24, 9, 8,
	19, 7, 6, 5, 4, 1, 0, 0, 0, 0,
	0, 0, 0, 0, 20, 14, 15, 17, 0, 23,
	16, 21, 0, 22, 24,
}
var yyPact = [...]int{

	86, -1000, 86, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 42, -33, -2, 8, -3, -15,
	70, -27, 5, -31, -6, -1000, -1000, 59, 42, 69,
	-22, 53, -3, -1000, -1000, -13, 38, -1000, -34, 65,
	-1000, -1, 67, -1000, -1000, -1000, 5, 60, -1000, 10,
	-1000, -1000, 4, -1000, -3, -3, 21, -1000, -13, -13,
	-13, -13, -1000, -1000, 28, 3, -1000, 62, -1000, -1000,
	42, 39, 21, 21, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -3, 42, 53, -1000, -1000,
}
var yyPgo = [...]int{

	0, 105, 14, 2, 104, 103, 102, 101, 99, 98,
	91, 83, 82, 81, 0, 3, 80, 1, 4, 79,
	78,
}
var yyR1 = [...]int{

	0, 1, 2, 2, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 14, 14, 8, 16, 16, 16,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	17, 17, 18, 18, 4, 5, 6, 7, 9, 9,
	10, 11, 13, 12, 19, 20,
}
var yyR2 = [...]int{

	0, 1, 1, 2, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 2, 3, 4, 0, 4, 3,
	3, 3, 2, 1, 1, 2, 3, 3, 3, 3,
	1, 2, 3, 1, 2, 3, 3, 5, 3, 2,
	4, 3, 3, 2, 1, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -5, -6, -7, -8, -9,
	-10, -11, -12, -13, 29, 30, 34, 31, 9, 14,
	28, 35, 37, 33, 38, -3, -14, 4, 40, 25,
	23, -15, 14, 21, 20, 19, -17, 25, 33, 28,
	7, 36, -18, 23, -19, -20, 39, 27, 5, -2,
	-14, 7, 32, -14, 12, 13, -15, -17, 15, 16,
	17, 18, 40, 7, 26, 8, 7, -18, 7, 5,
	23, -16, -15, -15, -17, -17, -17, -17, 7, 23,
	7, -14, 10, 11, -15, -14, -14,
}
var yyDef = [...]int{

	0, -2, 1, 2, 4, 5, 6, 7, 8, 9,
	10, 11, 12, 13, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 3, 34, 0, 0, 0,
	0, 0, 0, 23, 24, 0, 0, 30, 0, 0,
	39, 0, 0, 33, 43, 44, 0, 0, 14, 0,
	35, 36, 0, 17, 0, 0, 22, 25, 0, 0,
	0, 0, 31, 38, 0, 0, 41, 0, 42, 15,
	0, 16, 20, 21, 26, 27, 28, 29, 40, 32,
	45, 37, 0, 0, 0, 19, 18,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:59
		{
			l := yylex.(*lexer)
			l.dirtyHackReturn = yyDollar[1].statementList
		}
	case 2:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:66
		{
			yyVAL.statementList = []Statement{yyDollar[1].statement}
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:70
		{
			yyVAL.statementList = append(yyVAL.statementList, yyDollar[2].statement)
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:93
		{
			yyVAL.statementList = yyDollar[2].statementList
		}
	case 16:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:100
		{
			cs := ConditionalStatement{
				Operator:        ConditionIf,
				Condition:       yyDollar[2].boolExpr,
				Statements:      yyDollar[3].statementList,
				SubConditionals: yyDollar[4].subConditionals,
			}
			yyVAL.statement = cs
		}
	case 18:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:113
		{
			cs := ConditionalStatement{
				Operator:   ConditionElsif,
				Condition:  yyDollar[3].boolExpr,
				Statements: yyDollar[3].statementList,
			}
			yyVAL.subConditionals = append(yyVAL.subConditionals, cs)
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:122
		{
			cs := ConditionalStatement{
				Operator:   ConditionElse,
				Statements: yyDollar[2].statementList,
			}
			yyVAL.subConditionals = append(yyVAL.subConditionals, cs)
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:132
		{
			yyVAL.boolExpr = BooleanExpression{
				Operator:  BoolAnd,
				BoolTerms: []BooleanExpression{yyDollar[1].boolExpr, yyDollar[3].boolExpr},
			}
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:139
		{
			yyVAL.boolExpr = BooleanExpression{
				Operator:  BoolOr,
				BoolTerms: []BooleanExpression{yyDollar[1].boolExpr, yyDollar[3].boolExpr},
			}
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:146
		{
			yyVAL.boolExpr = BooleanExpression{
				Operator:  BoolNot,
				BoolTerms: []BooleanExpression{yyDollar[2].boolExpr},
			}
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:153
		{
			yyVAL.boolExpr = BooleanExpression{
				Operator: BoolStatic,
			}
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:159
		{
			yyVAL.boolExpr = BooleanExpression{
				Operator: BoolKnown,
			}
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:165
		{
			yyVAL.boolExpr = BooleanExpression{
				Operator:  BoolExists,
				DataTerms: []fmt.Stringer{yyDollar[2].dataTerm},
			}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:172
		{
			yyVAL.boolExpr = BooleanExpression{
				Operator:  BoolEqual,
				DataTerms: []fmt.Stringer{yyDollar[1].dataTerm, yyDollar[3].dataTerm},
			}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:179
		{
			yyVAL.boolExpr = BooleanExpression{
				Operator:  BoolInequal,
				DataTerms: []fmt.Stringer{yyDollar[1].dataTerm, yyDollar[3].dataTerm},
			}
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:186
		{
			yyVAL.boolExpr = BooleanExpression{
				Operator:  BoolRegexMatch,
				DataTerms: []fmt.Stringer{yyDollar[1].dataTerm, yyDollar[3].dataTerm},
			}
		}
	case 29:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:193
		{
			yyVAL.boolExpr = BooleanExpression{
				Operator:  BoolRegexIMatch,
				DataTerms: []fmt.Stringer{yyDollar[1].dataTerm, yyDollar[3].dataTerm},
			}
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:202
		{
			yyVAL.dataTerm = StringConstTerm(yyDollar[1].str)
		}
	case 31:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:206
		{
			yyVAL.dataTerm = PacketOptionTerm{
				optionName: yyDollar[2].str,
			}
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:224
		{
			yyVAL.ipList = append(yyVAL.ipList, net.ParseIP(yyDollar[3].str))
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:228
		{
			yyVAL.ipList = []net.IP{net.ParseIP(yyDollar[1].str)}
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:234
		{
			gs := GroupStatement{
				Statements: yyDollar[2].statementList,
			}
			yyVAL.statement = gs
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:242
		{
			hs := HostStatement{
				Hostname:   yyDollar[2].str,
				Statements: yyDollar[3].statementList,
			}
			yyVAL.statement = hs
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:251
		{
			is := IncludeStatement{
				Filename: yyDollar[2].str,
			}
			yyVAL.statement = is
		}
	case 37:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:259
		{
			sns := SubnetStatement{
				SubnetNumber: net.ParseIP(yyDollar[2].str),
				Netmask:      net.ParseIP(yyDollar[4].str),
				Statements:   yyDollar[5].statementList,
			}
			yyVAL.statement = sns
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:271
		{
			yyVAL.statement = AuthoritativeStatement(false)
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:275
		{
			yyVAL.statement = AuthoritativeStatement(true)
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:281
		{
			yyVAL.statement = HardwareStatement{
				HardwareType:    "ethernet",
				HardwareAddress: yyDollar[3].str,
			}
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:290
		{
			yyVAL.statement = FixedAddressStatement(yyDollar[2].ipList)
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:296
		{
			val := false
			cmpText := strings.ToLower(yyDollar[2].str)
			if cmpText == "on" || cmpText == "true" {
				val = true
			}
			yyVAL.statement = UseHostDeclNamesStatement(val)
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:307
		{
			yyVAL.statement = yyDollar[2].statement
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line src/github.com/sayotte/iscdhcp/parse.y:315
		{
			yyVAL.statement = DomainNameServersOption(yyDollar[2].ipList)
		}
	}
	goto yystack /* stack new state and value */
}

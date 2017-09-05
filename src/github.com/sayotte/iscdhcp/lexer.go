package iscdhcp

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type token struct {
	data []byte
	typ  int
}

func (t token) String() string {
	return fmt.Sprintf("token{data: []byte{%q}, typ: %d}", string(t.data), t.typ)
}

var cidrRegexp = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\/\d{1,2}$`)
var ipAddrRegexp = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
var macAddrRegexp = regexp.MustCompile(`^[a-fA-F0-9]{1,2}:[a-fA-F0-9]{1,2}:[a-fA-F0-9]{1,2}:[a-fA-F0-9]{1,2}:[a-fA-F0-9]{1,2}:[a-fA-F0-9]{1,2}$`)
var numRegexp = regexp.MustCompile(`\d+`)

var stringTokenMap = map[string]int{
	// declarations
	"group":   groupTok,
	"host":    hostTok,
	"subnet":  subnetTok,
	"netmask": netmaskTok,
	// parameters
	"authoritative":       authoritativeTok,
	"domain-name-servers": optDomainNameServersTok,
	"ethernet":            ethernetTok,
	"fixed-address":       fixedAddrTok,
	"hardware":            hardwareTok,
	"include":             includeTok,
	"option":              optionTok,
	"use-host-decl-names": useHostDeclNamesTok,
	// parameter states
	"on":    stateTok,
	"off":   stateTok,
	"true":  stateTok,
	"false": stateTok,
	// boolean comparison operators
	"if":    ConditionIf,
	"elsif": ConditionElsif,
	"else":  ConditionElse,
	"and":   BoolAnd,
	"or":    BoolOr,
	"not":   BoolNot,
	// data comparison operators
	"=":  BoolEqual,
	"!=": BoolInequal,
	"~=": BoolRegexMatch,
	"~~": BoolRegexIMatch,
	// boolean simple operators
	"exists": BoolExists,
	"known":  BoolKnown,
	"static": BoolStatic,
}

const (
	tokenTypeIdentifier = iota
	tokenTypeString
	tokenTypeWhiteSpace
	tokenTypeComment
	tokenTypeBlockStart
	tokenTypeBlockEnd
	tokenTypeSemicolon
	tokenTypeComma
)

func newLexer(r io.Reader) *lexer {
	t := &lexer{
		dataStream: r,
		scanner:    &scanner{},
	}
	t.scanner.init()
	return t
}

type lexer struct {
	dataStream      io.Reader
	scanner         *scanner
	currentToken    token
	wipToken        token
	dirtyHackReturn []Statement
	hasErrored      bool
	err             error
}

func (l *lexer) Error(s string) {
	// This method is called by yyParser.Parse() when the yacc-generated
	// parser hits a snag. It doesn't give us much context.
	l.err = fmt.Errorf("error at (or just before) token %q: %s", l.currentToken.data, s)
	l.hasErrored = true
}

func (l *lexer) Lex(lval *yySymType) int {
	for {
		tok, err := l.nextToken()
		if err != nil && err != io.EOF {
			l.err = err
			l.hasErrored = true
			return 0
		}
		if len(tok.data) == 0 {
			return 0
		}

		// assign l.currentToken, for use in error-message generation
		l.currentToken = tok

		txt := string(tok.data)
		cmpTxt := strings.ToLower(txt)
		//fmt.Printf("found token: %q\n", txt)

		// Quick-check for types the scanner/tokenizer identify at a lower
		// level.
		switch tok.typ {
		case tokenTypeWhiteSpace:
			continue
		case tokenTypeSemicolon:
			return semicolon
		case tokenTypeComma:
			return comma
		case tokenTypeBlockStart:
			return openBrace
		case tokenTypeBlockEnd:
			return closeBrace
		case tokenTypeString:
			lval.str = strings.TrimPrefix(strings.TrimSuffix(txt, "\""), "\"")
			return stringConst
		}

		// Simple string lookups
		if tok, found := stringTokenMap[cmpTxt]; found {
			lval.str = txt
			return tok
		}

		// Pattern matches
		if macAddrRegexp.MatchString(cmpTxt) {
			lval.str = txt
			return macAddr
		} else if cidrRegexp.MatchString(cmpTxt) {
			lval.str = txt
			return cidr
		} else if ipAddrRegexp.MatchString(cmpTxt) {
			lval.str = txt
			return ipAddr
		} else if len(tok.data) == 1 {
			return int(tok.data[0])
		}

		// Everything else is just a "WORD"
		lval.str = txt
		return word
	}
}

func (l *lexer) nextToken() (token, error) {
	if l.scanner == nil {
		return token{}, contextError("call on uninitialized lexer")
	}

	var retToken token
	workbuf := make([]byte, 1)
	var readErr error

	for {
		var bytesRead int
		bytesRead, readErr = l.dataStream.Read(workbuf)
		if bytesRead != 1 && readErr == io.EOF {
			// If we got an error or io.EOF, we may've still received data that
			// we need to process. If we didn't receive data though, return
			// immediately.
			return token{}, readErr
		}

		// Ask the scanner whether whether the byte we received is the boundary
		// of a new token, and if so what type.
		b := workbuf[0]
		code, err := l.scanner.step(b)
		if err != nil {
			return token{}, err
		}

		// Create or append to an existing token based on what the scanner told
		// us.
		switch code {
		case codeContinue:
			// codeContinue means "add to the existing token"
			l.wipToken.data = append(l.wipToken.data, b)
		case codeIdentifierBegin:
			// We've found the start of a new identifier token. If we were
			// previously working on another token-type, go ahead and return
			// that token, but save the first byte of this new identifier in
			// l.wipToken
			if len(l.wipToken.data) != 0 {
				retToken = l.wipToken
			}
			l.wipToken = token{
				typ:  tokenTypeIdentifier,
				data: []byte{b},
			}
		case codeIdentifierEnd:
			fallthrough
		case codeWhitespace:
			if l.wipToken.typ != tokenTypeWhiteSpace {
				if len(l.wipToken.data) != 0 {
					retToken = l.wipToken
					l.wipToken = token{
						typ:  tokenTypeWhiteSpace,
						data: []byte{b},
					}
				} else {
					l.wipToken = token{
						typ:  tokenTypeWhiteSpace,
						data: []byte{b},
					}
				}
			} else {
				l.wipToken.data = append(l.wipToken.data, b)
			}
		case codeBlockBegin:
			if len(l.wipToken.data) != 0 {
				retToken = l.wipToken
			}
			l.wipToken = token{
				typ:  tokenTypeBlockStart,
				data: []byte{b},
			}
		case codeBlockEnd:
			if len(l.wipToken.data) != 0 {
				retToken = l.wipToken
			}
			l.wipToken = token{
				typ:  tokenTypeBlockEnd,
				data: []byte{b},
			}
		case codeSemicolon:
			if len(l.wipToken.data) != 0 {
				retToken = l.wipToken
			}
			l.wipToken = token{
				typ:  tokenTypeSemicolon,
				data: []byte{b},
			}
		case codeComma:
			if len(l.wipToken.data) != 0 {
				retToken = l.wipToken
			}
			l.wipToken = token{
				typ:  tokenTypeComma,
				data: []byte{b},
			}
		case codeCommentBegin:
			if len(l.wipToken.data) != 0 {
				retToken = l.wipToken
			}
			l.wipToken = token{
				typ:  tokenTypeComment,
				data: []byte{b},
			}
		case codeCommentEnd:
			retToken = l.wipToken
			l.wipToken = token{
				typ:  tokenTypeWhiteSpace,
				data: []byte{b},
			}
		case codeStringBegin:
			if len(l.wipToken.data) != 0 {
				retToken = l.wipToken
			}
			l.wipToken = token{
				typ:  tokenTypeString,
				data: []byte{b},
			}
		case codeStringEnd:
			// This code is somewhat unique, since we receive it *on* the
			// boundary of a token whose length isn't exactly 1, rather than
			// immediately following (e.g. for an identifier, we get
			// codeIdentifierEnd or codeSemicolon or similar for the first byte
			// *after* the end of the identifier). This means that the next
			// token should start with a code that is different from
			// codeContinue, and can be trusted to notice that it should not
			// be appending to a token of tokenTypeString.
			//
			// Taking advantage of this, we're not going to assign
			// retToken = l.wipToken; instead we'll let the handler for
			// the first byte of the next token do that for us.
			l.wipToken.data = append(l.wipToken.data, b)
		}

		// If we found a completed token, break out so we can return it. Any
		// work-in-progress should be saved in l.wipToken for the next
		// call.
		if len(retToken.data) != 0 {
			break
		}
	}

	if readErr != nil && readErr != io.EOF {
		return token{}, contextErrorf("dataStream.Read(): %s", readErr.Error())
	}

	// If we dropped off the end of the data stream, we won't have
	// affirmatively assigned a value to retToken, but t.currentToken should
	// contain something.
	if len(retToken.data) == 0 {
		retToken = l.wipToken
	}

	return retToken, readErr
}

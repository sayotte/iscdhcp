package iscdhcp

import (
	"regexp"
)

// scanner analyzes a stream of bytes, one byte at a time, and emits codes
// denoting the beginning/end of lexical symbols
type scanner struct {
	stateStack []int
}

func (l *scanner) init() {
	l.stateStack = []int{scanStateFindAnyBegin}
}

func (l *scanner) step(b byte) (int, error) {
	for re, spec := range lexMap[l.stateStack[0]] {
		if re.Match([]byte{b}) {
			for _, newState := range spec.newStates {
				switch newState {
				case scanSameState:
					continue
				case scanPopState:
					l.stateStack = l.stateStack[1:]
					if len(l.stateStack) == 0 {
						return codeContinue, contextError("state stack is empty")
					}
				default:
					l.stateStack = append([]int{newState}, l.stateStack...)
				}
			}
			return spec.code, nil
		}
	}

	return codeContinue, nil
}

type transitionSpec struct {
	code      int
	newStates []int
}

var lexMap = map[int]map[*regexp.Regexp]transitionSpec{
	scanStateFindAnyBegin: map[*regexp.Regexp]transitionSpec{
		regexp.MustCompile(`\{`): {
			code:      codeBlockBegin,
			newStates: []int{scanSameState},
		},
		regexp.MustCompile(`}`): {
			code:      codeBlockEnd,
			newStates: []int{scanSameState},
		},
		regexp.MustCompile("[0-9a-zA-Z!=~]"): {
			code:      codeIdentifierBegin,
			newStates: []int{scanStateFindIdentifierEnd},
		},
		regexp.MustCompile(`[\s\n]`): {
			code:      codeWhitespace,
			newStates: []int{scanSameState},
		},
		regexp.MustCompile(`"`): {
			code:      codeStringBegin,
			newStates: []int{scanStateFindStringEnd},
		},
		regexp.MustCompile(`#`): {
			code:      codeCommentBegin,
			newStates: []int{scanStateFindCommentEnd},
		},
		regexp.MustCompile(`;`): {
			code:      codeSemicolon,
			newStates: []int{scanSameState},
		},
	},
	scanStateFindIdentifierEnd: map[*regexp.Regexp]transitionSpec{
		regexp.MustCompile(`[\s]`): {
			code:      codeIdentifierEnd,
			newStates: []int{scanPopState},
		},
		regexp.MustCompile(`{`): {
			code:      codeBlockBegin,
			newStates: []int{scanPopState},
		},
		regexp.MustCompile(`"`): {
			code:      codeStringBegin,
			newStates: []int{scanPopState, scanStateFindStringEnd},
		},
		regexp.MustCompile(`;`): {
			code:      codeSemicolon,
			newStates: []int{scanPopState},
		},
		regexp.MustCompile(`,`): {
			code:      codeComma,
			newStates: []int{scanPopState},
		},
	},
	scanStateFindStringEnd: map[*regexp.Regexp]transitionSpec{
		regexp.MustCompile(`"`): {
			code:      codeStringEnd,
			newStates: []int{scanPopState},
		},
	},
	scanStateFindCommentEnd: map[*regexp.Regexp]transitionSpec{
		regexp.MustCompile(`\n`): {
			code:      codeCommentEnd,
			newStates: []int{scanPopState},
		},
	},
}

const (
	scanPopState = iota
	scanSameState
	scanStateFindAnyBegin
	scanStateFindIdentifierEnd
	scanStateFindStringEnd
	scanStateFindCommentEnd
)

const (
	codeContinue = iota
	codeWhitespace
	codeBlockBegin
	codeBlockEnd
	codeStringBegin
	codeStringEnd
	codeIdentifierBegin
	codeIdentifierEnd
	codeCommentBegin
	codeCommentEnd
	codeSemicolon
	codeComma
)

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
			switch spec.newState {
			case scanSameState:
				break
			case scanPopState:
				l.stateStack = l.stateStack[1:]
				if len(l.stateStack) == 0 {
					return codeContinue, contextError("state stack is empty")
				}
			default:
				l.stateStack = append([]int{spec.newState}, l.stateStack...)
			}
			return spec.code, nil
		}
	}

	return codeContinue, nil
}

type transitionSpec struct {
	code     int
	newState int
}

var lexMap = map[int]map[*regexp.Regexp]transitionSpec{
	scanStateFindAnyBegin: map[*regexp.Regexp]transitionSpec{
		regexp.MustCompile(`\{`): {
			code:     codeBlockBegin,
			newState: scanSameState,
		},
		regexp.MustCompile(`}`): {
			code:     codeBlockEnd,
			newState: scanSameState,
		},
		regexp.MustCompile("[0-9a-zA-Z]"): {
			code:     codeIdentifierBegin,
			newState: scanStateFindIdentifierEnd,
		},
		regexp.MustCompile(`[\s\n]`): {
			code:     codeWhitespace,
			newState: scanSameState,
		},
		regexp.MustCompile(`"`): {
			code:     codeStringBegin,
			newState: scanStateFindStringEnd,
		},
		regexp.MustCompile(`#`): {
			code:     codeCommentBegin,
			newState: scanStateFindCommentEnd,
		},
		regexp.MustCompile(`=`): {
			code:     codeAssignment,
			newState: scanSameState,
		},
		regexp.MustCompile(`;`): {
			code:     codeStatementEnd,
			newState: scanSameState,
		},
	},
	scanStateFindIdentifierEnd: map[*regexp.Regexp]transitionSpec{
		regexp.MustCompile(`[\s,]`): {
			code:     codeIdentifierEnd,
			newState: scanPopState,
		},
		regexp.MustCompile(`;`): {
			code:     codeStatementEnd,
			newState: scanPopState,
		},
	},
	scanStateFindStringEnd: map[*regexp.Regexp]transitionSpec{
		regexp.MustCompile(`"`): {
			code:     codeStringEnd,
			newState: scanPopState,
		},
	},
	scanStateFindCommentEnd: map[*regexp.Regexp]transitionSpec{
		regexp.MustCompile(`\n`): {
			code:     codeCommentEnd,
			newState: scanPopState,
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
	codeAssignment
	codeStatementEnd
)

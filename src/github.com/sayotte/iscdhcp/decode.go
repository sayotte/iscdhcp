package iscdhcp

import (
	"io"
)

// Decode analyzes a slice of bytes, constructing primitive ISC-DHCP config
// objects.
func Decode(dataStream io.Reader) ([]Statement, error) {
	l := newLexer(dataStream)

	parser := yyNewParser()
	exitCode := parser.Parse(l)
	if exitCode != 0 || l.hasErrored {
		return nil, l.err
	}

	return l.dirtyHackReturn, nil
}

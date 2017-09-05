package iscdhcp

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestLexer_nextToken(t *testing.T) {
	// start with a string that includes some funky spacing, and all the most
	// common identifiers
	data := `
		group {
			authoritative ;
			option domain-name-servers 1.2.3.4, 5.6.7.8;
			if "foo"="foo"{
				host serverA.myDomain.tld{
					hardware ethernet 0:01:a2:3b:4:50;
					fixed-address 1.2.3.4, 5.6.7.8;
				}
			}
		}`

	l := newLexer(bytes.NewReader([]byte(data)))
	var tokens []token
	for {
		tok, err := l.nextToken()
		if err != nil && err != io.EOF {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(tok.data) != 0 {
			tokens = append(tokens, tok)
		}
		if err == io.EOF {
			break
		}
	}

	expected := []token{
		{[]byte("\n\t\t"), tokenTypeWhiteSpace},
		{[]byte("group"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("{"), tokenTypeBlockStart},
		{[]byte("\n\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("authoritative"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte(";"), tokenTypeSemicolon},
		{[]byte("\n\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("option"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("domain-name-servers"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("1.2.3.4"), tokenTypeIdentifier},
		{[]byte(","), tokenTypeComma},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("5.6.7.8"), tokenTypeIdentifier},
		{[]byte(";"), tokenTypeSemicolon},
		{[]byte("\n\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("if"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("\"foo\""), tokenTypeString},
		{[]byte("="), tokenTypeIdentifier},
		{[]byte("\"foo\""), tokenTypeString},
		{[]byte("{"), tokenTypeBlockStart},
		{[]byte("\n\t\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("host"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("serverA.myDomain.tld"), tokenTypeIdentifier},
		{[]byte("{"), tokenTypeBlockStart},
		{[]byte("\n\t\t\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("hardware"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("ethernet"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("0:01:a2:3b:4:50"), tokenTypeIdentifier},
		{[]byte(";"), tokenTypeSemicolon},
		{[]byte("\n\t\t\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("fixed-address"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("1.2.3.4"), tokenTypeIdentifier},
		{[]byte(","), tokenTypeComma},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("5.6.7.8"), tokenTypeIdentifier},
		{[]byte(";"), tokenTypeSemicolon},
		{[]byte("\n\t\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("}"), tokenTypeBlockEnd},
		{[]byte("\n\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("}"), tokenTypeBlockEnd},
	}

	lesserLen := len(expected)
	if len(tokens) < lesserLen {
		lesserLen = len(tokens)
	}
	for i := 0; i < lesserLen; i++ {
		if !reflect.DeepEqual(tokens[i], expected[i]) {
			t.Errorf("index %d: unequal, expected %q (type %d) got %q (type %d)", i, string(expected[i].data), expected[i].typ, string(tokens[i].data), tokens[i].typ)
		}
	}
}

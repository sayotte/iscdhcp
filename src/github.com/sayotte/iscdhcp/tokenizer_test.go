package iscdhcp

import (
	"reflect"
	"testing"
)

func TestTokenizer_Tokenize(t *testing.T) {
	data := `
		group {
			authoritative;
			option domain-name-servers 1.2.3.4, 5.6.7.8;
			host serverA.myDomain.tld {
				hardware ethernet 0:1:2:3:4:5;
				fixed-address 1.2.3.4, 5.6.7.8;
			}
		}`

	var toker tokenizer
	tokens, err := toker.Tokenize([]byte(data))
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expected := []token{
		{[]byte("\n\t\t"), tokenTypeWhiteSpace},
		{[]byte("group"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("{"), tokenTypeBlockStart},
		{[]byte("\n\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("authoritative"), tokenTypeIdentifier},
		{[]byte(";"), tokenTypeStatementEnd},
		{[]byte("\n\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("option"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("domain-name-servers"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("1.2.3.4"), tokenTypeIdentifier},
		{[]byte(", "), tokenTypeWhiteSpace},
		{[]byte("5.6.7.8"), tokenTypeIdentifier},
		{[]byte(";"), tokenTypeStatementEnd},
		{[]byte("\n\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("host"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("serverA.myDomain.tld"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("{"), tokenTypeBlockStart},
		{[]byte("\n\t\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("hardware"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("ethernet"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("0:1:2:3:4:5"), tokenTypeIdentifier},
		{[]byte(";"), tokenTypeStatementEnd},
		{[]byte("\n\t\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("fixed-address"), tokenTypeIdentifier},
		{[]byte(" "), tokenTypeWhiteSpace},
		{[]byte("1.2.3.4"), tokenTypeIdentifier},
		{[]byte(", "), tokenTypeWhiteSpace},
		{[]byte("5.6.7.8"), tokenTypeIdentifier},
		{[]byte(";"), tokenTypeStatementEnd},
		{[]byte("\n\t\t\t"), tokenTypeWhiteSpace},
		{[]byte("}"), tokenTypeBlockEnd},
		{[]byte("\n\t\t"), tokenTypeWhiteSpace},
		{[]byte("}"), tokenTypeBlockEnd},
	}

	if len(expected) != len(tokens) {
		t.Fatal("different lengths")
	}

	for i := range expected {
		if !reflect.DeepEqual(tokens[i], expected[i]) {
			t.Errorf("index %d: unequal, expected %q (type %d) got %q (type %d)", i, string(expected[i].data), expected[i].typ, string(tokens[i].data), tokens[i].typ)
		}
	}
}

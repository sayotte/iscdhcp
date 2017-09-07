package iscdhcp

import (
	"net"
	"reflect"
	"testing"
)

func TestComplex_roundTrip(t *testing.T) {
	hs := &hostStatement{}
	hs.hostname = "serverA.myDomain.tld"
	hs.params = append(hs.params, &hardwareStatement{"ethernet", "0:1:2:3:4:5"})
	hs.params = append(hs.params, &fixedAddressStatement{net.ParseIP("1.2.3.4"), net.ParseIP("5.6.7.8")})
	hs.params = append(hs.params, &includeStatement{"filename.cfg"})

	subnetStmt := &subnetStatement{
		subnetNumber: net.ParseIP("1.2.3.0"),
		netmask:      net.ParseIP("255.255.255.0"),
	}
	subnetStmt.params = append(subnetStmt.params, hs)

	gs := &groupStatement{}
	trueVal := true
	as := (*authoritativeStatement)(&trueVal)
	uhdns := (*useHostDeclNamesStatement)(&trueVal)
	gs.params = append(gs.params, as)
	gs.params = append(gs.params, uhdns)
	gs.params = append(gs.params, &domainNameServersOption{net.ParseIP("1.2.3.4"), net.ParseIP("5.6.7.8")})
	gs.params = append(gs.params, subnetStmt)

	toker := &tokenizer{}
	tokens, err := toker.Tokenize([]byte(gs.indentedString("")))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	tStream := &tokenStream{tokens: tokens}
	d := Decoder{}
	newStatements, err := d.Decode(tStream)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(gs, newStatements[0]) {
		t.Error("expected != actual")
	}
}

func TestStatements_roundtrip(t *testing.T) {
	trueReal := true
	stringReal := "foo"
	intReal := 0
	statements := []statement{
		&hostStatement{hostname: "serverA.myDomain.tld"},
		(*adaptiveLeaseThresholdStatement)(&intReal),
		(*alwaysBroadcastStatement)(&trueReal),
		(*alwaysReplyRFC1048)(&trueReal),
		(*authoritativeStatement)(&trueReal),
		(*bootUnknownClientsStatement)(&trueReal),
		(*ddnsDomainNameStatement)(&stringReal),
		(*ddnsHostNameStatement)(&stringReal),
		(*ddnsRevDomainNameStatement)(&stringReal),
		(*ddnsUpdateStyleStatement)(&intReal),
		(*ddnsUpdatesStatement)(&trueReal),
		(*defaultLeaseTimeStatement)(&intReal),
		(*delayedAckStatement)(&intReal),
		(*doForwardUpdatesStatement)(&trueReal),
		&dynamicBootpLeaseCutoffStatement{
			dayOfWeek:  1,
			year:       1999,
			month:      12,
			dayOfMonth: 31,
			hours:      23,
			minutes:    59,
			seconds:    59,
		},
		(*maxAckDelayStatement)(&intReal),
		(*useHostDeclNamesStatement)(&trueReal),
	}

	for _, statement := range statements {
		toker := &tokenizer{}
		tokens, err := toker.Tokenize([]byte(statement.indentedString("")))
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		tStream := &tokenStream{tokens: tokens}
		d := Decoder{}
		newStatements, err := d.Decode(tStream)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if len(newStatements) != 1 {
			t.Fatalf("expected exactly 1 statement, got %d", len(newStatements))
		}
		if !reflect.DeepEqual(newStatements[0], statement) {
			t.Error("actual != expected")
			t.Logf("actual: %s", newStatements[0].indentedString(""))
			t.Logf("expected: %s", statement.indentedString(""))
		}
	}
}

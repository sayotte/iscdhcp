package iscdhcp

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"testing"
)

func TestComplex_roundTrip(t *testing.T) {
	hs := HostStatement{}
	hs.Hostname = "serverA.myDomain.tld"
	hs.Statements = append(hs.Statements, HardwareStatement{"ethernet", "0:1:2:3:4:5"})
	hs.Statements = append(hs.Statements, FixedAddressStatement{net.ParseIP("1.2.3.4")})
	hs.Statements = append(hs.Statements, IncludeStatement{"filename.cfg"})

	subnetStmt := SubnetStatement{
		SubnetNumber: net.ParseIP("1.2.3.0"),
		Netmask:      net.ParseIP("255.255.255.0"),
	}
	subnetStmt.Statements = append(subnetStmt.Statements, hs)

	gs := GroupStatement{}
	//trueVal := true
	//as := (*authoritativeStatement)(&trueVal)
	//uhdns := (*useHostDeclNamesStatement)(&trueVal)
	//gs.params = append(gs.params, as)
	//gs.params = append(gs.params, uhdns)
	//gs.params = append(gs.params, &domainNameServersOption{net.ParseIP("1.2.3.4"), net.ParseIP("5.6.7.8")})
	gs.Statements = append(gs.Statements, subnetStmt)

	newStatements, err := Decode(strings.NewReader(gs.IndentedString("")))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(gs, newStatements[0]) {
		t.Error("expected != actual")
	}
}

func TestStatements_roundtrip(t *testing.T) {
	//trueReal := true
	//stringReal := "foo"
	//intReal := 0
	ip1 := net.ParseIP("1.2.3.2")
	ip2 := net.ParseIP("4.5.6.2")
	statements := []Statement{
		AuthoritativeStatement(false),
		AuthoritativeStatement(true),
		HostStatement{Hostname: "serverA.myDomain.tld"},
		//		(*adaptiveLeaseThresholdStatement)(&intReal),
		//		(*alwaysBroadcastStatement)(&trueReal),
		//		(*alwaysReplyRFC1048)(&trueReal),
		//		(*authoritativeStatement)(&trueReal),
		//		(*bootUnknownClientsStatement)(&trueReal),
		//		(*ddnsDomainNameStatement)(&stringReal),
		//		(*ddnsHostNameStatement)(&stringReal),
		//		(*ddnsRevDomainNameStatement)(&stringReal),
		//		(*ddnsUpdateStyleStatement)(&intReal),
		//		(*ddnsUpdatesStatement)(&trueReal),
		//		(*defaultLeaseTimeStatement)(&intReal),
		//		(*delayedAckStatement)(&intReal),
		//		(*doForwardUpdatesStatement)(&trueReal),
		//		&dynamicBootpLeaseCutoffStatement{
		//			dayOfWeek:  1,
		//			year:       1999,
		//			month:      12,
		//			dayOfMonth: 31,
		//			hours:      23,
		//			minutes:    59,
		//			seconds:    59,
		//		},
		FixedAddressStatement{ip1, ip2},
		HardwareStatement{HardwareType: "ethernet", HardwareAddress: "1:2:3:4:5:6"},
		IncludeStatement{"filename"},
		//		(*maxAckDelayStatement)(&intReal),
		UseHostDeclNamesStatement(true),
		DomainNameServersOption{ip1, ip2},
	}

	for _, statement := range statements {

		newStatements, err := Decode(strings.NewReader(statement.IndentedString("")))
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if len(newStatements) != 1 {
			t.Fatalf("expected exactly 1 statement, got %d", len(newStatements))
		}
		if !reflect.DeepEqual(newStatements[0], statement) {
			t.Error("actual != expected")
		}
	}
}

func TestConditionalStatement_roundtrip(t *testing.T) {
	// Create this:
	//   if ("foo" != "foo") or ((option domain = "foo") and not static) { }
	//   elsif known { }
	//   elsif "foo" ~~ "FOO" { }
	//   else { }
	// This touches on all control-flow modifiers, and all but one boolean
	// operator. It also ensures that operator associativity works as expected.
	//
	// We will round-trip by generating a string from this, then parse that
	// string back into in-memory structs, and then compare those structs to
	// what we started with expecting them to be equal.
	expected := ConditionalStatement{
		Operator: ConditionIf,
		Condition: BooleanExpression{
			Operator: BoolOr,
			BoolTerms: []BooleanExpression{
				{
					Operator: BoolInequal,
					DataTerms: []fmt.Stringer{
						StringConstTerm("foo"),
						StringConstTerm("foo"),
					},
				},
				{
					Operator: BoolAnd,
					BoolTerms: []BooleanExpression{
						{
							Operator: BoolEqual,
							DataTerms: []fmt.Stringer{
								PacketOptionTerm{"domain"},
								StringConstTerm("foo"),
							},
						},
						{
							Operator: BoolNot,
							BoolTerms: []BooleanExpression{
								{
									Operator: BoolStatic,
								},
							},
						},
					},
				},
			},
		},
		SubConditionals: []ConditionalStatement{
			{
				Operator:  ConditionElsif,
				Condition: BooleanExpression{Operator: BoolKnown},
			},
			{
				Operator: ConditionElsif,
				Condition: BooleanExpression{
					Operator: BoolRegexIMatch,
					DataTerms: []fmt.Stringer{
						StringConstTerm("foo"),
						StringConstTerm("FOO"),
					},
				},
			},
			{
				Operator: ConditionElse,
			},
		},
	}

	// First verify that we generate a sane string, not just one that survives
	// the round-trip further down
	expectedSubstring :=
		`if "foo" != "foo" or option domain = "foo" and not static {
}
elsif known {
}
elsif "foo" ~~ "FOO" {
}
else {
}`
	actualString := expected.IndentedString("")
	if !strings.Contains(actualString, expectedSubstring) {
		t.Errorf("expected string containing %q, got %q", expectedSubstring, actualString)
	}

	// Now round-trip the statements to a string and back, and verify we're
	// left with the same thing
	newStatements, err := Decode(strings.NewReader(actualString))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(newStatements) != 1 {
		t.Fatalf("expected exactly 1 statement, got %d", len(newStatements))
	}

	if !reflect.DeepEqual(expected, newStatements[0]) {
		t.Error("expected != actual")
	}
}

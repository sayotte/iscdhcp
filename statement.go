package iscdhcp

import (
	"fmt"
	"net"
)

const defaultIndent = "  "

// A Statement represents an ISC-DHCP configuration statement.
// See dhcpd.conf(5) dhcp-options(5) and dhcpd-eval(5) for canonical types.
// For implementations, see "go doc iscdhcp | grep Statement".
type Statement interface {
	// IndentedString produces string representation of the config statement,
	// in the form expected by dhcpd. The "prefix" argument is used to indent
	// nested statements; it is not a standard indent-depth, but an explicit
	// prefix for this particular statement's string representation.
	// FIXME It'd really be nice to be able to configure the standard depth.
	IndentedString(prefix string) string
}

// base behaviors

type block []Statement

func (b block) IndentedString(prefix string) string {
	var s string
	for _, decl := range b {
		s += decl.IndentedString(prefix + defaultIndent)
	}
	return s
}

type intDecl int

func (id intDecl) IndentedString(prefix, identifier string) string {
	return prefix + fmt.Sprintf("%s %d;\n", identifier, id)
}

type onOffBool bool

func (oob onOffBool) IndentedString(prefix, identifier string) string {
	onOrOff := "off"
	if oob {
		onOrOff = "on"
	}

	return prefix + identifier + " " + onOrOff + ";\n"
}

type stringDecl string

func (sd stringDecl) IndentedString(prefix, identifier string) string {
	return prefix + identifier + " \"" + string(sd) + "\";\n"
}

// DECLARATIONS

// A GroupStatement represents a group declaration.
// See "The group statement" in dhcpd.conf(5)
type GroupStatement struct {
	Statements []Statement
}

// IndentedString implements the method of the same name in the Statement interface
func (gs GroupStatement) IndentedString(prefix string) string {
	return prefix + "group {\n" + block(gs.Statements).IndentedString(prefix+defaultIndent) +
		prefix + "}\n"
}

// A HostStatement represents a host declaration.
// See "The host statement" in dhcpd.conf(5)
type HostStatement struct {
	Hostname   string
	Statements []Statement
}

// IndentedString implements the method of the same name in the Statement interface
func (hs HostStatement) IndentedString(prefix string) string {
	return prefix + "host " + hs.Hostname + " {\n" +
		block(hs.Statements).IndentedString(prefix+defaultIndent) +
		prefix + "}\n"
}

// An IncludeStatement represents an include-file declaration.
// See "The include statement" in dhcpd.conf(5)
type IncludeStatement struct {
	Filename string
}

// IndentedString implements the method of the same name in the Statement interface
func (is IncludeStatement) IndentedString(prefix string) string {
	return prefix + "include \"" + is.Filename + "\";\n"
}

type sharedNetworkStatement struct {
	Name       string
	Statements []Statement
}

func (sns sharedNetworkStatement) IndentedString(prefix string) string {
	return prefix + "shared-network " + sns.Name + " {\n" +
		block(sns.Statements).IndentedString(prefix+defaultIndent) +
		prefix + "}\n"
}

// A SubnetStatement represents a subnet declaration.
// See "The subnet statement" in dhcpd.conf(5)
type SubnetStatement struct {
	SubnetNumber net.IP
	Netmask      net.IP
	Statements   []Statement
}

// IndentedString implements the method of the same name in the Statement interface
func (sns SubnetStatement) IndentedString(prefix string) string {
	return prefix + "subnet " + sns.SubnetNumber.String() +
		" netmask " + sns.Netmask.String() + " {\n" +
		block(sns.Statements).IndentedString(prefix+defaultIndent) +
		prefix + "}\n"
}

// PARAMETERS

type adaptiveLeaseThresholdStatement int

func (alts adaptiveLeaseThresholdStatement) IndentedString(prefix string) string {
	return intDecl(alts).IndentedString(prefix, "adaptive-lease-time-threshold")
}

type alwaysBroadcastStatement bool

func (abs alwaysBroadcastStatement) IndentedString(prefix string) string {
	return onOffBool(abs).IndentedString(prefix, "always-broadcast")
}

type alwaysReplyRFC1048 bool

func (arrfc alwaysReplyRFC1048) IndentedString(prefix string) string {
	return onOffBool(arrfc).IndentedString(prefix, "always-reply-rfc1048")
}

// An AuthoritativeStatement represents an "authoritative" parameter.
// See "The authoritative statement" in dhcpd.conf(5)
type AuthoritativeStatement bool

// IndentedString implements the method of the same name in the Statement interface
func (as AuthoritativeStatement) IndentedString(prefix string) string {
	if as {
		return prefix + "authoritative;\n"
	}
	return prefix + "not authoritative;\n"
}

type bootUnknownClientsStatement bool

func (bucs bootUnknownClientsStatement) IndentedString(prefix string) string {
	return onOffBool(bucs).IndentedString(prefix, "boot-unknown-clients")
}

type dbTimeFormatStatement struct {
	localTime bool
}

func (dtfs dbTimeFormatStatement) IndentedString(prefix string) string {
	var format string
	if dtfs.localTime {
		format = "local"
	} else {
		format = "default"
	}
	return prefix + "db-time-format " + format + ";\n"
}

type ddnsDomainNameStatement string

func (ddnsds ddnsDomainNameStatement) IndentedString(prefix string) string {
	return stringDecl(ddnsds).IndentedString(prefix, "ddns-domainname")
}

type ddnsHostNameStatement string

func (ddnshns ddnsHostNameStatement) IndentedString(prefix string) string {
	return stringDecl(ddnshns).IndentedString(prefix, "ddns-hostname")
}

type ddnsRevDomainNameStatement string

func (ddnsrdns ddnsRevDomainNameStatement) IndentedString(prefix string) string {
	return stringDecl(ddnsrdns).IndentedString(prefix, "ddns-rev-domainname")
}

const (
	ddnsUpdateStyleNone = iota
	ddnsUpdateStyleAdHoc
	ddnsUpdateStyleInterim
)

type ddnsUpdateStyleStatement int

func (ddnsuss ddnsUpdateStyleStatement) IndentedString(prefix string) string {
	var style string
	switch ddnsuss {
	case ddnsUpdateStyleAdHoc:
		style = "ad-hoc"
	case ddnsUpdateStyleInterim:
		style = "interim"
	default:
		style = "none"
	}
	return prefix + "ddns-update-style " + style + ";\n"
}

type ddnsUpdatesStatement bool

func (ddnsus ddnsUpdatesStatement) IndentedString(prefix string) string {
	return onOffBool(ddnsus).IndentedString(prefix, "ddns-updates")
}

type defaultLeaseTimeStatement int

func (dlts defaultLeaseTimeStatement) IndentedString(prefix string) string {
	return intDecl(dlts).IndentedString(prefix, "default-lease-time")
}

type delayedAckStatement int

func (das delayedAckStatement) IndentedString(prefix string) string {
	return intDecl(das).IndentedString(prefix, "delayed-ack")
}

type doForwardUpdatesStatement bool

func (dfus doForwardUpdatesStatement) IndentedString(prefix string) string {
	return onOffBool(dfus).IndentedString(prefix, "do-forward-updates")
}

type dynamicBootpLeaseCutoffStatement struct {
	dayOfWeek  int
	year       int
	month      int
	dayOfMonth int
	hours      int
	minutes    int
	seconds    int
}

func (dblcs dynamicBootpLeaseCutoffStatement) IndentedString(prefix string) string {
	// W YYYY/MM/DD HH:MM:SS
	return prefix + fmt.Sprintf(
		"dynamic-bootp-lease-cutoff %1d %04d/%02d/%02d %02d:%02d:%02d;\n",
		dblcs.dayOfWeek,
		dblcs.year,
		dblcs.month,
		dblcs.dayOfMonth,
		dblcs.hours,
		dblcs.minutes,
		dblcs.seconds,
	)
}

// A FixedAddressStatement represents a fixed-address parameter.
// See "The fixed-address declaration" in dhcpd.conf(5)
type FixedAddressStatement []net.IP

// IndentedString implements the method of the same name in the Statement interface
func (fas FixedAddressStatement) IndentedString(prefix string) string {
	s := prefix + "fixed-address "
	for i := 0; i < len(fas)-1; i++ {
		s += fas[i].String() + ", "
	}
	return s + fas[len(fas)-1].String() + ";\n"
}

// A HardwareStatement represents a hardware parameter.
// See "The hardware statement" in dhcpd.conf(5)
type HardwareStatement struct {
	HardwareType    string
	HardwareAddress string
}

// IndentedString implements the method of the same name in the Statement interface
func (hs HardwareStatement) IndentedString(prefix string) string {
	return prefix + "hardware " + hs.HardwareType + " " + hs.HardwareAddress + ";\n"
}

type maxAckDelayStatement int

func (mads maxAckDelayStatement) IndentedString(prefix string) string {
	return intDecl(mads).IndentedString(prefix, "max-ack-delay")
}

// A UseHostDeclNamesStatement represents a "use-host-decl-names" parameter.
// See "The use-host-decl-names statement" in dhcpd.conf(5)
type UseHostDeclNamesStatement bool

// IndentedString implements the method of the same name in the Statement interface
func (uhdns UseHostDeclNamesStatement) IndentedString(prefix string) string {
	return onOffBool(uhdns).IndentedString(prefix, "use-host-decl-names")
}

// OPTIONS

// A DomainNameServersOption represents a domain-name-servers option parameter.
// See "option domain-name-servers" in dhcp-options(5)
type DomainNameServersOption []net.IP

// IndentedString implements the method of the same name in the Statement interface
func (dnso DomainNameServersOption) IndentedString(prefix string) string {
	s := prefix + "option domain-name-servers "
	for i := 0; i < len(dnso)-1; i++ {
		s += dnso[i].String() + ", "
	}
	return s + dnso[len(dnso)-1].String() + ";\n"
}

// CONDITIONALS

var conditionOpStrings = map[int]string{
	ConditionIf:    "if",
	ConditionElsif: "elsif",
	ConditionElse:  "else",
}

// A ConditionalStatement represents a conditional-evaluation statement in a
// config file, a la dhcp-eval(5).
//
// Each top-level ConditionalStatement must be an if-statement; it may
// optionally have else-if and else sub ConditionalStatements.
//
// Each ConditionalStatement also contains a block of other Statements.
//
// Example usage:
//	cs := ConditionalStatement{
//		Operator: ConditionIf,
//		Condition: BooleanExpression {
//			Operator: BoolEqual,
//			DataTerms: []fmt.Stringer {
//				PacketOptionTerm{"user-class"},
//				StringConstTerm("iPXE"),
//			},
//		},
//		SubConditionals: []ConditionalStatement {
//			{
//				Operator: ConditionElse,
//			},
//		},
//	}
//	cs.Statements = []Statement{IncludeStatement{"ipxe.conf"}}
//	cs.SubConditionals[0].Statements = []Statement{IncludeStatement{"non-ipxe.conf"}}
//
// This example corresponds to the following config-file text:
//
//	if option user-class = "iPXE" {
//		include "ipxe.conf";
//	}
//	else {
//		include "non-ipxe.conf";
//	}
type ConditionalStatement struct {
	// Operator is one of if/elsif/else, represented by the integers
	// {ConditionIf, ConditionElsif, ConditionElse}.
	Operator int
	// Condition is the expression immediately following the "if/elsif" word,
	// which determines if the other Statements within this ConditionalStatement
	// are evaluated. It is unnecessary / will be ignored for "else"
	// sub-conditionals.
	Condition BooleanExpression
	// Statements is the list of Statements which will be evaluated if Condition
	// proves True.
	Statements []Statement
	// SubConditionals is a list of elsif/else ConditionalStatements, whose
	// evaluation is predicated on the top-level statement not being evaluated,
	// and their own Conditions proving True.
	SubConditionals []ConditionalStatement
}

// IndentedString implements the method of the same name in the Statement interface
func (cs ConditionalStatement) IndentedString(prefix string) string {
	var s string
	if cs.Operator == ConditionIf || cs.Operator == ConditionElsif {
		s = fmt.Sprintf(
			"%s%s %s {\n%s}\n",
			prefix,
			conditionOpStrings[cs.Operator],
			cs.Condition.string(),
			block(cs.Statements).IndentedString(prefix+defaultIndent))
	} else {
		s = fmt.Sprintf(
			"%s%s {\n%s}\n",
			prefix,
			conditionOpStrings[cs.Operator],
			block(cs.Statements).IndentedString(prefix+defaultIndent))
	}
	for _, sc := range cs.SubConditionals {
		s += sc.IndentedString(prefix)
	}

	return s
}

var boolOpStrings = map[int]string{
	BoolAnd:         "and",
	BoolOr:          "or",
	BoolNot:         "not",
	BoolStatic:      "static",
	BoolKnown:       "known",
	BoolExists:      "exists",
	BoolEqual:       "=",
	BoolInequal:     "!=",
	BoolRegexMatch:  "~=",
	BoolRegexIMatch: "~~",
}

// A BooleanExpression is an expression of terms and operators which can be
// evaluated into a single true or false value. BooleanExpressions are used
// to determine if the statements within a ConditionalStatement will be
// evaluated. See dhcp-eval(5) for more information.
type BooleanExpression struct {
	// Operator is one of the boolean operators specified in dhcp-eval(5),
	// represented by the integer constants in this package beginning with
	// "Bool", e.g. BoolAnd, BoolOr, BoolEqual etc.
	//
	// The Operator specified will apply to either the BoolTerms or the
	// DataTerms supplied with the BooleanExpression.
	Operator int // one of boolOpStrings(...)
	// A BooleanExpression may apply an operator to the output of sub-
	// -expressions, e.g. the "and" operator in the algebraic expression
	// (x > 1) and (x < 10). In this case, the sub-expressions should be
	// specified under the BoolTerms field.
	BoolTerms []BooleanExpression
	// A BooleanExpression will otherwise apply to comparisons or operations
	// on data values, e.g. "if x = 'foo'". In this case, the sub-expressions
	// should be specified under the DataTerms field.
	DataTerms []fmt.Stringer
}

func (be BooleanExpression) string() string {
	switch be.Operator {
	case BoolAnd:
		fallthrough
	case BoolOr:
		if len(be.BoolTerms) != 2 {
			return "ERROR_INCORRECT_NUMBER_OF_TERMS (a)"
		}
		return fmt.Sprintf("%s %s %s", be.BoolTerms[0].string(), boolOpStrings[be.Operator], be.BoolTerms[1].string())
	case BoolNot:
		if len(be.BoolTerms) != 1 {
			return "ERROR_INCORRECT_NUMBER_OF_TERMS (b)"
		}
		return fmt.Sprintf("%s %s", boolOpStrings[be.Operator], be.BoolTerms[0].string())
	case BoolStatic:
		fallthrough
	case BoolKnown:
		return boolOpStrings[be.Operator]
	case BoolExists:
		if len(be.DataTerms) != 1 {
			return "ERROR_INCORRECT_NUMBER_OF_TERMS (c)"
		}
		return fmt.Sprintf("%s %s", boolOpStrings[be.Operator], be.DataTerms[0].String())
	case BoolEqual:
		fallthrough
	case BoolInequal:
		fallthrough
	case BoolRegexMatch:
		fallthrough
	case BoolRegexIMatch:
		if len(be.DataTerms) != 2 {
			return "ERROR_INCORRECT_NUMBER_OF_TERMS (d)"
		}
		return fmt.Sprintf("%s %s %s", be.DataTerms[0].String(), boolOpStrings[be.Operator], be.DataTerms[1].String())
	}

	return "ERROR_INVALID_OPERATOR"
}

// A StringConstTerm is a data-term used in a BooleanExpression. It represents
// a quote-enclosed arbitrary string, e.g. ``"foo"''.
type StringConstTerm string

func (sct StringConstTerm) String() string {
	return "\"" + string(sct) + "\""
}

// A PacketOptionTerm is a data-term used in a BooleanExpression. It represents
// an option in a DHCP packet, and is stringified for config file syntax like
// ``option user-class''.
type PacketOptionTerm struct {
	optionName string
}

func (pot PacketOptionTerm) String() string {
	return "option " + pot.optionName
}

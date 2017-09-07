package iscdhcp

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

const defaultIndent = "  "

var ipAddressRegexp = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)

var hwAddressRegexp = regexp.MustCompile(`([0-9a-fA-f]{1,2}:)+[0-9a-fA-F]{1,2}$`)

type statement interface {
	indentedString(prefix string) string
	fromTokenStream(tStream *tokenStream) error
}

// base behaviors

type baseScope struct {
	params       []statement
	declarations []statement
}

func (bs baseScope) indentedString(prefix string) string {
	var s string
	for _, param := range bs.params {
		s += param.indentedString(prefix + defaultIndent)
	}
	for _, decl := range bs.declarations {
		s += decl.indentedString(prefix + defaultIndent)
	}
	return s
}

func (bs *baseScope) fromTokenStream(tStream *tokenStream) error {
	// Find / ensure we have a block-start
	for {
		tok, err := tStream.next()
		if err != nil {
			return contextErrorf("error searching for start-of-block '{': %s", err)
		}
		if tok.typ == tokenTypeBlockStart {
			break
		}
	}

	// Get all the tokens from inside the block
	var blockTokens []token
	var subBlockDepth int
OUTER:
	for {
		tok, err := tStream.next()
		if err != nil {
			return err
		}
		switch tok.typ {
		case tokenTypeBlockStart:
			blockTokens = append(blockTokens, tok)
			subBlockDepth++
		case tokenTypeBlockEnd:
			if subBlockDepth == 0 {
				break OUTER
			}
			blockTokens = append(blockTokens, tok)
			subBlockDepth--
		default:
			blockTokens = append(blockTokens, tok)
		}
	}

	// Now decode all the tokens inside the block
	d := Decoder{}
	var err error
	bs.params, err = d.Decode(&tokenStream{tokens: blockTokens})
	return err
}

type intDecl int

func (id intDecl) indentedString(prefix, identifier string) string {
	return prefix + fmt.Sprintf("%s %d;\n", identifier, id)
}

func (id *intDecl) fromTokenStream(tStream *tokenStream, identifier string) error {
	err := expectIdentifier(tStream, identifier)
	if err != nil {
		return err
	}

	err = expectWhitespace(tStream)
	if err != nil {
		return err
	}

	tok, err := tStream.next()
	if err != nil {
		return err
	}
	digit, err := strconv.Atoi(string(tok.data))
	if err != nil {
		return contextErrorf("strconv.Atoi(%q): %s", string(tok.data), err)
	}
	*id = intDecl(digit)

	return nil
}

type onOffBool bool

func (oob onOffBool) indentedString(prefix, identifier string) string {
	onOrOff := "off"
	if oob {
		onOrOff = "on"
	}

	return prefix + identifier + " " + onOrOff + ";\n"
}

func (oob *onOffBool) fromTokenStream(tStream *tokenStream, identifier string) error {
	err := expectIdentifier(tStream, identifier)
	if err != nil {
		return err
	}

	err = expectWhitespace(tStream)
	if err != nil {
		return err
	}

	tok, err := tStream.next()
	if err != nil {
		return err
	}
	tokVal := strings.ToLower(string(tok.data))
	if tok.typ == tokenTypeIdentifier &&
		(tokVal == "on" || tokVal == "off" || tokVal == "true" || tokVal == "false") {
		*oob = false
		if tokVal == "on" || tokVal == "false" {
			*oob = true
		}
	} else {
		return contextErrorf("expected 'on' or 'off' identifier, got %q", string(tok.data))
	}

	return nil
}

type stringDecl string

func (sd stringDecl) indentedString(prefix, identifier string) string {
	return prefix + identifier + " \"" + string(sd) + "\";\n"
}

func (sd *stringDecl) fromTokenStream(tStream *tokenStream, identifier string) error {
	err := expectIdentifier(tStream, identifier)
	if err != nil {
		return err
	}

	err = expectWhitespace(tStream)
	if err != nil {
		return err
	}

	tok, err := tStream.next()
	if err != nil {
		return err
	}
	if tok.typ != tokenTypeString {
		return contextErrorf("expected string constant, got %q", string(tok.data))
	}
	*sd = stringDecl(strings.TrimPrefix(strings.TrimSuffix(string(tok.data), "\""), "\""))

	return nil
}

// DECLARATIONS

type groupStatement struct {
	baseScope
}

func (gs groupStatement) indentedString(prefix string) string {
	return prefix + "group {\n" + gs.baseScope.indentedString(prefix+defaultIndent) +
		prefix + "}\n"
}

func (gs *groupStatement) fromTokenStream(tStream *tokenStream) error {
	// Expect the term "group"
	err := expectIdentifier(tStream, "group")
	if err != nil {
		return contextError(err.Error())
	}

	return gs.baseScope.fromTokenStream(tStream)
}

type hostStatement struct {
	hostname string
	baseScope
}

func (hs hostStatement) indentedString(prefix string) string {
	return prefix + "host " + hs.hostname + " {\n" +
		hs.baseScope.indentedString(prefix+defaultIndent) +
		prefix + "}\n"
}

func (hs *hostStatement) fromTokenStream(tStream *tokenStream) error {
	// Expect the term "host"
	err := expectIdentifier(tStream, "host")
	if err != nil {
		return contextError(err.Error())
	}

	// Expect whitespace
	err = expectWhitespace(tStream)
	if err != nil {
		return contextError(err.Error())
	}

	// Find the hostname identifier
	tok, err := tStream.next()
	if err != nil {
		return err
	}
	if tok.typ != tokenTypeIdentifier {
		return contextErrorf("expected host name, got %q", string(tok.data))
	}
	hs.hostname = string(tok.data)

	return hs.baseScope.fromTokenStream(tStream)
}

type includeStatement struct {
	filename string
}

func (is includeStatement) indentedString(prefix string) string {
	return prefix + "include \"" + is.filename + "\";\n"
}

func (is *includeStatement) fromTokenStream(tStream *tokenStream) error {
	err := expectIdentifier(tStream, "include")
	if err != nil {
		return err
	}

	err = expectWhitespace(tStream)
	if err != nil {
		return err
	}

	tok, err := tStream.next()
	if err != nil {
		return err
	}
	if tok.typ != tokenTypeString {
		return contextErrorf("expected filename string, got %q", string(tok.data))
	}
	is.filename = strings.TrimPrefix(strings.TrimSuffix(string(tok.data), "\""), "\"")

	return nil
}

type sharedNetworkStatement struct {
	baseScope
	name string
}

func (sns sharedNetworkStatement) indentedString(prefix string) string {
	return prefix + "shared-network " + sns.name + " {\n" +
		sns.baseScope.indentedString(prefix+defaultIndent) +
		prefix + "}\n"
}

func (sns *sharedNetworkStatement) fromTokenStream(tStream *tokenStream) error {
	// Expect the term "shared-network"
	err := expectIdentifier(tStream, "shared-network")
	if err != nil {
		return contextError(err.Error())
	}

	// Expect whitespace
	err = expectWhitespace(tStream)
	if err != nil {
		return contextError(err.Error())
	}

	// Find the network name identifier
	tok, err := tStream.next()
	if err != nil {
		return err
	}
	if tok.typ != tokenTypeIdentifier {
		return contextErrorf("expected shared-network name, got %q", string(tok.data))
	}
	sns.name = string(tok.data)

	return sns.baseScope.fromTokenStream(tStream)
}

type subnetStatement struct {
	subnetNumber net.IP
	netmask      net.IP
	baseScope
}

func (sns subnetStatement) indentedString(prefix string) string {
	return prefix + "subnet " + sns.subnetNumber.String() +
		" netmask " + sns.netmask.String() + " {\n" +
		sns.baseScope.indentedString(prefix+defaultIndent) +
		prefix + "}\n"
}

func (sns *subnetStatement) fromTokenStream(tStream *tokenStream) error {
	// Expect the word "subnet"
	err := expectIdentifier(tStream, "subnet")
	if err != nil {
		return contextError(err.Error())
	}

	// Expect whitespace
	err = expectWhitespace(tStream)
	if err != nil {
		return contextError(err.Error())
	}

	// Expect the subnet number IP
	tok, err := tStream.next()
	if err != nil {
		return err
	}
	if !ipAddressRegexp.Match(tok.data) {
		return contextErrorf("malformed IP address: %q", string(tok.data))
	}
	sns.subnetNumber = net.ParseIP(string(tok.data))

	// Expect whitespace
	err = expectWhitespace(tStream)
	if err != nil {
		return contextError(err.Error())
	}

	// Expect the "netmask" identifier
	err = expectIdentifier(tStream, "netmask")
	if err != nil {
		return contextError(err.Error())
	}

	// Expect whitespace
	err = expectWhitespace(tStream)
	if err != nil {
		return contextError(err.Error())
	}

	// Expect the netmask address
	tok, err = tStream.next()
	if err != nil {
		return err
	}
	if !ipAddressRegexp.Match(tok.data) {
		return contextErrorf("malformed IP address: %q", string(tok.data))
	}
	sns.netmask = net.ParseIP(string(tok.data))

	return sns.baseScope.fromTokenStream(tStream)
}

// PARAMETERS

type adaptiveLeaseThresholdStatement int

func (alts adaptiveLeaseThresholdStatement) indentedString(prefix string) string {
	return intDecl(alts).indentedString(prefix, "adaptive-lease-time-threshold")
}

func (alts *adaptiveLeaseThresholdStatement) fromTokenStream(tStream *tokenStream) error {
	return (*intDecl)(alts).fromTokenStream(tStream, "adaptive-lease-time-threshold")
}

type alwaysBroadcastStatement bool

func (abs alwaysBroadcastStatement) indentedString(prefix string) string {
	return onOffBool(abs).indentedString(prefix, "always-broadcast")
}

func (abs *alwaysBroadcastStatement) fromTokenStream(tStream *tokenStream) error {
	return (*onOffBool)(abs).fromTokenStream(tStream, "always-broadcast")
}

type alwaysReplyRFC1048 bool

func (arrfc alwaysReplyRFC1048) indentedString(prefix string) string {
	return onOffBool(arrfc).indentedString(prefix, "always-reply-rfc1048")
}

func (arrfc *alwaysReplyRFC1048) fromTokenStream(tStream *tokenStream) error {
	return (*onOffBool)(arrfc).fromTokenStream(tStream, "always-reply-rfc1048")
}

type authoritativeStatement bool

func (as authoritativeStatement) indentedString(prefix string) string {
	if as == true {
		return prefix + "authoritative;\n"
	}
	return prefix + "not authoritative;\n"
}

func (as *authoritativeStatement) fromTokenStream(tStream *tokenStream) error {
	// Expect the "authoritative" identifier
	tok, err := tStream.next()
	if err != nil {
		return err
	}
	if strings.ToLower(string(tok.data)) != "authoritative" {
		return contextErrorf("identifier != 'authoritative': %q", string(tok.data))
	}

	// Rewind twice, look for a "not"
	tStream.undo()          // undo "authoritative"
	tStream.undo()          // undo whitespace
	tStream.undo()          // undo possibly "not"
	tok, _ = tStream.next() // re-consume "not"
	*as = true
	if tok.typ == tokenTypeIdentifier && strings.ToLower(string(tok.data)) == "not" {
		*as = false
	}

	// Don't forget to re-consume whitespace + "authoritative"
	tStream.next()
	tStream.next()

	return nil
}

type bootUnknownClientsStatement bool

func (bucs bootUnknownClientsStatement) indentedString(prefix string) string {
	return onOffBool(bucs).indentedString(prefix, "boot-unknown-clients")
}

func (bucs *bootUnknownClientsStatement) fromTokenStream(tStream *tokenStream) error {
	return (*onOffBool)(bucs).fromTokenStream(tStream, "boot-unknown-clients")
}

type dbTimeFormatStatement struct {
	localTime bool
}

func (dtfs dbTimeFormatStatement) indentedString(prefix string) string {
	var format string
	if dtfs.localTime {
		format = "local"
	} else {
		format = "default"
	}
	return prefix + "db-time-format " + format + ";\n"
}

func (dtfs *dbTimeFormatStatement) fromTokenStream(tStream *tokenStream) error {
	err := expectIdentifier(tStream, "db-time-format")
	if err != nil {
		return err
	}

	err = expectWhitespace(tStream)
	if err != nil {
		return err
	}

	return nil
}

type ddnsDomainNameStatement string

func (ddnsds ddnsDomainNameStatement) indentedString(prefix string) string {
	return stringDecl(ddnsds).indentedString(prefix, "ddns-domainname")
}

func (ddnsds *ddnsDomainNameStatement) fromTokenStream(tStream *tokenStream) error {
	return (*stringDecl)(ddnsds).fromTokenStream(tStream, "ddns-domainname")
}

type ddnsHostNameStatement string

func (ddnshns ddnsHostNameStatement) indentedString(prefix string) string {
	return stringDecl(ddnshns).indentedString(prefix, "ddns-hostname")
}

func (ddnshns *ddnsHostNameStatement) fromTokenStream(tStream *tokenStream) error {
	return (*stringDecl)(ddnshns).fromTokenStream(tStream, "ddns-hostname")
}

type ddnsRevDomainNameStatement string

func (ddnsrdns ddnsRevDomainNameStatement) indentedString(prefix string) string {
	return stringDecl(ddnsrdns).indentedString(prefix, "ddns-rev-domainname")
}

func (ddnsrdns *ddnsRevDomainNameStatement) fromTokenStream(tStream *tokenStream) error {
	return (*stringDecl)(ddnsrdns).fromTokenStream(tStream, "ddns-rev-domainname")
}

const (
	ddnsUpdateStyleNone = iota
	ddnsUpdateStyleAdHoc
	ddnsUpdateStyleInterim
)

type ddnsUpdateStyleStatement int

func (ddnsuss ddnsUpdateStyleStatement) indentedString(prefix string) string {
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

func (ddnsuss *ddnsUpdateStyleStatement) fromTokenStream(tStream *tokenStream) error {
	err := expectIdentifier(tStream, "ddns-update-style")
	if err != nil {
		return err
	}

	err = expectWhitespace(tStream)
	if err != nil {
		return err
	}

	tok, err := tStream.next()
	tokVal := strings.ToLower(string(tok.data))
	if tok.typ != tokenTypeIdentifier ||
		(tokVal != "ad-hoc" && tokVal != "interim" && tokVal != "none") {
		return contextErrorf("expected identifer, one of 'ad-hoc|interim|none', got %q", tokVal)
	}
	switch tokVal {
	case "ad-hoc":
		*ddnsuss = ddnsUpdateStyleAdHoc
	case "interim":
		*ddnsuss = ddnsUpdateStyleInterim
	default:
		*ddnsuss = ddnsUpdateStyleNone
	}
	return nil
}

type ddnsUpdatesStatement bool

func (ddnsus ddnsUpdatesStatement) indentedString(prefix string) string {
	return onOffBool(ddnsus).indentedString(prefix, "ddns-updates")
}

func (ddnsus *ddnsUpdatesStatement) fromTokenStream(tStream *tokenStream) error {
	return (*onOffBool)(ddnsus).fromTokenStream(tStream, "ddns-updates")
}

type defaultLeaseTimeStatement int

func (dlts defaultLeaseTimeStatement) indentedString(prefix string) string {
	return intDecl(dlts).indentedString(prefix, "default-lease-time")
}

func (dlts *defaultLeaseTimeStatement) fromTokenStream(tStream *tokenStream) error {
	return (*intDecl)(dlts).fromTokenStream(tStream, "default-lease-time")
}

type delayedAckStatement int

func (das delayedAckStatement) indentedString(prefix string) string {
	return intDecl(das).indentedString(prefix, "delayed-ack")
}

func (das *delayedAckStatement) fromTokenStream(tStream *tokenStream) error {
	return (*intDecl)(das).fromTokenStream(tStream, "delayed-ack")
}

type doForwardUpdatesStatement bool

func (dfus doForwardUpdatesStatement) indentedString(prefix string) string {
	return onOffBool(dfus).indentedString(prefix, "do-forward-updates")
}

func (dfus *doForwardUpdatesStatement) fromTokenStream(tStream *tokenStream) error {
	return (*onOffBool)(dfus).fromTokenStream(tStream, "do-forward-updates")
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

func (dblcs dynamicBootpLeaseCutoffStatement) indentedString(prefix string) string {
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

func (dblcs *dynamicBootpLeaseCutoffStatement) fromTokenStream(tStream *tokenStream) error {
	err := expectIdentifier(tStream, "dynamic-bootp-lease-cutoff")
	if err != nil {
		return err
	}

	err = expectWhitespace(tStream)
	if err != nil {
		return err
	}

	// W YYYY/MM/DD HH:MM:SS

	// Expect a day-of-week digit
	var dateBytes []byte
	tok, err := tStream.next()
	if err != nil {
		return err
	}
	if tok.typ != tokenTypeIdentifier {
		return contextErrorf("expected identifier, got %q", tok.data)
	}
	dateBytes = append(dateBytes, tok.data...)
	//Expect whitespace
	err = expectWhitespace(tStream)
	if err != nil {
		return err
	}
	dateBytes = append(dateBytes, []byte(" ")...)
	// Expect YYYY/MM/DD
	tok, err = tStream.next()
	if err != nil {
		return err
	}
	if tok.typ != tokenTypeIdentifier {
		return contextErrorf("expected identifier, got %q", tok.data)
	}
	dateBytes = append(dateBytes, tok.data...)
	// Expect whitespace
	err = expectWhitespace(tStream)
	if err != nil {
		return err
	}
	dateBytes = append(dateBytes, []byte(" ")...)
	// Expect HH:MM:SS
	tok, err = tStream.next()
	if err != nil {
		return err
	}
	if tok.typ != tokenTypeIdentifier {
		return contextErrorf("expected identifier, got %q", tok.data)
	}
	dateBytes = append(dateBytes, tok.data...)

	// Scan the string into our struct vals
	_, err = fmt.Sscanf(
		string(dateBytes),
		"%1d %04d/%02d/%02d %02d:%02d:%02d",
		&dblcs.dayOfWeek,
		&dblcs.year,
		&dblcs.month,
		&dblcs.dayOfMonth,
		&dblcs.hours,
		&dblcs.minutes,
		&dblcs.seconds,
	)
	if err != nil {
		return err
	}

	if dblcs.dayOfWeek < 0 || dblcs.dayOfWeek > 6 {
		return contextErrorf("unreasoable day of week value %d", dblcs.dayOfWeek)
	}
	if dblcs.year < 0 || dblcs.year > 9999 {
		return contextErrorf("unreasonable year value %d", dblcs.year)
	}
	if dblcs.month < 0 || dblcs.month > 12 {
		return contextErrorf("unreasonable month value %d", dblcs.month)
	}
	if dblcs.dayOfMonth < 0 || dblcs.dayOfMonth > 31 {
		return contextErrorf("unreasonable day of month value %d", dblcs.dayOfMonth)
	}
	if dblcs.hours < 0 || dblcs.hours > 23 {
		return contextErrorf("unreasonable hours value %d", dblcs.hours)
	}
	if dblcs.minutes < 0 || dblcs.minutes > 59 {
		return contextErrorf("unreasonable minutes value %d", dblcs.minutes)
	}
	if dblcs.seconds < 0 || dblcs.seconds > 59 {
		return contextErrorf("unreasonable seconds value %d", dblcs.seconds)
	}

	return nil
}

type fixedAddressStatement []net.IP

func (fas fixedAddressStatement) indentedString(prefix string) string {
	s := prefix + "fixed-address "
	for i := 0; i < len(fas)-1; i++ {
		s += fas[i].String() + ", "
	}
	return s + fas[len(fas)-1].String() + ";\n"
}

func (fas *fixedAddressStatement) fromTokenStream(tStream *tokenStream) error {
	// Consume "fixed-address"
	tok, err := tStream.next()
	if err != nil {
		return err
	}
	if strings.ToLower(string(tok.data)) != "fixed-address" {
		return contextErrorf("identifer != 'fixed-address': %q", string(tok.data))
	}

	// Consume whitespace + IP addresses, until we hit a semicolon
OUTER:
	for {
		tok, err = tStream.next()
		if err != nil {
			return err
		}
		switch tok.typ {
		case tokenTypeIdentifier:
			tokVal := strings.ToLower(string(tok.data))
			tokVal = strings.TrimSuffix(tokVal, ",")
			if !ipAddressRegexp.MatchString(tokVal) {
				return contextErrorf("malformed IP address %q", tokVal)
			}
			*fas = append(*fas, net.ParseIP(tokVal))
		case tokenTypeStatementEnd:
			break OUTER
		case tokenTypeWhiteSpace:
			fallthrough
		default:
			continue
		}
	}

	return nil
}

type hardwareStatement struct {
	hardwareType    string
	hardwareAddress string
}

func (hs hardwareStatement) indentedString(prefix string) string {
	return prefix + "hardware " + hs.hardwareType + " " + hs.hardwareAddress + ";\n"
}

func (hs *hardwareStatement) fromTokenStream(tStream *tokenStream) error {
	// Expect "hardware" identifier
	tok, err := tStream.next()
	if err != nil {
		return err
	}
	if strings.ToLower(string(tok.data)) != "hardware" {
		return contextErrorf("identifer != 'hardware': %q", string(tok.data))
	}

	// Expect whitespace
	err = expectWhitespace(tStream)
	if err != nil {
		return contextError(err.Error())
	}

	// Expect a hardware type, either "ethernet" or "token-ring"
	tok, err = tStream.next()
	if err != nil {
		return err
	}
	tokVal := strings.ToLower(string(tok.data))
	if tokVal != "ethernet" && tokVal != "token-ring" {
		return contextErrorf("expected either 'ethernet' or 'token-ring', got %q", string(tok.data))
	}
	hs.hardwareType = tokVal

	// Expect whitespace
	err = expectWhitespace(tStream)
	if err != nil {
		return contextError(err.Error())
	}

	// Expect a hardware address
	tok, err = tStream.next()
	if err != nil {
		return err
	}
	if !hwAddressRegexp.Match(tok.data) {
		return contextErrorf("malformed hardware address %q", string(tok.data))
	}
	hs.hardwareAddress = string(tok.data)

	return nil
}

type maxAckDelayStatement int

func (mads maxAckDelayStatement) indentedString(prefix string) string {
	return intDecl(mads).indentedString(prefix, "max-ack-delay")
}

func (mads *maxAckDelayStatement) fromTokenStream(tStream *tokenStream) error {
	return (*intDecl)(mads).fromTokenStream(tStream, "max-ack-delay")
}

type useHostDeclNamesStatement bool

func (uhdns useHostDeclNamesStatement) indentedString(prefix string) string {
	return onOffBool(uhdns).indentedString(prefix, "use-host-decl-names")
}

func (uhdns *useHostDeclNamesStatement) fromTokenStream(tStream *tokenStream) error {
	return (*onOffBool)(uhdns).fromTokenStream(tStream, "use-host-decl-names")
}

// OPTIONS

type domainNameServersOption []net.IP

func (dnso domainNameServersOption) indentedString(prefix string) string {
	s := prefix + "option domain-name-servers "
	for i := 0; i < len(dnso)-1; i++ {
		s += dnso[i].String() + ", "
	}
	return s + dnso[len(dnso)-1].String() + ";\n"
}

func (dnso *domainNameServersOption) fromTokenStream(tStream *tokenStream) error {
	// Expect "domain-name-servers" identifier
	tok, err := tStream.next()
	if err != nil {
		return err
	}
	if strings.ToLower(string(tok.data)) != "domain-name-servers" {
		return contextErrorf("identifer != 'domain-name-servers': %q", string(tok.data))
	}

	// Expect whitespace
	err = expectWhitespace(tStream)
	if err != nil {
		return contextError(err.Error())
	}

	// Consume whitespace + IP addresses, until we hit a semicolon
OUTER:
	for {
		tok, err = tStream.next()
		if err != nil {
			return err
		}
		switch tok.typ {
		case tokenTypeIdentifier:
			tokVal := strings.ToLower(string(tok.data))
			tokVal = strings.TrimSuffix(tokVal, ",")
			if !ipAddressRegexp.MatchString(tokVal) {
				return contextErrorf("malformed IP address %q", tokVal)
			}
			*dnso = append(*dnso, net.ParseIP(tokVal))
		case tokenTypeStatementEnd:
			break OUTER
		case tokenTypeWhiteSpace:
			fallthrough
		default:
			continue
		}
	}

	return nil
}

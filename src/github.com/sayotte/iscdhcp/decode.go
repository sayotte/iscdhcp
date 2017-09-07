package iscdhcp

import (
	"fmt"
	"strings"
)

// Decoder analyzes a slice of bytes, constructing primitive ISC-DHCP config
// objects.
type Decoder struct{}

func (d Decoder) Decode(tStream *tokenStream) ([]statement, error) {
	var statements []statement
	for tStream.index < len(tStream.tokens) {
		tok, err := tStream.next()
		if err != nil {
			return nil, err
		}
		if tok.typ != tokenTypeIdentifier {
			continue
		}

		var statement statement
		tokVal := strings.ToLower(string(tok.data))
		switch tokVal {
		case "adaptive-lease-time-threshold":
			var altts adaptiveLeaseThresholdStatement
			statement = &altts
		case "authoritative":
			var asReal authoritativeStatement
			statement = &asReal
		case "always-broadcast":
			var abs alwaysBroadcastStatement
			statement = &abs
		case "always-reply-rfc1048":
			var arrfc alwaysReplyRFC1048
			statement = &arrfc
		case "boot-unknown-clients":
			var bucsReal bootUnknownClientsStatement
			statement = &bucsReal
		case "ddns-domainname":
			var ddnsDomainname ddnsDomainNameStatement
			statement = &ddnsDomainname
		case "ddns-hostname":
			var ddnsHostname ddnsHostNameStatement
			statement = &ddnsHostname
		case "ddns-rev-domainname":
			var ddnsRevDomainname ddnsRevDomainNameStatement
			statement = &ddnsRevDomainname
		case "ddns-update-style":
			var ddnsUpdateStyle ddnsUpdateStyleStatement
			statement = &ddnsUpdateStyle
		case "ddns-updates":
			var ddnsuReal ddnsUpdatesStatement
			statement = &ddnsuReal
		case "default-lease-time":
			var dlts defaultLeaseTimeStatement
			statement = &dlts
		case "delayed-ack":
			var das delayedAckStatement
			statement = &das
		case "do-forward-updates":
			var dfus doForwardUpdatesStatement
			statement = &dfus
		case "dynamic-bootp-lease-cutoff":
			statement = &dynamicBootpLeaseCutoffStatement{}
		case "fixed-address":
			statement = &fixedAddressStatement{}
		case "group":
			statement = &groupStatement{}
		case "hardware":
			statement = &hardwareStatement{}
		case "host":
			statement = &hostStatement{}
		case "include":
			statement = &includeStatement{}
		case "max-ack-delay":
			var mads maxAckDelayStatement
			statement = &mads
		case "option":
			tStream.undo()
			statement, err := d.decodeOption(tStream)
			if err != nil {
				return nil, err
			}
			statements = append(statements, statement)
			continue
		case "use-host-decl-names":
			var uhdnsReal useHostDeclNamesStatement
			statement = &uhdnsReal
		case "shared-network":
			statement = &sharedNetworkStatement{}
		case "subnet":
			statement = &subnetStatement{}
		default:
			return nil, contextErrorf("unsupported identifier %q", tokVal)
		}
		tStream.undo()
		err = statement.fromTokenStream(tStream)
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	return statements, nil
}

func (d Decoder) decodeOption(tStream *tokenStream) (statement, error) {
	var statement statement

	// Expect "option" identifier
	tok, err := tStream.next()
	if err != nil {
		return statement, err
	}
	if strings.ToLower(string(tok.data)) != "option" {
		return statement, contextErrorf("identifier != 'option': %q", string(tok.data))
	}

	// Expect whitespace
	err = expectWhitespace(tStream)
	if err != nil {
		return statement, contextError(err.Error())
	}

	// Determine what kind of statement this is, populate it
	tok, err = tStream.next()
	if err != nil {
		return statement, err
	}
	tokVal := strings.ToLower(string(tok.data))
	switch tokVal {
	case "domain-name-servers":
		statement = &domainNameServersOption{}
	default:
		return statement, contextErrorf("unrecognized option %q", string(tok.data))
	}
	tStream.undo()
	err = statement.fromTokenStream(tStream)

	return statement, err
}

func expectIdentifier(tStream *tokenStream, lowerCaseIdentifier string) error {
	tok, err := tStream.next()
	if err != nil {
		return err
	}
	if tok.typ != tokenTypeIdentifier {
		return fmt.Errorf("expected an identifier, got %q", string(tok.data))
	}
	lcTokVal := strings.ToLower(string(tok.data))
	if lcTokVal != lowerCaseIdentifier {
		return fmt.Errorf("expected identifier %q, got %q", lowerCaseIdentifier, lcTokVal)
	}
	return nil
}

func expectWhitespace(tStream *tokenStream) error {
	tok, err := tStream.next()
	if err != nil {
		return err
	}
	if tok.typ != tokenTypeWhiteSpace {
		return fmt.Errorf("expected whitespace, got %q", string(tok.data))
	}
	return nil
}

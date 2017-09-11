[![Build Status](https://travis-ci.org/sayotte/iscdhcp.svg?branch=master)](https://travis-ci.org/sayotte/iscdhcp)

# Brief
This is a config-file parser and generator library for ISC-DHCP, written in the Go programming language.

The parser is generated using the *goyacc* utility.

Supported statements to date include a common subset of those used by DHCPd. The code should be easy
to extend.

# Usage
### Parsing
```go
fd, err := os.Open(fileName)
if err != nil {
    log.Fatalf("os.Open(%q): %s", fileName, err)
}

statements, err := iscdhcp.Decode(fd)
if err != nil {
    log.Fatalf("iscdhcp.Decode(): %s", err)
}
```

### Generating
```go
hs := HostStatement{
    Hostname: "serverA.myDomain.tld",
}
hs.Statements = append(hs.Statements, HardwareStatement{"ethernet", "0:1:2:3:4:5"})
hs.Statements = append(hs.Statements, FixedAddressStatement{net.ParseIP("1.2.3.4")})
hs.Statements = append(hs.Statements, IncludeStatement{"filename.cfg"})

subnetStmt := SubnetStatement{
    SubnetNumber: net.ParseIP("1.2.3.0"),
    Netmask:      net.ParseIP("255.255.255.0"),
}
subnetStmt.Statements = append(subnetStmt.Statements, hs)

fmt.Println(subnetStmt.IndentedString(""))
```

The above yields this:
```
subnet 1.2.3.0 netmask 255.255.255.0 {
    host serverA.myDomain.tld {
        hardware ethernet 0:1:2:3:4:5;
        fixed-address 1.2.3.4;
        include "filename.cfg";
    }
}
```

# Modifying
New statements must first be defined in `statements.go`, and satisfy the
`iscdhcp.Statement` interface.

To update the parser, edit the *goyacc* grammar definitions in the `parser.y`
file. After updating these you must call `goyacc parser.y` to regenerate the
`y.go` file, after which you can rebuild the package.

To update generator code, edit the `IndentedString()` method of the statements
with which you're concerned.

It's best to modify both parser- and generator-code in lockstep, and add 
"round-trip" tests for any new statements introduced to ensure that the
generated output is understandable to the parser. Round-trip tests are found in
`statements_test.go`. 

# Contributing
Contributions should be submitted as pull requests from a named (i.e. not
*master*) branch.

Running `make` (and the unit tests and linter checks it entails) on the result
of merging or rebasing your PR branch over current *master* branch must exit
successfully.


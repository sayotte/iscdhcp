%{
package iscdhcp

import (
    "fmt"
    "net"
    "strings"
)

%}

%token openBrace closeBrace quote semicolon comma

// conditional + boolean logic stuff
%token ConditionIf ConditionElsif ConditionElse
// Boolean combinations should be right-associative, so that e.g. this statement:
//    if x and y or not z
// ... will be grouped like this:
//    if x and (y or not z)
// ... as opposed to this:
//    if (x and y) or not z
// This gives short-circuit behavior if x!=true, which is what most users will
// expect. I have not verified this is actually how dhcpd behaves, though.
%right BoolAnd BoolOr BoolNot
%token BoolEqual BoolInequal BoolRegexMatch BoolRegexIMatch
%token BoolExists BoolKnown BoolStatic

// simple types
%token number ipAddr cidr stringConst macAddr

// reserved words
%token stateTok authoritativeTok
%token groupTok hostTok subnetTok netmaskTok optionTok includeTok
%token hardwareTok ethernetTok fixedAddrTok
%token useHostDeclNamesTok
%token optDomainNameServersTok

// everything else
%token word

// compound yyType definition, so we can turn any token into any of these types
// of object
%union {
    num int
    str string
    strList []string
    ipList []net.IP
    ip net.IP
    statement Statement
    statementList []Statement
    dataTerm fmt.Stringer
    boolExpr BooleanExpression
    subConditionals []ConditionalStatement
}

%%
// Primitives
config: statements
    {
        l := yylex.(*lexer)
        l.dirtyHackReturn = $1.statementList
    };

statements:
    statement
    {
        $$.statementList = []Statement{$1.statement}
    }
    | statements statement
    {
        $$.statementList = append($$.statementList, $2.statement)
    };

statement:
    // Statements can be either declarations...
    groupdecl
    | hostdecl
    | includedecl
    | subnetdecl
    | conditionalDecl

    // or parameters
    | authoritativeParam
    | hardwareparam
    | fixedaddressparam
    | optionparam
    | useHostDeclNamesParam
    ;

block:
      openBrace closeBrace // empty block
    | openBrace statements closeBrace
    {
        $$.statementList = $2.statementList
    };

// Conditionals
conditionalDecl:
    ConditionIf booleanExpr block subConditionList
    {
        cs := ConditionalStatement {
            Operator:        ConditionIf,
            Condition:       $2.boolExpr,
            Statements:      $3.statementList,
            SubConditionals: $4.subConditionals,
        }
        $$.statement = cs
    }

subConditionList:
    // empty list is one possibility
    | subConditionList ConditionElsif booleanExpr block
    {
        cs := ConditionalStatement {
            Operator:   ConditionElsif,
            Condition:  $3.boolExpr,
            Statements: $3.statementList,
        }
        $$.subConditionals = append($$.subConditionals, cs)
    }
    | subConditionList ConditionElse block
    {
        cs := ConditionalStatement {
            Operator:   ConditionElse,
            Statements: $2.statementList,
        }
        $$.subConditionals = append($$.subConditionals, cs)
    };

booleanExpr:
    booleanExpr BoolAnd booleanExpr
    {
        $$.boolExpr = BooleanExpression {
            Operator: BoolAnd,
            BoolTerms: []BooleanExpression{$1.boolExpr, $3.boolExpr},
        }
    }
    | booleanExpr BoolOr booleanExpr
    {
        $$.boolExpr = BooleanExpression {
            Operator: BoolOr,
            BoolTerms: []BooleanExpression{$1.boolExpr, $3.boolExpr},
        }
    }
    | BoolNot booleanExpr
    {
        $$.boolExpr = BooleanExpression {
            Operator: BoolNot,
            BoolTerms: []BooleanExpression{$2.boolExpr},
        }
    }
    | BoolStatic
    {
        $$.boolExpr = BooleanExpression {
            Operator: BoolStatic,
        }
    }
    | BoolKnown
    {
        $$.boolExpr = BooleanExpression {
            Operator: BoolKnown,
        }
    }
    | BoolExists dataTerm
    {
        $$.boolExpr = BooleanExpression {
            Operator: BoolExists,
            DataTerms: []fmt.Stringer{$2.dataTerm},
        }
    }
    | dataTerm BoolEqual dataTerm
    {
        $$.boolExpr = BooleanExpression {
            Operator: BoolEqual,
            DataTerms: []fmt.Stringer{$1.dataTerm, $3.dataTerm},
        }
    }
    | dataTerm BoolInequal dataTerm
    {
        $$.boolExpr = BooleanExpression {
            Operator: BoolInequal,
            DataTerms: []fmt.Stringer{$1.dataTerm, $3.dataTerm},
        }
    }
    | dataTerm BoolRegexMatch dataTerm
    {
        $$.boolExpr = BooleanExpression {
            Operator: BoolRegexMatch,
            DataTerms: []fmt.Stringer{$1.dataTerm, $3.dataTerm},
        }
    }
    | dataTerm BoolRegexIMatch dataTerm
    {
        $$.boolExpr = BooleanExpression {
            Operator: BoolRegexIMatch,
            DataTerms: []fmt.Stringer{$1.dataTerm, $3.dataTerm},
        }
    };

dataTerm:
    stringConst
    {
        $$.dataTerm = StringConstTerm($1.str)
    }
    | optionTok word
    {
        $$.dataTerm = PacketOptionTerm{
            optionName: $2.str,
        }
    };

//wordList:
//    wordList word
//    {
//        $$.strList = append($$.strList, $2.str)
//    }
//    | word
//    {
//        $$.strList = []string{$1.str}
//    };

ipList:
    ipList comma ipAddr
    {
        $$.ipList = append($$.ipList, net.ParseIP($3.str))
    }
    | ipAddr
    {
        $$.ipList = []net.IP{net.ParseIP($1.str)}
    };

// Declarations that include a block
groupdecl: groupTok block
    {
        gs := GroupStatement{
            Statements: $2.statementList,
        }
        $$.statement = gs
    };

hostdecl: hostTok word block
    {
        hs := HostStatement {
            Hostname:   $2.str,
            Statements: $3.statementList,
        }
        $$.statement = hs
    };

includedecl: includeTok stringConst semicolon
    {
        is := IncludeStatement {
            Filename: $2.str,
        }
        $$.statement = is
    };

subnetdecl: subnetTok ipAddr netmaskTok ipAddr block
    {
        sns := SubnetStatement {
            SubnetNumber: net.ParseIP($2.str),
            Netmask:      net.ParseIP($4.str),
            Statements:   $5.statementList,
        }
        $$.statement = sns
    };

// Parameters found within a block
authoritativeParam:
    BoolNot authoritativeTok semicolon
    {
        $$.statement = AuthoritativeStatement(false)
    }
    | authoritativeTok semicolon
    {
        $$.statement = AuthoritativeStatement(true)
    };

hardwareparam:
    hardwareTok ethernetTok macAddr semicolon
    {
        $$.statement = HardwareStatement {
            HardwareType: "ethernet",
            HardwareAddress: $3.str,
        }
    };

fixedaddressparam:
    fixedAddrTok ipList semicolon
    {
        $$.statement = FixedAddressStatement($2.ipList)
    };

useHostDeclNamesParam:
    useHostDeclNamesTok stateTok semicolon
    {
        val := false
        cmpText := strings.ToLower($2.str)
        if cmpText == "on" || cmpText == "true" {
            val = true
        }
        $$.statement = UseHostDeclNamesStatement(val)
    };

// Options, because they're weird
optionparam: optionTok optionClause
    {
        $$.statement = $2.statement
    };

optionClause:
    nameserversOptClause;

nameserversOptClause: optDomainNameServersTok ipList semicolon
    {
        $$.statement = DomainNameServersOption($2.ipList)
    };
%%

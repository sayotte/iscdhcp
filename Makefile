SHELL=/bin/bash
ROOT = $(shell pwd)
export GOPATH := ${ROOT}
export PATH := ${PATH}:${ROOT}/bin
TEST_FLAGS ?= -cover
YACC = bin/goyacc
YACC_FLAGS = -l

.PHONY: fmt build test lint

all: fmt build test lint

fmt:
	rm -f y.go
	go fmt *.go

build: y.go
	go build

y.go: parse.y ${YACC}
	${YACC} -o y.go ${YACC_FLAGS} parse.y

${YACC}:
	go get golang.org/x/tools/cmd/goyacc

lint: bin/gometalinter
	bin/gometalinter . -D gocyclo --exclude='y.go' --exclude='unused'

bin/gometalinter: 
	go get github.com/alecthomas/gometalinter
	bin/gometalinter --install

test:
	go test ${TEST_FLAGS}

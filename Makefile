SHELL=/bin/bash
ROOT = $(shell pwd)
export GOPATH := ${ROOT}
export PATH := ${PATH}:${ROOT}/bin
PROJ_ROOT = github.com/sayotte/iscdhcp
SRC_ROOT = src/${PROJ_ROOT}
TEST_FLAGS ?= -cover
YACC = bin/goyacc
YACC_FLAGS =

.PHONY: fmt install test lint

all: fmt install test lint

fmt:
	rm -f ${SRC_ROOT}/y.go
	go fmt ${PROJ_ROOT}/...

install: ${SRC_ROOT}/y.go
	go install ${PROJ_ROOT}/...

${SRC_ROOT}/y.go: ${SRC_ROOT}/parse.y ${YACC}
	${YACC} -o ${SRC_ROOT}/y.go ${YACC_FLAGS} ${SRC_ROOT}/parse.y

${YACC}:
	go get golang.org/x/tools/cmd/goyacc

lint: bin/gometalinter
	bin/gometalinter ${SRC_ROOT}/... -D gocyclo --exclude='y.go' --exclude='unused'

bin/gometalinter: 
	go get github.com/alecthomas/gometalinter
	bin/gometalinter --install

test:
	go test ${TEST_FLAGS} ${PROJ_ROOT}/...


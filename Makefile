SHELL=/bin/bash
ROOT = $(shell pwd)
export GOPATH := ${ROOT}
PROJ_ROOT = github.com/sayotte/iscdhcp
TEST_FLAGS ?= -cover

test:
	go test ${TEST_FLAGS} ${PROJ_ROOT}/...

install:
	go fmt ${PROJ_ROOT}/...
	go install ${PROJ_ROOT}/...


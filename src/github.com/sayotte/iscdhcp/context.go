package iscdhcp

import (
	"fmt"
	"runtime"
)

func callContext(skip int) string {
	pc, file, lineNum, ok := runtime.Caller(skip + 1)
	if !ok {
		return "(unknown-func)"
	}
	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("[%s:%d]:%s()", file, lineNum, funcName)
}

func contextError(msg string) error {
	return fmt.Errorf("%s: %s", callContext(1), msg)
}

func contextErrorf(msg string, a ...interface{}) error {
	msgInner := fmt.Sprintf("%s: %s", callContext(1), msg)
	return fmt.Errorf(msgInner, a)
}

package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []string
}

func (e *MultiError) Error() string {
	return fmt.Sprintf("%d errors occured:\n\t* %s\n", len(e.errs), strings.Join(e.errs, "\t* "))
}

func Append(err error, errs ...error) *MultiError {
	var (
		me *MultiError
		ok bool
	)

	if err == nil {
		me = &MultiError{}
	} else if me, ok = err.(*MultiError); !ok {
		me = &MultiError{errs: []string{err.Error()}}
	}

	for _, e := range errs {
		me.errs = append(me.errs, e.Error())
	}

	return me
}

func TestMultiError(t *testing.T) {
	var err = errors.New("error 0")
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))
	err = Append(err, errors.New("error 3"))

	expectedMessage := "4 errors occured:\n\t* error 0\t* error 1\t* error 2\t* error 3\n"
	assert.EqualError(t, err, expectedMessage)
}

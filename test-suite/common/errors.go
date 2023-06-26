package common

import (
	"github.com/pkg/errors"
	"log"
	"strings"
)

func NoError(err error, msgAndArgs ...interface{}) {
	if err != nil {
		log.Panicf("expect no error, got %v", err)
	}
}

func IsError(actualError, expectedError error, msgAndArgs ...interface{}) {
	if errors.Is(actualError, expectedError) {
		log.Panicf("expect error %v, got %v", expectedError, actualError)
	}
}

func ErrorContains(err error, msg string, msgAndArgs ...interface{}) {
	if err == nil {
		log.Panicf("expect an error, got nil")
	}
	if !strings.Contains(err.Error(), msg) {
		log.Panicf("expect error %v to contain message `%v`", err, msg)
	}
}

package common

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"strings"
)

func parseMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}

func NoError(err error, msgAndArgs ...interface{}) {
	if err != nil {
		log.Panicf("expect no error, got %v\n%v", err, parseMsgAndArgs(msgAndArgs))
	}
}

func True(in bool, msgAndArgs ...interface{}) {
	if !in {
		log.Panicf("expect true, got false\n%v", parseMsgAndArgs(msgAndArgs))
	}
}

func IsError(actualError, expectedError error, msgAndArgs ...interface{}) {
	if errors.Is(actualError, expectedError) {
		log.Panicf("expect error %v, got %v\n%v", expectedError, actualError, parseMsgAndArgs(msgAndArgs))
	}
}

func ErrorContains(err error, msg string, msgAndArgs ...interface{}) {
	if err == nil {
		log.Panicf("expect an error, got nil\n%v", parseMsgAndArgs(msgAndArgs))
	}
	if !strings.Contains(err.Error(), msg) {
		log.Panicf("expect error %v to contain message `%v`\n%v", err, msg, parseMsgAndArgs(msgAndArgs))
	}
}

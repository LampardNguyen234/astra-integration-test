package common

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type ComparisonMatcher struct {
	expected interface{}
	op       CompareOP
}

func (matcher *ComparisonMatcher) Match(a interface{}) (success bool, err error) {
	if a == nil && matcher.expected == nil {
		return false, fmt.Errorf("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.")
	}
	ret := false
	switch matcher.expected.(type) {
	case sdk.Int:
		if tmp, ok := a.(sdk.Int); ok {
			ret = matcher.op.CompareSdkInt(tmp, matcher.expected.(sdk.Int))
		}
	case sdk.Dec:
		if tmp, ok := a.(sdk.Dec); ok {
			ret = matcher.op.CompareSdkDec(tmp, matcher.expected.(sdk.Dec))
		}
	case sdk.Coin:
		if tmp, ok := a.(sdk.Coin); ok {
			ret = matcher.op.CompareSdkInt(tmp.Amount, matcher.expected.(sdk.Coin).Amount)
			ret = ret && tmp.Denom == matcher.expected.(sdk.Coin).Denom
		}
	case uint64:
		if tmp, ok := a.(uint64); ok {
			ret = matcher.op.CompareUint(tmp, matcher.expected.(uint64))
		}
	case float64:
		if tmp, ok := a.(float64); ok {
			ret = matcher.op.CompareFloat64(tmp, matcher.expected.(float64))
		}
	default:
	}

	return ret, nil
}

func (matcher *ComparisonMatcher) FailureMessage(actual interface{}) (message string) {
	actualString, actualOK := actual.(string)
	expectedString, expectedOK := matcher.expected.(string)
	if actualOK && expectedOK {
		return format.MessageWithDiff(actualString, fmt.Sprintf("to %v", matcher.op.String()), expectedString)
	}

	return format.Message(actual, fmt.Sprintf("to %v", matcher.op.String()), matcher.expected)
}

func (matcher *ComparisonMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, fmt.Sprintf("not to %v", matcher.op.String()), matcher.expected)
}

func EQ(expected interface{}) types.GomegaMatcher {
	return &ComparisonMatcher{
		expected: expected,
		op:       OpEQ,
	}
}

func LTE(expected interface{}) types.GomegaMatcher {
	return &ComparisonMatcher{
		expected: expected,
		op:       OpLTE,
	}
}

func LT(expected interface{}) types.GomegaMatcher {
	return &ComparisonMatcher{
		expected: expected,
		op:       OpLT,
	}
}

func GT(expected interface{}) types.GomegaMatcher {
	return &ComparisonMatcher{
		expected: expected,
		op:       OpGT,
	}
}

func GTE(expected interface{}) types.GomegaMatcher {
	return &ComparisonMatcher{
		expected: expected,
		op:       OpGTE,
	}
}

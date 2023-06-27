package assert

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"log"
)

type CompareOP uint8

const (
	OpLT CompareOP = iota
	OpLTE
	OpEQ
	OpGTE
	OpGT
)

func (op CompareOP) String() string {
	switch op {
	case OpLT:
		return "less than"
	case OpLTE:
		return "less than or equal to"
	case OpEQ:
		return "equal to"
	case OpGTE:
		return "greater than or equal to"
	case OpGT:
		return "greater than"
	default:
		return ""
	}
}

func (op CompareOP) CompareFloat64(a, b float64) bool {
	switch op {
	case OpLT:
		return a < b
	case OpLTE:
		return a <= b
	case OpEQ:
		return a == b
	case OpGTE:
		return a >= b
	case OpGT:
		return a > b
	default:
		return false
	}
}

func (op CompareOP) CompareUint(a, b uint64) bool {
	return op.CompareFloat64(float64(a), float64(b))
}

func (op CompareOP) CompareSdkInt(a, b sdk.Int) bool {
	switch op {
	case OpLT:
		return a.LT(b)
	case OpLTE:
		return a.LTE(b)
	case OpEQ:
		return a.Equal(b)
	case OpGTE:
		return a.GTE(b)
	case OpGT:
		return a.GT(b)
	default:
		return false
	}
}

func (op CompareOP) CompareSdkDec(a, b sdk.Dec) bool {
	switch op {
	case OpLT:
		return a.LT(b)
	case OpLTE:
		return a.LTE(b)
	case OpEQ:
		return a.Equal(b)
	case OpGTE:
		return a.GTE(b)
	case OpGT:
		return a.GT(b)
	default:
		return false
	}
}

// Compare performs op.Compare(a, b).
func Compare(a, b interface{}, op CompareOP) {
	ret := false
	if _, ok := a.(sdk.Int); ok {
		ret = op.CompareSdkInt(a.(sdk.Int), b.(sdk.Int))
	}
	if _, ok := a.(sdk.Dec); ok {
		ret = op.CompareSdkDec(a.(sdk.Dec), b.(sdk.Dec))
	}
	if _, ok := a.(uint64); ok {
		ret = op.CompareUint(a.(uint64), b.(uint64))
	}
	if _, ok := a.(float64); ok {
		ret = op.CompareFloat64(a.(float64), b.(float64))
	}

	if !ret {
		log.Panicf("%v is not %v %v", a, op.String(), b)
	}
}

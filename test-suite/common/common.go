package common

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

func (op CompareOP) Compare(a, b float64) bool {
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

package common

import (
	"encoding/hex"
	"github.com/LampardNguyen234/astra-go-sdk/account"
	"math/big"
)

// RandBytes returns random a byte-slice given its length.
func RandBytes(n uint) []byte {
	ret := make([]byte, 0)
	if n == 0 {
		return ret
	}

	for uint(len(ret)) < n {
		ret = append(ret, RandomHash().Bytes()...)
	}

	return ret[:n]
}

// RandInt returns a random int.
func RandInt() int {
	return int(RandInt64())
}

// RandInt64 returns a random int64.
func RandInt64() int64 {
	return new(big.Int).SetBytes(RandBytes(32)).Int64()
}

// RandInterval returns a random int64 x where a <= x < b && a < b.
func RandInterval(a, b int64) int64 {
	if b <= a {
		return 0
	}
	return a + RandInt64()%(b-a)
}

// RandUint64 returns a random uint64.
func RandUint64() uint64 {
	return new(big.Int).SetBytes(RandBytes(32)).Uint64()
}

// RandUint returns a random uint.
func RandUint() uint {
	return uint(RandUint64())
}

// RandUInterval returns a random uint64 x where a <= x < b && a < b.
func RandUInterval(a, b uint64) uint64 {
	if b <= a {
		return 0
	}
	return a + RandUint64()%(b-a)
}

// RandKeyInfo returns a random key info.
func RandKeyInfo() *account.KeyInfo {
	ret, _ := account.NewKeyInfoFromPrivateKey(hex.EncodeToString(RandBytes(32)))

	return ret
}

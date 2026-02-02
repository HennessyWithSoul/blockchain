package util

import (
	"goblockchain/types"
	"math/rand"
)

func RandomBytes(size int) []byte {
	buf := make([]byte, size)
	rand.Read(buf)
	return buf
}
func RandomHash() types.Hash {
	return types.HashFromBytes(RandomBytes(32))
}

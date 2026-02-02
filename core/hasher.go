package core

import (
	"crypto/sha256"
	"encoding/binary"
	"goblockchain/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}
type BlockHasher struct{}

func (BlockHasher) Hash(header *Header) types.Hash {
	h := sha256.Sum256(header.Bytes())
	return types.Hash(h)
}

type TxHasher struct {
}

func (TxHasher) Hash(tx *Transaction) types.Hash {

	hasher := sha256.New()

	// 1. 交易类型
	binary.Write(hasher, binary.BigEndian, int32(tx.Type))
	// 4. 金额
	binary.Write(hasher, binary.BigEndian, tx.Value)
	// 5. 数据
	if len(tx.Data) > 0 {
		hasher.Write(tx.Data)
	}
	// 6. 发送方公钥
	fromBytes := tx.From.ToSlice()
	hasher.Write(fromBytes)

	// 7. 接收方公钥
	toBytes := tx.To.ToSlice()
	hasher.Write(toBytes)
	var hash types.Hash
	copy(hash[:], hasher.Sum(nil))
	tx.hash = hash // 缓存计算结果

	return hash
	
}

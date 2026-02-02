package core

import (
	"fmt"
	"goblockchain/types"
)

type Storage interface {
	Put(*Block) error
	PutCollection(hash types.Hash, tx *CollectionTx) error
	GetCollection(hash types.Hash) (*CollectionTx, error)
}
type MemoryStorage struct {
	collectionState map[types.Hash]*CollectionTx
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		collectionState: make(map[types.Hash]*CollectionTx),
	}

}
func (ms *MemoryStorage) Put(b *Block) error {
	return nil
}
func (ms *MemoryStorage) PutCollection(hash types.Hash, tx *CollectionTx) error {
	ms.collectionState[hash] = tx
	return nil
}
func (ms *MemoryStorage) GetCollection(hash types.Hash) (*CollectionTx, error) {
	if ms.collectionState[hash] == nil {
		return nil, fmt.Errorf("collection does not exist")
	}
	return ms.collectionState[hash], nil
}

package core

import (
	"github.com/stretchr/testify/assert"
	"goblockchain/types"
	"testing"
)

func newBlockchainWithGenesis(t *testing.T) *BlockChain {
	bc, err := NewBlockChain(RandomBlock(t, 0, types.Hash{}))
	assert.Nil(t, err)
	return bc
}
func TestNewBlockChain(t *testing.T) {
	//bc, err := NewBlockChain(RandomBlock(0))
	bc := newBlockchainWithGenesis(t)
	//assert.Nil(t, err)
	assert.NotNil(t, bc.validator)
	assert.NotNil(t, bc.Height(), uint32(0))
}
func TestHasBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.True(t, bc.HasBlock(0))
}
func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	//fmt.Println(bc.GetHeader(0))
	lenblocks := 1000
	for i := 0; i < lenblocks; i++ {
		block := RandomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
	}
	assert.Equal(t, uint32(lenblocks), bc.Height())
	assert.NotNil(t, bc.AddBlock(RandomBlock(t, uint32(lenblocks+1), types.Hash{})))
}
func TestAddBlockToHigh(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc.AddBlock(randomBlockWithSignature(t, uint32(3), types.Hash{})))
}
func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenblocks := 1000
	for i := 0; i < lenblocks; i++ {
		block := RandomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(uint32(i + 1))
		assert.Nil(t, err)
		assert.Equal(t, block.Header, header)
	}
}
func getPrevBlockHash(t *testing.T, bc *BlockChain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevHeader)
}

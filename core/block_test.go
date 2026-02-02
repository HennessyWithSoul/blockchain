package core

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"goblockchain/crypto"
	"goblockchain/types"
	"testing"
	"time"
)

func randomBlockWithSignature(t *testing.T, height uint32, preBlockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	b := RandomBlock(t, height, preBlockHash)
	tx := randomTxWithSignature(t)
	b.AddTransaction(&tx)
	assert.Nil(t, b.Sign(privKey))
	return b
}
func TestHashBlock(t *testing.T) {
	b := RandomBlock(t, 3, types.Hash{})
	fmt.Println(b.Hash(BlockHasher{}))
}
func TestSignBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := RandomBlock(t, 0, types.Hash{})
	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
	assert.Nil(t, b.Verify())
	//fmt.Println(b.HeaderData())
	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t, b.Verify())
	//	b.Validator = privKey.PublicKey()
	b.Header.Height = 100
	//fmt.Println(b.HeaderData())
	assert.NotNil(t, b.Verify())

}
func TestDecodeEncodeBlock(t *testing.T) {
	b := RandomBlock(t, 1, types.Hash{})
	buf := &bytes.Buffer{}
	assert.Nil(t, b.Encode(NewGobBlockEncoder(buf)))
	bDecode := new(Block)
	assert.Nil(t, bDecode.Decode(NewGobBlockDecoder(buf)))
	assert.Equal(t, b, bDecode)
}
func RandomBlock(t *testing.T, height uint32, preBlockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	tx := randomTxWithSignature(t)
	header := &Header{
		Version:       1,
		PrevBlockHash: preBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}
	b := NewBlock(header, []*Transaction{&tx})
	dataHash, err := CalculateDataHash(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash
	assert.Nil(t, b.Sign(privKey))
	//tx := Transaction{}
	return b
}

//func TestHeader_Encode_Decode(t *testing.T) {
//	h := &Header{
//		Version:   1,
//		PrevBlock: types.RandomHash(),
//		Timestamp: uint64(time.Now().UnixNano()),
//		Height:    10,
//		Nonce:     989394,
//	}
//	buf := &bytes.Buffer{}
//	assert.Nil(t, h.EncodeBinary(buf))
//	hDecode := &Header{}
//	assert.Nil(t, hDecode.DecodeBinary(buf))
//	assert.Equal(t, h, hDecode)
//}
//func TestBlock_Encode_Decode(t *testing.T) {
//	b := &Block{
//		Header: Header{
//			Version:   1,
//			PrevBlock: types.RandomHash(),
//			Timestamp: uint64(time.Now().UnixNano()),
//			Height:    10,
//			Nonce:     989394,
//		},
//		Transactions: nil,
//	}
//	buf := &bytes.Buffer{}
//	assert.Nil(t, b.EncodeBinary(buf))
//	bDecode := &Block{}
//	assert.Nil(t, bDecode.DecodeBinary(buf))
//	assert.Equal(t, b, bDecode)
//	fmt.Printf("%+v\n", b)
//}
//func TestBlockHash(t *testing.T) {
//	b := Block{
//		Header: Header{
//			Version:   1,
//			PrevBlock: types.RandomHash(),
//			Timestamp: uint64(time.Now().UnixNano()),
//			Height:    10,
//		},
//		Transactions: []Transaction{},
//	}
//	h := b.Hash()
//	fmt.Println(h)
//	assert.False(t, h.IsZero())
//}

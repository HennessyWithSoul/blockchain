package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"goblockchain/crypto"
	"testing"
)

func TestSignTransaction(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	data := []byte("hello world")
	tx := &Transaction{
		Data: data,
	}
	assert.Nil(t, tx.Sign(privateKey))
	assert.NotNil(t, tx.Signature)
}
func TestVerifyTransaction(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	data := []byte("hello world")
	tx := &Transaction{
		Data: data,
	}
	assert.Nil(t, tx.Sign(privateKey))
	assert.NotNil(t, tx.Signature)
	assert.Nil(t, tx.Verify())
	otherPrivateKey := crypto.GeneratePrivateKey()
	tx.From = otherPrivateKey.PublicKey()
	assert.NotNil(t, tx.Verify())
}
func TestTxEncodeDecode(t *testing.T) {
	//tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}
	txx := NewTransaction([]byte("123"))
	privateKey := crypto.GeneratePrivateKey()
	txx.Sign(privateKey)
	assert.Nil(t, txx.Encode(NewGobTxEncoder(buf)))
	//assert.Nil(t, tx.Encode(NewGobTxEncoder(buf)))
	//txDecoded := new(Transaction)
	//assert.Nil(t, txDecoded.Decode(NewGobTxDecoder(buf)))
	//assert.Equal(t, tx, txDecoded)

}
func TestNativeTransferTransaction(t *testing.T) {
	from := crypto.GeneratePrivateKey()
	to := crypto.GeneratePrivateKey()
	tx := &Transaction{
		To:    to.PublicKey(),
		Value: 666,
	}
	assert.Nil(t, tx.Sign(from))
	hash := tx.Hash(TxHasher{})
	fmt.Println(hash)
}
func TestNFTTransaction(t *testing.T) {
	collectionTx := CollectionTx{
		Fee:      200,
		MetaData: []byte("the beginning of a new collection"),
	}
	privateKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Type:    TxTypeCollection,
		TxInner: collectionTx,
	}
	tx.Sign(privateKey)
	buf := new(bytes.Buffer)
	assert.Nil(t, gob.NewEncoder(buf).Encode(tx))
	txDecode := &Transaction{}
	assert.Nil(t, gob.NewDecoder(buf).Decode(txDecode))
	assert.Equal(t, tx, txDecode)
	fmt.Printf("tx:%+v\n", txDecode)
}
func randomTxWithSignature(t *testing.T) Transaction {
	privateKey := crypto.GeneratePrivateKey()
	tx := Transaction{
		Data: []byte("Hello World"),
	}
	assert.Nil(t, tx.Sign(privateKey))
	return tx
}

package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"goblockchain/crypto"
	"goblockchain/types"
	"time"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
	Nonce         uint64
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)
	return buf.Bytes()
	//	buf := &bytes.Buffer{}
	//	enc := gob.NewEncoder(buf)
	//	enc.Encode(b.Header)
	//	return buf.Bytes()
}

type Block struct {
	//Version string
	*Header
	Transactions []*Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature
	hash         types.Hash
}

func NewBlock(header *Header, transactions []*Transaction) *Block {
	return &Block{
		Header:       header,
		Transactions: transactions,
	}
}
func NewBlockFromPrevHeader(h *Header, txx []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(txx)
	fmt.Println(dataHash)
	if err != nil {
		return nil, err
	}
	header := &Header{
		Version:       1,
		Height:        h.Height + 1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(h),
		Timestamp:     time.Now().UnixNano(),
	}
	return NewBlock(header, txx), nil
}
func (b *Block) AddTransaction(transaction *Transaction) {
	b.Transactions = append(b.Transactions, transaction)
}
func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}
	return b.hash
}
func (b *Block) Sign(privateKey crypto.PrivateKey) error {
	sig, err := privateKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}
	b.Validator = privateKey.PublicKey()
	b.Signature = sig
	return nil
}
func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block signature is nil")
	}
	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block signature is invalid")
	}
	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}
	datahash, err := CalculateDataHash(b.Transactions)
	fmt.Println(b.Transactions)
	if err != nil {
		return err
	}
	if datahash != b.DataHash {

		return fmt.Errorf("data hash does not match; expected %x, got %x", datahash, b.DataHash)
	}
	//fmt.Println(b.Signature.Verify(b.Validator, b.HeaderData()))
	return nil
}
func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}
func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}
func CalculateDataHash(txx []*Transaction) (types.Hash, error) {

	buf := &bytes.Buffer{}

	for _, tx := range txx {
		// 使用二进制编码器替代 gob 编码器
		encoder := NewBinaryTxEncoder(buf)
		if err := tx.Encode(encoder); err != nil {
			return types.Hash{}, err
		}
	}

	hash := sha256.Sum256(buf.Bytes())
	return hash, nil

}

func GeneisisBlock() *Block {
	header := &Header{
		Version:   1,
		DataHash:  types.Hash{},
		Height:    0,
		Timestamp: 00000,
	}
	return NewBlock(header, nil)
}

package core

import (
	"encoding/gob"
	"fmt"
	"goblockchain/crypto"
	"goblockchain/types"
	"io"
)

type TxType byte

const (
	TxTypeCollection TxType = iota
	TxTypeMint
)

type CollectionTx struct {
	Fee      int64
	MetaData []byte
}
type MintTx struct {
	Fee             int64
	NFT             types.Hash
	Collection      types.Hash
	MetaData        []byte
	CollectionOwner crypto.PublicKey
	Signature       crypto.Signature
}
type Transaction struct {
	Type      TxType
	TxInner   any //interface{}
	To        crypto.PublicKey
	Value     uint64
	Data      []byte
	From      crypto.PublicKey
	Signature *crypto.Signature
	hash      types.Hash
	//第一次被本地看到时的时间戳
	firstseen int64
}

func (tx *Transaction) DecodeBinary(r io.Reader) error {
	return nil
}
func (tx *Transaction) EncodeBinary(w io.Writer) error {
	return nil
}
func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}
	return tx.hash
}
func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}
func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(tx.Data)
	if err != nil {
		return err
	}
	tx.From = privKey.PublicKey()
	tx.Signature = sig
	return nil
}
func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transcation has no signature")
	}
	if tx.Signature.Verify(tx.From, tx.Data) {
		return nil
	} else {
		return fmt.Errorf("transcation has invalid signature")
	}
}
func (tx *Transaction) SetFirstSeen(t int64) {
	tx.firstseen = t
}
func (tx *Transaction) FirstSeen() int64 {
	return tx.firstseen
}
func init() {
	gob.Register(CollectionTx{})
	gob.Register(MintTx{})
}

// Encode 使用指定的编码器编码交易
func (tx *Transaction) Encode(encoder TxEncoder) error {
	return encoder.Encode(tx)
}

// Decode 使用指定的解码器解码交易
func (tx *Transaction) Decode(decoder TxDecoder) error {
	decodedTx, err := decoder.Decode()
	if err != nil {
		return err
	}
	*tx = *decodedTx
	return nil
}

// 为了向后兼容，保留 GobTxEncoder（但重写实现）

// safeTransaction 用于 gob 序列化的安全交易结构
type safeTransaction struct {
	Type       TxType
	To         []byte
	Value      uint64
	Data       []byte
	From       []byte
	SignatureR []byte
	SignatureS []byte
	FirstSeen  int64
}

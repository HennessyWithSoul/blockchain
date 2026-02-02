package core

import (
	"crypto/elliptic"
	"encoding/gob"
	"io"
)

type Encoder[T any] interface {
	Encode(T) error
}
type Decoder[T any] interface {
	Decode(T) error
}
type GobTxEncoder struct {
	w io.Writer
}

func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	gob.Register(elliptic.P256())
	return &GobTxEncoder{
		w: w,
	}
}

func (e *GobTxEncoder) Encode(tx *Transaction) error {
	// 创建一个可序列化的交易副本
	safeTx := &safeTransaction{
		Type:      tx.Type,
		Value:     tx.Value,
		Data:      tx.Data,
		FirstSeen: tx.firstseen,
	}

	// 只序列化公钥的字节，不序列化完整的公钥结构
	if tx.To.Key != nil {
		safeTx.To = tx.To.Bytes()
	}
	if tx.From.Key != nil {
		safeTx.From = tx.From.Bytes()
	}
	if tx.Signature != nil {
		safeTx.SignatureR = tx.Signature.R.Bytes()
		safeTx.SignatureS = tx.Signature.S.Bytes()
	}

	return gob.NewEncoder(e.w).Encode(safeTx)
}

type GobTxDecoder struct {
	r io.Reader
}

func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	gob.Register(elliptic.P256())
	return &GobTxDecoder{
		r: r,
	}
}
func (e *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(e.r).Decode(tx)
}

type GobBlockEncoder struct {
	w io.Writer
}

func NewGobBlockEncoder(w io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder{
		w: w,
	}
}
func (e *GobBlockEncoder) Encode(b *Block) error {
	return gob.NewEncoder(e.w).Encode(b)
}

type GobBlockDecoder struct {
	r io.Reader
}

func NewGobBlockDecoder(r io.Reader) *GobBlockDecoder {
	//gob.Register(elliptic.P256())
	return &GobBlockDecoder{
		r: r,
	}
}
func (e *GobBlockDecoder) Decode(b *Block) error {
	return gob.NewDecoder(e.r).Decode(b)
}
func init() {
	gob.Register(elliptic.P256())
}

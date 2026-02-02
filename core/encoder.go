package core

import (
	"encoding/binary"
	"goblockchain/crypto"
	"io"
	"math/big"
)

// TxEncoder 交易编码器接口
type TxEncoder interface {
	Encode(*Transaction) error
}

// TxDecoder 交易解码器接口
type TxDecoder interface {
	Decode() (*Transaction, error)
}

// BinaryTxEncoder 二进制交易编码器
type BinaryTxEncoder struct {
	w io.Writer
}

func NewBinaryTxEncoder(w io.Writer) *BinaryTxEncoder {
	return &BinaryTxEncoder{w: w}
}

func (e *BinaryTxEncoder) Encode(tx *Transaction) error {
	// 编码交易类型
	if err := binary.Write(e.w, binary.BigEndian, tx.Type); err != nil {
		return err
	}

	// 编码 To 公钥
	if err := e.encodePublicKey(tx.To); err != nil {
		return err
	}

	// 编码 Value
	if err := binary.Write(e.w, binary.BigEndian, tx.Value); err != nil {
		return err
	}

	// 编码 Data 长度和数据
	if err := binary.Write(e.w, binary.BigEndian, uint32(len(tx.Data))); err != nil {
		return err
	}
	if len(tx.Data) > 0 {
		if _, err := e.w.Write(tx.Data); err != nil {
			return err
		}
	}

	// 编码 From 公钥
	if err := e.encodePublicKey(tx.From); err != nil {
		return err
	}

	// 编码签名
	if err := e.encodeSignature(tx.Signature); err != nil {
		return err
	}

	// 编码 firstseen
	if err := binary.Write(e.w, binary.BigEndian, tx.firstseen); err != nil {
		return err
	}

	return nil
}

func (e *BinaryTxEncoder) encodePublicKey(pk crypto.PublicKey) error {
	pubBytes := pk.Bytes()
	if pubBytes == nil {
		// 写入空标记
		return binary.Write(e.w, binary.BigEndian, uint32(0))
	}

	// 写入公钥长度和数据
	if err := binary.Write(e.w, binary.BigEndian, uint32(len(pubBytes))); err != nil {
		return err
	}
	_, err := e.w.Write(pubBytes)
	return err
}

func (e *BinaryTxEncoder) encodeSignature(sig *crypto.Signature) error {
	if sig == nil {
		// 写入空标记
		return binary.Write(e.w, binary.BigEndian, uint32(0))
	}

	// 写入 R 和 S
	rBytes := sig.R.Bytes()
	sBytes := sig.S.Bytes()

	// 写入 R 长度和数据
	if err := binary.Write(e.w, binary.BigEndian, uint32(len(rBytes))); err != nil {
		return err
	}
	if _, err := e.w.Write(rBytes); err != nil {
		return err
	}

	// 写入 S 长度和数据
	if err := binary.Write(e.w, binary.BigEndian, uint32(len(sBytes))); err != nil {
		return err
	}
	_, err := e.w.Write(sBytes)
	return err
}

// BinaryTxDecoder 二进制交易解码器
type BinaryTxDecoder struct {
	r io.Reader
}

func NewBinaryTxDecoder(r io.Reader) *BinaryTxDecoder {
	return &BinaryTxDecoder{r: r}
}

func (d *BinaryTxDecoder) Decode() (*Transaction, error) {
	tx := &Transaction{}

	// 解码交易类型
	if err := binary.Read(d.r, binary.BigEndian, &tx.Type); err != nil {
		return nil, err
	}

	// 解码 To 公钥
	toPub, err := d.decodePublicKey()
	if err != nil {
		return nil, err
	}
	tx.To = toPub

	// 解码 Value
	if err := binary.Read(d.r, binary.BigEndian, &tx.Value); err != nil {
		return nil, err
	}

	// 解码 Data
	var dataLen uint32
	if err := binary.Read(d.r, binary.BigEndian, &dataLen); err != nil {
		return nil, err
	}
	if dataLen > 0 {
		tx.Data = make([]byte, dataLen)
		if _, err := io.ReadFull(d.r, tx.Data); err != nil {
			return nil, err
		}
	}

	// 解码 From 公钥
	fromPub, err := d.decodePublicKey()
	if err != nil {
		return nil, err
	}
	tx.From = fromPub

	// 解码签名
	sig, err := d.decodeSignature()
	if err != nil {
		return nil, err
	}
	tx.Signature = sig

	// 解码 firstseen
	if err := binary.Read(d.r, binary.BigEndian, &tx.firstseen); err != nil {
		return nil, err
	}

	return tx, nil
}

func (d *BinaryTxDecoder) decodePublicKey() (crypto.PublicKey, error) {
	var pubLen uint32
	if err := binary.Read(d.r, binary.BigEndian, &pubLen); err != nil {
		return crypto.PublicKey{}, err
	}

	if pubLen == 0 {
		return crypto.PublicKey{}, nil
	}

	pubBytes := make([]byte, pubLen)
	if _, err := io.ReadFull(d.r, pubBytes); err != nil {
		return crypto.PublicKey{}, err
	}

	return crypto.PublicKeyFromBytes(pubBytes, nil)
}

func (d *BinaryTxDecoder) decodeSignature() (*crypto.Signature, error) {
	var rLen uint32
	if err := binary.Read(d.r, binary.BigEndian, &rLen); err != nil {
		return nil, err
	}

	if rLen == 0 {
		return nil, nil
	}

	rBytes := make([]byte, rLen)
	if _, err := io.ReadFull(d.r, rBytes); err != nil {
		return nil, err
	}

	var sLen uint32
	if err := binary.Read(d.r, binary.BigEndian, &sLen); err != nil {
		return nil, err
	}

	sBytes := make([]byte, sLen)
	if _, err := io.ReadFull(d.r, sBytes); err != nil {
		return nil, err
	}

	return &crypto.Signature{
		R: newIntFromBytes(rBytes),
		S: newIntFromBytes(sBytes),
	}, nil
}

func newIntFromBytes(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}

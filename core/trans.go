package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"goblockchain/crypto"
	"goblockchain/rpc"
	"goblockchain/types"
	"math/big"

	"google.golang.org/protobuf/types/known/anypb"
)

// ==================== PublicKey 转换 ====================
// ToGRPCPublicKey 将 string 格式的 ECDSA 公钥转换为  ECDSAPublicKey

func stringToECDSAPublicKey(publicKeyStr string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not an ECDSA public key")
	}

	return ecdsaPub, nil
}

func PublicKeyToBytes(publicKey crypto.PublicKey) []byte {
	// 将公钥编码为 DER 格式
	keyBytes, err := x509.MarshalPKIXPublicKey(publicKey.Key)
	if err != nil {
		panic(err)
	}

	// 创建 PEM 块
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	}

	// 编码为 PEM 字节数组
	return pem.EncodeToMemory(pemBlock)
}

// ToGRPCPublicKey 将 crypto.PublicKey 转换为 gRPC ECDSAPublicKey
func ToGRPCPublicKey(pk crypto.PublicKey) *rpc.SerializedPublicKey {
	if pk.Key == nil {
		return nil
	}

	return &rpc.SerializedPublicKey{
		CompressedData: pk.Bytes(), // 使用压缩格式的字节
		Curve:          "P-256",    // 假设使用 P-256 曲线
	}
}

// FromGRPCPublicKey 将 gRPC ECDSAPublicKey 转换为 crypto.PublicKey
func FromGRPCPublicKey(grpcPk *rpc.SerializedPublicKey) crypto.PublicKey {
	if grpcPk == nil || len(grpcPk.CompressedData) == 0 {
		return crypto.PublicKey{}
	}

	curve := getCurveFromName(grpcPk.Curve)
	pk, err := crypto.PublicKeyFromBytes(grpcPk.CompressedData, curve)
	if err != nil {
		// 记录错误，但返回空公钥
		return crypto.PublicKey{}
	}

	return pk
}

// ==================== Signature 转换 ====================

// ToGRPCSignature 将 crypto.Signature 转换为 gRPC ECDSASignature
func ToGRPCSignature(sig *crypto.Signature) *rpc.SerializedSignature {
	if sig == nil {
		return nil
	}

	return &rpc.SerializedSignature{
		R: sig.R.Bytes(), // 大端字节表示
		S: sig.S.Bytes(),
	}
}

// FromGRPCSignature 将 gRPC ECDSASignature 转换为 crypto.Signature
func FromGRPCSignature(grpcSig *rpc.SerializedSignature) *crypto.Signature {
	if grpcSig == nil {
		return nil
	}

	return &crypto.Signature{
		R: new(big.Int).SetBytes(grpcSig.R),
		S: new(big.Int).SetBytes(grpcSig.S),
	}
}

// ==================== Hash 转换 ====================

// ToGRPCHash 将 types.Hash 转换为 gRPC Hash
func ToGRPCHash(h types.Hash) *rpc.Hash {
	return &rpc.Hash{
		Data: h[:],
	}
}

// FromGRPCHash 将 gRPC Hash 转换为 types.Hash
func FromGRPCHash(grpcHash *rpc.Hash) types.Hash {
	var h types.Hash
	if grpcHash != nil && len(grpcHash.Data) == len(h) {
		copy(h[:], grpcHash.Data)
	}
	return h
}

// ==================== Transaction 转换 ====================

// ToGRPCTransaction 将 core.Transaction 转换为 gRPC Transaction
func ToGRPCTransaction(tx *Transaction) (*rpc.Transaction, error) {
	if tx == nil {
		return nil, nil
	}

	// 注意：这里不再调用 tx.Hash()，避免序列化问题
	grpcTx := &rpc.Transaction{
		To:        ToGRPCPublicKey(tx.To),
		Value:     tx.Value,
		Data:      tx.Data,
		From:      ToGRPCPublicKey(tx.From),
		Signature: ToGRPCSignature(tx.Signature),
		FirstSeen: tx.firstseen,
	}

	// 手动计算哈希并设置

	// 处理 TxInner
	if tx.TxInner != nil {
		anyInner, err := encodeTxInnerToAny(tx.TxInner)
		if err != nil {
			return nil, fmt.Errorf("failed to encode TxInner: %w", err)
		}
		grpcTx.TxInner = anyInner
	}

	return grpcTx, nil
}

// FromGRPCTransaction 将 gRPC Transaction 转换为 core.Transaction
func FromGRPCTransaction(grpcTx *rpc.Transaction) (*Transaction, error) {
	if grpcTx == nil {
		return nil, nil
	}

	tx := &Transaction{
		Type:      TxType(grpcTx.Type),
		To:        FromGRPCPublicKey(grpcTx.To),
		Value:     grpcTx.Value,
		Data:      grpcTx.Data,
		From:      FromGRPCPublicKey(grpcTx.From),
		Signature: FromGRPCSignature(grpcTx.Signature),
		hash:      FromGRPCHash(grpcTx.Hash),
		firstseen: grpcTx.FirstSeen,
	}

	// 处理 TxInner
	if grpcTx.TxInner != nil {
		inner, err := decodeAnyToTxInner(grpcTx.TxInner)
		if err != nil {
			return nil, fmt.Errorf("failed to decode TxInner: %w", err)
		}
		tx.TxInner = inner
	}

	return tx, nil
}

// ToGRPCBlockHeader 将 Header 转换为 gRPC BlockHeader
func ToGRPCBlockHeader(header *Header) *rpc.BlockHeader {
	if header == nil {
		return nil
	}

	return &rpc.BlockHeader{
		Version:       header.Version,
		DataHash:      ToGRPCHash(header.DataHash),
		PrevBlockHash: ToGRPCHash(header.PrevBlockHash),
		Timestamp:     header.Timestamp,
		Height:        header.Height,
		Nonce:         header.Nonce,
	}
}

// FromGRPCBlockHeader 将 gRPC BlockHeader 转换为 Header
func FromGRPCBlockHeader(grpcHeader *rpc.BlockHeader) *Header {
	if grpcHeader == nil {
		return nil
	}

	return &Header{
		Version:       grpcHeader.Version,
		DataHash:      FromGRPCHash(grpcHeader.DataHash),
		PrevBlockHash: FromGRPCHash(grpcHeader.PrevBlockHash),
		Timestamp:     grpcHeader.Timestamp,
		Height:        grpcHeader.Height,
		Nonce:         grpcHeader.Nonce,
	}
}

// ToGRPCBlock 将 Block 转换为 gRPC Block
func ToGRPCBlock(block *Block) (*rpc.Block, error) {
	if block == nil {
		return nil, nil
	}

	grpcBlock := &rpc.Block{
		Header:    ToGRPCBlockHeader(block.Header),
		Validator: ToGRPCPublicKey(block.Validator),
		Signature: ToGRPCSignature(block.Signature),
		Hash:      ToGRPCHash(block.hash),
	}
	for _, tx := range block.Transactions {
		grpcTx, err := ToGRPCTransaction(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to convert transaction: %w", err)
		}
		grpcBlock.Transactions = append(grpcBlock.Transactions, grpcTx)
	}

	return grpcBlock, nil
}

// FromGRPCBlock 将 gRPC Block 转换为 Block
func FromGRPCBlock(grpcBlock *rpc.Block) (*Block, error) {
	if grpcBlock == nil {
		return nil, nil
	}

	block := &Block{
		Header:    FromGRPCBlockHeader(grpcBlock.Header),
		Validator: FromGRPCPublicKey(grpcBlock.Validator),
		Signature: FromGRPCSignature(grpcBlock.Signature),
		hash:      FromGRPCHash(grpcBlock.Hash),
	}

	// 转换交易列表
	for _, grpcTx := range grpcBlock.Transactions {
		tx, err := FromGRPCTransaction(grpcTx)
		if err != nil {
			return nil, fmt.Errorf("failed to convert transaction: %w", err)
		}
		block.Transactions = append(block.Transactions, tx)
	}

	return block, nil
}

// ==================== TxInner 任意类型处理 ====================

// encodeTxInnerToAny 将任意类型的 TxInner 编码为 Any
func encodeTxInnerToAny(inner interface{}) (*anypb.Any, error) {
	switch v := inner.(type) {
	case *rpc.TransferTxInner:
		return anypb.New(v)
	case *rpc.CallTxInner:
		return anypb.New(v)
	case *rpc.DeployTxInner:
		return anypb.New(v)
	case map[string]interface{}:
		// 对于 map 类型，序列化为 JSON 并包装为 CallTxInner
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return anypb.New(&rpc.CallTxInner{
			Arguments: jsonData,
		})
	default:
		// 通用处理：序列化为 JSON
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return &anypb.Any{
			TypeUrl: "internal/" + fmt.Sprintf("%T", v),
			Value:   jsonData,
		}, nil
	}
}

// decodeAnyToTxInner 将 Any 解码为具体的 TxInner 类型
func decodeAnyToTxInner(anyInner *anypb.Any) (interface{}, error) {
	// 尝试已知类型
	switch anyInner.TypeUrl {
	case "type.googleapis.com/blockchain.TransferTxInner":
		var transferTx rpc.TransferTxInner
		if err := anyInner.UnmarshalTo(&transferTx); err != nil {
			return nil, err
		}
		return &transferTx, nil

	case "type.googleapis.com/blockchain.CallTxInner":
		var callTx rpc.CallTxInner
		if err := anyInner.UnmarshalTo(&callTx); err != nil {
			return nil, err
		}
		return &callTx, nil

	case "type.googleapis.com/blockchain.DeployTxInner":
		var deployTx rpc.DeployTxInner
		if err := anyInner.UnmarshalTo(&deployTx); err != nil {
			return nil, err
		}
		return &deployTx, nil

	default:
		// 内部类型，使用 JSON 反序列化
		if len(anyInner.TypeUrl) > 0 && anyInner.TypeUrl[:9] == "internal/" {
			var result map[string]interface{}
			if err := json.Unmarshal(anyInner.Value, &result); err != nil {
				return nil, err
			}
			return result, nil
		}

		return nil, errors.New("unknown TxInner type: " + anyInner.TypeUrl)
	}
}

// ==================== 辅助函数 ====================

// getCurveName 获取曲线名称

// 辅助函数
func uint64ToBytes(n uint64) []byte {
	bytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		bytes[7-i] = byte(n >> (i * 8))
	}
	return bytes
}

func int64ToBytes(n int64) []byte {
	return uint64ToBytes(uint64(n))
}

// CreateTransferTxInner 创建转账交易内部数据
func CreateTransferTxInner(memo string, gasLimit, gasPrice uint64) *rpc.TransferTxInner {
	return &rpc.TransferTxInner{
		Memo:     memo,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
	}
}

// CreateCallTxInner 创建合约调用内部数据
func CreateCallTxInner(contractAddress []byte, functionSelector []byte, args []byte, gasLimit, gasPrice uint64) *rpc.CallTxInner {
	return &rpc.CallTxInner{
		ContractAddress:  contractAddress,
		FunctionSelector: functionSelector,
		Arguments:        args,
		GasLimit:         gasLimit,
		GasPrice:         gasPrice,
	}
}

func getCurveFromName(curveName string) elliptic.Curve {
	switch curveName {
	case "P256":
		return elliptic.P256()
	case "P384":
		return elliptic.P384()
	case "P521":
		return elliptic.P521()
	default:
		return nil
	}
}

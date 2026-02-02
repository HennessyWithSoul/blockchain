package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"goblockchain/types"
	"math/big"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

type SerializablePublicKey struct {
	X     *big.Int
	Y     *big.Int
	Curve string // 曲线名称: "P256", "P384", "P521"
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}
	return &Signature{
		r,
		s,
	}, nil
}
func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return PrivateKey{
		key: key,
	}

}
func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: &k.key.PublicKey,
	}
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y)
}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())
	return types.AddressFrombytes(h[len(h)-20:])
}

type Signature struct {
	R, S *big.Int
}

func (sig Signature) Verify(publicKey PublicKey, data []byte) bool {
	//fmt.Println(publicKey, len(data), sig.r, sig.s, ecdsa.Verify(publicKey.key, data, sig.r, sig.s))
	//fmt.Printf("SHA-256: %x\n", sha256.Sum256(data))
	//data1 := []byte(fmt.Sprintf("123"))
	//data2 := []byte(fmt.Sprintf("234"))
	//fmt.Println(ecdsa.Verify(publicKey.key, data1, sig.r, sig.s), ecdsa.Verify(publicKey.key, data2, sig.r, sig.s), ecdsa.Verify(publicKey.key, data, sig.r, sig.s))
	return ecdsa.Verify(publicKey.Key, data, sig.R, sig.S)
}

func (k PublicKey) String() string {
	return hex.EncodeToString(k.ToSlice())
}

func StringToPublicKey(keyStr string) (PublicKey, error) {
	// 解码 PEM 块
	block, _ := pem.Decode([]byte(keyStr))
	if block == nil {
		return PublicKey{}, errors.New("failed to parse PEM block containing the public key")
	}

	if block.Type != "PUBLIC KEY" {
		return PublicKey{}, errors.New("invalid PEM block type, expected 'PUBLIC KEY'")
	}

	// 解析公钥
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return PublicKey{}, err
	}

	// 类型断言为 ECDSA 公钥
	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return PublicKey{}, errors.New("not an ECDSA public key")
	}

	return PublicKey{Key: ecdsaPubKey}, nil
}

func (pk PublicKey) ToSerializable() *SerializablePublicKey {
	if pk.Key == nil {
		return nil
	}

	return &SerializablePublicKey{
		X:     pk.Key.X,
		Y:     pk.Key.Y,
		Curve: getCurveName(pk.Key.Curve),
	}
}

//func FromSerializable(spk *SerializablePublicKey) PublicKey {
//	if spk == nil {
//		return PublicKey{}
//	}
//
//	curve := getCurveFromName(spk.Curve)
//	if curve == nil {
//		return PublicKey{}
//	}
//
//	return PublicKey{
//		Key: &ecdsa.PublicKey{
//			Curve: curve,
//			X:     spk.X,
//			Y:     spk.Y,
//		},
//	}
//}

// SerializableSignature 可序列化的签名
type SerializableSignature struct {
	R *big.Int
	S *big.Int
}

// ToSerializable 将 Signature 转换为可序列化的格式
func (sig *Signature) ToSerializable() *SerializableSignature {
	if sig == nil {
		return nil
	}

	return &SerializableSignature{
		R: sig.R,
		S: sig.S,
	}
}

// FromSerializable 从可序列化格式恢复 Signature
//func FromSerializable(ss *SerializableSignature) *Signature {
//	if ss == nil {
//		return nil
//	}
//
//	return &Signature{
//		R: ss.R,
//		S: ss.S,
//	}
//}

// 辅助函数
func getCurveName(curve elliptic.Curve) string {
	switch curve {
	case elliptic.P256():
		return "P256"
	case elliptic.P384():
		return "P384"
	case elliptic.P521():
		return "P521"
	default:
		return "Unknown"
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

// Bytes 返回公钥的压缩格式字节
func (pk PublicKey) Bytes() []byte {
	if pk.Key == nil {
		return nil
	}
	return elliptic.MarshalCompressed(pk.Key.Curve, pk.Key.X, pk.Key.Y)
}

// FromBytes 从字节恢复公钥
func PublicKeyFromBytes(data []byte, curve elliptic.Curve) (PublicKey, error) {
	if curve == nil {
		curve = elliptic.P256() // 默认曲线
	}

	x, y := elliptic.UnmarshalCompressed(curve, data)
	if x == nil {
		return PublicKey{}, errors.New("failed to unmarshal public key")
	}

	return PublicKey{
		Key: &ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
	}, nil
}

func PrivateKeyFromHex(hexKey string) (*PrivateKey, error) {
	if hexKey == "" {
		return nil, errors.New("private key string is empty")
	}

	// 解码十六进制字符串
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
	}

	return PrivateKeyFromBytes(keyBytes)
}

// PrivateKeyFromBytes 从字节创建私钥
func PrivateKeyFromBytes(keyBytes []byte) (*PrivateKey, error) {
	if len(keyBytes) != 32 {
		return nil, errors.New("invalid private key length")
	}

	// 创建大整数
	d := new(big.Int).SetBytes(keyBytes)

	// 创建 ECDSA 私钥
	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
		},
		D: d,
	}

	// 计算公钥坐标
	privateKey.PublicKey.X, privateKey.PublicKey.Y = elliptic.P256().ScalarBaseMult(d.Bytes())

	return &PrivateKey{key: privateKey}, nil
}

// String 返回公钥的十六进制字符串表示

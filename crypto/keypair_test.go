package crypto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePrivateKey(t *testing.T) {
	privatekey := GeneratePrivateKey()
	fmt.Println(privatekey.key)
	publicKey := privatekey.PublicKey()
	address := publicKey.Address()
	fmt.Println(address.String())
	msg := []byte("hello world")
	sig, err := privatekey.Sign(msg)
	assert.Nil(t, err)
	fmt.Println(sig)
	assert.True(t, sig.Verify(publicKey, msg))
}
func TestKeypairSignVerifySuccess(t *testing.T) {
	privatekey := GeneratePrivateKey()
	publicKey := privatekey.PublicKey()
	msg := []byte("hello world")
	sig, err := privatekey.Sign(msg)
	assert.Nil(t, err)
	assert.True(t, sig.Verify(publicKey, msg))
}

func TestKeypairSignVerifyFail(t *testing.T) {
	privatekey := GeneratePrivateKey()
	publicKey := privatekey.PublicKey()
	msg := []byte("hello world")
	sig, err := privatekey.Sign(msg)
	assert.Nil(t, err)
	otherPrivateKey := GeneratePrivateKey()
	otherPublicKey := otherPrivateKey.PublicKey()

	assert.False(t, sig.Verify(otherPublicKey, msg))
	assert.False(t, sig.Verify(publicKey, []byte("xns")))
}

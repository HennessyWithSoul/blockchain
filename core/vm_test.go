package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"goblockchain/crypto"
	"goblockchain/util"
	"testing"
)

func TestVM(t *testing.T) {
	data := []byte{0x03, 0x0a, 0x02, 0x0a, 0x0e}
	vm := NewVM(data, NewState())
	vm.Run()
	//fmt.Println(vm.stack)
	assert.Equal(t, 1, vm.stack.Pop().(int))
}
func TestStackByte(t *testing.T) {
	data := []byte{0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x03, 0x0a, 0x0d}
	vm := NewVM(data, NewState())
	assert.Nil(t, vm.Run())
	result := vm.stack.Pop().([]byte)
	fmt.Println(string(result))

}
func TestVMStore(t *testing.T) {
	//data := []byte{0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x03, 0x0a, 0x0d, 0x03, 0x0a, 0x02, 0x0e}
	data := []byte{0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x03, 0x0a, 0x0d, 0x03, 0x0a, 0x02, 0x0a, 0x0e, 0x0f}
	dataother := []byte{0x46, 0x0c, 0x4f, 0x0c, 0x4e, 0x0c, 0x03, 0x0a, 0x0d, 0x03, 0x0a, 0x02, 0x0a, 0x0e, 0x0f}
	data = append(data, dataother...)
	vm := NewVM(data, NewState())
	assert.Nil(t, vm.Run())

	//fmt.Println(vm.stack.data)
	//fmt.Printf("%+v\n", vm.contractState)
	keyBytes, _ := vm.contractState.Get([]byte("FOO"))
	key := util.DeserializeInt64(keyBytes)
	assert.Equal(t, key, int64(1))
}
func TestVMGet(t *testing.T) {
	//data := []byte{0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x03, 0x0a, 0x0d, 0x03, 0x0a, 0x02, 0x0e}
	data := []byte{0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x03, 0x0a, 0x0d, 0x03, 0x0a, 0x02, 0x0a, 0x0e, 0x0f,
		0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x03, 0x0a, 0x0d, 0xae,
	}
	vm := NewVM(data, NewState())
	assert.Nil(t, vm.Run())
	fmt.Printf("%+v\n%d\n", vm.stack, vm.stack.sp)
	value := vm.stack.Pop().([]byte)
	//fmt.Printf("%d\n", value.(uint8))
	//Printf("%+v\n", vm.stack)
	//assert.Nil(t, vm.stack.Pop())
	assert.Equal(t, util.DeserializeInt64(value), int64(1))
	//fmt.Println(vm.stack.data)
	//fmt.Printf("%+v\n", vm.contractState)
	//keyBytes, _ := vm.contractState.Get([]byte("FOO"))
	//key := util.DeserializeInt64(keyBytes)
	//assert.Equal(t, key, int64(1))
}

func TestKeyTrans(t *testing.T) {
	pri := crypto.GeneratePrivateKey()
	pub := pri.PublicKey()
	fmt.Println("pub", pub)
	fmt.Printf("%x\n", PublicKeyToBytes(pub))
	fmt.Println(pub.String())
	_, err := stringToECDSAPublicKey(pub.String())
	fmt.Println(err)
}

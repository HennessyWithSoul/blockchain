package network

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	tra.Connect(trb)
	trb.Connect(tra)
	assert.Equal(t, tra.Peers[trb.Addr()], trb)
	assert.Equal(t, trb.Peers[tra.Addr()], tra)
}
func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	tra.Connect(trb)
	trb.Connect(tra)
	msg := []byte("Hello World")
	assert.Nil(t, tra.SendMessage(trb.addr, msg))
	rpc := <-trb.Consume()
	buf := make([]byte, len(msg))
	n, err := rpc.Payload.Read(buf)
	//assert(rpc.Payload.Read(buf),)
	assert.Equal(t, n, len(msg))
	assert.Nil(t, err)
	assert.Equal(t, buf, msg)
	assert.Equal(t, rpc.From, tra.addr)

}
func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")
	tra.Connect(trb)
	tra.Connect(trc)
	msg := []byte("Hello World")
	assert.Nil(t, tra.Boardcast(msg))
	rpcb := <-trb.Consume()
	b, err := ioutil.ReadAll(rpcb.Payload)
	assert.Nil(t, err)
	assert.Equal(t, msg, b)
	rpcc := <-trc.Consume()
	c, err := ioutil.ReadAll(rpcc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, msg, c)

}

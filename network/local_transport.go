package network

import (
	"bytes"
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr         NetAddr
	lock         sync.RWMutex
	consumeCh    chan RPC
	Peers        map[NetAddr]*LocalTransport
	tcpAddr      TCPAddr
	TCPTransport *TCPTransport
}

func NewLocalTransport(addr NetAddr, tcpAddr TCPAddr) *LocalTransport {
	ch := make(chan RPC, 1024)
	return &LocalTransport{
		addr:         addr,
		tcpAddr:      tcpAddr,
		consumeCh:    ch,
		Peers:        make(map[NetAddr]*LocalTransport),
		TCPTransport: NewTCPTransport(string(tcpAddr), ch),
	}
}
func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}
func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.Peers[tr.Addr()] = tr.(*LocalTransport)
	return nil
}
func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}
func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	//fmt.Println("1111")
	defer t.lock.RUnlock()
	//fmt.Printf(" %s Sending message to %s\n", t.addr, to)
	if t.Addr() == to {
		return nil
	}
	peer, ok := t.Peers[to]
	if !ok {
		fmt.Errorf("could not send message to %s", t.addr)
	}
	//fmt.Printf("%s peer is%s  \n", t.Addr(), to, peer.Addr())
	peer.consumeCh <- RPC{
		From:    t.Addr(),
		Payload: bytes.NewReader(payload),
	}
	return nil
}

func (t *LocalTransport) Boardcast(payload []byte) error {
	for _, peer := range t.Peers {
		if err := t.SendMessage(peer.addr, payload); err != nil {
			return err
		}
	}
	return nil
}
func (t *LocalTransport) TCPAddr() TCPAddr {
	return t.tcpAddr
}

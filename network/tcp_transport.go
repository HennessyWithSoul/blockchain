package network

import (
	"bytes"
	"fmt"
	"net"
)

type TCPTransport struct {
	ListenAddr string       //端口号
	listener   net.Listener //接收者
	tcpchan    chan RPC
}
type TCPPeer struct {
	conn net.Conn
}

func NewTCPTransport(addr string, ch chan RPC) *TCPTransport {
	return &TCPTransport{
		ListenAddr: addr,
		tcpchan:    ch,
	}
}
func (t *TCPTransport) readLoop(peer *TCPPeer) error {
	buf := make([]byte, 2048)

	n, err := peer.conn.Read(buf)
	if err != nil {
		fmt.Printf("%d read error:%s", n, err.Error())
		//continue
	}
	msg := buf[:n]
	//fmt.Println(string(msg))

	t.tcpchan <- RPC{
		From:    "",
		Payload: bytes.NewReader(NewMessage(MessageTypeBlock, msg).Bytes()),
	}
	//fmt.Printf("read ans send ok :%s\n", string(msg))

	return nil
}
func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			//fmt.Printf("accept error from %+v\n", conn)
			continue
		}
		peer := &TCPPeer{
			conn: conn,
		}
		go t.readLoop(peer)
		//fmt.Printf("new incoming TCP connection =>%+v\n", conn)

	}
}
func (t *TCPTransport) Start() error {
	ln, err := net.Listen("tcp", t.ListenAddr)
	if err != nil {
		fmt.Printf("listen tcp err:%+v\n", err)
		return err
	}
	t.listener = ln
	go t.acceptLoop()
	fmt.Println("TCP listening to port:", t.ListenAddr)
	return nil
}

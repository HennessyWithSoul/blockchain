package network

type NetAddr string
type TCPAddr string
type Transport interface {
	Connect(Transport) error
	Consume() <-chan RPC
	SendMessage(NetAddr, []byte) error
	Boardcast([]byte) error
	Addr() NetAddr
	TCPAddr() TCPAddr
}

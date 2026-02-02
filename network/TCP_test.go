package network

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestTCP(t *testing.T) {
	tr := NewTCPTransport(":12214")
	go tr.Start()

	time.Sleep(1 * time.Second)
	for i := 0; i < 10; i++ {
		go tcpTester()
	}

	time.Sleep(5 * time.Second)
}
func tcpTester() {
	conn, err := net.Dial("tcp", ":12214")
	if err != nil {
		fmt.Println(err)
	}
	_, err = conn.Write([]byte("hello world"))
	if err != nil {
		fmt.Println(err)
	}

}

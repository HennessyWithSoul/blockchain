package network

import (
	"bytes"
	"fmt"
	"goblockchain/core"
)

func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		go s.Transport.SendMessage(tr.Addr(), payload)
		fmt.Println(tr.Addr())

	}
	return nil
}
func (s *Server) broadcastBlock(b *core.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(core.NewGobBlockEncoder(buf)); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeBlock, buf.Bytes())
	return s.broadcast(msg.Bytes())

}
func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeTx, buf.Bytes())
	s.broadcast(msg.Bytes())

	return s.broadcast(msg.Bytes())
}

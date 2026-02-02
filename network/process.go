package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"goblockchain/core"
	"net"
	"time"
)

func (s *Server) ProcessGetStatusMessage(from NetAddr, data *GetStatusMessage) error {
	//fmt.Println("ProcessGetStatusMessage")
	statusMessage := &StatusMessage{
		CurrentHeight: s.chain.Height(),
		ID:            s.ID,
	}
	//fmt.Println(statusMessage)
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(statusMessage); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeStatus, buf.Bytes())
	//fmt.Println(s.Transport.Addr(), from)
	return s.Transport.SendMessage(from, msg.Bytes())
}
func (s *Server) ProcessTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})
	if s.memPool.Has(hash) {
		//logrus.WithFields(
		//	logrus.Fields{
		//		"hash": hash,
		//	}).Info("transaction already in mempool")
		return nil
	}
	if err := tx.Verify(); err != nil {
		return err
	}
	tx.SetFirstSeen(time.Now().UnixNano())
	//s.Logger.Log("msg", "adding new transaction", "hash", hash, "mempoollen", s.memPool.Len())
	go s.broadcastTx(tx)

	return s.memPool.Add(tx)
}
func (s *Server) ProcessBlock(b *core.Block) error {
	//fmt.Println(s.ID, "get block", b.DataHash, b.Height)
	if err := s.chain.AddBlock(b); err != nil {
		fmt.Println(err)
		return err
	}
	go s.broadcastBlock(b)
	return nil
}
func (s *Server) ProcessStatusMessage(from NetAddr, data *StatusMessage) error {
	//fmt.Println("received status message from %s ==> %+v\n", from, msg)
	if data.CurrentHeight <= s.chain.Height() {
		return nil
	}
	getblockmessage := GetBlockMessage{
		From: s.chain.Height() + 1,
		To:   data.CurrentHeight,
	}
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(getblockmessage); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())
	if err := s.Transport.SendMessage(from, msg.Bytes()); err != nil {
		return err
	}

	return nil
}
func (s *Server) ProcessGetBlocksMessage(from NetAddr, data *GetBlockMessage) error {
	address := s.Transport.(*LocalTransport).Peers[from].TCPAddr()
	conn, err := net.Dial("tcp", string(address))
	if err != nil {
		fmt.Println(err)
	}

	block, err := s.chain.GetBlock(data.From)
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	if err = block.Encode(core.NewGobBlockEncoder(buf)); err != nil {
		return err
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

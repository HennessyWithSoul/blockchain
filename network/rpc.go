package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"goblockchain/core"
	"io"
)

type MessageType byte

const (
	MessageTypeTx        MessageType = 0x1
	MessageTypeBlock     MessageType = 0x2
	MessageTypeGetBlocks MessageType = 0x3
	MessageTypeStatus    MessageType = 0x4
	MessageTypeGetStatus MessageType = 0x5
)

type RPC struct {
	From    NetAddr
	Payload io.Reader
}
type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}
func (msg *Message) Bytes() []byte {
	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

type DecodeMessage struct {
	From NetAddr
	Data any
}
type RPCDecodeFunc func(RPC) (*DecodeMessage, error)
type RPCHandler interface {
	HandleRPC(rpc RPC) error
}
type DefaultRPCHandler struct {
	p RPCProcessor
}

func NewDefaultRPCHandler(p RPCProcessor) *DefaultRPCHandler {
	return &DefaultRPCHandler{
		p: p,
	}
}
func DefaultRPCDecodeFunc(rpc RPC) (*DecodeMessage, error) {
	msg := &Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(msg); err != nil {
		return nil, fmt.Errorf("could not decode payload: %s", err)
	}
	//fmt.Println(msg.Header, "----------")
	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodeMessage{
			From: rpc.From,
			Data: tx,
		}, nil
	//h.p.ProcessTransaction(rpc.From, tx)
	case MessageTypeBlock:
		block := new(core.Block)
		if err := block.Decode(core.NewGobBlockDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodeMessage{
			From: rpc.From,
			Data: block,
		}, nil
	case MessageTypeGetStatus:
		return &DecodeMessage{
			From: rpc.From,
			Data: &GetStatusMessage{},
		}, nil
	case MessageTypeStatus:
		statusMessage := new(StatusMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(statusMessage); err != nil {
			return nil, err
		}
		return &DecodeMessage{
			From: rpc.From,
			Data: statusMessage,
		}, nil
	case MessageTypeGetBlocks:
		getblockmessage := new(GetBlockMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(getblockmessage); err != nil {
			return nil, err
		}
		return &DecodeMessage{
			From: rpc.From,
			Data: getblockmessage,
		}, nil
	default:
		return nil, fmt.Errorf("invalid message header %d", msg.Header)
	}
}

type RPCProcessor interface {
	ProcessMessage(*DecodeMessage) error
}

package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/go-kit/log"
	"goblockchain/core"
	"goblockchain/crypto"
	"goblockchain/types"
	"os"
	"time"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	ID            string
	Transport     Transport
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}
type Server struct {
	ServerOpts
	blocktime   time.Duration
	memPool     *TxPool
	chain       *core.BlockChain
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	chain, err := core.NewBlockChain(geneisisBlock())
	if err != nil {
		return nil, err
	}
	s := &Server{
		//ServerOpts:  opts,
		chain:       chain,
		blocktime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	if opts.RPCProcessor == nil {
		opts.RPCProcessor = s
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}
	s.ServerOpts = opts
	//

	return s, nil
}
func (s *Server) Start() {
	s.initTransports()
	//time.Sleep(1 * time.Second)
	if s.isValidator {
		go s.ValidatorLoop()
	}
	go func() {
		for {
			trans := s.Transport.(*LocalTransport)
			for _, peer := range trans.Peers {
				//fmt.Println(" --------", s.Transport.Addr(), peer.addr)
				f := 1
				for _, per := range s.Transports {
					if per == peer {
						f = 0
					}
				}
				if f == 1 {
					fmt.Printf("New peer %s\n", peer.Addr())
					s.Transports = append(s.Transports, peer)
				}
				//time.Sleep(1 * time.Second)
				//fmt.Println(peer.Addr())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	//go func() {
	//	for {
	go func() {
		for {
			s.boostrapNodes()
			time.Sleep(1 * time.Second)
		}
	}()
	//		time.Sleep(1 * time.Second)
	//	}
	//}()

free:
	for {
		//fmt.Println(s.ID, s.Transports)
		//fmt.Println()
		select {
		case rpc := <-s.rpcCh:
			fmt.Println("rpc received", rpc)
			msg, err := s.RPCDecodeFunc(rpc)
			fmt.Println(msg, err)
			if err != nil {
				//log.Printf("RPCDecodeFunc error: %v", err)
				fmt.Errorf(err.Error())
			}
			if err = s.RPCProcessor.ProcessMessage(msg); err != nil {
				fmt.Errorf(err.Error())
				//log.Printf("RPCProcessor error: %v", err)
			}

		case <-s.quitCh:
			break free

		}
	}
	fmt.Println("Server shutdown")
}
func (s *Server) boostrapNodes() {
	for _, tr := range s.Transports {
		//if tr.Addr() != s.Transport.Addr() {
		//if err := s.Transport.Connect(tr); err != nil {
		//	fmt.Println("error ", "can not connect to remote ", err)
		//}
		//}
		if err := s.sendGetStatusMessage(tr); err != nil {
			s.Logger.Log("error", err)
		}
		//s.Logger.Log(s.Transport.Addr(), "send get status message to ", tr.Addr())
	}
}
func (s *Server) ValidatorLoop() {
	//fmt.Println(s.ID, s.Transports)
	ticker := time.NewTicker(s.blocktime)
	//s.Logger.Log("msg", "Validator loop started", "blockTime", s.blocktime)
	for {
		<-ticker.C
		s.createNewBlock()
	}
}
func (s *Server) sendGetStatusMessage(tr Transport) error {
	//TODO
	//fmt.Println("sendGetStatusMessage---------ok")
	//getStatusMsg := new(GetStatusMessage)
	getStatusMsg := &GetStatusMessage{}
	buf := &bytes.Buffer{}
	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeGetStatus, buf.Bytes())
	//fmt.Println(msg.Bytes())
	if err := s.Transport.SendMessage(tr.Addr(), msg.Bytes()); err != nil {
		return err
	}
	//fmt.Printf("%s sending get status message to %s\n", s.Transport.Addr(), tr.Addr())
	return nil
}
func (s *Server) ProcessMessage(msg *DecodeMessage) error {
	//fmt.Println(msg)
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.ProcessTransaction(t)
	case *core.Block:
		//fmt.Println("block received", t)
		return s.ProcessBlock(t)
	case *GetStatusMessage:
		//fmt.Println("GetStatusMessage")
		//fmt.Println(s.Transport.Addr(), msg.From)
		return s.ProcessGetStatusMessage(msg.From, t)
	case *StatusMessage:
		return s.ProcessStatusMessage(msg.From, t)
	case *GetBlockMessage:
		fmt.Println("okk")
		return s.ProcessGetBlocksMessage(msg.From, t)
	default:
		fmt.Println("unknown msg type: %T", t)
	}
	return nil
}
func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}
	txx := s.memPool.Transactions()
	//fmt.Println("-----", len(txx))
	block, err := core.NewBlockFromPrevHeader(currentHeader, txx)
	//fmt.Println("------", len(block.Transactions))
	if err != nil {
		return err
	}
	if err = block.Sign(*s.PrivateKey); err != nil {
		return err
	}
	if err = s.chain.AddBlock(block); err != nil {
		return err
	}
	fmt.Println(s.ID, "created new block", block.DataHash)
	s.memPool.Flush()
	go s.broadcastBlock(block)
	//fmt.Println("createNewBlock")
	return nil
}
func (s *Server) initTransports() {
	//for _, tr := range s.Transports {
	go func(tr Transport) {
		for rpc := range tr.Consume() {
			//handle
			s.rpcCh <- rpc
		}
	}(s.Transport)
	//}
}
func geneisisBlock() *core.Block {
	header := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Height:    0,
		Timestamp: 00000,
	}
	return core.NewBlock(header, nil)
}

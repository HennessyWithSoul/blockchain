package core

import (
	"context"

	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"goblockchain/crypto"
	"goblockchain/rpc"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
	"time"
)

type GRPCServer struct {
	rpc.UnimplementedBlockchainServiceServer
	port       string
	grpcServer *grpc.Server
	peers      sync.Map // publicKey to  grpcPort
	pub        crypto.PublicKey
	pri        crypto.PrivateKey
	client     *clientv3.Client
	bc         *BlockChain
	tp         *TxnPool
	lg         *zap.Logger
}

func NewGRPCServer(port string, pubk crypto.PublicKey, prik crypto.PrivateKey, blockchain *BlockChain, etcdClient *clientv3.Client, lg *zap.Logger) *GRPCServer {
	RawGrpcServer := grpc.NewServer()
	return &GRPCServer{
		port:       port,
		grpcServer: RawGrpcServer,
		pub:        pubk,
		pri:        prik,
		client:     etcdClient,
		peers:      sync.Map{},
		bc:         blockchain,
		lg:         lg,
	}
}

func (s *GRPCServer) Start() error {

	rpc.RegisterBlockchainServiceServer(s.grpcServer, s)
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		s.lg.Error("failed to listen", zap.Error(err))
		return err
	}
	s.lg.Info("grpc server listening", zap.String("port", s.port))
	go func() {
		if err = s.grpcServer.Serve(lis); err != nil {
			s.lg.Info("failed to serve", zap.Error(err))
		}
	}()
	return nil
}

func (s *GRPCServer) Stop() error {
	return nil
}

func (s *GRPCServer) SubmitTransaction(ctx context.Context, req *rpc.SubmitTransactionRequest) (*rpc.Empty, error) {
	return &rpc.Empty{}, nil
}

func (s *GRPCServer) SubmitBlock(ctx context.Context, req *rpc.SubmitBlockRequest) (*rpc.Empty, error) {
	return &rpc.Empty{}, nil
}

func (s *GRPCServer) GetBlock(ctx context.Context, req *rpc.GetBlocksRequest) (*rpc.GetBlocksResponse, error) {
	return &rpc.GetBlocksResponse{}, nil
}

func (s *GRPCServer) GetStatus(ctx context.Context, req *rpc.GetStatusRequest) (*rpc.GetStatusResponse, error) {
	return &rpc.GetStatusResponse{}, nil
}

func (s *GRPCServer) HandleTxn(value string, address string) error {
	int_value, err := strconv.Atoi(value)
	if err != nil {
		s.lg.Error("invalid value", zap.String("value", value), zap.Error(err))
		return err
	}
	publicKey, err := stringToECDSAPublicKey(address)
	if err != nil {
		s.lg.Error("invalid address", zap.String("address", address), zap.Error(err))
		return err
	}
	txn := &Transaction{
		To:    crypto.PublicKey{publicKey},
		From:  s.pub,
		Value: uint64(int_value),
	}
	txn.SetFirstSeen(time.Now().Unix())
	txn.Sign(s.pri)

	s.lg.Info("Created transaction", zap.String("to", address), zap.Uint64("value", uint64(int_value)))
	return nil
}

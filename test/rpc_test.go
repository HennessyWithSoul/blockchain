package test

import (
	"context"
	"fmt"
	"goblockchain/core"
	"goblockchain/crypto"
	"goblockchain/rpc"
	"goblockchain/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"testing"
	"time"
)

type blockchainServer struct {
	rpc.UnimplementedBlockchainServiceServer
}

func NewBlockchainServer() *blockchainServer {
	return &blockchainServer{}
}

func (s *blockchainServer) SubmitTransaction(ctx context.Context, req *rpc.SubmitTransactionRequest) (*rpc.Empty, error) {
	grpcTx := req.GetTransaction()
	if grpcTx == nil {
		return nil, status.Error(codes.InvalidArgument, "transaction is required")
	}

	// 转换为内部交易类型
	tx, err := core.FromGRPCTransaction(grpcTx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid transaction: "+err.Error())
	}
	log.Println(tx.From)
	log.Println(tx.To)
	log.Println(tx.Value)
	log.Println(string(tx.Data))

	return &rpc.Empty{}, nil
}

func (s *blockchainServer) SubmitBlock(ctx context.Context, req *rpc.SubmitBlockRequest) (*rpc.Empty, error) {
	grpcBlock := req.GetBlock()
	if grpcBlock == nil {
		return nil, status.Error(codes.InvalidArgument, "block is required")
	}

	// 转换为内部区块类型
	block, err := core.FromGRPCBlock(grpcBlock)
	fmt.Printf("++++++%+v\n", block.Transactions[0])
	fmt.Printf("++++++%+v\n", block.Transactions[1])
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid block: "+err.Error())
	}

	// 验证区块
	if err := block.Verify(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "block validation failed: "+err.Error())
	}
	return &rpc.Empty{}, nil
	// 添加到区块链

}

func TestTxnRpc(t *testing.T) {
	pri1 := crypto.GeneratePrivateKey()
	pub1 := pri1.PublicKey()
	pri2 := crypto.GeneratePrivateKey()
	pub2 := pri2.PublicKey()
	log.Println("发送的pub1", pub1)
	log.Println("发送的pub2", pub2)
	go func() {
		time.Sleep(5 * time.Second)
		conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		client := rpc.NewBlockchainServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		txn := &core.Transaction{
			To:    pub2,
			From:  pub1,
			Value: 1000,
			Data:  []byte("test transaction"),
		}
		txn.SetFirstSeen(time.Now().Unix())
		txn.Sign(pri1)
		grpcTx, err := core.ToGRPCTransaction(txn)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := client.SubmitTransaction(ctx, &rpc.SubmitTransactionRequest{
			Transaction: grpcTx,
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Transaction submitted successfully:", resp)
	}()
	go func() {
		server := NewBlockchainServer()
		grpcServer := grpc.NewServer()
		rpc.RegisterBlockchainServiceServer(grpcServer, server)
		lis, err := net.Listen("tcp", ":50053")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Println("Blockchain gRPC server started on :50053")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	time.Sleep(10 * time.Second)

}

func TestBlockRpc(t *testing.T) {
	pri1 := crypto.GeneratePrivateKey()
	pub1 := pri1.PublicKey()
	pri2 := crypto.GeneratePrivateKey()
	pub2 := pri2.PublicKey()
	log.Println("发送的pub1", pub1)
	log.Println("发送的pub2", pub2)
	go func() {
		bc, err := core.NewBlockChain(geneisisBlock())
		conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		client := rpc.NewBlockchainServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		txn1 := &core.Transaction{
			To:    pub2,
			From:  pub1,
			Value: 1000,
			Data:  []byte("test transaction 1"),
		}
		txn1.SetFirstSeen(time.Now().Unix())
		txn1.Sign(pri1)
		txn2 := &core.Transaction{
			To:    pub1,
			From:  pub2,
			Value: 500,
			Data:  []byte("test transaction 2"),
		}
		txn2.SetFirstSeen(time.Now().Unix())
		txn2.Sign(pri2)
		preheader, err := bc.GetHeader(bc.Height())
		transBlock, err := core.NewBlockFromPrevHeader(preheader, []*core.Transaction{txn1, txn2})
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(5 * time.Second)
		transBlock.Sign(pri1)
		grpcBlock, err := core.ToGRPCBlock(transBlock)
		//fmt.Printf("++++++%+v\n", transBlock)
		transBlock, err = core.FromGRPCBlock(grpcBlock)
		//fmt.Printf("++++++%+v\n", transBlock)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := client.SubmitBlock(ctx, &rpc.SubmitBlockRequest{
			Block: grpcBlock,
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Transaction submitted successfully:", resp)
	}()
	go func() {
		server := NewBlockchainServer()
		grpcServer := grpc.NewServer()
		rpc.RegisterBlockchainServiceServer(grpcServer, server)
		lis, err := net.Listen("tcp", ":50053")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Println("Blockchain gRPC server started on :50053")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	time.Sleep(10 * time.Second)

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

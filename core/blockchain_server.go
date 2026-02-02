package core

type BlockchainServer struct {
	bc         *BlockChain
	grpcServer *GRPCServer
	httpServer *HttpServer
	utxo       *UTXOSet
	calculator *HashCalculator
}

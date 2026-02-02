package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"goblockchain/core"
	"goblockchain/crypto"
	"log"
	"os"
)

var (
	name          string
	privateKey    string
	grpcPort      string
	httpPort      string
	etcdEndpoints []string
	logLevel      string
)
var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "A brief description of your application",
	Long:  "A longer description...",
	Run:   execute,
}

func execute(cmd *cobra.Command, args []string) {
	lg, err := zap.NewProduction()
	if err != nil {
		// 如果 zap 初始化失败，使用标准库日志
		log.Printf("Failed to create logger: %v", err)
		panic("cannot create logger: " + err.Error())
	}
	defer lg.Sync()

	pri, err := crypto.PrivateKeyFromHex(privateKey)
	if err != nil {
		lg.Error("cannot parse private key", zap.Error(err), zap.String("given pri key:", privateKey))
		panic("cannot parse private key: " + err.Error())
	}
	geneisis := core.GeneisisBlock()
	blockchain, err := core.NewBlockChain(geneisis)
	if err != nil {
		lg.Error("cannot create blockchain", zap.Error(err))
		panic("cannot create blockchain: " + err.Error())
	}
	grpcserver := core.NewGRPCServer(grpcPort, pri.PublicKey(), blockchain, nil, lg)
	err = grpcserver.Start()
	if err != nil {
		lg.Error("cannot start grpc server", zap.Error(err))
		panic("cannot start grpc server: " + err.Error())
	}
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String(name, "world", "name to greet")
	rootCmd.Flags().StringVarP(&privateKey, "private-key", "k", "4c7c9a7f8b2e3d4a5f6e1c8d9a0b3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d", "Private key string")
	rootCmd.Flags().StringVarP(&grpcPort, "grpc-port", "g", ":50351", "gRPC server port")
	rootCmd.Flags().StringVarP(&httpPort, "http-port", "p", ":8080", "HTTP server port")
	rootCmd.Flags().StringSliceVar(&etcdEndpoints, "etcd-endpoints", []string{"127.0.0.1:2379"}, "etcd endpoints")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
}

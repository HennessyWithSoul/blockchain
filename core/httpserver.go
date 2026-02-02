package core

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HttpServer struct {
	//TODO
	httpport   string
	grpcserver *GRPCServer
	blockchain *BlockChain
	lg         *zap.Logger
}

func NewHttpServer(port string, server *GRPCServer, chain *BlockChain, lg *zap.Logger) *HttpServer {
	return &HttpServer{
		httpport:   port,
		grpcserver: server,
		blockchain: chain,
		lg:         lg,
	}
}

func (hs *HttpServer) Start() error {
	r := gin.Default()
	r.POST("/AddTxn", hs.addTxn())
	return nil
}

func (hs *HttpServer) addTxn() gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.PostForm("value")
		if value == "" {
			hs.lg.Warn("value is empty")
			c.JSON(200, gin.H{
				"status":  "error",
				"message": "value is empty",
			})
			return
		}
		address := c.PostForm("address")
		if address == "" {
			hs.lg.Warn("address is empty")
			c.JSON(200, gin.H{
				"status":  "error",
				"message": "address is empty",
			})
			return
		}

		hs.lg.Info("Received transaction", zap.String("value", value), zap.String("address", address))
		hs.grpcserver.HandleTxn(value, address)
	}
}

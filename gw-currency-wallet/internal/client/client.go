package client

import (
	"gw-currency-wallet/pkg/logger"
	"os"

	pb "github.com/Stevenrunner4343/proto-exchange/exchange"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	pb.ExchangeServiceClient
	conn *grpc.ClientConn
}

func NewGrpcClient() *GrpcClient {
	grpcHost := os.Getenv("GRPC_HOST")

	conn, err := grpc.NewClient(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Ошибка при создании gRPC клиента", zap.Error(err))
		return nil
	}

	return &GrpcClient{
		ExchangeServiceClient: pb.NewExchangeServiceClient(conn),
		conn:                  conn,
	}
}

func (c *GrpcClient) Close() error {
	return c.conn.Close()
}

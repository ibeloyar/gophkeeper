package grpcclient

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

type GRPCClient struct {
	conn *grpc.ClientConn
	Cmd  gophkeeperv1.GophkeeperClient
}

func New(addr string) *GRPCClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	return &GRPCClient{
		conn: conn,
		Cmd:  gophkeeperv1.NewGophkeeperClient(conn),
	}
}

func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

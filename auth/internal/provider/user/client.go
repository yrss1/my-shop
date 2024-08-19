package user

import (
	pb "github.com/yrss1/proto-definitions/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

func New(address string) (client *Client, err error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	client = &Client{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
	}

	return
}

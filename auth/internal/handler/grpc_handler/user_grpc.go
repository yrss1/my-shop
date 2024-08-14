package grpc_handler

import (
	"context"
	"github.com/yrss1/my-shop/tree/main/user/internal/service/shop"
	"github.com/yrss1/my-shop/tree/main/user/pb"
	"log"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	shopService *shop.Service
}

func NewUserServiceServer(s *shop.Service) *UserServiceServer {
	return &UserServiceServer{shopService: s}
}

func (s *UserServiceServer) SayHello(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	log.Printf("message from client: %s", in.Body)
	return &pb.Message{Body: "hello my name is Test"}, nil
}

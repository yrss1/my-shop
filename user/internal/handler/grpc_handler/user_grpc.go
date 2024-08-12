package grpc_handler

import (
	"context"
	"github.com/yrss1/my-shop/tree/main/user/internal/service/shop"
	"log"
)

type UserServiceServer struct {
	proto.UnimplementedUserServiceServer
	shopService *shop.Service
}

//func NewUserServiceServer(s *shop.Service) *UserServiceServer {
//	return &UserServiceServer{shopService: s}
//}

func (s *UserServiceServer) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Получено сообщение от клиента: %s", in.Body)
	return &Message{Body: "Привет от сервера!"}, nil
}

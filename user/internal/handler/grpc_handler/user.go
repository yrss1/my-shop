package grpc_handler

import (
	"context"
	"errors"
	pb "github.com/yrss1/my-shop/pkg/pb/user"
	"github.com/yrss1/my-shop/user/internal/domain/user"
	"github.com/yrss1/my-shop/user/internal/service/userService"
	"github.com/yrss1/my-shop/user/pkg/helpers"
	"github.com/yrss1/my-shop/user/pkg/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService *userService.Service
}

func NewUserServiceServer(s *userService.Service) *UserServiceServer {
	return &UserServiceServer{userService: s}
}

func (s *UserServiceServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (res *pb.UserResponse, err error) {
	user, err := s.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, req.Email)
		default:
			return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
		}
		return
	}

	res = &pb.UserResponse{
		User: &pb.Response{
			Id:      user.ID,
			Name:    user.Name,
			Email:   user.Email,
			Address: user.Address,
			Role:    user.Role,
		},
	}

	return
}

func (s *UserServiceServer) RegisterUser(ctx context.Context, req *pb.UserRequest) (res *pb.UserResponse, err error) {
	if req == nil || req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid input: user request is nil")
	}

	userRequest := user.Request{
		Name:     helpers.GetStringPtr(req.User.Name),
		Email:    helpers.GetStringPtr(req.User.Email),
		Password: helpers.GetStringPtr(req.User.Password),
	}
	if req.User.Address != "" {
		userRequest.Address = helpers.GetStringPtr(req.User.Address)
	}
	if req.User.Role != "" {
		userRequest.Role = helpers.GetStringPtr(req.User.Role)
	}

	createdUser, err := s.userService.CreateUser(ctx, userRequest)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	res = &pb.UserResponse{User: &pb.Response{
		Id:      createdUser.ID,
		Name:    createdUser.Name,
		Email:   createdUser.Email,
		Address: createdUser.Address,
		Role:    createdUser.Role,
	}}

	return

}

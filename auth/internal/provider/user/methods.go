package user

import (
	"context"
	pb "github.com/yrss1/proto-definitions/user"
	"time"
)

func (c *Client) GetUserByEmail(ctx context.Context, email string) (res *pb.UserResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pb.GetUserByEmailRequest{Email: email}
	res, err = c.client.GetUserByEmail(ctx, req)
	if err != nil {
		return
	}

	return
}

func (c *Client) RegisterUser(ctx context.Context, user *pb.UserRequest) (*pb.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.RegisterUser(ctx, user)
}

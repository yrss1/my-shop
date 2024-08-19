package authService

import (
	"context"
	"github.com/yrss1/my-shop/auth/pkg/log"
	pb "github.com/yrss1/proto-definitions/user"
	"go.uber.org/zap"
)

func (s *Service) GetUserByEmail(ctx context.Context, email string) (res *pb.UserResponse, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetUserByEmail").With(zap.String("email", email))

	res, err = s.userClient.GetUserByEmail(ctx, email)
	if err != nil {
		logger.Error("failed to delete by id", zap.Error(err))
		return
	}

	return
}

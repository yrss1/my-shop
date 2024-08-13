package shop

import (
	"context"
	"errors"
	"github.com/yrss1/my-shop/tree/main/user/internal/domain/user"
	"github.com/yrss1/my-shop/tree/main/user/pkg/log"
	"github.com/yrss1/my-shop/tree/main/user/pkg/store"
	"go.uber.org/zap"
)

func (s *Service) ListUsers(ctx context.Context) (res []user.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListUsers")

	data, err := s.userRepository.List(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}

	res = user.ParseFromEntities(data)

	return
}

func (s *Service) CreateUser(ctx context.Context, req user.Request) (res user.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("CreateUser")

	data := user.Entity{
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
		Role:    req.Role,
	}

	data.ID, err = s.userRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}

	res = user.ParseFromEntity(data)

	return
}

func (s *Service) GetUser(ctx context.Context, id string) (res user.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetUser").With(zap.String("id", id))

	data, err := s.userRepository.Get(ctx, id)
	if err != nil {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}

	res = user.ParseFromEntity(data)

	return
}

func (s *Service) UpdateUser(ctx context.Context, id string, req user.Request) (err error) {
	logger := log.LoggerFromContext(ctx).Named("UpdateBook").With(zap.String("id", id))

	data := user.Entity{
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
		Role:    req.Role,
	}

	err = s.userRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to update by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) DeleteUser(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("DeleteUser").With(zap.String("id", id))

	err = s.userRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to delete by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) SearchUser(ctx context.Context, req user.Request) (res []user.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("SearchUser")

	if req.Name != nil {
		logger = logger.With(zap.String("name", *(req.Name)))
	}
	if req.Email != nil {
		logger = logger.With(zap.String("email", *(req.Email)))
	}

	searchData := user.Entity{
		Name:  req.Name,
		Email: req.Email,
	}
	data, err := s.userRepository.Search(ctx, searchData)
	if err != nil {
		logger.Error("failed to search users", zap.Error(err))
		return
	}

	res = user.ParseFromEntities(data)

	return
}

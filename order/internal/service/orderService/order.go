package orderService

import (
	"context"
	"errors"
	"github.com/yrss1/my-shop/tree/main/order/internal/domain/order"
	"github.com/yrss1/my-shop/tree/main/order/pkg/log"
	"github.com/yrss1/my-shop/tree/main/order/pkg/store"
	"go.uber.org/zap"
)

func (s *Service) ListOrders(ctx context.Context) (res []order.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListOrders")

	data, err := s.orderRepository.List(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}

	res = order.ParseFromEntities(data)

	return
}

func (s *Service) CreateOrder(ctx context.Context, req order.Request) (res order.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("CreateOrder")

	data := order.Entity{
		UserID:     req.UserID,
		Products:   req.Products,
		TotalPrice: req.TotalPrice,
		Status:     req.Status,
	}

	data.ID, err = s.orderRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}

	res = order.ParseFromEntity(data)

	return
}

func (s *Service) GetOrder(ctx context.Context, id string) (res order.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetOrder").With(zap.String("id", id))

	data, err := s.orderRepository.Get(ctx, id)
	if err != nil {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}

	res = order.ParseFromEntity(data)

	return
}

func (s *Service) UpdateOrder(ctx context.Context, id string, req order.Request) (err error) {
	logger := log.LoggerFromContext(ctx).Named("UpdateOrder").With(zap.String("id", id))

	data := order.Entity{
		UserID:     req.UserID,
		Products:   req.Products,
		TotalPrice: req.TotalPrice,
		Status:     req.Status,
	}

	err = s.orderRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to update by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) DeleteOrder(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("DeleteUser").With(zap.String("id", id))

	err = s.orderRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to delete by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) SearchOrder(ctx context.Context, req order.Request) (res []order.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("SearchOrder")

	if req.UserID != nil {
		logger = logger.With(zap.String("user_id", *(req.UserID)))
	}
	if req.Status != nil {
		logger = logger.With(zap.String("status", *(req.Status)))
	}

	searchData := order.Entity{
		UserID: req.UserID,
		Status: req.Status,
	}

	data, err := s.orderRepository.Search(ctx, searchData)
	if err != nil {
		logger.Error("failed to search orders", zap.Error(err))
		return
	}

	res = order.ParseFromEntities(data)

	return
}

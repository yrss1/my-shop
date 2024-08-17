package epayment

import (
	"context"
	"errors"
	"github.com/yrss1/my-shop/payment/internal/domain/payment"
	"github.com/yrss1/my-shop/payment/pkg/log"
	"github.com/yrss1/my-shop/payment/pkg/store"
	"go.uber.org/zap"
)

func (s *Service) ListPayments(ctx context.Context) (res []payment.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListPayments")

	data, err := s.paymentRepository.List(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}

	res = payment.ParseFromEntities(data)

	return
}

func (s *Service) CreatePayment(ctx context.Context, req payment.Request) (res payment.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("CreatePayment")

	data := payment.Entity{
		UserID:  req.UserID,
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Status:  req.Status,
	}

	data.ID, err = s.paymentRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}

	res = payment.ParseFromEntity(data)

	return
}

func (s *Service) GetPayment(ctx context.Context, id string) (res payment.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetPayment").With(zap.String("id", id))

	data, err := s.paymentRepository.Get(ctx, id)
	if err != nil {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}

	res = payment.ParseFromEntity(data)

	return
}

func (s *Service) UpdatePayment(ctx context.Context, id string, req payment.Request) (err error) {
	logger := log.LoggerFromContext(ctx).Named("UpdatePayment").With(zap.String("id", id))

	data := payment.Entity{
		UserID:  req.UserID,
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Status:  req.Status,
	}

	err = s.paymentRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to update by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) DeletePayment(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("DeletePayment").With(zap.String("id", id))

	err = s.paymentRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to delete by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) SearchPayment(ctx context.Context, req payment.Request) (res []payment.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("SearchUser")

	if req.UserID != nil {
		logger = logger.With(zap.String("user_id", *(req.UserID)))
	}
	if req.OrderID != nil {
		logger = logger.With(zap.String("order_id", *(req.OrderID)))
	}
	if req.Status != nil {
		logger = logger.With(zap.String("status", *(req.Status)))
	}

	searchData := payment.Entity{
		UserID:  req.UserID,
		OrderID: req.OrderID,
		Status:  req.Status,
	}
	data, err := s.paymentRepository.Search(ctx, searchData)
	if err != nil {
		logger.Error("failed to search payments", zap.Error(err))
		return
	}

	res = payment.ParseFromEntities(data)

	return
}

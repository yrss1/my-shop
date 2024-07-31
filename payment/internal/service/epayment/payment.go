package epayment

import (
	"context"
	"errors"
	"fmt"
)

func (s *Service) ListPayments(ctx context.Context) (res []payment.Response, err error) {
	data, err := s.paymentRepository.List(ctx)
	if err != nil {
		fmt.Printf("failed to select: %v\n", err)
		return
	}

	res = payment.ParseFromEntities(data)

	return
}

func (s *Service) CreatePayment(ctx context.Context, req payment.Request) (res payment.Response, err error) {
	data := payment.Entity{
		UserID:  req.UserID,
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Status:  req.Status,
	}

	data.ID, err = s.paymentRepository.Add(ctx, data)
	if err != nil {
		fmt.Printf("faled to create: %v\n", err)
		return
	}

	res = payment.ParseFromEntity(data)

	return
}

func (s *Service) GetPayment(ctx context.Context, id string) (res payment.Response, err error) {
	data, err := s.paymentRepository.Get(ctx, id)
	if err != nil {
		fmt.Printf("failed to get by id: %v\n", err)
		return
	}

	res = payment.ParseFromEntity(data)

	return
}

func (s *Service) UpdatePayment(ctx context.Context, id string, req payment.Request) (err error) {
	data := payment.Entity{
		UserID:  req.UserID,
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Status:  req.Status,
	}

	err = s.paymentRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to update by id: %v\n", err)
		return
	}

	return
}

func (s *Service) DeletePayment(ctx context.Context, id string) (err error) {
	err = s.paymentRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to delete by id: %v\n", err)
		return
	}

	return
}

func (s *Service) SearchPayment(ctx context.Context, req payment.Request) (res []payment.Response, err error) {
	searchData := payment.Entity{
		UserID:  req.UserID,
		OrderID: req.OrderID,
		Status:  req.Status,
	}
	data, err := s.paymentRepository.Search(ctx, searchData)
	if err != nil {
		fmt.Printf("failed to search payments: %v\n", err)
		return
	}

	res = payment.ParseFromEntities(data)

	return
}

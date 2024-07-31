package shop

import (
	"context"
	"errors"
	"fmt"
)

func (s *Service) ListOrders(ctx context.Context) (res []order.Response, err error) {
	data, err := s.orderRepository.List(ctx)
	if err != nil {
		fmt.Printf("failed to select: %v\n", err)
		return
	}

	res = order.ParseFromEntities(data)

	return
}

func (s *Service) CreateOrder(ctx context.Context, req order.Request) (res order.Response, err error) {
	data := order.Entity{
		UserID:     req.UserID,
		Products:   req.Products,
		TotalPrice: req.TotalPrice,
		Status:     req.Status,
	}

	data.ID, err = s.orderRepository.Add(ctx, data)
	if err != nil {
		fmt.Printf("faled to create: %v\n", err)
		return
	}

	res = order.ParseFromEntity(data)

	return
}

func (s *Service) GetOrder(ctx context.Context, id string) (res order.Response, err error) {
	data, err := s.orderRepository.Get(ctx, id)
	if err != nil {
		fmt.Printf("failed to get by id: %v\n", err)
		return
	}

	res = order.ParseFromEntity(data)

	return
}

func (s *Service) UpdateOrder(ctx context.Context, id string, req order.Request) (err error) {
	data := order.Entity{
		UserID:     req.UserID,
		Products:   req.Products,
		TotalPrice: req.TotalPrice,
		Status:     req.Status,
	}

	err = s.orderRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to update by id: %v\n", err)
		return
	}

	return
}

func (s *Service) DeleteOrder(ctx context.Context, id string) (err error) {
	err = s.orderRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to delete by id: %v\n", err)
		return
	}

	return
}

func (s *Service) SearchOrder(ctx context.Context, req order.Request) (res []order.Response, err error) {
	searchData := order.Entity{
		UserID: req.UserID,
		Status: req.Status,
	}

	data, err := s.orderRepository.Search(ctx, searchData)
	if err != nil {
		fmt.Printf("failed to search products: %v\n", err)
		return
	}

	res = order.ParseFromEntities(data)

	return
}

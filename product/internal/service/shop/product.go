package shop

import (
	"context"
	"errors"
	"fmt"
	"github.com/yrss1/my-shop/tree/main/product/internal/domain/product"
	"github.com/yrss1/my-shop/tree/main/product/pkg/store"
)

func (s *Service) ListProducts(ctx context.Context) (res []product.Response, err error) {
	data, err := s.productRepository.List(ctx)
	if err != nil {
		fmt.Printf("failed to select: %v\n", err)
		return
	}

	res = product.ParseFromEntities(data)

	return
}

func (s *Service) CreateProduct(ctx context.Context, req product.Request) (res product.Response, err error) {
	data := product.Entity{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Quantity:    req.Quantity,
	}

	data.ID, err = s.productRepository.Add(ctx, data)
	if err != nil {
		fmt.Printf("faled to create: %v\n", err)
		return
	}

	res = product.ParseFromEntity(data)

	return
}

func (s *Service) GetProduct(ctx context.Context, id string) (res product.Response, err error) {
	data, err := s.productRepository.Get(ctx, id)
	if err != nil {
		fmt.Printf("failed to get by id: %v\n", err)
		return
	}

	res = product.ParseFromEntity(data)

	return
}

func (s *Service) UpdateProduct(ctx context.Context, id string, req product.Request) (err error) {
	data := product.Entity{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Quantity:    req.Quantity,
	}

	err = s.productRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to update by id: %v\n", err)
		return
	}

	return
}

func (s *Service) DeleteProduct(ctx context.Context, id string) (err error) {
	err = s.productRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to delete by id: %v\n", err)
		return
	}

	return
}

func (s *Service) SearchProduct(ctx context.Context, req product.Request) (res []product.Response, err error) {
	searchData := product.Entity{
		Name:     req.Name,
		Category: req.Category,
	}
	data, err := s.productRepository.Search(ctx, searchData)
	if err != nil {
		fmt.Printf("failed to search products: %v\n", err)
		return
	}

	res = product.ParseFromEntities(data)

	return
}

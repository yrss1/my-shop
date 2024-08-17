package productService

import (
	"context"
	"errors"
	"github.com/yrss1/my-shop/product/internal/domain/product"
	"github.com/yrss1/my-shop/product/pkg/log"
	"github.com/yrss1/my-shop/product/pkg/store"
	"go.uber.org/zap"
)

func (s *Service) ListProducts(ctx context.Context) (res []product.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("ListProducts")

	data, err := s.productRepository.List(ctx)
	if err != nil {
		logger.Error("failed to select", zap.Error(err))
		return
	}

	res = product.ParseFromEntities(data)

	return
}

func (s *Service) CreateProduct(ctx context.Context, req product.Request) (res product.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("CreateProduct")

	data := product.Entity{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Quantity:    req.Quantity,
	}

	data.ID, err = s.productRepository.Add(ctx, data)
	if err != nil {
		logger.Error("failed to create", zap.Error(err))
		return
	}

	res = product.ParseFromEntity(data)

	return
}

func (s *Service) GetProduct(ctx context.Context, id string) (res product.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("GetProduct").With(zap.String("id", id))

	data, err := s.productRepository.Get(ctx, id)
	if err != nil {
		logger.Error("failed to get by id", zap.Error(err))
		return
	}

	res = product.ParseFromEntity(data)

	return
}

func (s *Service) UpdateProduct(ctx context.Context, id string, req product.Request) (err error) {
	logger := log.LoggerFromContext(ctx).Named("UpdateProduct").With(zap.String("id", id))

	data := product.Entity{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Quantity:    req.Quantity,
	}

	err = s.productRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to update by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) DeleteProduct(ctx context.Context, id string) (err error) {
	logger := log.LoggerFromContext(ctx).Named("DeleteProduct").With(zap.String("id", id))

	err = s.productRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		logger.Error("failed to delete by id", zap.Error(err))
		return
	}

	return
}

func (s *Service) SearchProduct(ctx context.Context, req product.Request) (res []product.Response, err error) {
	logger := log.LoggerFromContext(ctx).Named("SearchProduct")

	if req.Name != nil {
		logger = logger.With(zap.String("name", *(req.Name)))
	}
	if req.Category != nil {
		logger = logger.With(zap.String("category", *(req.Category)))
	}

	searchData := product.Entity{
		Name:     req.Name,
		Category: req.Category,
	}
	data, err := s.productRepository.Search(ctx, searchData)
	if err != nil {
		logger.Error("failed to search products", zap.Error(err))
		return
	}

	res = product.ParseFromEntities(data)

	return
}

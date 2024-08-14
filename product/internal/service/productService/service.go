package productService

import "github.com/yrss1/my-shop/tree/main/product/internal/domain/product"

type Configuration func(s *Service) error

type Service struct {
	productRepository product.Repository
}

func New(configs ...Configuration) (s *Service, err error) {
	s = &Service{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func WithProductRepository(productRepository product.Repository) Configuration {
	return func(s *Service) error {
		s.productRepository = productRepository
		return nil
	}
}

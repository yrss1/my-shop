package shop

import "github.com/yrss1/my-shop/tree/main/order/internal/domain/order"

type Configuration func(s *Service) error

type Service struct {
	orderRepository order.Repository
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

func WithOrderRepository(orderRepository order.Repository) Configuration {
	return func(s *Service) error {
		s.orderRepository = orderRepository
		return nil
	}
}

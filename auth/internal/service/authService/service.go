package authService

import (
	"github.com/yrss1/my-shop/auth/internal/provider/user"
)

type Configuration func(s *Service) error

type Service struct {
	userClient *user.Client
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

func WithUserClient(userClient *user.Client) Configuration {
	return func(s *Service) error {
		s.userClient = userClient
		return nil
	}
}

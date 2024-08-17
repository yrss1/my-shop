package repository

import (
	"github.com/yrss1/my-shop/order/internal/domain/order"
	"github.com/yrss1/my-shop/order/internal/repository/postgres"
	"github.com/yrss1/my-shop/order/pkg/store"
)

type Configuration func(r *Repository) error

type Repository struct {
	postgres store.SQLX

	Order order.Repository
}

func New(configs ...Configuration) (s *Repository, err error) {
	s = &Repository{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func WithPostgresStore(dbName string) Configuration {
	return func(r *Repository) (err error) {
		r.postgres, err = store.New(dbName)
		if err != nil {
			return
		}
		//if err = store.Migrate(dbName); err != nil {
		//	return
		//}

		r.Order = postgres.NewOrderRepository(r.postgres.Client)

		return
	}
}

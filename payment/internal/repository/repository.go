package repository

import (
	"github.com/yrss1/my-shop/payment/internal/domain/payment"
	"github.com/yrss1/my-shop/payment/internal/repository/postgres"
	"github.com/yrss1/my-shop/payment/pkg/store"
)

type Configuration func(r *Repository) error

type Repository struct {
	postgres store.SQLX

	Payment payment.Repository
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
		r.Payment = postgres.NewPaymentRepository(r.postgres.Client)

		return
	}
}

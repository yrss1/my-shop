package repository

import (
	"fmt"
	"github.com/yrss1/my-shop/tree/main/product/internal/domain/product"
	"github.com/yrss1/my-shop/tree/main/product/internal/repository/postgres"
	"github.com/yrss1/my-shop/tree/main/product/pkg/store"
)

type Configuration func(r *Repository) error

type Repository struct {
	postgres store.SQLX

	Product product.Repository
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
			fmt.Println("postgres")
			return
		}
		//if err = store.Migrate(dbName); err != nil {
		//	return
		//}

		r.Product = postgres.NewProductRepository(r.postgres.Client)

		return
	}
}

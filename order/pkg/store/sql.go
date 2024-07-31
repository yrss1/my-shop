package store

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strings"
)

type SQLX struct {
	Client *sqlx.DB
}

func New(dbSource string) (store SQLX, err error) {
	driverName := strings.ToLower(strings.Split(dbSource, "://")[0])
	store.Client, err = sqlx.Connect(driverName, dbSource)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v", err)
		return
	}
	store.Client.SetMaxOpenConns(20)

	return
}

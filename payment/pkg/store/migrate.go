package store

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"strings"
)

func Migrate(dataSourceName string) (err error) {
	if !strings.Contains(dataSourceName, "://") {
		err = errors.New("store: undefined data source name " + dataSourceName)
		return
	}

	migrations, err := migrate.New("file://db/migrations", dataSourceName)
	if err != nil {
		return
	}

	if err = migrations.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
	}

	return
}

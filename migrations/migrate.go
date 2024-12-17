package migrations

import (
	"errors"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func MigrateDb(files fs.FS, filePath, dbURL string) error {
	driver, err := iofs.New(files, filePath)
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", driver, dbURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}
	return nil
}

package database

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/wb-go/wbf/dbpg"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Connect(masterDSN string, slaveDSNs []string, opts *dbpg.Options) (*dbpg.DB, error) {
	db, err := dbpg.New(masterDSN, slaveDSNs, opts)
	if err != nil {
		return nil, err
	}

	if err := db.Master.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(migrationDir string, dbURL string) error {
	m, err := migrate.New("file://"+migrationDir, dbURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

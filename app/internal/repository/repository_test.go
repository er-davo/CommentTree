//go:build integration
// +build integration

package repository_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"comment-tree/internal/database"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/wb-go/wbf/dbpg"
)

var db *dbpg.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		tc.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(10*time.Second)),
	)
	if err != nil {
		log.Fatal(err)
	}

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate("../../../migrations", dsn); err != nil {
		log.Fatal(err)
	}

	db, err = dbpg.New(dsn, []string{}, &dbpg.Options{
		MaxOpenConns:    2,
		MaxIdleConns:    1,
		ConnMaxLifetime: time.Minute,
	})
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	db.Master.Close()
	_ = pgContainer.Terminate(ctx)
	os.Exit(code)
}

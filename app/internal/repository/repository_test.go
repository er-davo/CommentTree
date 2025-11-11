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
	"comment-tree/internal/models"
	"comment-tree/internal/repository"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
)

var db *dbpg.DB

var strategy = retry.Strategy{
	Attempts: 1,
	Delay:    1 * time.Second,
	Backoff:  1,
}

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

func TestCommentsRepository_CUD(t *testing.T) {
	repo := repository.NewCommentsRepository(db, strategy)

	com := models.Comment{
		Content:   "Test Content",
		CreatedAt: time.Now().UTC(),
		ParentID:  nil,
	}

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(t.Context(), &com)
		require.NoError(t, err)
	})

	com.Content = "Updated Content"

	t.Run("Update", func(t *testing.T) {
		err := repo.Update(t.Context(), &com)
		require.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(t.Context(), com.ID)
		require.NoError(t, err)
	})
}

func TestCommentsRepository_Get(t *testing.T) {
	repo := repository.NewCommentsRepository(db, strategy)

	rootCom := models.Comment{
		Content:   "Test Content",
		CreatedAt: time.Now().UTC(),
		ParentID:  nil,
	}
	err := repo.Create(t.Context(), &rootCom)
	require.NoError(t, err)

	rootComChildren1 := models.Comment{
		Content:   "Test Content",
		CreatedAt: time.Now().UTC(),
		ParentID:  &rootCom.ID,
	}

	err = repo.Create(t.Context(), &rootComChildren1)
	require.NoError(t, err)

	rootComChildren2 := models.Comment{
		Content:   "Test Content",
		CreatedAt: time.Now().UTC(),
		ParentID:  &rootCom.ID,
	}

	err = repo.Create(t.Context(), &rootComChildren2)
	require.NoError(t, err)

	coms, err := repo.GetByParent(t.Context(), nil, 10, 0)
	expectedRootCom := []*models.Comment{&rootCom}

	require.NoError(t, err)
	require.Len(t, coms, len(expectedRootCom))

	chlComs, err := repo.GetByParent(t.Context(), &rootCom.ID, 10, 0)

	expectedChlComs := []*models.Comment{&rootComChildren1, &rootComChildren2}

	require.NoError(t, err)
	require.Len(t, chlComs, len(expectedChlComs))
}

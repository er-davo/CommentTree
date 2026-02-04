package service_test

import (
	"context"
	"errors"
	"testing"

	"comment-tree/internal/mocks"
	"comment-tree/internal/models"
	"comment-tree/internal/service"

	"github.com/stretchr/testify/require"
	"github.com/wb-go/wbf/zlog"
	"go.uber.org/mock/gomock"
)

func newTestService(t *testing.T) (*service.CommentsService, *mocks.MockCommentsRepository, context.Context) {
	t.Helper()

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mocks.NewMockCommentsRepository(ctrl)
	log := &zlog.Zerolog{}
	svc := service.NewCommentsService(repo, log)

	return svc, repo, context.Background()
}

func TestCommentsService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, repo, ctx := newTestService(t)

		com := &models.Comment{Content: "test"}

		repo.EXPECT().
			Create(ctx, com).
			Return(nil)

		err := svc.Create(ctx, com)
		require.NoError(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		svc, repo, ctx := newTestService(t)

		com := &models.Comment{Content: "test"}
		expErr := errors.New("db error")

		repo.EXPECT().
			Create(ctx, com).
			Return(expErr)

		err := svc.Create(ctx, com)
		require.ErrorIs(t, err, expErr)
	})
}

func TestCommentsService_Update(t *testing.T) {
	svc, repo, ctx := newTestService(t)

	com := &models.Comment{ID: 1, Content: "updated"}

	repo.EXPECT().
		Update(ctx, com).
		Return(nil)

	err := svc.Update(ctx, com)
	require.NoError(t, err)
}

func TestCommentsService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, repo, ctx := newTestService(t)

		repo.EXPECT().
			Delete(ctx, int64(1)).
			Return(nil)

		err := svc.Delete(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		svc, repo, ctx := newTestService(t)

		expErr := errors.New("delete failed")

		repo.EXPECT().
			Delete(ctx, int64(1)).
			Return(expErr)

		err := svc.Delete(ctx, 1)
		require.ErrorIs(t, err, expErr)
	})
}

func TestCommentsService_GetByParent(t *testing.T) {
	svc, repo, ctx := newTestService(t)

	parentID := int64(10)
	limit, offset := int64(20), int64(0)

	expected := []*models.Comment{
		{ID: 1, ParentID: &parentID, Content: "c1"},
		{ID: 2, ParentID: &parentID, Content: "c2"},
	}

	repo.EXPECT().
		GetByParent(ctx, &parentID, limit, offset).
		Return(expected, nil)

	res, err := svc.GetByParent(ctx, &parentID, limit, offset)
	require.NoError(t, err)
	require.Equal(t, expected, res)
}

func TestCommentsService_Search(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, repo, ctx := newTestService(t)

		query := "hello"
		limit, offset := int64(10), int64(0)

		expected := []*models.Comment{
			{ID: 1, Content: "hello world"},
		}

		repo.EXPECT().
			Search(ctx, query, limit, offset).
			Return(expected, nil)

		res, err := svc.Search(ctx, query, limit, offset)
		require.NoError(t, err)
		require.Equal(t, expected, res)
	})

	t.Run("repo error", func(t *testing.T) {
		svc, repo, ctx := newTestService(t)

		expErr := errors.New("search failed")

		repo.EXPECT().
			Search(ctx, "q", int64(5), int64(0)).
			Return(nil, expErr)

		res, err := svc.Search(ctx, "q", 5, 0)
		require.Nil(t, res)
		require.ErrorIs(t, err, expErr)
	})
}

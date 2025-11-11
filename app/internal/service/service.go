package service

import (
	"comment-tree/internal/models"
	"context"

	"github.com/wb-go/wbf/zlog"
)

type CommentsRepository interface {
	Create(ctx context.Context, com *models.Comment) error
	Update(ctx context.Context, com *models.Comment) error
	Delete(ctx context.Context, id int64) error
	GetByParent(ctx context.Context, parentID *int64, limit, offset int64) ([]*models.Comment, error)
}

type CommentsService struct {
	repo CommentsRepository
	log  *zlog.Zerolog
}

func NewCommentsService(repo CommentsRepository, log *zlog.Zerolog) *CommentsService {
	return &CommentsService{
		repo: repo,
		log:  log,
	}
}

func (s *CommentsService) Create(ctx context.Context, com *models.Comment) error {
	return s.repo.Create(ctx, com)
}

func (s *CommentsService) Update(ctx context.Context, com *models.Comment) error {
	return s.repo.Update(ctx, com)
}

func (s *CommentsService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *CommentsService) GetByParent(ctx context.Context, parentID *int64, limit, offset int64) ([]*models.Comment, error) {
	return s.repo.GetByParent(ctx, parentID, limit, offset)
}

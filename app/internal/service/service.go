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
	Search(ctx context.Context, query string, limit, offset int64) ([]*models.Comment, error)
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
	if err := s.repo.Create(ctx, com); err != nil {
		s.log.Error().
			Err(err).
			Msg("failed to create comment")
		return err
	}
	return nil
}

func (s *CommentsService) Update(ctx context.Context, com *models.Comment) error {
	if err := s.repo.Update(ctx, com); err != nil {
		s.log.Error().
			Err(err).
			Msg("failed to update comment")
		return err
	}
	return nil
}

func (s *CommentsService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error().
			Err(err).
			Msg("failed to delete comment")
		return err
	}
	return nil
}

func (s *CommentsService) GetByParent(ctx context.Context, parentID *int64, limit, offset int64) ([]*models.Comment, error) {
	coms, err := s.repo.GetByParent(ctx, parentID, limit, offset)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("failed to get comments")
		return nil, err
	}
	return coms, nil
}

func (s *CommentsService) Search(ctx context.Context, query string, limit, offset int64) ([]*models.Comment, error) {
	coms, err := s.repo.Search(ctx, query, limit, offset)
	if err != nil {
		s.log.Error().
			Err(err).
			Msg("failed to search comments")
		return nil, err
	}
	return coms, nil
}

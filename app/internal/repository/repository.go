package repository

import (
	"context"

	"comment-tree/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
)

type CommentRepository struct {
	db       *dbpg.DB
	strategy retry.Strategy
	sb       squirrel.StatementBuilderType
}

func NewCommentRepository(db *dbpg.DB, strategy retry.Strategy) *CommentRepository {
	return &CommentRepository{
		db:       db,
		strategy: strategy,
		sb:       squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *CommentRepository) Create(ctx context.Context, com *models.Comment) error {
	if com == nil {
		return ErrNilValue
	}

	query := r.sb.Insert("comments").
		Columns(
			"parent_id", "content", "created_at",
		).Values(
		com.ParentID, com.Content, com.CreatedAt,
	).Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	row, err := r.db.QueryRowWithRetry(ctx, r.strategy, sql, args...)
	if err != nil {
		return wrapDBError(err)
	}

	return wrapDBError(
		row.Scan(&com.ID),
	)
}

func (r *CommentRepository) Update(ctx context.Context, com *models.Comment) error {
	if com == nil {
		return ErrNilValue
	}

	query := r.sb.Update("comments").
		Set("content", com.Content).
		Where(squirrel.Eq{"id": com.ID})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecWithRetry(ctx, r.strategy, sql, args...)

	return wrapDBError(err)
}

func (r *CommentRepository) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return ErrNilValue
	}

	query := r.sb.Delete("comments").
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecWithRetry(ctx, r.strategy, sql, args...)

	return wrapDBError(err)
}

func (r *CommentRepository) GetByParent(ctx context.Context, parentID *int64, limit, offset int64) ([]*models.Comment, error) {
	query := r.sb.
		Select("id", "parent_id", "content", "created_at").
		From("comments").
		OrderBy("created_at").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	if parentID == nil {
		// parent_id IS NULL
		query = query.Where("parent_id IS NULL")
	} else {
		// parent_id = $1
		query = query.Where(squirrel.Eq{"parent_id": *parentID})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryWithRetry(ctx, r.strategy, sql, args...)
	if err != nil {
		return nil, wrapDBError(err)
	}
	defer rows.Close()

	var result []*models.Comment
	for rows.Next() {
		c := &models.Comment{}
		if err := rows.Scan(&c.ID, &c.ParentID, &c.Content, &c.CreatedAt); err != nil {
			return nil, wrapDBError(err)
		}
		result = append(result, c)
	}

	return result, nil
}

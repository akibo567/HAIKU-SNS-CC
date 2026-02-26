package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HaikuPost struct {
	ID          string
	UserID      string
	Username    string
	DisplayName string
	Ku1         string
	Ku2         string
	Ku3         string
	LikeCount   int
	CreatedAt   string
}

type HaikuRepository struct {
	pool *pgxpool.Pool
}

func NewHaikuRepository(pool *pgxpool.Pool) *HaikuRepository {
	return &HaikuRepository{pool: pool}
}

func (r *HaikuRepository) Create(ctx context.Context, userID, ku1, ku2, ku3 string) (*HaikuPost, error) {
	var p HaikuPost
	err := r.pool.QueryRow(ctx,
		`INSERT INTO haiku_posts (user_id, ku1, ku2, ku3)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, ku1, ku2, ku3, like_count, created_at`,
		userID, ku1, ku2, ku3,
	).Scan(&p.ID, &p.UserID, &p.Ku1, &p.Ku2, &p.Ku3, &p.LikeCount, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create haiku: %w", err)
	}
	return &p, nil
}

func (r *HaikuRepository) List(ctx context.Context, cursor string, limit int) ([]HaikuPost, error) {
	var query string
	var args []any

	if cursor == "" {
		query = `
			SELECT p.id, p.user_id, u.username, u.display_name, p.ku1, p.ku2, p.ku3, p.like_count, p.created_at
			FROM haiku_posts p
			JOIN users u ON u.id = p.user_id
			ORDER BY p.created_at DESC
			LIMIT $1`
		args = []any{limit}
	} else {
		query = `
			SELECT p.id, p.user_id, u.username, u.display_name, p.ku1, p.ku2, p.ku3, p.like_count, p.created_at
			FROM haiku_posts p
			JOIN users u ON u.id = p.user_id
			WHERE p.created_at < (SELECT created_at FROM haiku_posts WHERE id = $1)
			ORDER BY p.created_at DESC
			LIMIT $2`
		args = []any{cursor, limit}
	}

	pgxRows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list haiku: %w", err)
	}
	defer pgxRows.Close()

	var posts []HaikuPost
	for pgxRows.Next() {
		var p HaikuPost
		if err := pgxRows.Scan(&p.ID, &p.UserID, &p.Username, &p.DisplayName, &p.Ku1, &p.Ku2, &p.Ku3, &p.LikeCount, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *HaikuRepository) FindByID(ctx context.Context, id string) (*HaikuPost, error) {
	var p HaikuPost
	err := r.pool.QueryRow(ctx,
		`SELECT p.id, p.user_id, u.username, u.display_name, p.ku1, p.ku2, p.ku3, p.like_count, p.created_at
		 FROM haiku_posts p
		 JOIN users u ON u.id = p.user_id
		 WHERE p.id = $1`,
		id,
	).Scan(&p.ID, &p.UserID, &p.Username, &p.DisplayName, &p.Ku1, &p.Ku2, &p.Ku3, &p.LikeCount, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find haiku by id: %w", err)
	}
	return &p, nil
}

func (r *HaikuRepository) ListByUserID(ctx context.Context, userID string, cursor string, limit int) ([]HaikuPost, error) {
	var query string
	var args []any

	if cursor == "" {
		query = `
			SELECT p.id, p.user_id, u.username, u.display_name, p.ku1, p.ku2, p.ku3, p.like_count, p.created_at
			FROM haiku_posts p
			JOIN users u ON u.id = p.user_id
			WHERE p.user_id = $1
			ORDER BY p.created_at DESC
			LIMIT $2`
		args = []any{userID, limit}
	} else {
		query = `
			SELECT p.id, p.user_id, u.username, u.display_name, p.ku1, p.ku2, p.ku3, p.like_count, p.created_at
			FROM haiku_posts p
			JOIN users u ON u.id = p.user_id
			WHERE p.user_id = $1
			  AND p.created_at < (SELECT created_at FROM haiku_posts WHERE id = $2)
			ORDER BY p.created_at DESC
			LIMIT $3`
		args = []any{userID, cursor, limit}
	}

	pgxRows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list haiku by user: %w", err)
	}
	defer pgxRows.Close()

	var posts []HaikuPost
	for pgxRows.Next() {
		var p HaikuPost
		if err := pgxRows.Scan(&p.ID, &p.UserID, &p.Username, &p.DisplayName, &p.Ku1, &p.Ku2, &p.Ku3, &p.LikeCount, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *HaikuRepository) Delete(ctx context.Context, id, userID string) (bool, error) {
	result, err := r.pool.Exec(ctx,
		`DELETE FROM haiku_posts WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	if err != nil {
		return false, fmt.Errorf("delete haiku: %w", err)
	}
	return result.RowsAffected() > 0, nil
}

func (r *HaikuRepository) AddLike(ctx context.Context, userID, postID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, postID,
	)
	if err != nil {
		return fmt.Errorf("add like: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE haiku_posts SET like_count = like_count + 1 WHERE id = $1`,
		postID,
	)
	if err != nil {
		return fmt.Errorf("update like count: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *HaikuRepository) RemoveLike(ctx context.Context, userID, postID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx,
		`DELETE FROM likes WHERE user_id = $1 AND post_id = $2`,
		userID, postID,
	)
	if err != nil {
		return fmt.Errorf("remove like: %w", err)
	}

	if result.RowsAffected() > 0 {
		_, err = tx.Exec(ctx,
			`UPDATE haiku_posts SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1`,
			postID,
		)
		if err != nil {
			return fmt.Errorf("update like count: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *HaikuRepository) IsLikedByUser(ctx context.Context, userID, postID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2)`,
		userID, postID,
	).Scan(&exists)
	return exists, err
}

func (r *HaikuRepository) LikedPostIDs(ctx context.Context, userID string, postIDs []string) (map[string]bool, error) {
	if len(postIDs) == 0 {
		return map[string]bool{}, nil
	}

	rows, err := r.pool.Query(ctx,
		`SELECT post_id FROM likes WHERE user_id = $1 AND post_id = ANY($2)`,
		userID, postIDs,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]bool)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result[id] = true
	}
	return result, nil
}

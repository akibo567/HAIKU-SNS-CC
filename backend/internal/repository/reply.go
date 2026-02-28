package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Reply struct {
	ID          string
	PostID      string
	UserID      string
	Username    string
	DisplayName string
	Ku1         string
	Ku2         string
	Ku3         string
	CreatedAt   time.Time
}

type ReplyRepository struct {
	pool *pgxpool.Pool
}

func NewReplyRepository(pool *pgxpool.Pool) *ReplyRepository {
	return &ReplyRepository{pool: pool}
}

func (r *ReplyRepository) Create(ctx context.Context, postID, userID, ku1, ku2, ku3 string) (*Reply, error) {
	var reply Reply
	err := r.pool.QueryRow(ctx,
		`INSERT INTO replies (post_id, user_id, ku1, ku2, ku3)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, post_id, user_id, ku1, ku2, ku3, created_at`,
		postID, userID, ku1, ku2, ku3,
	).Scan(&reply.ID, &reply.PostID, &reply.UserID, &reply.Ku1, &reply.Ku2, &reply.Ku3, &reply.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create reply: %w", err)
	}
	return &reply, nil
}

func (r *ReplyRepository) ListByPostID(ctx context.Context, postID string) ([]Reply, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT r.id, r.post_id, r.user_id, u.username, u.display_name,
		        r.ku1, r.ku2, r.ku3, r.created_at
		 FROM replies r
		 JOIN users u ON u.id = r.user_id
		 WHERE r.post_id = $1
		 ORDER BY r.created_at ASC`,
		postID,
	)
	if err != nil {
		return nil, fmt.Errorf("list replies: %w", err)
	}
	defer rows.Close()

	var replies []Reply
	for rows.Next() {
		var reply Reply
		if err := rows.Scan(
			&reply.ID, &reply.PostID, &reply.UserID,
			&reply.Username, &reply.DisplayName,
			&reply.Ku1, &reply.Ku2, &reply.Ku3,
			&reply.CreatedAt,
		); err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}
	return replies, nil
}

package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntisrangga142/chat/internals/models"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// Ambil semua notifikasi unread untuk user (owner akun)
func (r *NotificationRepository) GetUnreadNotifications(ctx context.Context, userID int) (models.NotificationList, error) {
	query := `
		-- FOLLOW notifications
		SELECT 'follow' AS type,
			   f.follower_id AS from_id,
			   p.fullname AS from_name,
			   NULL AS post_id,
			   p.fullname || ' followed you' AS message,
			   f.created_at
		FROM followers f
		JOIN profiles p ON p.id = f.follower_id
		WHERE f.account_id = $1 AND (f.read = false OR f.read IS NULL)

		UNION ALL

		-- LIKE notifications
		SELECT 'like' AS type,
			   l.account_id AS from_id,
			   p.fullname AS from_name,
			   l.post_id AS post_id,
			   p.fullname || ' liked your post' AS message,
			   l.created_at
		FROM likes l
		JOIN posts ps ON ps.id = l.post_id
		JOIN profiles p ON p.id = l.account_id
		WHERE ps.account_id = $1 AND (l.read = false OR l.read IS NULL)

		UNION ALL

		-- COMMENT notifications
		SELECT 'comment' AS type,
			   c.account_id AS from_id,
			   p.fullname AS from_name,
			   c.post_id AS post_id,
			   p.fullname || ' commented: ' || c.comment AS message,
			   c.created_at
		FROM comments c
		JOIN posts ps ON ps.id = c.post_id
		JOIN profiles p ON p.id = c.account_id
		WHERE ps.account_id = $1 AND (c.read = false OR c.read IS NULL)

		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications models.NotificationList
	for rows.Next() {
		var n models.Notification
		var postID *int
		if err := rows.Scan(&n.Type, &n.FromID, &n.FromName, &postID, &n.Message, &n.CreatedAt); err != nil {
			return nil, err
		}
		if postID != nil {
			n.PostID = postID
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

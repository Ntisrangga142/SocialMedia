package models

import "time"

type Like struct {
	ID        int        `json:"id"`
	AccountID int        `json:"account_id"`
	PostID    int        `json:"post_id"`
	Read      bool       `json:"read"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CreateLikeRequest struct {
	PostID int `json:"post_id"`
}

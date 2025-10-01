package models

import "time"

type Notification struct {
	Type      string    `json:"type"`              // follow, like, comment
	FromID    int       `json:"from_id"`           // id user yang melakukan aksi
	FromName  string    `json:"from_name"`         // nama user yang melakukan aksi
	PostID    *int      `json:"post_id,omitempty"` // kalau like/comment, ada post_id
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationList []Notification

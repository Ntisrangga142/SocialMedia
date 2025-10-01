package models

import "time"

type Post struct {
	ID        int        `json:"id"`
	AccountID int        `json:"account_id"`
	Caption   string     `json:"caption"`
	Images    []PostImg  `json:"images,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PostImg struct {
	ID        int        `json:"id"`
	PostID    int        `json:"post_id"`
	Img       string     `json:"img"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CreatePostRequest struct {
	Caption string `form:"caption" json:"caption"`
	Images  []string
}

// Post Feed
type PostFeed struct {
	ID           int      `json:"id"`
	Fullname     string   `json:"fullname"`
	Caption      string   `json:"caption"`
	Images       []string `json:"images"`
	LikeCount    int      `json:"like_count"`
	CommentCount int      `json:"comment_count"`
}

// Post Detail
type PostDetail struct {
	ID        int              `json:"id"`
	Caption   string           `json:"caption"`
	CreatedAt time.Time        `json:"created_at"`
	Author    AuthorProfile    `json:"author"`
	Images    []string         `json:"images"`
	Likes     int              `json:"likes"`
	Comments  []CommentPreview `json:"comments"`
}

// Info Author
type AuthorProfile struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Img      string `json:"img"`
}

// Preview Comment
type CommentPreview struct {
	ID        int       `json:"id"`
	AccountID int       `json:"account_id"`
	Fullname  string    `json:"fullname"`
	Img       string    `json:"img"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

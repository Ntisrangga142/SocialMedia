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
	ID           int      `json:"id" example:"101"`
	Fullname     string   `json:"fullname" example:"Rangga Putra"`
	Caption      string   `json:"caption" example:"Liburan di pantai bareng teman-teman!"`
	Images       []string `json:"images" example:"['public/post/1.jpg','public/post/2.jpg']"`
	LikeCount    int      `json:"like_count" example:"123"`
	CommentCount int      `json:"comment_count" example:"45"`
}

// Post Detail
type AuthorProfile struct {
	ID       int    `json:"id" example:"1"`
	Fullname string `json:"fullname" example:"Rangga Putra"`
	Img      string `json:"img" example:"public/profile/rangga.jpg"`
}

type CommentPreview struct {
	ID        int       `json:"id" example:"501"`
	Fullname  string    `json:"fullname" example:"Siti Amelia"`
	Comment   string    `json:"comment" example:"Keren banget fotonya!"`
	CreatedAt time.Time `json:"created_at" example:"2025-09-20T14:30:00Z"`
}

type PostDetail struct {
	ID        int              `json:"id" example:"101"`
	Caption   string           `json:"caption" example:"Liburan di pantai bareng teman-teman!"`
	CreatedAt time.Time        `json:"created_at" example:"2025-09-20T12:00:00Z"`
	Author    AuthorProfile    `json:"author"`
	Images    []string         `json:"images" example:"['public/post/1.jpg','public/post/2.jpg']"`
	Likes     int              `json:"likes" example:"123"`
	Comments  []CommentPreview `json:"comments"`
}

package models

type Comment struct {
	ID        int    `json:"id"`
	AccountID int    `json:"account_id"`
	PostID    int    `json:"post_id"`
	Comment   string `json:"comment"`
}

type CreateCommentRequest struct {
	PostID  int    `json:"post_id"`
	Comment string `json:"comment"`
}

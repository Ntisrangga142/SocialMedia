package models

type ResponsePostList struct {
	Success bool       `json:"success" example:"true"`
	Message string     `json:"message" example:"Success Get Post Followings"`
	Data    []PostFeed `json:"data"`
}

type ResponsePostDetail struct {
	Success bool       `json:"success" example:"true"`
	Message string     `json:"message" example:"Success Get Post Detail"`
	Data    PostDetail `json:"data"`
}

type ResponseCreatePost struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success Created Post"`
	Data    Post   `json:"data"`
}

type ResponseGetComment struct {
	Success bool      `json:"success" example:"true"`
	Message string    `json:"message" example:"Succes Get Comment by Id Post"`
	Data    []Comment `json:"data"`
}

type ResponseAny struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Succes Get Comment by Id Post"`
	Data    any    `json:"data"`
}

package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ntisrangga142/chat/internals/models"
	"github.com/ntisrangga142/chat/internals/repositories"
	"github.com/ntisrangga142/chat/internals/utils"
	"github.com/redis/go-redis/v9"
)

type PostHandler struct {
	repo *repositories.PostRepository
	rdb  *redis.Client
}

func NewPostHandler(repo *repositories.PostRepository, rdb *redis.Client) *PostHandler {
	return &PostHandler{repo: repo, rdb: rdb}
}

// GetFollowingPosts godoc
// @Summary Get Following Posts
// @Description Get posts from accounts that the user follows
// @Tags Posts
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.ResponsePostList
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/following [get]
func (h *PostHandler) GetFollowingPosts(ctx *gin.Context) {
	uid, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "invalid token", err)
		return
	}

	posts, err := h.repo.GetFollowingPosts(ctx, uid)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Error", "failed to get posts", err)
		return
	}

	ctx.JSON(http.StatusOK, models.Response[any]{
		Success: true,
		Message: "Success Get Post Followings",
		Data:    posts,
	})
}

// GetPostDetail godoc
// @Summary Get Post Detail
// @Description Get detail of a single post including author, images, like count, and top comments
// @Tags Posts
// @Security BearerAuth
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.ResponsePostDetail
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id} [get]
func (h *PostHandler) GetPostDetail(ctx *gin.Context) {
	postID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid ID", "post id must be number", err)
		return
	}

	var cachedData models.PostDetail
	var redisKey = fmt.Sprintf("Chat-PostDetail-%d", postID)
	if err := utils.CacheHit(ctx.Request.Context(), h.rdb, redisKey, &cachedData); err == nil {
		ctx.JSON(http.StatusOK, models.Response[any]{
			Success: true,
			Message: "Success Get Profile User (from cache)",
			Data:    cachedData,
		})
		return
	}

	post, err := h.repo.GetPostDetail(ctx.Request.Context(), postID)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed", "cannot get post detail", err)
		return
	}

	if err := utils.RenewCache(ctx.Request.Context(), h.rdb, redisKey, post, 10); err != nil {
		log.Println("Failed to set redis cache:", err)
	}

	ctx.JSON(http.StatusOK, models.Response[any]{
		Success: true,
		Message: fmt.Sprintf("Success Get Post Detail - %d", postID),
		Data:    post,
	})
}

// CreatePost godoc
// @Summary Create Post
// @Description Create a new post with caption and images (form-data)
// @Tags Posts
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param caption formData string true "Post Caption"
// @Param images formData []file true "Post Images"
// @Success 201 {object} models.ResponseCreatePost
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts [post]
func (h *PostHandler) CreatePost(ctx *gin.Context) {
	uid, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "invalid token", err)
		return
	}

	// Ambil caption
	caption := ctx.PostForm("caption")

	// Ambil file images
	form, err := ctx.MultipartForm()
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "failed to parse form-data", err)
		return
	}

	files := form.File["images"]
	var imagePaths []string

	for _, file := range files {
		// Simpan file ke folder uploads
		path := "public/post/" + file.Filename
		if err := ctx.SaveUploadedFile(file, path); err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Internal Error", "failed to save file", err)
			return
		}
		imagePaths = append(imagePaths, path)
	}

	req := models.CreatePostRequest{
		Caption: caption,
		Images:  imagePaths,
	}

	post, err := h.repo.CreatePost(ctx, req, uid)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Error", "failed to create post", err)
		return
	}

	ctx.JSON(http.StatusCreated, models.Response[any]{
		Success: true,
		Message: "Success Created Post",
		Data:    post,
	})
}

// LikePost godoc
// @Summary Like a Post
// @Description Like a post by ID
// @Tags Posts
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Success 204
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id}/like [post]
func (h *PostHandler) LikePost(ctx *gin.Context) {
	uid, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "invalid token", err)
		return
	}

	postID, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.repo.DeleteLike(ctx, uid, postID); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed to unlike post", err)
		return
	}

	if err := h.repo.CreateLike(ctx, uid, postID); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed to like post", err)
		return
	}

	var redisKey = fmt.Sprintf("Chat-PostDetail-%d", postID)
	if err := utils.InvalidateCache(ctx, h.rdb, redisKey); err != nil {
		log.Println("Failed invalidate cache:", err)
	}

	ctx.Status(http.StatusNoContent)
}

// UnlikePost godoc
// @Summary Unlike a Post
// @Description Remove like from a post by ID
// @Tags Posts
// @Security BearerAuth
// @Param id path int true "Post ID"
// @Success 204
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id}/like [delete]
func (h *PostHandler) UnlikePost(ctx *gin.Context) {
	uid, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "invalid token", err)
		return
	}

	postID, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.repo.DeleteLike(ctx, uid, postID); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Sever Error", "failed to unlike post", err)
		return
	}

	var redisKey = fmt.Sprintf("Chat-PostDetail-%d", postID)
	if err := utils.InvalidateCache(ctx, h.rdb, redisKey); err != nil {
		log.Println("Failed invalidate cache:", err)
	}

	ctx.Status(http.StatusNoContent)
}

// CreateComment godoc
// @Summary Create Comment
// @Description Create a comment for a post
// @Tags Posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.CreateCommentRequest true "Comment Request"
// @Success 201
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id}/comments [post]
func (h *PostHandler) CreateComment(ctx *gin.Context) {
	uid, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "invalid token", err)
		return
	}

	var req models.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "invalid request", err)
		return
	}

	if err := h.repo.CreateComment(ctx, uid, req); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed to create comment", err)
		return
	}

	var redisKey = fmt.Sprintf("Chat-PostDetail-%d", req.PostID)
	if err := utils.InvalidateCache(ctx, h.rdb, redisKey); err != nil {
		log.Println("Failed invalidate cache:", err)
	}

	ctx.Status(http.StatusCreated)
}

// GetAllCommentsByPost godoc
// @Summary Get All Comments by Post
// @Description Get all comments from a post by ID
// @Tags Posts
// @Security BearerAuth
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.ResponseGetComment
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id}/comments [get]
func (h *PostHandler) GetAllCommentsByPost(ctx *gin.Context) {
	postID, _ := strconv.Atoi(ctx.Param("id"))
	comments, err := h.repo.GetAllCommentsByPost(ctx, postID)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed to fetch comments", err)
		return
	}
	ctx.JSON(http.StatusOK, models.Response[any]{
		Success: true,
		Message: fmt.Sprintf("Succes Get Comment by Id Post : %d", postID),
		Data:    comments,
	})
}

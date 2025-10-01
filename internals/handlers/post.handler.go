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

// Get Following Posts
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

// Get Post Detail
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

// Create Post
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

// Like Post
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

// Unlike Post
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

// Comment Post
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

// Get Comment Post
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

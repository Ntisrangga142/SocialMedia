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

type UserHandler struct {
	repo *repositories.UserRepository
	rdb  *redis.Client
}

func NewUserHandler(repo *repositories.UserRepository, rdb *redis.Client) *UserHandler {
	return &UserHandler{repo: repo, rdb: rdb}
}

// Get Profile
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	uid, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "unable to get user's token", err)
		return
	}

	var cachedData models.Profile
	var redisKey = fmt.Sprintf("Chat-Profile-%d", uid)
	if err := utils.CacheHit(ctx.Request.Context(), h.rdb, redisKey, &cachedData); err == nil {
		ctx.JSON(http.StatusOK, models.Response[models.Profile]{
			Success: true,
			Message: "Success Get Profile User (from cache)",
			Data:    cachedData,
		})
		return
	}

	profile, err := h.repo.GetProfile(ctx.Request.Context(), uid)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "unable get profile user", err)
		// return
	}

	if err := utils.RenewCache(ctx.Request.Context(), h.rdb, redisKey, profile, 10); err != nil {
		log.Println("Failed to set redis cache:", err)
	}

	ctx.JSON(http.StatusOK, models.Response[models.Profile]{
		Success: true,
		Message: "Success Get Profile User",
		Data:    *profile,
	})
}

// Update Profile
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	uid, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "invalid token", err)
		return
	}

	updates := make(map[string]any)

	// ambil field dari form-data (jika ada)
	if fullname := ctx.PostForm("fullname"); fullname != "" {
		updates["fullname"] = fullname
	}
	if phone := ctx.PostForm("phone"); phone != "" {
		updates["phone"] = phone
	}

	// Upload image jika ada
	file, err := ctx.FormFile("img")
	if err == nil {
		destDir := "public/profile"
		filename := fmt.Sprintf("profile_%d", uid)

		path, saveErr := utils.SaveUploadedFile(ctx, file, destDir, filename)
		if saveErr != nil {
			utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "Upload Failed", saveErr)
			return
		}

		updates["img"] = path
	}

	// Update ke DB
	if err := h.repo.UpdateProfile(ctx.Request.Context(), uid, updates); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed to update profile", err)
		return
	}

	var redisKey = fmt.Sprintf("Chat-Profile-%d", uid)
	if err := utils.InvalidateCache(ctx, h.rdb, redisKey); err != nil {
		log.Println("Failed invalidate cache:", err)
	}

	ctx.JSON(http.StatusOK, models.Response[any]{
		Success: true,
		Message: "Profile updated successfully",
	})
}

// Follow user
func (h *UserHandler) Follow(ctx *gin.Context) {
	uid, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "invalid token", err)
		return
	}

	targetID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "invalid id", err)
		return
	}

	if targetID == uid {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "cannot follow yourself", nil)
		return
	}

	if err := h.repo.Follow(ctx, targetID, uid); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Error", "failed to follow user", err)
		return
	}

	ctx.JSON(http.StatusCreated, models.Response[any]{
		Success: true,
		Message: "Success Followed",
	})
}

// Unfollow user
func (h *UserHandler) Unfollow(ctx *gin.Context) {
	uid, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "invalid token", err)
		return
	}

	targetID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "invalid id", err)
		return
	}

	if err := h.repo.Unfollow(ctx, targetID, uid); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Error", "failed to unfollow user", err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Get Followers (daftar pengikut saya)
func (h *UserHandler) GetFollowers(ctx *gin.Context) {
	uid, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "invalid token", err)
		return
	}

	followers, err := h.repo.GetFollowers(ctx, uid)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Error", "failed to get followers", err)
		return
	}

	ctx.JSON(http.StatusOK, models.Response[any]{
		Success: true,
		Message: "Success Get Followers",
		Data:    followers,
	})
}

// Get Following (daftar siapa saja yang saya ikuti)
func (h *UserHandler) GetFollowing(ctx *gin.Context) {
	uid, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "invalid token", err)
		return
	}

	followings, err := h.repo.GetFollowing(ctx, uid)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Error", "failed to get following", err)
		return
	}

	ctx.JSON(http.StatusOK, models.Response[any]{
		Success: true,
		Message: "Success Get Followings",
		Data:    followings,
	})
}

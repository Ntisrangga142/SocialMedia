package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ntisrangga142/chat/internals/models"
	"github.com/ntisrangga142/chat/internals/repositories"
	"github.com/ntisrangga142/chat/internals/utils"
	"github.com/ntisrangga142/chat/pkg"
	"github.com/redis/go-redis/v9"
)

type AuthHandler struct {
	repo *repositories.Auth
	rdb  *redis.Client
}

func NewAuthHandler(repo *repositories.Auth, rdb *redis.Client) *AuthHandler {
	return &AuthHandler{repo: repo, rdb: rdb}
}

// Register godoc
// @Summary Register a new account
// @Description Create new account with email & password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.AuthRequest true "Register Request"
// @Success 201 {object} models.ErrorResponse "Register successful"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req models.AuthRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "failed binding data", err)
		return
	}

	// Validasi email
	if valid := utils.ValidateEmail(req.Email); !valid {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "invalid email", fmt.Errorf("invalid email"))
		return
	}

	// Validasi password
	if valid := utils.ValidatePassword(req.Password); !valid {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "invalid password", fmt.Errorf("invalid password"))
		return
	}

	// Hash Password
	hashConfig := pkg.NewHashConfig()
	hashConfig.UseRecommended()
	hashedPassword, err := hashConfig.GenHash(req.Password)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed hashed password", err)
		return
	}

	// Repository Register
	if err := h.repo.Register(ctx, req.Email, hashedPassword); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "Email is already registered", err)
		return
	}

	ctx.JSON(http.StatusCreated, models.Response[any]{
		Success: true,
		Message: "Register account successful",
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email & password, return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.AuthRequest true "Login Request"
// @Success 200 {object} models.ResponseLogin "Login successful"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req models.AuthRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "failed binding data", err)
		return
	}

	// Cari akun
	userID, hashedPassword, err := h.repo.Login(ctx.Request.Context(), req.Email)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "user not found", err)
		return
	}
	if userID == 0 {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "user not found", err)
		return
	}

	// Verifikasi password
	hashConfig := pkg.NewHashConfig()
	match, err := hashConfig.ComparePasswordAndHash(req.Password, hashedPassword)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed compare password", err)
		return
	}
	if !match {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "invalid username or password", errors.New("invalid password"))
		return
	}

	// Generate JWT
	claims := pkg.NewJWTClaims(userID)
	token, err := claims.GenToken()
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed generate token", err)
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseLogin{
		Success: true,
		Message: "Login successful",
		Token:   token,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate JWT token by blacklisting it
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.ErrorResponse "Successfully logged out"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(ctx *gin.Context) {
	token, err := utils.GetToken(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "failed get token", err)
		return
	}

	expiresAt, err := utils.GetExpiredFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "failed get expired time token", err)
		return
	}

	expiresIn := time.Until(expiresAt)
	if expiresIn <= 0 {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "token already expired", err)
		return
	}

	redisKey := fmt.Sprintf("Blacklist:%s", token)
	if err := utils.RenewCache(ctx.Request.Context(), h.rdb, redisKey, token, expiresIn); err != nil {
		log.Println("Failed to set redis cache:", err)
	}

	ctx.JSON(http.StatusOK, models.Response[any]{
		Success: true,
		Message: "Successfully logged out",
	})
}

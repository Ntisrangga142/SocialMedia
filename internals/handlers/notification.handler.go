package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntisrangga142/chat/internals/models"
	"github.com/ntisrangga142/chat/internals/repositories"
	"github.com/ntisrangga142/chat/internals/utils"
)

type NotificationHandler struct {
	repo *repositories.NotificationRepository
}

func NewNotificationHandler(repo *repositories.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

// Get unread notifications
func (h *NotificationHandler) GetUnreadNotifications(ctx *gin.Context) {
	uid, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "invalid token", err)
		return
	}

	notifications, err := h.repo.GetUnreadNotifications(ctx, uid)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed to get notifications", err)
		return
	}

	ctx.JSON(http.StatusOK, models.Response[models.NotificationList]{
		Success: true,
		Message: "Success Get Unread Notifications",
		Data:    notifications,
	})
}

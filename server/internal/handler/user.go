package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wenyuwang1985/community/internal/middleware"
	"github.com/wenyuwang1985/community/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type updateProfileRequest struct {
	Nickname  string `json:"nickname" binding:"required,max=20"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,max=500"`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	userID := c.GetInt64("userID")
	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, req.Nickname, req.AvatarURL)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	middleware.Success(c, gin.H{
		"id":         user.ID,
		"nickname":   user.Nickname,
		"avatar_url": user.AvatarURL,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "用户 ID 格式错误")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		middleware.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	middleware.Success(c, gin.H{
		"id":           user.ID,
		"nickname":     user.Nickname,
		"avatar_url":   user.AvatarURL,
		"credit_score": user.CreditScore,
		"created_at":   user.CreatedAt,
	})
}

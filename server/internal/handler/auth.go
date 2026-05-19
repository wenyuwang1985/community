package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wenyuwang1985/community/internal/config"
	"github.com/wenyuwang1985/community/internal/middleware"
	"github.com/wenyuwang1985/community/internal/service"
	"github.com/wenyuwang1985/community/pkg/jwt"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtCfg      config.JWTConfig
}

func NewAuthHandler(authService *service.AuthService, jwtCfg config.JWTConfig) *AuthHandler {
	return &AuthHandler{authService: authService, jwtCfg: jwtCfg}
}

type registerRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, gin.H{
		"id":    user.ID,
		"phone": user.Phone,
	})
}

type loginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.Login(c.Request.Context(), req.Phone, req.Password)
	if err != nil {
		middleware.Error(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID, h.jwtCfg.Secret, h.jwtCfg.AccessExpire)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, 500, "生成 Token 失败")
		return
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID, h.jwtCfg.Secret, h.jwtCfg.RefreshExpire)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, 500, "生成 Token 失败")
		return
	}

	middleware.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	claims, err := jwt.ParseToken(req.RefreshToken, h.jwtCfg.Secret)
	if err != nil {
		middleware.Error(c, http.StatusUnauthorized, 401, "refresh_token 无效或已过期")
		return
	}

	if claims.TokenType != "refresh" {
		middleware.Error(c, http.StatusUnauthorized, 401, "token 类型错误")
		return
	}

	accessToken, err := jwt.GenerateAccessToken(claims.UserID, h.jwtCfg.Secret, h.jwtCfg.AccessExpire)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, 500, "生成 Token 失败")
		return
	}

	middleware.Success(c, gin.H{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID := c.GetInt64("userID")

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		middleware.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	// 手机号脱敏
	phone := user.Phone
	if len(phone) >= 7 {
		phone = phone[:3] + "****" + phone[len(phone)-4:]
	}

	middleware.Success(c, gin.H{
		"id":           user.ID,
		"phone":        phone,
		"nickname":     user.Nickname,
		"avatar_url":   user.AvatarURL,
		"credit_score": user.CreditScore,
		"created_at":   user.CreatedAt,
	})
}

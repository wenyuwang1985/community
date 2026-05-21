package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wenyuwang1985/community/internal/middleware"
	"github.com/wenyuwang1985/community/internal/service"
)

type CommunityHandler struct {
	communityService *service.CommunityService
	chatService      *service.ChatService
}

func NewCommunityHandler(communityService *service.CommunityService, chatService *service.ChatService) *CommunityHandler {
	return &CommunityHandler{communityService: communityService, chatService: chatService}
}

func (h *CommunityHandler) Search(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		middleware.Success(c, []interface{}{})
		return
	}

	communities, err := h.communityService.SearchCommunities(c.Request.Context(), keyword)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	middleware.Success(c, communities)
}

func (h *CommunityHandler) Subscribe(c *gin.Context) {
	userID := c.GetInt64("userID")
	communityIDStr := c.Param("id")
	communityID, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "社区 ID 格式错误")
		return
	}

	ctx := c.Request.Context()
	if err := h.communityService.Subscribe(ctx, userID, communityID); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	// 自动加入社区频道
	if h.chatService != nil {
		_ = h.chatService.JoinChannelConversation(ctx, userID, communityID)
	}

	middleware.Success(c, nil)
}

func (h *CommunityHandler) Unsubscribe(c *gin.Context) {
	userID := c.GetInt64("userID")
	communityIDStr := c.Param("id")
	communityID, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "社区 ID 格式错误")
		return
	}

	ctx := c.Request.Context()
	// 先离开社区频道
	if h.chatService != nil {
		_ = h.chatService.LeaveChannelConversation(ctx, userID, communityID)
	}

	if err := h.communityService.Unsubscribe(ctx, userID, communityID); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, nil)
}

func (h *CommunityHandler) GetUserSubscriptions(c *gin.Context) {
	userID := c.GetInt64("userID")
	list, err := h.communityService.GetUserSubscriptionsWithPrimary(c.Request.Context(), userID)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	middleware.Success(c, list)
}

func (h *CommunityHandler) SetPrimary(c *gin.Context) {
	userID := c.GetInt64("userID")
	communityIDStr := c.Param("id")
	communityID, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "社区 ID 格式错误")
		return
	}

	if err := h.communityService.SetPrimaryCommunity(c.Request.Context(), userID, communityID); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, nil)
}

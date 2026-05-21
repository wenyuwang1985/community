package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wenyuwang1985/community/internal/middleware"
	"github.com/wenyuwang1985/community/internal/service"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

// ListConversations 获取我的会话列表
func (h *ChatHandler) ListConversations(c *gin.Context) {
	userID := c.GetInt64("userID")
	list, err := h.chatService.ListUserConversations(c.Request.Context(), userID)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	middleware.Success(c, list)
}

// ListMessages 获取会话历史消息
func (h *ChatHandler) ListMessages(c *gin.Context) {
	userID := c.GetInt64("userID")
	convIDStr := c.Param("id")
	convID, err := strconv.ParseInt(convIDStr, 10, 64)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "会话 ID 格式错误")
		return
	}

	ctx := c.Request.Context()
	// 检查用户是否属于该会话
	isMember, err := h.chatService.IsConversationParticipant(ctx, convID, userID)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	if !isMember {
		middleware.Error(c, http.StatusForbidden, 403, "无权访问该会话")
		return
	}

	lastIDStr := c.Query("last_id")
	lastID, _ := strconv.ParseInt(lastIDStr, 10, 64)
	limitStr := c.Query("limit")
	limit, _ := strconv.Atoi(limitStr)

	messages, err := h.chatService.ListMessages(ctx, convID, lastID, limit)
	if err != nil {
		middleware.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	// 组装发送者信息
	var result []gin.H
	for _, msg := range messages {
		sender, _ := h.chatService.GetUserByID(ctx, msg.SenderID)
		senderInfo := gin.H{}
		if sender != nil {
			senderInfo = gin.H{
				"id":         sender.ID,
				"nickname":   sender.Nickname,
				"avatar_url": sender.AvatarURL,
			}
		}
		result = append(result, gin.H{
			"id":              msg.ID,
			"conversation_id": msg.ConversationID,
			"sender_id":       msg.SenderID,
			"content":         msg.Content,
			"type":            msg.Type,
			"created_at":      msg.CreatedAt,
			"sender":          senderInfo,
		})
	}

	middleware.Success(c, result)
}

type createPrivateConversationRequest struct {
	TargetUserID int64 `json:"target_user_id" binding:"required"`
}

// CreatePrivateConversation 发起/获取私聊会话
func (h *ChatHandler) CreatePrivateConversation(c *gin.Context) {
	var req createPrivateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error())
		return
	}

	userID := c.GetInt64("userID")
	convID, err := h.chatService.GetOrCreatePrivateConversation(c.Request.Context(), userID, req.TargetUserID)
	if err != nil {
		middleware.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	middleware.Success(c, gin.H{"conversation_id": convID})
}

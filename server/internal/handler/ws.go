package handler

import (
	"encoding/json"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wenyuwang1985/community/internal/service"
	"github.com/wenyuwang1985/community/internal/ws"
	"github.com/wenyuwang1985/community/pkg/jwt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 开发阶段允许所有跨域
		return true
	},
}

type WSMessage struct {
	Type           string `json:"type"`
	ConversationID int64  `json:"conversation_id"`
	Content        string `json:"content"`
}

type WSHandler struct {
	hub         *ws.Hub
	chatService *service.ChatService
	jwtSecret   string
}

func NewWSHandler(hub *ws.Hub, chatService *service.ChatService, jwtSecret string) *WSHandler {
	return &WSHandler{
		hub:         hub,
		chatService: chatService,
		jwtSecret:   jwtSecret,
	}
}

// ServeHTTP 处理 WebSocket 连接请求（Gin handler）
func (h *WSHandler) ServeHTTP(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "缺少 token"})
		return
	}

	claims, err := jwt.ParseToken(token, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token 无效"})
		return
	}
	if claims.TokenType != "access" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token 类型错误"})
		return
	}
	userID := claims.UserID

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)
		return
	}

	client := &ws.Client{
		Hub:    h.hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userID,
	}

	h.hub.Register(client)

	// 启动读写 goroutine
	go client.WritePump()
	go client.ReadPump(func(uid int64, data []byte) {
		h.handleMessage(uid, data)
	})
}

// handleMessage 处理客户端发来的消息
func (h *WSHandler) handleMessage(userID int64, data []byte) {
	var msg WSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		h.sendError(userID, "消息格式错误")
		return
	}

	switch msg.Type {
	case "ping":
		h.sendPong(userID)

	case "send_message":
		if msg.ConversationID == 0 || msg.Content == "" {
			h.sendError(userID, "conversation_id 和 content 不能为空")
			return
		}

		ctx := context.Background()
		// 检查用户是否属于该会话
		isMember, err := h.chatService.IsConversationParticipant(ctx, msg.ConversationID, userID)
		if err != nil {
			h.sendError(userID, "服务端错误")
			return
		}
		if !isMember {
			h.sendError(userID, "无权向该会话发送消息")
			return
		}

		// 保存消息
		savedMsg, err := h.chatService.SaveMessage(ctx, msg.ConversationID, userID, msg.Content, "text")
		if err != nil {
			h.sendError(userID, "发送消息失败")
			return
		}

		// 获取发送者信息
		sender, _ := h.chatService.GetUserByID(ctx, userID)
		senderInfo := map[string]interface{}{}
		if sender != nil {
			senderInfo = map[string]interface{}{
				"id":         sender.ID,
				"nickname":   sender.Nickname,
				"avatar_url": sender.AvatarURL,
			}
		}

		// 组装推送消息
		payload := map[string]interface{}{
			"type": "message",
			"data": map[string]interface{}{
				"id":              savedMsg.ID,
				"conversation_id": savedMsg.ConversationID,
				"sender_id":       savedMsg.SenderID,
				"content":         savedMsg.Content,
				"type":            savedMsg.Type,
				"created_at":      savedMsg.CreatedAt,
				"sender":          senderInfo,
			},
		}
		payloadBytes, _ := json.Marshal(payload)

		// 推送给会话所有在线参与者
		participants, err := h.chatService.GetConversationParticipants(ctx, msg.ConversationID)
		if err != nil {
			log.Printf("获取参与者失败: %v", err)
			return
		}

		for _, pid := range participants {
			h.hub.SendToUser(pid, payloadBytes)
		}

	default:
		h.sendError(userID, "未知消息类型: "+msg.Type)
	}
}

func (h *WSHandler) sendError(userID int64, errMsg string) {
	payload, _ := json.Marshal(map[string]interface{}{
		"type": "error",
		"msg":  errMsg,
	})
	h.hub.SendToUser(userID, payload)
}

func (h *WSHandler) sendPong(userID int64) {
	payload, _ := json.Marshal(map[string]interface{}{
		"type": "pong",
	})
	h.hub.SendToUser(userID, payload)
}

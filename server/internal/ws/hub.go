package ws

import (
	"sync"
)

// Hub 管理所有 WebSocket 连接
type Hub struct {
	clients map[int64]*Client // userID -> Client
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int64]*Client),
	}
}

// Register 注册客户端
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	// 如果该用户已有连接，先关闭旧连接
	if old, ok := h.clients[client.UserID]; ok {
		old.Close()
	}
	h.clients[client.UserID] = client
	h.mu.Unlock()
}

// Unregister 注销客户端
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	if c, ok := h.clients[client.UserID]; ok && c == client {
		delete(h.clients, client.UserID)
	}
	h.mu.Unlock()
}

// SendToUser 向指定用户发送消息
func (h *Hub) SendToUser(userID int64, message []byte) bool {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()
	if !ok {
		return false
	}
	select {
	case client.Send <- message:
		return true
	default:
		// 发送缓冲区满，关闭连接
		client.Close()
		return false
	}
}

// IsOnline 检查用户是否在线
func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	_, ok := h.clients[userID]
	h.mu.RUnlock()
	return ok
}

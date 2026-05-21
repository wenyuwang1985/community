package model

import "time"

type Conversation struct {
	ID          int64     `json:"id"`
	Type        string    `json:"type"`
	CommunityID *int64    `json:"community_id,omitempty"`
	Name        string    `json:"name"`
	CreatedBy   *int64    `json:"created_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type ConversationParticipant struct {
	ConversationID int64     `json:"conversation_id"`
	UserID         int64     `json:"user_id"`
	JoinedAt       time.Time `json:"joined_at"`
}

type ChatMessage struct {
	ID             int64     `json:"id"`
	ConversationID int64     `json:"conversation_id"`
	SenderID       int64     `json:"sender_id"`
	Content        string    `json:"content"`
	Type           string    `json:"type"`
	CreatedAt      time.Time `json:"created_at"`
}

package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wenyuwang1985/community/internal/model"
	"github.com/wenyuwang1985/community/pkg/snowflake"
)

type ChatService struct {
	db *pgxpool.Pool
}

func NewChatService(db *pgxpool.Pool) *ChatService {
	return &ChatService{db: db}
}

// GetOrCreatePrivateConversation 获取或创建私聊会话
func (s *ChatService) GetOrCreatePrivateConversation(ctx context.Context, userA, userB int64) (int64, error) {
	if userA == userB {
		return 0, errors.New("不能与自己建立私聊")
	}

	// 查找是否已存在私聊会话
	var convID int64
	err := s.db.QueryRow(ctx,
		`SELECT c.id FROM conversations c
		 JOIN conversation_participants p1 ON c.id = p1.conversation_id AND p1.user_id = $1
		 JOIN conversation_participants p2 ON c.id = p2.conversation_id AND p2.user_id = $2
		 WHERE c.type = 'private'`,
		userA, userB).Scan(&convID)
	if err == nil {
		return convID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("查询私聊会话失败: %w", err)
	}

	// 创建新会话
	convID = snowflake.Generate()
	_, err = s.db.Exec(ctx,
		"INSERT INTO conversations (id, type, name, created_by) VALUES ($1, 'private', '', $2)",
		convID, userA)
	if err != nil {
		return 0, fmt.Errorf("创建私聊会话失败: %w", err)
	}

	// 添加参与者
	_, err = s.db.Exec(ctx,
		"INSERT INTO conversation_participants (conversation_id, user_id) VALUES ($1, $2), ($1, $3)",
		convID, userA, userB)
	if err != nil {
		return 0, fmt.Errorf("添加会话参与者失败: %w", err)
	}

	return convID, nil
}

// GetOrCreateChannelConversation 获取或创建社区频道会话
func (s *ChatService) GetOrCreateChannelConversation(ctx context.Context, communityID int64) (int64, error) {
	var convID int64
	err := s.db.QueryRow(ctx,
		"SELECT id FROM conversations WHERE type = 'channel' AND community_id = $1",
		communityID).Scan(&convID)
	if err == nil {
		return convID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("查询社区频道失败: %w", err)
	}

	// 创建社区频道
	convID = snowflake.Generate()
	_, err = s.db.Exec(ctx,
		"INSERT INTO conversations (id, type, community_id, name) VALUES ($1, 'channel', $2, '')",
		convID, communityID)
	if err != nil {
		return 0, fmt.Errorf("创建社区频道失败: %w", err)
	}

	return convID, nil
}

// JoinChannelConversation 用户加入社区频道
func (s *ChatService) JoinChannelConversation(ctx context.Context, userID, communityID int64) error {
	convID, err := s.GetOrCreateChannelConversation(ctx, communityID)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(ctx,
		"INSERT INTO conversation_participants (conversation_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		convID, userID)
	if err != nil {
		return fmt.Errorf("加入社区频道失败: %w", err)
	}
	return nil
}

// LeaveChannelConversation 用户离开社区频道
func (s *ChatService) LeaveChannelConversation(ctx context.Context, userID, communityID int64) error {
	var convID int64
	err := s.db.QueryRow(ctx,
		"SELECT id FROM conversations WHERE type = 'channel' AND community_id = $1",
		communityID).Scan(&convID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("查询社区频道失败: %w", err)
	}

	_, err = s.db.Exec(ctx,
		"DELETE FROM conversation_participants WHERE conversation_id = $1 AND user_id = $2",
		convID, userID)
	if err != nil {
		return fmt.Errorf("离开社区频道失败: %w", err)
	}
	return nil
}

// IsConversationParticipant 检查用户是否属于会话
func (s *ChatService) IsConversationParticipant(ctx context.Context, conversationID, userID int64) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM conversation_participants WHERE conversation_id = $1 AND user_id = $2)",
		conversationID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("查询参与者失败: %w", err)
	}
	return exists, nil
}

// SaveMessage 保存消息
func (s *ChatService) SaveMessage(ctx context.Context, conversationID, senderID int64, content, msgType string) (*model.ChatMessage, error) {
	msg := &model.ChatMessage{
		ID:             snowflake.Generate(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		Type:           msgType,
	}

	_, err := s.db.Exec(ctx,
		"INSERT INTO messages (id, conversation_id, sender_id, content, type) VALUES ($1, $2, $3, $4, $5)",
		msg.ID, msg.ConversationID, msg.SenderID, msg.Content, msg.Type)
	if err != nil {
		return nil, fmt.Errorf("保存消息失败: %w", err)
	}

	// 回填 created_at
	err = s.db.QueryRow(ctx,
		"SELECT created_at FROM messages WHERE id = $1",
		msg.ID).Scan(&msg.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("查询消息时间失败: %w", err)
	}

	return msg, nil
}

// GetConversationParticipants 获取会话所有参与者 userID 列表
func (s *ChatService) GetConversationParticipants(ctx context.Context, conversationID int64) ([]int64, error) {
	rows, err := s.db.Query(ctx,
		"SELECT user_id FROM conversation_participants WHERE conversation_id = $1",
		conversationID)
	if err != nil {
		return nil, fmt.Errorf("查询参与者失败: %w", err)
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var uid int64
		if err := rows.Scan(&uid); err != nil {
			return nil, fmt.Errorf("扫描参与者失败: %w", err)
		}
		userIDs = append(userIDs, uid)
	}

	return userIDs, nil
}

// ListMessages 分页获取历史消息
func (s *ChatService) ListMessages(ctx context.Context, conversationID, lastID int64, limit int) ([]*model.ChatMessage, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	var rows pgx.Rows
	var err error
	if lastID > 0 {
		rows, err = s.db.Query(ctx,
			"SELECT id, conversation_id, sender_id, content, type, created_at FROM messages WHERE conversation_id = $1 AND id < $2 ORDER BY id DESC LIMIT $3",
			conversationID, lastID, limit)
	} else {
		rows, err = s.db.Query(ctx,
			"SELECT id, conversation_id, sender_id, content, type, created_at FROM messages WHERE conversation_id = $1 ORDER BY id DESC LIMIT $2",
			conversationID, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("查询消息失败: %w", err)
	}
	defer rows.Close()

	var messages []*model.ChatMessage
	for rows.Next() {
		m := &model.ChatMessage{}
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.Type, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("扫描消息失败: %w", err)
		}
		messages = append(messages, m)
	}

	// 反转顺序，使前端按时间正序展示
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// ListUserConversations 获取用户的会话列表
func (s *ChatService) ListUserConversations(ctx context.Context, userID int64) ([]map[string]interface{}, error) {
	rows, err := s.db.Query(ctx,
		`SELECT c.id, c.type, c.community_id, c.name, c.created_at
		 FROM conversations c
		 JOIN conversation_participants p ON c.id = p.conversation_id
		 WHERE p.user_id = $1
		 ORDER BY c.created_at DESC`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("查询会话列表失败: %w", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var conv model.Conversation
		var communityID *int64
		if err := rows.Scan(&conv.ID, &conv.Type, &communityID, &conv.Name, &conv.CreatedAt); err != nil {
			return nil, fmt.Errorf("扫描会话失败: %w", err)
		}
		conv.CommunityID = communityID
		result = append(result, map[string]interface{}{
			"id":           conv.ID,
			"type":         conv.Type,
			"community_id": conv.CommunityID,
			"name":         conv.Name,
			"created_at":   conv.CreatedAt,
		})
	}

	return result, nil
}

// GetUserByID 查询用户（给 handler 组装发送者信息）
func (s *ChatService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}
	err := s.db.QueryRow(ctx,
		"SELECT id, phone, nickname, avatar_url, credit_score, created_at, updated_at FROM users WHERE id = $1",
		id).Scan(&user.ID, &user.Phone, &user.Nickname, &user.AvatarURL, &user.CreditScore, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return user, nil
}

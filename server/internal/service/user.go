package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wenyuwang1985/community/internal/model"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

func (s *UserService) UpdateProfile(ctx context.Context, userID int64, nickname, avatarURL string) (*model.User, error) {
	_, err := s.db.Exec(ctx,
		"UPDATE users SET nickname = $1, avatar_url = $2, updated_at = now() WHERE id = $3",
		nickname, avatarURL, userID)
	if err != nil {
		return nil, fmt.Errorf("更新用户资料失败: %w", err)
	}

	return s.GetUserByID(ctx, userID)
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
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

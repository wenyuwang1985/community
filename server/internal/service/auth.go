package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wenyuwang1985/community/internal/model"
	"github.com/wenyuwang1985/community/pkg/snowflake"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Register(ctx context.Context, phone, password string) (*model.User, error) {
	// 检查手机号是否已注册
	var exists bool
	err := s.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)", phone).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	if exists {
		return nil, errors.New("该手机号已注册")
	}

	// 密码哈希
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 生成 ID 并插入
	user := &model.User{
		ID:    snowflake.Generate(),
		Phone: phone,
	}

	_, err = s.db.Exec(ctx,
		"INSERT INTO users (id, phone, password) VALUES ($1, $2, $3)",
		user.ID, user.Phone, string(hashed))
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, phone, password string) (*model.User, error) {
	user := &model.User{}
	err := s.db.QueryRow(ctx,
		"SELECT id, phone, nickname, avatar_url, password, credit_score, created_at, updated_at FROM users WHERE phone = $1",
		phone).Scan(&user.ID, &user.Phone, &user.Nickname, &user.AvatarURL, &user.Password, &user.CreditScore, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("手机号或密码错误")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("手机号或密码错误")
	}

	return user, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
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

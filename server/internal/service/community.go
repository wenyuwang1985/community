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

type CommunityService struct {
	db *pgxpool.Pool
}

func NewCommunityService(db *pgxpool.Pool) *CommunityService {
	return &CommunityService{db: db}
}

func (s *CommunityService) SearchCommunities(ctx context.Context, keyword string) ([]*model.Community, error) {
	rows, err := s.db.Query(ctx,
		"SELECT id, name, province, city, district, adcode, created_at FROM communities WHERE name ILIKE $1 LIMIT 20",
		"%"+keyword+"%")
	if err != nil {
		return nil, fmt.Errorf("搜索社区失败: %w", err)
	}
	defer rows.Close()

	var communities []*model.Community
	for rows.Next() {
		c := &model.Community{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Province, &c.City, &c.District, &c.Adcode, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("扫描社区数据失败: %w", err)
		}
		communities = append(communities, c)
	}

	return communities, nil
}

func (s *CommunityService) Subscribe(ctx context.Context, userID, communityID int64) error {
	// 检查社区是否存在
	var exists bool
	err := s.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM communities WHERE id = $1)", communityID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("查询社区失败: %w", err)
	}
	if !exists {
		return errors.New("社区不存在")
	}

	// 检查是否已订阅
	var subscribed bool
	err = s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM user_community_subscriptions WHERE user_id = $1 AND community_id = $2)",
		userID, communityID).Scan(&subscribed)
	if err != nil {
		return fmt.Errorf("查询订阅关系失败: %w", err)
	}
	if subscribed {
		return errors.New("已订阅该社区")
	}

	// 插入订阅记录，如果是第一个订阅则设为主社区
	_, err = s.db.Exec(ctx,
		`INSERT INTO user_community_subscriptions (id, user_id, community_id, is_primary)
		 VALUES ($1, $2, $3, NOT EXISTS(SELECT 1 FROM user_community_subscriptions WHERE user_id = $2))`,
		snowflake.Generate(), userID, communityID)
	if err != nil {
		return fmt.Errorf("订阅社区失败: %w", err)
	}

	return nil
}

func (s *CommunityService) Unsubscribe(ctx context.Context, userID, communityID int64) error {
	// 查询当前订阅数量
	var count int
	err := s.db.QueryRow(ctx,
		"SELECT COUNT(*) FROM user_community_subscriptions WHERE user_id = $1",
		userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("查询订阅数量失败: %w", err)
	}
	if count <= 1 {
		return errors.New("至少保留一个订阅社区")
	}

	// 删除订阅关系
	result, err := s.db.Exec(ctx,
		"DELETE FROM user_community_subscriptions WHERE user_id = $1 AND community_id = $2",
		userID, communityID)
	if err != nil {
		return fmt.Errorf("取消订阅失败: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("未订阅该社区")
	}

	return nil
}

func (s *CommunityService) GetUserSubscriptions(ctx context.Context, userID int64) ([]*model.Community, error) {
	rows, err := s.db.Query(ctx,
		`SELECT c.id, c.name, c.province, c.city, c.district, c.adcode, c.created_at, s.is_primary
		 FROM communities c
		 JOIN user_community_subscriptions s ON c.id = s.community_id
		 WHERE s.user_id = $1
		 ORDER BY s.is_primary DESC, s.created_at ASC`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("查询订阅列表失败: %w", err)
	}
	defer rows.Close()

	type CommunityWithPrimary struct {
		model.Community
		IsPrimary bool `json:"is_primary"`
	}

	var communities []*model.Community
	for rows.Next() {
		c := &model.Community{}
		var isPrimary bool
		if err := rows.Scan(&c.ID, &c.Name, &c.Province, &c.City, &c.District, &c.Adcode, &c.CreatedAt, &isPrimary); err != nil {
			return nil, fmt.Errorf("扫描订阅数据失败: %w", err)
		}
		// 简单处理：把 is_primary 放到一个扩展结构里返回给 handler
		// 这里先直接返回 community，is_primary 由 handler 额外处理
		_ = isPrimary
		communities = append(communities, c)
	}

	return communities, nil
}

func (s *CommunityService) GetUserSubscriptionsWithPrimary(ctx context.Context, userID int64) ([]map[string]interface{}, error) {
	rows, err := s.db.Query(ctx,
		`SELECT c.id, c.name, c.province, c.city, c.district, c.adcode, c.created_at, s.is_primary
		 FROM communities c
		 JOIN user_community_subscriptions s ON c.id = s.community_id
		 WHERE s.user_id = $1
		 ORDER BY s.is_primary DESC, s.created_at ASC`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("查询订阅列表失败: %w", err)
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		c := &model.Community{}
		var isPrimary bool
		if err := rows.Scan(&c.ID, &c.Name, &c.Province, &c.City, &c.District, &c.Adcode, &c.CreatedAt, &isPrimary); err != nil {
			return nil, fmt.Errorf("扫描订阅数据失败: %w", err)
		}
		result = append(result, map[string]interface{}{
			"id":         c.ID,
			"name":       c.Name,
			"province":   c.Province,
			"city":       c.City,
			"district":   c.District,
			"adcode":     c.Adcode,
			"created_at": c.CreatedAt,
			"is_primary": isPrimary,
		})
	}

	return result, nil
}

func (s *CommunityService) SetPrimaryCommunity(ctx context.Context, userID, communityID int64) error {
	// 检查是否已订阅
	var exists bool
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM user_community_subscriptions WHERE user_id = $1 AND community_id = $2)",
		userID, communityID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("查询订阅关系失败: %w", err)
	}
	if !exists {
		return errors.New("未订阅该社区")
	}

	// 开启事务，将其他设为非主社区，目标设为主社区
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		"UPDATE user_community_subscriptions SET is_primary = false WHERE user_id = $1",
		userID)
	if err != nil {
		return fmt.Errorf("重置主社区失败: %w", err)
	}

	_, err = tx.Exec(ctx,
		"UPDATE user_community_subscriptions SET is_primary = true WHERE user_id = $1 AND community_id = $2",
		userID, communityID)
	if err != nil {
		return fmt.Errorf("设置主社区失败: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

func (s *CommunityService) GetCommunityByID(ctx context.Context, id int64) (*model.Community, error) {
	c := &model.Community{}
	err := s.db.QueryRow(ctx,
		"SELECT id, name, province, city, district, adcode, created_at FROM communities WHERE id = $1",
		id).Scan(&c.ID, &c.Name, &c.Province, &c.City, &c.District, &c.Adcode, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("社区不存在")
		}
		return nil, fmt.Errorf("查询社区失败: %w", err)
	}
	return c, nil
}

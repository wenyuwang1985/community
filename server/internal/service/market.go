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

type MarketService struct {
	db *pgxpool.Pool
}

func NewMarketService(db *pgxpool.Pool) *MarketService {
	return &MarketService{db: db}
}

// CreateItem 发布商品
func (s *MarketService) CreateItem(ctx context.Context, sellerID, communityID int64, title string, price int, condition, category string, images []string) (*model.Item, error) {
	// 检查用户是否已订阅该社区
	var subscribed bool
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM user_community_subscriptions WHERE user_id = $1 AND community_id = $2)",
		sellerID, communityID).Scan(&subscribed)
	if err != nil {
		return nil, fmt.Errorf("查询订阅关系失败: %w", err)
	}
	if !subscribed {
		return nil, errors.New("未订阅该社区，无法发布商品")
	}

	if len(images) > 9 {
		images = images[:9]
	}

	item := &model.Item{
		ID:          snowflake.Generate(),
		SellerID:    sellerID,
		CommunityID: communityID,
		Title:       title,
		Price:       price,
		Condition:   condition,
		Category:    category,
		Images:      images,
		Status:      "selling",
	}

	_, err = s.db.Exec(ctx,
		"INSERT INTO items (id, seller_id, community_id, title, price, condition, category, images, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		item.ID, item.SellerID, item.CommunityID, item.Title, item.Price, item.Condition, item.Category, item.Images, item.Status)
	if err != nil {
		return nil, fmt.Errorf("发布商品失败: %w", err)
	}

	return item, nil
}

// GetItemByID 获取商品详情
func (s *MarketService) GetItemByID(ctx context.Context, id int64) (*model.Item, error) {
	item := &model.Item{}
	err := s.db.QueryRow(ctx,
		"SELECT id, seller_id, community_id, title, price, condition, category, images, status, created_at, updated_at FROM items WHERE id = $1",
		id).Scan(&item.ID, &item.SellerID, &item.CommunityID, &item.Title, &item.Price, &item.Condition, &item.Category, &item.Images, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("商品不存在")
		}
		return nil, fmt.Errorf("查询商品失败: %w", err)
	}
	return item, nil
}

// UpdateItem 修改商品信息（仅卖家）
func (s *MarketService) UpdateItem(ctx context.Context, sellerID, itemID int64, title string, price int, condition, category string, images []string) (*model.Item, error) {
	if len(images) > 9 {
		images = images[:9]
	}

	result, err := s.db.Exec(ctx,
		"UPDATE items SET title = $1, price = $2, condition = $3, category = $4, images = $5, updated_at = now() WHERE id = $6 AND seller_id = $7 AND status = 'selling'",
		title, price, condition, category, images, itemID, sellerID)
	if err != nil {
		return nil, fmt.Errorf("更新商品失败: %w", err)
	}
	if result.RowsAffected() == 0 {
		return nil, errors.New("商品不存在或无权修改")
	}

	return s.GetItemByID(ctx, itemID)
}

// ListItems 分页获取商品列表（仅出售中）
func (s *MarketService) ListItems(ctx context.Context, communityID int64, category string, lastID int64, limit int) ([]*model.Item, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	var rows pgx.Rows
	var err error

	if category != "" {
		if lastID > 0 {
			rows, err = s.db.Query(ctx,
				"SELECT id, seller_id, community_id, title, price, condition, category, images, status, created_at, updated_at FROM items WHERE community_id = $1 AND category = $2 AND status = 'selling' AND id < $3 ORDER BY id DESC LIMIT $4",
				communityID, category, lastID, limit)
		} else {
			rows, err = s.db.Query(ctx,
				"SELECT id, seller_id, community_id, title, price, condition, category, images, status, created_at, updated_at FROM items WHERE community_id = $1 AND category = $2 AND status = 'selling' ORDER BY id DESC LIMIT $3",
				communityID, category, limit)
		}
	} else {
		if lastID > 0 {
			rows, err = s.db.Query(ctx,
				"SELECT id, seller_id, community_id, title, price, condition, category, images, status, created_at, updated_at FROM items WHERE community_id = $1 AND status = 'selling' AND id < $2 ORDER BY id DESC LIMIT $3",
				communityID, lastID, limit)
		} else {
			rows, err = s.db.Query(ctx,
				"SELECT id, seller_id, community_id, title, price, condition, category, images, status, created_at, updated_at FROM items WHERE community_id = $1 AND status = 'selling' ORDER BY id DESC LIMIT $2",
				communityID, limit)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("查询商品列表失败: %w", err)
	}
	defer rows.Close()

	var items []*model.Item
	for rows.Next() {
		item := &model.Item{}
		if err := rows.Scan(&item.ID, &item.SellerID, &item.CommunityID, &item.Title, &item.Price, &item.Condition, &item.Category, &item.Images, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("扫描商品数据失败: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// MarkSold 标记已售
func (s *MarketService) MarkSold(ctx context.Context, sellerID, itemID int64) error {
	result, err := s.db.Exec(ctx,
		"UPDATE items SET status = 'sold', updated_at = now() WHERE id = $1 AND seller_id = $2 AND status = 'selling'",
		itemID, sellerID)
	if err != nil {
		return fmt.Errorf("标记已售失败: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("商品不存在或无权操作")
	}
	return nil
}

// MarkOff 下架商品
func (s *MarketService) MarkOff(ctx context.Context, sellerID, itemID int64) error {
	result, err := s.db.Exec(ctx,
		"UPDATE items SET status = 'off', updated_at = now() WHERE id = $1 AND seller_id = $2 AND status = 'selling'",
		itemID, sellerID)
	if err != nil {
		return fmt.Errorf("下架商品失败: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("商品不存在或无权操作")
	}
	return nil
}

// GetUserByID 查询用户（给 handler 组装卖家信息）
func (s *MarketService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
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

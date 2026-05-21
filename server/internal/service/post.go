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

type PostService struct {
	db *pgxpool.Pool
}

func NewPostService(db *pgxpool.Pool) *PostService {
	return &PostService{db: db}
}

// CreatePost 创建帖子
func (s *PostService) CreatePost(ctx context.Context, userID, communityID int64, tag, content string, images []string) (*model.Post, error) {
	// 检查用户是否已订阅该社区
	var subscribed bool
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM user_community_subscriptions WHERE user_id = $1 AND community_id = $2)",
		userID, communityID).Scan(&subscribed)
	if err != nil {
		return nil, fmt.Errorf("查询订阅关系失败: %w", err)
	}
	if !subscribed {
		return nil, errors.New("未订阅该社区，无法发帖")
	}

	if len(images) > 9 {
		images = images[:9]
	}

	post := &model.Post{
		ID:          snowflake.Generate(),
		UserID:      userID,
		CommunityID: communityID,
		Tag:         tag,
		Content:     content,
		Images:      images,
		Status:      "normal",
	}

	_, err = s.db.Exec(ctx,
		"INSERT INTO posts (id, user_id, community_id, tag, content, images, status) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		post.ID, post.UserID, post.CommunityID, post.Tag, post.Content, post.Images, post.Status)
	if err != nil {
		return nil, fmt.Errorf("创建帖子失败: %w", err)
	}

	return post, nil
}

// GetPostByID 获取帖子详情
func (s *PostService) GetPostByID(ctx context.Context, postID int64) (*model.Post, error) {
	post := &model.Post{}
	err := s.db.QueryRow(ctx,
		"SELECT id, user_id, community_id, tag, content, images, status, like_count, comment_count, created_at, updated_at FROM posts WHERE id = $1 AND status = 'normal'",
		postID).Scan(&post.ID, &post.UserID, &post.CommunityID, &post.Tag, &post.Content, &post.Images, &post.Status, &post.LikeCount, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("帖子不存在")
		}
		return nil, fmt.Errorf("查询帖子失败: %w", err)
	}
	return post, nil
}

// DeletePost 删除帖子（只能删除自己的）
func (s *PostService) DeletePost(ctx context.Context, userID, postID int64) error {
	result, err := s.db.Exec(ctx,
		"UPDATE posts SET status = 'deleted', updated_at = now() WHERE id = $1 AND user_id = $2 AND status = 'normal'",
		postID, userID)
	if err != nil {
		return fmt.Errorf("删除帖子失败: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("帖子不存在或无权删除")
	}
	return nil
}

// ListPosts 分页获取社区 Feed 流
func (s *PostService) ListPosts(ctx context.Context, communityID, lastID int64, limit int) ([]*model.Post, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	var rows pgx.Rows
	var err error
	if lastID > 0 {
		rows, err = s.db.Query(ctx,
			"SELECT id, user_id, community_id, tag, content, images, status, like_count, comment_count, created_at, updated_at FROM posts WHERE community_id = $1 AND status = 'normal' AND id < $2 ORDER BY id DESC LIMIT $3",
			communityID, lastID, limit)
	} else {
		rows, err = s.db.Query(ctx,
			"SELECT id, user_id, community_id, tag, content, images, status, like_count, comment_count, created_at, updated_at FROM posts WHERE community_id = $1 AND status = 'normal' ORDER BY id DESC LIMIT $2",
			communityID, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("查询帖子列表失败: %w", err)
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		p := &model.Post{}
		if err := rows.Scan(&p.ID, &p.UserID, &p.CommunityID, &p.Tag, &p.Content, &p.Images, &p.Status, &p.LikeCount, &p.CommentCount, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("扫描帖子数据失败: %w", err)
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// IsLiked 检查用户是否点赞
func (s *PostService) IsLiked(ctx context.Context, postID, userID int64) (bool, error) {
	var liked bool
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM post_likes WHERE post_id = $1 AND user_id = $2)",
		postID, userID).Scan(&liked)
	if err != nil {
		return false, fmt.Errorf("查询点赞状态失败: %w", err)
	}
	return liked, nil
}

// Like 点赞
func (s *PostService) Like(ctx context.Context, postID, userID int64) error {
	_, err := s.db.Exec(ctx,
		"INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		postID, userID)
	if err != nil {
		return fmt.Errorf("点赞失败: %w", err)
	}

	// 更新点赞数
	_, err = s.db.Exec(ctx,
		"UPDATE posts SET like_count = like_count + 1 WHERE id = $1 AND status = 'normal'",
		postID)
	if err != nil {
		return fmt.Errorf("更新点赞数失败: %w", err)
	}

	return nil
}

// Unlike 取消点赞
func (s *PostService) Unlike(ctx context.Context, postID, userID int64) error {
	result, err := s.db.Exec(ctx,
		"DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2",
		postID, userID)
	if err != nil {
		return fmt.Errorf("取消点赞失败: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("未点赞")
	}

	// 更新点赞数
	_, err = s.db.Exec(ctx,
		"UPDATE posts SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1 AND status = 'normal'",
		postID)
	if err != nil {
		return fmt.Errorf("更新点赞数失败: %w", err)
	}

	return nil
}

// CreateComment 发表评论
func (s *PostService) CreateComment(ctx context.Context, postID, userID int64, content string) (*model.Comment, error) {
	// 检查帖子是否存在且未删除
	var exists bool
	err := s.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1 AND status = 'normal')",
		postID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("查询帖子失败: %w", err)
	}
	if !exists {
		return nil, errors.New("帖子不存在")
	}

	comment := &model.Comment{
		ID:      snowflake.Generate(),
		PostID:  postID,
		UserID:  userID,
		Content: content,
	}

	_, err = s.db.Exec(ctx,
		"INSERT INTO comments (id, post_id, user_id, content) VALUES ($1, $2, $3, $4)",
		comment.ID, comment.PostID, comment.UserID, comment.Content)
	if err != nil {
		return nil, fmt.Errorf("发表评论失败: %w", err)
	}

	// 更新评论数
	_, err = s.db.Exec(ctx,
		"UPDATE posts SET comment_count = comment_count + 1 WHERE id = $1",
		postID)
	if err != nil {
		return nil, fmt.Errorf("更新评论数失败: %w", err)
	}

	return comment, nil
}

// ListComments 获取评论列表
func (s *PostService) ListComments(ctx context.Context, postID int64) ([]*model.Comment, error) {
	rows, err := s.db.Query(ctx,
		"SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = $1 ORDER BY created_at ASC",
		postID)
	if err != nil {
		return nil, fmt.Errorf("查询评论列表失败: %w", err)
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		c := &model.Comment{}
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("扫描评论数据失败: %w", err)
		}
		comments = append(comments, c)
	}

	return comments, nil
}

// GetUserByID 复用查询用户（给 handler 组装作者信息）
func (s *PostService) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
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

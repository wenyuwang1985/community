package model

import "time"

type Post struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	CommunityID  int64     `json:"community_id"`
	Tag          string    `json:"tag"`
	Content      string    `json:"content"`
	Images       []string  `json:"images"`
	Status       string    `json:"status"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

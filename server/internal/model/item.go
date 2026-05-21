package model

import "time"

type Item struct {
	ID          int64     `json:"id"`
	SellerID    int64     `json:"seller_id"`
	CommunityID int64     `json:"community_id"`
	Title       string    `json:"title"`
	Price       int       `json:"price"`
	Condition   string    `json:"condition"`
	Category    string    `json:"category"`
	Images      []string  `json:"images"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

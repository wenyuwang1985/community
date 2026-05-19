package model

import "time"

type Community struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Province  string    `json:"province"`
	City      string    `json:"city"`
	District  string    `json:"district"`
	Adcode    string    `json:"adcode"`
	CreatedAt time.Time `json:"created_at"`
}

type UserCommunitySubscription struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CommunityID int64     `json:"community_id"`
	IsPrimary   bool      `json:"is_primary"`
	CreatedAt   time.Time `json:"created_at"`
}

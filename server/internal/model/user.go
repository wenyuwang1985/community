package model

import "time"

type User struct {
	ID          int64     `json:"id"`
	Phone       string    `json:"phone"`
	Nickname    string    `json:"nickname"`
	AvatarURL   string    `json:"avatar_url"`
	Password    string    `json:"-"`
	CreditScore int16     `json:"credit_score"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

package entity

import "time"

type UserID int64

func (id UserID) Int64() int64 { return int64(id) }

type User struct {
	ID       UserID `json:"id"`
	Username string `json:"username"`
	Password string `json:"-" gorm:"column:password_hash"`
	CreatedAt time.Time `json:"-" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at;autoUpdateTime"`
}

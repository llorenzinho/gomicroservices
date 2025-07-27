package models

import "time"

type User struct {
	Id        uint      `gorm:"primaryKey;column:id" json:"id"`
	Username  string    `gorm:"column:username;uniqueIndex" json:"username"`
	Email     string    `gorm:"column:email;uniqueIndex" json:"email"`
	Password  string    `gorm:"column:password" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}

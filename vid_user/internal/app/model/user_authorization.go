package model

import (
	"gorm.io/gorm"
	"vid_user/internal/app/uuid"
)

type UserAuthorization struct {
	ID           int64  `gorm:"primaryKey" json:"id"`
	UserID       int64  `gorm:"not null;column:user_id" json:"user_id"`
	IdentityType string `gorm:"not null;column:identity_type;default:phone;comment:登录类型(手机号/邮箱) 或第三方应用名称 (微信/微博等)" json:"identity_type"`
	Identifier   string `gorm:"not null;column:identifier;comment:唯一标识" json:"identifier"`
	Credential   string `gorm:"column:credential;comment:密码凭证或token" json:"credential"`
	UpdatedAt    int64  `json:"updated_at"`
	CreatedAt    int64  `gorm:"autoCreateTime" json:"created_at"`
}

func (u *UserAuthorization) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.GenerateID()
	return
}

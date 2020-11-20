package model

import (
	"gorm.io/gorm"
	"vid_user/internal/app/uuid"
)

type User struct {
	ID        int64  `gorm:"primaryKey" json:"id"`
	NickName  string `gorm:"column:nickname" json:"nickname"`
	Avatar    string `gorm:"column:avatar;default:http://oss.seefs.cn/avatar.jpg" json:"avatar"`
	Gender    int32  `gorm:"column:gender;default:0;comment:0为未知,1为男,2为女" json:"gender"`
	UpdatedAt int64  `json:"updated_at"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.GenerateID()
	return
}

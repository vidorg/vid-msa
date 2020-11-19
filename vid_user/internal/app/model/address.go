package model

type Address struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Content   string `json:"content"`
	Phone     string `json:"phone"`
	OpenID    string `json:"openid"`
	Sort      *int   `json:"sort"`
	IsDefault bool   `json:"is_default"`
	Nickname  string `gorm:"column:nickname" json:"nickname"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"create_time"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"update_time"`
}

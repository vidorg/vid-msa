package model

type User struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName  string `gorm:"column:username" json:"username"`
	Password  string `gorm:"column:password" json:"password"`
	Avatar    string `gorm:"column:avatar;default:http://oss.seefs.cn/avatar.jpg" json:"avatar"`
	Gender    int32  `gorm:"column:gender;default:0;comment:0为未知,1为男,2为女" json:"gender"`
	UpdatedAt int64  `json:"updated_at"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
}

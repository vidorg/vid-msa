package jwt

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/kataras/jwt"
)

var secret string
var appName string

// Init ...
func Init() {
	secret = conf.GetString("jwt.secret")
	appName = conf.GetString("jwt.app_name")
}

// Claims JWT
type Claims struct {
	UID int64 `json:"uid"`
}

// GenerateToken 生成token
func GenerateToken(uid int64) (string, error) {
	claims := Claims{
		UID: uid,
	}
	//now := time.Now()
	//standardClaims := jwt.Claims{
	//	Expiry:   expireTime.Unix(),
	//	IssuedAt: now.Unix(),
	//	Issuer:   "my-app",
	//}
	token, err := jwt.Sign(jwt.HS256, []byte(secret), claims)
	return string(token), err
}

// ParseToken 转换token
func ParseToken(token string) (*Claims, error) {
	verify, err := jwt.Verify(jwt.HS256, []byte(secret), []byte(token))
	var claims *Claims
	err = verify.Claims(&claims)

	return claims, err
}

// RedisTokenConcat Redis中键名
func RedisTokenConcat(uid int64, token string) string {
	return fmt.Sprintf("%s-%d-%s", appName, uid, token)
}

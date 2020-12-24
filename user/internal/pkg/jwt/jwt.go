package jwt

import (
	"fmt"
	"github.com/kataras/jwt"
	"strings"
	"time"
)

// Claims JWT
type Claims struct {
	UID int64 `json:"uid"`
}

// GenerateToken 生成token
func GenerateToken(uid int64, secret []byte) (string, error) {
	claims := Claims{
		UID: uid,
	}
	//now := time.Now()
	//standardClaims := jwt.Claims{
	//	Expiry:   expireTime.Unix(),
	//	IssuedAt: now.Unix(),
	//	Issuer:   "my-app",
	//}
	token, err := jwt.Sign(jwt.HS256, secret, claims, jwt.MaxAge(1500*time.Hour))
	return string(token), err
}

// ParseToken 转换token
func ParseToken(token string, secret []byte) (Claims, error) {
	verify, err := jwt.Verify(jwt.HS256, secret, []byte(token))
	if err != nil {
		return Claims{}, err
	}
	var claims Claims
	err = verify.Claims(&claims)

	return claims, err
}

// RedisTokenConcat Redis中键名
func RedisTokenConcat(uid, token string, appName string) string {
	return fmt.Sprintf("%s-%s-%s", appName, uid, token)
}

// GetTokenFromRedisPattern 从redis存储的pattern取token
func GetTokenFromRedisPattern(pattern string) string {
	split := strings.Split(pattern, "-")
	return split[2]
}

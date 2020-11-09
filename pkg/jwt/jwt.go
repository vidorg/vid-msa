package jwt

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
)

// Claims JWT
type Claims struct {
	jwt.StandardClaims
	UID int
}

// GenerateToken 生成token
func GenerateToken(uid int) (string, error) {
	expireTime := time.Now().Add(time.Duration(viper.GetInt("jwt.expire")) * time.Second)
	claims := Claims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    viper.GetString("jwt.issuer"),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString("jwt.secret")))
}

// ParseToken 转换token
func ParseToken(token string) (*Claims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	}
	tokenObj, err := jwt.ParseWithClaims(token, &Claims{}, keyFunc)
	if err != nil {
		return nil, err
	}
	claims, ok := tokenObj.Claims.(*Claims)
	if !ok || !tokenObj.Valid {
		return nil, jwt.ValidationError{Errors: jwt.ValidationErrorClaimsInvalid}
	}
	return claims, nil
}

// IsTokenExpired 是否过期
func IsTokenExpired(err error) bool {
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return true
		}
	}
	return false
}

// RedisTokenConcat Redis中键名
func RedisTokenConcat(uid, token string) string {
	return fmt.Sprintf("%s-%s-%s", viper.GetString("meta.name"), uid, token)
}

package app

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-programming-tour/blog-service/global"
	"github.com/go-programming-tour/blog-service/pkg/util"
)

type Claims struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
	jwt.StandardClaims
}

func GetJWTSecret() []byte {
	return []byte(global.JWTSetting.Secret)
}

func GenerateToken(appKey, appSecret string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(global.JWTSetting.Expire)
	claims := Claims{
		AppKey:    util.EncodeMD5(appKey),
		AppSecret: util.EncodeMD5(appSecret),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    global.JWTSetting.Issuer,
		},
	}
	// 生成jwt token
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名字符串
	token, err := tokenClaims.SignedString(GetJWTSecret())
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

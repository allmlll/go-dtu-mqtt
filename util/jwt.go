package util

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserClaims struct {
	Id                   primitive.ObjectID `json:"id"`
	jwt.RegisteredClaims                    // 内嵌标准的声明
}

// tokenExpireDuration token过期时间
const tokenExpireDuration = time.Hour * 72

// secret 用于加盐的字符串
var secret = []byte("3a3deb17-ab80-41ae-8729-80a8e4f34381")

const TokenPrefix = "Bearer"

// CreateToken 生成JWT
func CreateToken(Id primitive.ObjectID) (tokenString string, err error) {
	// 创建一个我们自己的声明
	claims := UserClaims{
		Id, // 自定义字段
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpireDuration)), // 设置过期时间
			Issuer:    "allmlll",                                               // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	tokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, err = tokenObject.SignedString(secret)
	return tokenString, err
}

// ParseToken 解析token
func ParseToken(tokenString string) (*UserClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

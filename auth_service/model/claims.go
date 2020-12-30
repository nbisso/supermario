package model

import "github.com/dgrijalva/jwt-go"

type SecurityClaims struct {
	UserName string `json:"username"`
	UserID   int64  `json:"user_id"`
	jwt.StandardClaims
}

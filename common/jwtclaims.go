package common

import "github.com/dgrijalva/jwt-go"

type JWTCustomClaims struct {
	Username  string `json:"username"`
	IsAdmin bool   `json:"isAdmin"`
	jwt.StandardClaims
}

package common

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type JWTCustomClaims struct {
	Username  string `json:"username"`
	IsAdmin bool   `json:"isAdmin"`
	jwt.StandardClaims
}

// IsAdmin returns true if the user is an admin
func IsAdmin(ctx echo.Context) bool {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(JWTCustomClaims)
	return claims.IsAdmin
}

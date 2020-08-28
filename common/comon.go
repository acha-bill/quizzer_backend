package common

import (
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type JWTCustomClaims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	Id       string `json:"_id"`
	jwt.StandardClaims
}

// IsAdmin returns true if the user is an admin
func IsAdmin(ctx echo.Context) bool {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(JWTCustomClaims)
	return claims.IsAdmin
}

// GetUsername returns the username of the user holding this context
func GetUsername(ctx echo.Context) string {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(JWTCustomClaims)
	return claims.Username
}

// IsDevelopment returns true if the server is running in dev mode.
func IsDevelopment() bool {
	return strings.HasPrefix(os.Getenv("ENV"), "d")
}

package socketserver

import (
	"errors"

	"github.com/acha-bill/quizzer_backend/common"
	userService "github.com/acha-bill/quizzer_backend/packages/dblayer/user"
	"github.com/dgrijalva/jwt-go"
)

var (
	ErrSocketNotAuthenticated = errors.New("socket not authorized")
)

func handleAuthMessage(wsConnection *WsConnection, msg SocketMessageAuth) {
	token, err := jwt.Parse(msg.Token, func(token *jwt.Token) (interface{}, error) {
		return token, nil
	})
	if token == nil || err != nil {
		ServerManager().WriteConnection(wsConnection, NewSocketResponseAuth("err parsing jwt"))
		return
	}
	if !token.Valid {
		ServerManager().WriteConnection(wsConnection, NewSocketResponseAuth("Invalid jwt"))
		return
	}
	claims := token.Claims.(common.JWTCustomClaims)
	userId := claims.Id
	user := userService.FindById(userId)
	wsConnection.Context.User = user
	wsConnection.Context.Ready = true

	ServerManager().WriteConnection(wsConnection, NewSocketResponseAuth(""))
}

// SocketMessageAuth is the auth message
type SocketMessageAuth struct {
	Token string
}

const authResponseType = "auth"

// SocketResponseAuth is the auth response
type SocketResponseAuth struct {
	Type  string `json:"type"`
	Error string `json:"error,omitempty"`
}

func NewSocketResponseAuth(err string) SocketResponseAuth {
	return SocketResponseAuth{
		Type:  authResponseType,
		Error: err,
	}
}

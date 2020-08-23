package auth

import (
	"github.com/acha-bill/quizzer_backend/plugins"
	"github.com/labstack/echo/v4"
	"net/http"
	"sync"
)

const (
	PluginName = "auth"
)
var (
	plugin Auth
	once sync.Once
)

type Auth struct {
	name string
	handlers []plugins.PluginHandler
}

func (plugin Auth) AddHandler(method string, path string, handler func(echo.Context) error) {
	pluginHandler := plugins.PluginHandler{
		Path:    path,
		Handler: handler,
		Method: method,
	}
	plugin.handlers = append(plugin.handlers, pluginHandler)
}

func (plugin Auth) Handlers() []plugins.PluginHandler {
	return plugin.handlers
}

func (plugin Auth) Name() string {
	return plugin.name
}

func NewPlugin() Auth {
	plugin := Auth{
		name: PluginName,
	}
	return plugin
}

func Plugin() Auth{
	once.Do(func() {
		plugin = NewPlugin()
	})
	return plugin
}

func init() {
	auth := Plugin()
	auth.AddHandler(http.MethodPost, "/login", login)
	auth.AddHandler(http.MethodPost, "/register", register)
}


///// handlers
func login(ctx echo.Context) error {
	var req LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, LoginResponse{
			Error: err.Error(),
		})
	}

	//check db
	//create token
	//return jwt
	return ctx.JSON(http.StatusOK, LoginResponse{
		JWT:   "adfasdfsaf",
	})
}
func register(cxt echo.Context) error {
	return nil
}


type LoginRequest struct {
	Username  string `json:"name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Error string `json:"error,omitempty"`
	JWT string `json:"jwt,omitempty"`
}


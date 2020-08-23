package plugins

import (
	"github.com/acha-bill/quizzer_backend/plugins/auth"
	"github.com/labstack/echo/v4"
)

type PluginHandler struct {
	Path string
	Handler func(echo.Context) error
	Method string
}

type Plugin interface {
	Name() string
	AddHandler(method string, path string, handler func(echo.Context) error)
	Handlers() []PluginHandler
}

var (
	Plugins = []Plugin {
		auth.Plugin(),
	}
)
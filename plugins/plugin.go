package plugins

import (
	"github.com/labstack/echo/v4"
)

// PluginHandler represents the PluginHandler structure
type PluginHandler struct {
	Path      string
	Handler   func(echo.Context) error
	Method    string
	AuthLevel AuthLevel
}

// AuthLevel represents our Auth level enum
type AuthLevel int

const (
	AuthLevelNone AuthLevel = iota
	AuthLevelUser
	AuthLevelAdmin
)

// Plugin interface
type Plugin interface {
	Name() string
	AddHandler(method string, path string, handler func(echo.Context) error, authLevel ...AuthLevel)
	Handlers() []*PluginHandler
}

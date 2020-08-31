package search

import (
	"errors"
	"net/http"
	"sync"

	"github.com/acha-bill/quizzer_backend/common"
	"github.com/acha-bill/quizzer_backend/packages/socketserver"

	"github.com/acha-bill/quizzer_backend/plugins"
	"github.com/labstack/echo/v4"
)

const (
	// PluginName defines the name of the plugin
	PluginName = "search"
)

var (
	plugin              *Search
	once                sync.Once
	ErrUserNotConnected = errors.New("user not connected")
)

const (
	// GameLength defines the number of questions in a game
	GameLength = 10
)

// Search structure
type Search struct {
	name     string
	handlers []*plugins.PluginHandler
}

// AddHandler Method definition from interface
func (plugin *Search) AddHandler(method string, path string, handler func(echo.Context) error, authLevel ...plugins.AuthLevel) {
	pluginHandler := &plugins.PluginHandler{
		Path:      path,
		Handler:   handler,
		Method:    method,
		AuthLevel: plugins.AuthLevelUser,
	}
	if len(authLevel) > 0 {
		pluginHandler.AuthLevel = authLevel[0]
	}
	plugin.handlers = append(plugin.handlers, pluginHandler)
}

// Handlers Method definition from interface
func (plugin *Search) Handlers() []*plugins.PluginHandler {
	return plugin.handlers
}

// Name defines name of the plugin
func (plugin *Search) Name() string {
	return plugin.name
}

// NewPlugin returns the new plugin
func NewPlugin() *Search {
	plugin := &Search{
		name: PluginName,
	}
	return plugin
}

// Plugin returns an instance of the plugin
func Plugin() *Search {
	once.Do(func() {
		plugin = NewPlugin()
	})
	return plugin
}

func init() {
	search := Plugin()
	search.AddHandler(http.MethodGet, "/search", findOpponent)
}

// @Summary search for a random opponent
// @Accept json
// @produce json
// @Router /search [get]
// @Tags Search
// @Success 200 {object} SearchOpponentResponse
func findOpponent(ctx echo.Context) error {
	gameMgr := socketserver.GameManager()
	serverMgr := socketserver.ServerManager()
	username := common.GetUsername(ctx)
	wsConn := serverMgr.GetByUsername(username)
	if wsConn == nil {
		return ctx.JSON(http.StatusNotFound, SearchOpponentResponse{
			Error: ErrUserNotConnected.Error(),
		})
	}
	err := gameMgr.AddSearcher(wsConn)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, SearchOpponentResponse{
			Error: err.Error(),
		})
	}

	player1, player2 := gameMgr.GetPair()
	if player1 != nil && player2 != nil {
		if err := gameMgr.NewGame(player1, player2, GameLength); err != nil {
			return ctx.JSON(http.StatusBadRequest, SearchOpponentResponse{
				Error: err.Error(),
			})
		}
	}

	return ctx.JSON(http.StatusOK, SearchOpponentResponse{})
}

// SearchOpponentResponse represents the Search Response
type SearchOpponentResponse struct {
	Error string `json:"error,omitempty"`
}

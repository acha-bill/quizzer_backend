package search

import (
	"net/http"
	"sync"

	"github.com/acha-bill/quizzer_backend/models"
	userService "github.com/acha-bill/quizzer_backend/packages/dblayer/user"
	"github.com/acha-bill/quizzer_backend/plugins"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// PluginName defines the name of the plugin
	PluginName = "search"
)

var (
	plugin *Search
	once   sync.Once
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
	search.AddHandler(http.MethodGet, "/:userId", findOpponent)
}

// @Summary search for opponent for user with id=userId
// @Accept ID
// @produce json
// @Router /search/:userId [get]
// @Tags Search
// @Param search userId
// @Success 200 {Object} SearchResponse
func findOpponent(ctx echo.Context) error {
	// Find user
	filter := bson.D{primitive.E{Key: "_id", Value: ctx.Param("userId")}}
	users, err := userService.Find(filter)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error: err.Error(),
		})
	}

	if len(users) == 0 {
		return ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error: "User doesn't exist!",
		})
	}

	// Put user in search mode
	users[0].IsSearching = true
	userService.UpdateByID(ctx.Param("userId"), *(users[0]))
	// Get other users in search mode
	filter = bson.D{primitive.E{Key: "isSearching", Value: true}}
	users, err = userService.Find(filter)
	if err != nil {
		return ctx.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error: err.Error(),
		})
	}
	// Pair user with first in search mode
	if len(users) > 0 {
		return ctx.JSON(http.StatusOK, users[0])
	}
	// Or return emtpy list for no opponents found
	return ctx.JSON(http.StatusOK, users)
}

// ErrorResponse represents the Search Response for Errors
type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

// SearchResponse represents the Response object for Search
type SearchResponse models.User

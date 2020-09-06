package category

import (
	"net/http"
	"sync"
	"time"

	"github.com/acha-bill/quizzer_backend/common"
	"github.com/acha-bill/quizzer_backend/models"
	categoryService "github.com/acha-bill/quizzer_backend/packages/dblayer/category"
	"github.com/acha-bill/quizzer_backend/plugins"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PluginName = "question"
)

var (
	plugin *Category
	once   sync.Once
)

type Category struct {
	name     string
	handlers []*plugins.PluginHandler
}

func (plugin *Category) AddHandler(method string, path string, handler func(echo.Context) error, authLevel ...plugins.AuthLevel) {
	pluginHandler := &plugins.PluginHandler{
		Path:      path,
		Handler:   handler,
		Method:    method,
		AuthLevel: plugins.AuthLevelAdmin,
	}
	if len(authLevel) > 0 {
		pluginHandler.AuthLevel = authLevel[0]
	}
	plugin.handlers = append(plugin.handlers, pluginHandler)
}

func (plugin *Category) Handlers() []*plugins.PluginHandler {
	return plugin.handlers
}

func (plugin *Category) Name() string {
	return plugin.name
}

func NewPlugin() *Category {
	plugin := &Category{
		name: PluginName,
	}
	return plugin
}

func Plugin() *Category {
	once.Do(func() {
		plugin = NewPlugin()
	})
	return plugin
}

func init() {
	category := Plugin()
	category.AddHandler(http.MethodPost, "/", create)
	category.AddHandler(http.MethodGet, "/", find)
	category.AddHandler(http.MethodPut, "/:id", edit)
	// TODO: implement
	//category.AddHandler(http.MethodGet, "/:id", remove)
}

// @Summary list all categories
// @Accept  json
// @Produce  json
// @Router /category/ [get]
// @Tags Category
// @Success 201 {object} FindCategoryResponse
func find(ctx echo.Context) error {
	if !common.IsAdmin(ctx) {
		return ctx.JSON(http.StatusUnauthorized, FindCategoryResponse{
			Error: "Unauthorized",
		})
	}
	res, err := categoryService.FindAll()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, FindCategoryResponse{
			Error: err.Error(),
		})
	}
	return ctx.JSON(http.StatusBadRequest, FindCategoryResponse{
		Categories: res,
	})
}

// @Summary create category
// @Accept  json
// @Produce  json
// @Router /category/ [post]
// @Tags Category
// @Param question body CreateCategoryRequest true "create"
// @Success 201 {object} CreateCategoryResponse
func create(ctx echo.Context) error {
	if !common.IsAdmin(ctx) {
		return ctx.JSON(http.StatusUnauthorized, CreateCategoryResponse{
			Error: "Unauthorized",
		})
	}
	var req CreateCategoryRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, CreateCategoryResponse{
			Error: err.Error(),
		})
	}

	q := models.Category{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := categoryService.Create(q)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, CreateCategoryResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusBadRequest, created)
}

// @Summary edit category
// @Accept  json
// @Produce  json
// @Router /category/:id [put]
// @Tags Category
// @Param question body EditCategoryRequest true "create"
// @Success 201 {object} EditCategoryResponse
func edit(ctx echo.Context) error {
	if !common.IsAdmin(ctx) {
		return ctx.JSON(http.StatusUnauthorized, EditCategoryResponse{
			Error: "Unauthorized",
		})
	}
	var req EditCategoryRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, EditCategoryResponse{
			Error: err.Error(),
		})
	}

	categoryID := ctx.Param("id")
	category, err := categoryService.FindById(categoryID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, EditCategoryResponse{
			Error: err.Error(),
		})
	}
	if category == nil {
		return ctx.JSON(http.StatusBadRequest, EditCategoryResponse{
			Error: "Not found",
		})
	}

	category.Name = req.Name
	err = categoryService.UpdateById(categoryID, *category)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, EditCategoryResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusBadRequest, category)
}

// FindCategoryResponse is the find response
type FindCategoryResponse struct {
	Error      string             `json:"error,omitempty"`
	Categories []*models.Category `json:"categories"`
}

type EditCategoryRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EditCategoryResponse struct {
	Error    string           `json:"error,omitempty"`
	Category *models.Category `json:"category"`
}

// CreateCategoryRequest is the request for create category
type CreateCategoryRequest struct {
	Name string `json:"name"`
}

// CreateCategoryResponse is the response for create category
type CreateCategoryResponse struct {
	Error    string           `json:"error,omitempty"`
	Category *models.Category `json:"category"`
}

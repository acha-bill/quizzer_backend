package question

import (
	"net/http"
	"sync"
	"time"

	"github.com/acha-bill/quizzer_backend/common"
	"github.com/acha-bill/quizzer_backend/models"
	questionService "github.com/acha-bill/quizzer_backend/packages/dblayer/question"
	"github.com/acha-bill/quizzer_backend/plugins"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PluginName = "question"
)

var (
	plugin *Question
	once   sync.Once
)

type Question struct {
	name     string
	handlers []*plugins.PluginHandler
}

func (plugin *Question) AddHandler(method string, path string, handler func(echo.Context) error, authLevel ...plugins.AuthLevel) {
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

func (plugin *Question) Handlers() []*plugins.PluginHandler {
	return plugin.handlers
}

func (plugin *Question) Name() string {
	return plugin.name
}

func NewPlugin() *Question {
	plugin := &Question{
		name: PluginName,
	}
	return plugin
}

func Plugin() *Question {
	once.Do(func() {
		plugin = NewPlugin()
	})
	return plugin
}

func init() {
	auth := Plugin()
	auth.AddHandler(http.MethodPost, "/", create)
	auth.AddHandler(http.MethodGet, "/", find)
	// TODO: add these
	//auth.AddHandler(http.MethodPut, "/:id", edit)
	//auth.AddHandler(http.MethodDelete, "/:id", find)

}

// @Summary list all questions
// @Accept  json
// @Produce  json
// @Router /question/ [get]
// @Tags Question
// @Success 201 {object} FindQuestionsResponse
func find(ctx echo.Context) error {
	if !common.IsAdmin(ctx) {
		return ctx.JSON(http.StatusUnauthorized, FindQuestionsResponse{
			Error: "Unauthorized",
		})
	}
	qs, err := questionService.FindAll()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, FindQuestionsResponse{
			Error: err.Error(),
		})
	}
	return ctx.JSON(http.StatusBadRequest, FindQuestionsResponse{
		Questions: qs,
	})
}

// @Summary create question
// @Accept  json
// @Produce  json
// @Router /question/ [post]
// @Tags Question
// @Param question body CreateQuestionRequest true "create"
// @Success 201 {object} CreateQuestionResponse
func create(ctx echo.Context) error {
	if !common.IsAdmin(ctx) {
		return ctx.JSON(http.StatusUnauthorized, CreateQuestionErrorResponse{
			Error: "Unauthorized",
		})
	}
	var req CreateQuestionRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, CreateQuestionErrorResponse{
			Error: err.Error(),
		})
	}

	q := models.Question{
		ID:            primitive.NewObjectID(),
		Question:      req.Question,
		Answers:       req.Answers,
		CorrectAnswer: req.CorrectAnswer,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	created, err := questionService.Create(q)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, CreateQuestionErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusBadRequest, created)
}

type CreateQuestionRequest struct {
	Question      string   `json:"question"`
	Answers       []string `json:"answers"`
	CorrectAnswer string   `json:"correctAnswer"`
}

type CreateQuestionResponse models.Question
type CreateQuestionErrorResponse struct {
	Error string `json:"error"`
}

type FindQuestionsResponse struct {
	Error     string             `json:"error,omitempty"`
	Questions []*models.Question `json:"questions"`
}

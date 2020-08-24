package auth

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/acha-bill/quizzer_backend/common"
	"github.com/acha-bill/quizzer_backend/models"
	userService "github.com/acha-bill/quizzer_backend/packages/dblayer/user"
	"github.com/acha-bill/quizzer_backend/plugins"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	// PluginName defines the name of the plugin
	PluginName = "auth"
)

var (
	plugin *Auth
	once   sync.Once
)

// Auth structure
type Auth struct {
	name     string
	handlers []*plugins.PluginHandler
}

// AddHandler Method definition from interface
func (plugin *Auth) AddHandler(method string, path string, handler func(echo.Context) error) {
	pluginHandler := &plugins.PluginHandler{
		Path:    path,
		Handler: handler,
		Method:  method,
	}
	plugin.handlers = append(plugin.handlers, pluginHandler)
}

// Handlers Method definition from interface
func (plugin *Auth) Handlers() []*plugins.PluginHandler {
	return plugin.handlers
}

// Name defines the name of the plugin
func (plugin *Auth) Name() string {
	return plugin.name
}

// NewPlugin returns the new plugin
func NewPlugin() *Auth {
	plugin := &Auth{
		name: PluginName,
	}
	return plugin
}

// Plugin returns an instance of the plugin
func Plugin() *Auth {
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
// @Summary Login user
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} LoginResponse
// @Router /auth/login [post]
// @Tags Auth
// @Param login body LoginRequest true "login"
func login(ctx echo.Context) error {
	var req LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, LoginResponse{
			Error: err.Error(),
		})
	}

	filter := bson.D{primitive.E{Key: "username", Value: req.Username}}
	users, err := userService.Find(filter)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, LoginResponse{
			Error: err.Error(),
		})
	}
	if len(users) == 0 {
		return ctx.JSON(http.StatusBadRequest, LoginResponse{
			Error: "user not found",
		})
	}
	u := users[0]

	claims := &common.JWTCustomClaims{
		Username: u.Username,
		IsAdmin:  u.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return ctx.JSON(http.StatusOK, LoginResponse{
		Token: t,
	})
}

// @Summary register user
// @Accept  application/json
// @Produce  application/json
// @Router /auth/register [post]
// @Tags Auth
// @Param register body RegisterRequest true "register"
// @Success 201 {object} RegisterResponse
func register(ctx echo.Context) error {
	var req RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, LoginResponse{
			Error: err.Error(),
		})
	}

	// Basic validation
	if len(req.Username) <= 0 && len(req.Password) <= 0 {
		return ctx.JSON(http.StatusBadRequest, RegisterErrorResponse{
			Error: "Empty values for username and password",
		})
	}

	filter := bson.D{primitive.E{Key: "username", Value: req.Username}}
	users, err := userService.Find(filter)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, RegisterErrorResponse{
			Error: err.Error(),
		})
	}
	if len(users) != 0 {
		return ctx.JSON(http.StatusBadRequest, RegisterErrorResponse{
			Error: "username already taken",
		})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	u := models.User{
		ID:         primitive.NewObjectID(),
		Username:   req.Username,
		Password:   string(hashedPassword),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ProfileURL: "",
		IsAdmin:    req.IsAdmin,
	}
	created, err := userService.Create(u)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, RegisterErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, created)
}

// LoginRequest represents the Request object for Login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the Response object for Login
type LoginResponse struct {
	Error string `json:"error,omitempty"`
	Token string `json:"token,omitempty"`
}

// RegisterRequest represents the Request object for Register
type RegisterRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	ProfileURL string `json:"profileURL"`
	IsAdmin    bool   `json:"isAdmin"`
}

// RegisterErrorResponse represents the Error Response object for Register
type RegisterErrorResponse struct {
	Error string `json:"error,omitempty"`
}

// RegisterResponse represents the Response object for Register
type RegisterResponse models.User

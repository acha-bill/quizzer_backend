package server

import (
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/acha-bill/quizzer_backend/common"
	"github.com/acha-bill/quizzer_backend/plugins"
	"github.com/acha-bill/quizzer_backend/plugins/auth"
	"github.com/acha-bill/quizzer_backend/plugins/question"
	"github.com/acha-bill/quizzer_backend/plugins/search"
	"github.com/gobuffalo/packr/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	once      sync.Once
	server    *echo.Echo
	jwtSecret string
)

var (
	appBox = packr.New("app", "../../frontend/quizzer/build")
	//assetsBox = packr.New("appAssets", "../frontend/quizzer/src/assets")
)

var (
	Plugins = []plugins.Plugin{
		auth.Plugin(),
		question.Plugin(),
		search.Plugin(),
	}
)

const (
	// ApiVersion is the latest version of the API
	ApiVersion = "v1"
)

// Instance creates and returns the echo server instance
func Instance() *echo.Echo {
	once.Do(func() {
		server = instance()
	})
	return server
}

func instance() *echo.Echo {
	// Echo instance
	e := echo.New()
	jwtSecret = os.Getenv("JWT_SECRET")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	//enable debugging
	//debuggingEnabled, err := strconv.ParseBool(os.Getenv("DEBUGGING_ENABLED"))
	//if debuggingEnabled && err == nil {
	//	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	//		Format: "method=${method}, uri=${uri}, status=${status}\n",
	//	}))
	//}

	if !common.IsDevelopment() {
		for _, res := range appBox.List() {
			e.GET("/"+res, fromFile)
		}
	}

	// index route
	e.GET("/", indexRoute)

	// Plugin Routes
	for _, plugin := range Plugins {
		for _, handler := range plugin.Handlers() {
			path := "api/" + ApiVersion + "/" + plugin.Name() + handler.Path
			e.Add(handler.Method, path, handler.Handler, middleware.JWTWithConfig(middleware.JWTConfig{
				Skipper: func(ctx echo.Context) bool {
					return handler.AuthLevel == plugins.AuthLevelNone
				},
				Claims:     &common.JWTCustomClaims{},
				SigningKey: []byte(jwtSecret),
			}))
		}
	}

	// other routes
	e.GET("*", indexRoute)
	return e
}

func fromFile(e echo.Context) error {
	reqPath := e.Request().URL.Path
	bytes, err := appBox.Find(reqPath)
	if err != nil {
		return err
	}

	parts := strings.Split(reqPath, ".")
	ext := ""
	if len(parts) >= 0 {
		ext = "." + parts[len(parts)-1]
	}
	mimeType := mime.TypeByExtension(ext)

	return e.Blob(http.StatusOK, mimeType, bytes)
}

func indexRoute(e echo.Context) error {
	if !common.IsDevelopment() {
		indexHTML, err := appBox.Find("index.html")
		if err != nil {
			return err
		}
		return e.HTMLBlob(http.StatusOK, indexHTML)
	}
	res, err := http.Get("http://127.0.0.1:3000" + e.Request().URL.Path)
	if err != nil {
		return err
	}
	devIndexHTML, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return e.Blob(http.StatusOK, res.Header.Get("Content-Type"), devIndexHTML)

}

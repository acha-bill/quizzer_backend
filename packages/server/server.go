package server

import (
	"fmt"
	"github.com/acha-bill/quizzer_backend/common"
	"github.com/acha-bill/quizzer_backend/plugins"
	"github.com/acha-bill/quizzer_backend/plugins/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	once sync.Once
	server *echo.Echo
	jwtSecret string
)

var (
	Plugins = []plugins.Plugin {
		auth.Plugin(),
	}
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
	debuggingEnabled, err := strconv.ParseBool(os.Getenv("DEBUGGING_ENABLED"))
	if  debuggingEnabled && err == nil {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}\n",
		}))
	}


	// Routes
	for _, plugin := range Plugins {
		fmt.Println(plugin)
		for _, handler := range plugin.Handlers() {
			path := plugin.Name() + handler.Path
			log.Info(path)
			e.Add(handler.Method, path , handler.Handler,middleware.JWTWithConfig(middleware.JWTConfig{
				Skipper: func(ctx echo.Context) bool {
					log.Info(ctx.Path())
					return strings.HasPrefix(ctx.Path(), "/auth")
				},
				Claims:     &common.JWTCustomClaims{},
				SigningKey: []byte(jwtSecret),
			}))
		}
	}

	log.Info(e.Routes())

	return e
}

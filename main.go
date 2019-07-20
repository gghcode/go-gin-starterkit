package main

import (
	"github.com/defval/inject"
	"github.com/gghcode/go-gin-starterkit/api"
	"github.com/gghcode/go-gin-starterkit/api/auth"
	"github.com/gghcode/go-gin-starterkit/api/common"
	"github.com/gghcode/go-gin-starterkit/api/todo"
	"github.com/gghcode/go-gin-starterkit/api/user"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"
	_ "github.com/gghcode/go-gin-starterkit/docs"
	"github.com/gghcode/go-gin-starterkit/middleware"
	"github.com/gghcode/go-gin-starterkit/service"
	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

const (
	envPrefix = "REST"
)

// @title Go Gin Starter API
// @version 1.0
// @description This is a sample gin starter server.
// @termsOfService http://swagger.io/terms/
// @securitydefinitions.apikey ApiKeyAuth
// @BasePath /api
// @in header
// @name Authorization
// @contact.name API Support
// @contact.email gyuhwan.a.kim@gmail.com
// @license.name MIT
// @license.url https://github.com/gghcode/go-gin-starterkit/blob/master/LICENSE
func main() {
	container, err := inject.New(
		inject.Provide(config.NewBuilder().
			AddConfigFile("config.yaml", true).
			BindEnvs(envPrefix).
			Build),

		inject.Provide(db.NewConn),
		inject.Provide(service.NewPassport),

		inject.Provide(common.NewController, inject.As(api.IController)),
		inject.Provide(user.NewRepository),
		inject.Provide(user.NewController, inject.As(api.IController)),

		inject.Provide(todo.NewRepository),
		inject.Provide(todo.NewController, inject.As(api.IController)),

		inject.Provide(auth.NewService),
		inject.Provide(auth.NewController, inject.As(api.IController)),
	)

	if err != nil {
		panic(err)
	}

	var conf config.Configuration
	if err := container.Extract(&conf); err != nil {
		panic(err)
	}

	var controllers []api.Controller
	if err := container.Extract(&controllers); err != nil {
		panic(err)
	}

	router := gin.New()
	router.Use(middleware.AddAuthHandler(conf.Jwt))
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiRouter := router.Group("api/")
	for _, controller := range controllers {
		controller.RegisterRoutes(apiRouter)
	}

	if err := router.Run(conf.Addr); err != nil {
		panic(err)
	}
}

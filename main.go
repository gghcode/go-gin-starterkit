package main

import (
	"github.com/gghcode/go-gin-starterkit/app"
	"github.com/gghcode/go-gin-starterkit/config"
)

const (
	envPrefix = "REST"
)

// @title Go Gin Starter API
// @version 1.0
// @description This is a sample gin starter server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email gyuhwan.a.kim@gmail.com

// @license.name MIT
// @license.url https://github.com/gghcode/go-gin-starterkit/blob/master/LICENSE
func main() {
	conf, err := config.NewBuilder().
		AddConfigFile("config.yaml", true).
		BindEnvs(envPrefix).
		Build()

	if err != nil {
		panic(err)
	}

	server := app.NewServer(conf)
	if err := server.Run(); err != nil {
		panic(err)
	}
}

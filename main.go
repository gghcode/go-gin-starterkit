package main

import (
	"github.com/gyuhwankim/go-gin-starterkit/app"
	"github.com/gyuhwankim/go-gin-starterkit/config"
)

const (
	envPrefix = "REST"
)

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

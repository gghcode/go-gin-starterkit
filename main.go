package main

import (
	"fmt"
	"github.com/pkg/errors"
	"golang-webapi-starterkit/api"
	"golang-webapi-starterkit/app"
	"golang-webapi-starterkit/config"
	"golang-webapi-starterkit/todo"
)

func main() {
	builder := config.NewViperBuilder()
	builder.BasePath(".")
	builder.JsonFile("config")
	//builder.AddEnvironmentVariables()

	configuration, err := builder.Build()
	if err != nil {
		panic(errors.Wrap(err, "Configuration build failed."))
	}

	controllers := []api.Controller{
		todo.NewController(configuration),
	}

	server := app.NewServer(configuration, controllers)
	if err := server.Run(); err != nil {
		panic(err)
	}

	fmt.Println("This is golang webapi starterkit.")
}

package main

import (
	"fmt"
	"github.com/pkg/errors"
	"golang-webapi-starterkit/config"
)

func main() {
	builder := config.NewViperBuilder()
	builder.SetBasePath(".")
	builder.AddJsonFile("config")
	//builder.AddEnvironmentVariables()

	_, err := builder.Build()
	if err != nil {
		panic(errors.Wrap(err, "Configuration build failed."))
	}

	fmt.Println("This is golang webapi starterkit...")
}

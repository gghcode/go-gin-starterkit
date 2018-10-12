package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-webapi-starterkit/api"
	"golang-webapi-starterkit/config"
)

// ApiServer is server that Application level
type ApiServer struct {
	engine        *gin.Engine
	configuration config.Configuration
	controllers   []api.Controller
}

func NewServer(configuration config.Configuration, controllers []api.Controller) *ApiServer {
	server := ApiServer{
		engine:        gin.New(),
		configuration: configuration,
		controllers:   controllers,
	}

	return &server
}

func (apiServer *ApiServer) Run() error {
	api.RegisterControllers(apiServer.engine, apiServer.controllers)

	listenPort := apiServer.configuration.ListenPort
	listenAddr := getAddrString(listenPort)

	return apiServer.engine.Run(listenAddr)
}

func getAddrString(port int) string {
	return fmt.Sprintf(":%d", port)
}

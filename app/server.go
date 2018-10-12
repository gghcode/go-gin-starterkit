package app

import (
	"fmt"
	"golang-webapi-starterkit/config"
)

// ApiServer is server that Application level
type ApiServer struct {
	engine        ServerEngine
	configuration config.Configuration
	controllers   []ApiController
}

func newServer(configuration config.Configuration) *ApiServer {
	return &ApiServer{
		configuration: configuration,
	}
}

func (apiServer *ApiServer) Initialize() {
	engine := NewGinEngine()

	RegisterControllers(engine, apiServer.controllers)

	apiServer.engine = engine
}

func (apiServer *ApiServer) Run() error {
	listenPort := apiServer.configuration.ListenPort
	listenAddr := getAddrString(listenPort)

	return apiServer.engine.Run(listenAddr)
}

func getAddrString(port int) string {
	return fmt.Sprintf(":%d", port)
}

package app

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"golang-webapi-starterkit/config"
	"golang-webapi-starterkit/todo"
)

type ApiServer struct {
	engine ServerEngine
	configuration config.Configuration
	controllers   []ApiController
	logger        *logrus.Entry
}

func NewServer(configuration config.Configuration) *ApiServer {
	logger := logrus.New().WithField("host", "server")
	controllers := []ApiController{
		todo.NewController(configuration),
	}

	return &ApiServer{
		configuration: configuration,
		controllers:   controllers,
		logger:        logger,
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

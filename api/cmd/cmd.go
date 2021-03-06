package main

import (
	"fmt"
	"net/http"

	"api/server"
	"api/util/configs"
)

func main() {
	// Get configs of db, server and any other services
	config, err := configs.LoadConfigs()
	if err != nil {
		panic(server.ConfigsLoadErr)
	}

	// pass configs to server connection
	srv := server.New(config)
	appServer := server.InitializeServerRoutes(srv)
	service := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.ServerConf.Port),
		Handler: appServer,
	}

	if err := service.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(server.StartupErr)
	}
}

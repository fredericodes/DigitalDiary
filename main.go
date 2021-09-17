package main

import (
	"fmt"
	"net/http"

	"diarynote/server"
	"diarynote/util/configs"
)

func main() {
	// Get configs of db, server and any other services
	configs, err := configs.LoadConfigs()
	if err != nil {
		panic(server.ConfigsLoadErr)
	}

	// pass configs to server connection
	srv := server.New(configs)
	appServer := server.InitializeServerRoutes(srv)
	service := &http.Server{
		Addr:    fmt.Sprintf(":%d", configs.ServerConf.Port),
		Handler: appServer,
	}

	if err := service.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(server.StartupErr)
	}
}

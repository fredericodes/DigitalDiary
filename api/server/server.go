package server

import (
	"api/repository"
	"api/util/configs"
)

const (
	ConfigsLoadErr = "configs could not be loaded"
	StartupErr     = "server could not start up"

	errFormDecodingErr           = "form decoding failure"
	errPasswordHashGenFailure    = "password could not be hashed"
	errServerProcessingErr       = "service is facing issue processing the request"
	errUserEmailExistsErr        = "user email exists, try login as user"
	errUserDoesntExistsErr  = "user doesn't exists, register as user"
	errInvalidCredentials        = "invalid credentials are provided"
	errQueryParamDoesntExistsErr = "query params not provided"
)

type Server struct {
	DB repository.DB
}

func New(configs *configs.Configs) *Server {
	return &Server{
		DB: repository.New(configs.DbConf),
	}
}

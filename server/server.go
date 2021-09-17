package server

import (
	"fmt"
	"net/http"

	"diarynote/repositories"
	"diarynote/util/configs"
	"diarynote/util/logger"
)

const (
	ConfigsLoadErr = "configs could not be loaded"
	StartupErr     = "server could not start up"

	errFormDecodingErr            = "form decoding failure"
	errPasswordHashGenFailure     = "password could not be hashed"
	errServerProcessingErr        = "service is facing issue processing the request"
	errUserEmailExistsErr         = "user email exists, try login as user"
	errUserEmailDoesntExistsErr   = "user email doesn't exists, register as user"
	errInvalidCredentials         = "invalid credentials are provided"
	errNotAuthorizedToViewContent = "not authorized to view content"
)

type Server struct {
	RedisDb *repositories.Db
	Logger  *logger.Logger
}

func New(configs *configs.Configs) *Server {
	return &Server{
		RedisDb: repositories.New(configs.RedisDbConf),
		Logger:  logger.New(),
	}
}

func requestLogInfo(r *http.Request, reqBody []byte) string {
	return fmt.Sprintf("%v %v %v %v\n%v\n%v\n", r.Host, r.RemoteAddr, r.Method, r.URL, r.Header, string(reqBody))
}

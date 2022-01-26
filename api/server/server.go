package server

import (
	"fmt"
	"net/http"

	"github.com/FreddyJilson/diarynote/repository"
	"github.com/FreddyJilson/diarynote/util/configs"
	"github.com/FreddyJilson/diarynote/util/logger"
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
)

type Server struct {
	RedisDb *repository.Db
	Logger  *logger.Logger
}

func New(configs *configs.Configs) *Server {
	return &Server{
		RedisDb: repository.New(configs.RedisDbConf),
		Logger:  logger.New(),
	}
}

func requestLogInfo(r *http.Request, reqBody []byte) string {
	return fmt.Sprintf("%v %v %v %v\n%v\n%v\n", r.Host, r.RemoteAddr, r.Method, r.URL, r.Header, string(reqBody))
}

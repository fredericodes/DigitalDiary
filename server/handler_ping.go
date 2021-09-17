package server

import (
	"net/http"
)

func (srv *Server) HandleAuthPing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

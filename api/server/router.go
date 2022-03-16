package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func InitializeServerRoutes(srv *Server) http.Handler {
	router := mux.NewRouter()

	handlersCors := cors.New(cors.Options{
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "pragma", "X-Organization"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
		ExposedHeaders:   []string{"X-Total-Count"},
		MaxAge:           300,
	})

	// Authentication routes
	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", srv.HandleLogin).Methods(http.MethodPost)
	auth.HandleFunc("/register", srv.HandleRegister).Methods(http.MethodPost)

	// Api v1 routes
	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/journal", Authenticate(srv.HandleListJournal)).Methods(http.MethodGet)
	apiV1.HandleFunc("/journal", Authenticate(srv.HandleCreateOrUpdateJournal)).Methods(http.MethodPost)

	apiV1.HandleFunc("/server-auth-ping", Authenticate(srv.HandleAuthPing)).Methods(http.MethodGet)

	appRouter := handlersCors.Handler(router)
	return appRouter
}

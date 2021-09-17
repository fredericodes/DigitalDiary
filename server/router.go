package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	PermissionTypeUser  = "user"
	PermissionTypeAdmin = "admin"

	Auth                 = "Auth"
	HandleLogin          = "HandleLogin"
	HandleRegister       = "HandleRegister"
	HandleChangePassword = "HandleChangePassword"

	ApiV1                             = "ApiV1"
	HandleListDiaryJournals           = "HandleListDiaryJournals"
	HandleCreateOrUpdateDiaryJournals = "HandleCreateOrUpdateDiaryJournals"
	HandleAuthPing                    = "HandleAuthPing"
)

type Routes struct {
	Path                string
	Policy              string
	AuthorizedUserRoles map[string]string
}

func RoutesMap() map[string]Routes {
	var routesMap = map[string]Routes{
		Auth: {
			"/auth",
			"",
			nil,
		},
		HandleLogin: {
			"/login",
			http.MethodPost,
			map[string]string{
				PermissionTypeUser:  PermissionTypeUser,
				PermissionTypeAdmin: PermissionTypeAdmin,
			},
		},
		HandleRegister: {
			"/register",
			http.MethodPost,
			map[string]string{
				PermissionTypeUser:  PermissionTypeUser,
				PermissionTypeAdmin: PermissionTypeAdmin,
			},
		},
		HandleChangePassword: {
			"/change-password",
			http.MethodPut,
			map[string]string{
				PermissionTypeUser:  PermissionTypeUser,
				PermissionTypeAdmin: PermissionTypeAdmin,
			},
		},
		ApiV1: {
			"/api/v1",
			"",
			nil,
		},
		HandleListDiaryJournals: {
			"/journal-entries",
			http.MethodGet,
			map[string]string{
				PermissionTypeUser:  PermissionTypeUser,
				PermissionTypeAdmin: PermissionTypeAdmin,
			},
		},
		HandleCreateOrUpdateDiaryJournals: {
			"/journal-entries",
			http.MethodPost,
			map[string]string{
				PermissionTypeUser:  PermissionTypeUser,
				PermissionTypeAdmin: PermissionTypeAdmin,
			},
		},
		HandleAuthPing: {
			"/server-auth-ping",
			http.MethodGet,
			map[string]string{
				PermissionTypeUser:  PermissionTypeUser,
				PermissionTypeAdmin: PermissionTypeAdmin,
			},
		},
	}

	return routesMap
}

func IsAuthorized(handlerName string, role string) bool {
	routesMap := RoutesMap()
	_, exists := routesMap[handlerName].AuthorizedUserRoles[role]
	if !exists {
		return false
	}

	return true
}

func InitializeServerRoutes(srv *Server) http.Handler {
	router := mux.NewRouter()

	routesMap := RoutesMap()

	handlersCors := cors.New(cors.Options{
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "pragma", "X-Organization"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
		ExposedHeaders:   []string{"X-Total-Count"},
		MaxAge:           300,
	})

	// Authentication routes
	auth := router.PathPrefix(routesMap[Auth].Path).Subrouter()
	auth.HandleFunc(routesMap[HandleLogin].Path, srv.HandleLogin).Methods(routesMap[HandleLogin].Policy)
	auth.HandleFunc(routesMap[HandleRegister].Path, srv.HandleRegister).Methods(routesMap[HandleRegister].Policy)
	auth.HandleFunc(routesMap[HandleChangePassword].Path, Authenticate(srv.HandleChangePassword)).Methods(routesMap[HandleChangePassword].Policy)

	// Api v1 routes
	apiV1 := router.PathPrefix(routesMap[ApiV1].Path).Subrouter()
	apiV1.HandleFunc(routesMap[HandleListDiaryJournals].Path, Authenticate(srv.HandleListDiaryJournals)).Methods(routesMap[HandleListDiaryJournals].Policy)
	apiV1.HandleFunc(routesMap[HandleCreateOrUpdateDiaryJournals].Path, Authenticate(srv.HandleCreateOrUpdateDiaryJournals)).Methods(routesMap[HandleCreateOrUpdateDiaryJournals].Policy)

	apiV1.HandleFunc(routesMap[HandleAuthPing].Path, Authenticate(srv.HandleAuthPing)).Methods(routesMap[HandleAuthPing].Policy)

	appRouter := handlersCors.Handler(router)
	return appRouter
}

package routes

import (
	"github.com/alibaba0010/postgres-api/internal/controllers"
	"github.com/gorilla/mux"
)

func AuthRoutes() *mux.Router {
	route := mux.NewRouter()
	route.HandleFunc("/signup", controllers.SignupHandler).Methods("POST")
	route.HandleFunc("/verify", controllers.ActivateUserHandler).Methods("GET")
	// route.HandleFunc("/signin", controllers.SigninHandler).Methods("POST")

	return route
}
// @title Book API
// @version 1.0
// @description This is a sample REST API built with Go, mux, zap, and bun.
// @termsOfService http://swagger.io/terms/
// @contact.name Ali Baba
// @contact.email ali@example.com
// @license.name MIT
// @BasePath /api/v1
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// "path/filepath"

	"github.com/alibaba0010/postgres-api/api/database"
	"github.com/alibaba0010/postgres-api/api/errors"
	_ "github.com/alibaba0010/postgres-api/docs" // swag doc
	"github.com/alibaba0010/postgres-api/logger"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)


func main(){
	err := godotenv.Load()
	if err != nil {
		// log.Println("No .env file found, using system environment variables...")
		logger.Log.Warn("No .env file found", zap.Error(err))
	}
	port := os.Getenv("PORT")
	if port == "" {
		port ="3001"
		fmt.Println("Add port to env file")
	}
	logger.InitLogger()
	// defer sync to flush logs on program exit
	defer logger.Sync()

	database.ConnectDB()
	defer database.CloseDB()

	route := mux.NewRouter()
	// Add recovery middleware early so panics are caught and do not print stack traces.
	route.Use(errors.RecoverMiddleware)
	route.Use(logger.Logger)
	
	// Serve Swagger UI at /swagger/
	route.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	// API routes
	// API routes
	route.HandleFunc("/getUser", getUserHandler).Methods("GET")
	route.HandleFunc("/getBook", GetBookHandler).Methods("GET")
	route.HandleFunc("/", httpHandler).Methods("GET")

	// Healthcheck
	route.HandleFunc("/api/v1/healthcheck", HealthCheckHandler).Methods("GET")

	route.NotFoundHandler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		errors.ErrorResponse(writer, request, errors.RouteNotExist())
	})
	logger.Log.Info("🚀 Server starting", zap.String("address", ":"+port))
	if  err:= http.ListenAndServe(":"+port, route); err != nil {
		log.Fatal(err)
	}
}

// http.HandleFunc("/getUser", getUserHandler)
// 	http.HandleFunc("/getBook", GetBookHandler)
// 	http.HandleFunc("/", httpHandler)
// if  err:= http.ListenAndServe(":"+port, nil); err != nil {

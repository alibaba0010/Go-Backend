package main

import (
	"log"
	"net/http"

	"github.com/alibaba0010/postgres-api/api/config"
	"github.com/alibaba0010/postgres-api/api/database"
	_ "github.com/alibaba0010/postgres-api/docs" // swag doc
	"github.com/alibaba0010/postgres-api/logger"
	"github.com/alibaba0010/postgres-api/api/routes"
	"go.uber.org/zap"
)


func main(){

	logger.InitLogger()
	// defer sync to flush logs on program exit
	defer logger.Sync()

	database.ConnectDB()
	defer database.CloseDB()
	port := config.LoadConfig().Port
	route := routes.ApiRouter()

	
	logger.Log.Info("🚀 Server starting", zap.String("address", ":"+port))
	if  err:= http.ListenAndServe(":"+port, route); err != nil {
		log.Fatal(err)
	}
}

// http.HandleFunc("/getUser", getUserHandler)
// 	http.HandleFunc("/getBook", GetBookHandler)
// 	http.HandleFunc("/", httpHandler)
// if  err:= http.ListenAndServe(":"+port, nil); err != nil {

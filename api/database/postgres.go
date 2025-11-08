package database

import (
	"context"
	"fmt"
	"time"

	"github.com/alibaba0010/postgres-api/logger"
	"github.com/alibaba0010/postgres-api/api/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)
var Pool *pgxpool.Pool
func ConnectDB(){
	cfg:= config.LoadConfig()
	host:= cfg.DB_HOST
	port:= cfg.DB_PORT
	user:= cfg.DB_USERNAME
	password:= cfg.DB_PASSWORD
	dbname:= cfg.DB_NAME

	connectionURL:= fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	context, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	var err error
	Pool, err = pgxpool.New(context, connectionURL)
	if err != nil {
		logger.Log.Fatal("Unable to connect to database", zap.Error(err))
	}
	// db := bun.NewDB(Pool, pgdialect.New())

	
		// Test connection
	err = Pool.Ping(context)
	if err != nil {
		logger.Log.Fatal("Database ping failed", zap.Error(err))
	}

	logger.Log.Info("✅ Connected to PostgreSQL database")

}

// Close connection when shutting down
func CloseDB() {
	if Pool != nil {
		Pool.Close()
		logger.Log.Info("🔌 Database connection closed")
	}
}
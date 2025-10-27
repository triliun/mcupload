package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/triliun/mcupload/backend/config"
	"github.com/triliun/mcupload/backend/database"
	"github.com/triliun/mcupload/backend/logger"
	"go.uber.org/zap"
	_ "golang.org/x/time/rate"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Initialize logger
	loggerConfig := logger.NewProductionConfig()
	if os.Getenv("ENV") == "development" {
		loggerConfig = logger.DefaultConfig()
	}

	err = logger.Init(loggerConfig)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	var db *sqlx.DB
	db, err = database.NewPostgresConnection(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		logger.Fatal("Database connection failed", zap.Error(err))
	}
	defer db.Close()

	router := config.SetupRouter(db)

	server := http.Server{
		Addr:    os.Getenv("SERVER_ADDRESS"),
		Handler: router,
	}

	go func() {
		logger.Info("Starting server", zap.String("address", server.Addr))
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancle := context.WithTimeout(context.Background(), time.Second*30)
	defer cancle()

	err = server.Shutdown(ctx)
	if err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

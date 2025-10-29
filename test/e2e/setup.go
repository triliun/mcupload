package e2e

import (
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	"github.com/triliun/mcupload/backend/config"
	"github.com/triliun/mcupload/backend/database"
	"github.com/triliun/mcupload/backend/logger"
	"github.com/triliun/mcupload/backend/shared"
	"go.uber.org/zap"
)

var router *chi.Mux

func setup() (*sqlx.DB, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}

	// Initialize logger
	loggerConfig := logger.NewProductionConfig()
	if os.Getenv("ENV") == "development" {
		loggerConfig = logger.DefaultConfig()
	}

	err = logger.Init(loggerConfig)
	if err != nil {
		return nil, err
	}

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
		return nil, err
	}
	// defer db.Close()

	err = cleanDatabase(db)
	if err != nil {
		logger.Error("Failed to cleanup database", zap.Error(err))
		os.Exit(1)
	}

	// create user account with role ceo
	err = createUserWithRoleCEO(db)
	if err != nil {
		return nil, err
	}

	// Setup dependencies
	router = config.SetupRouter(db)

	// Setup cron job
	c := cron.New()
	_, err = c.AddFunc("0 2 * * *", func() {
		// Implement cleanup if needed, but avoid in test environment
		logger.Info("Cleanup job skipped in test environment")
	})
	if err != nil {
		logger.Error("Failed to schedule cleanup job", zap.Error(err))
	}
	c.Start()

	logger.Info("E2E Test setup completed")

	return db, nil
}

// GetRouter returns the router for use in tests
func GetRouter() *chi.Mux {
	return router
}

func cleanup() {
	logger.Info("Cleaning up after tests")
	logger.Sync()
	// db.Close()
}

func createUserWithRoleCEO(db *sqlx.DB) error {
	user := UserWithRoleCEO

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	now := time.Now()

	user.ID = id
	user.CreatedAt = now
	user.UpdatedAt = now

	user.Password, err = shared.Crypto.HashPassword(user.Password)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO users (
		id, username, email, password, role, created_at, updated_at
	) VALUES (
		:id, :username, :email, :password, :role, :created_at, :updated_at
	)
	`

	_, err = db.NamedExec(query, &user)

	return err
}

func cleanDatabase(db *sqlx.DB) error {
	// Hapus semua data sebelum test dijalankan
	var err error

	tables := []string{
		"users",
		"minecraft_versions",
		"categories",
		"resource_packs",
		"pack_categories",
		"pack_files",
		"revoked_tokens",
	}

	for _, table := range tables {
		_, err = db.Exec("DELETE FROM " + table)
		if err != nil {
			return err
		}
	}

	return err
}

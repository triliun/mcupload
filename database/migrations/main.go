package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/triliun/mcupload/database"
	"github.com/triliun/mcupload/logger"
	"go.uber.org/zap"
)

// # Jalankan semua migrations (UP)
// go run main.go up

// # atau
// go run main.go

// # Rollback migration tertentu
// go run main.go down 005_create_resource_packs_table

// # Lihat status migrations
// psql -U your_user -d your_db -c "SELECT * FROM migrations ORDER BY applied_at;"

// Migration represents a database migrations
type Migration struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	AppliedAt time.Time `db:"applied_at"`
}

const MIGRATIONS_FILE_DIR string = "sql"

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	// Initialize logger
	config := logger.NewProductionConfig()
	if os.Getenv("ENV") == "development" {
		config = logger.DefaultConfig()
	}

	err = logger.Init(config)
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

	fmt.Println("Successfully connected to database!")

	// Create migrations table if it doesn't exist
	err = createMigrationsTable(db)
	if err != nil {
		log.Fatal("Error creating migrations table:", err)
	}

	// Check command line arguments for migration direction
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "up":
			err = runMigrations(db, "up")
		case "down":
			if len(os.Args) > 2 {
				err = runDownMigration(db, os.Args[2])
			} else {
				logger.Fatal("Please specify migration name to rollback")
			}
		default:
			logger.Fatal("Usage: go run main.go [up|down <migration_name>]")
		}
	} else {
		// Default: run up migrations
		err = runMigrations(db, "up")
	}

	if err != nil {
		logger.Fatal("Error running migrations", zap.Error(err))
	}

	fmt.Println("Migrations completed successfully!")
}

func createMigrationsTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`
	_, err := db.Exec(query)
	return err
}

func runMigrations(db *sqlx.DB, direction string) error {
	// Baca semua file migration dari folder
	migrationFiles, err := readMigrationFiles(MIGRATIONS_FILE_DIR, direction)
	if err != nil {
		return err
	}

	// Urutkan migrations berdasarkan nama file
	sort.Strings(migrationFiles)

	for _, migrationFile := range migrationFiles {
		migrationName := strings.TrimSuffix(filepath.Base(migrationFile), "."+direction+".sql")

		if direction == "up" {
			// Check if migration has already been applied
			var existingMigration Migration
			err := db.Get(&existingMigration, "SELECT id, name, applied_at FROM migrations WHERE name = $1", migrationName)

			if err == sql.ErrNoRows {
				// Migration not applied yet, run it
				fmt.Printf("Applying migration: %s\n", migrationName)

				// Baca file SQL
				query, err := ioutil.ReadFile(migrationFile)
				if err != nil {
					return fmt.Errorf("error reading migration file %s: %v", migrationFile, err)
				}

				tx, err := db.Beginx()
				if err != nil {
					return fmt.Errorf("error starting transaction for migration %s: %v", migrationName, err)
				}

				// Execute migration query dari file
				_, err = tx.Exec(string(query))
				if err != nil {
					tx.Rollback()
					return fmt.Errorf("error executing migration %s: %v", migrationName, err)
				}

				// Record migration in migrations table
				_, err = tx.Exec("INSERT INTO migrations (name) VALUES ($1)", migrationName)
				if err != nil {
					tx.Rollback()
					return fmt.Errorf("error recording migration %s: %v", migrationName, err)
				}

				// Commit transaction
				err = tx.Commit()
				if err != nil {
					return fmt.Errorf("error committing migration %s: %v", migrationName, err)
				}

				fmt.Printf("✅ Successfully applied migration: %s\n", migrationName)
			} else if err != nil {
				return fmt.Errorf("error checking migration %s: %v", migrationName, err)
			} else {
				fmt.Printf("⏩ Migration already applied: %s (applied at: %v)\n", migrationName, existingMigration.AppliedAt)
			}
		}
	}

	return nil
}

func runDownMigration(db *sqlx.DB, migrationName string) error {
	downFile := fmt.Sprintf("%s/%s.down.sql", MIGRATIONS_FILE_DIR, migrationName)

	// Check if down file exists
	if _, err := os.Stat(downFile); os.IsNotExist(err) {
		return fmt.Errorf("down migration file not found for %s", migrationName)
	}

	// Baca file down SQL
	query, err := ioutil.ReadFile(downFile)
	if err != nil {
		return fmt.Errorf("error reading down migration file %s: %v", downFile, err)
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Execute down migration
	_, err = tx.Exec(string(query))
	if err != nil {
		tx.Rollback()
		return err
	}

	// Hapus record migration
	_, err = tx.Exec("DELETE FROM migrations WHERE name = $1", migrationName)
	if err != nil {
		tx.Rollback()
		return err
	}

	fmt.Printf("✅ Successfully rolled back migration: %s\n", migrationName)
	return tx.Commit()
}

func readMigrationFiles(dir string, direction string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "."+direction+".sql") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

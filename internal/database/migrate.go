package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// RunMigrations runs all pending database migrations from the migrations directory.
// It's safe to run multiple times - already applied migrations are skipped.
func RunMigrations(cfg Config) error {
	// Build connection string
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	// Open database connection using database/sql (required by migrate)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	// Create migrate instance with file source
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run migrations
	slog.Info("Running database migrations...")
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			slog.Info("Database schema is up to date")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	slog.Info("Database migrations completed successfully")
	return nil
}

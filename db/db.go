package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var RM *sql.DB

func ConnectionAndMigrate() error {
	err := godotenv.Load("app.env")
	if err != nil {
		return fmt.Errorf("error loading .env: %w", err)
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening DB: %w", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("error pinging DB: %w", err)
	}

	RM = DB

	if err := MigrateUp(DB); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	log.Println("Database connected and migrated successfully.")
	return nil
}

func MigrateUp(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migration",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func ShutDownDBN() error {
	if RM != nil {
		return RM.Close()
	}
	return nil
}
func Tx(tx *sql.Tx, err *error) {
	if r := recover(); r != nil {
		_ = tx.Rollback()
		panic(r)
	} else if err != nil && *err != nil {
		_ = tx.Rollback()
	} else {
		_ = tx.Commit()
	}
}

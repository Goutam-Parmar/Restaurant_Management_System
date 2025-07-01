package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"
	"os"
)

var RM *sql.DB

func ConnectionAndMigrate() error {
	err := godotenv.Load()
	if err != nil {
		return err
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
		return err
	}
	err = DB.Ping()
	if err != nil {
		return err
	}
	RM = DB
	return MigrateUp(DB)
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

func Tx(tx *sql.Tx, err *error) {
	if r := recover(); r != nil {
		_ = tx.Rollback()
		panic(r)
	} else if err != nil {
		_ = tx.Rollback()
	} else {
		_ = tx.Commit()
	}
}

func ShutDownDBN() error {
	return RM.Close()
}

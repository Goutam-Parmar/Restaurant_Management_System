package main

import (
	"RMS/Routes"
	"RMS/db"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatal("Error loading app.env file:", err)
	}

	dbConn, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	))
	if err != nil {
		log.Fatal(" Database connection failed:", err)
	}
	defer dbConn.Close()

	// Ping database
	if err := dbConn.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}
	fmt.Println("Connected to the database successfully!")

	// Run migrations
	if err := db.MigrateUp(dbConn); err != nil {
		log.Fatal("Migration failed:", err)
	}
	fmt.Println("Migration successful!")

	// Initialize all routes using gorilla/mux
	router := Routes.InitRoutes(dbConn)

	// Start server
	fmt.Println("Server running at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", router))

}

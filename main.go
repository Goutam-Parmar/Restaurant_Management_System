package main

import (
	"RMS/Routes"
	"RMS/db"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	err := db.ConnectionAndMigrate()
	if err != nil {
		log.Fatal(" Database connection failed:", err)
	}
	db.ShutDownDBN()

	router := Routes.InitRoutes()

	// Start server
	fmt.Println("Server running at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", router))

}

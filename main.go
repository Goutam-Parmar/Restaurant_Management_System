package main

import (
	"RMS/Routes"
	"RMS/db"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutdownTimeout = 10 * time.Second

func main() {

	err := db.ConnectionAndMigrate()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	} else {
		log.Println("Database connection succeeded")
	}

	router := Routes.InitRoutes()
	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	//  goroutine for faster
	go func() {
		fmt.Println("Server running at http://localhost:8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	//  listen signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received...")

	// set contex time shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	if err := db.ShutDownDBN(); err != nil {
		log.Printf("Error closing DB: %v", err)
	}

	log.Println("Server stop gracefully")
}

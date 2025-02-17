package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vault-app/handler"
	"vault-app/repository"
	"vault-app/usecase"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {

	viper.AutomaticEnv()

	// 3. Build the connection string
	// e.g. host=localhost port=5432 user=postgres password=mysecret dbname=mydb sslmode=disable
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_SSLMODE"),
	)

	// 4. Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	animalRepo := repository.NewrepositoryAnimal(db)
	animalUcase := usecase.NewusecaseAnimal(&animalRepo)

	router := mux.NewRouter()

	router.Use(loggingMiddleware)

	handler.NewHandlerAnimal(router, &animalUcase)
	// 3. Create an HTTP server, specifying the mux as the handler
	srv := &http.Server{
		Addr:    ":8080", // or any port you prefer
		Handler: router,
	}

	// 4. Start the server in a goroutine so that it doesn’t block
	go func() {
		log.Printf("Starting server on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// 5. Set up channel on which to send signal notifications.
	quit := make(chan os.Signal, 1)
	// SIGINT (Ctrl+C) or SIGTERM is typically used to request a graceful shutdown
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 6. Block until we receive our signal
	<-quit
	log.Println("Shutdown signal received...")

	// 7. Create a deadline to wait for existing connections to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 8. Attempt the graceful shutdown
	log.Println("Shutting down server gracefully...")
	if err := srv.Shutdown(ctx); err != nil {
		// Forced shutdown if graceful shutdown fails
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.RequestURI)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

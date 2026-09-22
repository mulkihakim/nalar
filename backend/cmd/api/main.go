package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mulkihakim/nalar/backend/internal/db"
	"github.com/mulkihakim/nalar/backend/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	gormDB, err := db.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	userRepo := user.NewRepository(gormDB)
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		userHandler.RegisterRoutes(r)
	})

	log.Println("server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
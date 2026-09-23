package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/mulkihakim/nalar/backend/internal/db"
	"github.com/mulkihakim/nalar/backend/internal/middleware"
	"github.com/mulkihakim/nalar/backend/internal/user"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "nalar-default-secret-change-in-production"
	}

	gormDB, err := db.Connect(db.GetDSN())
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	userRepo := user.NewRepository(gormDB)
	userSvc := user.NewService(userRepo, []byte(jwtSecret))

	authMiddleware := middleware.Auth([]byte(jwtSecret))
	// Rate limit login: maks 5 percobaan per menit per IP
	loginRateLimiter := middleware.RateLimit(5, time.Minute)

	userHandler := user.NewHandler(userSvc, authMiddleware, loginRateLimiter)

	r := chi.NewRouter()
	r.Use(middleware.CORS())

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
package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/mulkihakim/nalar/backend/internal/class"
	"github.com/mulkihakim/nalar/backend/internal/db"
	"github.com/mulkihakim/nalar/backend/internal/exam"
	"github.com/mulkihakim/nalar/backend/internal/material"
	"github.com/mulkihakim/nalar/backend/internal/middleware"
	"github.com/mulkihakim/nalar/backend/internal/session"
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

	// Jalankan AutoMigrate dari GORM Model
	if err := db.AutoMigrate(gormDB); err != nil {
		log.Printf("warning: auto migrate returned error: %v", err)
	} else {
		log.Println("database auto-migrated successfully from models")
	}


	authMiddleware := middleware.Auth([]byte(jwtSecret))
	// Rate limit login: maks 5 percobaan per menit per IP
	loginRateLimiter := middleware.RateLimit(5, time.Minute)

	// User domain
	userRepo := user.NewRepository(gormDB)
	userSvc := user.NewService(userRepo, []byte(jwtSecret))
	userHandler := user.NewHandler(userSvc, authMiddleware, loginRateLimiter)

	// Class domain
	classRepo := class.NewRepository(gormDB)
	classSvc := class.NewService(classRepo, userRepo)
	classHandler := class.NewHandler(classSvc, authMiddleware)

	// Material domain
	materialRepo := material.NewRepository(gormDB)
	materialSvc := material.NewService(materialRepo)
	materialHandler := material.NewHandler(materialSvc, authMiddleware)

	// Exam domain
	examRepo := exam.NewRepository(gormDB)
	examSvc := exam.NewService(examRepo, materialRepo)
	examHandler := exam.NewHandler(examSvc, authMiddleware)

	// Session domain
	sessionRepo := session.NewRepository(gormDB)
	sessionSvc := session.NewService(sessionRepo)
	sessionHandler := session.NewHandler(sessionSvc, authMiddleware)

	r := chi.NewRouter()
	r.Use(middleware.CORS())

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		userHandler.RegisterRoutes(r)
		classHandler.RegisterRoutes(r)
		materialHandler.RegisterRoutes(r)
		examHandler.RegisterRoutes(r)
		sessionHandler.RegisterRoutes(r)
	})

	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := net.JoinHostPort(host, port)

	log.Printf("server running on http://%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
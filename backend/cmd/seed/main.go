package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mulkihakim/nalar/backend/internal/db"
	"github.com/mulkihakim/nalar/backend/internal/user"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	dsn := db.GetDSN()
	if dsn == "" {
		log.Fatal("database configuration (DATABASE_URL or DB_*) must be set")
	}

	gormDB, err := db.Connect(dsn)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	repo := user.NewRepository(gormDB)

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	seedUsers := []user.User{
		{
			Name:         "Super Admin",
			Username:     "admin",
			PasswordHash: string(hash),
			Role:         "admin",
			IsActive:     true,
		},
		{
			Name:         "Budi Asesor",
			Username:     "asesor1",
			PasswordHash: string(hash),
			Role:         "asesor",
			IsActive:     true,
		},
		{
			Name:         "Siti Siswa",
			Username:     "siswa1",
			PasswordHash: string(hash),
			Role:         "siswa",
			IsActive:     true,
		},
	}

	for _, u := range seedUsers {
		existing, err := repo.FindByUsername(u.Username)
		if err != nil {
			log.Printf("error checking user %s: %v", u.Username, err)
			continue
		}
		if existing != nil {
			log.Printf("user %s already exists, skipping", u.Username)
			continue
		}

		userToCreate := u
		if err := repo.Create(&userToCreate); err != nil {
			log.Printf("failed to create user %s: %v", u.Username, err)
		} else {
			log.Printf("successfully created user: %s (role: %s, password: password123)", u.Username, u.Role)
		}
	}

	log.Println("seeding completed!")
}

package db

import (
	"fmt"
	"os"
	"time"

	"github.com/mulkihakim/nalar/backend/internal/class"
	"github.com/mulkihakim/nalar/backend/internal/exam"
	"github.com/mulkihakim/nalar/backend/internal/material"
	"github.com/mulkihakim/nalar/backend/internal/session"
	"github.com/mulkihakim/nalar/backend/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetDSN mengambil connection string dari DATABASE_URL atau menggabungkan variabel DB_*
func GetDSN() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}

	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	user := getEnvOrDefault("DB_USER", "postgres")
	password := os.Getenv("DB_PASSWORD")
	dbname := getEnvOrDefault("DB_NAME", "nalar")
	sslmode := getEnvOrDefault("DB_SSLMODE", "disable")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func Connect(dsn string) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return gormDB, nil
}

// AutoMigrate mendaftarkan semua struct model GORM agar tabel dibuat/diperbarui otomatis
func AutoMigrate(gormDB *gorm.DB) error {
	return gormDB.AutoMigrate(
		&user.User{},
		&class.Class{},
		&class.ClassMember{},
		&material.Material{},
		&material.Argument{},
		&material.Option{},
		&exam.Exam{},
		&exam.ExamClass{},
		&exam.ExamStudent{},
		&session.Session{},
		&session.SessionArgument{},
		&session.AttemptLog{},
		&session.ArgumentProgress{},
	)
}

// DropAllTables membersihkan seluruh tabel di schema public (cocok untuk PostgreSQL)
func DropAllTables(gormDB *gorm.DB) error {
	return gormDB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error
}

// ResetAndMigrate melakukan drop all tables dan migrasi ulang dari model (mirip php artisan migrate:fresh)
func ResetAndMigrate(gormDB *gorm.DB) error {
	if err := DropAllTables(gormDB); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}
	if err := AutoMigrate(gormDB); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}
	return nil
}
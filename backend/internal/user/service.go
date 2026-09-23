package user

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mulkihakim/nalar/backend/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameTaken      = errors.New("username sudah dipakai")
	ErrForbiddenRole      = errors.New("tidak berwenang membuat role ini")
	ErrInvalidCredentials = errors.New("username atau password salah")
	ErrInactiveUser       = errors.New("akun tidak aktif")
	ErrUserNotFound       = errors.New("user tidak ditemukan")
)

type CreateUserInput struct {
	Name     string
	Username string
	Password string
	Role     string // "admin" | "asesor" | "siswa"
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type Service interface {
	CreateUser(input CreateUserInput, creatorRole string, creatorID uint) (*User, error)
	Authenticate(username, password string) (*User, error)
	Login(username, password string) (*LoginResponse, error)
	GetMe(userID uint) (*User, error)
	GetByID(id uint) (*User, error)
	ListUsers() ([]User, error)
}

type service struct {
	repo      Repository
	jwtSecret []byte
	jwtExpiry time.Duration
}

func NewService(repo Repository, jwtSecret []byte) Service {
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("nalar-default-secret-change-in-production")
	}
	return &service{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExpiry: 24 * time.Hour,
	}
}

func (s *service) CreateUser(input CreateUserInput, creatorRole string, creatorID uint) (*User, error) {
	if !canCreateRole(creatorRole, input.Role) {
		return nil, ErrForbiddenRole
	}

	existing, err := s.repo.FindByUsername(input.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &User{
		Name:         input.Name,
		Username:     input.Username,
		PasswordHash: string(hash),
		Role:         input.Role,
		CreatedBy:    &creatorID,
		IsActive:     true,
	}

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *service) Authenticate(username, password string) (*User, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	if !u.IsActive {
		return nil, ErrInactiveUser
	}
	return u, nil
}

func (s *service) Login(username, password string) (*LoginResponse, error) {
	u, err := s.Authenticate(username, password)
	if err != nil {
		return nil, err
	}

	token, err := s.generateJWT(u.ID, u.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  u,
	}, nil
}

func (s *service) generateJWT(userID uint, role string) (string, error) {
	claims := middleware.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *service) GetMe(userID uint) (*User, error) {
	u, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrUserNotFound
	}
	if !u.IsActive {
		return nil, ErrInactiveUser
	}
	return u, nil
}

func (s *service) GetByID(id uint) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *service) ListUsers() ([]User, error) {
	return s.repo.List()
}

// canCreateRole menegakkan tabel hak akses di 01-overview.md §3:
// Admin -> boleh buat asesor & siswa. Asesor -> boleh buat siswa saja.
func canCreateRole(creatorRole, targetRole string) bool {
	switch creatorRole {
	case "admin":
		return targetRole == "asesor" || targetRole == "siswa"
	case "asesor":
		return targetRole == "siswa"
	default:
		return false
	}
}
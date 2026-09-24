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
	ErrUnauthorizedAccess = errors.New("tidak memiliki izin untuk mengubah data pengguna ini")
	ErrUserHasExamRecords = errors.New("pengguna tidak dapat dihapus permanen karena telah memiliki riwayat pengerjaan ujian")
)

type CreateUserInput struct {
	Name     string
	Username string
	Password string
	Role     string // "admin" | "asesor" | "siswa"
}

type UpdateUserInput struct {
	Name     *string `json:"name"`
	Password *string `json:"password"`
	IsActive *bool   `json:"is_active"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type Service interface {
	CreateUser(input CreateUserInput, creatorRole string, creatorID uint) (*User, error)
	UpdateUser(id uint, input UpdateUserInput, updaterRole string, updaterID uint) (*User, error)
	DeleteUser(id uint, requesterRole string, requesterID uint) error
	Authenticate(username, password string) (*User, error)
	Login(username, password string) (*LoginResponse, error)
	GetMe(userID uint) (*User, error)
	GetByID(id uint) (*User, error)
	ListUsers(roleFilter string, requesterRole ...string) ([]User, error)
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

func (s *service) UpdateUser(id uint, input UpdateUserInput, updaterRole string, updaterID uint) (*User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrUserNotFound
	}

	// Otorisasi:
	// Admin boleh ubah siapa pun.
	// Asesor hanya boleh ubah akun siswa yang ia buat (created_by == updaterID).
	if updaterRole == "asesor" {
		if u.Role != "siswa" || u.CreatedBy == nil || *u.CreatedBy != updaterID {
			return nil, ErrUnauthorizedAccess
		}
	} else if updaterRole != "admin" {
		return nil, ErrUnauthorizedAccess
	}

	if input.Name != nil && *input.Name != "" {
		u.Name = *input.Name
	}
	if input.IsActive != nil {
		u.IsActive = *input.IsActive
	}
	if input.Password != nil && *input.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		u.PasswordHash = string(hash)
	}

	if err := s.repo.Update(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *service) DeleteUser(id uint, requesterRole string, requesterID uint) error {
	// Hanya role admin yang boleh hard delete
	if requesterRole != "admin" {
		return ErrUnauthorizedAccess
	}

	// Tidak boleh menghapus diri sendiri
	if id == requesterID {
		return errors.New("tidak dapat menghapus akun Anda sendiri")
	}

	u, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if u == nil {
		return ErrUserNotFound
	}

	// Cek apakah pengguna telah memiliki riwayat sesi ujian
	hasSessions, err := s.repo.HasExamSessions(id)
	if err != nil {
		return err
	}
	if hasSessions {
		return ErrUserHasExamRecords
	}

	return s.repo.Delete(id)
}

func (s *service) ListUsers(roleFilter string, requesterRole ...string) ([]User, error) {
	users, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	reqRole := ""
	if len(requesterRole) > 0 {
		reqRole = requesterRole[0]
	}

	var filtered []User
	for _, u := range users {
		// Jika requester adalah asesor, JANGAN tampilkan admin di atasnya!
		if reqRole == "asesor" && u.Role == "admin" {
			continue
		}

		if roleFilter != "" && u.Role != roleFilter {
			continue
		}

		filtered = append(filtered, u)
	}

	return filtered, nil
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
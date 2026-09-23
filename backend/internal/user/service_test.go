package user_test

import (
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mulkihakim/nalar/backend/internal/middleware"
	"github.com/mulkihakim/nalar/backend/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	usersByID       map[uint]*user.User
	usersByUsername map[string]*user.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		usersByID:       make(map[uint]*user.User),
		usersByUsername: make(map[string]*user.User),
	}
}

func (m *mockUserRepository) Create(u *user.User) error {
	if _, exists := m.usersByUsername[u.Username]; exists {
		return errors.New("duplicate username")
	}
	u.ID = uint(len(m.usersByID) + 1)
	m.usersByID[u.ID] = u
	m.usersByUsername[u.Username] = u
	return nil
}

func (m *mockUserRepository) FindByID(id uint) (*user.User, error) {
	u, exists := m.usersByID[id]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepository) FindByUsername(username string) (*user.User, error) {
	u, exists := m.usersByUsername[username]
	if !exists {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepository) List() ([]user.User, error) {
	var list []user.User
	for _, u := range m.usersByID {
		list = append(list, *u)
	}
	return list, nil
}

func (m *mockUserRepository) Update(u *user.User) error {
	m.usersByID[u.ID] = u
	m.usersByUsername[u.Username] = u
	return nil
}

func setupTestService(t *testing.T) (user.Service, *mockUserRepository, []byte) {
	t.Helper()
	repo := newMockUserRepository()
	jwtSecret := []byte("test-jwt-secret-key-12345")
	svc := user.NewService(repo, jwtSecret)

	// Seed user aktif & nonaktif
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	activeUser := &user.User{
		Name:         "Budi Asesor",
		Username:     "asesor_budi",
		PasswordHash: string(hashedPassword),
		Role:         "asesor",
		IsActive:     true,
	}
	if err := repo.Create(activeUser); err != nil {
		t.Fatalf("failed to seed active user: %v", err)
	}

	inactiveUser := &user.User{
		Name:         "Siti Siswa",
		Username:     "siswa_siti",
		PasswordHash: string(hashedPassword),
		Role:         "siswa",
		IsActive:     false,
	}
	if err := repo.Create(inactiveUser); err != nil {
		t.Fatalf("failed to seed inactive user: %v", err)
	}

	return svc, repo, jwtSecret
}

func TestLogin_Success(t *testing.T) {
	svc, _, secret := setupTestService(t)

	res, err := svc.Login("asesor_budi", "password123")
	if err != nil {
		t.Fatalf("expected login to succeed, got error: %v", err)
	}

	if res == nil {
		t.Fatal("expected non-nil response")
	}

	if res.User == nil || res.User.Username != "asesor_budi" {
		t.Errorf("expected user asesor_budi, got %v", res.User)
	}

	if res.Token == "" {
		t.Fatal("expected non-empty JWT token")
	}

	// Verifikasi token JWT dan claim role + user_id
	claims := &middleware.JWTClaims{}
	token, err := jwt.ParseWithClaims(res.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("token failed validation: %v", err)
	}

	if claims.UserID != res.User.ID {
		t.Errorf("expected claim user_id %d, got %d", res.User.ID, claims.UserID)
	}
	if claims.Role != "asesor" {
		t.Errorf("expected claim role 'asesor', got '%s'", claims.Role)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, _, _ := setupTestService(t)

	_, err := svc.Login("asesor_budi", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}

	if !errors.Is(err, user.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, _, _ := setupTestService(t)

	_, err := svc.Login("nonexistent_user", "password123")
	if err == nil {
		t.Fatal("expected error for nonexistent user, got nil")
	}

	if !errors.Is(err, user.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_InactiveUser(t *testing.T) {
	svc, _, _ := setupTestService(t)

	_, err := svc.Login("siswa_siti", "password123")
	if err == nil {
		t.Fatal("expected error for inactive user, got nil")
	}

	if !errors.Is(err, user.ErrInactiveUser) {
		t.Errorf("expected ErrInactiveUser, got %v", err)
	}
}

func TestGetMe_Success(t *testing.T) {
	svc, repo, _ := setupTestService(t)

	u, _ := repo.FindByUsername("asesor_budi")
	me, err := svc.GetMe(u.ID)
	if err != nil {
		t.Fatalf("expected GetMe to succeed, got %v", err)
	}

	if me.ID != u.ID || me.Username != u.Username {
		t.Errorf("expected user id %d, got %d", u.ID, me.ID)
	}
}

func TestGetMe_NotFound(t *testing.T) {
	svc, _, _ := setupTestService(t)

	_, err := svc.GetMe(99999)
	if err == nil {
		t.Fatal("expected error for nonexistent user, got nil")
	}

	if !errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestGetMe_InactiveUser(t *testing.T) {
	svc, repo, _ := setupTestService(t)

	u, _ := repo.FindByUsername("siswa_siti")
	_, err := svc.GetMe(u.ID)
	if err == nil {
		t.Fatal("expected error for inactive user, got nil")
	}

	if !errors.Is(err, user.ErrInactiveUser) {
		t.Errorf("expected ErrInactiveUser, got %v", err)
	}
}

package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mulkihakim/nalar/backend/internal/middleware"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	secret := []byte("test-secret-key")
	claims := middleware.JWTClaims{
		UserID: 42,
		Role:   "asesor",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, okID := middleware.GetUserID(r.Context())
		role, okRole := middleware.GetUserRole(r.Context())

		if !okID || userID != 42 {
			t.Errorf("expected user_id 42, got %d (ok=%v)", userID, okID)
		}
		if !okRole || role != "asesor" {
			t.Errorf("expected role 'asesor', got '%s' (ok=%v)", role, okRole)
		}
		w.WriteHeader(http.StatusOK)
	})

	mw := middleware.Auth(secret)
	handler := mw(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	mw := middleware.Auth([]byte("test-secret-key"))
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestRequireRole(t *testing.T) {
	secret := []byte("test-secret-key")
	claims := middleware.JWTClaims{
		UserID: 1,
		Role:   "siswa",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(secret)

	adminOnly := middleware.RequireRole("admin")
	handler := middleware.Auth(secret)(adminOnly(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status 403 Forbidden for siswa on admin endpoint, got %d", rec.Code)
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	limiter := middleware.RateLimit(2, time.Minute)
	handler := limiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 1st request -> ok
	req1 := httptest.NewRequest("POST", "/auth/login", nil)
	req1.RemoteAddr = "192.168.1.100:1234"
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Errorf("expected 200 on 1st request, got %d", rec1.Code)
	}

	// 2nd request -> ok
	req2 := httptest.NewRequest("POST", "/auth/login", nil)
	req2.RemoteAddr = "192.168.1.100:1234"
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Errorf("expected 200 on 2nd request, got %d", rec2.Code)
	}

	// 3rd request -> 429 Too Many Requests
	req3 := httptest.NewRequest("POST", "/auth/login", nil)
	req3.RemoteAddr = "192.168.1.100:1234"
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429 on 3rd request, got %d", rec3.Code)
	}
}

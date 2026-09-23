package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type clientRecord struct {
	count     int
	resetTime time.Time
}

type IPRateLimiter struct {
	mu      sync.Mutex
	records map[string]*clientRecord
	limit   int
	window  time.Duration
}

func NewIPRateLimiter(limit int, window time.Duration) *IPRateLimiter {
	limiter := &IPRateLimiter{
		records: make(map[string]*clientRecord),
		limit:   limit,
		window:  window,
	}

	// Pembersihan rutin data IP yang telah kedaluwarsa setiap 5 menit
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			limiter.cleanup()
		}
	}()

	return limiter
}

func (l *IPRateLimiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	for ip, record := range l.records {
		if now.After(record.resetTime) {
			delete(l.records, ip)
		}
	}
}

func (l *IPRateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	rec, exists := l.records[ip]
	if !exists || now.After(rec.resetTime) {
		l.records[ip] = &clientRecord{
			count:     1,
			resetTime: now.Add(l.window),
		}
		return true
	}

	if rec.count < l.limit {
		rec.count++
		return true
	}

	return false
}

// RateLimit membuat HTTP middleware chi untuk membatasi request per IP address
func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	limiter := NewIPRateLimiter(limit, window)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)
			if !limiter.Allow(ip) {
				writeJSONError(w, http.StatusTooManyRequests, "terlalu banyak percobaan login, coba lagi beberapa saat lagi")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

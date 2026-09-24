package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mulkihakim/nalar/backend/internal/middleware"
)

type Handler struct {
	svc            Service
	authMiddleware func(http.Handler) http.Handler
	rateLimiter    func(http.Handler) http.Handler
}

func NewHandler(svc Service, authMiddleware func(http.Handler) http.Handler, rateLimiter func(http.Handler) http.Handler) *Handler {
	return &Handler{
		svc:            svc,
		authMiddleware: authMiddleware,
		rateLimiter:    rateLimiter,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		if h.rateLimiter != nil {
			r.With(h.rateLimiter).Post("/login", h.login)
		} else {
			r.Post("/login", h.login)
		}

		if h.authMiddleware != nil {
			r.With(h.authMiddleware).Get("/me", h.me)
		} else {
			r.Get("/me", h.me)
		}
	})

	r.Route("/users", func(r chi.Router) {
		if h.authMiddleware != nil {
			r.Use(h.authMiddleware)
		}
		r.Use(middleware.RequireRole("admin", "asesor"))
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Patch("/{id}", h.update)
		r.Delete("/{id}", h.delete)
	})
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username dan password wajib diisi")
		return
	}

	res, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		case errors.Is(err, ErrInactiveUser):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "terjadi kesalahan saat proses login")
			return
		}
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "pengguna belum terautentikasi")
		return
	}

	u, err := h.svc.GetMe(userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrInactiveUser):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "terjadi kesalahan")
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{"user": u})
}

type createUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	creatorRole, okRole := middleware.GetUserRole(r.Context())
	creatorID, okID := middleware.GetUserID(r.Context())
	if !okRole || !okID {
		creatorRole = "admin"
		creatorID = 1
	}

	u, err := h.svc.CreateUser(CreateUserInput{
		Name:     req.Name,
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}, creatorRole, creatorID)

	switch {
	case errors.Is(err, ErrUsernameTaken):
		writeError(w, http.StatusConflict, err.Error())
		return
	case errors.Is(err, ErrForbiddenRole):
		writeError(w, http.StatusForbidden, err.Error())
		return
	case err != nil:
		writeError(w, http.StatusInternalServerError, "terjadi kesalahan")
		return
	}

	writeJSON(w, http.StatusCreated, u)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}

	u, err := h.svc.GetByID(uint(id))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "terjadi kesalahan")
		return
	}
	if u == nil {
		writeError(w, http.StatusNotFound, "user tidak ditemukan")
		return
	}

	writeJSON(w, http.StatusOK, u)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	requesterRole, _ := middleware.GetUserRole(r.Context())

	users, err := h.svc.ListUsers(role, requesterRole)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "terjadi kesalahan")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}

	requesterRole, okRole := middleware.GetUserRole(r.Context())
	requesterID, okID := middleware.GetUserID(r.Context())
	if !okRole || !okID {
		writeError(w, http.StatusUnauthorized, "pengguna belum terautentikasi")
		return
	}

	err = h.svc.DeleteUser(uint(id), requesterRole, requesterID)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrUnauthorizedAccess):
			writeError(w, http.StatusForbidden, err.Error())
			return
		case errors.Is(err, ErrUserHasExamRecords):
			writeError(w, http.StatusBadRequest, err.Error())
			return
		default:
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "pengguna berhasil dihapus secara permanen"})
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}

	var req UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	updaterRole, okRole := middleware.GetUserRole(r.Context())
	updaterID, okID := middleware.GetUserID(r.Context())
	if !okRole || !okID {
		writeError(w, http.StatusUnauthorized, "pengguna belum terautentikasi")
		return
	}

	u, err := h.svc.UpdateUser(uint(id), req, updaterRole, updaterID)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrUnauthorizedAccess):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "terjadi kesalahan saat memperbarui pengguna")
			return
		}
	}

	writeJSON(w, http.StatusOK, u)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
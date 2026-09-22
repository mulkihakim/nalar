package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
	})
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

	// Sementara hardcode — nanti diganti ambil dari JWT claim setelah auth dibuat.
	creatorRole := "admin"
	var creatorID uint = 1

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
	users, err := h.svc.ListUsers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "terjadi kesalahan")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
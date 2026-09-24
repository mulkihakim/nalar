package class

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
}

func NewHandler(svc Service, authMiddleware func(http.Handler) http.Handler) *Handler {
	return &Handler{
		svc:            svc,
		authMiddleware: authMiddleware,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/classes", func(r chi.Router) {
		if h.authMiddleware != nil {
			r.Use(h.authMiddleware)
		}
		r.Use(middleware.RequireRole("admin", "asesor"))

		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
		r.Patch("/{id}", h.update)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.delete)

		r.Get("/{id}/members", h.getMembers)
		r.Post("/{id}/members", h.addMember)
		r.Delete("/{id}/members/{userId}", h.removeMember)
	})
}

type classPayload struct {
	Name string `json:"name"`
}

type addMemberPayload struct {
	UserID uint `json:"user_id"`
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	classes, err := h.svc.ListClasses(role, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memuat daftar kelas")
		return
	}
	writeJSON(w, http.StatusOK, classes)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())

	var req classPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	c, err := h.svc.CreateClass(req.Name, userID)
	if err != nil {
		if errors.Is(err, ErrEmptyClassName) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal membuat kelas")
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	c, err := h.svc.GetClassByID(uint(id), role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrClassNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal memuat kelas")
			return
		}
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	var req classPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	c, err := h.svc.UpdateClass(uint(id), req.Name, role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyClassName):
			writeError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, ErrClassNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal memperbarui kelas")
			return
		}
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	err = h.svc.DeleteClass(uint(id), role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrClassNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal menghapus kelas")
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "kelas berhasil dihapus"})
}

func (h *Handler) getMembers(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	members, err := h.svc.GetMembers(uint(id), role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrClassNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal memuat anggota kelas")
			return
		}
	}
	writeJSON(w, http.StatusOK, members)
}

func (h *Handler) addMember(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	var req addMemberPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	err = h.svc.AddMember(uint(id), req.UserID, role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrClassNotFound), errors.Is(err, ErrStudentNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		case errors.Is(err, ErrInvalidMemberRole), errors.Is(err, ErrMemberAlreadyExists):
			writeError(w, http.StatusBadRequest, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal menambahkan anggota")
			return
		}
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "anggota berhasil ditambahkan"})
}

func (h *Handler) removeMember(w http.ResponseWriter, r *http.Request) {
	classID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id kelas tidak valid")
		return
	}
	studentID, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id siswa tidak valid")
		return
	}

	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	err = h.svc.RemoveMember(uint(classID), uint(studentID), role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrClassNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal menghapus anggota")
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "anggota berhasil dihapus"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

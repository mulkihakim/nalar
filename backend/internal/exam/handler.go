package exam

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
	r.Route("/exams", func(r chi.Router) {
		if h.authMiddleware != nil {
			r.Use(h.authMiddleware)
		}
		r.Use(middleware.RequireRole("admin", "asesor"))

		r.Get("/", h.listExams)
		r.Post("/", h.createExam)
		r.Get("/{id}", h.getExamByID)
		r.Patch("/{id}", h.updateExam)
		r.Put("/{id}", h.updateExam)
		r.Delete("/{id}", h.deleteExam)

		r.Put("/{id}/access", h.setAccess)
		r.Patch("/{id}/status", h.setStatus)
	})
}

func (h *Handler) listExams(w http.ResponseWriter, r *http.Request) {
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	exams, err := h.svc.ListExams(role, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memuat daftar ujian")
		return
	}
	writeJSON(w, http.StatusOK, exams)
}

func (h *Handler) createExam(w http.ResponseWriter, r *http.Request) {
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	var input CreateExamInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	e, err := h.svc.CreateExam(input, userID, role)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyTitle), errors.Is(err, ErrInvalidArgumentsPerSess):
			writeError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, ErrMaterialNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal membuat ujian")
			return
		}
	}
	writeJSON(w, http.StatusCreated, e)
}

func (h *Handler) getExamByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id ujian tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	e, err := h.svc.GetExamByID(uint(id), role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrExamNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal memuat ujian")
			return
		}
	}
	writeJSON(w, http.StatusOK, e)
}

func (h *Handler) updateExam(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id ujian tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	var input UpdateExamInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	e, err := h.svc.UpdateExam(uint(id), input, role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrExamNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		case errors.Is(err, ErrMaterialNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal memperbarui ujian")
			return
		}
	}
	writeJSON(w, http.StatusOK, e)
}

func (h *Handler) deleteExam(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id ujian tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	err = h.svc.DeleteExam(uint(id), role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrExamNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal menghapus ujian")
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "ujian berhasil dihapus"})
}

func (h *Handler) setAccess(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id ujian tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	var input AccessInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	err = h.svc.SetAccess(uint(id), input, role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrExamNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal mengatur akses ujian")
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "akses ujian berhasil diperbarui"})
}

type statusPayload struct {
	IsActive bool `json:"is_active"`
}

func (h *Handler) setStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id ujian tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	var input statusPayload
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	err = h.svc.SetStatus(uint(id), input.IsActive, role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrExamNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal mengubah status ujian")
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "status ujian berhasil diperbarui"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

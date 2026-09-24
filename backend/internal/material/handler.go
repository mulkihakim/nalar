package material

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
	r.Route("/materials", func(r chi.Router) {
		if h.authMiddleware != nil {
			r.Use(h.authMiddleware)
		}
		r.Use(middleware.RequireRole("admin", "asesor"))

		r.Get("/", h.listMaterials)
		r.Post("/", h.createMaterial)
		r.Get("/{id}", h.getMaterialByID)
		r.Patch("/{id}", h.updateMaterial)
		r.Put("/{id}", h.updateMaterial)
		r.Delete("/{id}", h.deleteMaterial)

		r.Get("/{id}/arguments", h.listArguments)
		r.Post("/{id}/arguments", h.createArgument)
		r.Get("/{id}/arguments/{aid}", h.getArgumentByID)
		r.Put("/{id}/arguments/{aid}", h.updateArgument)
		r.Patch("/{id}/arguments/{aid}", h.updateArgument)
		r.Delete("/{id}/arguments/{aid}", h.deleteArgument)
	})
}

func (h *Handler) listMaterials(w http.ResponseWriter, r *http.Request) {
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	materials, err := h.svc.ListMaterials(role, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memuat daftar materi")
		return
	}
	writeJSON(w, http.StatusOK, materials)
}

func (h *Handler) createMaterial(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r.Context())

	var input CreateMaterialInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	m, err := h.svc.CreateMaterial(input, userID)
	if err != nil {
		if errors.Is(err, ErrEmptyMaterialData) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal membuat materi")
		return
	}
	writeJSON(w, http.StatusCreated, m)
}

func (h *Handler) getMaterialByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id materi tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	m, err := h.svc.GetMaterialByID(uint(id), role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrMaterialNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal memuat materi")
			return
		}
	}
	writeJSON(w, http.StatusOK, m)
}

func (h *Handler) updateMaterial(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id materi tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	var input CreateMaterialInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	m, err := h.svc.UpdateMaterial(uint(id), input, role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyMaterialData):
			writeError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, ErrMaterialNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal memperbarui materi")
			return
		}
	}
	writeJSON(w, http.StatusOK, m)
}

func (h *Handler) deleteMaterial(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id materi tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	err = h.svc.DeleteMaterial(uint(id), role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrMaterialNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal menghapus materi")
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "materi berhasil dihapus"})
}

func (h *Handler) listArguments(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id materi tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	args, err := h.svc.ListArguments(uint(id), role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrMaterialNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal memuat argumen")
			return
		}
	}
	writeJSON(w, http.StatusOK, args)
}

func (h *Handler) createArgument(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id materi tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	var input ArgumentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	arg, err := h.svc.CreateArgument(uint(id), input, role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrMaterialNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		case errors.Is(err, ErrInvalidOptionsCount), errors.Is(err, ErrEmptyClaimText), errors.Is(err, ErrEmptyOptionText):
			writeError(w, http.StatusBadRequest, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal membuat argumen")
			return
		}
	}
	writeJSON(w, http.StatusCreated, arg)
}

func (h *Handler) getArgumentByID(w http.ResponseWriter, r *http.Request) {
	mid, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id materi tidak valid")
		return
	}
	aid, err := strconv.ParseUint(chi.URLParam(r, "aid"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id argumen tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	args, err := h.svc.ListArguments(uint(mid), role, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memuat argumen")
		return
	}
	for _, a := range args {
		if a.ID == uint(aid) {
			writeJSON(w, http.StatusOK, a)
			return
		}
	}
	writeError(w, http.StatusNotFound, "argumen tidak ditemukan")
}

func (h *Handler) updateArgument(w http.ResponseWriter, r *http.Request) {
	mid, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id materi tidak valid")
		return
	}
	aid, err := strconv.ParseUint(chi.URLParam(r, "aid"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id argumen tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	var input ArgumentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	arg, err := h.svc.UpdateArgument(uint(mid), uint(aid), input, role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrMaterialNotFound), errors.Is(err, ErrArgumentNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		case errors.Is(err, ErrInvalidOptionsCount), errors.Is(err, ErrEmptyClaimText), errors.Is(err, ErrEmptyOptionText):
			writeError(w, http.StatusBadRequest, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal memperbarui argumen")
			return
		}
	}
	writeJSON(w, http.StatusOK, arg)
}

func (h *Handler) deleteArgument(w http.ResponseWriter, r *http.Request) {
	mid, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id materi tidak valid")
		return
	}
	aid, err := strconv.ParseUint(chi.URLParam(r, "aid"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id argumen tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	err = h.svc.DeleteArgument(uint(mid), uint(aid), role, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrMaterialNotFound), errors.Is(err, ErrArgumentNotFound):
			writeError(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
			return
		default:
			writeError(w, http.StatusInternalServerError, "gagal menghapus argumen")
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "argumen berhasil dihapus"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

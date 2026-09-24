package session

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
	// 1. Student routes
	r.Group(func(r chi.Router) {
		if h.authMiddleware != nil {
			r.Use(h.authMiddleware)
		}
		r.Use(middleware.RequireRole("siswa"))

		r.Get("/my/exams", h.getMyExams)
		r.Get("/exams/{id}/session-status", h.getSessionStatus)
		r.Post("/exams/{id}/start", h.startSession)
		r.Get("/sessions/{id}", h.getSessionDetails)
		r.Post("/sessions/{id}/arguments/{aid}/drops", h.recordDrop)
		r.Post("/sessions/{id}/arguments/{aid}/confirm", h.confirmArgument)
		r.Get("/sessions/{id}/monitoring/{aid}", h.getMonitoring)
		r.Get("/sessions/{id}/analysis", h.getAnalysis)
		r.Get("/my/exams/{id}/sessions", h.getMyExamSessions)
	})

	// 2. Staff routes (admin / asesor)
	r.Group(func(r chi.Router) {
		if h.authMiddleware != nil {
			r.Use(h.authMiddleware)
		}
		r.Use(middleware.RequireRole("admin", "asesor"))

		r.Get("/exams/{id}/results", h.getExamResults)
		r.Get("/exams/{id}/logs", h.getExamLogs)
		r.Get("/exams/{id}/analysis", h.getStaffAnalysis)
	})
}

func (h *Handler) getMyExams(w http.ResponseWriter, r *http.Request) {
	studentID, _ := middleware.GetUserID(r.Context())
	exams, err := h.svc.GetStudentExams(studentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memuat daftar ujian")
		return
	}
	writeJSON(w, http.StatusOK, exams)
}

func (h *Handler) getSessionStatus(w http.ResponseWriter, r *http.Request) {
	examID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id ujian tidak valid")
		return
	}
	studentID, _ := middleware.GetUserID(r.Context())

	status, err := h.svc.GetSessionStatus(uint(examID), studentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memeriksa status sesi")
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func (h *Handler) startSession(w http.ResponseWriter, r *http.Request) {
	examID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id ujian tidak valid")
		return
	}
	studentID, _ := middleware.GetUserID(r.Context())

	var input StartSessionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		// Boleh body kosong, default standard
		input.Mode = ModeStandard
	}

	detail, err := h.svc.StartOrContinueSession(uint(examID), studentID, input)
	if err != nil {
		var conflictErr *ModeConflictError
		if errors.As(err, &conflictErr) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error":        conflictErr.Message,
				"current_mode": conflictErr.CurrentMode,
			})
			return
		}

		switch {
		case errors.Is(err, ErrExamNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrExamNotActive):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

func (h *Handler) getSessionDetails(w http.ResponseWriter, r *http.Request) {
	sessionID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id sesi tidak valid")
		return
	}
	studentID, _ := middleware.GetUserID(r.Context())

	detail, err := h.svc.GetSessionDetails(uint(sessionID), studentID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			writeError(w, http.StatusForbidden, err.Error())
			return
		}
		writeError(w, http.StatusNotFound, "sesi tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

func (h *Handler) recordDrop(w http.ResponseWriter, r *http.Request) {
	sessionID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id sesi tidak valid")
		return
	}
	argumentID, err := strconv.ParseUint(chi.URLParam(r, "aid"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id argumen tidak valid")
		return
	}
	studentID, _ := middleware.GetUserID(r.Context())

	var input DropInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	if err := h.svc.RecordDrop(uint(sessionID), uint(argumentID), input, studentID); err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, ErrArgumentNotCurrent), errors.Is(err, ErrInvalidSlot), errors.Is(err, ErrArgumentCompleted), errors.Is(err, ErrSessionCompleted), errors.Is(err, ErrOptionNotFound):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "drop tercatat"})
}

func (h *Handler) confirmArgument(w http.ResponseWriter, r *http.Request) {
	sessionID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id sesi tidak valid")
		return
	}
	argumentID, err := strconv.ParseUint(chi.URLParam(r, "aid"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id argumen tidak valid")
		return
	}
	studentID, _ := middleware.GetUserID(r.Context())

	var input ConfirmInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "payload tidak valid")
		return
	}

	res, err := h.svc.ConfirmArgument(uint(sessionID), uint(argumentID), input, studentID)
	if err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, ErrArgumentNotCurrent), errors.Is(err, ErrSessionCompleted):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) getMonitoring(w http.ResponseWriter, r *http.Request) {
	sessionID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id sesi tidak valid")
		return
	}
	argumentID, err := strconv.ParseUint(chi.URLParam(r, "aid"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id argumen tidak valid")
		return
	}
	studentID, _ := middleware.GetUserID(r.Context())

	res, err := h.svc.GetMonitoringAnalytics(uint(sessionID), uint(argumentID), studentID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			writeError(w, http.StatusForbidden, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) getAnalysis(w http.ResponseWriter, r *http.Request) {
	sessionID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id sesi tidak valid")
		return
	}
	studentID, _ := middleware.GetUserID(r.Context())

	res, err := h.svc.GetAnalysisAnalytics(uint(sessionID), studentID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			writeError(w, http.StatusForbidden, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) getMyExamSessions(w http.ResponseWriter, r *http.Request) {
	examID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id ujian tidak valid")
		return
	}
	studentID, _ := middleware.GetUserID(r.Context())

	history, err := h.svc.GetStudentHistory(uint(examID), studentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memuat riwayat ujian")
		return
	}
	writeJSON(w, http.StatusOK, history)
}

func (h *Handler) getExamResults(w http.ResponseWriter, r *http.Request) {
	examID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id ujian tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	res, err := h.svc.GetExamResults(uint(examID), role, userID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			writeError(w, http.StatusForbidden, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *Handler) getExamLogs(w http.ResponseWriter, r *http.Request) {
	examID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id ujian tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	logs, err := h.svc.GetExamLogs(uint(examID), role, userID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			writeError(w, http.StatusForbidden, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, logs)
}

func (h *Handler) getStaffAnalysis(w http.ResponseWriter, r *http.Request) {
	examID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id ujian tidak valid")
		return
	}
	role, _ := middleware.GetUserRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	res, err := h.svc.GetStaffAnalytics(uint(examID), role, userID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			writeError(w, http.StatusForbidden, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

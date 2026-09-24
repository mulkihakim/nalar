package session

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/mulkihakim/nalar/backend/internal/exam"
	"github.com/mulkihakim/nalar/backend/internal/material"
)

const MinPeersDefault = 3

var (
	ErrForbidden            = errors.New("akses ditolak")
	ErrExamNotActive        = errors.New("ujian tidak aktif")
	ErrExamNotFound         = errors.New("ujian tidak ditemukan")
	ErrModeConflict         = errors.New("mode belajar berbeda dengan sesi berjalan")
	ErrInvalidSlot          = errors.New("slot tidak valid atau tidak cocok dengan tipe opsi")
	ErrArgumentNotCurrent   = errors.New("argumen harus diselesaikan berurutan")
	ErrArgumentCompleted    = errors.New("argumen sudah selesai")
	ErrSessionCompleted     = errors.New("sesi sudah selesai")
	ErrOptionNotFound       = errors.New("opsi tidak ditemukan pada argumen ini")
)

type ModeConflictError struct {
	CurrentMode SessionMode `json:"current_mode"`
	Message     string      `json:"message"`
}

func (e *ModeConflictError) Error() string {
	return e.Message
}

type StartSessionInput struct {
	Mode               SessionMode `json:"mode"`
	ConfirmModeChange  bool        `json:"confirm_mode_change"`
}

type DropInput struct {
	OptionID uint   `json:"option_id"`
	Slot     string `json:"slot"` // "ground" | "warrant"
}

type ConfirmInput struct {
	GroundOptionID  uint `json:"ground_option_id"`
	WarrantOptionID uint `json:"warrant_option_id"`
}

type ConfirmResult struct {
	IsCorrect          bool   `json:"is_correct"`
	Completed          bool   `json:"completed"`
	IsSessionCompleted bool   `json:"is_session_completed"`
	Message            string `json:"message"`
	GroundCorrect      *bool  `json:"ground_correct,omitempty"`
	WarrantCorrect     *bool  `json:"warrant_correct,omitempty"`
}

type OptionResponse struct {
	ID         uint   `json:"id"`
	ArgumentID uint   `json:"argument_id"`
	Type       string `json:"type"`
	Text       string `json:"text"`
	IsCorrect  *bool  `json:"is_correct,omitempty"`
}

type ArgumentResponse struct {
	ID        uint             `json:"id"`
	ClaimText string           `json:"claim_text"`
	OrderNo   int              `json:"order_no"`
	Completed bool             `json:"completed"`
	Options   []OptionResponse `json:"options"`
}

type SessionDetailResponse struct {
	ID                 uint               `json:"id"`
	ExamID             uint               `json:"exam_id"`
	ExamTitle          string             `json:"exam_title"`
	MaterialTitle      string             `json:"material_title"`
	MaterialContent    string             `json:"material_content"`
	StudentID          uint               `json:"student_id"`
	AttemptNo          int                `json:"attempt_no"`
	Mode               SessionMode        `json:"mode"`
	Status             SessionStatus      `json:"status"`
	CurrentArgumentID  *uint              `json:"current_argument_id,omitempty"`
	Arguments          []ArgumentResponse `json:"arguments"`
	CompletedArguments int                `json:"completed_arguments"`
	TotalArguments     int                `json:"total_arguments"`
}

type SessionStatusResponse struct {
	HasActiveSession bool        `json:"has_active_session"`
	SessionID        *uint       `json:"session_id,omitempty"`
	Mode             SessionMode `json:"mode,omitempty"`
}

type MonitoringOptionStat struct {
	OptionID uint   `json:"option_id"`
	Type     string `json:"type"`
	Text     string `json:"text"`
	X        int64  `json:"x"`
	Y        string `json:"y"` // number as string or "-"
}

type MonitoringAnalyticsResponse struct {
	ArgumentID       uint                   `json:"argument_id"`
	ClaimText        string                 `json:"claim_text"`
	HasEnoughPeers   bool                   `json:"has_enough_peers"`
	MinPeers         int                    `json:"min_peers"`
	TotalPeers       int64                  `json:"total_peers"`
	Options          []MonitoringOptionStat `json:"options"`
}

type StudentChooserInfo struct {
	StudentID   uint      `json:"student_id"`
	StudentName string    `json:"student_name"`
	Username    string    `json:"username"`
	Attempts    int64     `json:"attempts"`
	LastAttempt time.Time `json:"last_attempt"`
}

type AnalysisOptionStat struct {
	OptionID       uint                 `json:"option_id"`
	Type           string               `json:"type"`
	Text           string               `json:"text"`
	TotalAttempts  int64                `json:"total_attempts"` // A
	UniqueStudents int64                `json:"unique_students"` // B
	Ratio          string               `json:"ratio"`           // "A:B"
	Students       []StudentChooserInfo `json:"students,omitempty"`
}

type AnalysisArgumentStat struct {
	ArgumentID uint                 `json:"argument_id"`
	ClaimText  string               `json:"claim_text"`
	OrderNo    int                  `json:"order_no"`
	Options    []AnalysisOptionStat `json:"options"`
}

type AnalysisAnalyticsResponse struct {
	ExamID         uint                   `json:"exam_id"`
	HasEnoughPeers bool                   `json:"has_enough_peers"`
	MinPeers       int                    `json:"min_peers"`
	TotalPeers     int64                  `json:"total_peers"`
	Arguments      []AnalysisArgumentStat `json:"arguments"`
}

type ArgumentHistoryItem struct {
	ArgumentID uint   `json:"argument_id"`
	ClaimText  string `json:"claim_text"`
	Attempts   int    `json:"attempts"`
	IsClean    bool   `json:"is_clean"` // exactly 2 attempts
}

type SessionHistoryItem struct {
	ID                 uint                  `json:"id"`
	AttemptNo          int                   `json:"attempt_no"`
	Mode               SessionMode           `json:"mode"`
	Status             SessionStatus         `json:"status"`
	StartedAt          time.Time             `json:"started_at"`
	CompletedAt        *time.Time            `json:"completed_at,omitempty"`
	TotalAttempts      int                   `json:"total_attempts"`
	CompletedArguments int                   `json:"completed_arguments"`
	TotalArguments     int                   `json:"total_arguments"`
	ArgumentDetails    []ArgumentHistoryItem `json:"argument_details"`
}

type Service interface {
	GetStudentExams(studentID uint) ([]exam.Exam, error)
	GetSessionStatus(examID, studentID uint) (*SessionStatusResponse, error)
	StartOrContinueSession(examID, studentID uint, input StartSessionInput) (*SessionDetailResponse, error)
	GetSessionDetails(sessionID, studentID uint) (*SessionDetailResponse, error)
	RecordDrop(sessionID, argumentID uint, input DropInput, studentID uint) error
	ConfirmArgument(sessionID, argumentID uint, input ConfirmInput, studentID uint) (*ConfirmResult, error)
	GetMonitoringAnalytics(sessionID, argumentID, studentID uint) (*MonitoringAnalyticsResponse, error)
	GetAnalysisAnalytics(sessionID, studentID uint) (*AnalysisAnalyticsResponse, error)
	GetStudentHistory(examID, studentID uint) ([]SessionHistoryItem, error)
	
	// Staff
	GetExamResults(examID uint, requesterRole string, requesterID uint) (any, error)
	GetExamLogs(examID uint, requesterRole string, requesterID uint) (any, error)
	GetStaffAnalytics(examID uint, requesterRole string, requesterID uint) (*AnalysisAnalyticsResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetStudentExams(studentID uint) ([]exam.Exam, error) {
	return s.repo.GetStudentExams(studentID)
}

func (s *service) GetSessionStatus(examID, studentID uint) (*SessionStatusResponse, error) {
	active, err := s.repo.GetActiveSession(examID, studentID)
	if err != nil {
		return nil, err
	}
	if active == nil {
		return &SessionStatusResponse{HasActiveSession: false}, nil
	}
	return &SessionStatusResponse{
		HasActiveSession: true,
		SessionID:        &active.ID,
		Mode:             active.Mode,
	}, nil
}

func (s *service) StartOrContinueSession(examID, studentID uint, input StartSessionInput) (*SessionDetailResponse, error) {
	if input.Mode == "" {
		input.Mode = ModeStandard
	}
	if input.Mode != ModeStandard && input.Mode != ModeHelp && input.Mode != ModeSocial {
		input.Mode = ModeStandard
	}

	// 1. Validasi akses exam
	e, err := s.repo.GetExamByID(examID)
	if err != nil {
		return nil, ErrExamNotFound
	}
	if !e.IsActive {
		return nil, ErrExamNotActive
	}

	// Cek apakah student berhak mengakses exam
	studentExams, err := s.repo.GetStudentExams(studentID)
	if err != nil {
		return nil, err
	}
	hasAccess := false
	for _, se := range studentExams {
		if se.ID == examID {
			hasAccess = true
			break
		}
	}
	if !hasAccess {
		return nil, ErrForbidden
	}

	// 2. Cek active session
	active, err := s.repo.GetActiveSession(examID, studentID)
	if err != nil {
		return nil, err
	}

	if active != nil {
		// Sudah ada sesi in_progress
		if active.Mode != input.Mode {
			if !input.ConfirmModeChange {
				return nil, &ModeConflictError{
					CurrentMode: active.Mode,
					Message:     fmt.Sprintf("Sesi sedang berjalan menggunakan mode %s. Apakah Anda yakin ingin mengubah ke mode %s?", active.Mode, input.Mode),
				}
			}
			// Update mode
			if err := s.repo.UpdateSessionMode(active.ID, input.Mode); err != nil {
				return nil, err
			}
			active.Mode = input.Mode
		}
		return s.buildSessionDetailResponse(active, false)
	}

	// 3. Buat sesi baru
	allArgs, err := s.repo.GetArgumentsByMaterialID(e.MaterialID)
	if err != nil {
		return nil, err
	}
	if len(allArgs) == 0 {
		return nil, errors.New("materi belum memiliki argumen")
	}

	// Acak argumen
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]material.Argument, len(allArgs))
	copy(shuffled, allArgs)
	randGen.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	count := e.ArgumentsPerSession
	if count <= 0 {
		count = 3
	}
	if count > len(shuffled) {
		count = len(shuffled)
	}

	var selectedArgIDs []uint
	for i := 0; i < count; i++ {
		selectedArgIDs = append(selectedArgIDs, shuffled[i].ID)
	}

	maxAttempt, err := s.repo.GetMaxAttemptNo(examID, studentID)
	if err != nil {
		return nil, err
	}

	newSession := &Session{
		ExamID:    examID,
		StudentID: studentID,
		AttemptNo: maxAttempt + 1,
		Mode:      input.Mode,
		Status:    StatusInProgress,
		StartedAt: time.Now(),
	}

	if err := s.repo.CreateSession(newSession, selectedArgIDs); err != nil {
		return nil, err
	}

	// Load full session
	created, err := s.repo.GetSessionByID(newSession.ID)
	if err != nil {
		return nil, err
	}

	return s.buildSessionDetailResponse(created, false)
}

func (s *service) GetSessionDetails(sessionID, studentID uint) (*SessionDetailResponse, error) {
	sess, err := s.repo.GetSessionByID(sessionID)
	if err != nil {
		return nil, err
	}
	if sess.StudentID != studentID {
		return nil, ErrForbidden
	}
	return s.buildSessionDetailResponse(sess, false)
}

func (s *service) buildSessionDetailResponse(sess *Session, allowIsCorrect bool) (*SessionDetailResponse, error) {
	// Map progress
	completedMap := make(map[uint]bool)
	for _, p := range sess.Progress {
		completedMap[p.ArgumentID] = true
	}

	var argResponses []ArgumentResponse
	var currentArgID *uint

	for _, sa := range sess.Arguments {
		if sa.Argument == nil {
			continue
		}
		isComp := completedMap[sa.ArgumentID]
		if !isComp && currentArgID == nil {
			curID := sa.ArgumentID
			currentArgID = &curID
		}

		var groundOpts []OptionResponse
		var warrantOpts []OptionResponse
		for _, o := range sa.Argument.Options {
			optResp := OptionResponse{
				ID:         o.ID,
				ArgumentID: o.ArgumentID,
				Type:       o.Type,
				Text:       o.Text,
			}
			// Sembunyikan is_correct dari klien siswa, KECUALI jika argumen sudah diselesaikan pada sesi ini
			if allowIsCorrect || isComp {
				correct := o.IsCorrect
				optResp.IsCorrect = &correct
			}
			if o.Type == "ground" {
				groundOpts = append(groundOpts, optResp)
			} else {
				warrantOpts = append(warrantOpts, optResp)
			}
		}

		// Acak urutan opsi ground dan warrant secara deterministik (stabil per sesi & argumen)
		// Sesuai F-15: "Urutan opsi diacak, stabil per sesi"
		optSeed := int64(sess.ID*10007 + sa.ArgumentID*37 + 101)
		rGround := rand.New(rand.NewSource(optSeed))
		rGround.Shuffle(len(groundOpts), func(i, j int) {
			groundOpts[i], groundOpts[j] = groundOpts[j], groundOpts[i]
		})

		rWarrant := rand.New(rand.NewSource(optSeed + 53))
		rWarrant.Shuffle(len(warrantOpts), func(i, j int) {
			warrantOpts[i], warrantOpts[j] = warrantOpts[j], warrantOpts[i]
		})

		var opts []OptionResponse
		opts = append(opts, groundOpts...)
		opts = append(opts, warrantOpts...)

		argResponses = append(argResponses, ArgumentResponse{
			ID:        sa.ArgumentID,
			ClaimText: sa.Argument.ClaimText,
			OrderNo:   sa.OrderNo,
			Completed: isComp,
			Options:   opts,
		})
	}

	resp := &SessionDetailResponse{
		ID:                 sess.ID,
		ExamID:             sess.ExamID,
		StudentID:          sess.StudentID,
		AttemptNo:          sess.AttemptNo,
		Mode:               sess.Mode,
		Status:             sess.Status,
		CurrentArgumentID:  currentArgID,
		Arguments:          argResponses,
		CompletedArguments: len(sess.Progress),
		TotalArguments:     len(sess.Arguments),
	}

	if sess.Exam != nil {
		resp.ExamTitle = sess.Exam.Title
		if sess.Exam.Material != nil {
			resp.MaterialTitle = sess.Exam.Material.Title
			resp.MaterialContent = sess.Exam.Material.Content
		}
	}

	return resp, nil
}

func (s *service) RecordDrop(sessionID, argumentID uint, input DropInput, studentID uint) error {
	sess, err := s.repo.GetSessionByID(sessionID)
	if err != nil {
		return err
	}
	if sess.StudentID != studentID {
		return ErrForbidden
	}
	if sess.Status != StatusInProgress {
		return ErrSessionCompleted
	}

	// Validasi argument ada dalam sesi
	inSession := false
	for _, sa := range sess.Arguments {
		if sa.ArgumentID == argumentID {
			inSession = true
			break
		}
	}
	if !inSession {
		return errors.New("argumen bukan bagian dari sesi ini")
	}

	// Validasi argumen belum selesai
	completed, err := s.repo.IsArgumentCompleted(sessionID, argumentID)
	if err != nil {
		return err
	}
	if completed {
		return ErrArgumentCompleted
	}

	// Validasi urutan pengerjaan (kriteria #5)
	currentArgID, err := s.getCurrentActiveArgumentID(sess)
	if err != nil {
		return err
	}
	if currentArgID != nil && *currentArgID != argumentID {
		return ErrArgumentNotCurrent
	}

	// Validasi slot dan option
	arg, err := s.repo.GetArgumentWithOptions(argumentID)
	if err != nil {
		return err
	}
	var targetOption *material.Option
	for i := range arg.Options {
		if arg.Options[i].ID == input.OptionID {
			targetOption = &arg.Options[i]
			break
		}
	}
	if targetOption == nil {
		return ErrOptionNotFound
	}
	if input.Slot != "ground" && input.Slot != "warrant" {
		return ErrInvalidSlot
	}
	if targetOption.Type != input.Slot {
		return ErrInvalidSlot
	}

	// Simpan drop attempt
	log := &AttemptLog{
		SessionID:  sessionID,
		ArgumentID: argumentID,
		OptionID:   input.OptionID,
		Slot:       input.Slot,
		CreatedAt:  time.Now(),
	}
	return s.repo.LogAttempt(log)
}

func (s *service) ConfirmArgument(sessionID, argumentID uint, input ConfirmInput, studentID uint) (*ConfirmResult, error) {
	sess, err := s.repo.GetSessionByID(sessionID)
	if err != nil {
		return nil, err
	}
	if sess.StudentID != studentID {
		return nil, ErrForbidden
	}
	if sess.Status != StatusInProgress {
		return nil, ErrSessionCompleted
	}

	// Validasi urutan pengerjaan
	currentArgID, err := s.getCurrentActiveArgumentID(sess)
	if err != nil {
		return nil, err
	}
	if currentArgID != nil && *currentArgID != argumentID {
		return nil, ErrArgumentNotCurrent
	}

	// Validasi opsi
	arg, err := s.repo.GetArgumentWithOptions(argumentID)
	if err != nil {
		return nil, err
	}

	var groundOpt, warrantOpt *material.Option
	for i := range arg.Options {
		if arg.Options[i].ID == input.GroundOptionID {
			groundOpt = &arg.Options[i]
		}
		if arg.Options[i].ID == input.WarrantOptionID {
			warrantOpt = &arg.Options[i]
		}
	}

	if groundOpt == nil || warrantOpt == nil {
		return nil, errors.New("opsi ground atau warrant tidak ditemukan pada argumen ini")
	}
	if groundOpt.Type != "ground" || warrantOpt.Type != "warrant" {
		return nil, errors.New("tipe opsi tidak sesuai slot")
	}

	isGroundCorrect := groundOpt.IsCorrect
	isWarrantCorrect := warrantOpt.IsCorrect
	allCorrect := isGroundCorrect && isWarrantCorrect

	res := &ConfirmResult{
		IsCorrect: allCorrect,
		Completed: allCorrect,
	}

	if allCorrect {
		res.Message = "Tepat sekali! Susunan argumen Anda sudah benar."
		if err := s.repo.MarkArgumentComplete(sessionID, argumentID); err != nil {
			return nil, err
		}

		// Cek apakah ini argumen terakhir di sesi
		progress, err := s.repo.GetArgumentProgress(sessionID)
		if err != nil {
			return nil, err
		}
		if len(progress) >= len(sess.Arguments) {
			res.IsSessionCompleted = true
			if err := s.repo.CompleteSession(sessionID); err != nil {
				return nil, err
			}
		}
	} else {
		res.Message = "Susunan argumen belum tepat. Silakan coba lagi."
	}

	// Feedback khusus mode bantuan
	if sess.Mode == ModeHelp {
		res.GroundCorrect = &isGroundCorrect
		res.WarrantCorrect = &isWarrantCorrect
	}

	return res, nil
}

func (s *service) getCurrentActiveArgumentID(sess *Session) (*uint, error) {
	completedMap := make(map[uint]bool)
	for _, p := range sess.Progress {
		completedMap[p.ArgumentID] = true
	}
	for _, sa := range sess.Arguments {
		if !completedMap[sa.ArgumentID] {
			id := sa.ArgumentID
			return &id, nil
		}
	}
	return nil, nil
}

func (s *service) GetMonitoringAnalytics(sessionID, argumentID, studentID uint) (*MonitoringAnalyticsResponse, error) {
	sess, err := s.repo.GetSessionByID(sessionID)
	if err != nil {
		return nil, err
	}
	if sess.StudentID != studentID {
		return nil, ErrForbidden
	}

	uniquePeers, err := s.repo.CountUniquePeersWithSessions(sess.ExamID)
	if err != nil {
		return nil, err
	}

	arg, err := s.repo.GetArgumentWithOptions(argumentID)
	if err != nil {
		return nil, err
	}

	resp := &MonitoringAnalyticsResponse{
		ArgumentID:     argumentID,
		ClaimText:      arg.ClaimText,
		MinPeers:       MinPeersDefault,
		TotalPeers:     uniquePeers,
		HasEnoughPeers: uniquePeers >= MinPeersDefault,
	}

	for _, opt := range arg.Options {
		x, err := s.repo.CountStudentAttemptsByOption(sess.ExamID, studentID, opt.ID)
		if err != nil {
			return nil, err
		}

		yStr := "-"
		if resp.HasEnoughPeers {
			totalAtt, uniqUsers, err := s.repo.GetGroupAttemptsByOption(sess.ExamID, opt.ID)
			if err != nil {
				return nil, err
			}
			if uniqUsers > 0 {
				// Y = round(total / uniqUsers) half away from zero
				yVal := math.Round(float64(totalAtt) / float64(uniqUsers))
				yStr = fmt.Sprintf("%.0f", yVal)
			}
		}

		resp.Options = append(resp.Options, MonitoringOptionStat{
			OptionID: opt.ID,
			Type:     opt.Type,
			Text:     opt.Text,
			X:        x,
			Y:        yStr,
		})
	}

	return resp, nil
}

func (s *service) GetAnalysisAnalytics(sessionID, studentID uint) (*AnalysisAnalyticsResponse, error) {
	sess, err := s.repo.GetSessionByID(sessionID)
	if err != nil {
		return nil, err
	}
	if sess.StudentID != studentID {
		return nil, ErrForbidden
	}

	uniquePeers, err := s.repo.CountUniquePeersWithSessions(sess.ExamID)
	if err != nil {
		return nil, err
	}

	resp := &AnalysisAnalyticsResponse{
		ExamID:         sess.ExamID,
		MinPeers:       MinPeersDefault,
		TotalPeers:     uniquePeers,
		HasEnoughPeers: uniquePeers >= MinPeersDefault,
	}

	for _, sa := range sess.Arguments {
		if sa.Argument == nil {
			continue
		}
		var optStats []AnalysisOptionStat
		for _, opt := range sa.Argument.Options {
			var totalAtt, uniqUsers int64
			var ratio string

			if resp.HasEnoughPeers {
				totalAtt, uniqUsers, err = s.repo.GetGroupAttemptsByOption(sess.ExamID, opt.ID)
				if err != nil {
					return nil, err
				}
				ratio = fmt.Sprintf("%d:%d", totalAtt, uniqUsers)
			} else {
				ratio = "-:-"
			}

			optStats = append(optStats, AnalysisOptionStat{
				OptionID:       opt.ID,
				Type:           opt.Type,
				Text:           opt.Text,
				TotalAttempts:  totalAtt,
				UniqueStudents: uniqUsers,
				Ratio:          ratio,
			})
		}

		resp.Arguments = append(resp.Arguments, AnalysisArgumentStat{
			ArgumentID: sa.ArgumentID,
			ClaimText:  sa.Argument.ClaimText,
			OrderNo:    sa.OrderNo,
			Options:    optStats,
		})
	}

	return resp, nil
}

func (s *service) GetStudentHistory(examID, studentID uint) ([]SessionHistoryItem, error) {
	sessions, err := s.repo.GetSessionsByExamAndStudent(examID, studentID)
	if err != nil {
		return nil, err
	}

	var result []SessionHistoryItem
	for _, sess := range sessions {
		logs, err := s.repo.GetAttemptLogsBySession(sess.ID)
		if err != nil {
			return nil, err
		}

		// Hitung percobaan per argumen
		attemptsPerArg := make(map[uint]int)
		for _, l := range logs {
			attemptsPerArg[l.ArgumentID]++
		}

		var argDetails []ArgumentHistoryItem
		for _, sa := range sess.Arguments {
			claim := ""
			if sa.Argument != nil {
				claim = sa.Argument.ClaimText
			}
			attCount := attemptsPerArg[sa.ArgumentID]
			argDetails = append(argDetails, ArgumentHistoryItem{
				ArgumentID: sa.ArgumentID,
				ClaimText:  claim,
				Attempts:   attCount,
				IsClean:    attCount == 2, // 1 ground + 1 warrant
			})
		}

		result = append(result, SessionHistoryItem{
			ID:                 sess.ID,
			AttemptNo:          sess.AttemptNo,
			Mode:               sess.Mode,
			Status:             sess.Status,
			StartedAt:          sess.StartedAt,
			CompletedAt:        sess.CompletedAt,
			TotalAttempts:      len(logs),
			CompletedArguments: len(sess.Progress),
			TotalArguments:     len(sess.Arguments),
			ArgumentDetails:    argDetails,
		})
	}

	return result, nil
}

func (s *service) GetExamResults(examID uint, requesterRole string, requesterID uint) (any, error) {
	e, err := s.repo.GetExamByID(examID)
	if err != nil {
		return nil, ErrExamNotFound
	}
	if requesterRole == "asesor" && e.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	sessions, err := s.repo.GetAllSessionsByExam(examID)
	if err != nil {
		return nil, err
	}

	type SessionSummary struct {
		ID                 uint          `json:"id"`
		StudentID          uint          `json:"student_id"`
		StudentName        string        `json:"student_name"`
		Username           string        `json:"username"`
		AttemptNo          int           `json:"attempt_no"`
		Mode               SessionMode   `json:"mode"`
		Status             SessionStatus `json:"status"`
		StartedAt          time.Time     `json:"started_at"`
		CompletedAt        *time.Time    `json:"completed_at,omitempty"`
		TotalAttempts      int           `json:"total_attempts"`
		CompletedArguments int           `json:"completed_arguments"`
		TotalArguments     int           `json:"total_arguments"`
	}

	var list []SessionSummary
	for _, sess := range sessions {
		logs, _ := s.repo.GetAttemptLogsBySession(sess.ID)
		name := ""
		uname := ""
		if sess.Student != nil {
			name = sess.Student.Name
			uname = sess.Student.Username
		}
		list = append(list, SessionSummary{
			ID:                 sess.ID,
			StudentID:          sess.StudentID,
			StudentName:        name,
			Username:           uname,
			AttemptNo:          sess.AttemptNo,
			Mode:               sess.Mode,
			Status:             sess.Status,
			StartedAt:          sess.StartedAt,
			CompletedAt:        sess.CompletedAt,
			TotalAttempts:      len(logs),
			CompletedArguments: len(sess.Progress),
			TotalArguments:     len(sess.Arguments),
		})
	}

	return list, nil
}

func (s *service) GetExamLogs(examID uint, requesterRole string, requesterID uint) (any, error) {
	e, err := s.repo.GetExamByID(examID)
	if err != nil {
		return nil, ErrExamNotFound
	}
	if requesterRole == "asesor" && e.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	sessions, err := s.repo.GetAllSessionsByExam(examID)
	if err != nil {
		return nil, err
	}

	type LogItem struct {
		ID          uint      `json:"id"`
		SessionID   uint      `json:"session_id"`
		AttemptNo   int       `json:"attempt_no"`
		StudentID   uint      `json:"student_id"`
		StudentName string    `json:"student_name"`
		ArgumentID  uint      `json:"argument_id"`
		OptionID    uint      `json:"option_id"`
		OptionText  string    `json:"option_text"`
		Slot        string    `json:"slot"`
		CreatedAt   time.Time `json:"created_at"`
	}

	var allLogs []LogItem
	for _, sess := range sessions {
		logs, err := s.repo.GetAttemptLogsBySession(sess.ID)
		if err != nil {
			continue
		}
		name := ""
		if sess.Student != nil {
			name = sess.Student.Name
		}
		for _, l := range logs {
			text := ""
			if l.Option != nil {
				text = l.Option.Text
			}
			allLogs = append(allLogs, LogItem{
				ID:          l.ID,
				SessionID:   sess.ID,
				AttemptNo:   sess.AttemptNo,
				StudentID:   sess.StudentID,
				StudentName: name,
				ArgumentID:  l.ArgumentID,
				OptionID:    l.OptionID,
				OptionText:  text,
				Slot:        l.Slot,
				CreatedAt:   l.CreatedAt,
			})
		}
	}

	return allLogs, nil
}

func (s *service) GetStaffAnalytics(examID uint, requesterRole string, requesterID uint) (*AnalysisAnalyticsResponse, error) {
	e, err := s.repo.GetExamByID(examID)
	if err != nil {
		return nil, ErrExamNotFound
	}
	if requesterRole == "asesor" && e.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	uniquePeers, err := s.repo.CountUniquePeersWithSessions(examID)
	if err != nil {
		return nil, err
	}

	args, err := s.repo.GetArgumentsByMaterialID(e.MaterialID)
	if err != nil {
		return nil, err
	}

	resp := &AnalysisAnalyticsResponse{
		ExamID:         examID,
		MinPeers:       0,
		TotalPeers:     uniquePeers,
		HasEnoughPeers: true, // Staff always sees full data
	}

	for _, a := range args {
		var optStats []AnalysisOptionStat
		for _, opt := range a.Options {
			totalAtt, uniqUsers, err := s.repo.GetGroupAttemptsByOption(examID, opt.ID)
			if err != nil {
				return nil, err
			}
			students, _ := s.repo.GetStudentsChoosingOption(examID, opt.ID)
			optStats = append(optStats, AnalysisOptionStat{
				OptionID:       opt.ID,
				Type:           opt.Type,
				Text:           opt.Text,
				TotalAttempts:  totalAtt,
				UniqueStudents: uniqUsers,
				Ratio:          fmt.Sprintf("%d:%d", totalAtt, uniqUsers),
				Students:       students,
			})
		}
		resp.Arguments = append(resp.Arguments, AnalysisArgumentStat{
			ArgumentID: a.ID,
			ClaimText:  a.ClaimText,
			OrderNo:    a.OrderNo,
			Options:    optStats,
		})
	}

	return resp, nil
}

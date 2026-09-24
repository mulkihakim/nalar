package session

import (
	"errors"
	"time"

	"github.com/mulkihakim/nalar/backend/internal/exam"
	"github.com/mulkihakim/nalar/backend/internal/material"
	"gorm.io/gorm"
)

var (
	ErrSessionNotFound = errors.New("sesi tidak ditemukan")
)

type Repository interface {
	GetActiveSession(examID, studentID uint) (*Session, error)
	GetSessionByID(sessionID uint) (*Session, error)
	GetMaxAttemptNo(examID, studentID uint) (int, error)
	CreateSession(session *Session, argumentIDs []uint) error
	UpdateSessionMode(sessionID uint, mode SessionMode) error
	CompleteSession(sessionID uint) error
	LogAttempt(log *AttemptLog) error
	GetArgumentProgress(sessionID uint) ([]ArgumentProgress, error)
	MarkArgumentComplete(sessionID, argumentID uint) error
	IsArgumentCompleted(sessionID, argumentID uint) (bool, error)
	CountStudentAttemptsByOption(examID, studentID, optionID uint) (int64, error)
	GetGroupAttemptsByOption(examID, optionID uint) (totalAttempts int64, uniqueStudents int64, err error)
	CountUniquePeersWithSessions(examID uint) (int64, error)
	GetSessionsByExamAndStudent(examID, studentID uint) ([]Session, error)
	GetAllSessionsByExam(examID uint) ([]Session, error)
	GetAttemptLogsBySession(sessionID uint) ([]AttemptLog, error)
	GetStudentExams(studentID uint) ([]exam.Exam, error)
	GetExamByID(examID uint) (*exam.Exam, error)
	GetArgumentWithOptions(argumentID uint) (*material.Argument, error)
	GetArgumentsByMaterialID(materialID uint) ([]material.Argument, error)
	GetStudentsChoosingOption(examID, optionID uint) ([]StudentChooserInfo, error)
	EnsureIndices() error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	r := &repository{db: db}
	_ = r.EnsureIndices()
	return r
}

func (r *repository) EnsureIndices() error {
	// Partial unique index: maksimal 1 session in_progress per (exam_id, student_id)
	return r.db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_sessions_exam_student_in_progress 
		ON sessions (exam_id, student_id) 
		WHERE status = 'in_progress'
	`).Error
}

func (r *repository) GetActiveSession(examID, studentID uint) (*Session, error) {
	var s Session
	err := r.db.Preload("Arguments.Argument.Options").
		Preload("Progress").
		Preload("Exam.Material").
		Where("exam_id = ? AND student_id = ? AND status = ?", examID, studentID, StatusInProgress).
		First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *repository) GetSessionByID(sessionID uint) (*Session, error) {
	var s Session
	err := r.db.Preload("Arguments.Argument.Options").
		Preload("Progress").
		Preload("Exam.Material").
		Preload("Student").
		First(&s, sessionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *repository) GetMaxAttemptNo(examID, studentID uint) (int, error) {
	var maxAttempt *int
	err := r.db.Model(&Session{}).
		Where("exam_id = ? AND student_id = ?", examID, studentID).
		Select("MAX(attempt_no)").
		Scan(&maxAttempt).Error
	if err != nil {
		return 0, err
	}
	if maxAttempt == nil {
		return 0, nil
	}
	return *maxAttempt, nil
}

func (r *repository) CreateSession(s *Session, argumentIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(s).Error; err != nil {
			return err
		}
		for i, argID := range argumentIDs {
			sa := SessionArgument{
				SessionID:  s.ID,
				ArgumentID: argID,
				OrderNo:    i + 1,
			}
			if err := tx.Create(&sa).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) UpdateSessionMode(sessionID uint, mode SessionMode) error {
	return r.db.Model(&Session{}).Where("id = ?", sessionID).Update("mode", mode).Error
}

func (r *repository) CompleteSession(sessionID uint) error {
	now := time.Now()
	return r.db.Model(&Session{}).
		Where("id = ?", sessionID).
		Updates(map[string]any{
			"status":       StatusCompleted,
			"completed_at": now,
		}).Error
}

func (r *repository) LogAttempt(log *AttemptLog) error {
	return r.db.Create(log).Error
}

func (r *repository) GetArgumentProgress(sessionID uint) ([]ArgumentProgress, error) {
	var list []ArgumentProgress
	err := r.db.Where("session_id = ?", sessionID).Find(&list).Error
	return list, err
}

func (r *repository) MarkArgumentComplete(sessionID, argumentID uint) error {
	prog := ArgumentProgress{
		SessionID:   sessionID,
		ArgumentID:  argumentID,
		CompletedAt: time.Now(),
	}
	// On conflict do nothing
	return r.db.Where(ArgumentProgress{SessionID: sessionID, ArgumentID: argumentID}).
		FirstOrCreate(&prog).Error
}

func (r *repository) IsArgumentCompleted(sessionID, argumentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&ArgumentProgress{}).
		Where("session_id = ? AND argument_id = ?", sessionID, argumentID).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) CountStudentAttemptsByOption(examID, studentID, optionID uint) (int64, error) {
	var count int64
	err := r.db.Model(&AttemptLog{}).
		Joins("JOIN sessions ON sessions.id = attempt_logs.session_id").
		Where("sessions.exam_id = ? AND sessions.student_id = ? AND attempt_logs.option_id = ?", examID, studentID, optionID).
		Count(&count).Error
	return count, err
}

func (r *repository) GetGroupAttemptsByOption(examID, optionID uint) (totalAttempts int64, uniqueStudents int64, err error) {
	// Total attempts on this option by all students in this exam
	err = r.db.Model(&AttemptLog{}).
		Joins("JOIN sessions ON sessions.id = attempt_logs.session_id").
		Where("sessions.exam_id = ? AND attempt_logs.option_id = ?", examID, optionID).
		Count(&totalAttempts).Error
	if err != nil {
		return 0, 0, err
	}

	// Unique students who chose this option in this exam
	err = r.db.Model(&AttemptLog{}).
		Joins("JOIN sessions ON sessions.id = attempt_logs.session_id").
		Where("sessions.exam_id = ? AND attempt_logs.option_id = ?", examID, optionID).
		Distinct("sessions.student_id").
		Count(&uniqueStudents).Error
	if err != nil {
		return 0, 0, err
	}

	return totalAttempts, uniqueStudents, nil
}

func (r *repository) CountUniquePeersWithSessions(examID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Session{}).
		Where("exam_id = ?", examID).
		Distinct("student_id").
		Count(&count).Error
	return count, err
}

func (r *repository) GetSessionsByExamAndStudent(examID, studentID uint) ([]Session, error) {
	var sessions []Session
	err := r.db.Preload("Arguments.Argument").
		Preload("Progress").
		Where("exam_id = ? AND student_id = ?", examID, studentID).
		Order("attempt_no ASC").
		Find(&sessions).Error
	return sessions, err
}

func (r *repository) GetAllSessionsByExam(examID uint) ([]Session, error) {
	var sessions []Session
	err := r.db.Preload("Student").
		Preload("Arguments.Argument").
		Preload("Progress").
		Where("exam_id = ?", examID).
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, err
}

func (r *repository) GetAttemptLogsBySession(sessionID uint) ([]AttemptLog, error) {
	var logs []AttemptLog
	err := r.db.Preload("Option").
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&logs).Error
	return logs, err
}

func (r *repository) GetStudentExams(studentID uint) ([]exam.Exam, error) {
	var exams []exam.Exam
	err := r.db.Preload("Material").
		Where("is_active = true AND (id IN (SELECT exam_id FROM exam_students WHERE student_id = ?) OR id IN (SELECT exam_id FROM exam_classes ec JOIN class_members cm ON ec.class_id = cm.class_id WHERE cm.user_id = ?))", studentID, studentID).
		Find(&exams).Error
	return exams, err
}

func (r *repository) GetExamByID(examID uint) (*exam.Exam, error) {
	var e exam.Exam
	err := r.db.Preload("Material.Arguments.Options").
		Preload("Classes").
		Preload("Students").
		First(&e, examID).Error
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repository) GetArgumentWithOptions(argumentID uint) (*material.Argument, error) {
	var a material.Argument
	err := r.db.Preload("Options").First(&a, argumentID).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repository) GetArgumentsByMaterialID(materialID uint) ([]material.Argument, error) {
	var list []material.Argument
	err := r.db.Preload("Options").
		Where("material_id = ?", materialID).
		Order("order_no ASC").
		Find(&list).Error
	return list, err
}

func (r *repository) GetStudentsChoosingOption(examID, optionID uint) ([]StudentChooserInfo, error) {
	var list []StudentChooserInfo
	err := r.db.Table("attempt_logs").
		Select("users.id as student_id, users.name as student_name, users.username, COUNT(attempt_logs.id) as attempts, MAX(attempt_logs.created_at) as last_attempt").
		Joins("JOIN sessions ON sessions.id = attempt_logs.session_id").
		Joins("JOIN users ON users.id = sessions.student_id").
		Where("sessions.exam_id = ? AND attempt_logs.option_id = ?", examID, optionID).
		Group("users.id, users.name, users.username").
		Order("attempts DESC, users.name ASC").
		Scan(&list).Error
	return list, err
}

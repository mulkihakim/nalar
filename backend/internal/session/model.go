package session

import (
	"time"

	"github.com/mulkihakim/nalar/backend/internal/exam"
	"github.com/mulkihakim/nalar/backend/internal/material"
	"github.com/mulkihakim/nalar/backend/internal/user"
)

type SessionMode string
const (
	ModeStandard SessionMode = "standard"
	ModeHelp     SessionMode = "help"
	ModeSocial   SessionMode = "social"
)

type SessionStatus string
const (
	StatusInProgress SessionStatus = "in_progress"
	StatusCompleted  SessionStatus = "completed"
)

type Session struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	ExamID      uint             `gorm:"not null;index" json:"exam_id"`
	Exam        *exam.Exam       `gorm:"foreignKey:ExamID;constraint:OnDelete:CASCADE" json:"exam,omitempty"`
	StudentID   uint             `gorm:"not null;index" json:"student_id"`
	Student     *user.User       `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"student,omitempty"`
	AttemptNo   int              `gorm:"not null;default:1" json:"attempt_no"`
	Mode        SessionMode      `gorm:"not null;default:'standard'" json:"mode"`
	Status      SessionStatus    `gorm:"not null;default:'in_progress'" json:"status"`
	StartedAt   time.Time        `gorm:"not null" json:"started_at"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
	Arguments   []SessionArgument `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE" json:"arguments,omitempty"`
	Progress    []ArgumentProgress `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE" json:"progress,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type SessionArgument struct {
	ID         uint               `gorm:"primaryKey" json:"id"`
	SessionID  uint               `gorm:"not null;uniqueIndex:idx_session_arg" json:"session_id"`
	ArgumentID uint               `gorm:"not null;uniqueIndex:idx_session_arg" json:"argument_id"`
	Argument   *material.Argument `gorm:"foreignKey:ArgumentID;constraint:OnDelete:CASCADE" json:"argument,omitempty"`
	OrderNo    int                `gorm:"not null;default:1" json:"order_no"`
	CreatedAt  time.Time          `json:"created_at"`
}

type AttemptLog struct {
	ID         uint             `gorm:"primaryKey" json:"id"`
	SessionID  uint             `gorm:"not null;index" json:"session_id"`
	ArgumentID uint             `gorm:"not null;index" json:"argument_id"`
	OptionID   uint             `gorm:"not null;index" json:"option_id"`
	Option     *material.Option `gorm:"foreignKey:OptionID;constraint:OnDelete:CASCADE" json:"option,omitempty"`
	Slot       string           `gorm:"not null" json:"slot"` // "ground" | "warrant"
	CreatedAt  time.Time        `json:"created_at"`
}

type ArgumentProgress struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	SessionID   uint               `gorm:"not null;uniqueIndex:idx_session_arg_prog" json:"session_id"`
	ArgumentID  uint               `gorm:"not null;uniqueIndex:idx_session_arg_prog" json:"argument_id"`
	Argument    *material.Argument `gorm:"foreignKey:ArgumentID;constraint:OnDelete:CASCADE" json:"argument,omitempty"`
	CompletedAt time.Time          `json:"completed_at"`
}

func (Session) TableName() string {
	return "sessions"
}

func (SessionArgument) TableName() string {
	return "session_arguments"
}

func (AttemptLog) TableName() string {
	return "attempt_logs"
}

func (ArgumentProgress) TableName() string {
	return "argument_progress"
}

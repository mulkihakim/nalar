package exam

import (
	"time"

	"github.com/mulkihakim/nalar/backend/internal/class"
	"github.com/mulkihakim/nalar/backend/internal/material"
	"github.com/mulkihakim/nalar/backend/internal/user"
)

type Exam struct {
	ID                  uint               `gorm:"primaryKey" json:"id"`
	Title               string             `gorm:"not null" json:"title"`
	MaterialID          uint               `gorm:"not null;index" json:"material_id"`
	Material            *material.Material `gorm:"foreignKey:MaterialID;constraint:OnDelete:CASCADE" json:"material,omitempty"`
	OwnerID             uint               `gorm:"not null;index" json:"owner_id"`
	Owner               *user.User         `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE" json:"owner,omitempty"`
	IsActive            bool               `gorm:"not null;default:true" json:"is_active"`
	ArgumentsPerSession int                `gorm:"not null;default:3" json:"arguments_per_session"`
	Classes             []class.Class      `gorm:"many2many:exam_classes;joinForeignKey:ExamID;joinReferences:ClassID" json:"classes,omitempty"`
	Students            []user.User        `gorm:"many2many:exam_students;joinForeignKey:ExamID;joinReferences:StudentID" json:"students,omitempty"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
}

type ExamClass struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	ExamID    uint         `gorm:"not null;uniqueIndex:idx_exam_class;index" json:"exam_id"`
	Exam      *Exam        `gorm:"foreignKey:ExamID;constraint:OnDelete:CASCADE" json:"-"`
	ClassID   uint         `gorm:"not null;uniqueIndex:idx_exam_class;index" json:"class_id"`
	Class     *class.Class `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time    `json:"created_at"`
}

type ExamStudent struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	ExamID    uint       `gorm:"not null;uniqueIndex:idx_exam_student;index" json:"exam_id"`
	Exam      *Exam      `gorm:"foreignKey:ExamID;constraint:OnDelete:CASCADE" json:"-"`
	StudentID uint       `gorm:"not null;uniqueIndex:idx_exam_student;index" json:"student_id"`
	Student   *user.User `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time  `json:"created_at"`
}

func (Exam) TableName() string {
	return "exams"
}

func (ExamClass) TableName() string {
	return "exam_classes"
}

func (ExamStudent) TableName() string {
	return "exam_students"
}

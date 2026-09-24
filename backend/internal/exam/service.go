package exam

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mulkihakim/nalar/backend/internal/material"
)

var (
	ErrExamNotFound             = errors.New("ujian tidak ditemukan")
	ErrForbidden                = errors.New("akses ditolak: bukan pemilik ujian")
	ErrEmptyTitle               = errors.New("judul ujian wajib diisi")
	ErrInvalidArgumentsPerSess = errors.New("arguments_per_session harus lebih besar dari 0")
	ErrMaterialNotFound         = errors.New("materi tidak ditemukan")
)

type CreateExamInput struct {
	Title               string `json:"title"`
	MaterialID          uint   `json:"material_id"`
	ArgumentsPerSession int    `json:"arguments_per_session"`
	IsActive            bool   `json:"is_active"`
	ClassIDs            []uint `json:"class_ids"`
	StudentIDs          []uint `json:"student_ids"`
}

type UpdateExamInput struct {
	Title               string `json:"title"`
	MaterialID          uint   `json:"material_id"`
	ArgumentsPerSession int    `json:"arguments_per_session"`
	IsActive            *bool  `json:"is_active"`
}

type AccessInput struct {
	ClassIDs   []uint `json:"class_ids"`
	StudentIDs []uint `json:"student_ids"`
}

type Service interface {
	CreateExam(input CreateExamInput, ownerID uint, requesterRole string) (*Exam, error)
	GetExamByID(id uint, requesterRole string, requesterID uint) (*Exam, error)
	ListExams(requesterRole string, requesterID uint) ([]Exam, error)
	UpdateExam(id uint, input UpdateExamInput, requesterRole string, requesterID uint) (*Exam, error)
	DeleteExam(id uint, requesterRole string, requesterID uint) error
	SetAccess(id uint, input AccessInput, requesterRole string, requesterID uint) error
	SetStatus(id uint, isActive bool, requesterRole string, requesterID uint) error
}

type service struct {
	repo         Repository
	materialRepo material.Repository
}

func NewService(repo Repository, materialRepo material.Repository) Service {
	return &service{
		repo:         repo,
		materialRepo: materialRepo,
	}
}

func (s *service) CreateExam(input CreateExamInput, ownerID uint, requesterRole string) (*Exam, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrEmptyTitle
	}

	argsPerSession := input.ArgumentsPerSession
	if argsPerSession <= 0 {
		argsPerSession = 3
	}

	// Pastikan materi ada
	mat, err := s.materialRepo.FindMaterialByID(input.MaterialID)
	if err != nil {
		return nil, err
	}
	if mat == nil {
		return nil, ErrMaterialNotFound
	}
	// DoD #1: Asesor hanya bisa membuat ujian menggunakan materi miliknya
	if requesterRole == "asesor" && mat.OwnerID != ownerID {
		return nil, ErrForbidden
	}

	e := &Exam{
		Title:               title,
		MaterialID:          input.MaterialID,
		OwnerID:             ownerID,
		IsActive:            input.IsActive,
		ArgumentsPerSession: argsPerSession,
	}

	if err := s.repo.Create(e, input.ClassIDs, input.StudentIDs); err != nil {
		return nil, err
	}
	return s.repo.FindByID(e.ID)
}

func (s *service) GetExamByID(id uint, requesterRole string, requesterID uint) (*Exam, error) {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, ErrExamNotFound
	}

	// DoD #1: Asesor hanya bisa melihat ujian miliknya
	if requesterRole == "asesor" && e.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	return e, nil
}

func (s *service) ListExams(requesterRole string, requesterID uint) ([]Exam, error) {
	var ownerIDFilter *uint
	// DoD #1: Asesor hanya melihat ujian miliknya
	if requesterRole == "asesor" {
		ownerIDFilter = &requesterID
	}
	return s.repo.List(ownerIDFilter)
}

func (s *service) UpdateExam(id uint, input UpdateExamInput, requesterRole string, requesterID uint) (*Exam, error) {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, ErrExamNotFound
	}

	// DoD #1: Asesor tidak bisa mengubah ujian milik asesor lain
	if requesterRole == "asesor" && e.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	title := strings.TrimSpace(input.Title)
	if title != "" {
		e.Title = title
	}

	if input.ArgumentsPerSession > 0 {
		e.ArgumentsPerSession = input.ArgumentsPerSession
	}

	if input.MaterialID > 0 && input.MaterialID != e.MaterialID {
		mat, err := s.materialRepo.FindMaterialByID(input.MaterialID)
		if err != nil {
			return nil, err
		}
		if mat == nil {
			return nil, ErrMaterialNotFound
		}
		if requesterRole == "asesor" && mat.OwnerID != requesterID {
			return nil, ErrForbidden
		}
		e.MaterialID = input.MaterialID
	}

	if input.IsActive != nil {
		e.IsActive = *input.IsActive
	}

	if err := s.repo.Update(e); err != nil {
		return nil, err
	}
	return s.repo.FindByID(e.ID)
}

func (s *service) DeleteExam(id uint, requesterRole string, requesterID uint) error {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if e == nil {
		return ErrExamNotFound
	}

	// DoD #1: Asesor tidak bisa menghapus ujian milik asesor lain
	if requesterRole == "asesor" && e.OwnerID != requesterID {
		return ErrForbidden
	}

	return s.repo.Delete(id)
}

func (s *service) SetAccess(id uint, input AccessInput, requesterRole string, requesterID uint) error {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if e == nil {
		return ErrExamNotFound
	}

	if requesterRole == "asesor" && e.OwnerID != requesterID {
		return ErrForbidden
	}

	// Cek kelas yang dicabut hak aksesnya
	newClassMap := make(map[uint]bool)
	for _, cid := range input.ClassIDs {
		newClassMap[cid] = true
	}
	for _, existingClass := range e.Classes {
		if !newClassMap[existingClass.ID] {
			taken, err := s.repo.HasClassTakenExam(id, existingClass.ID)
			if err != nil {
				return err
			}
			if taken {
				return fmt.Errorf("kelas '%s' tidak dapat dihapus dari akses karena siswa di dalamnya telah memulai/mengerjakan ujian ini", existingClass.Name)
			}
		}
	}

	// Cek siswa individu yang dicabut hak aksesnya
	newStudentMap := make(map[uint]bool)
	for _, sid := range input.StudentIDs {
		newStudentMap[sid] = true
	}
	for _, existingStudent := range e.Students {
		if !newStudentMap[existingStudent.ID] {
			taken, err := s.repo.HasStudentTakenExam(id, existingStudent.ID)
			if err != nil {
				return err
			}
			if taken {
				return fmt.Errorf("siswa '%s' (@%s) tidak dapat dihapus dari akses karena telah memulai/mengerjakan ujian ini", existingStudent.Name, existingStudent.Username)
			}
		}
	}

	return s.repo.SetAccess(id, input.ClassIDs, input.StudentIDs)
}

func (s *service) SetStatus(id uint, isActive bool, requesterRole string, requesterID uint) error {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if e == nil {
		return ErrExamNotFound
	}

	if requesterRole == "asesor" && e.OwnerID != requesterID {
		return ErrForbidden
	}

	return s.repo.SetStatus(id, isActive)
}

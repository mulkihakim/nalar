package class

import (
	"errors"

	"github.com/mulkihakim/nalar/backend/internal/user"
)

var (
	ErrClassNotFound        = errors.New("kelas tidak ditemukan")
	ErrForbidden            = errors.New("akses ditolak: bukan pemilik kelas")
	ErrInvalidMemberRole    = errors.New("hanya pengguna dengan role siswa yang dapat ditambahkan ke kelas")
	ErrMemberAlreadyExists  = errors.New("siswa sudah menjadi anggota kelas ini")
	ErrStudentNotFound      = errors.New("siswa tidak ditemukan")
	ErrEmptyClassName       = errors.New("nama kelas wajib diisi")
)

type Service interface {
	CreateClass(name string, ownerID uint) (*Class, error)
	GetClassByID(id uint, requesterRole string, requesterID uint) (*Class, error)
	ListClasses(requesterRole string, requesterID uint) ([]Class, error)
	UpdateClass(id uint, name string, requesterRole string, requesterID uint) (*Class, error)
	DeleteClass(id uint, requesterRole string, requesterID uint) error
	AddMember(classID uint, studentID uint, requesterRole string, requesterID uint) error
	RemoveMember(classID uint, studentID uint, requesterRole string, requesterID uint) error
	GetMembers(classID uint, requesterRole string, requesterID uint) ([]user.User, error)
}

type service struct {
	repo     Repository
	userRepo user.Repository
}

func NewService(repo Repository, userRepo user.Repository) Service {
	return &service{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *service) CreateClass(name string, ownerID uint) (*Class, error) {
	if name == "" {
		return nil, ErrEmptyClassName
	}
	c := &Class{
		Name:    name,
		OwnerID: ownerID,
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return s.repo.FindByID(c.ID)
}

func (s *service) GetClassByID(id uint, requesterRole string, requesterID uint) (*Class, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrClassNotFound
	}

	// DoD #1: Asesor hanya bisa melihat kelas miliknya
	if requesterRole == "asesor" && c.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	return c, nil
}

func (s *service) ListClasses(requesterRole string, requesterID uint) ([]Class, error) {
	var ownerIDFilter *uint
	// DoD #1: Asesor hanya bisa melihat kelas miliknya
	if requesterRole == "asesor" {
		ownerIDFilter = &requesterID
	}
	return s.repo.List(ownerIDFilter)
}

func (s *service) UpdateClass(id uint, name string, requesterRole string, requesterID uint) (*Class, error) {
	if name == "" {
		return nil, ErrEmptyClassName
	}
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrClassNotFound
	}

	// DoD #1: Asesor tidak bisa mengubah kelas milik asesor lain
	if requesterRole == "asesor" && c.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	c.Name = name
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return s.repo.FindByID(c.ID)
}

func (s *service) DeleteClass(id uint, requesterRole string, requesterID uint) error {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if c == nil {
		return ErrClassNotFound
	}

	// DoD #1: Asesor tidak bisa menghapus kelas milik asesor lain
	if requesterRole == "asesor" && c.OwnerID != requesterID {
		return ErrForbidden
	}

	return s.repo.Delete(id)
}

func (s *service) AddMember(classID uint, studentID uint, requesterRole string, requesterID uint) error {
	c, err := s.repo.FindByID(classID)
	if err != nil {
		return err
	}
	if c == nil {
		return ErrClassNotFound
	}

	if requesterRole == "asesor" && c.OwnerID != requesterID {
		return ErrForbidden
	}

	// Verifikasi calon anggota
	student, err := s.userRepo.FindByID(studentID)
	if err != nil {
		return err
	}
	if student == nil {
		return ErrStudentNotFound
	}
	if student.Role != "siswa" {
		return ErrInvalidMemberRole
	}

	isMember, err := s.repo.IsMember(classID, studentID)
	if err != nil {
		return err
	}
	if isMember {
		return ErrMemberAlreadyExists
	}

	return s.repo.AddMember(classID, studentID)
}

func (s *service) RemoveMember(classID uint, studentID uint, requesterRole string, requesterID uint) error {
	c, err := s.repo.FindByID(classID)
	if err != nil {
		return err
	}
	if c == nil {
		return ErrClassNotFound
	}

	if requesterRole == "asesor" && c.OwnerID != requesterID {
		return ErrForbidden
	}

	return s.repo.RemoveMember(classID, studentID)
}

func (s *service) GetMembers(classID uint, requesterRole string, requesterID uint) ([]user.User, error) {
	c, err := s.repo.FindByID(classID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrClassNotFound
	}

	if requesterRole == "asesor" && c.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	return s.repo.GetMembers(classID)
}

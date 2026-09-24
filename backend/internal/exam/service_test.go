package exam_test

import (
	"errors"
	"testing"

	"github.com/mulkihakim/nalar/backend/internal/exam"
	"github.com/mulkihakim/nalar/backend/internal/material"
	"github.com/mulkihakim/nalar/backend/internal/user"
)

type mockExamRepo struct {
	exams map[uint]*exam.Exam
}

func newMockExamRepo() *mockExamRepo {
	return &mockExamRepo{
		exams: make(map[uint]*exam.Exam),
	}
}

func (m *mockExamRepo) Create(e *exam.Exam, classIDs []uint, studentIDs []uint) error {
	e.ID = uint(len(m.exams) + 1)
	m.exams[e.ID] = e
	return nil
}

func (m *mockExamRepo) FindByID(id uint) (*exam.Exam, error) {
	e, ok := m.exams[id]
	if !ok {
		return nil, nil
	}
	return e, nil
}

func (m *mockExamRepo) List(ownerID *uint) ([]exam.Exam, error) {
	var res []exam.Exam
	for _, e := range m.exams {
		if ownerID == nil || e.OwnerID == *ownerID {
			res = append(res, *e)
		}
	}
	return res, nil
}

func (m *mockExamRepo) Update(e *exam.Exam) error {
	m.exams[e.ID] = e
	return nil
}

func (m *mockExamRepo) Delete(id uint) error {
	delete(m.exams, id)
	return nil
}

func (m *mockExamRepo) SetAccess(examID uint, classIDs []uint, studentIDs []uint) error {
	return nil
}

func (m *mockExamRepo) SetStatus(examID uint, isActive bool) error {
	if e, ok := m.exams[examID]; ok {
		e.IsActive = isActive
	}
	return nil
}

func (m *mockExamRepo) HasStudentTakenExam(examID uint, studentID uint) (bool, error) {
	// Misal student 4 sudah pernah mengerjakan exam 1
	if examID == 1 && studentID == 4 {
		return true, nil
	}
	return false, nil
}

func (m *mockExamRepo) HasClassTakenExam(examID uint, classID uint) (bool, error) {
	// Misal class 99 ada siswa yang sudah mengerjakan exam 1
	if examID == 1 && classID == 99 {
		return true, nil
	}
	return false, nil
}

type mockMaterialRepo struct {
	materials map[uint]*material.Material
}

func (m *mockMaterialRepo) CreateMaterial(mat *material.Material) error { return nil }
func (m *mockMaterialRepo) FindMaterialByID(id uint) (*material.Material, error) {
	return m.materials[id], nil
}
func (m *mockMaterialRepo) ListMaterials(ownerID *uint) ([]material.Material, error) { return nil, nil }
func (m *mockMaterialRepo) UpdateMaterial(mat *material.Material) error               { return nil }
func (m *mockMaterialRepo) DeleteMaterial(id uint) error                             { return nil }
func (m *mockMaterialRepo) CreateArgument(arg *material.Argument, options []material.Option) error {
	return nil
}
func (m *mockMaterialRepo) FindArgumentByID(id uint) (*material.Argument, error) { return nil, nil }
func (m *mockMaterialRepo) ListArguments(materialID uint) ([]material.Argument, error) {
	return nil, nil
}
func (m *mockMaterialRepo) UpdateArgument(arg *material.Argument, options []material.Option) error {
	return nil
}
func (m *mockMaterialRepo) DeleteArgument(id uint) error { return nil }

func TestExamAssessorIsolation_DoD1(t *testing.T) {
	examRepo := newMockExamRepo()
	materialRepo := &mockMaterialRepo{
		materials: map[uint]*material.Material{
			1: {ID: 1, Title: "Materi Asesor 2", OwnerID: 2},
			2: {ID: 2, Title: "Materi Asesor 3", OwnerID: 3},
		},
	}
	svc := exam.NewService(examRepo, materialRepo)

	// 1. Asesor 3 mencoba membuat ujian memakai Materi milik Asesor 2 -> Harus ErrForbidden
	_, err := svc.CreateExam(exam.CreateExamInput{
		Title:      "Ujian Ilegal",
		MaterialID: 1, // Milik asesor 2
	}, 3, "asesor")
	if !errors.Is(err, exam.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor creates exam using another's material, got %v", err)
	}

	// 2. Asesor 2 membuat ujian dengan materinya sendiri (MaterialID 1)
	e2, err := svc.CreateExam(exam.CreateExamInput{
		Title:      "Ujian Asesor 2",
		MaterialID: 1,
	}, 2, "asesor")
	if err != nil {
		t.Fatalf("expected create exam to succeed, got %v", err)
	}

	// 3. Asesor 3 mencoba melihat ujian milik Asesor 2 -> Harus ErrForbidden
	_, err = svc.GetExamByID(e2.ID, "asesor", 3)
	if !errors.Is(err, exam.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor 3 gets asesor 2's exam, got %v", err)
	}

	// 4. Asesor 3 mencoba update ujian milik Asesor 2 -> Harus ErrForbidden
	_, err = svc.UpdateExam(e2.ID, exam.UpdateExamInput{
		Title: "Hacked Exam",
	}, "asesor", 3)
	if !errors.Is(err, exam.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor 3 updates asesor 2's exam, got %v", err)
	}

	// 5. Asesor 3 mencoba set access pada ujian milik Asesor 2 -> Harus ErrForbidden
	err = svc.SetAccess(e2.ID, exam.AccessInput{ClassIDs: []uint{1}}, "asesor", 3)
	if !errors.Is(err, exam.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor 3 sets access on asesor 2's exam, got %v", err)
	}

	// 6. Asesor 3 mencoba set status pada ujian milik Asesor 2 -> Harus ErrForbidden
	err = svc.SetStatus(e2.ID, false, "asesor", 3)
	if !errors.Is(err, exam.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor 3 sets status on asesor 2's exam, got %v", err)
	}

	// 7. Asesor 3 mencoba menghapus ujian milik Asesor 2 -> Harus ErrForbidden
	err = svc.DeleteExam(e2.ID, "asesor", 3)
	if !errors.Is(err, exam.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor 3 deletes asesor 2's exam, got %v", err)
	}

	// 8. Admin boleh melihat dan mengatur ujian milik Asesor 2
	adminGet, err := svc.GetExamByID(e2.ID, "admin", 1)
	if err != nil || adminGet == nil {
		t.Errorf("expected admin to view exam, got err: %v", err)
	}
}

func TestExamAccessParticipantRestriction(t *testing.T) {
	examRepo := newMockExamRepo()
	materialRepo := &mockMaterialRepo{
		materials: map[uint]*material.Material{
			1: {ID: 1, Title: "Materi Asesor 2", OwnerID: 2},
		},
	}
	svc := exam.NewService(examRepo, materialRepo)

	// Buat exam dengan siswa 4 terdaftar
	e, err := svc.CreateExam(exam.CreateExamInput{
		Title:      "Ujian dengan Siswa Aktif",
		MaterialID: 1,
		StudentIDs: []uint{4},
	}, 2, "asesor")
	if err != nil {
		t.Fatalf("failed to create exam: %v", err)
	}

	// Update mock repo to have student 4 in e.Students
	e.Students = []user.User{{ID: 4, Name: "Siti Siswa"}}
	_ = examRepo.Update(e)

	// Coba hapus siswa 4 dari akses (kirim StudentIDs kosong) -> Harus gagal karena siswa 4 sudah pernah mengerjakan
	err = svc.SetAccess(e.ID, exam.AccessInput{
		StudentIDs: []uint{},
	}, "asesor", 2)
	if err == nil {
		t.Error("expected error when removing student who has taken exam, got nil")
	}
}

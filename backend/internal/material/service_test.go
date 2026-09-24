package material_test

import (
	"errors"
	"testing"

	"github.com/mulkihakim/nalar/backend/internal/material"
)

type mockMaterialRepo struct {
	materials map[uint]*material.Material
	arguments map[uint]*material.Argument
}

func newMockMaterialRepo() *mockMaterialRepo {
	return &mockMaterialRepo{
		materials: make(map[uint]*material.Material),
		arguments: make(map[uint]*material.Argument),
	}
}

func (m *mockMaterialRepo) CreateMaterial(mat *material.Material) error {
	mat.ID = uint(len(m.materials) + 1)
	m.materials[mat.ID] = mat
	return nil
}

func (m *mockMaterialRepo) FindMaterialByID(id uint) (*material.Material, error) {
	mat, ok := m.materials[id]
	if !ok {
		return nil, nil
	}
	return mat, nil
}

func (m *mockMaterialRepo) ListMaterials(ownerID *uint) ([]material.Material, error) {
	var res []material.Material
	for _, mat := range m.materials {
		if ownerID == nil || mat.OwnerID == *ownerID {
			res = append(res, *mat)
		}
	}
	return res, nil
}

func (m *mockMaterialRepo) UpdateMaterial(mat *material.Material) error {
	m.materials[mat.ID] = mat
	return nil
}

func (m *mockMaterialRepo) DeleteMaterial(id uint) error {
	delete(m.materials, id)
	return nil
}

func (m *mockMaterialRepo) CreateArgument(arg *material.Argument, options []material.Option) error {
	arg.ID = uint(len(m.arguments) + 1)
	arg.Options = options
	m.arguments[arg.ID] = arg
	return nil
}

func (m *mockMaterialRepo) FindArgumentByID(id uint) (*material.Argument, error) {
	arg, ok := m.arguments[id]
	if !ok {
		return nil, nil
	}
	return arg, nil
}

func (m *mockMaterialRepo) ListArguments(materialID uint) ([]material.Argument, error) {
	var res []material.Argument
	for _, arg := range m.arguments {
		if arg.MaterialID == materialID {
			res = append(res, *arg)
		}
	}
	return res, nil
}

func (m *mockMaterialRepo) UpdateArgument(arg *material.Argument, options []material.Option) error {
	arg.Options = options
	m.arguments[arg.ID] = arg
	return nil
}

func (m *mockMaterialRepo) DeleteArgument(id uint) error {
	delete(m.arguments, id)
	return nil
}

func buildValidOptions() []material.OptionInput {
	return []material.OptionInput{
		{Type: "ground", Text: "Ground 1 (Benar)", IsCorrect: true},
		{Type: "ground", Text: "Ground 2", IsCorrect: false},
		{Type: "ground", Text: "Ground 3", IsCorrect: false},
		{Type: "ground", Text: "Ground 4", IsCorrect: false},
		{Type: "warrant", Text: "Warrant 1 (Benar)", IsCorrect: true},
		{Type: "warrant", Text: "Warrant 2", IsCorrect: false},
		{Type: "warrant", Text: "Warrant 3", IsCorrect: false},
		{Type: "warrant", Text: "Warrant 4", IsCorrect: false},
	}
}

// Test DoD #2: Validasi minimal 3 ground (1 benar) dan minimal 3 warrant (1 benar)
func TestArgumentValidation_DoD2(t *testing.T) {
	repo := newMockMaterialRepo()
	svc := material.NewService(repo)

	m, err := svc.CreateMaterial(material.CreateMaterialInput{
		Title:   "Materi Logika Toulmin",
		Content: "Teks bacaan lengkap...",
	}, 1)
	if err != nil {
		t.Fatalf("failed to create material: %v", err)
	}

	// 1. Kasus Sukses: 3 ground (1 benar) dan 3 warrant (1 benar)
	min3Options := []material.OptionInput{
		{Type: "ground", Text: "Ground 1 (Benar)", IsCorrect: true},
		{Type: "ground", Text: "Ground 2", IsCorrect: false},
		{Type: "ground", Text: "Ground 3", IsCorrect: false},
		{Type: "warrant", Text: "Warrant 1 (Benar)", IsCorrect: true},
		{Type: "warrant", Text: "Warrant 2", IsCorrect: false},
		{Type: "warrant", Text: "Warrant 3", IsCorrect: false},
	}
	arg3, err := svc.CreateArgument(m.ID, material.ArgumentInput{
		ClaimText: "Klaim Minimal 3",
		Options:   min3Options,
	}, "admin", 1)
	if err != nil {
		t.Fatalf("expected valid 3+3 argument to succeed, got %v", err)
	}
	if arg3 == nil || len(arg3.Options) != 6 {
		t.Errorf("expected argument with 6 options, got %v", arg3)
	}

	// 2. Kasus Sukses: 4 ground (1 benar) dan 4 warrant (1 benar)
	validInput := material.ArgumentInput{
		ClaimText: "Klaim 4+4",
		Options:   buildValidOptions(),
	}
	arg4, err := svc.CreateArgument(m.ID, validInput, "admin", 1)
	if err != nil {
		t.Fatalf("expected valid 4+4 argument to succeed, got %v", err)
	}
	if arg4 == nil || len(arg4.Options) != 8 {
		t.Errorf("expected argument with 8 options, got %v", arg4)
	}

	// 3. Ground kurang dari 3 (cuma 2 ground, 3 warrant) -> Ditolak
	shortGround := []material.OptionInput{
		{Type: "ground", Text: "Ground 1 (Benar)", IsCorrect: true},
		{Type: "ground", Text: "Ground 2", IsCorrect: false},
		{Type: "warrant", Text: "Warrant 1 (Benar)", IsCorrect: true},
		{Type: "warrant", Text: "Warrant 2", IsCorrect: false},
		{Type: "warrant", Text: "Warrant 3", IsCorrect: false},
	}
	_, err = svc.CreateArgument(m.ID, material.ArgumentInput{
		ClaimText: "Klaim Kurang Ground",
		Options:   shortGround,
	}, "admin", 1)
	if !errors.Is(err, material.ErrInvalidOptionsCount) {
		t.Errorf("expected ErrInvalidOptionsCount for <3 ground options, got %v", err)
	}

	// 4. Ground tidak ada yang benar (0 benar) -> Ditolak
	noCorrectGround := min3Options
	noCorrectGround[0].IsCorrect = false
	_, err = svc.CreateArgument(m.ID, material.ArgumentInput{
		ClaimText: "Klaim Tanpa Ground Benar",
		Options:   noCorrectGround,
	}, "admin", 1)
	if !errors.Is(err, material.ErrInvalidOptionsCount) {
		t.Errorf("expected ErrInvalidOptionsCount when 0 ground is correct, got %v", err)
	}

	// 5. Warrant punya 2 jawaban benar -> Ditolak
	twoCorrectWarrant := buildValidOptions()
	twoCorrectWarrant[5].IsCorrect = true // warrant 1 dan 2 true
	_, err = svc.CreateArgument(m.ID, material.ArgumentInput{
		ClaimText: "Klaim 2 Warrant Benar",
		Options:   twoCorrectWarrant,
	}, "admin", 1)
	if !errors.Is(err, material.ErrInvalidOptionsCount) {
		t.Errorf("expected ErrInvalidOptionsCount when 2 warrants are correct, got %v", err)
	}
}

// Test DoD #1: Asesor tidak bisa melihat/mengubah materi milik asesor lain
func TestMaterialAssessorIsolation_DoD1(t *testing.T) {
	repo := newMockMaterialRepo()
	svc := material.NewService(repo)

	// Asesor 2 membuat materi
	m2, err := svc.CreateMaterial(material.CreateMaterialInput{
		Title:   "Materi Milik Asesor 2",
		Content: "Isi materi...",
	}, 2)
	if err != nil {
		t.Fatalf("failed to create material: %v", err)
	}

	// Asesor 3 mencoba get materi milik Asesor 2 -> ErrForbidden
	_, err = svc.GetMaterialByID(m2.ID, "asesor", 3)
	if !errors.Is(err, material.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor 3 gets asesor 2's material, got %v", err)
	}

	// Asesor 3 mencoba update materi milik Asesor 2 -> ErrForbidden
	_, err = svc.UpdateMaterial(m2.ID, material.CreateMaterialInput{
		Title:   "Hacked Title",
		Content: "Hacked Content",
	}, "asesor", 3)
	if !errors.Is(err, material.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor 3 updates asesor 2's material, got %v", err)
	}

	// Asesor 3 mencoba delete materi milik Asesor 2 -> ErrForbidden
	err = svc.DeleteMaterial(m2.ID, "asesor", 3)
	if !errors.Is(err, material.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor 3 deletes asesor 2's material, got %v", err)
	}

	// Asesor 3 mencoba tambah argumen ke materi milik Asesor 2 -> ErrForbidden
	_, err = svc.CreateArgument(m2.ID, material.ArgumentInput{
		ClaimText: "Klaim",
		Options:   buildValidOptions(),
	}, "asesor", 3)
	if !errors.Is(err, material.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor 3 adds argument to asesor 2's material, got %v", err)
	}

	// Admin bebas akses
	adminGet, err := svc.GetMaterialByID(m2.ID, "admin", 1)
	if err != nil || adminGet == nil {
		t.Errorf("expected admin to get material, got %v", err)
	}
}

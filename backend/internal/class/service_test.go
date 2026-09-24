package class_test

import (
	"errors"
	"testing"

	"github.com/mulkihakim/nalar/backend/internal/class"
	"github.com/mulkihakim/nalar/backend/internal/user"
)

type mockClassRepo struct {
	classes map[uint]*class.Class
	members map[uint]map[uint]bool // classID -> map[userID]bool
}

func newMockClassRepo() *mockClassRepo {
	return &mockClassRepo{
		classes: make(map[uint]*class.Class),
		members: make(map[uint]map[uint]bool),
	}
}

func (m *mockClassRepo) Create(c *class.Class) error {
	c.ID = uint(len(m.classes) + 1)
	m.classes[c.ID] = c
	return nil
}

func (m *mockClassRepo) FindByID(id uint) (*class.Class, error) {
	c, ok := m.classes[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (m *mockClassRepo) List(ownerID *uint) ([]class.Class, error) {
	var res []class.Class
	for _, c := range m.classes {
		if ownerID == nil || c.OwnerID == *ownerID {
			res = append(res, *c)
		}
	}
	return res, nil
}

func (m *mockClassRepo) Update(c *class.Class) error {
	m.classes[c.ID] = c
	return nil
}

func (m *mockClassRepo) Delete(id uint) error {
	delete(m.classes, id)
	delete(m.members, id)
	return nil
}

func (m *mockClassRepo) AddMember(classID uint, userID uint) error {
	if m.members[classID] == nil {
		m.members[classID] = make(map[uint]bool)
	}
	m.members[classID][userID] = true
	return nil
}

func (m *mockClassRepo) RemoveMember(classID uint, userID uint) error {
	if m.members[classID] != nil {
		delete(m.members[classID], userID)
	}
	return nil
}

func (m *mockClassRepo) GetMembers(classID uint) ([]user.User, error) {
	return nil, nil
}

func (m *mockClassRepo) IsMember(classID uint, userID uint) (bool, error) {
	if m.members[classID] == nil {
		return false, nil
	}
	return m.members[classID][userID], nil
}

type mockUserRepo struct {
	users map[uint]*user.User
}

func (m *mockUserRepo) Create(u *user.User) error                 { return nil }
func (m *mockUserRepo) FindByID(id uint) (*user.User, error)      { return m.users[id], nil }
func (m *mockUserRepo) FindByUsername(un string) (*user.User, error) { return nil, nil }
func (m *mockUserRepo) List() ([]user.User, error)               { return nil, nil }
func (m *mockUserRepo) Update(u *user.User) error                 { return nil }
func (m *mockUserRepo) Delete(id uint) error                       { return nil }
func (m *mockUserRepo) HasExamSessions(userID uint) (bool, error) { return false, nil }

func TestClassAssessorIsolation_DoD1(t *testing.T) {
	classRepo := newMockClassRepo()
	userRepo := &mockUserRepo{
		users: map[uint]*user.User{
			1: {ID: 1, Role: "admin"},
			2: {ID: 2, Role: "asesor"},
			3: {ID: 3, Role: "asesor"},
			4: {ID: 4, Role: "siswa"},
		},
	}
	svc := class.NewService(classRepo, userRepo)

	// Asesor 2 membuat kelas
	c2, err := svc.CreateClass("Kelas Asesor 2", 2)
	if err != nil {
		t.Fatalf("failed to create class: %v", err)
	}

	// 1. Asesor 3 mencoba melihat kelas Asesor 2 -> Harus ErrForbidden
	_, err = svc.GetClassByID(c2.ID, "asesor", 3)
	if !errors.Is(err, class.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor 3 gets asesor 2's class, got %v", err)
	}

	// 2. Admin boleh melihat kelas Asesor 2
	adminGet, err := svc.GetClassByID(c2.ID, "admin", 1)
	if err != nil || adminGet == nil {
		t.Errorf("expected admin to view class, got err: %v", err)
	}

	// 3. Asesor 3 mencoba mengupdate kelas Asesor 2 -> Harus ErrForbidden
	_, err = svc.UpdateClass(c2.ID, "Hack Name", "asesor", 3)
	if !errors.Is(err, class.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor 3 updates asesor 2's class, got %v", err)
	}

	// 4. Asesor 3 mencoba menghapus kelas Asesor 2 -> Harus ErrForbidden
	err = svc.DeleteClass(c2.ID, "asesor", 3)
	if !errors.Is(err, class.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor 3 deletes asesor 2's class, got %v", err)
	}

	// 5. Asesor 3 mencoba menambah anggota ke kelas Asesor 2 -> Harus ErrForbidden
	err = svc.AddMember(c2.ID, 4, "asesor", 3)
	if !errors.Is(err, class.ErrForbidden) {
		t.Errorf("expected ErrForbidden when asesor 3 adds member to asesor 2's class, got %v", err)
	}

	// 6. Asesor 2 boleh menambah anggota siswa (ID 4)
	err = svc.AddMember(c2.ID, 4, "asesor", 2)
	if err != nil {
		t.Errorf("expected asesor 2 to succeed adding member, got %v", err)
	}

	// 7. Mencoba menambah non-siswa (Asesor 3) sebagai anggota -> Harus ErrInvalidMemberRole
	err = svc.AddMember(c2.ID, 3, "asesor", 2)
	if !errors.Is(err, class.ErrInvalidMemberRole) {
		t.Errorf("expected ErrInvalidMemberRole when adding assessor as member, got %v", err)
	}
}

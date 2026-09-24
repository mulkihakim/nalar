package session

import (
	"errors"
	"testing"

	"github.com/mulkihakim/nalar/backend/internal/exam"
	"github.com/mulkihakim/nalar/backend/internal/material"
)

type mockRepository struct {
	Repository
	activeSession      *Session
	sessionByID        map[uint]*Session
	peerCount          int64
	studentAttempts    map[uint]int64
	groupTotalAttempts map[uint]int64
	groupUniquePeers   map[uint]int64
	argumentWithOptions map[uint]*material.Argument
	argumentsByMatID   map[uint][]material.Argument
	exams              map[uint]*exam.Exam
	studentExams       []exam.Exam
	createdSessions    []*Session
	updatedModes       map[uint]SessionMode
}

func (m *mockRepository) EnsureIndices() error { return nil }
func (m *mockRepository) GetActiveSession(examID, studentID uint) (*Session, error) {
	return m.activeSession, nil
}
func (m *mockRepository) GetSessionByID(sessionID uint) (*Session, error) {
	if s, ok := m.sessionByID[sessionID]; ok {
		return s, nil
	}
	return nil, ErrSessionNotFound
}
func (m *mockRepository) CountUniquePeersWithSessions(examID uint) (int64, error) {
	return m.peerCount, nil
}
func (m *mockRepository) CountStudentAttemptsByOption(examID, studentID, optionID uint) (int64, error) {
	return m.studentAttempts[optionID], nil
}
func (m *mockRepository) GetGroupAttemptsByOption(examID, optionID uint) (int64, int64, error) {
	return m.groupTotalAttempts[optionID], m.groupUniquePeers[optionID], nil
}
func (m *mockRepository) GetArgumentWithOptions(argumentID uint) (*material.Argument, error) {
	if a, ok := m.argumentWithOptions[argumentID]; ok {
		return a, nil
	}
	return nil, errors.New("not found")
}
func (m *mockRepository) GetArgumentsByMaterialID(materialID uint) ([]material.Argument, error) {
	return m.argumentsByMatID[materialID], nil
}
func (m *mockRepository) GetStudentsChoosingOption(examID, optionID uint) ([]StudentChooserInfo, error) {
	return nil, nil
}
func (m *mockRepository) GetExamByID(examID uint) (*exam.Exam, error) {
	if e, ok := m.exams[examID]; ok {
		return e, nil
	}
	return nil, ErrExamNotFound
}
func (m *mockRepository) GetStudentExams(studentID uint) ([]exam.Exam, error) {
	return m.studentExams, nil
}
func (m *mockRepository) GetMaxAttemptNo(examID, studentID uint) (int, error) {
	return 0, nil
}
func (m *mockRepository) CreateSession(session *Session, argumentIDs []uint) error {
	session.ID = 101
	m.createdSessions = append(m.createdSessions, session)
	return nil
}
func (m *mockRepository) UpdateSessionMode(sessionID uint, mode SessionMode) error {
	if m.updatedModes == nil {
		m.updatedModes = make(map[uint]SessionMode)
	}
	m.updatedModes[sessionID] = mode
	return nil
}

func TestRatioCalculation_Monitoring_Criteria7(t *testing.T) {
	// Kriteria #7: Untuk data 2, 3, 4 percobaan oleh 3 siswa unik:
	// total = 9, unik = 3 -> Y = round(9/3) = 3
	// Masing-masing siswa melihat X:Y -> 2:3, 3:3, 4:3
	optID := uint(10)
	arg := &material.Argument{
		ID: 1,
		Options: []material.Option{
			{ID: optID, Text: "Option 1", Type: "ground"},
		},
	}

	testCases := []struct {
		studentID       uint
		studentAttempts int64
		expectedX       int64
		expectedY       string
	}{
		{studentID: 1, studentAttempts: 2, expectedX: 2, expectedY: "3"},
		{studentID: 2, studentAttempts: 3, expectedX: 3, expectedY: "3"},
		{studentID: 3, studentAttempts: 4, expectedX: 4, expectedY: "3"},
	}

	for _, tc := range testCases {
		mockRepo := &mockRepository{
			sessionByID: map[uint]*Session{
				1: {ID: 1, ExamID: 1, StudentID: tc.studentID},
			},
			peerCount: 3, // >= MIN_PEERS
			studentAttempts: map[uint]int64{
				optID: tc.studentAttempts,
			},
			groupTotalAttempts: map[uint]int64{
				optID: 9,
			},
			groupUniquePeers: map[uint]int64{
				optID: 3,
			},
			argumentWithOptions: map[uint]*material.Argument{
				1: arg,
			},
		}

		svc := NewService(mockRepo)
		res, err := svc.GetMonitoringAnalytics(1, 1, tc.studentID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(res.Options) != 1 {
			t.Fatalf("expected 1 option, got %d", len(res.Options))
		}

		optStat := res.Options[0]
		if optStat.X != tc.expectedX {
			t.Errorf("expected X=%d, got %d", tc.expectedX, optStat.X)
		}
		if optStat.Y != tc.expectedY {
			t.Errorf("expected Y=%s, got %s", tc.expectedY, optStat.Y)
		}
	}
}

func TestRatioCalculation_Analysis_Criteria8(t *testing.T) {
	// Kriteria #8: Untuk total 9 percobaan oleh 4 siswa unik: analysis menampilkan 9:4
	optID := uint(10)
	arg := material.Argument{
		ID: 1,
		Options: []material.Option{
			{ID: optID, Text: "Option 1", Type: "ground"},
		},
	}

	mockRepo := &mockRepository{
		sessionByID: map[uint]*Session{
			1: {
				ID:        1,
				ExamID:    1,
				StudentID: 1,
				Arguments: []SessionArgument{
					{SessionID: 1, ArgumentID: 1, Argument: &arg},
				},
			},
		},
		peerCount: 4, // >= MIN_PEERS
		groupTotalAttempts: map[uint]int64{
			optID: 9,
		},
		groupUniquePeers: map[uint]int64{
			optID: 4,
		},
	}

	svc := NewService(mockRepo)
	res, err := svc.GetAnalysisAnalytics(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Arguments) != 1 || len(res.Arguments[0].Options) != 1 {
		t.Fatalf("unexpected result structure")
	}

	optStat := res.Arguments[0].Options[0]
	if optStat.Ratio != "9:4" {
		t.Errorf("expected ratio '9:4', got '%s'", optStat.Ratio)
	}
	if optStat.TotalAttempts != 9 {
		t.Errorf("expected TotalAttempts=9, got %d", optStat.TotalAttempts)
	}
	if optStat.UniqueStudents != 4 {
		t.Errorf("expected UniqueStudents=4, got %d", optStat.UniqueStudents)
	}
}

func TestNoSiblings_ZeroY_Criteria9(t *testing.T) {
	// Kriteria #9: Opsi yang belum dipilih siapa pun menampilkan Y = - tanpa error
	optID := uint(20)
	arg := &material.Argument{
		ID: 1,
		Options: []material.Option{
			{ID: optID, Text: "Unchosen Option", Type: "ground"},
		},
	}

	mockRepo := &mockRepository{
		sessionByID: map[uint]*Session{
			1: {ID: 1, ExamID: 1, StudentID: 1},
		},
		peerCount: 3,
		studentAttempts: map[uint]int64{
			optID: 0,
		},
		groupTotalAttempts: map[uint]int64{
			optID: 0,
		},
		groupUniquePeers: map[uint]int64{
			optID: 0,
		},
		argumentWithOptions: map[uint]*material.Argument{
			1: arg,
		},
	}

	svc := NewService(mockRepo)
	res, err := svc.GetMonitoringAnalytics(1, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Options[0].Y != "-" {
		t.Errorf("expected Y='-', got '%s'", res.Options[0].Y)
	}
}

func TestMinPeers_HidesGroupData_Criteria10(t *testing.T) {
	// Kriteria #10: Data kelompok tidak tampil ke siswa jika jumlah siswa pada ujian < MIN_PEERS
	optID := uint(30)
	arg := &material.Argument{
		ID: 1,
		Options: []material.Option{
			{ID: optID, Text: "Option", Type: "ground"},
		},
	}

	mockRepo := &mockRepository{
		sessionByID: map[uint]*Session{
			1: {ID: 1, ExamID: 1, StudentID: 1},
		},
		peerCount: 2, // < MIN_PEERS (3)
		studentAttempts: map[uint]int64{
			optID: 2,
		},
		groupTotalAttempts: map[uint]int64{
			optID: 5,
		},
		groupUniquePeers: map[uint]int64{
			optID: 2,
		},
		argumentWithOptions: map[uint]*material.Argument{
			1: arg,
		},
	}

	svc := NewService(mockRepo)
	res, err := svc.GetMonitoringAnalytics(1, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.HasEnoughPeers {
		t.Errorf("expected HasEnoughPeers=false when peers < 3")
	}
	if res.Options[0].Y != "-" {
		t.Errorf("expected group data Y to be hidden ('-'), got '%s'", res.Options[0].Y)
	}
}

func TestModeConflict_Criteria13(t *testing.T) {
	// Kriteria #13: Siswa dengan sesi in_progress yang memilih mode berbeda dan belum mengonfirmasi
	// mendapat mode conflict (409); setelah konfirmasi, mode berubah.
	activeSession := &Session{
		ID:        1,
		ExamID:    1,
		StudentID: 5,
		Mode:      ModeStandard,
		Status:    StatusInProgress,
	}

	mockRepo := &mockRepository{
		activeSession: activeSession,
		sessionByID: map[uint]*Session{
			1: activeSession,
		},
		exams: map[uint]*exam.Exam{
			1: {ID: 1, IsActive: true},
		},
		studentExams: []exam.Exam{
			{ID: 1, IsActive: true},
		},
	}

	svc := NewService(mockRepo)

	// Percobaan 1: Pilih mode bantuan tanpa konfirmasi
	_, err := svc.StartOrContinueSession(1, 5, StartSessionInput{
		Mode:              ModeHelp,
		ConfirmModeChange: false,
	})
	if err == nil {
		t.Fatalf("expected mode conflict error, got nil")
	}

	var conflictErr *ModeConflictError
	if !errors.As(err, &conflictErr) {
		t.Fatalf("expected ModeConflictError, got %T: %v", err, err)
	}
	if conflictErr.CurrentMode != ModeStandard {
		t.Errorf("expected current mode 'standard', got '%s'", conflictErr.CurrentMode)
	}

	// Percobaan 2: Pilih mode bantuan DENGAN konfirmasi
	res, err := svc.StartOrContinueSession(1, 5, StartSessionInput{
		Mode:              ModeHelp,
		ConfirmModeChange: true,
	})
	if err != nil {
		t.Fatalf("unexpected error when confirmed: %v", err)
	}
	if res.Mode != ModeHelp {
		t.Errorf("expected session mode updated to 'help', got '%s'", res.Mode)
	}
	if mockRepo.updatedModes[1] != ModeHelp {
		t.Errorf("expected repository to update mode to 'help'")
	}
}

func TestCompletedArgument_ExposesIsCorrect(t *testing.T) {
	mockRepo := &mockRepository{
		sessionByID: map[uint]*Session{
			10: {
				ID:        10,
				ExamID:    1,
				StudentID: 5,
				Mode:      ModeSocial,
				Status:    StatusInProgress,
				Arguments: []SessionArgument{
					{
						SessionID:  10,
						ArgumentID: 101,
						OrderNo:    1,
						Argument: &material.Argument{
							ID:        101,
							ClaimText: "Claim 1",
							Options: []material.Option{
								{ID: 1, ArgumentID: 101, Type: "ground", Text: "G1 Correct", IsCorrect: true},
								{ID: 2, ArgumentID: 101, Type: "ground", Text: "G2 Wrong", IsCorrect: false},
								{ID: 3, ArgumentID: 101, Type: "warrant", Text: "W1 Correct", IsCorrect: true},
								{ID: 4, ArgumentID: 101, Type: "warrant", Text: "W2 Wrong", IsCorrect: false},
							},
						},
					},
					{
						SessionID:  10,
						ArgumentID: 102,
						OrderNo:    2,
						Argument: &material.Argument{
							ID:        102,
							ClaimText: "Claim 2",
							Options: []material.Option{
								{ID: 5, ArgumentID: 102, Type: "ground", Text: "G3 Correct", IsCorrect: true},
								{ID: 6, ArgumentID: 102, Type: "ground", Text: "G4 Wrong", IsCorrect: false},
								{ID: 7, ArgumentID: 102, Type: "warrant", Text: "W3 Correct", IsCorrect: true},
								{ID: 8, ArgumentID: 102, Type: "warrant", Text: "W4 Wrong", IsCorrect: false},
							},
						},
					},
				},
				Progress: []ArgumentProgress{
					{SessionID: 10, ArgumentID: 101}, // Argumen 101 sudah selesai
				},
			},
		},
	}

	svc := NewService(mockRepo)
	detail, err := svc.GetSessionDetails(10, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(detail.Arguments) != 2 {
		t.Fatalf("expected 2 arguments, got %d", len(detail.Arguments))
	}

	// Argumen 101 (selesai): harus memiliki is_correct pada opsinya
	arg101 := detail.Arguments[0]
	if !arg101.Completed {
		t.Errorf("expected argument 101 to be completed")
	}
	for _, opt := range arg101.Options {
		if opt.IsCorrect == nil {
			t.Errorf("expected option %d on completed argument to have is_correct populated", opt.ID)
		}
	}

	// Argumen 102 (belum selesai): TIDAK boleh memiliki is_correct (harus nil agar siswa tidak curang)
	arg102 := detail.Arguments[1]
	if arg102.Completed {
		t.Errorf("expected argument 102 to be uncompleted")
	}
	for _, opt := range arg102.Options {
		if opt.IsCorrect != nil {
			t.Errorf("expected option %d on uncompleted argument to have is_correct=nil, got %v", opt.ID, *opt.IsCorrect)
		}
	}
}

package material

import (
	"errors"
	"strings"
)

var (
	ErrMaterialNotFound    = errors.New("materi tidak ditemukan")
	ErrArgumentNotFound    = errors.New("argumen tidak ditemukan")
	ErrForbidden           = errors.New("akses ditolak: bukan pemilik materi")
	ErrEmptyMaterialData   = errors.New("judul dan konten materi wajib diisi")
	ErrEmptyClaimText      = errors.New("teks claim argumen wajib diisi")
	ErrInvalidOptionsCount = errors.New("argumen harus memiliki minimal 3 opsi ground (1 benar) dan 3 opsi warrant (1 benar)")
	ErrEmptyOptionText     = errors.New("setiap opsi harus memiliki teks yang tidak kosong")
)

type CreateMaterialInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type OptionInput struct {
	Type      string `json:"type"` // "ground" | "warrant"
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type ArgumentInput struct {
	ClaimText string        `json:"claim_text"`
	OrderNo   int           `json:"order_no"`
	Options   []OptionInput `json:"options"`
}

type Service interface {
	CreateMaterial(input CreateMaterialInput, ownerID uint) (*Material, error)
	GetMaterialByID(id uint, requesterRole string, requesterID uint) (*Material, error)
	ListMaterials(requesterRole string, requesterID uint) ([]Material, error)
	UpdateMaterial(id uint, input CreateMaterialInput, requesterRole string, requesterID uint) (*Material, error)
	DeleteMaterial(id uint, requesterRole string, requesterID uint) error

	CreateArgument(materialID uint, input ArgumentInput, requesterRole string, requesterID uint) (*Argument, error)
	UpdateArgument(materialID uint, argumentID uint, input ArgumentInput, requesterRole string, requesterID uint) (*Argument, error)
	DeleteArgument(materialID uint, argumentID uint, requesterRole string, requesterID uint) error
	ListArguments(materialID uint, requesterRole string, requesterID uint) ([]Argument, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateMaterial(input CreateMaterialInput, ownerID uint) (*Material, error) {
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Content) == "" {
		return nil, ErrEmptyMaterialData
	}
	m := &Material{
		Title:   strings.TrimSpace(input.Title),
		Content: strings.TrimSpace(input.Content),
		OwnerID: ownerID,
	}
	if err := s.repo.CreateMaterial(m); err != nil {
		return nil, err
	}
	return s.repo.FindMaterialByID(m.ID)
}

func (s *service) GetMaterialByID(id uint, requesterRole string, requesterID uint) (*Material, error) {
	m, err := s.repo.FindMaterialByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, ErrMaterialNotFound
	}
	if requesterRole == "asesor" && m.OwnerID != requesterID {
		return nil, ErrForbidden
	}
	return m, nil
}

func (s *service) ListMaterials(requesterRole string, requesterID uint) ([]Material, error) {
	var ownerIDFilter *uint
	if requesterRole == "asesor" {
		ownerIDFilter = &requesterID
	}
	return s.repo.ListMaterials(ownerIDFilter)
}

func (s *service) UpdateMaterial(id uint, input CreateMaterialInput, requesterRole string, requesterID uint) (*Material, error) {
	if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Content) == "" {
		return nil, ErrEmptyMaterialData
	}
	m, err := s.repo.FindMaterialByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, ErrMaterialNotFound
	}
	if requesterRole == "asesor" && m.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	m.Title = strings.TrimSpace(input.Title)
	m.Content = strings.TrimSpace(input.Content)
	if err := s.repo.UpdateMaterial(m); err != nil {
		return nil, err
	}
	return s.repo.FindMaterialByID(m.ID)
}

func (s *service) DeleteMaterial(id uint, requesterRole string, requesterID uint) error {
	m, err := s.repo.FindMaterialByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return ErrMaterialNotFound
	}
	if requesterRole == "asesor" && m.OwnerID != requesterID {
		return ErrForbidden
	}
	return s.repo.DeleteMaterial(id)
}

func validateOptions(options []OptionInput) ([]Option, error) {
	var groundCount, warrantCount int
	var groundCorrect, warrantCorrect int
	var parsedOptions []Option

	for _, opt := range options {
		text := strings.TrimSpace(opt.Text)
		if text == "" {
			return nil, ErrEmptyOptionText
		}
		optType := strings.ToLower(strings.TrimSpace(opt.Type))
		switch optType {
		case "ground":
			groundCount++
			if opt.IsCorrect {
				groundCorrect++
			}
		case "warrant":
			warrantCount++
			if opt.IsCorrect {
				warrantCorrect++
			}
		default:
			return nil, ErrInvalidOptionsCount
		}

		parsedOptions = append(parsedOptions, Option{
			Type:      optType,
			Text:      text,
			IsCorrect: opt.IsCorrect,
		})
	}

	// Minimal 3 ground dan minimal 3 warrant, masing-masing tepat 1 is_correct = true
	if groundCount < 3 || warrantCount < 3 || groundCorrect != 1 || warrantCorrect != 1 {
		return nil, ErrInvalidOptionsCount
	}

	return parsedOptions, nil
}

func (s *service) CreateArgument(materialID uint, input ArgumentInput, requesterRole string, requesterID uint) (*Argument, error) {
	m, err := s.repo.FindMaterialByID(materialID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, ErrMaterialNotFound
	}
	if requesterRole == "asesor" && m.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	claimText := strings.TrimSpace(input.ClaimText)
	if claimText == "" {
		return nil, ErrEmptyClaimText
	}

	parsedOptions, err := validateOptions(input.Options)
	if err != nil {
		return nil, err
	}

	orderNo := input.OrderNo
	if orderNo <= 0 {
		orderNo = len(m.Arguments) + 1
	}

	arg := &Argument{
		MaterialID: materialID,
		ClaimText:  claimText,
		OrderNo:    orderNo,
	}

	if err := s.repo.CreateArgument(arg, parsedOptions); err != nil {
		return nil, err
	}
	return s.repo.FindArgumentByID(arg.ID)
}

func (s *service) UpdateArgument(materialID uint, argumentID uint, input ArgumentInput, requesterRole string, requesterID uint) (*Argument, error) {
	m, err := s.repo.FindMaterialByID(materialID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, ErrMaterialNotFound
	}
	if requesterRole == "asesor" && m.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	arg, err := s.repo.FindArgumentByID(argumentID)
	if err != nil {
		return nil, err
	}
	if arg == nil || arg.MaterialID != materialID {
		return nil, ErrArgumentNotFound
	}

	claimText := strings.TrimSpace(input.ClaimText)
	if claimText == "" {
		return nil, ErrEmptyClaimText
	}

	parsedOptions, err := validateOptions(input.Options)
	if err != nil {
		return nil, err
	}

	arg.ClaimText = claimText
	if input.OrderNo > 0 {
		arg.OrderNo = input.OrderNo
	}

	if err := s.repo.UpdateArgument(arg, parsedOptions); err != nil {
		return nil, err
	}
	return s.repo.FindArgumentByID(arg.ID)
}

func (s *service) DeleteArgument(materialID uint, argumentID uint, requesterRole string, requesterID uint) error {
	m, err := s.repo.FindMaterialByID(materialID)
	if err != nil {
		return err
	}
	if m == nil {
		return ErrMaterialNotFound
	}
	if requesterRole == "asesor" && m.OwnerID != requesterID {
		return ErrForbidden
	}

	arg, err := s.repo.FindArgumentByID(argumentID)
	if err != nil {
		return err
	}
	if arg == nil || arg.MaterialID != materialID {
		return ErrArgumentNotFound
	}

	return s.repo.DeleteArgument(argumentID)
}

func (s *service) ListArguments(materialID uint, requesterRole string, requesterID uint) ([]Argument, error) {
	m, err := s.repo.FindMaterialByID(materialID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, ErrMaterialNotFound
	}
	if requesterRole == "asesor" && m.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	return s.repo.ListArguments(materialID)
}

package material

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	CreateMaterial(m *Material) error
	FindMaterialByID(id uint) (*Material, error)
	ListMaterials(ownerID *uint) ([]Material, error)
	UpdateMaterial(m *Material) error
	DeleteMaterial(id uint) error

	CreateArgument(arg *Argument, options []Option) error
	FindArgumentByID(id uint) (*Argument, error)
	ListArguments(materialID uint) ([]Argument, error)
	UpdateArgument(arg *Argument, options []Option) error
	DeleteArgument(id uint) error
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateMaterial(m *Material) error {
	return r.db.Create(m).Error
}

func (r *gormRepository) FindMaterialByID(id uint) (*Material, error) {
	var m Material
	err := r.db.Preload("Owner").
		Preload("Arguments", func(db *gorm.DB) *gorm.DB {
			return db.Order("arguments.order_no ASC, arguments.id ASC")
		}).
		Preload("Arguments.Options").
		First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *gormRepository) ListMaterials(ownerID *uint) ([]Material, error) {
	var materials []Material
	q := r.db.Preload("Owner").
		Preload("Arguments", func(db *gorm.DB) *gorm.DB {
			return db.Order("arguments.order_no ASC, arguments.id ASC")
		}).
		Preload("Arguments.Options")
	if ownerID != nil {
		q = q.Where("owner_id = ?", *ownerID)
	}
	err := q.Order("id DESC").Find(&materials).Error
	return materials, err
}

func (r *gormRepository) UpdateMaterial(m *Material) error {
	return r.db.Save(m).Error
}

func (r *gormRepository) DeleteMaterial(id uint) error {
	return r.db.Delete(&Material{}, id).Error
}

func (r *gormRepository) CreateArgument(arg *Argument, options []Option) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(arg).Error; err != nil {
			return err
		}
		for i := range options {
			options[i].ArgumentID = arg.ID
			if err := tx.Create(&options[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *gormRepository) FindArgumentByID(id uint) (*Argument, error) {
	var arg Argument
	err := r.db.Preload("Options").First(&arg, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &arg, nil
}

func (r *gormRepository) ListArguments(materialID uint) ([]Argument, error) {
	var args []Argument
	err := r.db.Where("material_id = ?", materialID).
		Preload("Options").
		Order("order_no ASC, id ASC").
		Find(&args).Error
	return args, err
}

func (r *gormRepository) UpdateArgument(arg *Argument, options []Option) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(arg).Error; err != nil {
			return err
		}
		// Hapus opsi lama, simpan opsi baru
		if err := tx.Where("argument_id = ?", arg.ID).Delete(&Option{}).Error; err != nil {
			return err
		}
		for i := range options {
			options[i].ArgumentID = arg.ID
			if err := tx.Create(&options[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *gormRepository) DeleteArgument(id uint) error {
	return r.db.Delete(&Argument{}, id).Error
}

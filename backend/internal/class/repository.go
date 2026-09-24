package class

import (
	"errors"

	"github.com/mulkihakim/nalar/backend/internal/user"
	"gorm.io/gorm"
)

type Repository interface {
	Create(c *Class) error
	FindByID(id uint) (*Class, error)
	List(ownerID *uint) ([]Class, error)
	Update(c *Class) error
	Delete(id uint) error
	AddMember(classID uint, userID uint) error
	RemoveMember(classID uint, userID uint) error
	GetMembers(classID uint) ([]user.User, error)
	IsMember(classID uint, userID uint) (bool, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(c *Class) error {
	return r.db.Create(c).Error
}

func (r *gormRepository) FindByID(id uint) (*Class, error) {
	var c Class
	err := r.db.Preload("Owner").Preload("Members").First(&c, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *gormRepository) List(ownerID *uint) ([]Class, error) {
	var classes []Class
	q := r.db.Preload("Owner").Preload("Members")
	if ownerID != nil {
		q = q.Where("owner_id = ?", *ownerID)
	}
	err := q.Find(&classes).Error
	return classes, err
}

func (r *gormRepository) Update(c *Class) error {
	return r.db.Save(c).Error
}

func (r *gormRepository) Delete(id uint) error {
	return r.db.Delete(&Class{}, id).Error
}

func (r *gormRepository) AddMember(classID uint, userID uint) error {
	member := ClassMember{
		ClassID: classID,
		UserID:  userID,
	}
	return r.db.Create(&member).Error
}

func (r *gormRepository) RemoveMember(classID uint, userID uint) error {
	return r.db.Where("class_id = ? AND user_id = ?", classID, userID).Delete(&ClassMember{}).Error
}

func (r *gormRepository) GetMembers(classID uint) ([]user.User, error) {
	var members []user.User
	err := r.db.Table("users").
		Joins("JOIN class_members ON class_members.user_id = users.id").
		Where("class_members.class_id = ?", classID).
		Find(&members).Error
	return members, err
}

func (r *gormRepository) IsMember(classID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&ClassMember{}).
		Where("class_id = ? AND user_id = ?", classID, userID).
		Count(&count).Error
	return count > 0, err
}

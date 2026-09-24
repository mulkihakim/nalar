package user

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(u *User) error
	FindByID(id uint) (*User, error)
	FindByUsername(username string) (*User, error)
	List() ([]User, error)
	Update(u *User) error
	Delete(id uint) error
	HasExamSessions(userID uint) (bool, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(u *User) error {
	return r.db.Create(u).Error
}

func (r *gormRepository) FindByID(id uint) (*User, error) {
	var u User
	err := r.db.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *gormRepository) FindByUsername(username string) (*User, error) {
	var u User
	err := r.db.Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *gormRepository) List() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *gormRepository) Update(u *User) error {
	return r.db.Save(u).Error
}

func (r *gormRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

func (r *gormRepository) HasExamSessions(userID uint) (bool, error) {
	if !r.db.Migrator().HasTable("sessions") {
		return false, nil
	}
	var count int64
	err := r.db.Table("sessions").Where("student_id = ?", userID).Count(&count).Error
	return count > 0, err
}
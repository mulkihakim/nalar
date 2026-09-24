package exam

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(e *Exam, classIDs []uint, studentIDs []uint) error
	FindByID(id uint) (*Exam, error)
	List(ownerID *uint) ([]Exam, error)
	Update(e *Exam) error
	Delete(id uint) error
	SetAccess(examID uint, classIDs []uint, studentIDs []uint) error
	SetStatus(examID uint, isActive bool) error
	HasStudentTakenExam(examID uint, studentID uint) (bool, error)
	HasClassTakenExam(examID uint, classID uint) (bool, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(e *Exam, classIDs []uint, studentIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(e).Error; err != nil {
			return err
		}
		for _, cid := range classIDs {
			ec := ExamClass{ExamID: e.ID, ClassID: cid}
			if err := tx.Create(&ec).Error; err != nil {
				return err
			}
		}
		for _, sid := range studentIDs {
			es := ExamStudent{ExamID: e.ID, StudentID: sid}
			if err := tx.Create(&es).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *gormRepository) FindByID(id uint) (*Exam, error) {
	var e Exam
	err := r.db.Preload("Material").
		Preload("Owner").
		Preload("Classes").
		Preload("Students").
		First(&e, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *gormRepository) List(ownerID *uint) ([]Exam, error) {
	var exams []Exam
	q := r.db.Preload("Material").
		Preload("Owner").
		Preload("Classes").
		Preload("Students")
	if ownerID != nil {
		q = q.Where("owner_id = ?", *ownerID)
	}
	err := q.Order("id DESC").Find(&exams).Error
	return exams, err
}

func (r *gormRepository) Update(e *Exam) error {
	return r.db.Save(e).Error
}

func (r *gormRepository) Delete(id uint) error {
	return r.db.Delete(&Exam{}, id).Error
}

func (r *gormRepository) SetAccess(examID uint, classIDs []uint, studentIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Hapus relasi akses lama
		if err := tx.Where("exam_id = ?", examID).Delete(&ExamClass{}).Error; err != nil {
			return err
		}
		if err := tx.Where("exam_id = ?", examID).Delete(&ExamStudent{}).Error; err != nil {
			return err
		}

		for _, cid := range classIDs {
			ec := ExamClass{ExamID: examID, ClassID: cid}
			if err := tx.Create(&ec).Error; err != nil {
				return err
			}
		}
		for _, sid := range studentIDs {
			es := ExamStudent{ExamID: examID, StudentID: sid}
			if err := tx.Create(&es).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *gormRepository) SetStatus(examID uint, isActive bool) error {
	return r.db.Model(&Exam{}).Where("id = ?", examID).Update("is_active", isActive).Error
}

func (r *gormRepository) HasStudentTakenExam(examID uint, studentID uint) (bool, error) {
	if !r.db.Migrator().HasTable("sessions") {
		return false, nil
	}
	var count int64
	err := r.db.Table("sessions").Where("exam_id = ? AND student_id = ?", examID, studentID).Count(&count).Error
	return count > 0, err
}

func (r *gormRepository) HasClassTakenExam(examID uint, classID uint) (bool, error) {
	if !r.db.Migrator().HasTable("sessions") || !r.db.Migrator().HasTable("class_members") {
		return false, nil
	}
	var count int64
	err := r.db.Table("sessions").
		Joins("JOIN class_members ON class_members.user_id = sessions.student_id").
		Where("sessions.exam_id = ? AND class_members.class_id = ?", examID, classID).
		Count(&count).Error
	return count > 0, err
}

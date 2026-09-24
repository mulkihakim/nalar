package class

import (
	"time"

	"github.com/mulkihakim/nalar/backend/internal/user"
)

type Class struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	Name      string      `gorm:"not null" json:"name"`
	OwnerID   uint        `gorm:"not null;index" json:"owner_id"`
	Owner     *user.User  `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE" json:"owner,omitempty"`
	Members   []user.User `gorm:"many2many:class_members;joinForeignKey:ClassID;joinReferences:UserID" json:"members,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type ClassMember struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	ClassID   uint       `gorm:"not null;uniqueIndex:idx_class_user;index" json:"class_id"`
	Class     *Class     `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"-"`
	UserID    uint       `gorm:"not null;uniqueIndex:idx_class_user;index" json:"user_id"`
	User      *user.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time  `json:"created_at"`
}

func (Class) TableName() string {
	return "classes"
}

func (ClassMember) TableName() string {
	return "class_members"
}


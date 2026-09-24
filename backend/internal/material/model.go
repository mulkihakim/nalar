package material

import (
	"time"

	"github.com/mulkihakim/nalar/backend/internal/user"
)

type Material struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	Title     string      `gorm:"not null" json:"title"`
	Content   string      `gorm:"not null" json:"content"`
	OwnerID   uint        `gorm:"not null;index" json:"owner_id"`
	Owner     *user.User  `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE" json:"owner,omitempty"`
	Arguments []Argument  `gorm:"foreignKey:MaterialID;constraint:OnDelete:CASCADE" json:"arguments,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type Argument struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	MaterialID uint      `gorm:"not null;index" json:"material_id"`
	Material   *Material `gorm:"foreignKey:MaterialID;constraint:OnDelete:CASCADE" json:"-"`
	ClaimText  string    `gorm:"not null" json:"claim_text"`
	OrderNo    int       `gorm:"not null;default:1" json:"order_no"`
	Options    []Option  `gorm:"foreignKey:ArgumentID;constraint:OnDelete:CASCADE" json:"options,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Option struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ArgumentID uint      `gorm:"not null;index" json:"argument_id"`
	Argument   *Argument `gorm:"foreignKey:ArgumentID;constraint:OnDelete:CASCADE" json:"-"`
	Type       string    `gorm:"not null" json:"type"` // "ground" | "warrant"
	Text       string    `gorm:"not null" json:"text"`
	IsCorrect  bool      `gorm:"not null;default:false" json:"is_correct"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Material) TableName() string {
	return "materials"
}

func (Argument) TableName() string {
	return "arguments"
}

func (Option) TableName() string {
	return "options"
}

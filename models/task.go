package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"not null" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" gorm:"default:'pending'"`

	CreatedByID uint      `json:"created_by_id" gorm:"column:created_by_id;not null"`
	CreatedBy   User      `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" validate:"-"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`

	UpdatedByID uint      `json:"updated_by_id" gorm:"column:updated_by_id;not null"`
	UpdatedBy   User      `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" validate:"-"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	DeletedByID *uint          `json:"deleted_by_id" gorm:"column:deleted_by_id"`
	DeletedBy   *User          `gorm:"foreignKey:DeletedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" validate:"-"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

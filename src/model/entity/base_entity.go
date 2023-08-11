package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseEntity struct {
	ID *uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Timestamp
}

type Timestamp struct {
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;type:timestamp;not null;autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp ON UPDATE CURRENT_TIMESTAMP;null;autoUpdateTime"`
}

func (m *BaseEntity) SetUUID() *BaseEntity {
	m.ID = new(uuid.UUID)
	*m.ID = uuid.New()
	return m
}

func (m *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	m.SetUUID()
	return nil
}

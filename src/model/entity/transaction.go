package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	BaseEntity
	Type        string          `json:"type" gorm:"type:varchar(30);index;not null"`
	ReferenceID string          `json:"reference_id" gorm:"column:reference_id;type:varchar(100);index;not null"`
	At          *time.Time      `json:"at" gorm:"type:timestamp;null"`
	Status      string          `json:"status" gorm:"column:status;type:varchar(30);index;not null;default:'submitted'"`
	Amount      decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(64,15);not null;default:0"`
	StatusAt    *time.Time      `json:"status_at" gorm:"column:status_at;type:timestamp;null"`
	IsProcessed int             `json:"is_processed" gorm:"column:is_processed;not null;default:0"`
	ProcessedAt *time.Time      `json:"processed_at" gorm:"column:processed_at;type:timestamp;null"`
	IsCompleted int             `json:"is_completed" gorm:"column:is_completed;not null;default:0"`
	CompletedAt *time.Time      `json:"completed_at" gorm:"column:completed_at;type:timestamp;null"`
	IsFailed    int             `json:"is_failed" gorm:"column:is_failed;not null;default:0"`
	FailedAt    *time.Time      `json:"failed_at" gorm:"column:failed_at;type:timestamp;null"`
}

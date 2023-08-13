package entity

import (
	"time"

	"gorm.io/gorm"

	"github.com/shopspring/decimal"
)

const (
	DEPOSIT    = "deposit"
	WITHDRAWAL = "withdrawal"
)

const (
	TransactionStatusSubmitted = "submitted"
	TransactionStatusProcessed = "processed"
	TransactionStatusCompleted = "success"
	TransactionStatusFailed    = "failed"
)

type Transaction struct {
	BaseEntity
	CustomerXID string          `json:"customer_xid" gorm:"column:customer_xid;type:char(36);index;not null"`
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

func (m *Transaction) syncStatus() {
	if m.IsFailed == 1 {
		m.Status = TransactionStatusFailed
		m.StatusAt = m.FailedAt
	} else if m.IsCompleted == 1 {
		m.Status = TransactionStatusCompleted
		m.StatusAt = m.CompletedAt
	} else if m.IsProcessed == 1 {
		m.Status = TransactionStatusProcessed
		m.StatusAt = m.ProcessedAt
	} else {
		m.Status = TransactionStatusSubmitted
		m.StatusAt = &m.CreatedAt
	}
}

func (m *Transaction) processSyncDate() {
	if m.IsProcessed == 0 && m.ProcessedAt != nil {
		m.ProcessedAt = nil
	} else if m.IsProcessed == 1 && m.ProcessedAt == nil {
		nwTime := time.Now()
		m.ProcessedAt = &nwTime
	}
}

func (m *Transaction) completeSyncDate() {
	if m.IsCompleted == 0 && m.CompletedAt != nil {
		m.CompletedAt = nil
	} else if m.IsCompleted == 1 && m.CompletedAt == nil {
		nwTime := time.Now()
		m.CompletedAt = &nwTime
	}
}

func (m *Transaction) failSyncDate() {
	if m.IsFailed == 0 && m.FailedAt != nil {
		m.FailedAt = nil
	} else if m.IsFailed == 1 && m.FailedAt == nil {
		nwTime := time.Now()
		m.FailedAt = &nwTime
	}
}

func (m *Transaction) BeforeSave(db *gorm.DB) error {
	m.processSyncDate()
	m.completeSyncDate()
	m.failSyncDate()
	m.syncStatus()
	return nil
}

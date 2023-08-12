package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	BaseEntity
	CustomerXID      string          `json:"customer_xid" gorm:"column:customer_xid;type:char(36);unique;index;not null"`
	IsActive         int             `json:"is_active" gorm:"column:is_active;not null;default:0"`
	ActivatedAt      *time.Time      `json:"activated_at" gorm:"column:activated_at;type:timestamp;null"`
	BalanceAvailable decimal.Decimal `json:"balance_available" gorm:"column:balance_available;type:decimal(64,15);not null;default:0"`
	BalancePending   decimal.Decimal `json:"balance_pending" gorm:"column:balance_pending;type:decimal(64,15);not null;default:0"`
}

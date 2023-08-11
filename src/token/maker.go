package token

import (
	"time"
)

type Maker interface {
	CreateToken(customerXID string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

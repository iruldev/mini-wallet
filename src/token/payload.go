package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Payload struct {
	ID          uuid.UUID `json:"id"`
	CustomerXID string    `json:"customer_xid"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt.Add(-time.Second)), nil
}

func (p *Payload) GetIssuer() (string, error) {
	return p.CustomerXID, nil
}

func (p *Payload) GetSubject() (string, error) {
	return p.ID.String(), nil
}

func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	audience := jwt.ClaimStrings{}
	audience = append(audience, p.CustomerXID)
	return audience, nil
}

func NewPayload(customerXID string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:          tokenID,
		CustomerXID: customerXID,
		IssuedAt:    time.Now(),
		ExpiredAt:   time.Now().Add(duration),
	}
	return payload, nil
}

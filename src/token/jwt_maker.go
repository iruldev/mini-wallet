package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) Maker {
	return &JWTMaker{secretKey: secretKey}
}

func (t *JWTMaker) CreateToken(customerXID string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(customerXID, duration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(t.secretKey))
	return token, payload, err
}

func (t *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token is invalid")
		}
		return []byte(t.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, &Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errors.New("token is invalid")
	}

	return payload, nil
}

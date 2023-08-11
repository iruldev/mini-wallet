package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker := NewJWTMaker("julojayajayajaya")

	customerXID := "ea0212d3-abd6-406f-8c67-868e814a2436"
	duration := time.Minute * 5

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(customerXID, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, customerXID, payload.CustomerXID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	ISSUER = "hopper.bdreece.dev"
)

func CreateToken(sub, aud, secret string) (*string, *time.Time, error) {
	now := time.Now()
	expiration := now.Add(2 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    ISSUER,
		Subject:   sub,
		Audience:  jwt.ClaimStrings{aud},
		ExpiresAt: jwt.NewNumericDate(expiration),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        uuid.New().String(),
	})

	jwt, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, nil, err
	}

	return &jwt, &expiration, err
}

func DecodeApiKey(apiKey, secret string) (id, tenantId *uint64, err error) {
	token, err := jwt.Parse(apiKey, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Bad signing method")
		}
		return secret, nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.RegisteredClaims)
	if !ok || !token.Valid {
		err = errors.New("Invalid API key")
		return
	}

	*id, err = strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		return
	}

	aud := claims.Audience[0]
	*tenantId, err = strconv.ParseUint(aud, 10, 32)
	return
}

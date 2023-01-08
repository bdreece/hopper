/*
 * hopper - A gRPC API for collecting IoT device event messages
 * Copyright (C) 2022 Brian Reece

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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

var (
	ErrBadSigningMethod = errors.New("bad signing method")
	ErrInvalidApiKey    = errors.New("invalid API key")
	ErrInvalidToken     = errors.New("invalid token")
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

func DecodeToken(token, secret string) (sub, aud, iss *string, err error) {
	t, err := jwt.Parse(token, func(u *jwt.Token) (interface{}, error) {
		if _, ok := u.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrBadSigningMethod
		}

		return secret, nil
	})

	if err != nil {
		return nil, nil, nil, WrapError(ErrInvalidToken, err)
	}

	claims, ok := t.Claims.(jwt.RegisteredClaims)
	if !ok || !t.Valid {
		return nil, nil, nil, ErrInvalidToken
	}

	return &claims.Subject, &claims.Audience[0], &claims.Issuer, nil
}

func DecodeApiKey(apiKey, secret string) (id, tenantId *uint64, err error) {
	token, err := jwt.Parse(apiKey, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrBadSigningMethod
		}
		return secret, nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.RegisteredClaims)
	if !ok || !token.Valid {
		err = ErrInvalidApiKey
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

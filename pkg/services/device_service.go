package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	. "github.com/bdreece/hopper/pkg/models"
	pb "github.com/bdreece/hopper/pkg/proto"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const (
	ISSUER = "hopper.bdreece.dev"
)

type DeviceService struct {
	db     *gorm.DB
	secret string
}

func (s *DeviceService) AuthDevice(ctx context.Context, input *pb.AuthDeviceRequest) (*pb.AuthDeviceResponse, error) {
	var device *Device = nil

	s.db.Where("uuid = ?", input.Uuid).First(device)

	if device == nil {
		salt := make([]byte, 128)
		rand.Read(salt)

		saltedKey := append([]byte(input.ApiKey), salt...)
		hash := sha256.Sum256(saltedKey)

		hashString := base64.StdEncoding.EncodeToString(hash[:])
		saltString := base64.StdEncoding.EncodeToString(salt[:])

		device = &Device{
			Uuid: input.Uuid,
			Hash: hashString,
			Salt: saltString,
		}

		s.db.Save(device)
	} else {
		salt, err := base64.StdEncoding.DecodeString(device.Salt)
		if err != nil {
			return nil, err
		}

		saltedKey := append([]byte(input.ApiKey), salt...)
		hash := sha256.Sum256(saltedKey)

		hashString := base64.StdEncoding.EncodeToString(hash[:])

		if hashString != device.Hash {
			return nil, errors.New("Bad credentials!")
		}
	}

	now := time.Now()
	accessTokenExpiration := now.Add(2 * time.Hour)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(accessTokenExpiration),
		Issuer:    ISSUER,
		Subject:   fmt.Sprint(device.Entity.ID),
	})

	accessJwt, err := accessToken.SignedString([]byte(s.secret))
	if err != nil {
		return nil, err
	}

	refreshTokenExpiration := time.Now().Add(24 * 7 * 365 * 10 * time.Hour)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    ISSUER,
		Subject:   fmt.Sprint(device.Entity.ID),
		Audience:  jwt.ClaimStrings{fmt.Sprint(device.TenantID)},
		ExpiresAt: jwt.NewNumericDate(refreshTokenExpiration),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        uuid.New().String(),
	})

	refreshJwt, err := refreshToken.SignedString([]byte(s.secret))
	if err != nil {
		return nil, err
	}

	return &pb.AuthDeviceResponse{
		AccessToken: &pb.AuthDeviceResponse_Token{
			Token:      accessJwt,
			Expiration: timestamppb.New(accessTokenExpiration),
		},
		RefreshToken: &pb.AuthDeviceResponse_Token{
			Token:      refreshJwt,
			Expiration: timestamppb.New(refreshTokenExpiration),
		},
	}, nil
}

func (s *DeviceService) GetDevice(ctx context.Context, input *pb.GetDeviceRequest) (*pb.Device, error) {
	device := new(Device)

	switch t := input.Where.(type) {
	case *pb.GetDeviceRequest_Id:
		s.db.Where("id = ?", uint(t.Id)).First(device)
		break
	case *pb.GetDeviceRequest_Uuid:
		s.db.Where("uuid = ?", t.Uuid).First(device)
		break
	}

	if device == nil {
		return nil, errors.New("Device not found!")
	}

	return device.Marshal(), nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, input *pb.UpdateDeviceRequest) (*pb.Device, error) {
	device := new(Device)

	switch t := input.Where.Where.(type) {
	case *pb.GetDeviceRequest_Id:
		s.db.Where("id = ?", uint(t.Id)).First(device)
		break
	case *pb.GetDeviceRequest_Uuid:
		s.db.Where("uuid = ?", t.Uuid).First(device)
		break
	}

	if device == nil {
		return nil, errors.New("Device not found!")
	}

	device.Update(input)
	s.db.Save(device)

	return device.Marshal(), nil

}

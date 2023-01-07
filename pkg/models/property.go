package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Property struct {
	gorm.Model
	pb.Property
	Events []Event
}

func NewProperty(req pb.CreatePropertyRequest) Property {
	return Property{
		Property: pb.Property{
			Uuid: uuid.NewString(),
		},
	}
}

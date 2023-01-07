package models

import (
	"github.com/bdreece/hopper/pkg/proto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Property struct {
	gorm.Model
	proto.Property
	Events []Event
}

func NewProperty(req *proto.CreatePropertyRequest) Property {
	return Property{
		Property: proto.Property{
			Uuid:        uuid.NewString(),
			Name:        req.Name,
			Description: req.Description,
			TypeId:      req.TypeId,
		},
	}
}

func (p *Property) Update(in *proto.UpdatePropertyRequest) {
	if in.Name != nil {
		p.Name = *in.Name
	}
	if in.Description != nil {
		p.Description = in.Description
	}
	if in.TypeId != nil {
		p.TypeId = *in.TypeId
	}
}

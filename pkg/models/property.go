package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Property struct {
	NamedEntity

	PresentValue *string

	TypeID   uint
	DeviceID *uint
	ModelID  *uint

	Events []Event
}

func (p Property) Marshal() (bytes []byte, err error) {
	msg := &pb.Property{
		Id:           uint32(p.ID),
		CreatedAt:    timestamppb.New(p.CreatedAt),
		UpdatedAt:    timestamppb.New(p.UpdatedAt),
		Name:         p.Name,
		Description:  p.Description,
		PresentValue: p.PresentValue,
		TypeId:       uint32(p.TypeID),
	}

	bytes, err = proto.Marshal(msg)
	return
}

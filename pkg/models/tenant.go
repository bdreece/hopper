package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Tenant struct {
	NamedEntity

	Uuid string
	Hash string
	Salt string

	Devices []Device
}

func (t Tenant) Marshal() (bytes []byte, err error) {
	msg := &pb.Tenant{
		Id:          uint32(t.ID),
		CreatedAt:   timestamppb.New(t.CreatedAt),
		UpdatedAt:   timestamppb.New(t.UpdatedAt),
		Name:        t.Name,
		Description: t.Description,
		Uuid:        t.Uuid,
	}

	bytes, err = proto.Marshal(msg)
	return
}

package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Model struct {
	NamedEntity

	Uuid     string
	TenantID uint

	Firmwares  []Firmware
	Devices    []Device
	Properties []Property
}

func (m Model) Marshal() (bytes []byte, err error) {
	msg := &pb.Model{
		Id:          uint32(m.ID),
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timestamppb.New(m.UpdatedAt),
		Name:        m.Name,
		Description: m.Description,
		Uuid:        m.Uuid,
		TenantId:    uint32(m.TenantID),
	}

	bytes, err = proto.Marshal(msg)
	return
}

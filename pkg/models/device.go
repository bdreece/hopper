package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Device struct {
	NamedEntity

	Uuid string
	Hash string
	Salt string

	TenantID   uint
	ModelID    uint
	FirmwareID uint

	Properties []Property
	Events     []Event
}

func (d *Device) Update(input *pb.UpdateDeviceRequest) {
	if input.Name != nil {
		d.Name = *input.Name
	}
	if input.Description != nil {
		d.Description = input.Description
	}
	if input.FirmwareId != nil {
		d.FirmwareID = uint(*input.FirmwareId)
	}
	if input.ModelId != nil {
		d.ModelID = uint(*input.ModelId)
	}
}

func (d Device) Marshal() (bytes []byte, err error) {
	msg := &pb.Device{
		Id:          uint32(d.ID),
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
		Name:        d.Name,
		Description: d.Description,
		Uuid:        d.Uuid,
		TenantId:    uint32(d.TenantID),
		ModelId:     uint32(d.ModelID),
		FirmwareId:  uint32(d.FirmwareID),
	}
	bytes, err = proto.Marshal(msg)
	return
}

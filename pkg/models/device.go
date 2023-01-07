package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	pb.Device

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
		d.FirmwareId = *input.FirmwareId
	}
	if input.ModelId != nil {
		d.ModelId = *input.ModelId
	}
}

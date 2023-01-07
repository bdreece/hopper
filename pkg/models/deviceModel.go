package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"
	"gorm.io/gorm"
)

type DeviceModel struct {
	gorm.Model
	pb.DeviceModel
	Firmwares  []Firmware
	Devices    []Device
	Properties []Property
}

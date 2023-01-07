package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"
	"gorm.io/gorm"
)

type Firmware struct {
	gorm.Model
	pb.Firmware
	Devices []Device
}

package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"
	"gorm.io/gorm"
)

type Tenant struct {
	gorm.Model
	pb.Tenant
	Devices []Device
}

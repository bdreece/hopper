package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"
	"gorm.io/gorm"
)

type Unit struct {
	gorm.Model
	pb.Unit
	Types []Type
}

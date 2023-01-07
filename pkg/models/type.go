package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"
	"gorm.io/gorm"
)

type Type struct {
	gorm.Model
	pb.Type
	Properties []Property
}

package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Firmware struct {
	Entity

	VersionMajor uint
	VersionMinor uint
	VersionPatch uint
	Url          string

	ModelID uint
	Devices []Device
}

func (f Firmware) Marshal() (bytes []byte, err error) {
	msg := &pb.Firmware{
		Id:           uint32(f.ID),
		CreatedAt:    timestamppb.New(f.CreatedAt),
		UpdatedAt:    timestamppb.New(f.UpdatedAt),
		VersionMajor: uint32(f.VersionMajor),
		VersionMinor: uint32(f.VersionMinor),
		VersionPatch: uint32(f.VersionPatch),
		Url:          f.Url,
	}

	bytes, err = proto.Marshal(msg)
	return
}

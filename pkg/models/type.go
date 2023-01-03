package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Type struct {
	NamedEntity

	DataType string
	UnitID   uint

	Properties []Property
}

func (t Type) Marshal() (bytes []byte, err error) {
	msg := &pb.Type{
		Id:          uint32(t.ID),
		CreatedAt:   timestamppb.New(t.CreatedAt),
		UpdatedAt:   timestamppb.New(t.UpdatedAt),
		Name:        t.Name,
		Description: t.Description,
		DataType:    t.DataType,
		UnitId:      uint32(t.UnitID),
	}

	bytes, err = proto.Marshal(msg)
	return
}

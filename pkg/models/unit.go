package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Unit struct {
	NamedEntity

	Symbol string

	Types []Type
}

func (u Unit) Marshal() (bytes []byte, err error) {
	msg := &pb.Unit{
		Id:          uint32(u.ID),
		CreatedAt:   timestamppb.New(u.CreatedAt),
		UpdatedAt:   timestamppb.New(u.UpdatedAt),
		Name:        u.Name,
		Description: u.Description,
		Symbol:      u.Symbol,
	}

	bytes, err = proto.Marshal(msg)
	return
}

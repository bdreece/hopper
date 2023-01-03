package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Event struct {
	Entity

	DeviceID   uint
	PropertyID uint
	Value      string
}

func (e Event) Marshal() (bytes []byte, err error) {
	msg := &pb.Event{
		Id:         uint32(e.ID),
		CreatedAt:  timestamppb.New(e.CreatedAt),
		UpdatedAt:  timestamppb.New(e.UpdatedAt),
		DeviceId:   uint32(e.DeviceID),
		PropertyId: uint32(e.PropertyID),
		Value:      e.Value,
	}

	bytes, err = proto.Marshal(msg)
	return
}

package models

import (
	pb "github.com/bdreece/hopper/pkg/proto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	pb.Event
}

func NewEvent(deviceId uint32, req *pb.CreateEventRequest) Event {
	return Event{
		Event: pb.Event{
			Uuid:       uuid.NewString(),
			Timestamp:  req.GetTimestamp(),
			Value:      req.GetValue(),
			DeviceId:   deviceId,
			PropertyId: req.GetPropertyId(),
		},
	}
}

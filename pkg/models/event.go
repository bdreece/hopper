/*
 * hopper - A gRPC API for collecting IoT device event messages
 * Copyright (C) 2022 Brian Reece

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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

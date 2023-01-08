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

package services

import (
	"context"
	"errors"

	"github.com/bdreece/hopper/pkg/config"
	"github.com/bdreece/hopper/pkg/models"
	"github.com/bdreece/hopper/pkg/proto"
	"github.com/bdreece/hopper/pkg/proto/grpc"
	"github.com/bdreece/hopper/pkg/utils"
	"github.com/bdreece/hopper/pkg/utils/iter"

	"gorm.io/gorm"
)

var (
	ErrCreateEvent     = errors.New("failed creating event")
	ErrEventNotFound   = errors.New("event not found")
	ErrEventQuery      = errors.New("failed to query event")
	ErrMissingDeviceId = errors.New("missing device ID")
)

type EventService struct {
	db     *gorm.DB
	logger utils.Logger
	grpc.UnimplementedEventServiceServer
}

func NewEventService(cfg *config.Config) *EventService {
	return &EventService{
		db:     cfg.DB,
		logger: cfg.Logger.WithContext("EventService"),

		UnimplementedEventServiceServer: grpc.UnimplementedEventServiceServer{},
	}
}

func (s *EventService) CreateEvents(ctx context.Context, in *proto.CreateEventsRequest) (*proto.Events, error) {
	logger := s.logger.WithContext("CreateEvents")
	logger.Infoln("Resolving device ID from context...")
	deviceId, ok := ctx.Value("deviceId").(uint32)
	if !ok {
		return nil, s.handleError(ErrMissingDeviceId)
	}

	logger.Infoln("Creating events...")
	events := make([]models.Event, 0, 1)
	eventMsgs := make([]*proto.Event, 0, 1)

	for _, eventRequest := range in.Events {
		event := models.NewEvent(deviceId, eventRequest)
		eventMsgs = append(eventMsgs, &event.Event)
		events = append(events, event)
	}

	logger.Infof("Saving %d events to database...", len(events))
	result := s.db.Create(&events)
	if result.Error != nil {
		result.Error = utils.WrapError(ErrCreateEvent, result.Error)
		return nil, s.handleError(result.Error)
	}

	logger.Infoln("Events saved!")
	return &proto.Events{
		Events: eventMsgs,
	}, nil
}

func (s *EventService) GetEvent(ctx context.Context, in *proto.GetEventRequest) (*proto.Event, error) {
	logger := s.logger.WithContext("GetEvent")
	logger.Infoln("Querying event...")

	query := s.db
	switch t := in.Where.(type) {
	case *proto.GetEventRequest_Uuid:
		logger.Infoln("...by UUID")
		query = query.
			Where("uuid = ?", t.Uuid)

	case *proto.GetEventRequest_DeviceTimestamp:
		logger.Infoln("...by device and timestamp")
		query = query.
			Joins("inner join device on event.deviceId = device.Id").
			Where("device.uuid = ? AND event.timestamp = ?",
				t.DeviceTimestamp.Device.Uuid,
				t.DeviceTimestamp.Timestamp)
	}

	event := models.Event{}
	result := query.First(&event)
	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Event received!")
	return &event.Event, nil
}

func (s *EventService) GetEvents(ctx context.Context, in *proto.GetEventsRequest) (*proto.Events, error) {
	logger := s.logger.WithContext("GetEvents")
	logger.Infoln("Querying events by device UUID")

	events := make([]models.Event, 0, 1)
	result := s.db.
		Joins("inner join device on event.deviceId = device.Id").
		Where("device.deviceUuid = ?", in.Where.Device.Uuid).
		Scan(&events)

	if result.Error != nil {
		return nil, s.handleQueryError(result.Error)
	}

	logger.Infoln("Events received!")
	return &proto.Events{
		Events: iter.Collect(iter.NewMap(
			iter.FromSlice(&events),
			func(in *models.Event) *proto.Event {
				return &in.Event
			},
		)),
	}, nil
}

func (s *EventService) handleError(err error) error {
	s.logger.Errorf("An error occurred: %v\n", err)
	return err
}

func (s *EventService) handleQueryError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = utils.WrapError(ErrEventNotFound, err)
	}
	err = utils.WrapError(ErrEventQuery, err)
	return s.handleError(err)
}

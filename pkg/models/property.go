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
	"github.com/bdreece/hopper/pkg/proto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Property struct {
	gorm.Model
	proto.Property
	Events []Event
}

func NewProperty(req *proto.CreatePropertyRequest) Property {
	return Property{
		Property: proto.Property{
			Uuid:        uuid.NewString(),
			Name:        req.Name,
			Description: req.Description,
			TypeId:      req.TypeId,
		},
	}
}

func (p *Property) Update(in *proto.UpdatePropertyRequest) {
	if in.Name != nil {
		p.Name = *in.Name
	}
	if in.Description != nil {
		p.Description = in.Description
	}
	if in.TypeId != nil {
		p.TypeId = *in.TypeId
	}
}

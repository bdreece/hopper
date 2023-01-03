package models

import "gorm.io/gorm"

type IEntity interface {
	Marshal() ([]byte, error)
}

type Entity struct {
	gorm.Model
}

type NamedEntity struct {
	Entity

	Name        string
	Description *string
}

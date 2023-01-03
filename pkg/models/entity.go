package models

import "gorm.io/gorm"

type Entity struct {
	gorm.Model
}

type NamedEntity struct {
	Entity

	Name        string
	Description *string
}

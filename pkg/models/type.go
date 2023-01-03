package models

type Type struct {
	NamedEntity

	DataType string
	UnitID   uint

	Properties []Property
}

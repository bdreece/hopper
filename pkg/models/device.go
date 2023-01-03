package models

type Device struct {
	NamedEntity

	Uuid string

	ModelID    uint
	FirmwareID uint

	Properties []Property
	Events     []Event
}

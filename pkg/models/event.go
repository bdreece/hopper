package models

type Event struct {
	Entity

	DeviceID   uint
	PropertyID uint
	Value      string
}

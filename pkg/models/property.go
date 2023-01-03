package models

type Property struct {
	NamedEntity

	PresentValue *string

	TypeID   uint
	DeviceID *uint
	ModelID  *uint

	Events []Event
}

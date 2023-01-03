package models

type Model struct {
	NamedEntity

	Uuid string

	Firmwares  []Firmware
	Devices    []Device
	Properties []Property
}

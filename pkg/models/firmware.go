package models

type Firmware struct {
	Entity

	VersionMajor uint
	VersionMinor uint
	VersionPatch uint
	Url          string

	ModelID uint
	Devices []Device
}

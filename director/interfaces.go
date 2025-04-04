package director

import "github.com/hoopahmadness/sync-director/v2/web"

type Identifier interface {
	GetId() string
	GetFriendlyName() string
}

type Device interface {
	Identifier
	GetOrderedStatus() int // return the integer corresponding to the device's status
	GetConnectedDevices(sd *director) ([]Device, error)
	web.Pairable
}

type NetworkFolder interface {
	Identifier
	GetNumSharedDevices() int
}

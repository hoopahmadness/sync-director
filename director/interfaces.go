package director

type Identifier interface {
	GetId() string
	GetFriendlyName() string
}

// type Device interface {
// 	Identifier
// 	web.Pairable
// 	GetOrderedStatus() int // return the integer corresponding to the device's status
// 	GetConnectedDevices(sd *director) ([]Device, error)
// 	QueryFolders() ()
// }

// type Folder interface {
// }

// type NetworkFolder interface {
// 	Identifier
// 	GetNumSharedDevices() int
// }

// type DeviceWeb interface {
// 	NewDeviceConnection(dev, connectedDev Device)
// }

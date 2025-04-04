package device

import (
	"encoding/json"

	log "github.com/inconshreveable/log15"
)

// A Device contains a client for communicating with a specific instance of Syncthing
type Device struct {
	*Client
	log log.Logger
}

func NewDevice(nickname string, log log.Logger) *Device {
	d := &Device{}
	c := newClient(d, nickname)
	d.Client = c
	d.log = log
	return d
}

func (dev *Device) String() string {
	b, _ := json.Marshal(dev)
	return string(b)
}

// Queries the device client to get a complete list of all synced and pending folders
func (dev *Device) QueryFolders() (GetFolderResponse, GetPendingFoldersResponse, error) {
	syncedFolders, err := dev.Client.querySyncedFolders()
	if err != nil {
		return nil, nil, err
	}
	pendingFolders, err := dev.Client.queryPendingFolders()
	if err != nil {
		return nil, nil, err
	}
	return syncedFolders, pendingFolders, nil
}

// Queries the device client to get all connected devices.
// Also returns pending connections? not sure where I was going with this.
func (dev *Device) GetConnectedDevices(m deviceManager) ([]*Device, error) {
	resp, err := dev.Client.queryConnectedDevices()
	if err != nil {
		// fmt.Println(err.Error())
		return nil, err
	}
	deviceList := []*Device{}
	if resp == nil {
		return deviceList, nil
	}
	for id := range resp.Connections {
		var connectedDevice *Device
		var exists bool
		if connectedDevice, exists = m.GetDeviceById(id); !exists {
			connectedDevice = NewDevice("unknown-"+id[0:7], dev.log)
			connectedDevice.DeviceId = id
			if resp.Connections[id].Address == "" {
				connectedDevice.Status = OFFLINE
			} else {
				connectedDevice.Status = OUTOFNETWORK
				connectedDevice.IpAddress = resp.Connections[id].Address
			}
			// m.DevicesById[id] = connectedDevice
			m.SetDeviceById(connectedDevice)
		}
		// dev.ConnectedDevices[connectedDevice] = struct{ Pending bool }{Pending: false}
		// maybe I'm doing too much in this function and should just return the connected devices
		deviceList = append(deviceList, connectedDevice)
	}
	return deviceList, nil

}

func (dev *Device) GetId() string {
	return dev.DeviceId
}

func (dev *Device) GetFriendlyName() string {
	return dev.Nickname
}

// // Device pairs are bidirectional objects that show a relationship between two devices
// // If one of the devices has offerred to share a folder or connection but a second device has not accepted it,
// // the *offering* device will be put in the pending slot.
// type DevicePair struct {
// 	Dev1         *Device
// 	DevA         *Device
// 	OfferPending *Device
// 	GraphLink    *opts.GraphLink
// }

// func (dp *DevicePair) Other(given *Device) *Device {
// 	if dp.Dev1 == given {
// 		return dp.DevA
// 	} else if dp.DevA == given {
// 		return dp.Dev1
// 	}
// 	return nil
// }

// // If one of the devices has not accepted the folder then this returns the
// // device *offering* the folder for syncing. If both hosts are sharing then returns nil
// func (dp *DevicePair) GetPending() *Device {
// 	return dp.OfferPending
// }

package device

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/inconshreveable/log15"
)

// A Device contains a client for communicating with a specific instance of Syncthing
type Device struct {
	*Client
	log log.Logger
	// lastStatusCheck           time.Time
	lastConnectedDevicesCheck time.Time
	lastFolderCheck           time.Time
}

func NewDevice(nickname string, log log.Logger) *Device {
	d := &Device{}
	c := newClient(d, nickname)
	d.Client = c
	d.log = log
	return d
}

// func (d Device) timeForNewStatusCheck() bool {
// 	now := time.Now()
// 	if d.Status == OFFLINE {
// 		return now.After(d.lastStatusCheck.Add(10 * time.Second))
// 	}
// 	return now.After(d.lastStatusCheck.Add(10 * time.Minute))
// }

func (d Device) timeForNewConnectionsCheck() bool {
	now := time.Now()
	return now.After(d.lastConnectedDevicesCheck.Add(1 * time.Minute))
}

func (d Device) timeForNewFolderCheck() bool {
	now := time.Now()
	return now.After(d.lastFolderCheck.Add(1 * time.Minute))
}

func (dev *Device) AddNetworkInfo(deviceID, apiKey, ipAddress string) {
	dev.addToNetwork(deviceID, apiKey, ipAddress)
}

func (dev *Device) String() string {
	b, _ := json.Marshal(dev)
	return string(b)
}

// Queries the device client to get a complete list of all synced and pending folders
func (dev *Device) QueryFolders() (GetFolderResponse, GetPendingFoldersResponse, error) {
	dev.Client.ping()
	if !dev.timeForNewFolderCheck() {
		return nil, nil, fmt.Errorf("not getting folders for this device again, it's too soon")
	}
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
func (dev Device) GetConnectedDevices(m deviceManager) ([]*Device, error) {
	var resp *ConnectedDevicesResponse
	var err error
	if dev.timeForNewConnectionsCheck() {
		resp, err = dev.Client.queryConnectedDevices()
		if err != nil {
			dev.Status = OFFLINE
			return nil, err
		}
	} else {
		dev.log.Debug("Not going to get connected devices for this device because it's not been long enough",
			"device", dev,
		)
	}
	dev.lastConnectedDevicesCheck = time.Now()
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
		}
		deviceList = append(deviceList, connectedDevice)
	}
	return deviceList, nil

}

func (dev *Device) GetId() string {
	return dev.DeviceId
}

func (dev Device) GetFriendlyName() string {
	return dev.Nickname
}

func (dev Device) GetOrderedStatus() int {
	return orderedStatuses[dev.Status]
}

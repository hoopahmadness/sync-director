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
	Collapsed                 bool
	Hidden                    bool
	SelectedFolder            *bool
}

func NewDevice(log log.Logger) *Device {
	d := &Device{}
	c := newClient(d)
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

func (dev *Device) Offline() bool {
	return dev.Client.Status == OFFLINE
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
func (dev Device) GetConnectedDevices(m deviceManager) ([]*Device, error) {
	var connectedResp *ConnectedDevicesResponse
	var connectedErr error
	var configuredResp ConfiguredDevicesResponse
	var configuredErr error
	if dev.timeForNewConnectionsCheck() {
		connectedResp, connectedErr = dev.Client.queryConnectedDevices()
		if connectedErr != nil {
			dev.Status = OFFLINE
			return nil, connectedErr
		}
		configuredResp, configuredErr = dev.Client.queryConfiguredDevices()
		if configuredErr != nil {
			dev.Status = OFFLINE
			return nil, configuredErr
		}
	} else {
		dev.log.Debug("Not going to get connected devices for this device because it's not been long enough",
			"device", dev,
		)
	}
	dev.lastConnectedDevicesCheck = time.Now()
	deviceList := []*Device{}
	if connectedResp == nil {
		return deviceList, nil
	}
	for _, confDev := range configuredResp {
		id := confDev.DeviceID
		devName := confDev.Name
		var connectedDevice *Device
		var exists bool
		if connectedDevice, exists = m.GetDeviceById(id); !exists {
			connectedDevice = NewDevice(dev.log)
			connectedDevice.DeviceId = id
			if connectedResp.Connections[id].Address == "" {
				connectedDevice.Status = OFFLINE
			} else {
				connectedDevice.Status = OUTOFNETWORK
				connectedDevice.IpAddress = connectedResp.Connections[id].Address
			}
		}
		connectedDevice.AddNickname(devName)
		deviceList = append(deviceList, connectedDevice)
	}
	return deviceList, nil
}

func (dev *Device) GetId() string {
	return dev.DeviceId
}

func (dev Device) Name() string {
	if dev.Nickname == "" {
		return "unknown-" + dev.DeviceId[0:7]
	}
	return dev.Nickname
}

func (dev Device) GetOrderedStatus() int {
	return orderedStatuses[dev.Status]
}

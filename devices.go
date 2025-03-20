package main

import (
	"encoding/json"

	"github.com/go-echarts/go-echarts/v2/opts"
	log "github.com/inconshreveable/log15"
)

// A Device contains a client for communicating with a specific instance of Syncthing
type Device struct {
	*Client
	GraphNode *opts.GraphNode
	log       log.Logger
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
func (dev *Device) GetFolders() ([]*Folder, error) {
	syncedFolders, err := dev.Client.querySyncedFolders()
	if err != nil {
		return nil, err
	}
	pendingFolders, err := dev.Client.queryPendingFolders()
	if err != nil {
		return nil, err
	}

	// parse normal folders into a map
	idMap := map[string]*Folder{}

	for _, response := range syncedFolders {
		folder := &Folder{
			Id:            response.Id,
			Label:         response.Label,
			Path:          response.Path,
			Type:          response.Type,
			HostDevice:    dev.Client.DeviceId,
			SharedDevices: map[string]struct{ Pending bool }{},
		}
		for _, data := range response.Devices {
			folder.SharedDevices[data.DeviceID] = struct{ Pending bool }{Pending: false}
		}
		idMap[folder.Id] = folder
	}

	// parse pending folders into the map, updating existing folders as necesary
	for folderId, offeredBy := range pendingFolders {
		existingFolder, OK := idMap[folderId]
		if !OK {
			existingFolder = &Folder{
				Id:         folderId,
				Label:      "",
				Path:       "",
				Type:       "",
				HostDevice: dev.DeviceId,
				SharedDevices: map[string]struct {
					Pending bool
				}{},
			}
		}
		for offeringDeviceId, data := range offeredBy.OfferedBy {
			existingFolder.SharedDevices[offeringDeviceId] = struct{ Pending bool }{
				Pending: true,
			}
			existingFolder.Label = data.Label
		}
	}

	// push them all back into an array
	asArray := []*Folder{}
	for _, folder := range idMap {
		asArray = append(asArray, folder)
	}

	return asArray, nil

}

// Queries the device client to get all connected devices.
// Also returns pending connections? not sure where I was going with this.
func (dev *Device) GetConnectedDevices(m *model) ([]*Device, error) {
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
		if connectedDevice, exists = m.DevicesById[id]; !exists {
			connectedDevice = NewDevice("unknown-"+id[0:7], m.log)
			connectedDevice.DeviceId = id
			connectedDevice.Status = OUTOFNETWORK
			m.DevicesById[id] = connectedDevice
		}
		// dev.ConnectedDevices[connectedDevice] = struct{ Pending bool }{Pending: false}
		// maybe I'm doing too much in this function and should just return the connected devices
		deviceList = append(deviceList, connectedDevice)
	}
	return deviceList, nil

}

// Device pairs are bidirectional objects that show a relationship between two devices
// If one of the devices has offerred to share a folder or connection but a second device has not accepted it,
// the *offering* device will be put in the pending slot.
type DevicePair struct {
	Dev1         *Device
	DevA         *Device
	OfferPending *Device
	GraphLink    *opts.GraphLink
}

func (dp *DevicePair) Other(given *Device) *Device {
	if dp.Dev1 == given {
		return dp.DevA
	} else if dp.DevA == given {
		return dp.Dev1
	}
	return nil
}

// If one of the devices has not accepted the folder then this returns the
// device *offering* the folder for syncing. If both hosts are sharing then returns nil
func (dp *DevicePair) GetPending() *Device {
	return dp.OfferPending
}

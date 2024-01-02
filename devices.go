package main

import "fmt"

type Device struct {
	*Client
}

func (dev *Device) GetFolders() ([]*Folder, error) {
	fmt.Println("Getting folders")
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
			HostDevice:    dev.Client.deviceId,
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
				HostDevice: dev.deviceId,
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

func (dev *Device) GetConnectedDevices() ([]*Device, error) {

}

// Device pairs are bidirectional objects that show a relationship between two devices
// If one of the devices has offerred to share a folder but a second device has not accepted it,
// the *offering* device will be put in the pending slot.
type DevicePair struct {
	dev1         *Device
	devA         *Device
	offerPending *Device
}

func (dp *DevicePair) Other(given *Device) *Device {
	if dp.dev1 == given {
		return dp.devA
	} else if dp.devA == given {
		return dp.dev1
	}
	return nil
}

// If one of the devices has not accepted the folder then this returns the
// device *offering* the folder for syncing. If both hosts are sharing then returns nil
func (dp *DevicePair) GetPending() *Device {
	return dp.offerPending
}

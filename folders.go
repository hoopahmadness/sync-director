package main

import "encoding/json"

// This refers to the network-wide concept of a Folder that is shared on many devices.
// The NetworkFolder tracks the specific Folder instances corresponding to various Devices
// It also tracks the web of conections between devices that this folder is synced on.
type NetworkFolder struct {
	Id        string
	folders   map[*Device]*Folder
	DeviceWeb *DeviceWeb
}

func (nf *NetworkFolder) String() string {
	b, _ := json.Marshal(nf)
	return string(b)
}

func newNetworkFolder(id string) *NetworkFolder {
	nf := &NetworkFolder{
		Id:      id,
		folders: map[*Device]*Folder{},
	}
	nf.DeviceWeb = newDeviceWeb()
	return nf
}

func (nf *NetworkFolder) IngestFolder(folder *Folder) {
	// get the device for this folder
	hostDev := devicesById[folder.HostDevice]

	// get all the device pairs from this folder
	for sharedDevID, data := range folder.SharedDevices {
		sharedDev := devicesById[sharedDevID]
		nf.DeviceWeb.NewDevicePairForFolder(hostDev, sharedDev, data.Pending)
	}

	// add folder to map
	nf.folders[hostDev] = folder
}

// This refers to the knowledge a specific Device has about the folders it is watching, as well as
// pending folders that it has not saved. Many Folders can have the same ID since many devices have
// their own instances of that folder
// The Host Device obvious refers to the device that this is saved on, with a path, etc.
// The shared devices are the devices that are syncing this folder or *have offerred* to sync
// with the host device
type Folder struct {
	Id            string
	Label         string
	Path          string
	Type          string
	HostDevice    string
	SharedDevices map[string]struct {
		Pending bool
	}
}

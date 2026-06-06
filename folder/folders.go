package folder

import (
	"encoding/json"

	"github.com/hoopahmadness/sync-director/v2/device"
	"github.com/hoopahmadness/sync-director/v2/web"
	log "github.com/inconshreveable/log15"
)

// This refers to the network-wide concept of a Folder that is shared on many devices.
// The NetworkFolder tracks the specific Folder instances corresponding to various Devices
// It also tracks the web of conections between devices that this folder is synced on.
type NetworkFolder[D any, DPtr web.Pairable[D]] struct {
	Id        string
	Folders   map[DPtr]*Folder
	DeviceWeb *web.Web[D, DPtr]
}

func (nf *NetworkFolder[D, Device]) String() string {
	b, _ := json.Marshal(nf)
	return string(b)
}

func (nf *NetworkFolder[D, Device]) Name() string {
	return nf.Id
}

func (nf *NetworkFolder[D, Device]) GetId() string {
	return nf.Id
}

func NewNetworkFolder[D any, DPtr web.Pairable[D]](id string, log log.Logger) *NetworkFolder[D, DPtr] {
	nf := &NetworkFolder[D, DPtr]{
		Id:      id,
		Folders: map[DPtr]*Folder{},
	}
	nf.DeviceWeb = web.NewDeviceWeb[D, DPtr](log)
	return nf
}

func (nf *NetworkFolder[D, Device]) IngestFolder(folder *Folder, m folderManager[D, Device]) {
	// get the device for this folder
	hostDev, _ := m.GetDeviceById(folder.HostDevice)

	// get all the device pairs from this folder
	for sharedDevID, data := range folder.SharedDevices {
		sharedDev, _ := m.GetDeviceById(sharedDevID)
		nf.DeviceWeb.NewDevicePairForFolder(hostDev, sharedDev, data.Pending, nil)
	}

	// add folder to map
	nf.Folders[hostDev] = folder
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

func GetFolders(dev *device.Device, logger *log.Logger) ([]*Folder, error) {
	syncedFolders, pendingFolders, err := dev.QueryFolders()
	if err != nil {
		log.Info("Unable to query folders for this device", "err", err, "device", &dev)
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

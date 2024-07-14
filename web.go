package main

// A device web represents a set of links between pairs of devices regarding a given folder
type DeviceWeb struct {
	allPairs  map[*DevicePair]bool
	perDevice map[*Device][]*DevicePair
	// manager   *WebManager
}

// Given a device, find the other devices that are sharing the folder with it
// Does not care if a device is pending or not.
func (dw *DeviceWeb) getOtherDevices(aDevice *Device) []*Device {
	dw.initializeDevicePairsListIfNeeded(aDevice)
	pairs := dw.perDevice[aDevice]
	devices := []*Device{}
	for _, pair := range pairs {
		devices = append(devices, pair.Other(aDevice))
	}
	return devices
}

// A folder can add a new pairing between two of its devices.
// In general pairs are bidirectional so it doesn't matter which device is the host;
// adding the inverse pair will be compared to the existing pair and is a no-op.
// If the host device has not accepted syncing for a folder from the synced device, set 'pending'
// to true and the device pair will show that there is a pending folder.
func (dw *DeviceWeb) NewDevicePairForFolder(hostDevice, syncedDevice *Device, pending bool) {
	dp := dw.getDevicePair(hostDevice, syncedDevice)
	if pending {
		dp.offerPending = syncedDevice
	}
}

// Create a new connection between devices. All device pairs created this way are pending by defualt.
// When the inverse pair is created (in other words, when the synced device tries to add a connection
// to this device) then the connection will no longer be pending
func (dw *DeviceWeb) NewDeviceConnection(hostDevice, connectedDevice *Device) {
	dp := dw.getDevicePair(hostDevice, connectedDevice)
	if dp.offerPending == nil {
		dp.offerPending = hostDevice
	} else if dp.offerPending == connectedDevice {
		dp.offerPending = nil
	}
}

// Adds a new DevicePair to the web.
// Each DevicePair is kept in a map as well as kept in a list for each related device
// (that's three pointers per pair)
// Adding an existing pair again is a no-op
// For now let's assume this is only run internall by the web manager
func (dw *DeviceWeb) addPairing(pair *DevicePair) {
	if _, OK := dw.allPairs[pair]; OK {
		return
	}
	aDevice := pair.dev1
	anotherDevice := pair.devA
	dw.initializeDevicePairsListIfNeeded(aDevice, anotherDevice)
	listforADevice := dw.perDevice[aDevice]
	listforADevice = append(listforADevice, pair)
	dw.perDevice[aDevice] = listforADevice

	listforAnotherDevice := dw.perDevice[anotherDevice]
	listforAnotherDevice = append(listforAnotherDevice, pair)
	dw.perDevice[anotherDevice] = listforAnotherDevice

	dw.allPairs[pair] = true
}

func (dw *DeviceWeb) initializeDevicePairsListIfNeeded(devices ...*Device) {
	for _, aDevice := range devices {
		_, OK := dw.perDevice[aDevice]
		if !OK {
			dw.perDevice[aDevice] = []*DevicePair{}
		}
	}
}

// // The manager will be a global variable that tracks and disseminates device pairs to sub-webs on a per-folder basis
// // WebManager actually uses a "master" DeviceWeb internally to store all known DevicePair objects
// // WAIT: Why am I globalizing these? Wouldn't the pending field correspond to specific folders, since a device can be pending for one folder but not another?
// type WebManager struct {
// 	master     *DeviceWeb
// 	folderWebs map[*Folder]*DeviceWeb
// }

// This func should be used by a DeviceWeb to get a pointer to an existing DevicePair
// If the DevicePair for these devices can't be found it is created.
// This allows us to guarantee that any pair of devices regardless of order will yield the same object
func (dw *DeviceWeb) getDevicePair(aDevice, anotherDevice *Device) *DevicePair {

	// see if the pair already exists
	pairList := dw.perDevice[aDevice]
	for _, pair := range pairList {
		if pair.Other(aDevice) == anotherDevice {
			return pair
		}
	}

	// create it
	dp := &DevicePair{aDevice, anotherDevice, nil}

	dw.addPairing(dp)

	return dp
}

func newDeviceWeb() *DeviceWeb {
	newWeb := &DeviceWeb{}
	newWeb.perDevice = map[*Device][]*DevicePair{}
	newWeb.allPairs = map[*DevicePair]bool{}
	return newWeb
}

// func (wb *WebManager) ProcessNewFolder(folder *Folder) {
// 	newWeb := wb.NewDeviceWeb()
// 	if folder != nil {
// 		wb.folderWebs[folder] = newWeb
// 	}
// 	// don't we need to do more here? I think there should be some actual processing but I'm not sure off the top of my head what it is
// }

// func newWebManager() *WebManager {
// 	newMan := &WebManager{
// 		folderWebs: map[*Folder]*DeviceWeb{},
// 	}
// 	newMan.master = newMan.NewDeviceWeb()
// 	return newMan
// }

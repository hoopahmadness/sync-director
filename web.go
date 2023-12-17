package main

// A device web represents a set of links between pairs of devices.
type DeviceWeb struct {
	// allPairs  []*DevicePair
	perDevice map[*Device][]*DevicePair
	*WebManager
}

func (dw *DeviceWeb) getOtherDevices(aDevice *Device) []*Device {
	dw.initializeDevicePairsListIfNeeded(aDevice)
	pairs := dw.perDevice[aDevice]
	devices := []*Device{}
	for _, pair := range pairs {
		devices = append(devices, pair.Other(aDevice))
	}
	return devices
}

func (dw *DeviceWeb) getPairing(aDevice, anotherDevice *Device) *DevicePair {
	pairList := dw.perDevice[aDevice]
	for _, pair := range pairList {
		if pair.Other(aDevice) == anotherDevice {
			return pair
		}
	}
	return nil
}

func (dw *DeviceWeb) NewDevicePairForFolder(aDevice, anotherDevice *Device, folder *Folder) {
	dp := dw.master.getDevicePair(aDevice, anotherDevice)
	dw.addPairing(dp)
}
func (dw *DeviceWeb) addPairing(pair *DevicePair) {
	aDevice := pair.dev1
	anotherDevice := pair.devA
	dw.initializeDevicePairsListIfNeeded(aDevice, anotherDevice)
	listforADevice := dw.perDevice[aDevice]
	listforADevice = append(listforADevice, pair)
	dw.perDevice[aDevice] = listforADevice

	listforAnotherDevice := dw.perDevice[anotherDevice]
	listforAnotherDevice = append(listforAnotherDevice, pair)
	dw.perDevice[anotherDevice] = listforAnotherDevice

}

func (dw *DeviceWeb) initializeDevicePairsListIfNeeded(devices ...*Device) {
	for _, aDevice := range devices {
		_, OK := dw.perDevice[aDevice]
		if !OK {
			dw.perDevice[aDevice] = []*DevicePair{}
		}
	}
}

// The manager will be a global variable that tracks and disseminates device pairs to sub-webs on a per-folder basis
type WebManager struct {
	master *DeviceWeb
	// folderWebs map[*Folder]DeviceWeb
}

func (wm *WebManager) getDevicePair(aDevice, anotherDevice *Device) *DevicePair {

	// see if the pair already exists
	pair := wm.master.getPairing(aDevice, anotherDevice)
	if pair != nil {
		return pair
	}

	// create it
	dp := &DevicePair{aDevice, anotherDevice}
	// add it to the master
	wm.master.addPairing(dp)
	return dp
}

func (wb *WebManager) newDeviceWeb(folder *Folder) *DeviceWeb {
	newWeb := &DeviceWeb{}
	// newWeb.allPairs = []*DevicePair{}
	newWeb.perDevice = map[*Device][]*DevicePair{}
	newWeb.WebManager = wb
	// if nil then don't add to this map
	// wb.folderWebs[folder] = newWeb
	return newWeb
}

func newWebManager() *WebManager {
	newMan := &WebManager{}
	newMan.master = newMan.newDeviceWeb(nil)
	return newMan
}

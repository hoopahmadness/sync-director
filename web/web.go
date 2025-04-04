package web

import (
	"fmt"

	log "github.com/inconshreveable/log15"
)

// A device web represents a set of links between pairs of devices regarding a given folder
type Web[Device Pairable] struct {
	AllPairs  map[*Pair[Device]]bool
	PerDevice map[*Device][]*Pair[Device]
	log       log.Logger
}

func (dw *Web[Device]) String() string {
	webStr := "Device web:\n"
	for pairing, _ := range dw.AllPairs {
		firstOffering := ""
		secondOffering := ""
		if pairing.OfferPending == pairing.Dev1 {
			firstOffering = "*"
		}
		if pairing.OfferPending == pairing.DevA {
			secondOffering = "*"
		}
		x := *(pairing.Dev1)
		y := *(pairing.DevA)
		pairStr := fmt.Sprintf("%s%s paired with %s%s\n", x.FriendlyName(), firstOffering, y.FriendlyName(), secondOffering)
		webStr += pairStr
	}
	return webStr
}

// Given a device, find the other devices that are sharing the folder with it
// Does not care if a device is pending or not.
func (dw *Web[Device]) getOtherDevices(aDevice *Device) []*Device {
	dw.initializeDevicePairsListIfNeeded(aDevice)
	pairs := dw.PerDevice[aDevice]
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
func (dw *Web[Device]) NewDevicePairForFolder(hostDevice, syncedDevice *Device, pending bool) {
	dp, err := dw.getDevicePair(hostDevice, syncedDevice)
	if err != nil {
		// fmt.Println("A device can't share a folder with itself; skipping")
		return
	}
	if pending {
		dp.OfferPending = syncedDevice
	}
}

// Create a new connection between devices. All device pairs created this way are pending by defualt.
// When the inverse pair is created (in other words, when the synced device tries to add a connection
// to this device) then the connection will no longer be pending
func (dw *Web[Device]) NewDeviceConnection(hostDevice, connectedDevice *Device) {
	dp, err := dw.getDevicePair(hostDevice, connectedDevice)
	if err != nil {
		// fmt.Println("Can't pair device to itself; skipping")
		return
	}
	if dp.OfferPending == nil {
		dp.OfferPending = hostDevice
	} else if dp.OfferPending == connectedDevice {
		// fmt.Println("Accepting pairing offer")
		dp.OfferPending = nil
	}
}

// Adds a new DevicePair to the web.
// Each DevicePair is kept in a map as well as kept in a list for each related device
// (that's three pointers per pair)
// Adding an existing pair again is a no-op
// For now let's assume this is only run internall by the web manager
func (dw *Web[Device]) addPairing(pair *Pair[Device]) {
	if _, OK := dw.AllPairs[pair]; OK {
		return
	}
	aDevice := pair.Dev1
	anotherDevice := pair.DevA
	dw.initializeDevicePairsListIfNeeded(aDevice, anotherDevice)
	listforADevice := dw.PerDevice[aDevice]
	listforADevice = append(listforADevice, pair)
	dw.PerDevice[aDevice] = listforADevice

	listforAnotherDevice := dw.PerDevice[anotherDevice]
	listforAnotherDevice = append(listforAnotherDevice, pair)
	dw.PerDevice[anotherDevice] = listforAnotherDevice

	dw.AllPairs[pair] = true
}

func (dw *Web[Device]) initializeDevicePairsListIfNeeded(devices ...*Device) {
	dw.log.Debug("logging the contents of this DW!", "dw", dw.String())
	for _, aDevice := range devices {
		_, OK := dw.PerDevice[aDevice]
		if !OK {
			dw.PerDevice[aDevice] = []*Pair[Device]{}
		}
	}
}

// This func should be used by a DeviceWeb to get a pointer to an existing DevicePair
// If the DevicePair for these devices can't be found it is created.
// This allows us to guarantee that any pair of devices regardless of order will yield the same object
// Returns error when attempting to pair a device to itself
func (dw *Web[Device]) getDevicePair(aDevice, anotherDevice *Device) (*Pair[Device], error) {
	// make sure devices are not the same
	if aDevice == anotherDevice {
		return nil, fmt.Errorf("cannot pair two devices that are the same")
	}

	// see if the pair already exists
	pairList := dw.PerDevice[aDevice]
	for _, pair := range pairList {
		if pair.Other(aDevice) == anotherDevice {
			// fmt.Println("Found existing device pairing: " + aDevice.Nickname + " & " + anotherDevice.Nickname)
			return pair, nil
		}
	}

	// create it
	// fmt.Println("Creating new device pairing: " + aDevice.Nickname + " & " + anotherDevice.Nickname)
	dp := &Pair[Device]{
		Dev1:         aDevice,
		DevA:         anotherDevice,
		OfferPending: nil,
	}

	dw.addPairing(dp)

	return dp, nil
}

func NewDeviceWeb[Device Pairable](log log.Logger) *Web[Device] {
	newWeb := &Web[Device]{}
	newWeb.PerDevice = map[*Device][]*Pair[Device]{}
	newWeb.AllPairs = map[*Pair[Device]]bool{}
	newWeb.log = log
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
